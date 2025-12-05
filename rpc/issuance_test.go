// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trilitech/tzgo/tezos"
)

func TestGetIssuanceV023(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
			return
		}
		var content []byte
		if r.URL.Path != "/chains/main/blocks/BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb/context/issuance/expected_issuance" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown URL: %s", r.URL)
		} else {
			// v023: baking_reward_bonus_per_slot
			w.WriteHeader(http.StatusOK)
			content = []byte(`[
  {
    "cycle": 1909,
    "baking_reward_fixed_portion": "49242420",
    "baking_reward_bonus_per_slot": "21091",
    "attesting_reward_per_slot": "14033",
    "seed_nonce_revelation_tip": "1154079",
    "vdf_revelation_tip": "1154079",
    "dal_attesting_reward_per_shard": "42716"
  },
  {
    "cycle": 1910,
    "baking_reward_fixed_portion": "49224030",
    "baking_reward_bonus_per_slot": "21083",
    "attesting_reward_per_slot": "14028",
    "seed_nonce_revelation_tip": "1153648",
    "vdf_revelation_tip": "1153648",
    "dal_attesting_reward_per_shard": "42700"
  },
  {
    "cycle": 1911,
    "baking_reward_fixed_portion": "49205648",
    "baking_reward_bonus_per_slot": "21075",
    "attesting_reward_per_slot": "14023",
    "seed_nonce_revelation_tip": "1153217",
    "vdf_revelation_tip": "1153217",
    "dal_attesting_reward_per_shard": "42684"
  }
]`)
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, content); err != nil {
				panic(err)
			}
			w.Write(buffer.Bytes())
		}
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	res, e := c.GetIssuance(context.TODO(), tezos.MustParseBlockHash("BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb"))
	assert.NoError(t, e)
	assert.Len(t, res, 3)
	assert.Equal(t, int64(21091), *(res[0].BakingBonusPerSlot))
	assert.Nil(t, res[0].BakingBonusPerBlock)
	assert.Equal(t, int64(14033), *(res[0].AttestingRewardPerSlot))
	assert.Nil(t, res[0].AttestingRewardPerBlock)
}

func TestGetIssuanceV024(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
			return
		}
		var content []byte
		if r.URL.Path != "/chains/main/blocks/BLJodc1mt3TyFci3DvBZvryQUp2HB6QzDdZDED6rw6CCif4e4AA/context/issuance/expected_issuance" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown URL: %s", r.URL)
		} else {
			// v024: baking_reward_bonus_per_block
			w.WriteHeader(http.StatusOK)
			content = []byte(`[
  {
    "cycle": 1205,
    "baking_reward_fixed_portion": "302194513",
    "baking_reward_bonus_per_block": "302194513",
    "attesting_reward_per_block": "604389027",
    "seed_nonce_revelation_tip": "1475378",
    "vdf_revelation_tip": "1475378",
    "dal_attesting_reward_per_shard": "262131"
  },
  {
    "cycle": 1206,
    "baking_reward_fixed_portion": "302195663",
    "baking_reward_bonus_per_block": "302195663",
    "attesting_reward_per_block": "604391327",
    "seed_nonce_revelation_tip": "1475384",
    "vdf_revelation_tip": "1475384",
    "dal_attesting_reward_per_shard": "262132"
  },
  {
    "cycle": 1207,
    "baking_reward_fixed_portion": "302196813",
    "baking_reward_bonus_per_block": "302196813",
    "attesting_reward_per_block": "604393626",
    "seed_nonce_revelation_tip": "1475389",
    "vdf_revelation_tip": "1475389",
    "dal_attesting_reward_per_shard": "262133"
  }
]`)
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, content); err != nil {
				panic(err)
			}
			w.Write(buffer.Bytes())
		}
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	res, e := c.GetIssuance(context.TODO(), tezos.MustParseBlockHash("BLJodc1mt3TyFci3DvBZvryQUp2HB6QzDdZDED6rw6CCif4e4AA"))
	assert.NoError(t, e)
	assert.Len(t, res, 3)
	assert.Equal(t, int64(302194513), *(res[0].BakingBonusPerBlock))
	assert.Nil(t, res[0].BakingBonusPerSlot)
	assert.Equal(t, int64(604389027), *(res[0].AttestingRewardPerBlock))
	assert.Nil(t, res[0].AttestingRewardPerSlot)
}
