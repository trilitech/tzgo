package rpc

import (
	"encoding/json"
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

// TestConstantsV025Dal verifies the v025 DAL parametric constants decode from
// the dal_parametric object and are surfaced into tezos.Params.
func TestConstantsV025Dal(t *testing.T) {
	const blob = `{
		"blocks_per_cycle": 10800,
		"cache_layout_size": 5,
		"dal_parametric": {
			"feature_enable": true,
			"number_of_slots": 160,
			"slot_size": 380832,
			"page_size": 3967,
			"number_of_shards": 512,
			"redundancy_factor": 8,
			"attestation_lag": 5,
			"attestation_threshold": 66,
			"attestation_lags": [1, 2, 3, 4, 5]
		}
	}`

	var c Constants
	if err := json.Unmarshal([]byte(blob), &c); err != nil {
		t.Fatalf("unmarshal constants: %v", err)
	}
	if c.CacheLayoutSize != 5 {
		t.Errorf("CacheLayoutSize = %d, want 5", c.CacheLayoutSize)
	}
	if c.Dal.NumberOfSlots != 160 || c.Dal.SlotSize != 380832 || c.Dal.AttestationLag != 5 {
		t.Errorf("unexpected DAL constants: %+v", c.Dal)
	}
	if len(c.Dal.AttestationLags) != 5 || c.Dal.AttestationLags[4] != 5 {
		t.Errorf("AttestationLags = %v, want [1 2 3 4 5]", c.Dal.AttestationLags)
	}

	p := c.MapToChainParams()
	if p.DalNumberOfSlots != 160 || p.DalSlotSize != 380832 || p.DalAttestationLag != 5 {
		t.Errorf("DAL params not surfaced: slots=%d size=%d lag=%d", p.DalNumberOfSlots, p.DalSlotSize, p.DalAttestationLag)
	}
	if len(p.DalAttestationLags) != 5 {
		t.Errorf("DalAttestationLags = %v, want length 5", p.DalAttestationLags)
	}
}

// TestConstantsPreV025NoDal verifies that constants from chains without DAL
// parameters (no dal_parametric key, no cache_layout_size) decode cleanly and
// leave the DAL fields in tezos.Params zero-valued.
func TestConstantsPreV025NoDal(t *testing.T) {
	const blob = `{"blocks_per_cycle": 8192, "consensus_rights_delay": 2}`

	var c Constants
	if err := json.Unmarshal([]byte(blob), &c); err != nil {
		t.Fatalf("unmarshal pre-v025 constants: %v", err)
	}
	if c.CacheLayoutSize != 0 {
		t.Errorf("CacheLayoutSize = %d, want 0 on pre-v025 chains", c.CacheLayoutSize)
	}
	if c.Dal.NumberOfSlots != 0 || c.Dal.AttestationLag != 0 || c.Dal.AttestationLags != nil {
		t.Errorf("unexpected non-zero DAL constants: %+v", c.Dal)
	}

	p := c.MapToChainParams()
	if p.DalNumberOfSlots != 0 || p.DalSlotSize != 0 || p.DalAttestationLag != 0 || p.DalAttestationLags != nil {
		t.Errorf("DAL params must stay zero without dal_parametric: %+v", p)
	}
}
