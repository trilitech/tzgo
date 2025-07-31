// Copyright (c) 2025 Trilitech Ltd.
// Author: tzstats@trili.tech

package rpc

// Ensure DoubleConsensusOperationEvidence implements the TypedOperation interface.
var _ TypedOperation = (*DoubleConsensusOperationEvidence)(nil)

// DoubleConsensusOperationEvidence represents a double_consensus_operation_evidence operation
type DoubleConsensusOperationEvidence struct {
	Generic
	Slot int64              `json:"slot"`
	OP1  InlinedEndorsement `json:"op1"`
	OP2  InlinedEndorsement `json:"op2"`
}
