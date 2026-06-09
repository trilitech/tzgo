// Copyright (c) 2026 TriliTech Ltd.
// Author: tzstats@trili.tech

package contract

import (
	"testing"

	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/tezos"
)

// TestFA21TransferRoundTrip verifies that an FA2.1 (sTEZ) `transfer` parameter is
// produced and decoded entirely through TzGo's generic Micheline layer. FA2.1
// keeps the FA2/TZIP-12 `transfer` type unchanged, so the v025 (Ushuaia) sTEZ
// native contract requires no SDK-specific decoder. This is the regression test
// backing that conclusion (BIN-921).
func TestFA21TransferRoundTrip(t *testing.T) {
	from := tezos.MustParseAddress("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx")
	to := tezos.MustParseAddress("tz1burnburnburnburnburnburnburjAYjjX")

	var tokenId, amount tezos.Z
	tokenId.SetInt64(0) // sTEZ is a single-asset FA2.1 token (token_id 0)
	amount.SetInt64(1_000_000)

	// build the transfer parameters via the generic FA2 helper
	args := NewFA2TransferArgs().
		WithTransfer(from, to, tokenId, amount).
		WithSource(from).
		WithDestination(to)
	params := args.Parameters()

	if params.Entrypoint != "transfer" {
		t.Fatalf("entrypoint = %q, want transfer", params.Entrypoint)
	}

	// decode the built value back through the generic TZIP-12 type
	typ := micheline.ITzip12.TypeOf("transfer")
	val := micheline.NewValue(typ, params.Value)

	var decoded FA2TransferList
	if err := val.Unmarshal(&decoded); err != nil {
		t.Fatalf("decode transfer via generic micheline layer: %v", err)
	}

	if len(decoded) != 1 {
		t.Fatalf("decoded %d transfers, want 1", len(decoded))
	}
	got := decoded[0]
	if !got.From.Equal(from) {
		t.Errorf("from = %s, want %s", got.From, from)
	}
	if !got.To.Equal(to) {
		t.Errorf("to = %s, want %s", got.To, to)
	}
	if got.Amount.Int64() != 1_000_000 {
		t.Errorf("amount = %s, want 1000000", got.Amount.String())
	}
	if got.TokenId.Int64() != 0 {
		t.Errorf("token_id = %s, want 0", got.TokenId.String())
	}
}
