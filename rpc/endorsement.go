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
	Level          int64               `json:"level"`                 // <= v008, v012+
	Endorsement    *InlinedEndorsement `json:"endorsement,omitempty"` // v009+
	Slot           int                 `json:"slot"`                  // v009+
	Round          int                 `json:"round"`                 // v012+
	PayloadHash    tezos.PayloadHash   `json:"block_payload_hash"`    // v012+
	// DalAttestation is a raw bitset of attested DAL slots. v019+.
	// BREAKING in v025 (Ushuaia): the bit semantics changed (baker-attested vs
	// protocol-attested slots) and a multi-lag layout was introduced. The value
	// still decodes as a Z, but code that interprets individual bits must account
	// for the new layout for v025+ blocks.
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
