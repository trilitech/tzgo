// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package tezos

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

func MustDecodeString(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func TestAddress(t *testing.T) {
	type testcase struct {
		Address string
		Hash    string
		Type    AddressType
		Bytes   string
		Padded  string
	}

	cases := []testcase{
		// tz1
		{
			Address: "tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q",
			Hash:    "0b78887fdd0cd3bfbe75a717655728e0205bb958",
			Type:    AddressTypeEd25519,
			Bytes:   "000b78887fdd0cd3bfbe75a717655728e0205bb958",
			Padded:  "00000b78887fdd0cd3bfbe75a717655728e0205bb958",
		},
		// tz2
		{
			Address: "tz2VN9n2C56xGLykHCjhNvZQqUeTVisrHjxA",
			Hash:    "e6e7cfd00186c29ede318bef62ac85ddec8a50d5",
			Type:    AddressTypeSecp256k1,
			Bytes:   "01e6e7cfd00186c29ede318bef62ac85ddec8a50d5",
			Padded:  "0001e6e7cfd00186c29ede318bef62ac85ddec8a50d5",
		},
		// tz3
		{
			Address: "tz3Qa3kjWa6B3XgvZcVe24gTfjkc5WZRz59Q",
			Hash:    "2e8671595e32ddd3c1e3f229898e9bec727eca90",
			Type:    AddressTypeP256,
			Bytes:   "022e8671595e32ddd3c1e3f229898e9bec727eca90",
			Padded:  "00022e8671595e32ddd3c1e3f229898e9bec727eca90",
		},
		// KT1
		{
			Address: "KT1GyeRktoGPEKsWpchWguyy8FAf3aNHkw2T",
			Hash:    "5c149d65c5ca113bc2bc3c861ef6ea8030d71553",
			Type:    AddressTypeContract,
			Bytes:   "015c149d65c5ca113bc2bc3c861ef6ea8030d7155300",
			Padded:  "015c149d65c5ca113bc2bc3c861ef6ea8030d7155300",
		},
		// TODO: AddressTypeSapling
		// tz4
		{
			Address: "tz4HVR6aty9KwsQFHh81C1G7gBdhxT8kuytm",
			Hash:    "5d1497f39b87599983fe8f29599b679564be822d",
			Type:    AddressTypeBls12_381,
			Bytes:   "035d1497f39b87599983fe8f29599b679564be822d",
			Padded:  "00035d1497f39b87599983fe8f29599b679564be822d",
		},
		// tz5 (v025 Ushuaia, ML-DSA-44 PKH only)
		{
			Address: "tz5VWE3unqGsLVrYhxCGBxiVVYXDjHHbmbTY",
			Hash:    "5d1497f39b87599983fe8f29599b679564be822d",
			Type:    AddressTypeMlDsa44,
			Bytes:   "045d1497f39b87599983fe8f29599b679564be822d",
			Padded:  "00045d1497f39b87599983fe8f29599b679564be822d",
		},
		// txr1
		{
			Address: "txr1QVAMSfhGduYQoQwrWroJW5b2796Qmb9ej",
			Hash:    "202e50c8eed224f3961d83522039be4eee40633d",
			Type:    AddressTypeTxRollup,
			Bytes:   "02202e50c8eed224f3961d83522039be4eee40633d00",
			Padded:  "02202e50c8eed224f3961d83522039be4eee40633d00",
		},
		// sr1
		{
			Address: "sr1Fq8fPi2NjhWUXtcXBggbL6zFjZctGkmso",
			Hash:    "6b6209e8037138491d8d5d8ee340000d51b91581",
			Type:    AddressTypeSmartRollup,
			Bytes:   "036b6209e8037138491d8d5d8ee340000d51b9158100",
			Padded:  "036b6209e8037138491d8d5d8ee340000d51b9158100",
		},
	}

	for i, c := range cases {
		h := MustDecodeString(c.Hash)
		buf := MustDecodeString(c.Bytes)
		pad := MustDecodeString(c.Padded)

		// base58 must parse
		a, err := ParseAddress(c.Address)
		if err != nil {
			t.Fatalf("Case %d - parsing address %s: %v", i, c.Address, err)
		}

		// check type
		if got, want := a.Type(), c.Type; got != want {
			t.Errorf("Case %d - mismatched type got=%s want=%s", i, got, want)
		}

		// check hash
		if !bytes.Equal(a[1:], h) {
			t.Errorf("Case %d - mismatched hash got=%x want=%x", i, a[1:], h)
		}

		// check bytes
		if !bytes.Equal(a.Encode(), buf) {
			t.Errorf("Case %d - mismatched binary encoding got=%x want=%x", i, a.Encode(), buf)
		}

		// check padded bytes
		if !bytes.Equal(a.EncodePadded(), pad) {
			t.Errorf("Case %d - mismatched padded binary encoding got=%x want=%x", i, a.EncodePadded(), pad)
		}

		// marshal text
		out, err := a.MarshalText()
		if err != nil {
			t.Errorf("Case %d - marshal text unexpected error: %v", i, err)
		}

		if got, want := string(out), c.Address; got != want {
			t.Errorf("Case %d - mismatched text encoding got=%s want=%s", i, got, want)
		}

		// unmarshal from bytes
		var a2 Address
		err = a2.Decode(buf)
		if err != nil {
			t.Fatalf("Case %d - unmarshal binary %s: %v", i, c.Bytes, err)
		}

		if !a2.Equal(a) {
			t.Errorf("Case %d - mismatched address got=%s want=%s", i, a2, a)
		}

		// unmarshal from padded bytes
		err = a2.Decode(pad)
		if err != nil {
			t.Fatalf("Case %d - unmarshal binary %s: %v", i, c.Padded, err)
		}

		if !a2.Equal(a) {
			t.Errorf("Case %d - mismatched address got=%s want=%s", i, a2, a)
		}

		// unmarshal text
		err = a2.UnmarshalText([]byte(c.Address))
		if err != nil {
			t.Fatalf("Case %d - unmarshal text %s: %v", i, c.Address, err)
		}

		if !a2.Equal(a) {
			t.Errorf("Case %d - mismatched address got=%s want=%s", i, a2, a)
		}

		// marshal binary roundtrip
		out = a.Encode()
		err = a2.Decode(out)
		if err != nil {
			t.Fatalf("Case %d - binary roundtrip: %v", i, err)
		}

		if !a2.Equal(a) {
			t.Errorf("Case %d - mismatched binary roundtrip got=%s want=%s", i, a2, a)
		}
	}
}

func TestInvalidAddress(t *testing.T) {
	// invalid base58 string
	if _, err := ParseAddress("tz1KzpjBnunNJVABHBnzfG4iuLmphitExW2"); err == nil {
		t.Errorf("Expected error on invalid base58 string")
	}

	// init from invalid short hash
	hash := MustDecodeString("0b78887fdd0cd3bfbe75a717655728e0205bb9")
	a := NewAddress(AddressTypeEd25519, hash)
	if a.IsValid() {
		t.Errorf("Expected invalid address from short hash")
	}

	// init from invalid empty bytes
	a = NewAddress(AddressTypeEd25519, nil)
	if a.IsValid() {
		t.Errorf("Expected invalid address from nil hash")
	}

	// decode from short buffer
	err := a.Decode(MustDecodeString("000b78887fdd0cd3bfbe75a717655728e0205bb9"))
	if err == nil || a.IsValid() {
		t.Errorf("Expected unmarshal error from short buffer")
	}

	// decode from nil buffer
	err = a.Decode(nil)
	if err == nil || a.IsValid() {
		t.Errorf("Expected unmarshal error from short buffer")
	}

	// decode from invalid buffer (wrong type)
	err = a.Decode(MustDecodeString("00FF000b80d92ce17aa6070fde1a99288a4213a5b650"))
	if err == nil || a.IsValid() {
		t.Errorf("Expected unmarshal error from invalid buffer")
	}
}

func BenchmarkAddressDecode(b *testing.B) {
	b.SetBytes(21)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseAddress("tz3Qa3kjWa6B3XgvZcVe24gTfjkc5WZRz59Q")
	}
}

