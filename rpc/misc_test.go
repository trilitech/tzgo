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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
			return
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
		{
			Delegate:     tezos.MustParseAddress("tz1N29q5T3jJ2i1JEWHax7q1NRkDMADj6fof"),
			ConsensusPkh: tezos.MustParseAddress("tz1N29q5T3jJ2i1JEWHax7q1NRkDMADj6fof"),
			BakingPower:  47741065155535,
		},
		{
			Delegate:     tezos.MustParseAddress("tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"),
			ConsensusPkh: tezos.MustParseAddress("tz4UmdB1JufxWrapeLysYxK9J574oFgkaAYs"),
			BakingPower:  199293741687329,
		},
		{
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
		{
			input:    "[1, 2, 3]",
			errorMsg: "unexpected input size; outer array should have length 2",
		},
		{
			input:    "[1, 2]",
			errorMsg: "failed to parse total baking power: json: cannot unmarshal number into Go value of type string",
		},
		{
			input:    `["a", 2]`,
			errorMsg: "failed to parse total baking power: strconv.ParseInt: parsing \"a\": invalid syntax",
		},
		{
			input:    `["1", ["a"]]`,
			errorMsg: "failed to parse baking power delegates: json: cannot unmarshal string into Go value of type []json.RawMessage",
		},
		{
			input:    `["1", [["a"]]]`,
			errorMsg: "failed to parse delegate info: unexpected input size; inner delegate array should have length 2",
		},
		{
			input:    `["1", [["a", "b"]]]`,
			errorMsg: "failed to parse delegate info: json: cannot unmarshal string into Go value of type rpc.BakingPowerDelegate",
		},
		{
			input:    `["1", [[{"delegate": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"}, "b"]]]`,
			errorMsg: "failed to parse delegate info: invalid or missing addresses",
		},
		{
			input:    `["1", [[{"delegate": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF", "consensus_pkh":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"}, 1]]]`,
			errorMsg: "failed to parse delegate info: json: cannot unmarshal number into Go value of type string",
		},
		{
			input:    `["1", [[{"delegate": "tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF", "consensus_pkh":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"}, "b"]]]`,
			errorMsg: "failed to parse delegate info: strconv.ParseInt: parsing \"b\": invalid syntax",
		},
	} {
		e := json.Unmarshal([]byte(test.input), &dummy)
		assert.EqualError(t, e, test.errorMsg, "case %d error mismatch", i)
	}
}

func TestGetTz4BakerNumberRatio(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
			return
		}

		var content []byte
		switch r.URL.String() {
		case "/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/helpers/tz4_baker_number_ratio?cycle=42":
			content = []byte(`"12.34%"`)
			w.WriteHeader(http.StatusOK)
		default:
			content = []byte(fmt.Sprintf("\"Unexpected URL: %s\"", r.URL.String()))
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Write(content)
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetTz4BakerNumberRatio(context.TODO(), tezos.MustParseBlockHash("BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ"), 42)
	assert.NoError(t, e)
	assert.Equal(t, float64(12.34), value)
}

func TestGetDestinationIndex(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
			return
		}

		var content []byte
		status := http.StatusOK
		switch r.URL.String() {
		case "/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/context/destination/tz1ZvUkxJHPTy1tC7kF8Fg1Ko8jFvSumeENg/index":
			content = []byte(`"42"`)
		case "/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/context/destination/KT1TxqZ8QtKvLu3V3JH7Gx58n7Co8pgtpQU5/index":
			content = []byte(`null`)
		// this case shouldn't really exist normally; just to check the error message
		case "/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/context/destination/tz1XiVu2qdyg2Lwpq3rc9a46rdgjH2txuUsz/index":
			content = []byte(`"bad response"`)
		default:
			content = []byte(fmt.Sprintf("\"Unexpected URL: %s\"", r.URL.String()))
			status = http.StatusBadRequest
		}

		w.WriteHeader(status)
		w.Write(content)
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetDestinationIndex(context.TODO(), tezos.MustParseBlockHash("BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ"), tezos.MustParseAddress("tz1ZvUkxJHPTy1tC7kF8Fg1Ko8jFvSumeENg"))
	assert.NoError(t, e)
	assert.Equal(t, uint64(42), *value)

	value, e = c.GetDestinationIndex(context.TODO(), tezos.MustParseBlockHash("BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ"), tezos.MustParseAddress("KT1TxqZ8QtKvLu3V3JH7Gx58n7Co8pgtpQU5"))
	assert.NoError(t, e)
	assert.Nil(t, value)

	value, e = c.GetDestinationIndex(context.TODO(), tezos.MustParseBlockHash("BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ"), tezos.MustParseAddress("tz1XiVu2qdyg2Lwpq3rc9a46rdgjH2txuUsz"))
	assert.EqualError(t, e, "failed to parse index: strconv.ParseUint: parsing \"bad response\": invalid syntax")
	assert.Nil(t, value)
}

