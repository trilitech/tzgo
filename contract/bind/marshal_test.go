package bind

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/tezos"
)

func TestMarshalPrim_Scalars(t *testing.T) {
	t.Run("time", func(t *testing.T) {
		ts := time.Unix(1, 0).UTC()

		p, err := MarshalPrim(ts, false)
		require.NoError(t, err)
		require.Equal(t, micheline.PrimString, p.Type)
		require.Equal(t, "1970-01-01T00:00:01Z", p.String)

		p, err = MarshalPrim(ts, true)
		require.NoError(t, err)
		require.Equal(t, micheline.PrimInt, p.Type)
		require.Equal(t, big.NewInt(1), p.Int)
	})

	t.Run("address", func(t *testing.T) {
		addr := tezos.MustParseAddress("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb")

		p, err := MarshalPrim(addr, false)
		require.NoError(t, err)
		require.Equal(t, micheline.PrimString, p.Type)
		require.Equal(t, addr.String(), p.String)

		p, err = MarshalPrim(addr, true)
		require.NoError(t, err)
		require.Equal(t, micheline.PrimBytes, p.Type)
		require.Equal(t, addr.EncodePadded(), p.Bytes)
	})
}

func TestMarshalParams_FoldRightComb(t *testing.T) {
	p, err := MarshalParams(false, "a", "b", "c")
	require.NoError(t, err)

	// right comb: Pair("a", Pair("b","c"))
	require.Equal(t, micheline.D_PAIR, p.OpCode)
	require.Len(t, p.Args, 2)
	require.Equal(t, "a", p.Args[0].String)
	require.Equal(t, micheline.D_PAIR, p.Args[1].OpCode)
	require.Equal(t, "b", p.Args[1].Args[0].String)
	require.Equal(t, "c", p.Args[1].Args[1].String)
}

func TestMarshalParamsPath_InsertPrim(t *testing.T) {
	// Create a Pair tree and insert values at left and right.
	paths := [][]int{{0}, {1}}
	p, err := MarshalParamsPath(false, paths, "left", "right")
	require.NoError(t, err)

	require.Equal(t, micheline.D_PAIR, p.OpCode)
	require.Equal(t, "left", p.Args[0].String)
	require.Equal(t, "right", p.Args[1].String)
}
