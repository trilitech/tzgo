package tezos_test

import (
	"testing"

	"github.com/trilitech/tzgo/tezos"
)

// TestProtoV025Registration verifies that protocol v025 (Ushuaia) is registered
// with the expected version number and that ProtoAlpha is bumped past it.
func TestProtoV025Registration(t *testing.T) {
	if got := tezos.Versions[tezos.ProtoV025]; got != 25 {
		t.Errorf("ProtoV025 version = %d, want 25", got)
	}
	if got := tezos.Versions[tezos.ProtoAlpha]; got != 26 {
		t.Errorf("ProtoAlpha version = %d, want 26 (must not collide with v025)", got)
	}
	if !tezos.PsUshuai.Equal(tezos.ProtoV025) {
		t.Errorf("PsUshuai alias does not equal ProtoV025")
	}
}

// TestWithProtocolV025 verifies WithProtocol resolves v025 to version 25 and
// operation tags version 4 (v23+).
func TestWithProtocolV025(t *testing.T) {
	p := tezos.NewParams().WithProtocol(tezos.ProtoV025)
	if p.Version != 25 {
		t.Errorf("v025 params Version = %d, want 25", p.Version)
	}
	if p.OperationTagsVersion != 4 {
		t.Errorf("v025 params OperationTagsVersion = %d, want 4", p.OperationTagsVersion)
	}
}