func TestGetBlockValidatorsPreV024(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
			return
		}

		var content []byte
		status := http.StatusOK
		switch r.URL.String() {
		case "/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/helpers/validators":
			content = []byte(`[
  {
    "level": 17141799,
    "delegate": "tz3cpZyzsqEMk5RGHYL7m7obnLFtq7LxvCLA",
    "slots": [694, 813, 1050, 1271, 2093, 2580, 3669, 4042, 4396, 6078, 6611],
    "consensus_key": "tz3aYdfwVMHSdXjRfFDNR9ZFS1rqpMPj9uCp"
  },
  {
    "level": 17141799,
    "delegate": "tz3btDQsDkqq2G7eBdrrLqetaAfLVw6BnPez",
    "slots": [48, 55, 57, 185, 190, 248, 450, 456, 460, 517, 585, 649, 830, 849, 1065, 1074, 1079, 1101, 1170, 1223, 1317],
    "consensus_key": "tz3btDQsDkqq2G7eBdrrLqetaAfLVw6BnPez",
	"companion_key": "tz4Lc6qm8t6cS4UNdse9pAXm7pxiUmMruuvw"
  }
]`)
		default:
			content = []byte(fmt.Sprintf("\"Unexpected URL: %s\"", r.URL.String()))
			status = http.StatusBadRequest
		}

		w.WriteHeader(status)
		w.Write(content)
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockValidators(context.TODO(), tezos.MustParseBlockHash("BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ"), nil)
	assert.NoError(t, e)

	companionKey := tezos.MustParseAddress("tz4Lc6qm8t6cS4UNdse9pAXm7pxiUmMruuvw")
	preV024Value, e := value.AsPreV024Value()
	assert.NoError(t, e)
	assert.Equal(t, []GetBlockValidatorsResponsePreV024{{
		Level:        17141799,
		Delegate:     tezos.MustParseAddress("tz3cpZyzsqEMk5RGHYL7m7obnLFtq7LxvCLA"),
		Slots:        []uint16{694, 813, 1050, 1271, 2093, 2580, 3669, 4042, 4396, 6078, 6611},
		ConsensusKey: tezos.MustParseAddress("tz3aYdfwVMHSdXjRfFDNR9ZFS1rqpMPj9uCp"),
	}, {
		Level:        17141799,
		Delegate:     tezos.MustParseAddress("tz3btDQsDkqq2G7eBdrrLqetaAfLVw6BnPez"),
		Slots:        []uint16{48, 55, 57, 185, 190, 248, 450, 456, 460, 517, 585, 649, 830, 849, 1065, 1074, 1079, 1101, 1170, 1223, 1317},
		ConsensusKey: tezos.MustParseAddress("tz3btDQsDkqq2G7eBdrrLqetaAfLVw6BnPez"),
		CompanionKey: &companionKey,
	}}, *preV024Value)

	v024Value, e := value.AsV024Value()
	assert.Error(t, e)
	assert.Nil(t, v024Value)
}

func TestGetBlockValidatorsV024(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
			return
		}

		var content []byte
		status := http.StatusOK
		switch r.URL.String() {
		case "/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/helpers/validators":
			content = []byte(`[
  {
    "level": 766274,
    "consensus_threshold": "604936371385957",
    "consensus_committee": "907339747101888",
    "all_bakers_attest_activated": true,
    "delegates": [
      {
        "delegate": "tz4HG14YMihpZxynRXu7tK72hoz3mnnXZGzm",
        "rounds": [6938, 6980, 6984],
        "consensus_key": "tz4HG14YMihpZxynRXu7tK72hoz3mnnXZGzm",
        "companion_key": "tz4Lc6qm8t6cS4UNdse9pAXm7pxiUmMruuvw",
        "attesting_power": "55831337496322",
        "attestation_slot": 11
      },
      {
        "delegate": "tz3Q1fwk1vh3zm5LqyUV9e2wZBdaEXcovh2r",
        "rounds": [1, 71],
        "consensus_key": "tz3Q1fwk1vh3zm5LqyUV9e2wZBdaEXcovh2r",
        "attesting_power": "38734749288786",
        "attestation_slot": 10
      }
    ]
  }
]`)
		default:
			content = []byte(fmt.Sprintf("\"Unexpected URL: %s\"", r.URL.String()))
			status = http.StatusBadRequest
		}

		w.WriteHeader(status)
		w.Write(content)
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockValidators(context.TODO(), tezos.MustParseBlockHash("BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ"), nil)
	assert.NoError(t, e)

	preV024Value, e := value.AsPreV024Value()
	assert.Error(t, e)
	assert.Nil(t, preV024Value)

	companionKey := tezos.MustParseAddress("tz4Lc6qm8t6cS4UNdse9pAXm7pxiUmMruuvw")
	v024Value, e := value.AsV024Value()
	assert.NoError(t, e)
	assert.Equal(t, []GetBlockValidatorsResponseV024{{
		Level:                    766274,
		ConsensusThreshold:       uint64(604936371385957),
		ConsensusCommittee:       uint64(907339747101888),
		AllBakersAttestActivated: true,
		Delegates: []struct {
			Delegate        tezos.Address  `json:"delegate"`
			Rounds          []int32        `json:"rounds"`
			ConsensusKey    tezos.Address  `json:"consensus_key"`
			CompanionKey    *tezos.Address `json:"companion_key"`
			AttestingPower  int64          `json:"attesting_power,string"`
			AttestationSlot uint16         `json:"attestation_slot"`
		}{{
			Delegate:        tezos.MustParseAddress("tz4HG14YMihpZxynRXu7tK72hoz3mnnXZGzm"),
			Rounds:          []int32{6938, 6980, 6984},
			ConsensusKey:    tezos.MustParseAddress("tz4HG14YMihpZxynRXu7tK72hoz3mnnXZGzm"),
			CompanionKey:    &companionKey,
			AttestingPower:  int64(55831337496322),
			AttestationSlot: uint16(11),
		}, {
			Delegate:        tezos.MustParseAddress("tz3Q1fwk1vh3zm5LqyUV9e2wZBdaEXcovh2r"),
			Rounds:          []int32{1, 71},
			ConsensusKey:    tezos.MustParseAddress("tz3Q1fwk1vh3zm5LqyUV9e2wZBdaEXcovh2r"),
			AttestingPower:  int64(38734749288786),
			AttestationSlot: uint16(10),
		}},
	}}, *v024Value)
}
