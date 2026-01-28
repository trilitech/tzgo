// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech
//
// Integration tests for RPC client endpoints.
// These tests connect to a real RPC endpoint (Ghostnet) and validate
// the structure and types of returned payloads.

package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/tezos"
)

// public Mainnet RPC endpoint
var testRPCURL = "https://mainnet.smartpy.io/"

// TestBlockOperations_Integration tests GetBlockOperations endpoint
// and validates operation structure and types.
func TestBlockOperations_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	// Use a known block hash from Ghostnet
	blockID := tezos.MustParseBlockHash("BLuhVLfNYkkUrNeypkYp5u9Rpz47tc7xoggkCUn3rwE9KL9DCVZ")

	ops, err := c.GetBlockOperations(ctx, blockID)
	require.NoError(t, err)
	require.NotNil(t, ops)
	assert.Greater(t, len(ops), 0, "should have at least one operation list")

	// Test first operation list (validation pass 0)
	if len(ops) > 0 && len(ops[0]) > 0 {
		op := ops[0][0]
		assert.NotEmpty(t, op.Hash, "operation hash should not be empty")
		assert.NotEmpty(t, op.Branch, "operation branch should not be empty")
		assert.NotEmpty(t, op.Protocol, "operation protocol should not be empty")
		assert.NotEmpty(t, op.ChainID, "operation chain ID should not be empty")
		assert.NotEmpty(t, op.Signature, "operation signature should not be empty")
		assert.Greater(t, len(op.Contents), 0, "operation should have contents")

		// Test accessing operation contents
		firstContent := op.Contents.N(0)
		require.NotNil(t, firstContent, "first content should not be nil")

		// Validate operation kind
		kind := firstContent.Kind()
		assert.NotEqual(t, tezos.OpTypeInvalid, kind, "operation kind should be valid")

		// Validate metadata
		meta := firstContent.Meta()
		assert.NotNil(t, meta, "metadata should not be nil")

		// Test specific operation types
		switch typedOp := firstContent.(type) {
		case *Endorsement:
			validateEndorsement(t, typedOp, meta)
		case *Transaction:
			validateTransaction(t, typedOp, meta)
		case *Delegation:
			validateDelegation(t, typedOp, meta)
		case *Origination:
			validateOrigination(t, typedOp, meta)
		case *Reveal:
			validateReveal(t, typedOp, meta)
		}
	}
}

// TestBlock_Integration tests block-related endpoints.
func TestBlock_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	blockID := tezos.MustParseBlockHash("BLuhVLfNYkkUrNeypkYp5u9Rpz47tc7xoggkCUn3rwE9KL9DCVZ")

	// Test GetBlock
	block, err := c.GetBlock(ctx, blockID)
	require.NoError(t, err)
	require.NotNil(t, block)
	validateBlock(t, block)

	// Test GetBlockHeader
	header, err := c.GetBlockHeader(ctx, blockID)
	require.NoError(t, err)
	require.NotNil(t, header)
	validateBlockHeader(t, header)

	// Test GetBlockMetadata
	metadata, err := c.GetBlockMetadata(ctx, blockID)
	require.NoError(t, err)
	require.NotNil(t, metadata)
	validateBlockMetadata(t, metadata)

	// Test GetBlockHash
	hash, err := c.GetBlockHash(ctx, blockID)
	require.NoError(t, err)
	assert.Equal(t, blockID, hash, "block hash should match")

	// Test GetHeadBlock
	head, err := c.GetHeadBlock(ctx)
	require.NoError(t, err)
	require.NotNil(t, head)
	assert.Greater(t, head.Header.Level, int64(0), "head block level should be positive")
}

