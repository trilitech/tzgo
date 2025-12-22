// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/trilitech/tzgo/tezos"
)

type BakingPowerDelegate struct {
	Delegate     tezos.Address `json:"delegate"`
	ConsensusPkh tezos.Address `json:"consensus_pkh"`
	BakingPower  int64
}

type BakingPowerDistribution struct {
	TotalBakingPower int64
	Delegates        []BakingPowerDelegate
}

func (b *BakingPowerDistribution) UnmarshalJSON(data []byte) error {
	var rawRes []json.RawMessage
	if err := json.Unmarshal(data, &rawRes); err != nil {
		return err
	}
	if len(rawRes) != 2 {
		return fmt.Errorf("unexpected input size; outer array should have length 2")
	}

	totalStr := ""
	if err := json.Unmarshal(rawRes[0], &totalStr); err != nil {
		return fmt.Errorf("failed to parse total baking power: %v", err)
	}
	v, err := strconv.ParseInt(totalStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse total baking power: %v", err)
	}
	b.TotalBakingPower = v

	var rawDelegates [][]json.RawMessage
	if err := json.Unmarshal(rawRes[1], &rawDelegates); err != nil {
		return fmt.Errorf("failed to parse baking power delegates: %v", err)
	}
	for _, rawDelegate := range rawDelegates {
		if len(rawDelegate) != 2 {
			return fmt.Errorf("failed to parse delegate info: unexpected input size; inner delegate array should have length 2")
		}
		var d BakingPowerDelegate
		if err := json.Unmarshal(rawDelegate[0], &d); err != nil {
			return fmt.Errorf("failed to parse delegate info: %v", err)
		}
		if !d.ConsensusPkh.IsValid() || !d.Delegate.IsValid() {
			return fmt.Errorf("failed to parse delegate info: invalid or missing addresses")
		}
		powerStr := ""
		if err := json.Unmarshal(rawDelegate[1], &powerStr); err != nil {
			return fmt.Errorf("failed to parse delegate info: %v", err)
		}
		v, err := strconv.ParseInt(powerStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse delegate info: %v", err)
		}
		d.BakingPower = v
		b.Delegates = append(b.Delegates, d)
	}

	return nil
}

// GetAllBakersAttestActivationLevel returns the first level at which the All Bakers Attest feature activates. Nil is returned if the feature is not yet effective. RPC introduced in v024.
func (c *Client) GetAllBakersAttestActivationLevel(ctx context.Context, id BlockID) (*LevelInfo, error) {
	var level *LevelInfo
	u := fmt.Sprintf("chains/main/blocks/%s/helpers/all_bakers_attest_activation_level", id)
	if err := c.Get(ctx, u, &level); err != nil {
		return nil, err
	}

	return level, nil
}

// GetBakingPowerDistributionForCurrentCycle returns the total baking power and baking power distribution for all the active delegates. RPC introduced in v024.
func (c *Client) GetBakingPowerDistributionForCurrentCycle(ctx context.Context, id BlockID) (*BakingPowerDistribution, error) {
	var d *BakingPowerDistribution
	u := fmt.Sprintf("chains/main/blocks/%s/helpers/baking_power_distribution_for_current_cycle", id)
	if err := c.Get(ctx, u, &d); err != nil {
		return nil, err
	}

	return d, nil
}
