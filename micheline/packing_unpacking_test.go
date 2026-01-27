package micheline

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"testing"
)

func TestPrim_PackUnpack_Roundtrip(t *testing.T) {
	orig := NewString("hello")
	packed := Prim{Type: PrimBytes, Bytes: orig.Pack()}
	if !packed.IsPacked() {
		t.Fatalf("expected packed bytes to be IsPacked()")
	}
	up, err := packed.Unpack()
	if err != nil {
		t.Fatalf("Unpack err=%v", err)
	}
	if !up.WasPacked {
		t.Fatalf("expected WasPacked=true")
	}
	if up.Type != PrimString || up.String != "hello" {
		t.Fatalf("unexpected unpack result: %s", up.Dump())
	}
}

func TestMainnetPackedBytes_UnpackKeyAndValue(t *testing.T) {
	// This fixture is a bytes/bytes big_map where both keys and values are PACK'ed bytes.
	// It gives us real-world packed data (including packed Michelson code).
	const cat = "bigmap"
	scanTestFiles(t, cat)

	var (
		next int
		err  error
	)
	found := false
	for {
		var tests []testcase
		next, err = loadNextTestFile(cat, next, &tests)
		if err != nil {
			if err == io.EOF {
				break
			}
			if len(tests) == 0 {
				break
			}
		}
		for _, tc := range tests {
			// Look for one of the KT1PWx2... packed bigmap fixtures.
			// (The key is packed "paused", want_key is "paused".)
			if tc.WantKey == nil || tc.WantValue == nil {
				continue
			}
			var wantKey string
			if err := json.Unmarshal(tc.WantKey, &wantKey); err != nil || wantKey != "paused" {
				continue
			}

			found = true
			typ := checkTypeEncoding(t, tc)
			key := checkKeyEncoding(t, tc)
			val := checkValueEncoding(t, tc)

			// Key: bytes -> unpacked string
			k, err := NewKey(typ.Left(), key)
			if err != nil {
				t.Fatalf("NewKey err=%v", err)
			}
			if !k.IsPacked() {
				t.Fatalf("expected key to be packed")
			}
			uk, err := k.Unpack()
			if err != nil {
				t.Fatalf("key unpack err=%v", err)
			}
			kbuf, err := uk.MarshalJSON()
			if err != nil {
				t.Fatalf("key MarshalJSON err=%v", err)
			}
			if !jsonDiff(t, kbuf, tc.WantKey) {
				t.Fatalf("unpacked key mismatch")
			}

			// Value: bytes -> unpacked primitive rendered into want_value (bool here).
			v := Value{
				Type:   typ.Right(),
				Value:  val,
				Render: RENDER_TYPE_FAIL,
			}
			if !v.IsPackedAny() {
				t.Fatalf("expected value to be packed")
			}
			uv, err := v.UnpackAll()
			if err != nil {
				t.Fatalf("value unpack err=%v", err)
			}
			vbuf, err := uv.MarshalJSON()
			if err != nil {
				t.Fatalf("value MarshalJSON err=%v", err)
			}
			if !jsonDiff(t, vbuf, tc.WantValue) {
				t.Fatalf("unpacked value mismatch")
			}
			return
		}
	}
	if !found {
		t.Fatalf("did not find packed-bytes fixture testcase in %s corpus", cat)
	}
}

func TestPrim_UnpackAllAsciiStrings(t *testing.T) {
	p := NewPair(
		NewBytes([]byte("ascii")),
		NewBytes([]byte{0xff, 0x00}), // not ascii
	)
	up := p.UnpackAllAsciiStrings()
	if up.Args[0].Type != PrimString || up.Args[0].String != "ascii" {
		t.Fatalf("expected first arg to become string, got %s", up.Args[0].Dump())
	}
	if up.Args[1].Type != PrimBytes {
		t.Fatalf("expected second arg to remain bytes, got %s", up.Args[1].Dump())
	}
}

func TestPrim_Unpack_Errors(t *testing.T) {
	p := Prim{Type: PrimBytes, Bytes: []byte{0x00, 0x01}}
	if p.IsPacked() {
		t.Fatalf("unexpected IsPacked=true")
	}
	if _, err := p.Unpack(); err == nil {
		t.Fatalf("expected error for unpacking non-packed prim")
	}

	// malformed packed bytes should not panic; it should return an error
	bad := Prim{Type: PrimBytes, Bytes: []byte{0x05, 0xff}}
	if !bad.IsPacked() {
		// isPackedBytes requires b[1] <= 0x0A; force that
		bad.Bytes[1] = 0x0A
	}
	if _, err := bad.Unpack(); err == nil {
		t.Fatalf("expected error for malformed packed prim")
	}
}

func TestPrim_Pack_EncodesPrefix(t *testing.T) {
	p := NewInt64(1)
	b := p.Pack()
	if len(b) == 0 || b[0] != 0x05 {
		t.Fatalf("Pack prefix mismatch: %x", b)
	}
	// make sure the remainder is valid binary prim
	var out Prim
	if err := out.UnmarshalBinary(b[1:]); err != nil {
		t.Fatalf("packed content UnmarshalBinary err=%v", err)
	}
}

func TestPrim_IsPacked_AsciiAndAddressBytes(t *testing.T) {
	ascii := Prim{Type: PrimBytes, Bytes: []byte("abc")}
	if !ascii.IsPacked() {
		t.Fatalf("expected ascii bytes to be considered packed")
	}
	up, err := ascii.Unpack()
	if err != nil {
		t.Fatalf("ascii Unpack err=%v", err)
	}
	if up.Type != PrimString || up.String != "abc" {
		t.Fatalf("ascii unpack mismatch: %s", up.Dump())
	}

	// address bytes should unpack into a string address
	addrBytes, _ := hex.DecodeString("000082eea640431c731757a8b7ee226da7a784797dde")
	addr := Prim{Type: PrimBytes, Bytes: addrBytes}
	if !addr.IsPacked() {
		t.Fatalf("expected address bytes to be considered packed")
	}
	aup, err := addr.Unpack()
	if err != nil {
		t.Fatalf("address Unpack err=%v", err)
	}
	if aup.Type != PrimString || aup.String == "" {
		t.Fatalf("address unpack mismatch: %s", aup.Dump())
	}
}
