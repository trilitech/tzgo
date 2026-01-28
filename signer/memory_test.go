package signer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/trilitech/tzgo/codec"
	"github.com/trilitech/tzgo/signer"
	"github.com/trilitech/tzgo/tezos"
)

func TestMemorySigner_ListAddresses(t *testing.T) {
	sk := tezos.MustParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd")
	s := signer.NewFromKey(sk)

	addrs, err := s.ListAddresses(context.Background())
	if err != nil {
		t.Fatalf("ListAddresses err=%v", err)
	}
	if len(addrs) != 1 {
		t.Fatalf("ListAddresses len=%d want 1", len(addrs))
	}
	if got, want := addrs[0], sk.Address(); !got.Equal(want) {
		t.Fatalf("ListAddresses[0]=%s want %s", got, want)
	}
}

func TestMemorySigner_GetKey(t *testing.T) {
	sk := tezos.MustParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd")
	s := signer.NewFromKey(sk)
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		pk, err := s.GetKey(ctx, sk.Address())
		if err != nil {
			t.Fatalf("GetKey err=%v", err)
		}
		if !pk.IsValid() {
			t.Fatalf("GetKey returned invalid key")
		}
		if got, want := pk, sk.Public(); !got.IsEqual(want) {
			t.Fatalf("GetKey mismatch got=%s want=%s", got, want)
		}
	})

	t.Run("address_mismatch", func(t *testing.T) {
		other := tezos.MustParseAddress("tz1MKPxkZLfdw31LL7zi55aZEoyH9DPL7eh7")
		pk, err := s.GetKey(ctx, other)
		if !errors.Is(err, signer.ErrAddressMismatch) {
			t.Fatalf("GetKey error=%v want ErrAddressMismatch", err)
		}
		if pk.IsValid() {
			t.Fatalf("GetKey returned valid key on mismatch: %s", pk)
		}
		if pk.Type != tezos.KeyTypeInvalid || len(pk.Data) != 0 {
			t.Fatalf("GetKey returned non-invalid key on mismatch: type=%v len=%d", pk.Type, len(pk.Data))
		}
	})
}

func TestMemorySigner_SignMessage(t *testing.T) {
	sk := tezos.MustParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd")
	pk := sk.Public()
	s := signer.NewFromKey(sk)
	ctx := context.Background()

	cases := []struct {
		name    string
		addr    tezos.Address
		msg     string
		wantErr error
	}{
		{
			name: "ok_nonempty",
			addr: sk.Address(),
			msg:  "hello",
		},
		{
			name: "ok_empty",
			addr: sk.Address(),
			msg:  "",
		},
		{
			name:    "address_mismatch",
			addr:    tezos.MustParseAddress("tz1MKPxkZLfdw31LL7zi55aZEoyH9DPL7eh7"),
			msg:     "hello",
			wantErr: signer.ErrAddressMismatch,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			sig, err := s.SignMessage(ctx, tt.addr, tt.msg)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("SignMessage err=%v want %v", err, tt.wantErr)
				}
				if sig.IsValid() {
					t.Fatalf("SignMessage returned valid signature on error: %s", sig)
				}
				if sig.Type != tezos.SignatureTypeInvalid || len(sig.Data) != 0 {
					t.Fatalf("SignMessage returned non-invalid signature on error: type=%v len=%d", sig.Type, len(sig.Data))
				}
				return
			}
			if err != nil {
				t.Fatalf("SignMessage err=%v", err)
			}
			if !sig.IsValid() {
				t.Fatalf("SignMessage returned invalid signature")
			}

			// Must match MemorySigner implementation: sign Digest(op.Bytes()) where op is a failing noop.
			op := codec.NewOp().
				WithBranch(tezos.ZeroBlockHash).
				WithContents(&codec.FailingNoop{Arbitrary: tt.msg})
			d := tezos.Digest(op.Bytes())

			if err := pk.Verify(d[:], sig); err != nil {
				t.Fatalf("signature verification failed: %v", err)
			}
		})
	}
}

