package rpc

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trilitech/tzgo/tezos"
)

// Mainnet reference fixtures (fetched from tezos rpc node).
//
// Kept as embedded byte payloads (compile-time, fully offline) and used as ground truth
// for structural + semantic decoding assertions.
//
// NOTE: These JSON fixtures MUST preserve field order for operation contents, because
// rpc.OperationList has a fast-path unmarshaller that peeks into `{"kind":"..."}`.
//
//go:embed testdata/mainnet_block_level_11680082.json
var mainnetBlockLevel11680082JSON []byte

//go:embed testdata/mainnet_block_level_11610082.json
var mainnetBlockLevel11610082JSON []byte

func mustUnmarshalBlock(t *testing.T, raw []byte) *Block {
	t.Helper()
	var b Block
	err := json.Unmarshal(raw, &b)
	if err != nil {
		t.Fatalf("unmarshal block: %v", err)
	}
	return &b
}

func TestGetBlockHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		switch r.URL.Path {
		case "/chains/main/blocks/BLc35zTguA6svYv1o8P9RaJVBRHTEuPPcczjsHuJv7JiLhLoug3/header":
			// v023
			content := []byte(`{
  "protocol": "PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh",
  "chain_id": "NetXnHfVqm9iesp",
  "hash": "BLc35zTguA6svYv1o8P9RaJVBRHTEuPPcczjsHuJv7JiLhLoug3",
  "level": 16718489,
  "proto": 13,
  "predecessor": "BKirYQXBE7JpLresSpRP9gMoodVAKKA3opkHdicPEL4zsnXxJqJ",
  "timestamp": "2025-12-04T17:07:52Z",
  "validation_pass": 4,
  "operations_hash": "LLoZnRYK2PErb7dFnciFvHyXuRaqXqLJWQ7yZrAz3e4R12bFcsZMG",
  "fitness": [
    "02",
    "00ff1a99",
    "",
    "ffffffff",
    "00000000"
  ],
  "context": "CoVoctaEQeg222YUutru7iXofoyjsid8M5ZZy3MJSkHZoi6hHxLg",
  "payload_hash": "vh2XQ2N35VpDfWgX39NHy9GrXBUDPYbU72YwD4hJnt6kCk5dNVSH",
  "payload_round": 0,
  "proof_of_work_nonce": "13afca5d43820200",
  "liquidity_baking_toggle_vote": "on",
  "adaptive_issuance_vote": "pass",
  "signature": "BLsig9yXoMX9CodJwTeANdxF9eutKmaKtUqNz8QxcrQwv2hLxtjCq6M7HtQoUSj5d7LiqeFg6TA1vibZcoPqegAyi7jhVch2jmpkMfWKmApS7NFKUeYcwJUGHwU4kXPmH567fkWrm1qSp4"
}`)
			w.WriteHeader(http.StatusOK)
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, content); err != nil {
				panic(err)
			}
			w.Write(buffer.Bytes())
		case "/chains/main/blocks/BMZ9thR8HyS9JKbCBCydeKvXnuGWzcBB4cwn3nbCdaSrzbyLSRV/header":
			// v024
			content := []byte(`{
  "protocol": "PtTALLiNtPec7mE7yY4m3k26J8Qukef3E3ehzhfXgFZKGtDdAXu",
  "chain_id": "NetXe8DbhW9A1eS",
  "hash": "BMZ9thR8HyS9JKbCBCydeKvXnuGWzcBB4cwn3nbCdaSrzbyLSRV",
  "level": 340114,
  "proto": 2,
  "predecessor": "BL5YeeL6gJE8SMz7EaVhLow2pcbxPkoNTY651HQqAi1GHRxYJ3Q",
  "timestamp": "2025-12-04T16:19:47Z",
  "validation_pass": 4,
  "operations_hash": "LLoZdgAMEVtYDe7ceVPop35BYoMfKy2Ap8sga45Yjr9JQZWnAgDPX",
  "fitness": [
    "02",
    "00053092",
    "",
    "ffffffff",
    "00000000"
  ],
  "context": "CoVdmJtRQEn34rQwcTEE1RznKgYpcWjfJrvbrYEjVDXbMzGKe32Y",
  "payload_hash": "vh3TEnHL6vFNebqDtQnt4rAEfrorSt4kiqB7saK3B2FSZ4UpxGFT",
  "payload_round": 0,
  "proof_of_work_nonce": "60afd4da00000000",
  "liquidity_baking_toggle_vote": "pass",
  "signature": "sigdLTfKRoEENPECAk451xTaTThvmcu6ZzedtRU3Nt3EPA3M5sUdqUXLKB9FHX48aitXuekxemmJ6x4QSse6SNtV7aGeraL1"
}`)
			w.WriteHeader(http.StatusOK)
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, content); err != nil {
				panic(err)
			}
			w.Write(buffer.Bytes())
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "unexpected url: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	// v023
	value, e := c.GetBlockHeader(context.TODO(), tezos.MustParseBlockHash("BLc35zTguA6svYv1o8P9RaJVBRHTEuPPcczjsHuJv7JiLhLoug3"))
	assert.Nil(t, e)
	assert.Equal(t, tezos.FeatureVotePass, *value.AdaptiveIssuanceVote)

	// v024; AdaptiveIssuanceVote is removed
	value, e = c.GetBlockHeader(context.TODO(), tezos.MustParseBlockHash("BMZ9thR8HyS9JKbCBCydeKvXnuGWzcBB4cwn3nbCdaSrzbyLSRV"))
	assert.Nil(t, e)
	assert.Nil(t, value.AdaptiveIssuanceVote)
}

