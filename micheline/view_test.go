package micheline

import (
	"testing"
)

func TestView_Basics(t *testing.T) {
	p := NewCode(K_VIEW, NewString("my_view"), NewPrim(T_NAT), NewPrim(T_INT), NewSeq(NewCode(I_UNIT)))
	v := NewView(p)
	if !v.IsValid() {
		t.Fatalf("view invalid")
	}
	if v.Name != "my_view" {
		t.Fatalf("name=%q want my_view", v.Name)
	}
	v2 := v.Clone()
	if !v.IsEqual(v2) || !v.IsEqualWithAnno(v2) || !v.IsEqualWithCode(v2) {
		t.Fatalf("cloned view not equal")
	}
	if _, err := v.MarshalJSON(); err != nil {
		t.Fatalf("MarshalJSON err=%v", err)
	}
	td := v.Typedef()
	if td.Type != K_VIEW.String() {
		t.Fatalf("Typedef type=%q want %q", td.Type, K_VIEW.String())
	}
	if v.TypedefPtr("x") == nil {
		t.Fatalf("TypedefPtr returned nil")
	}
}
