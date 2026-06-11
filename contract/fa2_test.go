// Copyright (c) 2026 TriliTech Ltd.
// Author: tzstats@trili.tech

package contract

import (
	"encoding/json"
	"testing"

	"github.com/trilitech/tzgo/micheline"
	"github.com/trilitech/tzgo/tezos"
)

// stezTransferType is the parameter type of the v025 (Ushuaia) sTEZ native
// contract's %transfer entrypoint, transcribed from the protocol source
// (octez src/proto_025_PsUshuai/lib_protocol/script_native_types.ml, type
// `transfer` at line 262, `transfer_type` at line 403):
//
//	list (pair (address %from_)
//	           (list %txs (pair (address %to_) (nat %token_id) (nat %amount))))
//
// Note the inner element is a tup3 (right-comb 3-tuple), whereas classic
// FA2/TZIP-12 interfaces commonly write nested binary pairs — the same
// Michelson type by comb-pair equivalence. No on-chain sTEZ value exists yet
// on any network (the stez feature flag is disabled even on ushuaianet), so
// the protocol source is the authoritative fixture.
const stezTransferType = `{
	"prim": "list", "annots": ["%transfer"], "args": [
		{"prim": "pair", "args": [
			{"prim": "address", "annots": ["%from_"]},
			{"prim": "list", "annots": ["%txs"], "args": [
				{"prim": "pair", "args": [
					{"prim": "address", "annots": ["%to_"]},
					{"prim": "nat", "annots": ["%token_id"]},
					{"prim": "nat", "annots": ["%amount"]}
				]}
			]}
		]}
	]
}`

// TestStezTransferDecode verifies that transfer values shaped like the v025
// sTEZ native contract emits them decode through TzGo's generic Micheline
// layer into the existing FA2 helper types, across the value encodings a node
// may produce (nested binary pairs and n-ary comb pairs). Backs BIN-921.
func TestStezTransferDecode(t *testing.T) {
	var typPrim micheline.Prim
	if err := json.Unmarshal([]byte(stezTransferType), &typPrim); err != nil {
		t.Fatalf("parse sTEZ transfer type: %v", err)
	}
	typ := micheline.NewType(typPrim)

	from := "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"
	to := "tz1burnburnburnburnburnburnburjAYjjX"

	type wantTx struct {
		from, to string
		tokenId  int64
		amount   int64
	}
	cases := []struct {
		name  string
		value string
		want  []wantTx
	}{
		{
			// classic nested binary pairs in readable-mode JSON
			name: "nested-pairs",
			value: `[{"prim":"Pair","args":[{"string":"` + from + `"},[
				{"prim":"Pair","args":[{"string":"` + to + `"},{"prim":"Pair","args":[{"int":"0"},{"int":"1000000"}]}]}
			]]}]`,
			want: []wantTx{{from, to, 0, 1000000}},
		},
		{
			// n-ary comb pair for the tup3, matching the protocol's type shape
			name: "comb-pairs",
			value: `[{"prim":"Pair","args":[{"string":"` + from + `"},[
				{"prim":"Pair","args":[{"string":"` + to + `"},{"int":"0"},{"int":"1000000"}]}
			]]}]`,
			want: []wantTx{{from, to, 0, 1000000}},
		},
		{
			// multiple transfer groups and txs
			name: "multi-tx",
			value: `[
				{"prim":"Pair","args":[{"string":"` + from + `"},[
					{"prim":"Pair","args":[{"string":"` + to + `"},{"int":"0"},{"int":"1"}]},
					{"prim":"Pair","args":[{"string":"` + from + `"},{"int":"0"},{"int":"2"}]}
				]]},
				{"prim":"Pair","args":[{"string":"` + to + `"},[
					{"prim":"Pair","args":[{"string":"` + from + `"},{"int":"0"},{"int":"3"}]}
				]]}
			]`,
			want: []wantTx{
				{from, to, 0, 1},
				{from, from, 0, 2},
				{to, from, 0, 3},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var valPrim micheline.Prim
			if err := json.Unmarshal([]byte(c.value), &valPrim); err != nil {
				t.Fatalf("parse value: %v", err)
			}
			val := micheline.NewValue(typ, valPrim)

			var decoded FA2TransferList
			if err := val.Unmarshal(&decoded); err != nil {
				t.Fatalf("decode sTEZ transfer via generic micheline layer: %v", err)
			}
			if len(decoded) != len(c.want) {
				t.Fatalf("decoded %d transfers, want %d", len(decoded), len(c.want))
			}
			for i, w := range c.want {
				got := decoded[i]
				if got.From.String() != w.from {
					t.Errorf("tx %d: from = %s, want %s", i, got.From, w.from)
				}
				if got.To.String() != w.to {
					t.Errorf("tx %d: to = %s, want %s", i, got.To, w.to)
				}
				if got.TokenId.Int64() != w.tokenId {
					t.Errorf("tx %d: token_id = %d, want %d", i, got.TokenId.Int64(), w.tokenId)
				}
				if got.Amount.Int64() != w.amount {
					t.Errorf("tx %d: amount = %d, want %d", i, got.Amount.Int64(), w.amount)
				}
			}
		})
	}
}