func TestBlockHeaderProtocolData_NilAdaptiveIssuanceVote(t *testing.T) {
	// Regression test: starting with v024, `adaptive_issuance_vote` is removed from
	// block headers and unmarshals to nil. ProtocolData must not panic in this case.
	//
	// See incident report: core-inc-2604.
	var h BlockHeader
	err := json.Unmarshal([]byte(`{
		"payload_hash": "vh3TEnHL6vFNebqDtQnt4rAEfrorSt4kiqB7saK3B2FSZ4UpxGFT",
		"payload_round": 0,
		"proof_of_work_nonce": "60afd4da00000000",
		"liquidity_baking_toggle_vote": "pass"
	}`), &h)
	assert.NoError(t, err)
	assert.Nil(t, h.AdaptiveIssuanceVote)

	var data []byte
	assert.NotPanics(t, func() {
		data = h.ProtocolData()
	})

	// protocol_data without signature (invalid/missing) is:
	// payload_hash(32) + payload_round(4) + pow_nonce(8) + seed_nonce_presence(1) + votes(1)
	assert.Len(t, data, 32+4+8+1+1)
	assert.Equal(t, byte(0x0), data[32+4+8])          // seed nonce absent
	assert.Equal(t, h.LbVote().Tag(), data[32+4+8+1]) // ai vote missing => ai tag == 0
}

