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
