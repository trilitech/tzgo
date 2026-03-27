package contract

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

func TestTz16License_UnmarshalJSON(t *testing.T) {
	var l Tz16License

	require.NoError(t, json.Unmarshal([]byte(`"MIT"`), &l))
	require.Equal(t, "MIT", l.Name)

	require.NoError(t, json.Unmarshal([]byte(`{"name":"ISC","details":"x"}`), &l))
	require.Equal(t, "ISC", l.Name)
	require.Equal(t, "x", l.Details)
}

func TestResolveTz16Uri_TezosStorageAndSha256(t *testing.T) {
	addr := tezos.MustParseAddress("KT1RJ6PbjHpwc3M5rw5s2Nbmefwbuwbdxton")

	// storage type: big_map %metadata string bytes
	storageType := micheline.NewCodeAnno(
		micheline.T_BIG_MAP,
		"%metadata",
		micheline.NewCode(micheline.T_STRING),
		micheline.NewCode(micheline.T_BYTES),
	)
	// storage value: bigmap id (int)
	storageValue := micheline.NewInt64(123)

	// RPC bigmap payload: {"name":"MyContract"}
	metaJSON := []byte(`{"name":"MyContract"}`)
	metaPrim := micheline.NewBytes(metaJSON)
	metaPrimJSON, err := json.Marshal(metaPrim)
	require.NoError(t, err)

	// Run code response: storage = Some(5)
	runCodeResp := rpc.RunCodeResponse{
		Storage: micheline.NewCode(micheline.D_SOME, micheline.NewInt64(5)),
	}
	runCodeRespJSON, err := json.Marshal(runCodeResp)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/chains/main/blocks/head/context/big_maps/123/"):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(metaPrimJSON)
			return

		case r.Method == http.MethodPost && r.URL.Path == "/chains/main/blocks/head/helpers/scripts/run_code":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(runCodeRespJSON)
			return

		case r.Method == http.MethodGet && r.URL.Path == "/meta.json":
			// ensure ResolveTz16Uri sets expected headers
			require.Contains(t, r.Header.Get("Accept"), "text/plain")
			require.NotEmpty(t, r.Header.Get("User-Agent"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(metaJSON)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cli, err := rpc.NewClient(server.URL, nil)
	require.NoError(t, err)

	c := NewContract(addr, cli).
		WithScript(&micheline.Script{
			Code: micheline.Code{
				Param:   micheline.NewCode(micheline.K_PARAMETER, micheline.NewCode(micheline.T_UNIT)),
				Storage: micheline.NewCode(micheline.K_STORAGE, storageType),
				Code:    micheline.NewCode(micheline.K_CODE, micheline.NewSeq()),
			},
			Storage: storageValue,
		}).
		WithStorage(&storageValue)

	// tezos-storage:... resolves via bigmap lookup and JSON unmarshal
	var out Tz16
	require.NoError(t, c.ResolveTz16Uri(context.Background(), "tezos-storage:contents", &out, nil))
	require.Equal(t, "MyContract", out.Name)

	// sha256:// wrapper should validate checksum
	sum := sha256.Sum256(metaJSON)
	escaped := url.QueryEscape("tezos-storage:contents")
	uri := fmt.Sprintf("sha256://0x%s/%s", hex.EncodeToString(sum[:]), escaped)
	out = Tz16{}
	require.NoError(t, c.ResolveTz16Uri(context.Background(), uri, &out, nil))
	require.Equal(t, "MyContract", out.Name)

	// wrong checksum should error
	bad := bytes.Repeat([]byte{0x00}, 32)
	uri = fmt.Sprintf("sha256://0x%s/%s", hex.EncodeToString(bad), escaped)
	require.ErrorContains(t, c.ResolveTz16Uri(context.Background(), uri, &out, nil), "checksum mismatch")

	// http(s) URL should work (no network) via httptest handler
	out = Tz16{}
	require.NoError(t, c.ResolveTz16Uri(context.Background(), server.URL+"/meta.json", &out, nil))
	require.Equal(t, "MyContract", out.Name)

	// TZIP-16 storage view Run() returns resp.Storage.Args[0]
	view := &Tz16StorageView{
		// ParamType invalid => should be treated as unit
		ReturnType: micheline.NewCode(micheline.T_INT),
		Code:       micheline.NewSeq(micheline.NewCode(micheline.I_DROP), micheline.NewInt64(5)),
	}
	res, err := view.Run(context.Background(), c, micheline.InvalidPrim)
	require.NoError(t, err)
	require.Equal(t, int64(5), res.Int.Int64())
}