func TestGetBlockMetadataV023(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		switch r.URL.Path {
		case "/chains/main/blocks/BL1CW14TFJ2XHdqmMCB1yH126yQdYrXDhQAuTuY8r4io7AFwDQr/metadata":
			content := []byte(`{
  "protocol": "PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh",
  "next_protocol": "PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh",
  "test_chain_status": {
    "status": "not_running"
  },
  "max_operations_ttl": 150,
  "max_operation_data_length": 32768,
  "max_block_header_length": 289,
  "max_operation_list_length": [
    {
      "max_size": 4194304,
      "max_op": 2048
    },
    {
      "max_size": 32768
    },
    {
      "max_size": 135168,
      "max_op": 132
    },
    {
      "max_size": 524288
    }
  ],
  "proposer": "tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs",
  "baker": "tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs",
  "level_info": {
    "level": 3383167,
    "level_position": 3383166,
    "cycle": 11277,
    "cycle_position": 66,
    "expected_commitment": false
  },
  "voting_period_info": {
    "voting_period": {
      "index": 11277,
      "kind": "proposal",
      "start_position": 3383100
    },
    "position": 66,
    "remaining": 233
  },
  "nonce_hash": null,
  "deactivated": [],
  "balance_updates": [
    {
      "kind": "accumulator",
      "category": "block fees",
      "change": "-800",
      "origin": "block"
    },
    {
      "kind": "contract",
      "contract": "tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs",
      "change": "800",
      "origin": "block"
    },
    {
      "kind": "minted",
      "category": "baking rewards",
      "change": "-300784355",
      "origin": "block"
    },
    {
      "kind": "freezer",
      "category": "deposits",
      "staker": {
        "baker_own_stake": "tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs"
      },
      "change": "300784355",
      "origin": "block"
    },
    {
      "kind": "minted",
      "category": "baking rewards",
      "change": "-13028170",
      "origin": "block"
    },
    {
      "kind": "contract",
      "contract": "tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs",
      "change": "13028170",
      "origin": "block"
    },
    {
      "kind": "minted",
      "category": "baking bonuses",
      "change": "-300495226",
      "origin": "block"
    },
    {
      "kind": "freezer",
      "category": "deposits",
      "staker": {
        "baker_own_stake": "tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs"
      },
      "change": "300495226",
      "origin": "block"
    },
    {
      "kind": "minted",
      "category": "baking bonuses",
      "change": "-13015647",
      "origin": "block"
    },
    {
      "kind": "contract",
      "contract": "tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs",
      "change": "13015647",
      "origin": "block"
    }
  ],
  "liquidity_baking_toggle_ema": 0,
  "adaptive_issuance_vote_ema": 0,
  "adaptive_issuance_activation_cycle": 0,
  "implicit_operations_results": [
    {
      "kind": "transaction",
      "storage": [
        {
          "int": "1"
        },
        {
          "int": "1127720872378"
        },
        {
          "int": "100"
        },
        {
          "bytes": "01e927f00ef734dfc85919635e9afc9166c83ef9fc00"
        },
        {
          "bytes": "0115eb0104481a6d7921160bc982c5e0a561cd8a3a00"
        }
      ],
      "balance_updates": [
        {
          "kind": "minted",
          "category": "subsidy",
          "change": "-333333",
          "origin": "subsidy"
        },
        {
          "kind": "contract",
          "contract": "KT1TxqZ8QtKvLu3V3JH7Gx58n7Co8pgtpQU5",
          "change": "333333",
          "origin": "subsidy"
        }
      ],
      "consumed_milligas": "206582",
      "storage_size": "4632"
    }
  ],
  "proposer_consensus_key": "tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL",
  "baker_consensus_key": "tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL",
  "consumed_milligas": "3000000",
  "dal_attestation": "1"
}`)
			w.WriteHeader(http.StatusOK)
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, content); err != nil {
				panic(err)
			}
			w.Write(buffer.Bytes())
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "unexpected url: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockMetadata(context.TODO(), tezos.MustParseBlockHash("BL1CW14TFJ2XHdqmMCB1yH126yQdYrXDhQAuTuY8r4io7AFwDQr"))
	assert.NoError(t, e)
	assert.Nil(t, value.Attestations)
	assert.Nil(t, value.Preattestations)
	assert.Nil(t, value.AllBakersAttestActivationLevel)
}

