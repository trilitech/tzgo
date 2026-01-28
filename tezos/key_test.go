// Copyright (c) 2020-2022 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package tezos

import (
	"bytes"
	"encoding"
	"testing"
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
		// {
		//     Priv:    "BLsk1eGhiPQXKtvvkBeXzmtVVJs6KPhEF45drF7MLjoCDcSnTGuyjL",
		//     Pub:     "BLpk1ur5XXicWYMMzCVZZWyLZhybtyX8Zot2uCzDCZW8KcC5BdZiLVXRZvZzi4GuZYL9SarUvKpE",
		//     Address: MustParseAddress("tz4TFJdv9Jd44FtBMAxi3KQT7AtazhVyaPa6"),
		// },
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
		Priv    string
		Pub     string
		Msg     string
		Digest  string
		Sig     string
		Generic string
	}

	cases := []testcase{
		// ed25519 unencrypted
		{
			Priv: "edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd",
			Pub:  "edpkv45regue1bWtuHnCgLU8xWKLwa9qRqv4gimgJKro4LSc3C5VjV",
			Msg:  "hello",
		},
		// secp256k1 unencrypted
		{
			Priv: "spsk2oTAhiaSywh9ctt8yZLRxL3bo8Mayd3hKFi5iBaoqj2R8bx7ow",
			Pub:  "sppk7auhfZa5wAcR8hk3WCw47kHgG3Pp8zaP3ctdAqdDd2dBAeZBof1",
			Msg:  "hello",
		},
		// p256 unencrypted
		{
			Priv: "p2sk35q9MJHLN1SBHNhKq7oho1vnZL28bYfsSKDUrDn2e4XVcp6ohZ",
			Pub:  "p2pk64zMPtYav6yiaHV2DhSQ65gbKMr3gkLQtK7TTQCpJEVUhxxEnxo",
			Msg:  "hello",
		},
	}

	for i, c := range cases {
		digest := Digest([]byte(c.Msg))
		sk := MustParsePrivateKey(c.Priv)
		pk := sk.Public()
		sig, err := sk.Sign(digest[:])
		if err != nil {
			t.Errorf("Case %d - Signing failed: %v", i, err)
		}
		if !sig.IsValid() {
			t.Errorf("Case %d - Invalid signature %s", i, sig)
		}
		if err := pk.Verify(digest[:], sig); err != nil {
			t.Errorf("Case %d - Verify failed %v", i, err)
		}
		if err := pk.Verify(digest[:], MustParseSignature(sig.Generic())); err != nil {
			t.Errorf("Case %d - Verify generic failed %v", i, err)
		}
	}
}

