// Copyright (c) 2026 TriliTech Ltd.
// Author: tzstats@trili.tech

package rpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/trilitech/tzgo/tezos"
)

// This file wires up RPC endpoints introduced in protocol v025 (Ushuaia):
//   - DAL past parameters and attestation bitset decode/encode helpers
//   - SWRR (Stake-Weighted Round-Robin) baker selection helpers
//   - enshrined liquid staking (sTEZ) context queries
//
// Schemas follow the ushuaia-openapi specification.

// DalParameters mirrors the DAL parameters returned by the DAL parameter RPCs.
// Introduced together with dynamic attestation lag in v025 (Ushuaia).
type DalParameters struct {
	FeatureEnable             bool        `json:"feature_enable"`
	IncentivesEnable          bool        `json:"incentives_enable"`
	DynamicLagEnable          bool        `json:"dynamic_lag_enable"`
	NumberOfSlots             int64       `json:"number_of_slots"`
	AttestationLag            int64       `json:"attestation_lag"`
	AttestationLags           []int64     `json:"attestation_lags"`
	AttestationThreshold      int64       `json:"attestation_threshold"`
	MinimalParticipationRatio tezos.Ratio `json:"minimal_participation_ratio"`
	RewardsRatio              tezos.Ratio `json:"rewards_ratio"`
	TrapsFraction             tezos.Ratio `json:"traps_fraction"`
	RedundancyFactor          int64       `json:"redundancy_factor"`
	PageSize                  int64       `json:"page_size"`
	SlotSize                  int64       `json:"slot_size"`
	NumberOfShards            int64       `json:"number_of_shards"`
}

// GetDalPastParameters returns the DAL parameters that were active at the given
// level. RPC introduced in v025 (Ushuaia).
func (c *Client) GetDalPastParameters(ctx context.Context, id BlockID, level int64) (*DalParameters, error) {
	var params *DalParameters
	u := fmt.Sprintf("chains/main/blocks/%s/context/dal/past_parameters/%d", id, level)
	if err := c.Get(ctx, u, &params); err != nil {
		return nil, err
	}
	return params, nil
}

// DalAttestationByLag is the explicit set of attested DAL slot indices for a
// single attestation lag, as returned/consumed by the DAL attestation helpers.
type DalAttestationByLag struct {
	LagIndex    int   `json:"lag_index"`
	SlotIndices []int `json:"slot_indices"`
}

