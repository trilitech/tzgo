package contract

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/tezos"
)

func TestFA1Transfer_UnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"transfer": {
			"from": "tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb",
			"to": "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6",
			"value": "42"
		}
	}`)

	var x FA1Transfer
	require.NoError(t, json.Unmarshal(data, &x))
	require.True(t, x.From.Equal(tezos.MustParseAddress("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb")))
	require.True(t, x.To.Equal(tezos.MustParseAddress("tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6")))
	require.Equal(t, int64(42), x.Amount.Int64())
}

func TestFA1ApprovalArgs_Parameters(t *testing.T) {
	spender := tezos.MustParseAddress("tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6")
	amt := tezos.NewZ(10)

	a := NewFA1ApprovalArgs().Approve(spender, amt)
	p := a.Parameters()
	require.Equal(t, "approve", p.Entrypoint)
	require.Equal(t, micheline.D_PAIR, p.Value.OpCode)
	require.Equal(t, micheline.PrimBytes, p.Value.Args[0].Type)
	require.Equal(t, spender.EncodePadded(), p.Value.Args[0].Bytes)
	require.Equal(t, amt.Big(), p.Value.Args[1].Int)

	a = NewFA1ApprovalArgs().Revoke(spender)
	p = a.Parameters()
	require.Equal(t, int64(0), p.Value.Args[1].Int.Int64())
}

func TestFA1TransferArgs_Parameters(t *testing.T) {
	from := tezos.MustParseAddress("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb")
	to := tezos.MustParseAddress("tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6")
	amt := tezos.NewZ(123)

	a := NewFA1TransferArgs().WithTransfer(from, to, amt)
	p := a.Parameters()
	require.Equal(t, "transfer", p.Entrypoint)
	require.Equal(t, micheline.D_PAIR, p.Value.OpCode)
	require.Equal(t, from.EncodePadded(), p.Value.Args[0].Bytes)
	require.Equal(t, to.EncodePadded(), p.Value.Args[1].Args[0].Bytes)
	require.Equal(t, amt.Big(), p.Value.Args[1].Args[1].Int)
}
