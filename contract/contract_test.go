package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

func orRight(prims []micheline.Prim) micheline.Prim {
	switch len(prims) {
	case 0:
		return micheline.Prim{}
	case 1:
		return prims[0]
	default:
		return micheline.NewCode(micheline.T_OR, prims[0], orRight(prims[1:]))
	}
}

func scriptForInterface(i micheline.Interface) *micheline.Script {
	specs := micheline.InterfaceSpecs[i]
	prims := make([]micheline.Prim, 0, len(specs))
	for _, v := range specs {
		prims = append(prims, v)
	}
	paramType := orRight(prims)
	return &micheline.Script{
		Code: micheline.Code{
			Param:   micheline.NewCode(micheline.K_PARAMETER, paramType),
			Storage: micheline.NewCode(micheline.K_STORAGE, micheline.NewCode(micheline.T_UNIT)),
			Code:    micheline.NewCode(micheline.K_CODE, micheline.NewSeq()),
		},
		Storage: micheline.NewCode(micheline.D_UNIT),
	}
}

func TestTxArgs_Encode(t *testing.T) {
	src := tezos.MustParseAddress("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb")
	dst := tezos.MustParseAddress("KT1RJ6PbjHpwc3M5rw5s2Nbmefwbuwbdxton")
	amt := tezos.N(123)
	params := micheline.Parameters{Entrypoint: "default", Value: micheline.NewCode(micheline.D_UNIT)}

	a := NewTxArgs()
	a.WithSource(src)
	a.WithDestination(dst)
	a.WithAmount(amt)
	a.WithParameters(params)

	tx := a.Encode()
	require.True(t, tx.Source.Equal(src))
	require.True(t, tx.Destination.Equal(dst))
	require.Equal(t, int64(amt), int64(tx.Amount))
	require.NotNil(t, tx.Parameters)
	require.Equal(t, "default", tx.Parameters.Entrypoint)
}

func TestContract_TokenKind(t *testing.T) {
	fa2code, err := os.ReadFile("../examples/tzcompose/token/fa2_single_asset.json")
	require.NoError(t, err)
	var fa2 micheline.Code
	require.NoError(t, json.Unmarshal(fa2code, &fa2))

	c := NewEmptyContract(nil).WithScript(&micheline.Script{
		Code:    fa2,
		Storage: micheline.NewCode(micheline.D_UNIT),
	})
	require.True(t, c.IsFA2())
	require.True(t, c.IsToken())
	require.Equal(t, TokenKindFA2, c.TokenKind())

	fa12code, err := os.ReadFile("../examples/tzcompose/token/fa12_code.json")
	require.NoError(t, err)
	var fa12 micheline.Code
	require.NoError(t, json.Unmarshal(fa12code, &fa12))

	c = NewEmptyContract(nil).WithScript(&micheline.Script{
		Code:    fa12,
		Storage: micheline.NewCode(micheline.D_UNIT),
	})
	require.True(t, c.IsFA12())
	// Note: many FA1.2 contracts also satisfy the older FA1 interface, and
	// TokenKind() currently classifies those as TokenKindFA1.
	require.True(t, c.IsFA1())
	require.Equal(t, TokenKindFA1, c.TokenKind())
}

func TestContract_EntrypointAndView_NoScript(t *testing.T) {
	c := NewEmptyContract(nil)
	_, ok := c.Entrypoint("transfer")
	require.False(t, ok)

	_, ok = c.View("x")
	require.False(t, ok)
}

func TestContract_RPCBacked_ResolveReloadRunViewRunCallback(t *testing.T) {
	addr := tezos.MustParseAddress("KT1RJ6PbjHpwc3M5rw5s2Nbmefwbuwbdxton")

	script := scriptForInterface(micheline.ITzip12)
	scriptResp, err := json.Marshal(script)
	require.NoError(t, err)

	var storageBytes = []byte(`{"int":"1"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/chains/main/blocks/head/context/contracts/"+addr.String()+"/script/normalized":
			// ensure the request body contains the unparsing mode we expect
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(r.Body)
			require.True(t, strings.Contains(buf.String(), "Readable"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(scriptResp)
			return

		case r.Method == http.MethodGet && r.URL.Path == "/chains/main/blocks/head/context/contracts/"+addr.String()+"/storage":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(storageBytes)
			return

		case r.Method == http.MethodPost && r.URL.Path == "/chains/main/blocks/head/helpers/scripts/run_script_view":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"int":"2"}}`))
			return

		case r.Method == http.MethodPost && r.URL.Path == "/chains/main/blocks/head/helpers/scripts/run_view":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"int":"3"}}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cli, err := rpc.NewClient(server.URL, nil)
	require.NoError(t, err)

	c := NewContract(addr, cli)
	require.NoError(t, c.Resolve(context.Background()))
	require.NotNil(t, c.Script())
	require.NotNil(t, c.Storage())

	storageBytes = []byte(`{"int":"5"}`)
	require.NoError(t, c.Reload(context.Background()))
	require.Equal(t, int64(5), c.Storage().Int.Int64())

	prim, err := c.RunView(context.Background(), "my_view", micheline.NewCode(micheline.D_UNIT))
	require.NoError(t, err)
	require.Equal(t, int64(2), prim.Int.Int64())

	prim, err = c.RunCallback(context.Background(), "my_callback", micheline.NewCode(micheline.D_UNIT))
	require.NoError(t, err)
	require.Equal(t, int64(3), prim.Int.Int64())
}
