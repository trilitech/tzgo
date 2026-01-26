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