// DecodeDalAttestation decodes a DAL attestation bitset into an explicit
// representation of attested slots per lag, using the protocol parameters active
// at id. RPC introduced in v025 (Ushuaia).
func (c *Client) DecodeDalAttestation(ctx context.Context, id BlockID, bitset string) ([]DalAttestationByLag, error) {
	var res []DalAttestationByLag
	u := fmt.Sprintf("chains/main/blocks/%s/helpers/decode_dal_attestation/%s", id, bitset)
	if err := c.Get(ctx, u, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// EncodeDalAttestation encodes an explicit representation of attested slots per
// lag into a DAL attestation bitset (decimal big-number string), using the
// protocol parameters active at id. RPC introduced in v025 (Ushuaia).
func (c *Client) EncodeDalAttestation(ctx context.Context, id BlockID, slots []DalAttestationByLag) (tezos.Z, error) {
	var res tezos.Z
	u := fmt.Sprintf("chains/main/blocks/%s/helpers/encode_dal_attestation", id)
	if err := c.Post(ctx, u, slots, &res); err != nil {
		return tezos.Z{}, err
	}
	return res, nil
}

// SwrrSelectedBaker is one baker selected for round 0 of a cycle by the SWRR
// algorithm. Introduced in v025 (Ushuaia).
type SwrrSelectedBaker struct {
	ConsensusPk tezos.Key     `json:"consensus_pk"`
	Delegate    tezos.Address `json:"delegate"`
	CompanionPk tezos.Key     `json:"companion_pk"`
}

// GetSwrrSelectedBakers returns the bakers selected for round 0 of the cycle via
// the SWRR algorithm. If cycle < 0 the node's current cycle is used. Returns nil
// when the rights for the cycle have not been computed yet (or were pruned) and
// also nil (without error) when the swrr_new_baker_lottery_enable feature flag is
// not enabled. RPC introduced in v025 (Ushuaia).
func (c *Client) GetSwrrSelectedBakers(ctx context.Context, id BlockID, cycle int64) ([]SwrrSelectedBaker, error) {
	var res []SwrrSelectedBaker
	u := fmt.Sprintf("chains/main/blocks/%s/helpers/swrr_selected_bakers", id)
	if cycle >= 0 {
		u += fmt.Sprintf("?cycle=%d", cycle)
	}
	if err := c.Get(ctx, u, &res); err != nil {
		if isNonActivatedFeature(err) {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

// SwrrCredit is the remaining SWRR credit for a delegate. Introduced in v025.
type SwrrCredit struct {
	Delegate tezos.Address `json:"delegate"`
	Credit   tezos.Z       `json:"credit"`
}

// GetSwrrCredits returns the current SWRR state: the remaining credits for all
// active delegates. Returns nil (without error) when the
// swrr_new_baker_lottery_enable feature flag is not enabled. RPC introduced in
// v025 (Ushuaia).
func (c *Client) GetSwrrCredits(ctx context.Context, id BlockID) ([]SwrrCredit, error) {
	var res []SwrrCredit
	u := fmt.Sprintf("chains/main/blocks/%s/helpers/swrr_credits", id)
	if err := c.Get(ctx, u, &res); err != nil {
		if isNonActivatedFeature(err) {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

// GetStezTotalSupply returns the total supply of sTEZ tokens. RPC introduced in
// v025 (Ushuaia).
func (c *Client) GetStezTotalSupply(ctx context.Context, id BlockID) (tezos.Z, error) {
	return c.getStezBignum(ctx, id, "total_supply")
}

// GetStezTotalAmountOfTez returns the total amount of tez (in mutez) in the sTEZ
// staking ledger. RPC introduced in v025 (Ushuaia).
func (c *Client) GetStezTotalAmountOfTez(ctx context.Context, id BlockID) (tezos.Z, error) {
	return c.getStezBignum(ctx, id, "total_amount_of_tez")
}

func (c *Client) getStezBignum(ctx context.Context, id BlockID, field string) (tezos.Z, error) {
	var res tezos.Z
	u := fmt.Sprintf("chains/main/blocks/%s/context/stez/%s", id, field)
	if err := c.Get(ctx, u, &res); err != nil {
		return tezos.Z{}, err
	}
	return res, nil
}

// StezExchangeRate is the exchange rate between sTEZ and tez, expressed as the
// ratio total_amount_of_tez / total_supply (or 1/1 when total_supply is 0).
type StezExchangeRate struct {
	Numerator   tezos.Z `json:"numerator"`
	Denominator tezos.Z `json:"denominator"`
}

// GetStezExchangeRate returns the exchange rate between sTEZ and tez. RPC
// introduced in v025 (Ushuaia).
func (c *Client) GetStezExchangeRate(ctx context.Context, id BlockID) (*StezExchangeRate, error) {
	var res *StezExchangeRate
	u := fmt.Sprintf("chains/main/blocks/%s/context/stez/exchange_rate", id)
	if err := c.Get(ctx, u, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// isNonActivatedFeature reports whether err is a Tezos RPC error caused by a
// feature flag being disabled (error id ".../non_activated_feature").
func isNonActivatedFeature(err error) bool {
	var rpcErr RPCError
	if errors.As(err, &rpcErr) {
		for _, e := range rpcErr.Errors() {
			if strings.Contains(e.ErrorID(), "non_activated_feature") {
				return true
			}
		}
	}
	return false
}
