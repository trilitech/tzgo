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

func TestGetAllBakersAttestActivationLevel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		var content []byte
		switch r.URL.Path {
		// block where the feature is not yet activated
		case "/chains/main/blocks/BLCeBQsUGfMnEuYJJrBGBYpgJpt3mTu71MoL97SZkNyRTrFixsw/helpers/all_bakers_attest_activation_level":
			content = []byte(`null`)
			w.WriteHeader(http.StatusOK)
		// block where the feature is activated
		case "/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/helpers/all_bakers_attest_activation_level":
			content = []byte(`{
  "level": 15601,
  "level_position": 15600,
  "cycle": 52,
  "cycle_position": 0,
  "expected_commitment": false
}`)
			w.WriteHeader(http.StatusOK)
		default:
			content = []byte(fmt.Sprintf("\"Unexpected URL: %s\"", r.URL.Path))
			w.WriteHeader(http.StatusBadRequest)
		}

		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, content); err != nil {
			panic(err)
		}
		w.Write(buffer.Bytes())
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	// block where the feature is not yet activated
	value, e := c.GetAllBakersAttestActivationLevel(context.TODO(), tezos.MustParseBlockHash("BLCeBQsUGfMnEuYJJrBGBYpgJpt3mTu71MoL97SZkNyRTrFixsw"))
	assert.NoError(t, e)
	assert.Nil(t, value)

	// block where the feature is activated
	value, e = c.GetAllBakersAttestActivationLevel(context.TODO(), tezos.MustParseBlockHash("BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ"))
	assert.NoError(t, e)
	assert.Equal(t, &LevelInfo{
		Level:              15601,
		LevelPosition:      15600,
		Cycle:              52,
		CyclePosition:      0,
		ExpectedCommitment: false,
	}, value)
}

func TestGetBakingPowerDistributionForCurrentCycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		var content []byte
		switch r.URL.Path {
		case "/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/helpers/baking_power_distribution_for_current_cycle":
			content = []byte(`[
  "299309813454073",
  [
    [
      {
        "delegate": "tz1N29q5T3jJ2i1JEWHax7q1NRkDMADj6fof",
        "consensus_pkh": "tz1N29q5T3jJ2i1JEWHax7q1NRkDMADj6fof"
      },
      "47741065155535"
    ],
    [
      {
        "delegate": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF",
        "consensus_pkh": "tz4UmdB1JufxWrapeLysYxK9J574oFgkaAYs"
      },
      "199293741687329"
    ],
    [
      {
        "delegate": "tz4HG14YMihpZxynRXu7tK72hoz3mnnXZGzm",
        "consensus_pkh": "tz4HG14YMihpZxynRXu7tK72hoz3mnnXZGzm"
      },
      "52275006611209"
    ]
  ]
]`)
			w.WriteHeader(http.StatusOK)
		default:
			content = []byte(fmt.Sprintf("\"Unexpected URL: %s\"", r.URL.Path))
			w.WriteHeader(http.StatusBadRequest)
		}

		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, content); err != nil {
			panic(err)
		}
		w.Write(buffer.Bytes())
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBakingPowerDistributionForCurrentCycle(context.TODO(), tezos.MustParseBlockHash("BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ"))
	assert.NoError(t, e)
	assert.Equal(t, int64(299309813454073), value.TotalBakingPower)
	assert.Equal(t, []BakingPowerDelegate{
		BakingPowerDelegate{
			Delegate:     tezos.MustParseAddress("tz1N29q5T3jJ2i1JEWHax7q1NRkDMADj6fof"),
			ConsensusPkh: tezos.MustParseAddress("tz1N29q5T3jJ2i1JEWHax7q1NRkDMADj6fof"),
			BakingPower:  47741065155535,
		},
		BakingPowerDelegate{
			Delegate:     tezos.MustParseAddress("tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"),
			ConsensusPkh: tezos.MustParseAddress("tz4UmdB1JufxWrapeLysYxK9J574oFgkaAYs"),
			BakingPower:  199293741687329,
		},
		BakingPowerDelegate{
			Delegate:     tezos.MustParseAddress("tz4HG14YMihpZxynRXu7tK72hoz3mnnXZGzm"),
			ConsensusPkh: tezos.MustParseAddress("tz4HG14YMihpZxynRXu7tK72hoz3mnnXZGzm"),
			BakingPower:  52275006611209,
		},
	}, value.Delegates)
}

func TestUnmarshalBadBakingPowerDistribution(t *testing.T) {
	type testcase struct {
		input    string
		errorMsg string
	}

	var dummy *BakingPowerDistribution
	for i, test := range []testcase{
		testcase{
			input:    "[1, 2, 3]",
			errorMsg: "unexpected input size; outer array should have length 2",
		},
		testcase{
			input:    "[1, 2]",
			errorMsg: "failed to parse total baking power: json: cannot unmarshal number into Go value of type string",
		},
		testcase{
			input:    `["a", 2]`,
			errorMsg: "failed to parse total baking power: strconv.ParseInt: parsing \"a\": invalid syntax",
		},
		testcase{
			input:    `["1", ["a"]]`,
			errorMsg: "failed to parse baking power delegates: json: cannot unmarshal string into Go value of type []json.RawMessage",
		},
		testcase{
			input:    `["1", [["a"]]]`,
			errorMsg: "failed to parse delegate info: unexpected input size; inner delegate array should have length 2",
		},
		testcase{
			input:    `["1", [["a", "b"]]]`,
			errorMsg: "failed to parse delegate info: json: cannot unmarshal string into Go value of type rpc.BakingPowerDelegate",
		},
		testcase{
			input:    `["1", [[{"delegate": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"}, "b"]]]`,
			errorMsg: "failed to parse delegate info: invalid or missing addresses",
		},
		testcase{
			input:    `["1", [[{"delegate": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF", "consensus_pkh":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"}, 1]]]`,
			errorMsg: "failed to parse delegate info: json: cannot unmarshal number into Go value of type string",
		},
		testcase{
			input:    `["1", [[{"delegate": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF", "consensus_pkh":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"}, "b"]]]`,
			errorMsg: "failed to parse delegate info: strconv.ParseInt: parsing \"b\": invalid syntax",
		},
	} {
		e := json.Unmarshal([]byte(test.input), &dummy)
		assert.EqualError(t, e, test.errorMsg, "case %d error mismatch", i)
	}
}
