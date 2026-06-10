package tezos

import "testing"

func TestProtocolHistory_CloneAndAdd(t *testing.T) {
	h := ProtocolHistory{
		{Protocol: ProtoGenesis, StartHeight: 0, StartCycle: 0},
		{Protocol: ProtoBootstrap, StartHeight: 1, StartCycle: 0},
	}
	c := h.Clone()
	if len(c) != len(h) {
		t.Fatalf("Clone len=%d want %d", len(c), len(h))
	}
	// ensure it is a copy (slice backing array not shared)
	c[0].StartHeight = 999
	if h[0].StartHeight == 999 {
		t.Fatalf("Clone shares backing array")
	}

	// Add appends
	h.Add(Deployment{Protocol: PtLimaPt, StartHeight: 100, StartCycle: 10})
	if got := h.Last().Protocol; got != PtLimaPt {
		t.Fatalf("Last.Protocol=%s want %s", got, PtLimaPt)
	}
}

func TestProtocolHistory_Lookups(t *testing.T) {
	h := ProtocolHistory{
		{Protocol: ProtoGenesis, StartHeight: 0, EndHeight: 0, StartCycle: 0},
		{Protocol: ProtoBootstrap, StartHeight: 1, EndHeight: 9, StartCycle: 0},
		{Protocol: PtLimaPt, StartHeight: 10, EndHeight: -1, StartCycle: 1},
	}
	if got := h.AtBlock(0).Protocol; got != ProtoGenesis {
		t.Fatalf("AtBlock(0)=%s", got)
	}
	if got := h.AtBlock(5).Protocol; got != ProtoBootstrap {
		t.Fatalf("AtBlock(5)=%s", got)
	}
	if got := h.AtBlock(10).Protocol; got != PtLimaPt {
		t.Fatalf("AtBlock(10)=%s", got)
	}
	if got := h.AtCycle(0).Protocol; got != ProtoBootstrap {
		t.Fatalf("AtCycle(0)=%s", got)
	}
	if got := h.AtCycle(1).Protocol; got != PtLimaPt {
		t.Fatalf("AtCycle(1)=%s", got)
	}
	if got := h.AtProtocol(ProtoGenesis).Protocol; got != ProtoGenesis {
		t.Fatalf("AtProtocol(genesis)=%s", got)
	}
}

func TestParams_ActivationThresholdFloat(t *testing.T) {
	var p Params
	p.AllBakersAttestActivationThreshold = Ratio{Num: 1, Den: 0}
	if p.ActivationThresholdFloat() != 0 {
		t.Fatalf("den=0 expected 0")
	}
	p.AllBakersAttestActivationThreshold = Ratio{Num: 1, Den: 2}
	if p.ActivationThresholdFloat() != 0.5 {
		t.Fatalf("1/2 expected 0.5")
	}
}

func TestParams_WithNetworkAndWithBlock(t *testing.T) {
	p := NewParams().WithChainId(Mainnet).WithNetwork("x")
	// WithChainId sets default network when unknown; WithNetwork should not override it
	if p.Network != "Mainnet" {
		t.Fatalf("WithChainId network=%q want Mainnet", p.Network)
	}
	p.WithNetwork("y")
	if p.Network != "Mainnet" {
		t.Fatalf("WithNetwork should not override: %q", p.Network)
	}

	// WithBlock should set start/end/cycle based on deployments for chain id
	p2 := NewParams().WithChainId(Mainnet).WithBlock(1)
	if p2.StartHeight == 0 || p2.EndHeight == 0 {
		t.Fatalf("WithBlock did not set heights: start=%d end=%d", p2.StartHeight, p2.EndHeight)
	}
}
