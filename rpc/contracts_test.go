// Copyright (c) 2025 Trilitech Ltd.
// Author: tzstats@trili.tech

package rpc

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trilitech/tzgo/tezos"
)

// Offline fixtures for contract endpoints.
//
//go:embed testdata/mainnet_contract_KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd.json
var mainnetContractJSON []byte

//go:embed testdata/mainnet_contract_KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd_balance.json
var mainnetContractBalanceJSON []byte

//go:embed testdata/mainnet_contract_KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd_script.json
var mainnetContractScriptJSON []byte

//go:embed testdata/mainnet_contract_KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd_storage.json
var mainnetContractStorageJSON []byte

//go:embed testdata/mainnet_contract_KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd_entrypoints.json
var mainnetContractEntrypointsJSON []byte

//go:embed testdata/mainnet_bigmap_515_info_raw.json
var mainnetBigmap515InfoRawJSON []byte

//go:embed testdata/mainnet_bigmap_515_keys_raw.json
var mainnetBigmap515KeysRawJSON []byte

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

// TestContractsAPI_FromTestdataFixtures tests contract-related RPC client methods
// using embedded mainnet fixtures from examples/contract and examples/rpc.
func TestContractsAPI_FromTestdataFixtures(t *testing.T) {
	contractAddr := tezos.MustParseAddress("KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd")
	bigmapID := int64(515)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "application/json", r.Header.Get("Accept"))

		switch r.URL.Path {
		case "/chains/main/blocks/head/context/contracts/KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetContractJSON)
		case "/chains/main/blocks/head/context/contracts/KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd/balance":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetContractBalanceJSON)
		case "/chains/main/blocks/head/context/contracts/KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd/script":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetContractScriptJSON)
		case "/chains/main/blocks/head/context/contracts/KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd/storage":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetContractStorageJSON)
		case "/chains/main/blocks/head/context/contracts/KT1LRboPna9yQY9BrjtQYDS1DVxhKESK4VVd/entrypoints":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetContractEntrypointsJSON)
		case "/chains/main/blocks/head/context/raw/json/big_maps/index/515":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBigmap515InfoRawJSON)
		case "/chains/main/blocks/head/context/raw/json/big_maps/index/515/contents":
			w.WriteHeader(http.StatusOK)
			w.Write(mainnetBigmap515KeysRawJSON)
		default:
			t.Errorf("Unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c, err := NewClient(server.URL, nil)
	require.NoError(t, err)
	ctx := context.Background()

	t.Run("GetContract", func(t *testing.T) {
		contract, err := c.GetContract(ctx, contractAddr, Head)
		require.NoError(t, err)
		require.NotNil(t, contract)
		assert.Equal(t, int64(0), contract.Balance)
		assert.Equal(t, int64(0), contract.Counter)
	})

	t.Run("GetContractBalance", func(t *testing.T) {
		balance, err := c.GetContractBalance(ctx, contractAddr, Head)
		require.NoError(t, err)
		assert.Equal(t, "0", balance.String())
	})

	t.Run("GetContractScript", func(t *testing.T) {
		script, err := c.GetContractScript(ctx, contractAddr)
		require.NoError(t, err)
		require.NotNil(t, script)
		// Verify script has code and storage
		assert.NotNil(t, script.Code, "Script should have code")
		assert.NotNil(t, script.Storage, "Script should have storage")
	})

	t.Run("GetContractStorage", func(t *testing.T) {
		storage, err := c.GetContractStorage(ctx, contractAddr, Head)
		require.NoError(t, err)
		assert.NotNil(t, storage)
		// Verify storage is a valid Micheline primitive
		assert.True(t, storage.IsValid(), "Storage should be a valid Micheline primitive")
	})

	t.Run("GetContractEntrypoints", func(t *testing.T) {
		entrypoints, err := c.GetContractEntrypoints(ctx, contractAddr)
		require.NoError(t, err)
		if entrypoints == nil {
			t.Log("Entrypoints is nil")
			return
		}
		assert.Greater(t, len(entrypoints), 0, "Contract should have at least one entrypoint")
		// Verify "admin" entrypoint exists (from the fixture)
		adminEp, hasAdmin := entrypoints["admin"]
		assert.True(t, hasAdmin, "Contract should have an 'admin' entrypoint")
		if hasAdmin {
			assert.True(t, adminEp.IsValid(), "Admin entrypoint type should be valid")
		}
	})

	t.Run("GetBigmapInfo", func(t *testing.T) {
		info, err := c.GetBigmapInfo(ctx, bigmapID, Head)
		require.NoError(t, err)
		require.NotNil(t, info)
		assert.True(t, info.KeyType.IsValid(), "KeyType should be valid")
		assert.True(t, info.ValueType.IsValid(), "ValueType should be valid")
		assert.GreaterOrEqual(t, info.TotalBytes, int64(0), "TotalBytes should be non-negative")
	})

	t.Run("ListBigmapKeys", func(t *testing.T) {
		keys, err := c.ListBigmapKeys(ctx, bigmapID, Head)
		require.NoError(t, err)
		require.NotNil(t, keys)
		assert.Greater(t, len(keys), 0, "Bigmap should have at least one key")
		// Verify keys are valid expression hashes
		for _, key := range keys {
			assert.True(t, key.IsValid(), "Each key should be a valid expression hash")
		}
	})
}
