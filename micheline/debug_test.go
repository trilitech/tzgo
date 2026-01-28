package micheline

import (
	"testing"
)

func TestPrim_DebugDump_NoPanic(t *testing.T) {
	p := NewSeq(NewCode(I_UNIT), NewCode(I_NIL, NewPrim(T_OPERATION)))
	// Just ensure these functions run and return something.
	if s := p.Dump(); s == "" {
		t.Fatalf("Dump returned empty")
	}
	if s := p.DumpLimit(16); s == "" {
		t.Fatalf("DumpLimit returned empty")
	}
}
