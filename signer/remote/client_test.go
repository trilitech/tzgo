package remote_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/trilitech/tzgo/codec"
	"github.com/trilitech/tzgo/signer"
	"github.com/trilitech/tzgo/signer/remote"
	"github.com/trilitech/tzgo/tezos"
)

func TestRemoteSigner_ListAddresses_WithAddress_WithAuthKey(t *testing.T) {
	s, err := remote.New("http://example.com", http.DefaultClient)
	if err != nil {
		t.Fatalf("New err=%v", err)
	}
	a1 := tezos.MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q")
	a2 := tezos.MustParseAddress("tz1MKPxkZLfdw31LL7zi55aZEoyH9DPL7eh7")
	s.WithAddress(a1).WithAddress(a2)

	// WithAuthKey is a public method and must not break basic behavior.
	s.WithAuthKey(tezos.MustParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd"))

	addrs, err := s.ListAddresses(context.Background())
	if err != nil {
		t.Fatalf("ListAddresses err=%v", err)
	}
	if len(addrs) != 2 {
		t.Fatalf("ListAddresses len=%d want 2", len(addrs))
	}
	if !addrs[0].Equal(a1) || !addrs[1].Equal(a2) {
		t.Fatalf("ListAddresses order mismatch got=%v", addrs)
	}
}

func TestRemoteSigner_AuthorizedKeys(t *testing.T) {
	want := []string{
		"tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q",
		"tz1MKPxkZLfdw31LL7zi55aZEoyH9DPL7eh7",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method=%s want GET", r.Method)
		}
		if r.URL.Path != "/authorized_keys" {
			t.Fatalf("path=%s want /authorized_keys", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"authorized_keys": want,
		})
	}))
	defer srv.Close()

	s, err := remote.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New err=%v", err)
	}
	addrs, err := s.AuthorizedKeys(context.Background())
	if err != nil {
		t.Fatalf("AuthorizedKeys err=%v", err)
	}
	if len(addrs) != len(want) {
		t.Fatalf("AuthorizedKeys len=%d want %d", len(addrs), len(want))
	}
	for i := range want {
		if got := addrs[i].String(); got != want[i] {
			t.Fatalf("AuthorizedKeys[%d]=%q want %q", i, got, want[i])
		}
	}
}

func TestRemoteSigner_GetKey(t *testing.T) {
	addr := tezos.MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q")
	pub := "edpkv45regue1bWtuHnCgLU8xWKLwa9qRqv4gimgJKro4LSc3C5VjV"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method=%s want GET", r.Method)
		}
		if r.URL.Path != "/keys/"+addr.String() {
			t.Fatalf("path=%s want /keys/%s", r.URL.Path, addr)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"public_key": pub,
		})
	}))
	defer srv.Close()

	s, err := remote.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New err=%v", err)
	}
	got, err := s.GetKey(context.Background(), addr)
	if err != nil {
		t.Fatalf("GetKey err=%v", err)
	}
	if !got.IsValid() {
		t.Fatalf("GetKey returned invalid key")
	}
	if got.String() != pub {
		t.Fatalf("GetKey=%q want %q", got.String(), pub)
	}
}

