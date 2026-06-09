// Copyright (c) 2026 TriliTech Ltd.
// Author: tzstats@trili.tech

package rpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trilitech/tzgo/tezos"
)

const v025TestBlock = "BL9gWjpSWuNgJ7vm4pVgBg1QWhTcwGis3rDzFtEkLAvSSP8vbzB"

// newJSONServer returns a test server that routes by URL path suffix.
func newJSONServer(t *testing.T, routes map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for suffix, body := range routes {
			if strings.Contains(r.URL.Path, suffix) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(body))
				return
			}
		}
		t.Errorf("unexpected path: %s", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestGetDalPastParameters(t *testing.T) {
	body := `{
		"feature_enable": true, "incentives_enable": true, "dynamic_lag_enable": true,
		"number_of_slots": 160, "attestation_lag": 5, "attestation_lags": [1,2,3,4,5],
		"attestation_threshold": 66,
		"minimal_participation_ratio": {"numerator": 1, "denominator": 2},
		"rewards_ratio": {"numerator": 1, "denominator": 10},
		"traps_fraction": {"numerator": 1, "denominator": 2000},
		"redundancy_factor": 8, "page_size": 3967, "slot_size": 380832, "number_of_shards": 512
	}`
	server := newJSONServer(t, map[string]string{"/context/dal/past_parameters/100": body})
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	p, err := c.GetDalPastParameters(context.TODO(), tezos.MustParseBlockHash(v025TestBlock), 100)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, int64(160), p.NumberOfSlots)
	assert.Equal(t, int64(380832), p.SlotSize)
	assert.Equal(t, []int64{1, 2, 3, 4, 5}, p.AttestationLags)
	assert.Equal(t, 2, p.MinimalParticipationRatio.Den)
}

func TestDecodeEncodeDalAttestation(t *testing.T) {
	decodeBody := `[{"lag_index": 0, "slot_indices": [1, 5, 9]}, {"lag_index": 1, "slot_indices": [2]}]`
	server := newJSONServer(t, map[string]string{
		"/helpers/decode_dal_attestation/": decodeBody,
		"/helpers/encode_dal_attestation":  `"546"`,
	})
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	id := tezos.MustParseBlockHash(v025TestBlock)

	dec, err := c.DecodeDalAttestation(context.TODO(), id, "546")
	assert.NoError(t, err)
	assert.Len(t, dec, 2)
	assert.Equal(t, 0, dec[0].LagIndex)
	assert.Equal(t, []int{1, 5, 9}, dec[0].SlotIndices)

	enc, err := c.EncodeDalAttestation(context.TODO(), id, dec)
	assert.NoError(t, err)
	assert.Equal(t, "546", enc.String())
}

func TestGetSwrrCredits(t *testing.T) {
	body := `[{"delegate": "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx", "credit": "1000"}]`
	server := newJSONServer(t, map[string]string{"/helpers/swrr_credits": body})
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	credits, err := c.GetSwrrCredits(context.TODO(), tezos.MustParseBlockHash(v025TestBlock))
	assert.NoError(t, err)
	assert.Len(t, credits, 1)
	assert.Equal(t, "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx", credits[0].Delegate.String())
	assert.Equal(t, int64(1000), credits[0].Credit.Int64())
}

// TestGetSwrrCreditsNonActivated verifies that a non_activated_feature RPC error
// is mapped to (nil, nil) so callers can treat a disabled feature flag gracefully.
func TestGetSwrrCreditsNonActivated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`[{"kind":"permanent","id":"proto.025-PsUshuai.non_activated_feature"}]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	credits, err := c.GetSwrrCredits(context.TODO(), tezos.MustParseBlockHash(v025TestBlock))
	assert.NoError(t, err)
	assert.Nil(t, credits)
}

func TestGetStezSupplyAndRate(t *testing.T) {
	server := newJSONServer(t, map[string]string{
		"/context/stez/total_supply":        `"123456"`,
		"/context/stez/total_amount_of_tez": `"654321"`,
		"/context/stez/exchange_rate":       `{"numerator": "654321", "denominator": "123456"}`,
	})
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	id := tezos.MustParseBlockHash(v025TestBlock)

	supply, err := c.GetStezTotalSupply(context.TODO(), id)
	assert.NoError(t, err)
	assert.Equal(t, int64(123456), supply.Int64())

	amount, err := c.GetStezTotalAmountOfTez(context.TODO(), id)
	assert.NoError(t, err)
	assert.Equal(t, int64(654321), amount.Int64())

	rate, err := c.GetStezExchangeRate(context.TODO(), id)
	assert.NoError(t, err)
	assert.NotNil(t, rate)
	assert.Equal(t, int64(654321), rate.Numerator.Int64())
	assert.Equal(t, int64(123456), rate.Denominator.Int64())
}
