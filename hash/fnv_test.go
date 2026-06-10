package hash

import (
	"encoding/binary"
	"encoding/hex"
	"hash/fnv"
	"testing"
	"testing/quick"
)

// test implementation of FNV Hash64
func TestNewHash64(t *testing.T) {
	expect := func(data []byte) uint64 {
		h := fnv.New64a()
		if _, err := h.Write(data); err != nil {
			t.Fatal(err)
		}
		return h.Sum64()
	}
	hash := func(data []byte) uint64 {
		h := NewInlineFNV64a()
		h.Write(data)
		return h.Sum64()
	}
	if err := quick.CheckEqual(hash, expect, nil); err != nil {
		t.Fatal(err)
	}
}

func TestStaticHash64(t *testing.T) {
	expect := func(data []byte) uint64 {
		h := fnv.New64a()
		if _, err := h.Write(data); err != nil {
			t.Fatal(err)
		}
		return h.Sum64()
	}
	if err := quick.CheckEqual(Hash64, expect, nil); err != nil {
		t.Fatal(err)
	}
}

func TestInlineFNV64a_WriteStringSumReset(t *testing.T) {
	h := NewInlineFNV64a()
	if got := h.Sum64(); got != offset64 {
		t.Fatalf("initial Sum64=%d, want %d", got, uint64(offset64))
	}

	// WriteString delegates to Write([]byte(...)).
	if n, err := h.WriteString("abc"); err != nil || n != 3 {
		t.Fatalf("WriteString n=%d err=%v, want n=3 err=nil", n, err)
	}
	want := Hash64([]byte("abc"))
	if got := h.Sum64(); got != want {
		t.Fatalf("Sum64=%d, want %d", got, want)
	}

	sum := h.Sum()
	if len(sum) != 8 {
		t.Fatalf("Sum len=%d, want 8", len(sum))
	}
	if got := binary.BigEndian.Uint64(sum); got != want {
		t.Fatalf("Sum bytes=%d, want %d", got, want)
	}

	h.Reset()
	if got := h.Sum64(); got != offset64 {
		t.Fatalf("after Reset Sum64=%d, want %d", got, uint64(offset64))
	}
}

func BenchmarkNewHash64(b *testing.B) {
	buf, _ := hex.DecodeString("029d4ed3161d644bedccb8673f30c6682b6e0a11756a3f75d7a739dede1cf29e")
	b.SetBytes(32)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h := NewInlineFNV64a()
		h.Write(buf)
		_ = h.Sum64()
	}
}

func BenchmarkStaticHash64(b *testing.B) {
	buf, _ := hex.DecodeString("029d4ed3161d644bedccb8673f30c6682b6e0a11756a3f75d7a739dede1cf29e")
	b.SetBytes(32)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Hash64(buf)
	}
}
