package bind

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/micheline"
)

func TestOr_MarshalUnmarshalPrim(t *testing.T) {
	l := Left[string, *big.Int]("a")
	p, err := l.MarshalPrim(false)
	require.NoError(t, err)
	require.Equal(t, micheline.D_LEFT, p.OpCode)
	require.Len(t, p.Args, 1)
	require.Equal(t, "a", p.Args[0].String)

	var round Or[string, *big.Int]
	require.NoError(t, round.UnmarshalPrim(p))
	v, ok := round.Left()
	require.True(t, ok)
	require.Equal(t, "a", v)

	r := Right[string, *big.Int](big.NewInt(42))
	p, err = r.MarshalPrim(false)
	require.NoError(t, err)
	require.Equal(t, micheline.D_RIGHT, p.OpCode)
	require.Len(t, p.Args, 1)
	require.Equal(t, int64(42), p.Args[0].Int.Int64())

	require.NoError(t, round.UnmarshalPrim(p))
	v2, ok := round.Right()
	require.True(t, ok)
	require.Equal(t, big.NewInt(42), v2)
}