// TestBlockOperations_OperationTypes tests various operation types.
func TestBlockOperations_OperationTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	blockID := tezos.MustParseBlockHash("BLuhVLfNYkkUrNeypkYp5u9Rpz47tc7xoggkCUn3rwE9KL9DCVZ")

	ops, err := c.GetBlockOperations(ctx, blockID)
	require.NoError(t, err)

	// Collect operation types found
	opTypes := make(map[tezos.OpType]bool)
	for _, opList := range ops {
		for _, op := range opList {
			for _, content := range op.Contents {
				opTypes[content.Kind()] = true
			}
		}
	}

	// Validate that we found at least some operations
	assert.Greater(t, len(opTypes), 0, "should find at least one operation type")

	// Test individual operation access
	for _, opList := range ops {
		for _, op := range opList {
			for i := 0; i < op.Contents.Len(); i++ {
				content := op.Contents.N(i)
				require.NotNil(t, content, "content should not be nil")

				// Validate TypedOperation interface
				kind := content.Kind()
				assert.NotEqual(t, tezos.OpTypeInvalid, kind, "operation kind should be valid")

				meta := content.Meta()
				assert.NotNil(t, meta, "metadata should not be nil")

				result := content.Result()
				assert.NotNil(t, result, "result should not be nil")

				costs := content.Costs()
				assert.NotNil(t, costs, "costs should not be nil")

				limits := content.Limits()
				assert.NotNil(t, limits, "limits should not be nil")
			}
		}
	}
}

// TestContract_Integration tests contract-related endpoints.
func TestContract_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	// Use a known contract address (you may need to adjust this for Ghostnet)
	contractAddr := tezos.MustParseAddress("KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd")

	// Test GetContract
	contract, err := c.GetContract(ctx, contractAddr, Head)
	require.NoError(t, err)
	require.NotNil(t, contract)
	validateContractInfo(t, contract)

	// Test GetContractBalance
	balance, err := c.GetContractBalance(ctx, contractAddr, Head)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, balance.Int64(), int64(0), "balance should be non-negative")

	// Test GetContractStorage
	storage, err := c.GetContractStorage(ctx, contractAddr, Head)
	require.NoError(t, err)
	assert.NotNil(t, storage, "storage should not be nil")

	// Test GetContractScript
	script, err := c.GetContractScript(ctx, contractAddr)
	if err == nil {
		require.NotNil(t, script)
		assert.NotNil(t, script.Code, "script code should not be nil")
		assert.NotNil(t, script.Storage, "script storage should not be nil")
	}

	// Test GetContractEntrypoints
	entrypoints, err := c.GetContractEntrypoints(ctx, contractAddr)
	if err == nil {
		assert.NotNil(t, entrypoints, "entrypoints should not be nil")
	}
}

// TestChain_Integration tests chain-related endpoints.
func TestChain_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test GetChainId
	chainID, err := c.GetChainId(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, chainID, "chain ID should not be empty")

	// Test GetStatus
	status, err := c.GetStatus(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, status.SyncState, "sync state should not be empty")
}

// TestConstants_Integration tests protocol constants endpoint.
func TestConstants_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	// Test GetParams
	params, err := c.GetParams(ctx, Head)
	require.NoError(t, err)
	require.NotNil(t, params)
	validateParams(t, params)

	// Test GetConstants
	constants, err := c.GetConstants(ctx, Head)
	require.NoError(t, err)
	require.NotNil(t, constants)
	validateConstants(t, constants)
}

// TestRights_Integration tests baking and endorsing rights endpoints.
func TestRights_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	// Get current block to determine cycle
	head, err := c.GetHeadBlock(ctx)
	require.NoError(t, err)
	currentCycle := head.GetCycle()

	// Test ListBakingRights
	bakingRights, err := c.ListBakingRights(ctx, Head, 10)
	if err == nil {
		assert.NotNil(t, bakingRights, "baking rights should not be nil")
		if len(bakingRights) > 0 {
			validateBakingRight(t, &bakingRights[0])
		}
	}

	// Test ListBakingRightsCycle
	bakingRightsCycle, err := c.ListBakingRightsCycle(ctx, Head, currentCycle, 10)
	if err == nil {
		assert.NotNil(t, bakingRightsCycle, "baking rights cycle should not be nil")
	}

	// Test ListEndorsingRights
	endorsingRights, err := c.ListEndorsingRights(ctx, Head)
	if err == nil {
		assert.NotNil(t, endorsingRights, "endorsing rights should not be nil")
		if len(endorsingRights) > 0 {
			validateEndorsingRight(t, &endorsingRights[0])
		}
	}

	// Test ListEndorsingRightsCycle
	endorsingRightsCycle, err := c.ListEndorsingRightsCycle(ctx, Head, currentCycle)
	if err == nil {
		assert.NotNil(t, endorsingRightsCycle, "endorsing rights cycle should not be nil")
	}
}

