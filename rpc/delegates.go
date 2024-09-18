// Copyright (c) 2020-2024 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/trilitech/tzgo/tezos"
)

// v021+
type MisbehaviourKind string
type Vote string

const (
	MisbehaviourAttestation    MisbehaviourKind = "attestation"
	MisbehaviourBlock          MisbehaviourKind = "block"
	MisbehaviourPreattestation MisbehaviourKind = "preattestation"

	VoteNay  Vote = "nay"
	VoteYay  Vote = "yay"
	VotePass Vote = "pass"
)

// Delegate holds information about an active delegate
type Delegate struct {
	// extra info
	Delegate tezos.Address `json:"-"`
	Height   int64         `json:"-"`
	Block    string        `json:"-"`

	// tezos data
	Deactivated bool          `json:"deactivated"`
	GracePeriod int64         `json:"grace_period"`
	VotingPower Int64orString `json:"voting_power"`

	// -v011
	FrozenBalanceByCycle []CycleBalance `json:"frozen_balance_by_cycle"`
	FrozenBalance        int64          `json:"frozen_balance,string"`
	Balance              int64          `json:"balance,string"`

	// -v020
	StakingBalance     int64           `json:"staking_balance,string"`
	DelegatedContracts []tezos.Address `json:"delegated_contracts"`
	DelegatedBalance   int64           `json:"delegated_balance,string"`

	// v012-v020
	FullBalance           int64 `json:"full_balance,string"`
	FrozenDeposits        int64 `json:"frozen_deposits,string"`
	CurrentFrozenDeposits int64 `json:"current_frozen_deposits,string"`
	FrozenDepositsLimit   int64 `json:"frozen_deposits_limit,string"`

	// v015-v020
	ActiveConsensusKey   tezos.Address `json:"active_consensus_key"`
	PendingConsensusKeys []CycleKey    `json:"pending_consensus_keys"`

	// v019+
	MinDelegated struct {
		Amount int64     `json:"amount,string"`
		Level  LevelInfo `json:"level"`
	} `json:"min_delegated_in_current_cycle"`
	StakingDenominator int64 `json:"staking_denominator,string"`

	// v019-v020
	PendingDenunciations bool  `json:"pending_denunciations"`
	TotalDelegatedStake  int64 `json:"total_delegated_stake,string"`

	// v021+
	IsForbidden   bool `json:"is_forbidden"`
	Participation struct {
		ExpectedCycleActivity       int64  `json:"expected_cycle_activity"`
		MinimalCycleActivity        int64  `json:"minimal_cycle_activity"`
		MissedSlots                 int64  `json:"missed_slots"`
		MissedLevels                int64  `json:"missed_levels"`
		RemainingAllowedMissedSlots int64  `json:"remaining_allowed_missed_slots"`
		ExpectedAttestingRewards    uint64 `json:"expected_attesting_rewards,string"`
	} `json:"participation"`
	ActiveStakingParameters  StakingParameters `json:"active_staking_parameters"`
	PendingStakingParameters []struct {
		Cycle      int64             `json:"cycle"`
		Parameters StakingParameters `json:"parameters"`
	} `json:"pending_staking_parameters"`
	BakingPower           int64  `json:"baking_power,string"`
	TotalStaked           uint64 `json:"total_staked,string"`
	TotalDelegated        uint64 `json:"total_delegated,string"`
	OwnFullBalance        uint64 `json:"own_full_balance,string"`
	OwnStaked             uint64 `json:"own_staked,string"`
	OwnDelegated          uint64 `json:"own_delegated,string"`
	ExternalStaked        uint64 `json:"external_staked,string"`
	ExternalDelegated     uint64 `json:"external_delegated,string"`
	TotalUnstakedPerCycle struct {
		Cycle   int64  `json:"cycle"`
		Deposit uint64 `json:"deposit,string"`
	} `json:"total_unstaked_per_cycle"`
	Denunciations []struct {
		OperationHash tezos.OpHash  `json:"operation_hash"`
		Rewarded      tezos.Address `json:"rewarded"`
		Misbehaviour  struct {
			Level uint64           `json:"level"`
			Round int64            `json:"round"`
			Kind  MisbehaviourKind `json:"kind,string"`
		} `json:"misbehaviour"`
	} `json:"denunciations"`
	EstimatedSharedPendingSlashedAmount uint64 `json:"estimated_shared_pending_slashed_amount,string"`
	CurrentVotingPower                  uint64 `json:"current_voting_power,string"`
	VotingInfo                          struct {
		VotingPower        int64                `json:"voting_power"`
		CurrentBallot      Vote                 `json:"current_ballot"`
		CurrentProposals   []tezos.ProtocolHash `json:"current_proposals"`
		RemainingProposals int64                `json:"remaining_proposals"`
	} `json:"voting_info"`
	ConsensusKey struct {
		Active struct {
			Pkh tezos.Address `json:"pkh"`
			Pk  tezos.Key     `json:"pk"`
		} `json:"active"`
		Pendings []CycleKey `json:"pendings"`
	} `json:"consensus_key"`
	Stakers []struct {
		Staker        tezos.Address `json:"staker"`
		FrozenDeposit uint64        `json:"frozen_deposit,string"`
	} `json:"stakers"`
	Delegators []tezos.Address `json:"delegators"`
}

