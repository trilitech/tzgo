package tezos

import (
	"bytes"
	"testing"
)

func TestLimitsAndCosts_Add(t *testing.T) {
	l1 := Limits{Fee: 1, GasLimit: 2, StorageLimit: 3}
	l2 := Limits{Fee: 4, GasLimit: 5, StorageLimit: 6}
	sum := l1.Add(l2)
	if sum.Fee != 5 || sum.GasLimit != 7 || sum.StorageLimit != 9 {
		t.Fatalf("Limits.Add got=%+v", sum)
	}
	// inputs unchanged
	if l1.Fee != 1 || l2.Fee != 4 {
		t.Fatalf("Limits.Add modified inputs")
	}

	c1 := Costs{Fee: 1, Burn: 2, GasUsed: 3, StorageUsed: 4, StorageBurn: 5, AllocationBurn: 6}
	c2 := Costs{Fee: 10, Burn: 20, GasUsed: 30, StorageUsed: 40, StorageBurn: 50, AllocationBurn: 60}
	csum := c1.Add(c2)
	if csum.Fee != 11 || csum.Burn != 22 || csum.GasUsed != 33 || csum.StorageUsed != 44 || csum.StorageBurn != 55 || csum.AllocationBurn != 66 {
		t.Fatalf("Costs.Add got=%+v", csum)
	}
	if c1.Fee != 1 || c2.Fee != 10 {
		t.Fatalf("Costs.Add modified inputs")
	}
}

func TestAddressFilter(t *testing.T) {
	a1 := MustParseAddress("tz1LggX2HUdvJ1tF4Fvv8fjsrzLeW4Jr9t2Q")
	a2 := MustParseAddress("KT1GyeRktoGPEKsWpchWguyy8FAf3aNHkw2T")

	f := NewAddressFilter(a1)
	if !f.Contains(a1) || f.Contains(a2) {
		t.Fatalf("Contains mismatch")
	}
	if f.Len() != 1 {
		t.Fatalf("Len=%d want 1", f.Len())
	}
	if ok := f.AddUnique(a1); ok {
		t.Fatalf("AddUnique(existing)=true want false")
	}
	if ok := f.AddUnique(a2); !ok {
		t.Fatalf("AddUnique(new)=false want true")
	}
	if f.Len() != 2 {
		t.Fatalf("Len=%d want 2", f.Len())
	}
	f.Remove(a1)
	if f.Contains(a1) {
		t.Fatalf("Remove failed")
	}

	// Merge
	f2 := NewAddressFilter(a1)
	f.Merge(f2)
	if !f.Contains(a1) || !f.Contains(a2) {
		t.Fatalf("Merge failed")
	}
}

func TestPrivateKey_Encrypt_Roundtrip(t *testing.T) {
	// Covers encryptPrivateKey via PrivateKey.Encrypt (off-chain).
	sk, err := ParsePrivateKey("edsk4FTF78Qf1m2rykGpHqostAiq5gYW4YZEoGUSWBTJr2njsDHSnd")
	if err != nil {
		t.Fatalf("ParsePrivateKey err=%v", err)
	}
	enc, err := sk.Encrypt(func() ([]byte, error) { return []byte("password"), nil })
	if err != nil {
		t.Fatalf("Encrypt err=%v", err)
	}
	if !IsEncryptedKey(enc) {
		t.Fatalf("Encrypt produced non-encrypted key: %s", enc)
	}
	dec, err := ParseEncryptedPrivateKey(enc, func() ([]byte, error) { return []byte("password"), nil })
	if err != nil {
		t.Fatalf("ParseEncryptedPrivateKey err=%v", err)
	}
	if dec.Type != sk.Type || !bytes.Equal(dec.Data, sk.Data) {
		t.Fatalf("Encrypt/Decrypt roundtrip mismatch")
	}
}