func BenchmarkAddressEncode(b *testing.B) {
	a, _ := ParseAddress("tz3Qa3kjWa6B3XgvZcVe24gTfjkc5WZRz59Q")
	b.SetBytes(21)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = a.String()
	}
}

// TestMlDsa44Address verifies tz5 (v025 Ushuaia) PKH-only semantics: implicit
// account classification and binary tag 4 ownership.
func TestMlDsa44Address(t *testing.T) {
	a, err := ParseAddress("tz5VWE3unqGsLVrYhxCGBxiVVYXDjHHbmbTY")
	if err != nil {
		t.Fatalf("parse tz5: %v", err)
	}
	if !a.IsEOA() {
		t.Errorf("tz5 address must classify as EOA")
	}
	if a.IsContract() || a.IsRollup() {
		t.Errorf("tz5 address must not classify as contract or rollup")
	}
	if got := a.KeyType(); got != KeyTypeMlDsa44 {
		t.Errorf("tz5 KeyType = %v, want KeyTypeMlDsa44", got)
	}
	// ML-DSA key material is unsupported in the PKH-only slice
	if KeyTypeMlDsa44.PkHashType().Len != 0 || KeyTypeMlDsa44.PkPrefixBytes() != nil {
		t.Errorf("ML-DSA pk decoding must remain unsupported")
	}
	if DetectAddressType("tz5VWE3unqGsLVrYhxCGBxiVVYXDjHHbmbTY") != AddressTypeMlDsa44 {
		t.Errorf("DetectAddressType must recognize tz5")
	}

	// malformed inputs must fail explicitly
	if _, err := ParseAddress("tz5VWE3unqGsLVrYhxCGBxiVVYXDjHHbmbTZ"); err == nil {
		t.Errorf("expected checksum error on corrupted tz5 string")
	}
	if _, err := ParseAddress("tz5VWE3unq"); err == nil {
		t.Errorf("expected error on truncated tz5 string")
	}

	// binary ML-DSA public keys (tag 4) remain unsupported with an explicit error
	var k Key
	buf := append([]byte{4}, make([]byte, 32)...)
	err = k.UnmarshalBinary(buf)
	if err == nil || !strings.Contains(err.Error(), "ML-DSA-44") {
		t.Errorf("binary key tag 4 must fail with explicit unsupported error, got %v", err)
	}
}