// TestStezTransferTypeMatchesTzip12 asserts the sTEZ transfer parameter type
// (from the protocol source) unfolds to the same type definition as the
// classic FA2/TZIP-12 transfer type TzGo ships in micheline.ITzip12 — the
// basis for reusing the existing FA2 helpers for sTEZ.
func TestStezTransferTypeMatchesTzip12(t *testing.T) {
	var stez micheline.Prim
	if err := json.Unmarshal([]byte(stezTransferType), &stez); err != nil {
		t.Fatalf("parse sTEZ transfer type: %v", err)
	}
	stezTd := micheline.NewType(stez).Typedef("transfer")
	tzipTd := micheline.ITzip12.TypeOf("transfer").Typedef("transfer")

	// Compare modulo prim paths: the sTEZ type writes the inner triple as a
	// tup3 (right-comb) while TZIP-12 writes nested binary pairs, so element
	// paths differ ([...,2] vs [...,1,1]) while field names and types are
	// identical. Value decoding across both layouts is covered by
	// TestStezTransferDecode.
	stezJSON := marshalTypedefNoPaths(t, stezTd)
	tzipJSON := marshalTypedefNoPaths(t, tzipTd)
	if stezJSON != tzipJSON {
		t.Errorf("sTEZ transfer typedef differs from TZIP-12:\n stez=%s\n tzip=%s", stezJSON, tzipJSON)
	}
}

// marshalTypedefNoPaths renders a typedef as JSON with all prim paths removed,
// leaving only field names and types for structural comparison.
func marshalTypedefNoPaths(t *testing.T, td micheline.Typedef) string {
	t.Helper()
	buf, err := json.Marshal(td)
	if err != nil {
		t.Fatalf("marshal typedef: %v", err)
	}
	var v any
	if err := json.Unmarshal(buf, &v); err != nil {
		t.Fatalf("unmarshal typedef: %v", err)
	}
	var strip func(any)
	strip = func(node any) {
		switch n := node.(type) {
		case map[string]any:
			delete(n, "path")
			for _, c := range n {
				strip(c)
			}
		case []any:
			for _, c := range n {
				strip(c)
			}
		}
	}
	strip(v)
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("re-marshal typedef: %v", err)
	}
	return string(out)
}

// TestFA21TransferRoundTrip verifies that a transfer built with TzGo's generic
// FA2 helper decodes back through the TZIP-12 type — the encode side of the
// sTEZ compatibility claim (recipients can be paid via the existing helper).
func TestFA21TransferRoundTrip(t *testing.T) {
	from := tezos.MustParseAddress("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx")
	to := tezos.MustParseAddress("tz1burnburnburnburnburnburnburjAYjjX")

	var tokenId, amount tezos.Z
	tokenId.SetInt64(0) // sTEZ is a single-asset FA2.1 token (token_id 0)
	amount.SetInt64(1_000_000)

	args := NewFA2TransferArgs().
		WithTransfer(from, to, tokenId, amount).
		WithSource(from).
		WithDestination(to)
	params := args.Parameters()

	if params.Entrypoint != "transfer" {
		t.Fatalf("entrypoint = %q, want transfer", params.Entrypoint)
	}

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
	if !got.From.Equal(from) || !got.To.Equal(to) {
		t.Errorf("addresses mismatch: from=%s to=%s", got.From, got.To)
	}
	if got.Amount.Int64() != 1_000_000 || got.TokenId.Int64() != 0 {
		t.Errorf("amount/token_id mismatch: %d/%d", got.Amount.Int64(), got.TokenId.Int64())
	}
}
