// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trilitech/tzgo/tezos"
)

func TestParseDoubleConsensusOperationEvidence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations" {
			t.Errorf("Expected to request '/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		content := []byte(`[
  [],
  [],
  [
    {
      "protocol": "PsRiotumaAMotcRoDWW1bysEhQy2n1M5fy8JgRp8jjRfHGmfeA7",
      "chain_id": "NetXd56aBs1aeW3",
      "hash": "ooqgVxC8XDYHXSnznWSdkXktgUaA2PZdUs4azRbco9i6fhiY1ui",
      "branch": "BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb",
      "contents": [
        {
          "kind": "double_consensus_operation_evidence",
          "slot": 0,
          "op1": {
            "branch": "BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb",
            "operations": {
              "kind": "attestation",
              "level": 1331,
              "block_payload_hash": "vh1g87ZG6scSYxKhspAUzprQVuLAyoa5qMBKcUfjgnQGnFb3dJcG",
              "round": 0,
              "slot": 0
            },
            "signature": "sigbQ5ZNvkjvGssJgoAnUAfY4Wvvg3QZqawBYB1j1VDBNTMBAALnCzRHWzer34bnfmzgHg3EvwdzQKdxgSghB897cono6gbQ"
          },
          "op2": {
            "branch": "BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb",
            "operations": {
              "kind": "preattestation",
              "level": 1331,
              "block_payload_hash": "vh1g87ZG6scSYxKhspAUzprQVuLAyoa5qMBKcUfjgnQGnFb3dJcG",
              "round": 0,
              "slot": 0
            },
            "signature": "sigbQ5ZNvkjvGssJgoAnUAfY4Wvvg3QZqawBYB1j1VDBNTMBAALnCzRHWzer34bnfmzgHg3EvwdzQKdxgSghB897cono6gbQ"
          }
        }
      ]
    }
  ],
  []
]`)
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, content); err != nil {
			panic(err)
		}
		w.Write(buffer.Bytes())
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)
	assert.Len(t, value, 4)
	assert.Len(t, value[2], 1)
	assert.Equal(t, value[2][0].Contents.Len(), 1)
	op := value[2][0].Contents.Select(tezos.OpTypeDoubleConsensusOperationEvidence, 0).(*DoubleConsensusOperationEvidence)
	assert.Equal(t, tezos.OpTypeDoubleConsensusOperationEvidence, op.Kind())
	assert.Equal(t, tezos.OpTypeAttestation, op.Op1.Operations.OpKind)
	assert.Equal(t, "sigbQ5ZNvkjvGssJgoAnUAfY4Wvvg3QZqawBYB1j1VDBNTMBAALnCzRHWzer34bnfmzgHg3EvwdzQKdxgSghB897cono6gbQ", op.Op1.Signature.String())
	assert.Equal(t, tezos.OpTypePreattestation, op.Op2.Operations.OpKind)
}