func TestGetBlockMetadataV024(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		switch r.URL.Path {
		case "/chains/main/blocks/BMRzRcNz2NzBNYy9XAfA4siKYknDV2zaDJ9QeBQTbC2Bg1Cs6YF/metadata":
			content := []byte(`{
  "protocol": "PtTALLiNtPec7mE7yY4m3k26J8Qukef3E3ehzhfXgFZKGtDdAXu",
  "next_protocol": "PtTALLiNtPec7mE7yY4m3k26J8Qukef3E3ehzhfXgFZKGtDdAXu",
  "test_chain_status": {
    "status": "not_running"
  },
  "max_operations_ttl": 150,
  "max_operation_data_length": 32768,
  "max_block_header_length": 289,
  "max_operation_list_length": [
    {
      "max_size": 4194304,
      "max_op": 2048
    },
    {
      "max_size": 32768
    },
    {
      "max_size": 135168,
      "max_op": 132
    },
    {
      "max_size": 524288
    }
  ],
  "proposer": "tz1TnEtqDV9mZyts2pfMy6Jw1BTPs4LMjL8M",
  "baker": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF",
  "level_info": {
    "level": 429176,
    "level_position": 429175,
    "cycle": 1430,
    "cycle_position": 175,
    "expected_commitment": false
  },
  "voting_period_info": {
    "voting_period": {
      "index": 1430,
      "kind": "proposal",
      "start_position": 429000
    },
    "position": 175,
    "remaining": 124
  },
  "nonce_hash": null,
  "deactivated": [],
  "balance_updates": [
    {
      "kind": "accumulator",
      "category": "block fees",
      "change": "-800",
      "origin": "block"
    },
    {
      "kind": "contract",
      "contract": "tz1TnEtqDV9mZyts2pfMy6Jw1BTPs4LMjL8M",
      "change": "800",
      "origin": "block"
    },
    {
      "kind": "minted",
      "category": "baking rewards",
      "change": "-256955481",
      "origin": "block"
    },
    {
      "kind": "freezer",
      "category": "deposits",
      "staker": {
        "baker_own_stake": "tz1TnEtqDV9mZyts2pfMy6Jw1BTPs4LMjL8M"
      },
      "change": "256955481",
      "origin": "block"
    },
    {
      "kind": "minted",
      "category": "baking rewards",
      "change": "-45497615",
      "origin": "block"
    },
    {
      "kind": "contract",
      "contract": "tz1TnEtqDV9mZyts2pfMy6Jw1BTPs4LMjL8M",
      "change": "45497615",
      "origin": "block"
    },
    {
      "kind": "minted",
      "category": "baking bonuses",
      "change": "-251202404",
      "origin": "block"
    },
    {
      "kind": "freezer",
      "category": "deposits",
      "staker": {
        "baker_own_stake": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"
      },
      "change": "251202404",
      "origin": "block"
    },
    {
      "kind": "minted",
      "category": "baking bonuses",
      "change": "-198",
      "origin": "block"
    },
    {
      "kind": "contract",
      "contract": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF",
      "change": "198",
      "origin": "block"
    }
  ],
  "liquidity_baking_toggle_ema": 38325558,
  "implicit_operations_results": [],
  "proposer_consensus_key": "tz1TnEtqDV9mZyts2pfMy6Jw1BTPs4LMjL8M",
  "baker_consensus_key": "tz4UmdB1JufxWrapeLysYxK9J574oFgkaAYs",
  "consumed_milligas": "3000000",
  "dal_attestation": "256",
  "all_bakers_attest_activation_level": {
    "level": 15601,
    "level_position": 15600,
    "cycle": 52,
    "cycle_position": 0,
    "expected_commitment": false
  },
  "attestations": {
    "total_committee_power": "469521635844122",
    "threshold": "313036782067802",
    "recorded_power": "443005372003831"
  },
  "preattestations": {
    "total_committee_power": "469521635844122",
    "threshold": "313036782067802",
    "recorded_power": "468739947472511"
  }
}`)
			w.WriteHeader(http.StatusOK)
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, content); err != nil {
				panic(err)
			}
			w.Write(buffer.Bytes())
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "unexpected url: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	// v024; new fields `attestations` and `preattestations`
	value, e := c.GetBlockMetadata(context.TODO(), tezos.MustParseBlockHash("BMRzRcNz2NzBNYy9XAfA4siKYknDV2zaDJ9QeBQTbC2Bg1Cs6YF"))
	assert.NoError(t, e)
	assert.Equal(t, int64(469521635844122), value.Attestations.TotalCommitteePower)
	assert.Equal(t, int64(313036782067802), value.Attestations.Threshold)
	assert.Equal(t, int64(443005372003831), value.Attestations.RecordedPower)
	assert.Equal(t, int64(469521635844122), value.Preattestations.TotalCommitteePower)
	assert.Equal(t, int64(313036782067802), value.Preattestations.Threshold)
	assert.Equal(t, int64(468739947472511), value.Preattestations.RecordedPower)

	// v024; new field `all_bakers_attest_activation_level`
	// The field is actually optional in v24 but here that case is skipped as it's rather trivial
	assert.Equal(t, &LevelInfo{
		Level:              15601,
		LevelPosition:      15600,
		Cycle:              52,
		CyclePosition:      0,
		ExpectedCommitment: false,
	}, value.AllBakersAttestActivationLevel)
}

