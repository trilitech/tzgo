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

func TestGetContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		switch r.URL.Path {
		case "/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/context/contracts/tz2XbNrEZRJ8DeSGYbuZoRyDn1Qfj1rJoCLE":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"balance":"1000","counter":"741806","revealed":true}`))
		case "/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/context/contracts/KT18x7skHqt9hGYjrg3EJKceigfz1sJJPgZ8":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"balance":"1000000","storage": {"int": "0"}}`))
		default:
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	// Implicit accounts have the `Revealed` field
	value, e := c.GetContract(context.TODO(), tezos.MustParseAddress("tz2XbNrEZRJ8DeSGYbuZoRyDn1Qfj1rJoCLE"), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)
	assert.NotNil(t, value.Revealed)
	assert.True(t, *value.Revealed)

	// Regular contract accounts do not have the `Revealed` field
	value, e = c.GetContract(context.TODO(), tezos.MustParseAddress("KT18x7skHqt9hGYjrg3EJKceigfz1sJJPgZ8"), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)
	assert.Nil(t, value.Revealed)
}
