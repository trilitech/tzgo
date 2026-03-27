package bind

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/tezos"
)

var testAddress = tezos.MustParseAddress("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx")

func TestUnmarshalPrim(t *testing.T) {
	cases := map[string]struct {
		prim micheline.Prim
		dst  any
		want any
		wErr error
	}{
		"string":        {prim: micheline.NewString("hello"), dst: "", want: "hello"},
		"bigInt":        {prim: micheline.NewInt64(42), dst: &big.Int{}, want: big.NewInt(42)},
		"bytes":         {prim: micheline.NewBytes([]byte{4, 2}), dst: []byte{}, want: []byte{4, 2}},
		"address":       {prim: micheline.NewString(testAddress.String()), dst: tezos.Address{}, want: testAddress},
		"string slice":  {prim: micheline.NewSeq(micheline.NewString("1"), micheline.NewString("2")), dst: []string{}, want: []string{"1", "2"}},
		"struct":        {prim: micheline.NewPair(micheline.NewString("aaa"), micheline.NewPair(micheline.NewInt64(42), micheline.NewBytes([]byte{1, 2, 3}))), dst: (*unmarshaler)(nil), want: &unmarshaler{"aaa", big.NewInt(42), []byte{1, 2, 3}}},
		"nested struct": {prim: micheline.NewPair(micheline.NewPair(micheline.NewString("aaa"), micheline.NewPair(micheline.NewInt64(42), micheline.NewBytes([]byte{1, 2, 3}))), micheline.NewString("uuu")), dst: (*nestedUnmarshaler)(nil), want: &nestedUnmarshaler{&unmarshaler{"aaa", big.NewInt(42), []byte{1, 2, 3}}, "uuu"}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			val := reflect.New(reflect.TypeOf(c.dst))
			val.Elem().Set(reflect.ValueOf(c.dst))

			err := UnmarshalPrim(c.prim, val.Interface())
			if c.wErr != nil {
				require.ErrorIs(t, err, c.wErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, c.want, val.Elem().Interface())
		})
	}
}

type unmarshaler struct {
	A string
	B *big.Int
	C []byte
}

func (u *unmarshaler) UnmarshalPrim(prim micheline.Prim) error {
	return UnmarshalPrimPaths(prim, map[string]any{"l": &u.A, "r/l": &u.B, "r/r": &u.C})
}

type nestedUnmarshaler struct {
	U *unmarshaler
	S string
}

func (u *nestedUnmarshaler) UnmarshalPrim(prim micheline.Prim) error {
	return UnmarshalPrimPaths(prim, map[string]any{"l": &u.U, "r": &u.S})
}

func TestUnmarshalPrim_PointerInit(t *testing.T) {
	var i *big.Int
	require.Nil(t, i)

	err := UnmarshalPrim(micheline.NewInt64(42), &i)
	require.NoError(t, err)
	require.NotNil(t, i)
	require.Equal(t, big.NewInt(42), i)
}

func TestUnmarshalPrim_TimeFromString(t *testing.T) {
	var ts time.Time
	err := UnmarshalPrim(micheline.NewString("1970-01-01T00:00:01Z"), &ts)
	require.NoError(t, err)
	require.Equal(t, time.Unix(1, 0).UTC(), ts.UTC())
}

func TestUnmarshalPrim_AddressFromStringAndBytes(t *testing.T) {
	addr := tezos.MustParseAddress("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb")

	var a1 tezos.Address
	require.NoError(t, UnmarshalPrim(micheline.NewString(addr.String()), &a1))
	require.True(t, a1.Equal(addr))

	var a2 tezos.Address
	require.NoError(t, UnmarshalPrim(micheline.NewBytes(addr.EncodePadded()), &a2))
	require.True(t, a2.Equal(addr))
}

func TestUnmarshalPrim_SequenceIntoSlice(t *testing.T) {
	a1 := tezos.MustParseAddress("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb")
	a2 := tezos.MustParseAddress("tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6")
	prim := micheline.NewSeq(
		micheline.NewString(a1.String()),
		micheline.NewString(a2.String()),
	)

	var out []tezos.Address
	require.NoError(t, UnmarshalPrim(prim, &out))
	require.Len(t, out, 2)
	require.True(t, out[0].Equal(a1))
	require.True(t, out[1].Equal(a2))
}

func TestUnmarshalPrim_RequiresPointer(t *testing.T) {
	var s string
	require.Error(t, UnmarshalPrim(micheline.NewString("x"), s))
}
