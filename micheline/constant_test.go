package micheline

import (
	"testing"

	"github.com/trilitech/tzgo/tezos"
)

func TestConstantDict_And_PrimConstants(t *testing.T) {
	var d ConstantDict
	h := tezos.MustParseExprHash("exprtrNvsCL6MY7rfibVC6t8uqVJgAXfjGDyBzxonZeUpiunP5P9KC")
	val := NewCode(H_CONSTANT, NewString(h.String()))

	if d.Has(h) {
		t.Fatalf("Has on nil dict want false")
	}
	if _, ok := d.Get(h); ok {
		t.Fatalf("Get on nil dict want ok=false")
	}

	d.Add(h, val)
	if !d.Has(h) {
		t.Fatalf("Has after Add want true")
	}
	if got, ok := d.Get(h); !ok || !got.IsConstant() {
		t.Fatalf("Get after Add ok=%v got=%v", ok, got)
	}
	if got, ok := d.GetString(h.String()); !ok || !got.IsConstant() {
		t.Fatalf("GetString after Add ok=%v got=%v", ok, got)
	}

	// Prim.Constants walks the tree and extracts expr hashes.
	c := NewSeq(val, NewCode(D_UNIT))
	list := c.Constants()
	if len(list) != 1 || list[0] != h {
		t.Fatalf("Constants=%v want [%s]", list, h)
	}
}
