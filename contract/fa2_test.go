package contract

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/tezos"
)

func TestFA2Approval_UnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"add_operator": {
			"owner": "tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb",
			"operator": "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6",
			"token_id": "0"
		}
	}`)

	var a FA2Approval
	require.NoError(t, json.Unmarshal(data, &a))
	require.True(t, a.Add)
	require.True(t, a.Owner.Equal(tezos.MustParseAddress("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb")))
	require.True(t, a.Operator.Equal(tezos.MustParseAddress("tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6")))
	require.Equal(t, int64(0), a.TokenId.Int64())
}

func TestFA2ApprovalArgs_ParametersBranches(t *testing.T) {
	owner := tezos.MustParseAddress("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb")
	op := tezos.MustParseAddress("tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6")
	id := tezos.NewZ(7)

	args := NewFA2ApprovalArgs().
		AddOperator(owner, op, id).
		RemoveOperator(owner, op, id)

	p := args.Parameters()
	require.Equal(t, "update_operators", p.Entrypoint)
	require.Equal(t, micheline.PrimSequence, p.Value.Type)
	require.Len(t, p.Value.Args, 2)

	require.Equal(t, micheline.D_LEFT, p.Value.Args[0].OpCode)  // add_operator
	require.Equal(t, micheline.D_RIGHT, p.Value.Args[1].OpCode) // remove_operator
}

func TestFA2TransferArgs_OptimizeAndParameters(t *testing.T) {
	from1 := tezos.MustParseAddress("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb")
	from2 := tezos.MustParseAddress("tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6")
	to := tezos.MustParseAddress("tz1burnburnburnburnburnburnburjAYjjX")

	args := NewFA2TransferArgs().
		WithTransfer(from2, to, tezos.NewZ(0), tezos.NewZ(1)).
		WithTransfer(from1, to, tezos.NewZ(0), tezos.NewZ(2)).
		WithTransfer(from2, to, tezos.NewZ(1), tezos.NewZ(3)).
		Optimize()

	p := args.Parameters()
	require.Equal(t, "transfer", p.Entrypoint)
	require.Equal(t, micheline.PrimSequence, p.Value.Type)

	// grouped by from address after Optimize()
	require.Len(t, p.Value.Args, 2)
	require.Equal(t, micheline.D_PAIR, p.Value.Args[0].OpCode)
	require.Equal(t, micheline.D_PAIR, p.Value.Args[1].OpCode)

	// each group is Pair(from, Seq(...txs...))
	for _, grp := range p.Value.Args {
		require.Len(t, grp.Args, 2)
		require.Equal(t, micheline.PrimBytes, grp.Args[0].Type)
		require.Equal(t, micheline.PrimSequence, grp.Args[1].Type)
	}
}
