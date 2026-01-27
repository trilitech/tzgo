// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

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

	"github.com/stretchr/testify/assert"
	"github.com/trilitech/tzgo/tezos"
)

// Offline fixtures for operation endpoints at mainnet block level 11680082.
//
//go:embed testdata/mainnet_block_level_11680082_operation_hashes.json
var mainnetBlockLevel11680082OperationHashesJSON []byte

//go:embed testdata/mainnet_block_level_11680082_operation_hashes_0.json
var mainnetBlockLevel11680082OperationHashes0JSON []byte

//go:embed testdata/mainnet_block_level_11680082_operation_hashes_0_0.json
var mainnetBlockLevel11680082OperationHashes00JSON []byte

//go:embed testdata/mainnet_block_level_11680082_operations.json
var mainnetBlockLevel11680082OperationsJSON []byte

//go:embed testdata/mainnet_block_level_11680082_operations_0.json
var mainnetBlockLevel11680082Operations0JSON []byte

//go:embed testdata/mainnet_block_level_11680082_operations_0_0.json
var mainnetBlockLevel11680082Operations00JSON []byte

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

func TestBlockOperationsAPI_FromTestdataFixtures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
			return
		}

		switch r.URL.Path {
		case "/chains/main/blocks/11680082/operation_hashes":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBlockLevel11680082OperationHashesJSON)
		case "/chains/main/blocks/11680082/operation_hashes/0":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBlockLevel11680082OperationHashes0JSON)
		case "/chains/main/blocks/11680082/operation_hashes/0/0":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBlockLevel11680082OperationHashes00JSON)

		case "/chains/main/blocks/11680082/operations":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBlockLevel11680082OperationsJSON)
		case "/chains/main/blocks/11680082/operations/0":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBlockLevel11680082Operations0JSON)
		case "/chains/main/blocks/11680082/operations/0/0":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBlockLevel11680082Operations00JSON)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "unexpected url: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)

	// hashes
	h, err := c.GetBlockOperationHash(context.Background(), BlockLevel(11680082), 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, tezos.MustParseOpHash("opQ6YSNmjZi2XkM7GTeLXTJyzYZzdKCfUo9qcqPaXXcC82p77Qx"), h)

	hashes, err := c.GetBlockOperationHashes(context.Background(), BlockLevel(11680082))
	assert.NoError(t, err)
	assert.Len(t, hashes, 4)
	assert.Greater(t, len(hashes[0]), 0)
	assert.Equal(t, tezos.MustParseOpHash("opQ6YSNmjZi2XkM7GTeLXTJyzYZzdKCfUo9qcqPaXXcC82p77Qx"), hashes[0][0])

	listHashes, err := c.GetBlockOperationListHashes(context.Background(), BlockLevel(11680082), 0)
	assert.NoError(t, err)
	assert.Greater(t, len(listHashes), 0)
	assert.Equal(t, tezos.MustParseOpHash("opQ6YSNmjZi2XkM7GTeLXTJyzYZzdKCfUo9qcqPaXXcC82p77Qx"), listHashes[0])

	// operations
	op, err := c.GetBlockOperation(context.Background(), BlockLevel(11680082), 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, tezos.MustParseOpHash("opQ6YSNmjZi2XkM7GTeLXTJyzYZzdKCfUo9qcqPaXXcC82p77Qx"), op.Hash)
	assert.GreaterOrEqual(t, len(op.Contents), 1)
	assert.Equal(t, tezos.OpTypeAttestationWithDal, op.Contents[0].Kind())

	ops0, err := c.GetBlockOperationList(context.Background(), BlockLevel(11680082), 0)
	assert.NoError(t, err)
	assert.Greater(t, len(ops0), 0)
	assert.Equal(t, tezos.MustParseOpHash("opQ6YSNmjZi2XkM7GTeLXTJyzYZzdKCfUo9qcqPaXXcC82p77Qx"), ops0[0].Hash)

	ops, err := c.GetBlockOperations(context.Background(), BlockLevel(11680082))
	assert.NoError(t, err)
	assert.Len(t, ops, 4)
	assert.Greater(t, len(ops[0]), 0)
	assert.Equal(t, tezos.MustParseOpHash("opQ6YSNmjZi2XkM7GTeLXTJyzYZzdKCfUo9qcqPaXXcC82p77Qx"), ops[0][0].Hash)
}
