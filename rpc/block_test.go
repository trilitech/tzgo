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