func TestRemoteSigner_SignOperation_HTTPContract(t *testing.T) {
	ctx := context.Background()
	addr := tezos.MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q")
	sk := tezos.MustParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd")
	pk := sk.Public()

	branch := tezos.MustParseBlockHash("BKjS7rtCjysnMNWUuevZiF2a6NkUas9bnSsNQ3ibh5GfKNrQoGk")
	op := codec.NewOp().
		WithBranch(branch).
		WithContents(&codec.FailingNoop{Arbitrary: "hello"})
	wantHex := hex.EncodeToString(op.WatermarkedBytes())
	if !strings.HasPrefix(wantHex, "03") {
		t.Fatalf("expected watermarked op to start with 0x03, got %s", wantHex[:2])
	}
	wantSig, err := sk.Sign(op.Digest())
	if err != nil {
		t.Fatalf("Sign digest err=%v", err)
	}

	var gotBody string
	var calls int32
	var handlerErr atomic.Value // stores error
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		if r.Method != http.MethodPost {
			handlerErr.Store(fmt.Errorf("method=%s want POST", r.Method))
			http.Error(w, "bad method", http.StatusBadRequest)
			return
		}
		if r.URL.Path != "/keys/"+addr.String() {
			handlerErr.Store(fmt.Errorf("path=%s want /keys/%s", r.URL.Path, addr))
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		b, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(b, &gotBody); err != nil {
			handlerErr.Store(fmt.Errorf("failed to decode json body %q: %w", string(b), err))
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"signature": wantSig.String(),
		})
	}))
	defer srv.Close()

	s, err := remote.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New err=%v", err)
	}
	gotSig, err := s.SignOperation(ctx, addr, op)
	if err != nil {
		t.Fatalf("SignOperation err=%v", err)
	}
	if v := handlerErr.Load(); v != nil {
		t.Fatalf("server handler error: %v", v.(error))
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("server call count=%d want 1", got)
	}
	if gotBody != wantHex {
		t.Fatalf("request body hex mismatch\n got: %s\nwant: %s", gotBody, wantHex)
	}
	if !gotSig.Equal(wantSig) {
		t.Fatalf("signature mismatch got=%s want=%s", gotSig, wantSig)
	}
	if err := pk.Verify(op.Digest(), gotSig); err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func TestRemoteSigner_SignMessage_HTTPContract(t *testing.T) {
	ctx := context.Background()
	addr := tezos.MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q")
	sk := tezos.MustParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd")
	pk := sk.Public()

	msg := "arbitrary message"

	// RemoteSigner.SignMessage wraps msg into a failing noop operation with
	// tezos.ZeroBlockHash. Note that tezos.ZeroBlockHash is considered invalid,
	// so op.WatermarkedBytes() is nil and the resulting hex string is empty.
	op := codec.NewOp().
		WithBranch(tezos.ZeroBlockHash).
		WithContents(&codec.FailingNoop{Arbitrary: msg})
	wantHex := hex.EncodeToString(op.WatermarkedBytes())

	var gotBody string
	var calls int32
	var handlerErr atomic.Value // stores error
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		if r.Method != http.MethodPost {
			handlerErr.Store(fmt.Errorf("method=%s want POST", r.Method))
			http.Error(w, "bad method", http.StatusBadRequest)
			return
		}
		if r.URL.Path != "/keys/"+addr.String() {
			handlerErr.Store(fmt.Errorf("path=%s want /keys/%s", r.URL.Path, addr))
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		b, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(b, &gotBody); err != nil {
			handlerErr.Store(fmt.Errorf("failed to decode json body %q: %w", string(b), err))
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}

		// Sign exactly what the client sent (hex-decoded).
		payload, err := hex.DecodeString(gotBody)
		if err != nil {
			handlerErr.Store(fmt.Errorf("failed to decode hex %q: %w", gotBody, err))
			http.Error(w, "bad hex", http.StatusBadRequest)
			return
		}
		d := tezos.Digest(payload)
		sig, err := sk.Sign(d[:])
		if err != nil {
			handlerErr.Store(fmt.Errorf("failed to sign digest: %w", err))
			http.Error(w, "sign error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"signature": sig.String(),
		})
	}))
	defer srv.Close()

	s, err := remote.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New err=%v", err)
	}
	gotSig, err := s.SignMessage(ctx, addr, msg)
	if err != nil {
		t.Fatalf("SignMessage err=%v", err)
	}
	if v := handlerErr.Load(); v != nil {
		t.Fatalf("server handler error: %v", v.(error))
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("server call count=%d want 1", got)
	}
	if gotBody != wantHex {
		t.Fatalf("request body hex mismatch\n got: %q\nwant: %q", gotBody, wantHex)
	}

	// Verify returned signature against the digest of the payload actually sent.
	payload, err := hex.DecodeString(gotBody)
	if err != nil {
		t.Fatalf("decode hex err=%v", err)
	}
	d := tezos.Digest(payload)
	if err := pk.Verify(d[:], gotSig); err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func TestRemoteSigner_SignBlock_HTTPContract(t *testing.T) {
	ctx := context.Background()
	addr := tezos.MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q")
	sk := tezos.MustParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd")
	pk := sk.Public()

	chain := tezos.MustParseChainIdHash("NetXdQprcVkpaWU")
	head := &codec.BlockHeader{
		Level:            1,
		Proto:            1,
		Timestamp:        time.Unix(0, 0).UTC(),
		ValidationPass:   4,
		PayloadRound:     0,
		ProofOfWorkNonce: make([]byte, 8),
	}
	head.WithChainId(chain)

	wantHex := hex.EncodeToString(head.WatermarkedBytes())
	if !strings.HasPrefix(wantHex, "11"+hex.EncodeToString(chain.Bytes())) {
		t.Fatalf("expected watermarked block to start with 0x11 + chain_id bytes")
	}
	wantSig, err := sk.Sign(head.Digest())
	if err != nil {
		t.Fatalf("Sign digest err=%v", err)
	}

	var gotBody string
	var calls int32
	var handlerErr atomic.Value // stores error
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		if r.Method != http.MethodPost {
			handlerErr.Store(fmt.Errorf("method=%s want POST", r.Method))
			http.Error(w, "bad method", http.StatusBadRequest)
			return
		}
		if r.URL.Path != "/keys/"+addr.String() {
			handlerErr.Store(fmt.Errorf("path=%s want /keys/%s", r.URL.Path, addr))
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		b, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(b, &gotBody); err != nil {
			handlerErr.Store(fmt.Errorf("failed to decode json body %q: %w", string(b), err))
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"signature": wantSig.String(),
		})
	}))
	defer srv.Close()

	s, err := remote.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New err=%v", err)
	}
	gotSig, err := s.SignBlock(ctx, addr, head)
	if err != nil {
		t.Fatalf("SignBlock err=%v", err)
	}
	if v := handlerErr.Load(); v != nil {
		t.Fatalf("server handler error: %v", v.(error))
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("server call count=%d want 1", got)
	}
	if gotBody != wantHex {
		t.Fatalf("request body hex mismatch\n got: %s\nwant: %s", gotBody, wantHex)
	}
	if !gotSig.Equal(wantSig) {
		t.Fatalf("signature mismatch got=%s want=%s", gotSig, wantSig)
	}
	if err := pk.Verify(head.Digest(), gotSig); err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func TestRemoteSigner_ErrorOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "nope", http.StatusBadRequest)
	}))
	defer srv.Close()

	s, err := remote.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New err=%v", err)
	}
	_, err = s.AuthorizedKeys(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	// Ensure error contains HTTP status information (rpc client wraps this).
	if !strings.Contains(err.Error(), "400") && !strings.Contains(strings.ToLower(err.Error()), "bad request") {
		t.Fatalf("error=%q expected to mention status", err.Error())
	}

	// Ensure error is not spuriously wrapped as ErrAddressMismatch.
	if errors.Is(err, signer.ErrAddressMismatch) {
		t.Fatalf("unexpected address mismatch error")
	}
}