// TestMempool_Integration tests mempool endpoint.
func TestMempool_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test GetMempool
	mempool, err := c.GetMempool(ctx)
	if err == nil {
		require.NotNil(t, mempool)
		assert.NotNil(t, mempool.Applied, "applied operations should not be nil")
		assert.NotNil(t, mempool.Refused, "refused operations should not be nil")
		assert.NotNil(t, mempool.BranchRefused, "branch refused operations should not be nil")
		assert.NotNil(t, mempool.BranchDelayed, "branch delayed operations should not be nil")
	}
}

// TestDelegates_Integration tests delegate-related endpoints.
func TestDelegates_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	// Get a delegate address from the current block
	head, err := c.GetHeadBlock(ctx)
	require.NoError(t, err)
	baker := head.Metadata.Baker

	// Test GetDelegate
	delegate, err := c.GetDelegate(ctx, baker, Head)
	if err == nil {
		require.NotNil(t, delegate)
		assert.NotEmpty(t, delegate.Delegate, "delegate address should not be empty")
		assert.GreaterOrEqual(t, delegate.Balance, int64(0), "delegate balance should be non-negative")
	}

	// Test GetDelegateBalance
	balance, err := c.GetDelegateBalance(ctx, baker, Head)
	if err == nil {
		assert.GreaterOrEqual(t, balance, int64(0), "delegate balance should be non-negative")
	}

	// Test GetDelegateKey
	key, err := c.GetDelegateKey(ctx, baker, Head)
	if err == nil {
		assert.NotEmpty(t, key, "delegate key should not be empty")
	}

	// Test ListActiveDelegates
	delegates, err := c.ListActiveDelegates(ctx, Head)
	if err == nil {
		assert.NotNil(t, delegates, "delegates list should not be nil")
		if len(delegates) > 0 {
			assert.NotEmpty(t, delegates[0], "delegate address should not be empty")
		}
	}
}

// TestVotes_Integration tests voting-related endpoints.
func TestVotes_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	// Test ListVoters
	voters, err := c.ListVoters(ctx, Head)
	if err == nil {
		assert.NotNil(t, voters, "voters should not be nil")
	}

	// Test ListProposals
	proposals, err := c.ListProposals(ctx, Head)
	if err == nil {
		assert.NotNil(t, proposals, "proposals should not be nil")
	}

	// Test ListBallots
	ballots, err := c.ListBallots(ctx, Head)
	if err == nil {
		assert.NotNil(t, ballots, "ballots should not be nil")
	}

	// Test GetVoteQuorum
	quorum, err := c.GetVoteQuorum(ctx, Head)
	if err == nil {
		assert.GreaterOrEqual(t, quorum, 1, "quorum should be positive")
	}
}

// TestOperationHashes_Integration tests operation hash endpoints.
func TestOperationHashes_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	blockID := tezos.MustParseBlockHash("BLuhVLfNYkkUrNeypkYp5u9Rpz47tc7xoggkCUn3rwE9KL9DCVZ")

	// Test GetBlockOperationHashes
	hashes, err := c.GetBlockOperationHashes(ctx, blockID)
	require.NoError(t, err)
	assert.Len(t, hashes, 4, "should have 4 operation lists")
	if len(hashes) > 0 && len(hashes[0]) > 0 {
		assert.NotEmpty(t, hashes[0][0], "operation hash should not be empty")
	}

	// Test GetBlockOperationListHashes
	if len(hashes) > 0 {
		listHashes, err := c.GetBlockOperationListHashes(ctx, blockID, 0)
		if err == nil {
			assert.NotNil(t, listHashes, "list hashes should not be nil")
			if len(listHashes) > 0 {
				assert.NotEmpty(t, listHashes[0], "operation hash should not be empty")
			}
		}
	}

	// Test GetBlockOperationHash
	if len(hashes) > 0 && len(hashes[0]) > 0 {
		hash, err := c.GetBlockOperationHash(ctx, blockID, 0, 0)
		if err == nil {
			assert.NotEmpty(t, hash, "operation hash should not be empty")
			assert.Equal(t, hashes[0][0], hash, "hash should match")
		}
	}
}

