package rpc

import (
	"testing"

	"github.com/trilitech/tzgo/tezos"
)

func TestMapToChainParams_AllBakersAttestActivationThreshold(t *testing.T) {
	c := Constants{
		BlocksPerCycle:                     100,
		AllBakersAttestActivationThreshold: tezos.Ratio{Num: 1, Den: 2},
	}
	p := c.MapToChainParams()
	if p.AllBakersAttestActivationThreshold.Num != 1 || p.AllBakersAttestActivationThreshold.Den != 2 {
		t.Fatalf("expected ratio 1/2, got %d/%d", p.AllBakersAttestActivationThreshold.Num, p.AllBakersAttestActivationThreshold.Den)
	}
	if got := p.ActivationThresholdFloat(); got != 0.5 {
		t.Fatalf("expected float 0.5, got %v", got)
	}
}
