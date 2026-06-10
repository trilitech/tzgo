package tezos

import "testing"

func TestOpStatus_ParseAndText(t *testing.T) {
	cases := []struct {
		s    string
		want OpStatus
		ok   bool
	}{
		{"applied", OpStatusApplied, true},
		{"failed", OpStatusFailed, true},
		{"skipped", OpStatusSkipped, true},
		{"backtracked", OpStatusBacktracked, true},
		{"nope", OpStatusInvalid, false},
	}
	for _, c := range cases {
		got := ParseOpStatus(c.s)
		if got != c.want {
			t.Fatalf("ParseOpStatus(%q)=%v want %v", c.s, got, c.want)
		}
		if got.IsValid() != c.ok {
			t.Fatalf("IsValid(%q)=%v want %v", c.s, got.IsValid(), c.ok)
		}
		if got.IsSuccess() != (got == OpStatusApplied) {
			t.Fatalf("IsSuccess(%q) mismatch", c.s)
		}
		if got.IsValid() {
			b, err := got.MarshalText()
			if err != nil || string(b) != c.s {
				t.Fatalf("MarshalText(%q) got=%q err=%v", c.s, string(b), err)
			}
			var out OpStatus
			if err := out.UnmarshalText([]byte(c.s)); err != nil || out != got {
				t.Fatalf("UnmarshalText(%q) out=%v err=%v", c.s, out, err)
			}
		} else {
			var out OpStatus
			if err := out.UnmarshalText([]byte(c.s)); err == nil {
				t.Fatalf("UnmarshalText(%q) expected error", c.s)
			}
		}
	}
}

func TestOpType_ParseAndTags(t *testing.T) {
	// basic parse/marshal
	if got := ParseOpType("transaction"); got != OpTypeTransaction {
		t.Fatalf("ParseOpType(transaction)=%v", got)
	}
	if got := ParseOpType("nope"); got != OpTypeInvalid {
		t.Fatalf("ParseOpType(nope)=%v", got)
	}
	var ot OpType
	if err := ot.UnmarshalText([]byte("transaction")); err != nil || ot != OpTypeTransaction {
		t.Fatalf("OpType.UnmarshalText err=%v ot=%v", err, ot)
	}
	if _, err := OpTypeTransaction.MarshalText(); err != nil {
		t.Fatalf("OpType.MarshalText err=%v", err)
	}

	// tag functions: known types should have stable tags for their version
	if tag := OpTypeTransaction.TagVersion(2); tag != 108 {
		t.Fatalf("TagVersion(transaction,v2)=%d want 108", tag)
	}
	if tag := OpTypeTransaction.Tag(); tag != 108 {
		t.Fatalf("Tag(transaction)=%d want 108", tag)
	}
	if ParseOpTag(108) != OpTypeTransaction {
		t.Fatalf("ParseOpTag(108) mismatch")
	}
	if ParseOpTag(255) != OpTypeInvalid {
		t.Fatalf("ParseOpTag(255) expected invalid")
	}
	if ParseOpTagVersion(108, 2) != OpTypeTransaction {
		t.Fatalf("ParseOpTagVersion(108,v2) mismatch")
	}
	// In tag version 0, transaction uses tag 8 (not 108).
	if ParseOpTagVersion(8, 0) != OpTypeTransaction {
		t.Fatalf("ParseOpTagVersion(8,v0) mismatch")
	}
	if ParseOpTagVersion(108, 0) != OpTypeInvalid {
		t.Fatalf("ParseOpTagVersion(108,v0) expected invalid")
	}

	// sizes and list ids are used for off-chain encoding/validation
	if OpTypeTransaction.MinSizeVersion(3) == 0 || OpTypeTransaction.MinSize() == 0 {
		t.Fatalf("MinSize should be >0 for transaction")
	}
	if OpTypeTransaction.ListId() != 3 {
		t.Fatalf("ListId(transaction)=%d want 3", OpTypeTransaction.ListId())
	}
	if OpTypeInvalid.ListId() != -1 {
		t.Fatalf("ListId(invalid)=%d want -1", OpTypeInvalid.ListId())
	}
}
