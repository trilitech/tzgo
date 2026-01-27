package micheline

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"testing"
)

func TestMainnetFixtureCorpus_DecoderRobustness(t *testing.T) {
	// This is intentionally "robustness", not "correctness":
	// we only assert the decoders can consume the entire fixture corpus without panicking.
	categories := []string{"bigmap", "storage", "params"}
	for _, cat := range categories {
		t.Run(cat, func(t *testing.T) {
			var (
				next int
				err  error
				n    int
			)
			scanTestFiles(t, cat)
			for {
				var tests []testcase
				next, err = loadNextTestFile(cat, next, &tests)
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Errorf("loadNextTestFile(%s,%d): %v", cat, next, err)
					if len(tests) == 0 {
						break
					}
				}
				for _, tc := range tests {
					n++
					// json decode
					var tt Type
					if err := tt.UnmarshalJSON(tc.Type); err != nil {
						t.Fatalf("Type.UnmarshalJSON (%s): %v", tc.Name, err)
					}
					var pv Prim
					if err := pv.UnmarshalJSON(tc.Value); err != nil {
						t.Fatalf("Prim.UnmarshalJSON(value) (%s): %v", tc.Name, err)
					}
					if len(tc.Key) > 0 {
						var pk Prim
						if err := pk.UnmarshalJSON(tc.Key); err != nil {
							t.Fatalf("Prim.UnmarshalJSON(key) (%s): %v", tc.Name, err)
						}
					}

					// hex decode
					tb, err := hex.DecodeString(tc.TypeHex)
					if err != nil {
						t.Fatalf("TypeHex decode (%s): %v", tc.Name, err)
					}
					if err := tt.UnmarshalBinary(tb); err != nil {
						t.Fatalf("Type.UnmarshalBinary (%s): %v", tc.Name, err)
					}
					vb, err := hex.DecodeString(tc.ValueHex)
					if err != nil {
						t.Fatalf("ValueHex decode (%s): %v", tc.Name, err)
					}
					if err := pv.UnmarshalBinary(vb); err != nil {
						t.Fatalf("Prim.UnmarshalBinary(value) (%s): %v", tc.Name, err)
					}
					if len(tc.KeyHex) > 0 {
						kb, err := hex.DecodeString(tc.KeyHex)
						if err != nil {
							t.Fatalf("KeyHex decode (%s): %v", tc.Name, err)
						}
						var pk Prim
						if err := pk.UnmarshalBinary(kb); err != nil {
							t.Fatalf("Prim.UnmarshalBinary(key) (%s): %v", tc.Name, err)
						}
					}

					// encode paths should not error
					if _, err := tt.MarshalBinary(); err != nil {
						t.Fatalf("Type.MarshalBinary (%s): %v", tc.Name, err)
					}
					if _, err := pv.MarshalBinary(); err != nil {
						t.Fatalf("Prim.MarshalBinary(value) (%s): %v", tc.Name, err)
					}

					// ensure JSON marshaling is safe too
					if _, err := json.Marshal(tt); err != nil {
						t.Fatalf("json.Marshal(Type) (%s): %v", tc.Name, err)
					}
					if _, err := json.Marshal(pv); err != nil {
						t.Fatalf("json.Marshal(Prim) (%s): %v", tc.Name, err)
					}
				}
			}
			if n == 0 {
				t.Fatalf("no fixtures found for category %s", cat)
			}
		})
	}
}
