// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

package rpc

import (
	"context"
	"fmt"
)

// GetAllBakersAttestActivationLevel returns the first level at which the All Bakers Attest feature activates. Nil is returned if the feature is not yet effective. RPC introduced in v024.
func (c *Client) GetAllBakersAttestActivationLevel(ctx context.Context, id BlockID) (*LevelInfo, error) {
	var level *LevelInfo
	u := fmt.Sprintf("chains/main/blocks/%s/helpers/all_bakers_attest_activation_level", id)
	if err := c.Get(ctx, u, &level); err != nil {
		return nil, err
	}

	return level, nil
}
