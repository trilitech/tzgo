// Copyright (c) 2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package tezos

import (
	"bytes"
	"encoding"
	"encoding/binary"

	// "encoding/hex"
	"fmt"
	"testing"
)

type Marshallable interface {
	encoding.TextUnmarshaler
	encoding.BinaryUnmarshaler
	fmt.Stringer
}

func TestHash(t *testing.T) {
	type testcase struct {
		String string
		Bytes  []byte
		Type   HashType
		Val    Marshallable
	}

	cases := []testcase{
		// chain id
		{
			String: "NetXdQprcVkpaWU",
			Bytes:  MustDecodeString("7a06a770"),
			Type:   HashTypeChainId,
			Val:    &ChainIdHash{},
		},
		// block
		{
			String: "BKjS7rtCjysnMNWUuevZiF2a6NkUas9bnSsNQ3ibh5GfKNrQoGk",
			Bytes:  MustDecodeString("029d4ed3161d644bedccb8673f30c6682b6e0a11756a3f75d7a739dede1cf29e"),
			Type:   HashTypeBlock,
			Val:    &BlockHash{},
		},
		// protocol
		{
			String: "PtLimaPtLMwfNinJi9rCfDPWea8dFgTZ1MeJ9f1m2SRic6ayiwW",
			Bytes:  MustDecodeString("d57ed88be5a69815e39386a33f7dcad391f5f507e03b376e499272c86c6cf2a7"),
			Type:   HashTypeProtocol,
			Val:    &ProtocolHash{},
		},
		// op
		{
			String: "oogC8ju9tMDqeB6RiAXdch3hnt8u3Pbf2ZXyyhAmJAhjQ4q1wUS",
			Bytes:  MustDecodeString("88315c911f6b4c38b2d8ea27319cf91d3614e1dde486dc83d55cd47bfbc568b4"),
			Type:   HashTypeOperation,
			Val:    &OpHash{},
		},
		// op list list
		{
			String: "LLoZnUuxzhNESHg7HvXxoccUCvWPrVmAucjaJDfKBeea39LqyVKEP",
			Bytes:  MustDecodeString("3ccb42fba1b24ce6c4e99ed1c4674dc90542d9fbc5694b9220ed74feb0eb3507"),
			Type:   HashTypeOperationListList,
			Val:    &OpListListHash{},
		},
		// payload
		{
			String: "vh26kqZ6LeKwygSY9JYZbqSdhtRjphAcrSEJqJ6a8EgnQEMtQx8J",
			Bytes:  MustDecodeString("37eec128736b994b3ce44a36c81dd00b2eab68057c11900a8135eaa5bff606fa"),
			Type:   HashTypeBlockPayload,
			Val:    &PayloadHash{},
		},
		// expr
		{
			String: "expruPBWMccKybChcJmGF8oMo263Ri6HgbKbAJRS8j6GbmqZJPJfVG",
			Bytes:  MustDecodeString("72386e5b4dbe9cc415bb0d909fe63e3162bc4b662ea4bfc92c485ad12f6700e8"),
			Type:   HashTypeScriptExpr,
			Val:    &ExprHash{},
		},
		// nonce
		{
			String: "nceUcirB7QYmVgcNUYQvd1fTzqSHjyusoE8VmJX9SNgELeQ4ffdcr",
			Bytes:  MustDecodeString("33b55c290efcc31f235a3809198ff324b31c65088a18815b8c67bf5ed1567dcd"),
			Type:   HashTypeNonce,
			Val:    &NonceHash{},
		},
		// context
		{
			String: "CoWAJ6dKTDySvhrV5njZwHpckRSBgzb84vXZKrEd1AwyshxFb9vo",
			Bytes:  MustDecodeString("c7c8fad1f5d2c144edbc763dc4e009c1fca3637dd265ae2c0044168065986b2b"),
			Type:   HashTypeContext,
			Val:    &ContextHash{},
		},
	}

	for i, c := range cases {
		// base58 must parse
		if err := c.Val.UnmarshalText([]byte(c.String)); err != nil {
			t.Fatalf("Case %d - unmarshal %s hash %s: %v", i, c.Type, c.String, err)
		}

		// write binary
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, c.Val)
		if err != nil {
			t.Fatalf("Case %d - write binary for hash %s: %v", i, c.String, err)
		}

		// check binary
		if !bytes.Equal(buf.Bytes(), c.Bytes) {
			t.Errorf("Case %d - mismatched hash got=%x want=%x", i, buf.Bytes(), c.Bytes)
		}

		// unmarshal from bytes
		if err := c.Val.UnmarshalBinary(c.Bytes); err != nil {
			t.Fatalf("Case %d - unmarshal binary %s: %v", i, c.Bytes, err)
		}

		// marshal text
		s := c.Val.String()

		// check text
		if s != c.String {
			t.Errorf("Case %d - mismatched text encoding got=%x want=%x", i, s, c.String)
		}

		// extra off-chain helpers: Bytes/MarshalText/MarshalBinary/Clone/Set
		switch v := c.Val.(type) {
		case *ChainIdHash:
			if !bytes.Equal(v.Bytes(), c.Bytes) {
				t.Errorf("Case %d - Bytes mismatch got=%x want=%x", i, v.Bytes(), c.Bytes)
			}
			if b, err := v.MarshalText(); err != nil || string(b) != c.String {
				t.Errorf("Case %d - MarshalText got=%q err=%v", i, string(b), err)
			}
			if b, err := v.MarshalBinary(); err != nil || !bytes.Equal(b, c.Bytes) {
				t.Errorf("Case %d - MarshalBinary got=%x err=%v", i, b, err)
			}
			if got := v.Clone(); got != *v {
				t.Errorf("Case %d - Clone mismatch", i)
			}
			var vv ChainIdHash
			if err := vv.Set(c.String); err != nil || vv != *v {
				t.Errorf("Case %d - Set mismatch err=%v", i, err)
			}
		case *BlockHash:
			if !bytes.Equal(v.Bytes(), c.Bytes) {
				t.Errorf("Case %d - Bytes mismatch got=%x want=%x", i, v.Bytes(), c.Bytes)
			}
			if b, err := v.MarshalText(); err != nil || string(b) != c.String {
				t.Errorf("Case %d - MarshalText got=%q err=%v", i, string(b), err)
			}
			if b, err := v.MarshalBinary(); err != nil || !bytes.Equal(b, c.Bytes) {
				t.Errorf("Case %d - MarshalBinary got=%x err=%v", i, b, err)
			}
			if got := v.Clone(); got != *v {
				t.Errorf("Case %d - Clone mismatch", i)
			}
			var vv BlockHash
			if err := vv.Set(c.String); err != nil || vv != *v {
				t.Errorf("Case %d - Set mismatch err=%v", i, err)
			}
		case *ProtocolHash:
			if !bytes.Equal(v.Bytes(), c.Bytes) {
				t.Errorf("Case %d - Bytes mismatch got=%x want=%x", i, v.Bytes(), c.Bytes)
			}
			if b, err := v.MarshalText(); err != nil || string(b) != c.String {
				t.Errorf("Case %d - MarshalText got=%q err=%v", i, string(b), err)
			}
			if b, err := v.MarshalBinary(); err != nil || !bytes.Equal(b, c.Bytes) {
				t.Errorf("Case %d - MarshalBinary got=%x err=%v", i, b, err)
			}
			if got := v.Clone(); got != *v {
				t.Errorf("Case %d - Clone mismatch", i)
			}
			var vv ProtocolHash
			if err := vv.Set(c.String); err != nil || vv != *v {
				t.Errorf("Case %d - Set mismatch err=%v", i, err)
			}
		case *OpHash:
			if !bytes.Equal(v.Bytes(), c.Bytes) {
				t.Errorf("Case %d - Bytes mismatch got=%x want=%x", i, v.Bytes(), c.Bytes)
			}
			if b, err := v.MarshalText(); err != nil || string(b) != c.String {
				t.Errorf("Case %d - MarshalText got=%q err=%v", i, string(b), err)
			}
			if b, err := v.MarshalBinary(); err != nil || !bytes.Equal(b, c.Bytes) {
				t.Errorf("Case %d - MarshalBinary got=%x err=%v", i, b, err)
			}
			if got := v.Clone(); got != *v {
				t.Errorf("Case %d - Clone mismatch", i)
			}
			var vv OpHash
			if err := vv.Set(c.String); err != nil || vv != *v {
				t.Errorf("Case %d - Set mismatch err=%v", i, err)
			}
		case *OpListListHash:
			if !bytes.Equal(v.Bytes(), c.Bytes) {
				t.Errorf("Case %d - Bytes mismatch got=%x want=%x", i, v.Bytes(), c.Bytes)
			}
			if b, err := v.MarshalText(); err != nil || string(b) != c.String {
				t.Errorf("Case %d - MarshalText got=%q err=%v", i, string(b), err)
			}
			if b, err := v.MarshalBinary(); err != nil || !bytes.Equal(b, c.Bytes) {
				t.Errorf("Case %d - MarshalBinary got=%x err=%v", i, b, err)
			}
			if got := v.Clone(); got != *v {
				t.Errorf("Case %d - Clone mismatch", i)
			}
			var vv OpListListHash
			if err := vv.Set(c.String); err != nil || vv != *v {
				t.Errorf("Case %d - Set mismatch err=%v", i, err)
			}
		case *PayloadHash:
			if !bytes.Equal(v.Bytes(), c.Bytes) {
				t.Errorf("Case %d - Bytes mismatch got=%x want=%x", i, v.Bytes(), c.Bytes)
			}
			if b, err := v.MarshalText(); err != nil || string(b) != c.String {
				t.Errorf("Case %d - MarshalText got=%q err=%v", i, string(b), err)
			}
			if b, err := v.MarshalBinary(); err != nil || !bytes.Equal(b, c.Bytes) {
				t.Errorf("Case %d - MarshalBinary got=%x err=%v", i, b, err)
			}
			if got := v.Clone(); got != *v {
				t.Errorf("Case %d - Clone mismatch", i)
			}
			var vv PayloadHash
			if err := vv.Set(c.String); err != nil || vv != *v {
				t.Errorf("Case %d - Set mismatch err=%v", i, err)
			}
		case *ExprHash:
			if !bytes.Equal(v.Bytes(), c.Bytes) {
				t.Errorf("Case %d - Bytes mismatch got=%x want=%x", i, v.Bytes(), c.Bytes)
			}
			if b, err := v.MarshalText(); err != nil || string(b) != c.String {
				t.Errorf("Case %d - MarshalText got=%q err=%v", i, string(b), err)
			}
			if b, err := v.MarshalBinary(); err != nil || !bytes.Equal(b, c.Bytes) {
				t.Errorf("Case %d - MarshalBinary got=%x err=%v", i, b, err)
			}
			if got := v.Clone(); got != *v {
				t.Errorf("Case %d - Clone mismatch", i)
			}
			var vv ExprHash
			if err := vv.Set(c.String); err != nil || vv != *v {
				t.Errorf("Case %d - Set mismatch err=%v", i, err)
			}
		case *NonceHash:
			if !bytes.Equal(v.Bytes(), c.Bytes) {
				t.Errorf("Case %d - Bytes mismatch got=%x want=%x", i, v.Bytes(), c.Bytes)
			}
			if b, err := v.MarshalText(); err != nil || string(b) != c.String {
				t.Errorf("Case %d - MarshalText got=%q err=%v", i, string(b), err)
			}
			if b, err := v.MarshalBinary(); err != nil || !bytes.Equal(b, c.Bytes) {
				t.Errorf("Case %d - MarshalBinary got=%x err=%v", i, b, err)
			}
			if got := v.Clone(); got != *v {
				t.Errorf("Case %d - Clone mismatch", i)
			}
			var vv NonceHash
			if err := vv.Set(c.String); err != nil || vv != *v {
				t.Errorf("Case %d - Set mismatch err=%v", i, err)
			}
		}
	}
}