func TestBlock_MainnetFixtures_AllHelpers(t *testing.T) {
	// v024 (Tallinn)
	b24 := mustUnmarshalBlock(t, mainnetBlockLevel11680082JSON)
	assert.Equal(t, int64(11680082), b24.GetLevel())
	assert.True(t, b24.GetTimestamp().Equal(b24.Header.Timestamp))
	assert.Equal(t, 24, b24.GetVersion())
	assert.Equal(t, int64(1126), b24.GetCycle())
	assert.Equal(t, int64(1126), b24.GetLevelInfo().Cycle)
	assert.Equal(t, int64(166), b24.GetVotingPeriod())
	assert.Equal(t, tezos.VotingPeriodProposal, b24.GetVotingPeriodKind())
	vpi := b24.GetVotingInfo()
	assert.Equal(t, int64(166), vpi.VotingPeriod.Index)
	assert.Equal(t, tezos.VotingPeriodProposal, vpi.VotingPeriod.Kind)
	assert.Equal(t, int64(39793), vpi.Position)
	assert.Equal(t, int64(161806), vpi.Remaining)
	assert.False(t, b24.IsProtocolUpgrade())

	// v023 (Seoul)
	b23 := mustUnmarshalBlock(t, mainnetBlockLevel11610082JSON)
	assert.Equal(t, int64(11610082), b23.GetLevel())
	assert.True(t, b23.GetTimestamp().Equal(b23.Header.Timestamp))
	assert.Equal(t, 23, b23.GetVersion())
	assert.Equal(t, int64(1121), b23.GetCycle())
	assert.Equal(t, int64(1121), b23.GetLevelInfo().Cycle)
	assert.Equal(t, int64(165), b23.GetVotingPeriod())
	assert.Equal(t, tezos.VotingPeriodAdoption, b23.GetVotingPeriodKind())
	vpi = b23.GetVotingInfo()
	assert.Equal(t, int64(165), vpi.VotingPeriod.Index)
	assert.Equal(t, tezos.VotingPeriodAdoption, vpi.VotingPeriod.Kind)
	assert.Equal(t, int64(120993), vpi.Position)
	assert.Equal(t, int64(30206), vpi.Remaining)
	assert.False(t, b23.IsProtocolUpgrade())
}

func TestBlock_FallbackHelpers_UsingLegacyMetadataLevel(t *testing.T) {
	kind := tezos.VotingPeriodProposal
	b := Block{
		Header: BlockHeader{
			Level: 123,
			Proto: 7,
		},
		Metadata: BlockMetadata{
			Level: &LevelInfo{
				Level:                123,
				Cycle:                9,
				VotingPeriod:         3,
				VotingPeriodPosition: 5,
			},
			VotingPeriodKind: &kind,
		},
	}
	assert.Equal(t, int64(123), b.GetLevel())
	assert.Equal(t, 7, b.GetVersion())
	assert.Equal(t, int64(9), b.GetCycle())
	assert.Equal(t, int64(9), b.GetLevelInfo().Cycle)
	assert.Equal(t, int64(3), b.GetVotingPeriod())
	assert.Equal(t, tezos.VotingPeriodProposal, b.GetVotingPeriodKind())
	vpi := b.GetVotingInfo()
	assert.Equal(t, int64(5), vpi.Position)
	assert.Equal(t, int64(32768-5), vpi.Remaining)
	assert.Equal(t, int64(3), vpi.VotingPeriod.Index)
	assert.Equal(t, tezos.VotingPeriodProposal, vpi.VotingPeriod.Kind)
	assert.Equal(t, int64(3*32768), vpi.VotingPeriod.StartPosition)
}

func TestBlockMetadata_GetLevel(t *testing.T) {
	m := &BlockMetadata{}
	assert.Equal(t, int64(0), m.GetLevel())
	m.LevelInfo = &LevelInfo{Level: 42}
	assert.Equal(t, int64(42), m.GetLevel())
	m.LevelInfo = nil
	m.Level = &LevelInfo{Level: 7}
	assert.Equal(t, int64(7), m.GetLevel())
}

