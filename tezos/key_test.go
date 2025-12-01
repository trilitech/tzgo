// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package tezos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	type testcase struct {
		Address Address
		Priv    string
		Pub     string
		Pass    string
	}

	cases := []testcase{
		// ed25519 unencrypted
		{
			Priv:    "edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd",
			Pub:     "edpkv45regue1bWtuHnCgLU8xWKLwa9qRqv4gimgJKro4LSc3C5VjV",
			Address: MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q"),
		},
		// secp256k1 unencrypted
		{
			Priv:    "spsk2oTAhiaSywh9ctt8yZLRxL3bo8Mayd3hKFi5iBaoqj2R8bx7ow",
			Pub:     "sppk7auhfZa5wAcR8hk3WCw47kHgG3Pp8zaP3ctdAqdDd2dBAeZBof1",
			Address: MustParseAddress("tz2VN9n2C56xGLykHCjhNvZQqUeTVisrHjxA"),
		},
		// p256 unencrypted
		{
			Priv:    "p2sk35q9MJHLN1SBHNhKq7oho1vnZL28bYfsSKDUrDn2e4XVcp6ohZ",
			Pub:     "p2pk64zMPtYav6yiaHV2DhSQ65gbKMr3gkLQtK7TTQCpJEVUhxxEnxo",
			Address: MustParseAddress("tz3VCJEo1rRyyVejmpaRjbgGT9uE66sZmUtQ"),
		},
		// bls12_381 unencrypted
		{
			Priv:    "BLsk1eGhiPQXKtvvkBeXzmtVVJs6KPhEF45drF7MLjoCDcSnTGuyjL",
			Pub:     "BLpk1ur5XXicWYMMzCVZZWyLZhybtyX8Zot2uCzDCZW8KcC5BdZiLVXRZvZzi4GuZYL9SarUvKpE",
			Address: MustParseAddress("tz4TFJdv9Jd44FtBMAxi3KQT7AtazhVyaPa6"),
		},
		// ed25519 encrypted
		{
			Priv:    "edesk1uiM6BaysskGto8pRtzKQqFqsy1sea1QRjTzaQYuBxYNhuN6eqEU78TGRXZocsVRJYcN7AaU9JDykwUd8KW",
			Pub:     "edpkttVn1coEZNjcjjAF36jDXDB377imNiKCHqjdXSt85eVN779jfX",
			Address: MustParseAddress("tz1MKPxkZLfdw31LL7zi55aZEoyH9DPL7eh7"),
			Pass:    "foo",
		},
		// secp256k1 encrypted
		{
			Priv:    "spesk246GnDVaqGoYZvKbjrWM1g6xUXnyETXtwZgEYFnP8BQXcaS4rfQQco7C94D1yBmcL1v46Sqy8fXrhBSM7TW",
			Pub:     "sppk7aSJpAzeXNTaobig65si221WTqgPh8mJsCJSAiZU7asJkWBVGyx",
			Address: MustParseAddress("tz29QkiEM1xf3chaZj6DjL5udNLbUZ8d6QJ4"),
			Pass:    "foo",
		},
		// p256 encrypted
		{
			Priv:    "p2esk27ocLPLp1JkTWfxByXysGyB7MBDURYJAzAGJLR3XSEV9Nq8wFFdDVXVTwvCwR7Ne2dcUveamjXbvZf3on6T",
			Pub:     "p2pk66vAYU7rN1ckJMp38Z9pXCrkiZCVyi6KyeMwhY69h5WDPHdMecH",
			Address: MustParseAddress("tz3Qa3kjWa6B3XgvZcVe24gTfjkc5WZRz59Q"),
			Pass:    "foo",
		},
	}

	for i, c := range cases {
		if !IsPrivateKey(c.Priv) {
			t.Errorf("Case %d - Expected private key", i)
		}
		if !IsPublicKey(c.Pub) {
			t.Errorf("Case %d - Expected public key", i)
		}
		if c.Pass != "" && !IsEncryptedKey(c.Priv) {
			t.Errorf("Case %d - Expected encrypted key", i)
		}

		sk, err := ParseEncryptedPrivateKey(c.Priv, func() ([]byte, error) { return []byte(c.Pass), nil })
		if err != nil {
			t.Errorf("Case %d - Parsing key %s: %v", i, c.Priv, err)
		}
		if !sk.IsValid() {
			t.Errorf("Case %d - Expected valid key %s", i, c.Priv)
		}

		pk, err := ParseKey(c.Pub)
		if err != nil {
			t.Errorf("Case %d - Parsing pubkey %s: %v", i, c.Pub, err)
		}
		if !pk.IsValid() {
			t.Errorf("Case %d - Expected valid pubkey %s", i, c.Priv)
		}

		// generate pk from sk
		if check := sk.Public(); !check.IsEqual(pk) {
			t.Errorf("Case %d - Mismatch pk have=%s want=%s", i, check, pk)
		}

		// address from pk
		if got, want := pk.Address(), c.Address; !got.Equal(want) {
			t.Errorf("Case %d - Mismatch address got=%s want=%s", i, got, want)
		}
	}
}

