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

func TestOperationMetadataAddressRegistryDiff(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		content := `{
  "protocol": "PtTALLiNtPec7mE7yY4m3k26J8Qukef3E3ehzhfXgFZKGtDdAXu",
  "chain_id": "NetXe8DbhW9A1eS",
  "hash": "onhS8L7VH7h1nXTVUvcCqp9EDAFXUvPnBSCepmMuXhwmVGSMPFC",
  "branch": "BLBnHF7w1ZaMRJGXmtkd2Q66Z6faToq7MDaHDS4wiT769KYk8Ai",
  "contents": [
    {
      "kind": "transaction",
      "source": "tz1ZvUkxJHPTy1tC7kF8Fg1Ko8jFvSumeENg",
      "fee": "481",
      "counter": "646854",
      "gas_limit": "1846",
      "storage_limit": "22",
      "amount": "0",
      "destination": "KT19BqYoEU1dVRsVWcWQXneherXR4mLpf4LZ",
      "parameters": {
        "entrypoint": "default",
        "value": {
          "string": "tz1ZvUkxJHPTy1tC7kF8Fg1Ko8jFvSumeENg"
        }
      },
      "metadata": {
        "balance_updates": [
          {
            "kind": "contract",
            "contract": "tz1ZvUkxJHPTy1tC7kF8Fg1Ko8jFvSumeENg",
            "change": "-481",
            "origin": "block"
          },
          {
            "kind": "accumulator",
            "category": "block fees",
            "change": "481",
            "origin": "block"
          }
        ],
        "operation_result": {
          "status": "applied",
          "storage": {
            "prim": "Some",
            "args": [
              {
                "int": "5"
              }
            ]
          },
          "balance_updates": [
            {
              "kind": "contract",
              "contract": "tz1ZvUkxJHPTy1tC7kF8Fg1Ko8jFvSumeENg",
              "change": "-500",
              "origin": "block"
            },
            {
              "kind": "burned",
              "category": "storage fees",
              "change": "500",
              "origin": "block"
            }
          ],
          "consumed_milligas": "1745500",
          "storage_size": "46",
          "paid_storage_size_diff": "2",
          "address_registry_diff": [
            {
              "address": "tz1ZvUkxJHPTy1tC7kF8Fg1Ko8jFvSumeENg",
              "index": "5"
            }
          ]
        }
      }
    }
  ],
  "signature": "sigNLZww8NtqgKsEUuEKM6XXzwwsHkDobNpfBY2zrr7HVM7pwvnCz2RepaXpCRswJK7QzirUBce7yV1z2ahREjBxDWP4zGrZ"
}`
		if r.URL.Path != "/chains/main/blocks/BLBnHF7w1ZaMRJGXmtkd2Q66Z6faToq7MDaHDS4wiT769KYk8Ai/operations/3/0" {
			status = http.StatusBadRequest
			content =
				fmt.Sprintf("\"Expected to request '/chains/main/blocks/BLBnHF7w1ZaMRJGXmtkd2Q66Z6faToq7MDaHDS4wiT769KYk8Ai/operations/3/0', got: %s\"", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			status = http.StatusBadRequest
			content =
				fmt.Sprintf("\"Expected Accept: application/json header, got: %s\"", r.Header.Get("Accept"))
		}

		w.WriteHeader(status)
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, []byte(content)); err != nil {
			panic(err)
		}
		w.Write(buffer.Bytes())
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	c.MetadataMode = MetadataModeAlways
	value, e := c.GetBlockOperation(context.TODO(), tezos.MustParseBlockHash("BLBnHF7w1ZaMRJGXmtkd2Q66Z6faToq7MDaHDS4wiT769KYk8Ai"), 3, 0)
	assert.NoError(t, e)
	assert.Len(t, value.Contents, 1)
	assert.Equal(t, &[]AddressRegistryDiff{{Address: tezos.MustParseAddress("tz1ZvUkxJHPTy1tC7kF8Fg1Ko8jFvSumeENg"), Index: 5}}, value.Contents.N(0).Meta().Result.AddressRegistryDiff)
}
