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
