// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"context"
	"fmt"
)

type IssuanceParameters struct {
	Cycle                   int64  `json:"cycle"`
	BakingReward            int64  `json:"baking_reward_fixed_portion,string"`
	BakingBonusPerSlot      *int64 `json:"baking_reward_bonus_per_slot,string,omitempty"`  // -v023
	BakingBonusPerBlock     *int64 `json:"baking_reward_bonus_per_block,string,omitempty"` // v024
	AttestingRewardPerSlot  *int64 `json:"attesting_reward_per_slot,string,omitempty"`     // -v023
	AttestingRewardPerBlock *int64 `json:"attesting_reward_per_block,string,omitempty"`    // v024
	LBSubsidy               int64  `json:"liquidity_baking_subsidy,string"`
	SeedNonceTip            int64  `json:"seed_nonce_revelation_tip,string"`
	VdfTip                  int64  `json:"vdf_revelation_tip,string"`
}

// GetIssuance returns expected xtz issuance for known future cycles
func (c *Client) GetIssuance(ctx context.Context, id BlockID) ([]IssuanceParameters, error) {
	u := fmt.Sprintf("chains/main/blocks/%s/context/issuance/expected_issuance", id)
	p := make([]IssuanceParameters, 0, 5)
	if err := c.Get(ctx, u, &p); err != nil {
		return nil, err
	}
	return p, nil
}
