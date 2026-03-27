package bind

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/micheline"
)

func TestOption_MarshalUnmarshalPrim(t *testing.T) {
	o := Some("hello")
	p, err := o.MarshalPrim(false)
	require.NoError(t, err)
	require.Equal(t, micheline.D_SOME, p.OpCode)
	require.Len(t, p.Args, 1)
	require.Equal(t, "hello", p.Args[0].String)

	var round Option[string]
	require.NoError(t, round.UnmarshalPrim(p))
	v, ok := round.Get()
	require.True(t, ok)
	require.Equal(t, "hello", v)

	none := None[string]()
	p, err = none.MarshalPrim(false)
	require.NoError(t, err)
	require.Equal(t, micheline.D_NONE, p.OpCode)

	require.NoError(t, round.UnmarshalPrim(p))
	require.True(t, round.IsNone())
}

func TestOption_SetUntyped(t *testing.T) {
	var o Option[string]

	require.NoError(t, o.SetUntyped("x"))
	require.True(t, o.IsSome())
	require.Equal(t, "Some(x)", o.String())

	require.NoError(t, o.SetUntyped(nil))
	require.True(t, o.IsNone())
	require.Equal(t, "None", o.String())

	require.Error(t, o.SetUntyped(123))
}
