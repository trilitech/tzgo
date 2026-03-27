package rpc

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trilitech/tzgo/tezos"
)

//go:embed testdata/mainnet_chain_id.json
var mainnetChainIDJSON []byte

//go:embed testdata/mainnet_is_bootstrapped.json
var mainnetIsBootstrappedJSON []byte

//go:embed testdata/mainnet_version.json
var mainnetVersionJSON []byte

func TestChainAPI_FromTestdataFixtures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("missing Accept: application/json"))
			return
		}
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/chains/main/chain_id":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(mainnetChainIDJSON)
		case "/chains/main/is_bootstrapped":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(mainnetIsBootstrappedJSON)
		case "/version":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(mainnetVersionJSON)
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("unexpected url: " + r.URL.String()))
		}
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)

	chainID, err := c.GetChainId(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, tezos.MustParseChainIdHash("NetXdQprcVkpaWU"), chainID)

	status, err := c.GetStatus(context.Background())
	assert.NoError(t, err)
	assert.True(t, status.Bootstrapped)
	assert.Equal(t, "synced", status.SyncState)

	ver, err := c.GetVersionInfo(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 24, ver.NodeVersion.Major)
	assert.Equal(t, 0, ver.NodeVersion.Minor)
	assert.Equal(t, "TEZOS_MAINNET", ver.NetworkVersion.ChainName)
	assert.Equal(t, 2, ver.NetworkVersion.DistributedDbVersion)
	assert.Equal(t, 1, ver.NetworkVersion.P2pVersion)
	assert.Equal(t, "5a3ca147", ver.CommitInfo.CommitHash)
}