type CycleKey struct {
	Cycle int64         `json:"cycle"`
	Pkh   tezos.Address `json:"pkh"`
	Pk    tezos.Key     `json:"pk"`
}

type CycleBalance struct {
	Cycle   int64 `json:"cycle"`
	Deposit int64 `json:"deposit,string"`
	Fees    int64 `json:"fees,string"`
	Rewards int64 `json:"rewards,string"`
}

// DelegateList contains a list of delegates
type DelegateList []tezos.Address

// ListActiveDelegates returns information about all active delegates at a block.
func (c *Client) ListActiveDelegates(ctx context.Context, id BlockID) (DelegateList, error) {
	p, err := c.GetParams(ctx, id)
	if err != nil {
		return nil, err
	}
	selector := "active=true"
	if p.Version >= 13 {
		selector += "&with_minimal_stake=true"
	}
	delegates := make(DelegateList, 0)
	u := fmt.Sprintf("chains/main/blocks/%s/context/delegates?%s", id, selector)
	if err := c.Get(ctx, u, &delegates); err != nil {
		return nil, err
	}
	return delegates, nil
}

// GetDelegate returns information about a delegate at a specific height.
func (c *Client) GetDelegate(ctx context.Context, addr tezos.Address, id BlockID) (*Delegate, error) {
	delegate := &Delegate{
		Delegate: addr,
		Height:   id.Int64(),
		Block:    id.String(),
	}
	u := fmt.Sprintf("chains/main/blocks/%s/context/delegates/%s", id, addr)
	if err := c.Get(ctx, u, &delegate); err != nil {
		return nil, err
	}
	return delegate, nil
}

// GetDelegateBalance returns a delegate's balance
func (c *Client) GetDelegateBalance(ctx context.Context, addr tezos.Address, id BlockID) (int64, error) {
	u := fmt.Sprintf("chains/main/blocks/%s/context/delegates/%s/balance", id, addr)
	var bal string
	err := c.Get(ctx, u, &bal)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(bal, 10, 64)
}

// GetDelegateKey returns a delegate's current consensus key
func (c *Client) GetDelegateKey(ctx context.Context, addr tezos.Address, id BlockID) (tezos.Key, error) {
	u := fmt.Sprintf("chains/main/blocks/%s/context/delegates/%s/consensus_key", id, addr)
	type ActiveConsensusKey struct {
		Active struct {
			Pk tezos.Key `json:"pk"`
		} `json:"active"`
	}
	var key ActiveConsensusKey
	err := c.Get(ctx, u, &key)
	return key.Active.Pk, err
}