func TestInvalidHash(t *testing.T) {
	// invalid base58 string
	if _, err := ParseBlockHash("tz1KzpjBnunNJVABHBnzfG4iuLmphitExW2"); err == nil {
		t.Errorf("Expected error on invalid base58 string")
	}

	// decode from short buffer
	var b BlockHash
	err := b.UnmarshalBinary(MustDecodeString("000b78887fdd0cd3bfbe75a717655728e0205bb9"))
	if err == nil {
		t.Errorf("Expected unmarshal error from short buffer")
	}

	// decode from nil buffer is OK
	err = b.UnmarshalBinary(nil)
	if err != nil {
		t.Errorf("Expected no unmarshal error from nil buffer, got %v", err)
	}

	// decode from empty string is OK
	err = b.UnmarshalText(nil)
	if err != nil {
		t.Errorf("Expected no unmarshal error from null string, got %v", err)
	}

	// ParseNonceHashSafe ignores errors and returns zero value on invalid input
	if got := ParseNonceHashSafe("not-a-nonce"); got.IsValid() {
		t.Errorf("ParseNonceHashSafe(invalid) expected invalid")
	}
}

func BenchmarkHashDecode(b *testing.B) {
	b.SetBytes(32)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseBlockHash("BKjS7rtCjysnMNWUuevZiF2a6NkUas9bnSsNQ3ibh5GfKNrQoGk")
	}
}

