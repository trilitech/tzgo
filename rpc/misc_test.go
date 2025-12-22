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
