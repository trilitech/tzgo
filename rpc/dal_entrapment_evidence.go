// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

package rpc

import "encoding/json"

// Ensure DalEntrapmentEvidence implements the TypedOperation interface.
var _ TypedOperation = (*DalEntrapmentEvidence)(nil)

type RawShard struct {
	json.RawMessage
}

func (s RawShard) AsInt() (int, error) {
	var value int
	if s.RawMessage == nil {
		return value, nil
	}
	err := json.Unmarshal(s.RawMessage, &value)
	return value, err
}

func (s RawShard) AsStringSlice() ([]string, error) {
	var value []string
	if s.RawMessage == nil {
		return value, nil
	}
	err := json.Unmarshal(s.RawMessage, &value)
	return value, err
}

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
}
