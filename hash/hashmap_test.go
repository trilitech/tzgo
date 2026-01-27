package hash

import "testing"

func TestHashMap_AddContainsRemove(t *testing.T) {
	m := NewHashMap()

	// add first element
	a := []byte("a")
	if got := m.Add(a); got != 1 {
		t.Fatalf("Add(a) len=%d, want 1", got)
	}
	if !m.Contains(a) {
		t.Fatalf("Contains(a)=false, want true")
	}

	// add second element
	b := []byte("b")
	if got := m.Add(b); got != 2 {
		t.Fatalf("Add(b) len=%d, want 2", got)
	}
	if !m.Contains(b) {
		t.Fatalf("Contains(b)=false, want true")
	}

	// overwrite existing (same key+value) should not change size
	if got := m.Add(a); got != 2 {
		t.Fatalf("Add(a) overwrite len=%d, want 2", got)
	}

	// remove existing should delete
	if got := m.Remove(a); got != 1 {
		t.Fatalf("Remove(a) len=%d, want 1", got)
	}
	if m.Contains(a) {
		t.Fatalf("Contains(a)=true after Remove, want false")
	}

	// remove missing should be no-op
	if got := m.Remove(a); got != 1 {
		t.Fatalf("Remove(a) missing len=%d, want 1", got)
	}
}

func TestHashMap_RemoveDoesNotDeleteOnValueMismatch(t *testing.T) {
	m := NewHashMap()

	buf := []byte("payload")
	h := Hash64(buf)
	// Force a mismatched value for an existing hash key.
	m[h] = []byte("different")

	if got := m.Remove(buf); got != 1 {
		t.Fatalf("Remove(buf) len=%d, want 1", got)
	}
	if _, ok := m[h]; !ok {
		t.Fatalf("expected key to remain when value mismatches")
	}
}

func TestHashMap_ContainsValueMismatch(t *testing.T) {
	m := NewHashMap()

	buf := []byte("payload")
	h := Hash64(buf)
	m[h] = []byte("different")

	if m.Contains(buf) {
		t.Fatalf("Contains(buf)=true with mismatched value, want false")
	}
	if m.Contains([]byte("absent")) {
		t.Fatalf("Contains(absent)=true, want false")
	}
}