func TestKey_ExtraOffchainHelpers(t *testing.T) {
	pub := "edpkv45regue1bWtuHnCgLU8xWKLwa9qRqv4gimgJKro4LSc3C5VjV"
	priv := "edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd"
	encPriv := "edesk1uiM6BaysskGto8pRtzKQqFqsy1sea1QRjTzaQYuBxYNhuN6eqEU78TGRXZocsVRJYcN7AaU9JDykwUd8KW"

	// prefix checks
	if !HasKeyPrefix(pub) || !HasKeyPrefix(priv) || !HasKeyPrefix(encPriv) {
		t.Fatalf("HasKeyPrefix expected true for known keys")
	}
	if HasKeyPrefix("nope") {
		t.Fatalf("HasKeyPrefix(nope) expected false")
	}

	// ParseKeyType
	if kt, enc := ParseKeyType("edpk"); kt != KeyTypeInvalid || enc {
		t.Fatalf("ParseKeyType(edpk) expected invalid (not a sk prefix)")
	}
	if kt, enc := ParseKeyType(ED25519_ENCRYPTED_SEED_PREFIX); kt != KeyTypeEd25519 || !enc {
		t.Fatalf("ParseKeyType(edesk) mismatch got=%v enc=%v", kt, enc)
	}
	if kt, enc := ParseKeyType(ED25519_SEED_PREFIX); kt != KeyTypeEd25519 || enc {
		t.Fatalf("ParseKeyType(edsk) mismatch got=%v enc=%v", kt, enc)
	}

	// NewKey + Clone (deep copy)
	k, err := ParseKey(pub)
	if err != nil || !k.IsValid() {
		t.Fatalf("ParseKey err=%v", err)
	}
	k2 := NewKey(k.Type, k.Data)
	if !k2.IsEqual(k) {
		t.Fatalf("NewKey mismatch")
	}
	clone := k.Clone()
	if !clone.IsEqual(k) {
		t.Fatalf("Clone mismatch")
	}
	clone.Data[0] ^= 0xff
	if clone.IsEqual(k) {
		t.Fatalf("Clone unexpectedly shares backing array")
	}

	// String/MarshalText/UnmarshalText
	if k.String() != pub {
		t.Fatalf("Key.String mismatch got=%q want=%q", k.String(), pub)
	}
	var _ encoding.TextMarshaler = k
	var _ encoding.TextUnmarshaler = (*Key)(nil)
	tb, err := k.MarshalText()
	if err != nil || string(tb) != pub {
		t.Fatalf("MarshalText err=%v got=%q", err, string(tb))
	}
	var k3 Key
	if err := k3.UnmarshalText([]byte(pub)); err != nil || !k3.IsEqual(k) {
		t.Fatalf("UnmarshalText err=%v equal=%v", err, k3.IsEqual(k))
	}

	// Bytes/MarshalBinary/UnmarshalBinary + DecodeKey
	b := k.Bytes()
	if len(b) == 0 {
		t.Fatalf("Bytes returned empty")
	}
	bin, err := k.MarshalBinary()
	if err != nil || !bytes.Equal(bin, b) {
		t.Fatalf("MarshalBinary err=%v", err)
	}
	var k4 Key
	if err := k4.UnmarshalBinary(bin); err != nil || !k4.IsEqual(k) {
		t.Fatalf("UnmarshalBinary err=%v", err)
	}
	if _, err := DecodeKey(nil); err != nil {
		t.Fatalf("DecodeKey(nil) err=%v", err)
	}
	if _, err := DecodeKey([]byte{0x01}); err == nil {
		t.Fatalf("DecodeKey(short) expected error")
	}
	if err := k4.UnmarshalBinary([]byte{0x01}); err == nil {
		t.Fatalf("Key.UnmarshalBinary(short) expected error")
	}
	if err := k4.UnmarshalBinary(append([]byte{0xff}, make([]byte, 33)...)); err == nil {
		t.Fatalf("Key.UnmarshalBinary(invalid tag) expected error")
	}

	// MustParseKey + Set
	km := MustParseKey(pub)
	if !km.IsEqual(k) {
		t.Fatalf("MustParseKey mismatch")
	}
	var ks Key
	if err := ks.Set(pub); err != nil || !ks.IsEqual(k) {
		t.Fatalf("Key.Set err=%v", err)
	}

	// PrivateKey String/MarshalText/UnmarshalText/Address
	sk, err := ParsePrivateKey(priv)
	if err != nil || !sk.IsValid() {
		t.Fatalf("ParsePrivateKey err=%v", err)
	}
	if sk.String() != priv {
		t.Fatalf("PrivateKey.String mismatch got=%q want=%q", sk.String(), priv)
	}
	var _ encoding.TextMarshaler = sk
	var _ encoding.TextUnmarshaler = (*PrivateKey)(nil)
	var pks PrivateKey
	if err := pks.UnmarshalText([]byte(priv)); err != nil || pks.Type != sk.Type || !bytes.Equal(pks.Data, sk.Data) {
		t.Fatalf("PrivateKey.UnmarshalText err=%v", err)
	}
	if addr := sk.Address(); !addr.IsValid() {
		t.Fatalf("PrivateKey.Address invalid")
	}

	// GenerateKey (off-chain)
	gen, err := GenerateKey(KeyTypeEd25519)
	if err != nil || !gen.IsValid() || !gen.Public().IsValid() {
		t.Fatalf("GenerateKey err=%v valid=%v pub=%v", err, gen.IsValid(), gen.Public().IsValid())
	}
}

func TestKeyTypeHelpers(t *testing.T) {
	types := []KeyType{KeyTypeEd25519, KeyTypeSecp256k1, KeyTypeP256, KeyTypeBls12_381}
	for _, kt := range types {
		if !kt.IsValid() {
			t.Fatalf("KeyType %v expected valid", kt)
		}
		if kt.String() == "" {
			t.Fatalf("KeyType %v String empty", kt)
		}
		if len(kt.PkPrefixBytes()) == 0 || kt.PkPrefix() == "" {
			t.Fatalf("KeyType %v pk prefix empty", kt)
		}
		if len(kt.SkPrefixBytes()) == 0 || kt.SkPrefix() == "" {
			t.Fatalf("KeyType %v sk prefix empty", kt)
		}
		if kt.Tag() > 3 {
			t.Fatalf("KeyType %v unexpected tag %d", kt, kt.Tag())
		}
		if ParseKeyTag(kt.Tag()) != kt {
			t.Fatalf("ParseKeyTag(%d) mismatch", kt.Tag())
		}
	}
	if ParseKeyTag(255) != KeyTypeInvalid {
		t.Fatalf("ParseKeyTag(255) expected invalid")
	}
	if KeyTypeInvalid.IsValid() {
		t.Fatalf("KeyTypeInvalid expected invalid")
	}
	if KeyTypeInvalid.PkPrefix() != "" || KeyTypeInvalid.SkPrefix() != "" || KeyTypeInvalid.Tag() != 255 {
		t.Fatalf("KeyTypeInvalid helpers mismatch")
	}
}