func TestSign(t *testing.T) {
	type testcase struct {
		Priv        string
		Pub         string
		Msg         []byte
		Digest      string
		Sig         string
		Generic     string
		ExpectedSig string
	}

	cases := map[string]testcase{
		/*
			"ed25519 unencrypted": {
				Priv:        "edsk2qWwwqoVa2XCiXeaihYVg6BfLcicR1TwC6vf63dCDYPh3qjR2g",
				Pub:         "edpkvRy8feM9jRQaPNgtQYURAjWj3qrJtXkXp6hwfkDTtXDcfB2i5F",
				Msg:         []byte{5, 1, 0, 0, 0, 1, 97},
				ExpectedSig: "edsigtxuEw21vZaZYaamHwkdCb2XXVAG45hmaHE6QjEMGdNst3WJ1UYFK5R1HY7UtBJ7ZuU1ud52LPE7YPWPYSNbRDBArr39oAv",
			},
			"secp256k1 unencrypted": {
				Priv:        "spsk1nDmRj6hETy89DfJzHnmyFicx853ShpiLHLJAbg2Qu9gYdx35n",
				Pub:         "sppk7anmrSFCPfSKbm6GsARo1JRpethThozcxipErX4QtT8CBDojnaJ",
				Msg:         []byte{5, 1, 0, 0, 0, 1, 97},
				ExpectedSig: "spsig1DBG3ZMB5a7rwKMD4bXsxt7mD6QndfvZ6xATBAgdbajrnbohonsqYzLVQFWescq2JFF9PztcUbDaKeX89nxcXR7EYrHedF",
			},
			"p256 unencrypted": {
				Priv:        "p2sk3heCRmbfiArx8so4SBevK8t7mPGRqBN8eAYTzZJPWnu6LadRbM",
				Pub:         "p2pk66FyiYn3WDkJ5DEQptvaPy3gBEXGr7TTMFh94pZ5p3KALfzamqi",
				Msg:         []byte{5, 1, 0, 0, 0, 1, 97},
				ExpectedSig: "p2sigrpjz56TPL47L8mxcpdE9UpK2wZxuP2ZwkqezDXDNCMdLK4TGd9gTvbe2a9M9mj5nMWV7ZwxmwYe9S4mNc8JYRwNag39cB",
			},
		*/
		"bls unencrypted": {
			Priv:        "BLsk2rrqeqLp7ujt4NSCG8W6xDMFDBm6QHoTabeSQ4HjPDewaX4F6k",
			Pub:         "BLpk1wdBzZKshyhkdge3cXvWdTWhCWDsih8X1pbEdvjTapd1PvsESzTjMTwNWpephX8wyhshSFCp",
			Msg:         []byte{5, 1, 0, 0, 0, 1, 97},
			ExpectedSig: "BLsigBcuqVyAuyYiMhPmcf16B8BmvhujB8DPUAmgYb94ixJ9wkraLfxzJpt2eyWMePtzuRNhRWSF4LukEv39LxAi7nGiH83ihKc9jnyhjbLc76QKTs4h1sTzQQEKR15yF9tSyU39iEsyTx",
		},
	}

	for name, c := range cases {
		digest := Digest(c.Msg)
		sk := MustParsePrivateKey(c.Priv)
		assert.Equal(t, c.Priv, sk.String(), "Case %s - Private key round trip mismatch", name)
		pk := sk.Public()
		if sk.Type == KeyTypeBls12_381 {
			msg := append(pk.Data, c.Msg...)
			digest = Digest(msg)
		}
		sig, err := sk.Sign(digest[:])
		if err != nil {
			t.Errorf("Case %s - Signing failed: %v", name, err)
		}
		if !sig.IsValid() {
			t.Errorf("Case %s - Invalid signature %s", name, sig)
		}
		assert.Equal(t, c.ExpectedSig, sig.String(), "Case %s - Signature mismatch", name)
		if err := pk.Verify(digest[:], sig); err != nil {
			t.Errorf("Case %s - Verify failed %v", name, err)
		}
		if err := pk.Verify(digest[:], MustParseSignature(sig.Generic())); err != nil {
			t.Errorf("Case %s - Verify generic failed %v", name, err)
		}
	}
}
