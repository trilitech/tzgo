package micheline

import (
	"encoding/hex"
	"io"
	"testing"
)

func TestRoundtripInvariants_TypeAndPrim_BinaryAndJSON(t *testing.T) {
	categories := []string{"bigmap", "storage", "params"}
	for _, cat := range categories {
		t.Run(cat, func(t *testing.T) {
			var (
				next int
				err  error
			)
			scanTestFiles(t, cat)
			for {
				var tests []testcase
				next, err = loadNextTestFile(cat, next, &tests)
				if err != nil {
					if err == io.EOF {
						break
					}
					// keep going: we want best-effort coverage across the corpus
					t.Errorf("loadNextTestFile(%s,%d): %v", cat, next, err)
					if len(tests) == 0 {
						break
					}
				}
				for _, tc := range tests {
					tc := tc
					t.Run(tc.Name, func(t *testing.T) {
						// Type: hex -> Type -> bin -> Type (semantic equality with annots)
						tb, err := hex.DecodeString(tc.TypeHex)
						if err != nil {
							t.Fatalf("TypeHex decode: %v", err)
						}
						var t1 Type
						if err := t1.UnmarshalBinary(tb); err != nil {
							t.Fatalf("Type.UnmarshalBinary: %v", err)
						}
						tbin, err := t1.MarshalBinary()
						if err != nil {
							t.Fatalf("Type.MarshalBinary: %v", err)
						}
						var t2 Type
						if err := t2.UnmarshalBinary(tbin); err != nil {
							t.Fatalf("Type.UnmarshalBinary(roundtrip): %v", err)
						}
						if !t1.IsEqualWithAnno(t2) {
							t.Fatalf("Type binary roundtrip mismatch:\nA=%s\nB=%s", t1.Dump(), t2.Dump())
						}

						// Type: Micheline JSON -> Type (semantic equality with annots)
						// Use the Prim Micheline JSON encoder because encoding/json does not
						// serialize the Type wrapper in Micheline format.
						var tj1 Type
						if err := tj1.UnmarshalJSON(tc.Type); err != nil {
							t.Fatalf("Type.UnmarshalJSON: %v", err)
						}
						jb, err := tj1.Prim.MarshalJSON()
						if err != nil {
							t.Fatalf("Prim.MarshalJSON(type): %v", err)
						}
						var tj2 Type
						if err := tj2.UnmarshalJSON(jb); err != nil {
							t.Fatalf("Type.UnmarshalJSON(roundtrip): %v", err)
						}
						// Use Dump() equality here: JSON roundtrips may normalize the internal PrimType
						// (binary/unary/variadic), but the Micheline tree should stay identical.
						if tj1.Dump() != tj2.Dump() {
							t.Fatalf("Type json roundtrip mismatch:\nA=%s\nB=%s", tj1.Dump(), tj2.Dump())
						}

						// Value: hex -> Prim -> bin -> Prim (semantic equality with annots)
						vb, err := hex.DecodeString(tc.ValueHex)
						if err != nil {
							t.Fatalf("ValueHex decode: %v", err)
						}
						var p1 Prim
						if err := p1.UnmarshalBinary(vb); err != nil {
							t.Fatalf("Prim.UnmarshalBinary(value): %v", err)
						}
						pbin, err := p1.MarshalBinary()
						if err != nil {
							t.Fatalf("Prim.MarshalBinary(value): %v", err)
						}
						var p2 Prim
						if err := p2.UnmarshalBinary(pbin); err != nil {
							t.Fatalf("Prim.UnmarshalBinary(value roundtrip): %v", err)
						}
						if !p1.IsEqualWithAnno(p2) {
							t.Fatalf("value Prim binary roundtrip mismatch:\nA=%s\nB=%s", p1.Dump(), p2.Dump())
						}

						// Value: json -> Prim -> json -> Prim (semantic equality with annots)
						var pj1 Prim
						if err := pj1.UnmarshalJSON(tc.Value); err != nil {
							t.Fatalf("Prim.UnmarshalJSON(value): %v", err)
						}
						jvb, err := pj1.MarshalJSON()
						if err != nil {
							t.Fatalf("Prim.MarshalJSON(value): %v", err)
						}
						var pj2 Prim
						if err := pj2.UnmarshalJSON(jvb); err != nil {
							t.Fatalf("Prim.UnmarshalJSON(value roundtrip): %v", err)
						}
						if pj1.Dump() != pj2.Dump() {
							t.Fatalf("value Prim json roundtrip mismatch:\nA=%s\nB=%s", pj1.Dump(), pj2.Dump())
						}

						// Key (optional): only when present in testcase
						if len(tc.KeyHex) > 0 && len(tc.Key) > 0 {
							kb, err := hex.DecodeString(tc.KeyHex)
							if err != nil {
								t.Fatalf("KeyHex decode: %v", err)
							}
							var k1 Prim
							if err := k1.UnmarshalBinary(kb); err != nil {
								t.Fatalf("Prim.UnmarshalBinary(key): %v", err)
							}
							kbin, err := k1.MarshalBinary()
							if err != nil {
								t.Fatalf("Prim.MarshalBinary(key): %v", err)
							}
							var k2 Prim
							if err := k2.UnmarshalBinary(kbin); err != nil {
								t.Fatalf("Prim.UnmarshalBinary(key roundtrip): %v", err)
							}
							if !k1.IsEqualWithAnno(k2) {
								t.Fatalf("key Prim binary roundtrip mismatch:\nA=%s\nB=%s", k1.Dump(), k2.Dump())
							}
						}
					})
				}
			}
		})
	}
}