// TestOperationList_Integration tests operation list endpoints.
func TestOperationList_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	blockID := tezos.MustParseBlockHash("BLuhVLfNYkkUrNeypkYp5u9Rpz47tc7xoggkCUn3rwE9KL9DCVZ")

	// Test GetBlockOperationList
	ops, err := c.GetBlockOperationList(ctx, blockID, 0)
	if err == nil {
		assert.NotNil(t, ops, "operations should not be nil")
		if len(ops) > 0 {
			assert.NotEmpty(t, ops[0].Hash, "operation hash should not be empty")
			assert.Greater(t, len(ops[0].Contents), 0, "operation should have contents")
		}
	}

	// Test GetBlockOperation
	if len(ops) > 0 {
		op, err := c.GetBlockOperation(ctx, blockID, 0, 0)
		if err == nil {
			require.NotNil(t, op)
			assert.NotEmpty(t, op.Hash, "operation hash should not be empty")
			assert.Greater(t, len(op.Contents), 0, "operation should have contents")
		}
	}
}

// TestIssuance_Integration tests issuance endpoint.
func TestIssuance_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := NewClient(testRPCURL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.Init(ctx)
	require.NoError(t, err)

	// Test GetIssuance
	issuance, err := c.GetIssuance(ctx, Head)
	if err == nil {
		assert.NotNil(t, issuance, "issuance should not be nil")
		if len(issuance) > 0 {
			assert.Greater(t, issuance[0].Cycle, int64(0), "issuance cycle should be positive")
		}
	}
}

// Validation helper functions

func validateBlock(t *testing.T, block *Block) {
	t.Helper()
	assert.NotEmpty(t, block.Hash, "block hash should not be empty")
	assert.NotEmpty(t, block.Protocol, "block protocol should not be empty")
	assert.NotEmpty(t, block.ChainId, "block chain ID should not be empty")
	assert.GreaterOrEqual(t, block.Header.Level, int64(0), "block level should be non-negative")
	assert.NotEmpty(t, block.Header.Predecessor, "predecessor should not be empty")
	assert.NotZero(t, block.Header.Timestamp, "timestamp should not be zero")
	assert.Len(t, block.Operations, 4, "should have 4 operation lists")
}

func validateBlockHeader(t *testing.T, header *BlockHeader) {
	t.Helper()
	assert.GreaterOrEqual(t, header.Level, int64(0), "header level should be non-negative")
	assert.NotEmpty(t, header.Predecessor, "predecessor should not be empty")
	assert.NotZero(t, header.Timestamp, "timestamp should not be zero")
	assert.GreaterOrEqual(t, header.Proto, 0, "proto version should be non-negative")
	assert.NotEmpty(t, header.Context, "context hash should not be empty")
	assert.NotEmpty(t, header.Signature, "signature should not be empty")
}

func validateBlockMetadata(t *testing.T, metadata *BlockMetadata) {
	t.Helper()
	assert.NotEmpty(t, metadata.Protocol, "protocol should not be empty")
	assert.NotEmpty(t, metadata.NextProtocol, "next protocol should not be empty")
	assert.Greater(t, metadata.MaxOperationsTTL, 0, "max operations TTL should be positive")
	assert.NotEmpty(t, metadata.Baker, "baker should not be empty")
}

func validateEndorsement(t *testing.T, op *Endorsement, meta OperationMetadata) {
	t.Helper()
	assert.Greater(t, op.GetLevel(), int64(0), "endorsement level should be positive")
	assert.NotEmpty(t, meta.Delegate, "endorsement delegate should not be empty")
	assert.NotNil(t, meta.TotalConsensusPower, "endorsement slots should not be nil")
}