func TestMemorySigner_SignOperation(t *testing.T) {
	sk := tezos.MustParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd")
	pk := sk.Public()
	s := signer.NewFromKey(sk)
	ctx := context.Background()

	t.Run("ok_signs_digest", func(t *testing.T) {
		branch := tezos.MustParseBlockHash("BKjS7rtCjysnMNWUuevZiF2a6NkUas9bnSsNQ3ibh5GfKNrQoGk")
		op := codec.NewOp().
			WithBranch(branch).
			WithContents(&codec.FailingNoop{Arbitrary: "hello"})

		sig, err := s.SignOperation(ctx, sk.Address(), op)
		if err != nil {
			t.Fatalf("SignOperation err=%v", err)
		}
		if !sig.IsValid() {
			t.Fatalf("SignOperation returned invalid signature")
		}
		if !op.Signature.IsValid() {
			t.Fatalf("operation signature not set")
		}
		if !sig.Equal(op.Signature) {
			t.Fatalf("returned signature differs from op.Signature")
		}
		if err := pk.Verify(op.Digest(), sig); err != nil {
			t.Fatalf("signature verification failed: %v", err)
		}
	})

	t.Run("ok_noop_when_already_signed", func(t *testing.T) {
		branch := tezos.MustParseBlockHash("BKjS7rtCjysnMNWUuevZiF2a6NkUas9bnSsNQ3ibh5GfKNrQoGk")
		op := codec.NewOp().
			WithBranch(branch).
			WithContents(&codec.FailingNoop{Arbitrary: "hello"})
		if err := op.Sign(sk); err != nil {
			t.Fatalf("pre-sign err=%v", err)
		}
		prev := op.Signature

		sig, err := s.SignOperation(ctx, sk.Address(), op)
		if err != nil {
			t.Fatalf("SignOperation err=%v", err)
		}
		if !sig.Equal(prev) || !op.Signature.Equal(prev) {
			t.Fatalf("expected signature unchanged when already signed")
		}
		if err := pk.Verify(op.Digest(), sig); err != nil {
			t.Fatalf("signature verification failed: %v", err)
		}
	})

	t.Run("address_mismatch", func(t *testing.T) {
		branch := tezos.MustParseBlockHash("BKjS7rtCjysnMNWUuevZiF2a6NkUas9bnSsNQ3ibh5GfKNrQoGk")
		op := codec.NewOp().
			WithBranch(branch).
			WithContents(&codec.FailingNoop{Arbitrary: "hello"})
		other := tezos.MustParseAddress("tz1MKPxkZLfdw31LL7zi55aZEoyH9DPL7eh7")

		sig, err := s.SignOperation(ctx, other, op)
		if !errors.Is(err, signer.ErrAddressMismatch) {
			t.Fatalf("SignOperation err=%v want ErrAddressMismatch", err)
		}
		if sig.IsValid() {
			t.Fatalf("SignOperation returned valid signature on mismatch: %s", sig)
		}
		if sig.Type != tezos.SignatureTypeInvalid || len(sig.Data) != 0 {
			t.Fatalf("SignOperation returned non-invalid signature on mismatch: type=%v len=%d", sig.Type, len(sig.Data))
		}
	})

	t.Run("missing_branch_propagates_error", func(t *testing.T) {
		op := codec.NewOp().WithContents(&codec.FailingNoop{Arbitrary: "hello"})
		sig, err := s.SignOperation(ctx, sk.Address(), op)
		if err == nil {
			t.Fatalf("SignOperation expected error for missing branch")
		}
		if sig.IsValid() {
			t.Fatalf("SignOperation returned valid signature on error: %s", sig)
		}
	})

	t.Run("empty_contents_propagates_error", func(t *testing.T) {
		op := codec.NewOp().WithBranch(tezos.ZeroBlockHash)
		sig, err := s.SignOperation(ctx, sk.Address(), op)
		if err == nil {
			t.Fatalf("SignOperation expected error for empty contents")
		}
		if sig.IsValid() {
			t.Fatalf("SignOperation returned valid signature on error: %s", sig)
		}
	})
}

func TestMemorySigner_SignBlock(t *testing.T) {
	sk := tezos.MustParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd")
	pk := sk.Public()
	s := signer.NewFromKey(sk)
	ctx := context.Background()

	t.Run("panics_when_chain_id_missing", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic when ChainId is nil")
			}
		}()
		var head codec.BlockHeader
		_, _ = s.SignBlock(ctx, sk.Address(), &head)
	})

	t.Run("ok_signs_when_chain_id_present", func(t *testing.T) {
		chain := tezos.MustParseChainIdHash("NetXdQprcVkpaWU")
		head := &codec.BlockHeader{
			Level:            1,
			Proto:            1,
			Timestamp:        time.Now().UTC(),
			ValidationPass:   4,
			PayloadRound:     0,
			ProofOfWorkNonce: make([]byte, 8),
		}
		head.WithChainId(chain)

		sig, err := s.SignBlock(ctx, sk.Address(), head)
		if err != nil {
			t.Fatalf("SignBlock err=%v", err)
		}
		if !sig.IsValid() {
			t.Fatalf("SignBlock returned invalid signature")
		}
		if sig.Type != tezos.SignatureTypeGeneric {
			t.Fatalf("SignBlock signature type=%v want generic", sig.Type)
		}
		if !sig.Equal(head.Signature) {
			t.Fatalf("returned signature differs from head.Signature")
		}
		// Verify against the digest that was signed: BlockHeader.Sign() signs the
		// header digest computed *before* the signature is set, while Digest()
		// currently includes the signature bytes once present.
		unsigned := *head
		unsigned.Signature = tezos.InvalidSignature
		if err := pk.Verify(unsigned.Digest(), sig); err != nil {
			t.Fatalf("signature verification failed: %v", err)
		}
	})

	t.Run("address_mismatch", func(t *testing.T) {
		chain := tezos.MustParseChainIdHash("NetXdQprcVkpaWU")
		head := &codec.BlockHeader{ProofOfWorkNonce: make([]byte, 8)}
		head.WithChainId(chain)
		other := tezos.MustParseAddress("tz1MKPxkZLfdw31LL7zi55aZEoyH9DPL7eh7")

		sig, err := s.SignBlock(ctx, other, head)
		if !errors.Is(err, signer.ErrAddressMismatch) {
			t.Fatalf("SignBlock err=%v want ErrAddressMismatch", err)
		}
		if sig.IsValid() {
			t.Fatalf("SignBlock returned valid signature on mismatch: %s", sig)
		}
		if sig.Type != tezos.SignatureTypeInvalid || len(sig.Data) != 0 {
			t.Fatalf("SignBlock returned non-invalid signature on mismatch: type=%v len=%d", sig.Type, len(sig.Data))
		}
	})
}