func TestBlockHeader_LbVote_FallbackEscapeVote(t *testing.T) {
	h := BlockHeader{
		// LiquidityBakingToggleVote is invalid/zero -> fallback to escape vote.
		LiquidityBakingEscapeVote: true,
	}
	assert.Equal(t, tezos.FeatureVoteOff, h.LbVote())
	h.LiquidityBakingEscapeVote = false
	assert.Equal(t, tezos.FeatureVoteOn, h.LbVote())
}

func TestBlockHeader_ProtocolData_MainnetVotesAndLength(t *testing.T) {
	b24 := mustUnmarshalBlock(t, mainnetBlockLevel11680082JSON)
	h24 := b24.Header
	// v024 mainnet blocks do not include adaptive_issuance_vote in header JSON.
	assert.Nil(t, h24.AdaptiveIssuanceVote)
	data := h24.ProtocolData()
	// payload_hash(32) + payload_round(4) + pow_nonce(8) + seed_nonce_presence(1) + votes(1) + signature(64)
	assert.Len(t, data, 32+4+8+1+1+64)
	// votes byte: lb(on)=0, ai(nil)=0
	assert.Equal(t, byte(0), data[32+4+8+1])

	b23 := mustUnmarshalBlock(t, mainnetBlockLevel11610082JSON)
	h23 := b23.Header
	assert.NotNil(t, h23.AdaptiveIssuanceVote)
	data = h23.ProtocolData()
	assert.Len(t, data, 32+4+8+1+1+64)
	// votes byte: lb(pass)=2, ai(pass)=2 => 2 | (2<<2) = 10
	assert.Equal(t, byte(10), data[32+4+8+1])
}

