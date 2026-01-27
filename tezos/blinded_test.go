package tezos

import (
	"bytes"
	"testing"
)

func TestBlindedAddressHelpers(t *testing.T) {
	// Use a deterministic secret (blake2b key) and a known address hash.
	secret := []byte("tzgo-test-secret")
	a := MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q")

	// BlindHash
	bh, err := BlindHash(a.Hash(), secret)
	if err != nil {
		t.Fatalf("BlindHash err=%v", err)
	}
	if len(bh) != 20 {
		t.Fatalf("BlindHash len=%d want 20", len(bh))
	}

	// BlindAddress + MatchBlindedAddress
	ba, err := BlindAddress(a, secret)
	if err != nil {
		t.Fatalf("BlindAddress err=%v", err)
	}
	if ba.Type() != AddressTypeBlinded {
		t.Fatalf("BlindAddress type=%v want %v", ba.Type(), AddressTypeBlinded)
	}
	if !bytes.Equal(ba.Hash(), bh) {
		t.Fatalf("BlindAddress hash mismatch")
	}
	if !MatchBlindedAddress(a, ba, secret) {
		t.Fatalf("MatchBlindedAddress expected true")
	}

	// EncodeBlindedAddress / DecodeBlindedAddress
	enc, err := EncodeBlindedAddress(a.Hash(), secret)
	if err != nil {
		t.Fatalf("EncodeBlindedAddress err=%v", err)
	}
	dec, err := DecodeBlindedAddress(enc)
	if err != nil {
		t.Fatalf("DecodeBlindedAddress err=%v", err)
	}
	if !dec.Equal(ba) {
		t.Fatalf("DecodeBlindedAddress mismatch got=%s want=%s", dec, ba)
	}

	// error branch in BlindHash: blake2b keys > 64 bytes are rejected
	longSecret := bytes.Repeat([]byte{0x01}, 65)
	if _, err := BlindHash(a.Hash(), longSecret); err == nil {
		t.Fatalf("BlindHash(long secret) expected error")
	}
}