func BenchmarkHashEncode(b *testing.B) {
	b.SetBytes(32)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ZeroBlockHash.String()
	}
}

func TestHashTypeHelpers(t *testing.T) {
	if HashTypeInvalid.IsValid() {
		t.Fatalf("HashTypeInvalid.IsValid=true want false")
	}
	if !HashTypeBlock.IsValid() {
		t.Fatalf("HashTypeBlock.IsValid=false want true")
	}
	if HashTypeBlock.String() != HashTypeBlock.B58Prefix {
		t.Fatalf("HashType.String mismatch")
	}
	if !HashTypeBlock.Equal(HashType{B58Prefix: HashTypeBlock.B58Prefix}) {
		t.Fatalf("HashType.Equal mismatch")
	}
}

func TestHashHelpers_ExtraMethods(t *testing.T) {
	chain := MustParseChainIdHash("NetXdQprcVkpaWU")
	if !chain.IsValid() {
		t.Fatalf("chain id expected valid")
	}
	if got := chain.Clone(); got != chain {
		t.Fatalf("Clone mismatch")
	}
	if !bytes.Equal(chain.Bytes(), MustDecodeString("7a06a770")) {
		t.Fatalf("Bytes mismatch")
	}
	if b, err := chain.MarshalText(); err != nil || string(b) != chain.String() {
		t.Fatalf("MarshalText got=%q err=%v", string(b), err)
	}
	if b, err := chain.MarshalBinary(); err != nil || !bytes.Equal(b, chain[:]) {
		t.Fatalf("MarshalBinary mismatch err=%v", err)
	}
	if chain.Uint32() != binary.BigEndian.Uint32(chain[:]) {
		t.Fatalf("Uint32 mismatch")
	}
	var chain2 ChainIdHash
	if err := chain2.Set(chain.String()); err != nil || chain2 != chain {
		t.Fatalf("Set mismatch got=%v err=%v", chain2, err)
	}

	// Block hash extra helpers
	bh := MustParseBlockHash("BKjS7rtCjysnMNWUuevZiF2a6NkUas9bnSsNQ3ibh5GfKNrQoGk")
	if !bh.IsValid() || bh.Int64() != -1 {
		t.Fatalf("block hash helpers mismatch")
	}
	if got := bh.Clone(); got != bh {
		t.Fatalf("block Clone mismatch")
	}
	if b, err := bh.MarshalText(); err != nil || string(b) != bh.String() {
		t.Fatalf("block MarshalText err=%v", err)
	}
	if b, err := bh.MarshalBinary(); err != nil || !bytes.Equal(b, bh[:]) {
		t.Fatalf("block MarshalBinary err=%v", err)
	}
	var bh2 BlockHash
	if err := bh2.Set(bh.String()); err != nil || bh2 != bh {
		t.Fatalf("block Set mismatch err=%v", err)
	}

	// Protocol/Op/OpListList/Payload/Expr/Nonce parse helpers (previously untested Parse* and MustParse*)
	ph := MustParseProtocolHash("PtLimaPtLMwfNinJi9rCfDPWea8dFgTZ1MeJ9f1m2SRic6ayiwW")
	if got := ph.Clone(); got != ph {
		t.Fatalf("protocol Clone mismatch")
	}
	if b, err := ph.MarshalText(); err != nil || string(b) != ph.String() {
		t.Fatalf("protocol MarshalText err=%v", err)
	}
	if b, err := ph.MarshalBinary(); err != nil || !bytes.Equal(b, ph[:]) {
		t.Fatalf("protocol MarshalBinary err=%v", err)
	}
	var ph2 ProtocolHash
	if err := ph2.Set(ph.String()); err != nil || ph2 != ph {
		t.Fatalf("protocol Set mismatch err=%v", err)
	}

	if _, err := ParseOpHash("oogC8ju9tMDqeB6RiAXdch3hnt8u3Pbf2ZXyyhAmJAhjQ4q1wUS"); err != nil {
		t.Fatalf("ParseOpHash err=%v", err)
	}
	oh := MustParseOpHash("oogC8ju9tMDqeB6RiAXdch3hnt8u3Pbf2ZXyyhAmJAhjQ4q1wUS")
	if got := oh.Clone(); got != oh || !oh.IsValid() {
		t.Fatalf("op hash helpers mismatch")
	}
	if b, err := oh.MarshalText(); err != nil || string(b) != oh.String() {
		t.Fatalf("op MarshalText err=%v", err)
	}
	if b, err := oh.MarshalBinary(); err != nil || !bytes.Equal(b, oh[:]) {
		t.Fatalf("op MarshalBinary err=%v", err)
	}
	var oh2 OpHash
	if err := oh2.Set(oh.String()); err != nil || oh2 != oh {
		t.Fatalf("op Set mismatch err=%v", err)
	}

	if _, err := ParseOpListListHash("LLoZnUuxzhNESHg7HvXxoccUCvWPrVmAucjaJDfKBeea39LqyVKEP"); err != nil {
		t.Fatalf("ParseOpListListHash err=%v", err)
	}
	ll := MustParseOpListListHash("LLoZnUuxzhNESHg7HvXxoccUCvWPrVmAucjaJDfKBeea39LqyVKEP")
	if got := ll.Clone(); got != ll || !ll.IsValid() {
		t.Fatalf("oplistlist helpers mismatch")
	}
	var ll2 OpListListHash
	if err := ll2.Set(ll.String()); err != nil || ll2 != ll {
		t.Fatalf("oplistlist Set mismatch err=%v", err)
	}

	if _, err := ParsePayloadHash("vh26kqZ6LeKwygSY9JYZbqSdhtRjphAcrSEJqJ6a8EgnQEMtQx8J"); err != nil {
		t.Fatalf("ParsePayloadHash err=%v", err)
	}
	pl := MustParsePayloadHash("vh26kqZ6LeKwygSY9JYZbqSdhtRjphAcrSEJqJ6a8EgnQEMtQx8J")
	if got := pl.Clone(); got != pl || !pl.IsValid() {
		t.Fatalf("payload helpers mismatch")
	}
	var pl2 PayloadHash
	if err := pl2.Set(pl.String()); err != nil || pl2 != pl {
		t.Fatalf("payload Set mismatch err=%v", err)
	}

	if _, err := ParseExprHash("expruPBWMccKybChcJmGF8oMo263Ri6HgbKbAJRS8j6GbmqZJPJfVG"); err != nil {
		t.Fatalf("ParseExprHash err=%v", err)
	}
	ex := MustParseExprHash("expruPBWMccKybChcJmGF8oMo263Ri6HgbKbAJRS8j6GbmqZJPJfVG")
	if got := ex.Clone(); got != ex || !ex.IsValid() {
		t.Fatalf("expr helpers mismatch")
	}
	var ex2 ExprHash
	if err := ex2.Set(ex.String()); err != nil || ex2 != ex {
		t.Fatalf("expr Set mismatch err=%v", err)
	}

	if _, err := ParseNonceHash("nceUcirB7QYmVgcNUYQvd1fTzqSHjyusoE8VmJX9SNgELeQ4ffdcr"); err != nil {
		t.Fatalf("ParseNonceHash err=%v", err)
	}
	nn := MustParseNonceHash("nceUcirB7QYmVgcNUYQvd1fTzqSHjyusoE8VmJX9SNgELeQ4ffdcr")
	if got := nn.Clone(); got != nn || !nn.IsValid() {
		t.Fatalf("nonce helpers mismatch")
	}
	var nn2 NonceHash
	if err := nn2.Set(nn.String()); err != nil || nn2 != nn {
		t.Fatalf("nonce Set mismatch err=%v", err)
	}

	// ParseNonceHashSafe ignores errors and returns the zero value on invalid inputs.
	if got := ParseNonceHashSafe("nceUcirB7QYmVgcNUYQvd1fTzqSHjyusoE8VmJX9SNgELeQ4ffdcr"); !got.IsValid() {
		t.Fatalf("ParseNonceHashSafe(valid) returned invalid")
	}
	if got := ParseNonceHashSafe("not-a-nonce"); got.IsValid() {
		t.Fatalf("ParseNonceHashSafe(invalid) expected invalid")
	}

	// ContextHash + SmartRollup* hashes (construct from bytes, then roundtrip through parse/marshal).
	ctx := NewContextHash(bytes.Repeat([]byte{0x11}, 32))
	if !ctx.IsValid() {
		t.Fatalf("context expected valid")
	}
	ctxs := ctx.String()
	ctxp, err := ParseContextHash(ctxs)
	if err != nil || ctxp != ctx {
		t.Fatalf("ParseContextHash mismatch err=%v", err)
	}
	_ = MustParseContextHash(ctxs)
	var ctx2 ContextHash
	if err := ctx2.Set(ctxs); err != nil || ctx2 != ctx {
		t.Fatalf("context Set mismatch err=%v", err)
	}

	com := NewSmartRollupCommitHash(bytes.Repeat([]byte{0x22}, 32))
	commits := com.String()
	comp, err := ParseSmartRollupCommitHash(commits)
	if err != nil || comp != com {
		t.Fatalf("ParseSmartRollupCommitHash mismatch err=%v", err)
	}
	_ = MustParseSmartRollupCommitHash(commits)
	var com2 SmartRollupCommitHash
	if err := com2.Set(commits); err != nil || com2 != com {
		t.Fatalf("commit Set mismatch err=%v", err)
	}

	st := NewSmartRollupStateHash(bytes.Repeat([]byte{0x33}, 32))
	sts := st.String()
	stp, err := ParseSmartRollupStateHash(sts)
	if err != nil || stp != st {
		t.Fatalf("ParseSmartRollupStateHash mismatch err=%v", err)
	}
	_ = MustParseSmartRollupStateHash(sts)
	var st2 SmartRollupStateHash
	if err := st2.Set(sts); err != nil || st2 != st {
		t.Fatalf("state Set mismatch err=%v", err)
	}
}
