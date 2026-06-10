// Copyright (c) 2026 TriliTech Ltd.
// Author: tzstats@trili.tech

package codec

import (
	"bytes"
	"testing"

	"github.com/trilitech/tzgo/tezos"
)

// TestTransactionWithTz5Destination verifies that operations referencing a tz5
// (ML-DSA-44, v025 Ushuaia) address as destination round-trip through binary
// encode/decode without loss (PKH-only support, BIN-923).
func TestTransactionWithTz5Destination(t *testing.T) {
	src := tezos.MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q")
	dst := tezos.MustParseAddress("tz5VWE3unqGsLVrYhxCGBxiVVYXDjHHbmbTY")
	tx := &Transaction{
		Manager:     Manager{Source: src},
		Destination: dst,
		Amount:      tezos.N(1000),
	}

	buf := bytes.NewBuffer(nil)
	if err := tx.EncodeBuffer(buf, tezos.DefaultParams); err != nil {
		t.Fatalf("encode tx with tz5 destination: %v", err)
	}
	enc := append([]byte(nil), buf.Bytes()...)

	tx2 := &Transaction{}
	if err := tx2.DecodeBuffer(bytes.NewBuffer(enc), tezos.DefaultParams); err != nil {
		t.Fatalf("decode tx with tz5 destination: %v", err)
	}
	if got, want := tx2.Destination.String(), dst.String(); got != want {
		t.Errorf("destination = %s, want %s", got, want)
	}

	buf2 := bytes.NewBuffer(nil)
	if err := tx2.EncodeBuffer(buf2, tezos.DefaultParams); err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	if !bytes.Equal(enc, buf2.Bytes()) {
		t.Errorf("binary round-trip mismatch:\n enc=%x\n got=%x", enc, buf2.Bytes())
	}
}
