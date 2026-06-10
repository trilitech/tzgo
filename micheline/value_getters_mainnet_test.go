package micheline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func mustLoadMainnetTestcases(t *testing.T, relPath string) []testcase {
	t.Helper()
	buf, err := os.ReadFile(filepath.Join("testdata-mainnet", relPath))
	if err != nil {
		t.Fatal(err)
	}
	var tests []testcase
	if err := json.Unmarshal(buf, &tests); err != nil {
		t.Fatal(err)
	}
	if len(tests) == 0 {
		t.Fatalf("no testcases in %s", relPath)
	}
	return tests
}

func TestValue_Getters_FromMainnetFixtures(t *testing.T) {
	// Storage fixture has address, bool and int-like fields in want_value.
	st := mustLoadMainnetTestcases(t, filepath.Join("storage", "KT1AFA2mwNUMNd4SsujE1YYp29vd8BZejyKW.json"))[0]
	typ := checkTypeEncoding(t, st)
	val := checkValueEncoding(t, st)

	v := NewValue(typ, val)

	// Map + getters
	if _, err := v.Map(); err != nil {
		t.Fatalf("Map error: %v", err)
	}
	if got, ok := v.GetAddress("administrator"); !ok || !got.IsValid() {
		t.Fatalf("GetAddress(administrator) ok=%v got=%v", ok, got)
	}
	if got, ok := v.GetBool("paused"); !ok || got != false {
		t.Fatalf("GetBool(paused) ok=%v got=%v want=false", ok, got)
	}
	if got, ok := v.GetInt64("ledger"); !ok || got <= 0 {
		t.Fatalf("GetInt64(ledger) ok=%v got=%d want>0", ok, got)
	}

	// Walk + Unmarshal
	var sawPaused bool
	if err := v.Walk("", func(label string, _ interface{}) error {
		if label == "paused" {
			sawPaused = true
		}
		return nil
	}); err != nil {
		t.Fatalf("Walk error: %v", err)
	}
	if !sawPaused {
		t.Fatalf("Walk did not visit paused")
	}

	var out map[string]interface{}
	if err := v.Unmarshal(&out); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if _, ok := out["paused"]; !ok {
		t.Fatalf("Unmarshal output missing paused")
	}

	// Storage fixture contains keys.
	keys := mustLoadMainnetTestcases(t, filepath.Join("storage", "KT1QuofAgnsWffHzLA7D78rxytJruGHDe7XG.json"))[0]
	keysTyp := checkTypeEncoding(t, keys)
	keysVal := checkValueEncoding(t, keys)
	kv := NewValue(keysTyp, keysVal)
	if got, ok := kv.GetKey("key_info.key_groups.0.signatories.0"); !ok || !got.IsValid() {
		t.Fatalf("GetKey(key_info.key_groups.0.signatories.0) ok=%v got=%v", ok, got)
	}

	// Params fixture contains signatures.
	sigs := mustLoadMainnetTestcases(t, filepath.Join("params", "ooMwWU2b53uMg3cjPexf5H6pBkadfmiX2G546mxKNDCrTNVvGGJ.json"))[0]
	sigsTyp := checkTypeEncoding(t, sigs)
	sigsVal := checkValueEncoding(t, sigs)
	sv := NewValue(sigsTyp, sigsVal)
	if got, ok := sv.GetSignature("Action.signatures.0.1"); !ok || !got.IsValid() {
		t.Fatalf("GetSignature(Action.signatures.0.1) ok=%v got=%v", ok, got)
	}

	// GetBytes: use a small bytes-typed value to force hex string conversion.
	btyp := Type{Prim{Type: PrimNullary, OpCode: T_BYTES, Anno: []string{"%data"}}}
	bval := Prim{Type: PrimBytes, Bytes: []byte{0xde, 0xad, 0xbe, 0xef}}
	bv := NewValue(btyp, bval)
	if got, ok := bv.GetBytes("data"); !ok || len(got) != 4 {
		t.Fatalf("GetBytes(data) ok=%v got=%x", ok, got)
	}
}
