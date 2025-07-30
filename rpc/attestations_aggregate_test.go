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
		if r.URL.Path != "/chains/main/blocks/BKprStSDDvZQDu62cBQAtRnohsauwtMkQgBuKV3qv9StvRgkndd/operations" {
			t.Errorf("Expected to request '/chains/main/blocks/BKprStSDDvZQDu62cBQAtRnohsauwtMkQgBuKV3qv9StvRgkndd/operations', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[[{"protocol":"PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh","chain_id":"NetXd56aBs1aeW3","hash":"ooiGZe18gbU699rP8ChA9SpKjA1kL8oj45jauxSqiTakEKq1go4","branch":"BLX6FYnwSov7qsiQrSD6eHVorSz8Y91QARTxVeeNSBRVsPj2QuG","contents":[{"kind":"attestations_aggregate","consensus_content":{"level":410651,"round":0,"block_payload_hash":"vh3K7vaxgyKa6umEamKJt2CgiNGjopt6yfgiNmUcAa9UnHHGVHBH"},"committee":[{"slot":0,"dal_attestation":"1"},{"slot":2,"dal_attestation":"1"},{"slot":7,"dal_attestation":"1"},{"slot":16},{"slot":714},{"slot":1143},{"slot":5604},{"slot":5715}],"metadata":{"committee":[{"delegate":"tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs","consensus_pkh":"tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL","consensus_power":1808},{"delegate":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF","consensus_pkh":"tz4HRf3jYhfo8W3D4ppUCWjAi2rw7jPSZ234","consensus_power":1792},{"delegate":"tz3PgFHdYvEGEbUo1pUJmuNH8fgc8cwARKfC","consensus_pkh":"tz4EWkmNN93yE7HrjaRR6mGh22rUgSYJG1Sj","consensus_power":859},{"delegate":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_pkh":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_power":820},{"delegate":"tz4DLnQQ9XvmgrVGFZ72JFAckPcozpUWarXc","consensus_pkh":"tz4DLnQQ9XvmgrVGFZ72JFAckPcozpUWarXc","consensus_power":3},{"delegate":"tz1TGKSrZrBpND3PELJ43nVdyadoeiM1WMzb","consensus_pkh":"tz4Cv1k5zjnHFYB4QVpXcSFnegQQ1y7ozzXj","consensus_power":10},{"delegate":"tz4S6P3CqKKPxGwFXdDeLL1jLfrrs5WTcHk9","consensus_pkh":"tz4S6P3CqKKPxGwFXdDeLL1jLfrrs5WTcHk9","consensus_power":1},{"delegate":"tz1LmrwzCKUDibk7xaGC5RxvTbmUbCAtCA4a","consensus_pkh":"tz4QqpgVTG4CCETKMgV6YaUHyowE4Gkqwdfi","consensus_power":1}],"total_consensus_power":5294}}],"signature":"BLsigBavJFTR154dzBqF1We18A8eDgpinREtE8ayWohQmKXP3cwJoEBMhTa9EHnzyzsw6qGNVvoYkNiKD2azox6xJtA4cW5G4TUcUqgnwjPewfKG3wXqpVEn8wSSMWfCUmQYocYWxsy6EY"}],[],[],[]]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BKprStSDDvZQDu62cBQAtRnohsauwtMkQgBuKV3qv9StvRgkndd"))
	assert.Nil(t, e)
	assert.Equal(t, tezos.OpTypeAttestationsAggregate, value[0][0].Contents.N(0).Kind())
}
