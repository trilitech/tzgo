// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package tezos

import (
	"bytes"
	"encoding"
	"encoding/hex"
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
		// btz1
		{
			Address: "btz1LKs15uHQ4PgCoY3ZDq55CKJ5wDq9jQwfk",
			Hash:    "000b80d92ce17aa6070fde1a99288a4213a5b650",
			Type:    AddressTypeBlinded,
			Bytes:   "04000b80d92ce17aa6070fde1a99288a4213a5b650",
			Padded:  "0004000b80d92ce17aa6070fde1a99288a4213a5b650",
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

func TestAddressTypeHelpers(t *testing.T) {
	// ParseAddressType accepts either prefix (tz1/KT1/...) or canonical name (ed25519/...)
	if got := ParseAddressType("tz1"); got != AddressTypeEd25519 {
		t.Fatalf("ParseAddressType(tz1)=%v", got)
	}
	if got := ParseAddressType("ed25519"); got != AddressTypeEd25519 {
		t.Fatalf("ParseAddressType(ed25519)=%v", got)
	}
	if got := ParseAddressType("nope"); got != AddressTypeInvalid {
		t.Fatalf("ParseAddressType(nope)=%v", got)
	}

	// String/Prefix/KeyType
	if AddressTypeContract.String() == "" {
		t.Fatalf("AddressType.String empty")
	}
	if AddressTypeEd25519.Prefix() != "tz1" {
		t.Fatalf("AddressTypeEd25519.Prefix=%q", AddressTypeEd25519.Prefix())
	}
	if !AddressTypeEd25519.KeyType().IsValid() {
		t.Fatalf("expected ed25519 KeyType to be valid")
	}
	if AddressTypeContract.KeyType().IsValid() {
		t.Fatalf("expected contract KeyType to be invalid")
	}

	// Marshal/Unmarshal text for AddressType
	var _ encoding.TextMarshaler = AddressTypeEd25519
	var _ encoding.TextUnmarshaler = (*AddressType)(nil)
	b, err := AddressTypeSecp256k1.MarshalText()
	if err != nil || string(b) != AddressTypeSecp256k1.String() {
		t.Fatalf("MarshalText got=%q err=%v", string(b), err)
	}
	var at AddressType
	if err := at.UnmarshalText([]byte("tz3")); err != nil || at != AddressTypeP256 {
		t.Fatalf("UnmarshalText(tz3) got=%v err=%v", at, err)
	}
	if err := at.UnmarshalText([]byte("nope")); err == nil {
		t.Fatalf("UnmarshalText(nope) expected error")
	}

	// Prefix detection
	if !HasAddressPrefix("tz1burnburnburnburnburnburnburjAYjjX") {
		t.Fatalf("HasAddressPrefix(tz1...) expected true")
	}
	if HasAddressPrefix("not-an-address") {
		t.Fatalf("HasAddressPrefix(not...) expected false")
	}
}

func TestAddressMoreHelpers(t *testing.T) {
	tz1 := MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q")
	kt1 := MustParseAddress("KT1GyeRktoGPEKsWpchWguyy8FAf3aNHkw2T")
	txr1 := MustParseAddress("txr1QVAMSfhGduYQoQwrWroJW5b2796Qmb9ej")
	sr1 := MustParseAddress("sr1Fq8fPi2NjhWUXtcXBggbL6zFjZctGkmso")

	if !tz1.IsEOA() || tz1.IsContract() || tz1.IsRollup() {
		t.Fatalf("tz1 classification mismatch: eoa=%v contract=%v rollup=%v", tz1.IsEOA(), tz1.IsContract(), tz1.IsRollup())
	}
	if kt1.IsEOA() || !kt1.IsContract() || kt1.IsRollup() {
		t.Fatalf("kt1 classification mismatch: eoa=%v contract=%v rollup=%v", kt1.IsEOA(), kt1.IsContract(), kt1.IsRollup())
	}
	if !txr1.IsRollup() || !sr1.IsRollup() {
		t.Fatalf("rollup classification mismatch")
	}

	// Hash()/KeyType()
	if len(tz1.Hash()) != 20 {
		t.Fatalf("Hash len=%d want 20", len(tz1.Hash()))
	}
	if !tz1.KeyType().IsValid() {
		t.Fatalf("KeyType invalid for tz1")
	}

	// Clone() is a deep copy of bytes
	c := tz1.Clone()
	if !c.Equal(tz1) {
		t.Fatalf("Clone not equal")
	}
	c[1] ^= 0xff
	if c.Equal(tz1) {
		t.Fatalf("Clone shares backing array unexpectedly")
	}

	// MarshalBinary/UnmarshalBinary (tzgo 21-byte format)
	bin, err := tz1.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary err=%v", err)
	}
	var tz1b Address
	if err := tz1b.UnmarshalBinary(bin); err != nil {
		t.Fatalf("UnmarshalBinary err=%v", err)
	}
	if !tz1b.Equal(tz1) {
		t.Fatalf("binary roundtrip mismatch got=%s want=%s", tz1b, tz1)
	}
	// error branch
	if err := tz1b.UnmarshalBinary(bin[:20]); err == nil {
		t.Fatalf("UnmarshalBinary(short) expected error")
	}

	// IsAddressBytes detection
	if !IsAddressBytes(tz1.Encode()) || !IsAddressBytes(tz1.EncodePadded()) {
		t.Fatalf("IsAddressBytes expected true for tz1 encodings")
	}
	if IsAddressBytes([]byte{0x01, 0x02}) {
		t.Fatalf("IsAddressBytes expected false for short buffer")
	}
	if IsAddressBytes(make([]byte, 23)) {
		t.Fatalf("IsAddressBytes expected false for other lengths")
	}

	// Contract/rollup address string helpers (reinterpret the hash bytes)
	if got := tz1.ContractAddress(); got[:3] != "KT1" {
		t.Fatalf("ContractAddress prefix=%q", got[:3])
	}
	if got := tz1.TxRollupAddress(); got[:4] != "txr1" {
		t.Fatalf("TxRollupAddress prefix=%q", got[:4])
	}
	if got := tz1.SmartRollupAddress(); got[:3] != "sr1" {
		t.Fatalf("SmartRollupAddress prefix=%q", got[:3])
	}

	// flags.Value Set
	var a Address
	if err := a.Set(tz1.String()); err != nil {
		t.Fatalf("Address.Set err=%v", err)
	}
	if !a.Equal(tz1) {
		t.Fatalf("Address.Set mismatch got=%s want=%s", a, tz1)
	}

	// Address.UnmarshalText should ignore entrypoint suffix (e.g. KT1...%default)
	var a2 Address
	if err := a2.UnmarshalText([]byte(kt1.String() + "%default")); err != nil {
		t.Fatalf("Address.UnmarshalText(entrypoint) err=%v", err)
	}
	if !a2.Equal(kt1) {
		t.Fatalf("Address.UnmarshalText(entrypoint) mismatch got=%s want=%s", a2, kt1)
	}

	// Ensure EncodeAddress(AddressTypeInvalid, ...) is empty string.
	if EncodeAddress(AddressTypeInvalid, tz1.Hash()) != "" {
		t.Fatalf("EncodeAddress(invalid) expected empty")
	}

	// Decode should accept longer byte strings (suffix/padding) and ignore extra bytes.
	buf := append(tz1.EncodePadded(), bytes.Repeat([]byte{0x00}, 10)...)
	var a3 Address
	if err := a3.Decode(buf); err != nil {
		t.Fatalf("Decode(with suffix) err=%v", err)
	}
	if !a3.Equal(tz1) {
		t.Fatalf("Decode(with suffix) mismatch got=%s want=%s", a3, tz1)
	}
}
