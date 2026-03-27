package micheline

import (
	"encoding/json"
	"testing"

	"github.com/trilitech/tzgo/tezos"
)

func TestBuilderAndUtilHelpers(t *testing.T) {
	// builder funcs that were previously uncovered
	if NewZ(tezos.NewZ(123)).Type != PrimInt {
		t.Fatalf("NewZ type mismatch")
	}
	if NewMutez(tezos.N(123)).Type != PrimInt {
		t.Fatalf("NewMutez type mismatch")
	}
	if NewNat(tezos.NewZ(1).Big()).Type != PrimInt {
		t.Fatalf("NewNat type mismatch")
	}
	if NewMapType(NewPrim(T_STRING), NewPrim(T_BYTES)).OpCode != T_MAP {
		t.Fatalf("NewMapType opcode mismatch")
	}
	if NewSetType(NewPrim(T_NAT)).OpCode != T_SET {
		t.Fatalf("NewSetType opcode mismatch")
	}
	if NewOptType(NewPrim(T_NAT)).OpCode != T_OPTION {
		t.Fatalf("NewOptType opcode mismatch")
	}
	if NewCombPair(NewInt64(1), NewInt64(2)).Type != PrimSequence {
		t.Fatalf("NewCombPair type mismatch")
	}
	if NewCombPairType(NewPrim(T_INT), NewPrim(T_NAT)).OpCode != T_PAIR {
		t.Fatalf("NewCombPairType opcode mismatch")
	}
	if got := NewPrim(T_INT).WithAnno("%x"); len(got.Anno) != 1 {
		t.Fatalf("WithAnno failed")
	}

	// macros
	if ASSERT_CMPEQ().Type != PrimSequence {
		t.Fatalf("ASSERT_CMPEQ type mismatch")
	}
	if DUUP().Type != PrimSequence {
		t.Fatalf("DUUP type mismatch")
	}
	if IFCMPNEQ(NewSeq(), NewSeq()).Type != PrimSequence {
		t.Fatalf("IFCMPNEQ type mismatch")
	}
	if UNPAIR().Type != PrimSequence {
		t.Fatalf("UNPAIR type mismatch")
	}

	// util helpers
	if limit("abcdefghij", 5) != "abcde" || limit("abc", 5) != "abc" {
		t.Fatalf("limit mismatch")
	}
	if min(10, 5) != 5 || min(3, 5) != 3 {
		t.Fatalf("min mismatch")
	}
}

func TestFeatures_Prim(t *testing.T) {
	p := NewSeq(
		NewCode(I_CREATE_ACCOUNT),
		NewCode(I_CHAIN_ID),
		NewCode(I_TICKET),
		NewCode(H_CONSTANT, NewString("exprtrNvsCL6MY7rfibVC6t8uqVJgAXfjGDyBzxonZeUpiunP5P9KC")),
		NewCode(K_VIEW, NewString("v"), NewPrim(T_NAT), NewPrim(T_INT), NewSeq(NewCode(I_UNIT))),
		NewPrim(T_CHEST_KEY),
	)
	f := p.Features()
	if !f.Contains(FeatureAccountFactory) || !f.Contains(FeatureChainId) || !f.Contains(FeatureTicket) ||
		!f.Contains(FeatureGlobalConstant) || !f.Contains(FeatureView) || !f.Contains(FeatureTimelock) {
		t.Fatalf("unexpected features set: %v (%s)", uint16(f), f.String())
	}
	if got := f.Array(); len(got) == 0 {
		t.Fatalf("Array empty")
	}
	if s := f.String(); s == "" {
		t.Fatalf("String empty")
	}
	if _, err := json.Marshal(f); err != nil {
		t.Fatalf("MarshalJSON err=%v", err)
	}
}
