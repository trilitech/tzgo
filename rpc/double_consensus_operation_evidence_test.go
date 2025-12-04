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
		if r.URL.Path != "/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/operations" {
			t.Errorf("Expected to request '/chains/main/blocks/BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ/operations', got: %s", r.URL.Path)
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
      "branch": "BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ",
      "contents": [
        {
          "kind": "double_consensus_operation_evidence",
          "slot": 0,
          "op1": {
            "branch": "BLkHSfhY7kff5Jzo9Nvj7hCu2hi2SeDaNmSTmUU7SDaGGhEWQo9",
            "operations": {
              "kind": "attestation_with_dal",
              "slot": 0,
              "level": 150562,
              "round": 0,
              "block_payload_hash": "vh2Sqiiseutbb13hthkf2usRi3rko4wTcq5yQKn6zubBDwoCHJoT",
              "dal_attestation": "0"
            },
            "signature": "BLsig9tpbsg9VSzGaeZxvRRR9HJ62AhQmtEX1RFUp6AvBsP62aN3jFEkdNq6u2imXAjPUEQ4Z3Lt4H9h7nHhXkpUXF6WKymJFRxB3aZbm4AydnT7hNnSEWGRQ3wpjMBvQ1giAX7UTodyZi"
          },
          "op2": {
            "branch": "BLkHSfhY7kff5Jzo9Nvj7hCu2hi2SeDaNmSTmUU7SDaGGhEWQo9",
            "operations": {
              "kind": "attestation",
              "slot": 0,
              "level": 150562,
              "round": 0,
              "block_payload_hash": "vh1g87ZG6scSYxKhspAUzprQVuLAyoa5qMBKcUfjgnQGnFb3dJcG"
            },
            "signature": "BLsig9yLWpztFBMJRd4bacFXnE4mgyTruWPWCVNA4vYFP9sgv4t9NVtgVs5Tx1umTiSzYWtmbAsrZhYLCwjhpXQAssTaNNtCHeZAbmXZRHJhqsTee8GWzWftoH36VvDEpg9NnahmMvwwE1"
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
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BL2xnxQ6YPmySpE9CZgDUFPKD146Hku36aojiGGF1gWfbijeihZ"))
	assert.Nil(t, e)
	assert.Len(t, value, 4)
	assert.Len(t, value[2], 1)
	assert.Equal(t, value[2][0].Contents.Len(), 1)
	op := value[2][0].Contents.Select(tezos.OpTypeDoubleConsensusOperationEvidence, 0).(*DoubleConsensusOperationEvidence)
	assert.Equal(t, tezos.OpTypeDoubleConsensusOperationEvidence, op.Kind())
	assert.Equal(t, tezos.OpTypeAttestationWithDal, op.Op1.Operations.OpKind)
	assert.Equal(t, "BLsig9tpbsg9VSzGaeZxvRRR9HJ62AhQmtEX1RFUp6AvBsP62aN3jFEkdNq6u2imXAjPUEQ4Z3Lt4H9h7nHhXkpUXF6WKymJFRxB3aZbm4AydnT7hNnSEWGRQ3wpjMBvQ1giAX7UTodyZi", op.Op1.Signature.String())
	assert.Equal(t, tezos.OpTypeAttestation, op.Op2.Operations.OpKind)
	assert.Equal(t, "BLsig9yLWpztFBMJRd4bacFXnE4mgyTruWPWCVNA4vYFP9sgv4t9NVtgVs5Tx1umTiSzYWtmbAsrZhYLCwjhpXQAssTaNNtCHeZAbmXZRHJhqsTee8GWzWftoH36VvDEpg9NnahmMvwwE1", op.Op2.Signature.String())
}
