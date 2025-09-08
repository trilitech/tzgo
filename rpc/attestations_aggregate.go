package rpc

import (
	"github.com/trilitech/tzgo/tezos"
)

// Ensure AttestationsAggregate implements the TypedOperation interface.
var _ TypedOperation = (*AttestationsAggregate)(nil)

// AttestationsAggregate represents an attestations aggregate operation.
type AttestationsAggregate struct {
	Generic
	ConsensusContent ConsensusContent `json:"consensus_content"` // v023+
	Committee        []Committee      `json:"committee"`         // v023+
}

type ConsensusContent struct {
	Level       int64             `json:"level"`              // v023+
	Round       int               `json:"round"`              // v023+
	PayloadHash tezos.PayloadHash `json:"block_payload_hash"` // v023+
}

type Committee struct {
	Slot           int     `json:"slot"`            // v023+
	DalAttestation tezos.Z `json:"dal_attestation"` // v023+
}

// Ensure PreattestationsAggregate implements the TypedOperation interface.
var _ TypedOperation = (*PreattestationsAggregate)(nil)

// PreattestationsAggregate represents an pre-attestations aggregate operation.
type PreattestationsAggregate struct {
	Generic
	ConsensusContent ConsensusContent `json:"consensus_content"` // v023+
	Committee        []int            `json:"committee"`         // v023+
}
