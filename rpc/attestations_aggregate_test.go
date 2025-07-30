// Copyright (c) 2020-2021 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package rpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trilitech/tzgo/tezos"
)

func TestParseAttestationsAggregate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations" {
			t.Errorf("Expected to request '/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[[{"protocol":"PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh","chain_id":"NetXd56aBs1aeW3","hash":"ooqgVxC8XDYHXSnznWSdkXktgUaA2PZdUs4azRbco9i6fhiY1ui","branch":"BLwNFCbHuF21bF4S9KybZd52si5QQG6k29mQTshbY7fVnehrxbh","contents":[{"kind":"attestations_aggregate","consensus_content":{"level":109771,"round":0,"block_payload_hash":"vh3RiisbNp7QBvLkJ6HAyoYfMh4AoZLcJFVE5M34LLKnwAWTRv6J"},"committee":[{"slot":0,"dal_attestation":"0"},{"slot":4,"dal_attestation":"0"},{"slot":15},{"slot":21,"dal_attestation":"0"},{"slot":203}],"metadata":{"committee":[{"delegate":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF","consensus_pkh":"tz4X5GCfEHQCUnrBy9Qo1PSsmExYHXxiEkvp","consensus_power":1435},{"delegate":"tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs","consensus_pkh":"tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL","consensus_power":1465},{"delegate":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_pkh":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_power":707},{"delegate":"tz3PgFHdYvEGEbUo1pUJmuNH8fgc8cwARKfC","consensus_pkh":"tz4EWkmNN93yE7HrjaRR6mGh22rUgSYJG1Sj","consensus_power":730},{"delegate":"tz1LmrwzCKUDibk7xaGC5RxvTbmUbCAtCA4a","consensus_pkh":"tz4QqpgVTG4CCETKMgV6YaUHyowE4Gkqwdfi","consensus_power":1}],"total_consensus_power":4338}}],"signature":"BLsigB1uDeiiuW1NPKWsZ6WRAKL3aSGTXKmtsDH2xxWgWqM3QBr3mFW8QqzH6VWGsvPGsrii3VVw7KA9CvC9LjC3VxH3MgHSvcWVK6Z7rBbEY79sKXi4XrbbfY8QJpE38B4u6mteGKHnVj"}],[],[],[]]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)
	assert.Equal(t, tezos.OpTypeAttestationsAggregate, value[0][0].Contents.N(0).Kind())
}
