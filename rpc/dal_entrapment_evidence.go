// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

package rpc

import "encoding/json"

// Ensure DalEntrapmentEvidence implements the TypedOperation interface.
var _ TypedOperation = (*DalEntrapmentEvidence)(nil)

// RawShard holds one element of a shard_with_proof.shard array, which the
// node encodes either as an integer (the shard index) or as a list of hex
// strings (the shard data). Use AsInt or AsStringSlice to interpret it.
type RawShard struct {
	json.RawMessage
}

// AsInt interprets the raw shard element as an integer shard index.
func (s RawShard) AsInt() (int, error) {
	var value int
	if s.RawMessage == nil {
		return value, nil
	}
	err := json.Unmarshal(s.RawMessage, &value)
	return value, err
}

// AsStringSlice interprets the raw shard element as a list of hex strings.
func (s RawShard) AsStringSlice() ([]string, error) {
	var value []string
	if s.RawMessage == nil {
		return value, nil
	}
	err := json.Unmarshal(s.RawMessage, &value)
	return value, err
}

// ShardWithProof is a DAL shard together with its inclusion proof as carried
// by dal_entrapment_evidence operations.
type ShardWithProof struct {
	// `shard` can be either an integer or a list of strings
	// https://gitlab.com/tezos/tezos/-/blob/2c16b5170ad3305538cba6cbc636ef7560531f05/docs/api/seoul-openapi.json#L15621
	Shard []RawShard `json:"shard"`
	Proof string     `json:"proof"`
}

// DalEntrapmentEvidence represents a dal_entrapment_evidence operation
type DalEntrapmentEvidence struct {
	Generic
	Attestation    InlinedEndorsement `json:"attestation"`
	ConsensusSlot  uint16             `json:"consensus_slot"`
	SlotIndex      uint8              `json:"slot_index"`
	ShardWithProof ShardWithProof     `json:"shard_with_proof"`
	// LagIndex identifies which attestation lag the evidence refers to: an
	// index into the attestation_lags protocol constant (see the
	// dal_parametric block of the chain constants, surfaced as
	// tezos.Params.DalAttestationLags). Added in v025 (Ushuaia) together with
	// dynamic attestation lag; nil for pre-v025 blocks. Note that a present
	// zero value ("lag_index": 0) is distinct from an absent field.
	LagIndex *uint8 `json:"lag_index,omitempty"`
}
