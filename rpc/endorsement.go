// Copyright (c) 2020-2024 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"github.com/trilitech/tzgo/tezos"
)

// Ensure Endorsement implements the TypedOperation interface.
var _ TypedOperation = (*Endorsement)(nil)

// Endorsement represents an endorsement operation
type Endorsement struct {
	Generic
	Level       int64               `json:"level"`                 // <= v008, v012+
	Endorsement *InlinedEndorsement `json:"endorsement,omitempty"` // v009+
	Slot        int                 `json:"slot"`                  // v009+
	Round       int                 `json:"round"`                 // v012+
	PayloadHash tezos.PayloadHash   `json:"block_payload_hash"`    // v012+
	// DalAttestation is a raw bitset of attested DAL slots. v019+.
	//
	// BREAKING in v025 (Ushuaia): the bit semantics changed (baker-attested vs
	// protocol-attested slots) and a multi-lag layout was introduced. The value
	// still decodes as a Z, but the bit layout is protocol-dependent: for
	// v024 and earlier, bit i means slot i is attested; from v025 on, the
	// layout packs attested slots per attestation lag and individual bits
	// MUST NOT be interpreted without protocol context. TzGo deliberately
	// keeps this field opaque. To interpret a v025+ bitset, use the node's
	// helper RPCs GET .../helpers/decode_dal_attestation/<bitset> and
	// POST .../helpers/encode_dal_attestation, which translate between the
	// bitset and an explicit list of attested slot indices per lag using the
	// protocol parameters active at the queried block.
	DalAttestation tezos.Z `json:"dal_attestation"` // v019+
}

func (e Endorsement) GetLevel() int64 {
	if e.Endorsement != nil {
		return e.Endorsement.Operations.Level
	}
	return e.Level
}

// InlinedEndorsement represents and embedded endorsement
type InlinedEndorsement struct {
	Branch     tezos.BlockHash `json:"branch"`     // the double block
	Operations Endorsement     `json:"operations"` // only level and kind are set
	Signature  tezos.Signature `json:"signature"`
}