func TestClient_BlockAPI_AllMethods(t *testing.T) {
	b24 := mustUnmarshalBlock(t, mainnetBlockLevel11680082JSON)
	b23 := mustUnmarshalBlock(t, mainnetBlockLevel11610082JSON)

	// Build header payloads that include the extra fields returned by /header.
	hdrJSON := func(b *Block) []byte {
		h := b.Header
		h.Hash = b.Hash
		h.Protocol = b.Protocol
		h.ChainId = b.ChainId
		buf, err := json.Marshal(&h)
		if err != nil {
			t.Fatalf("marshal header: %v", err)
		}
		return buf
	}
	metaJSON := func(b *Block) []byte {
		buf, err := json.Marshal(&b.Metadata)
		if err != nil {
			t.Fatalf("marshal metadata: %v", err)
		}
		return buf
	}

	// Minimal genesis block payload.
	//
	// Important: omit fields like `header.signature` to avoid tezos.Signature JSON parsing errors
	// on empty/invalid values.
	genesisJSON := []byte(fmt.Sprintf(
		`{"protocol":%q,"chain_id":%q,"hash":%q,"header":{"level":0,"proto":0,"validation_pass":0},"metadata":{"protocol":%q,"next_protocol":%q},"operations":[[],[],[],[]]}`,
		b23.Protocol.String(),
		b23.ChainId.String(),
		"BLockGenesisGenesisGenesisGenesisGenesisf79b5d1CoW2",
		b23.Protocol.String(),
		b23.Protocol.String(),
	))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
			return
		}
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/chains/main/blocks/11680082":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBlockLevel11680082JSON)
		case "/chains/main/blocks/11610082":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBlockLevel11610082JSON)
		case "/chains/main/blocks/head":
			// use level 11680082 as deterministic "head"
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBlockLevel11680082JSON)
		case "/chains/main/blocks/genesis":
			w.WriteHeader(http.StatusOK)
			w.Write(genesisJSON)

		case "/chains/main/blocks/head/header":
			w.WriteHeader(http.StatusOK)
			w.Write(hdrJSON(b24))
		case "/chains/main/blocks/11680082/header":
			w.WriteHeader(http.StatusOK)
			w.Write(hdrJSON(b24))
		case "/chains/main/blocks/11610082/header":
			w.WriteHeader(http.StatusOK)
			w.Write(hdrJSON(b23))

		case "/chains/main/blocks/11680082/metadata":
			w.WriteHeader(http.StatusOK)
			w.Write(metaJSON(b24))
		case "/chains/main/blocks/11610082/metadata":
			w.WriteHeader(http.StatusOK)
			w.Write(metaJSON(b23))

		case "/chains/main/blocks/11680082/hash":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%q", b24.Hash.String())

		case "/chains/main/blocks":
			// used by GetTips and GetBlockPredHashes
			q := r.URL.Query()
			if q.Get("length") == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "missing length query")
				return
			}
			// return a deterministic list-of-lists of block hashes
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `[[%q,%q]]`, b24.Hash.String(), b23.Hash.String())

		case "/chains/main/invalid_blocks":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `[{"block":%q,"level":%d,"error":[{"id":"proto.test.invalid_block","kind":"temporary"}]}]`, b23.Hash.String(), b23.Header.Level)
		case "/chains/main/invalid_blocks/" + b23.Hash.String():
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"block":%q,"level":%d,"error":[{"id":"proto.test.invalid_block","kind":"temporary"}]}`, b23.Hash.String(), b23.Header.Level)

		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "unexpected url: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)

	// GetBlock / GetBlockHeight
	got, e := c.GetBlock(context.Background(), BlockLevel(11680082))
	assert.NoError(t, e)
	assert.Equal(t, b24.Hash, got.Hash)

	got, e = c.GetBlockHeight(context.Background(), 11610082)
	assert.NoError(t, e)
	assert.Equal(t, b23.Hash, got.Hash)

	// GetHeadBlock / GetGenesisBlock
	got, e = c.GetHeadBlock(context.Background())
	assert.NoError(t, e)
	assert.Equal(t, b24.Hash, got.Hash)

	got, e = c.GetGenesisBlock(context.Background())
	if !assert.NoError(t, e) {
		return
	}
	assert.Equal(t, int64(0), got.Header.Level)

	// GetTips
	tips, e := c.GetTips(context.Background(), 2, tezos.BlockHash{})
	assert.NoError(t, e)
	assert.Len(t, tips, 1)
	assert.Equal(t, b24.Hash, tips[0][0])

	// GetTipHeader / GetBlockHeader
	h, e := c.GetTipHeader(context.Background())
	assert.NoError(t, e)
	assert.Equal(t, b24.Hash, h.Hash)

	h, e = c.GetBlockHeader(context.Background(), BlockLevel(11610082))
	assert.NoError(t, e)
	assert.Equal(t, b23.Hash, h.Hash)

	// GetBlockMetadata
	m, e := c.GetBlockMetadata(context.Background(), BlockLevel(11680082))
	assert.NoError(t, e)
	assert.Equal(t, b24.Metadata.Protocol, m.Protocol)

	// GetBlockHash
	bh, e := c.GetBlockHash(context.Background(), BlockLevel(11680082))
	assert.NoError(t, e)
	assert.Equal(t, b24.Hash, bh)

	// GetBlockPredHashes (count <= 0 defaults to 1)
	preds, e := c.GetBlockPredHashes(context.Background(), b24.Hash, 0)
	assert.NoError(t, e)
	assert.Len(t, preds, 2)
	assert.Equal(t, b24.Hash, preds[0])

	// GetInvalidBlocks / GetInvalidBlock
	invalid, e := c.GetInvalidBlocks(context.Background())
	assert.NoError(t, e)
	assert.Len(t, invalid, 1)
	assert.Equal(t, b23.Hash, invalid[0].Block)
	assert.Equal(t, ErrorKindTemporary, invalid[0].Error.ErrorKind())

	one, e := c.GetInvalidBlock(context.Background(), b23.Hash)
	assert.NoError(t, e)
	assert.Equal(t, b23.Hash, one.Block)
	assert.Equal(t, ErrorKindTemporary, one.Error.ErrorKind())
}

func TestMainnetFixture_Sanity_Timestamps(t *testing.T) {
	// Guard against silent fixture drift when regenerating payloads.
	b24 := mustUnmarshalBlock(t, mainnetBlockLevel11680082JSON)
	ts, err := time.Parse(time.RFC3339, "2026-01-27T12:11:37Z")
	assert.NoError(t, err)
	assert.True(t, b24.Header.Timestamp.Equal(ts))

	b23 := mustUnmarshalBlock(t, mainnetBlockLevel11610082JSON)
	ts, err = time.Parse(time.RFC3339, "2026-01-21T20:38:20Z")
	assert.NoError(t, err)
	assert.True(t, b23.Header.Timestamp.Equal(ts))
}