// TestBlindedAddressNoBinaryTag verifies that blinded (btz1) addresses no longer
// alias on-chain binary tag 4, which belongs to ML-DSA-44 (tz5) since v025.
// Text parsing of btz1 addresses is unchanged.
func TestBlindedAddressNoBinaryTag(t *testing.T) {
	a, err := ParseAddress("btz1LKs15uHQ4PgCoY3ZDq55CKJ5wDq9jQwfk")
	if err != nil {
		t.Fatalf("parse btz1: %v", err)
	}
	if got, want := a.Type(), AddressTypeBlinded; got != want {
		t.Fatalf("type = %v, want %v", got, want)
	}
	// hash check preserved from the former TestAddress table entry; only the
	// binary tag round-trip is intentionally no longer supported since v025
	if want := MustDecodeString("000b80d92ce17aa6070fde1a99288a4213a5b650"); !bytes.Equal(a[1:], want) {
		t.Errorf("hash = %x, want %x", a[1:], want)
	}
	if got, want := a.String(), "btz1LKs15uHQ4PgCoY3ZDq55CKJ5wDq9jQwfk"; got != want {
		t.Errorf("text round-trip = %s, want %s", got, want)
	}
	// binary encoding must not emit tag 4 anymore
	if enc := a.Encode(); len(enc) > 0 && enc[0] == 4 {
		t.Errorf("blinded address must not binary-encode with tag 4 (tz5)")
	}
	// tag 4 binary decodes as tz5, not blinded
	var b Address
	buf := MustDecodeString("045d1497f39b87599983fe8f29599b679564be822d")
	if err := b.Decode(buf); err != nil {
		t.Fatalf("decode tag 4: %v", err)
	}
	if got, want := b.Type(), AddressTypeMlDsa44; got != want {
		t.Errorf("tag 4 decodes as %v, want %v", got, want)
	}
}
