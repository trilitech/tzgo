// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

package rpc

// Ensure DoubleConsensusOperationEvidence implements the TypedOperation interface.
var _ TypedOperation = (*DoubleConsensusOperationEvidence)(nil)

// DoubleConsensusOperationEvidence represents a double_consensus_operation_evidence operation
type DoubleConsensusOperationEvidence struct {
	Generic
	Slot int64              `json:"slot"`
	Op1  InlinedEndorsement `json:"op1"`
	Op2  InlinedEndorsement `json:"op2"`
}