func validateTransaction(t *testing.T, op *Transaction, _ OperationMetadata) {
	t.Helper()
	assert.NotEmpty(t, op.Source, "transaction source should not be empty")
	assert.NotEmpty(t, op.Destination, "transaction destination should not be empty")
	assert.GreaterOrEqual(t, op.Amount, int64(0), "transaction amount should be non-negative")
	assert.GreaterOrEqual(t, op.Fee, int64(0), "transaction fee should be non-negative")
	assert.Greater(t, op.Counter, int64(0), "transaction counter should be positive")
	assert.Greater(t, op.GasLimit, int64(0), "gas limit should be positive")
}

func validateDelegation(t *testing.T, op *Delegation, _ OperationMetadata) {
	t.Helper()
	assert.NotEmpty(t, op.Source, "delegation source should not be empty")
	assert.GreaterOrEqual(t, op.Fee, int64(0), "delegation fee should be non-negative")
}

func validateOrigination(t *testing.T, op *Origination, _ OperationMetadata) {
	t.Helper()
	assert.NotEmpty(t, op.Source, "origination source should not be empty")
	assert.GreaterOrEqual(t, op.Fee, int64(0), "origination fee should be non-negative")
}

func validateReveal(t *testing.T, op *Reveal, _ OperationMetadata) {
	t.Helper()
	assert.NotEmpty(t, op.Source, "reveal source should not be empty")
	assert.NotEmpty(t, op.PublicKey, "reveal public key should not be empty")
	assert.GreaterOrEqual(t, op.Fee, int64(0), "reveal fee should be non-negative")
}

func validateContractInfo(t *testing.T, info *ContractInfo) {
	t.Helper()
	assert.GreaterOrEqual(t, info.Balance, int64(0), "contract balance should be non-negative")
	assert.GreaterOrEqual(t, info.Counter, int64(0), "contract counter should be non-negative")
}

func validateParams(t *testing.T, params *tezos.Params) {
	t.Helper()
	assert.Greater(t, params.BlocksPerCycle, int64(0), "blocks per cycle should be positive")
	assert.Greater(t, params.BlocksPerCommitment, int64(0), "blocks per commitment should be positive")
	assert.Greater(t, params.MinimalBlockDelay.Nanoseconds(), int64(0), "minimal block delay should be positive")
}

func validateConstants(t *testing.T, constants Constants) {
	t.Helper()
	assert.Greater(t, constants.BlocksPerCycle, int64(0), "blocks per cycle should be positive")
	assert.Greater(t, constants.BlocksPerCommitment, int64(0), "blocks per commitment should be positive")
	assert.Greater(t, constants.NonceRevelationThreshold, int64(0), "nonce revelation threshold should be positive")
	assert.Greater(t, constants.SmartRollupChallengeWindowInBlocks, int64(0), "smart rollup challenge window in blocks should be positive")
	assert.Greater(t, constants.SmartRollupCommitmentPeriodInBlocks, int64(0), "smart rollup commitment period in blocks should be positive")
	assert.Greater(t, constants.SmartRollupMaxLookaheadInBlocks, int64(0), "smart rollup max lookahead in blocks should be positive")
	assert.Greater(t, constants.SmartRollupTimeoutPeriodInBlocks, int64(0), "smart rollup timeout period in blocks should be positive")
	assert.Greater(t, constants.AllBakersAttestActivationThreshold.Num, int(0), "all bakers attest activation threshold numerator should be positive")
	assert.Greater(t, constants.AllBakersAttestActivationThreshold.Den, int(0), "all bakers attest activation threshold denominator should be positive")
}

func validateBakingRight(t *testing.T, right *BakingRight) {
	t.Helper()
	assert.Greater(t, right.Level, int64(0), "baking right level should be positive")
	assert.NotEmpty(t, right.Delegate, "baking right delegate should not be empty")
	assert.GreaterOrEqual(t, right.Priority, 0, "baking right priority should be non-negative")
}

func validateEndorsingRight(t *testing.T, right *EndorsingRight) {
	t.Helper()
	assert.Greater(t, right.Level, int64(0), "endorsing right level should be positive")
	assert.NotEmpty(t, right.Delegate, "endorsing right delegate should not be empty")
	assert.Greater(t, len(right.Slots), 0, "endorsing right should have slots")
}
