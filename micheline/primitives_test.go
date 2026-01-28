// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

package micheline

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalIndexAddress(t *testing.T) {
	var prim Prim
	err := json.Unmarshal([]byte(`{"prim": "INDEX_ADDRESS"}`), &prim)
	assert.Nil(t, err)
	assert.Equal(t, I_INDEX_ADDRESS, prim.OpCode)

	err = json.Unmarshal([]byte(`{"prim": "Index_Address"}`), &prim)
	assert.NotNil(t, err)
}

func TestUnmarshalGetAddressIndex(t *testing.T) {
	var prim Prim
	err := json.Unmarshal([]byte(`{"prim": "GET_ADDRESS_INDEX"}`), &prim)
	assert.Nil(t, err)
	assert.Equal(t, I_GET_ADDRESS_INDEX, prim.OpCode)

	err = json.Unmarshal([]byte(`{"prim": "Get_address_index"}`), &prim)
	assert.NotNil(t, err)
}

func TestPrim_AnnoHelpers(t *testing.T) {
	p := Prim{Anno: []string{"%a", ":t", "@f"}}
	if !p.HasAnno() {
		t.Fatalf("HasAnno=false want true")
	}
	if !p.MatchesAnno("%a") || !p.MatchesAnno("a") {
		t.Fatalf("MatchesAnno failed for var anno")
	}
	if !p.HasTypeAnno() || p.GetTypeAnno() != "t" || p.GetTypeAnnoAny() != "t" {
		t.Fatalf("type anno helpers mismatch: has=%v get=%q any=%q", p.HasTypeAnno(), p.GetTypeAnno(), p.GetTypeAnnoAny())
	}
	if !p.HasVarAnno() || p.GetVarAnno() != "a" || p.GetVarAnnoAny() != "a" {
		t.Fatalf("var anno helpers mismatch: has=%v get=%q any=%q", p.HasVarAnno(), p.GetVarAnno(), p.GetVarAnnoAny())
	}
	if !p.HasFieldAnno() || p.GetFieldAnno() != "f" || p.GetFieldAnnoAny() != "f" {
		t.Fatalf("field anno helpers mismatch: has=%v get=%q any=%q", p.HasFieldAnno(), p.GetFieldAnno(), p.GetFieldAnnoAny())
	}
	if !p.HasVarOrFieldAnno() || p.GetVarOrFieldAnno() == "" {
		t.Fatalf("var/field anno helpers mismatch")
	}

	// StripAnno removes by name without prefix.
	p2 := Prim{Anno: []string{"%x", "%y"}}
	p2.StripAnno("x")
	if len(p2.Anno) != 1 || p2.Anno[0] != "%y" {
		t.Fatalf("StripAnno failed: %#v", p2.Anno)
	}
}
