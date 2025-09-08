// Copyright (c) 2025 Trilitech Ltd.
// Author: tzstats@trili.tech

package rpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trilitech/tzgo/tezos"
)

func TestParseUpdateCompanionKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations" {
			t.Errorf("Expected to request '/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[[{"protocol":"PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh","chain_id":"NetXd56aBs1aeW3","hash":"ooqgVxC8XDYHXSnznWSdkXktgUaA2PZdUs4azRbco9i6fhiY1ui","branch": "BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb","contents": [{"kind": "update_companion_key","source": "tz4XsUYECWrJ5c4gktAnr1nkFC9mjDEVn3qh","fee": "33","counter": "732","gas_limit": "9451117","storage_limit": "57024931117","pk": "BLpk1xXdveUYh7YFsyf6LwGWfv5zAfLvnMG71byiMFDZc4CkXzZPVko3Dz4sD43Ln5uFNvdjiQJY","proof": "BLsigA5hv9bjFDmrbvVfam6CaXM87xZv1GYhFmn3uaR4LTXahHpADaNk9LcEcQb2EF1w9b72AutNqPxzVT5EhRFQN1Ue43XA3HkMfmCYTCHa3WiLsJYLmuQvchTwoMy8zx8EibsZAJ3Dug"}],"signature": "BLsigAV587qg1S78p9hRyfmcqM6cVxsafJ9MVRWLLxfykmdB7LcpMWgLZjSQfhbuffvJovUsS9vJpesiwJm9xEhZPjW93DB5tjSUkUrXcQpTVGoJe4KsGR16BUdVEZQuXeW2XjUk9BqYSP"}],[],[],[]]`))
	}))
	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)
	assert.Len(t, value, 4)
	assert.Len(t, value[0], 1)
	assert.Equal(t, value[0][0].Contents.Len(), 1)
	op := value[0][0].Contents.N(0).(*UpdateConsensusKey)
	assert.Equal(t, tezos.OpTypeUpdateCompanionKey, op.Kind())
	assert.Equal(t, int64(732), op.Counter)
	assert.Equal(t, tezos.MustParseAddress("tz4XsUYECWrJ5c4gktAnr1nkFC9mjDEVn3qh"), op.Source)
}
