// Copyright (c) 2020-2024 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package micheline

import (
	"encoding/json"
	"testing"
)

func TestSeoulPrimitivesUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		opcode  OpCode
	}{
		// Positive cases - valid Seoul primitives
		{
			name:    "IS_IMPLICIT_ACCOUNT primitive",
			json:    `{"prim": "IS_IMPLICIT_ACCOUNT"}`,
			wantErr: false,
			opcode:  I_IS_IMPLICIT_ACCOUNT,
		},
		{
			name:    "IS_IMPLICIT_ACCOUNT with args",
			json:    `{"prim": "IS_IMPLICIT_ACCOUNT", "args": [{"prim": "PUSH", "args": [{"prim": "address"}, {"string": "tz1abc"}]}]}`,
			wantErr: false,
			opcode:  I_IS_IMPLICIT_ACCOUNT,
		},
		{
			name:    "IS_IMPLICIT_ACCOUNT with annotations",
			json:    `{"prim": "IS_IMPLICIT_ACCOUNT", "annots": ["%check_implicit"]}`,
			wantErr: false,
			opcode:  I_IS_IMPLICIT_ACCOUNT,
		},
		// Existing primitives for comparison
		{
			name:    "IMPLICIT_ACCOUNT primitive",
			json:    `{"prim": "IMPLICIT_ACCOUNT"}`,
			wantErr: false,
			opcode:  I_IMPLICIT_ACCOUNT,
		},
		{
			name:    "PACK primitive",
			json:    `{"prim": "PACK"}`,
			wantErr: false,
			opcode:  I_PACK,
		},
		// Negative cases - invalid primitives
		{
			name:    "unknown primitive",
			json:    `{"prim": "UNKNOWN_PRIMITIVE"}`,
			wantErr: true,
			opcode:  OpCode(255),
		},
		{
			name:    "case sensitive primitive",
			json:    `{"prim": "is_implicit_account"}`,
			wantErr: true,
			opcode:  OpCode(255),
		},
		{
			name:    "typo in primitive name",
			json:    `{"prim": "IS_IMPLICIT_ACCONT"}`,
			wantErr: true,
			opcode:  OpCode(255),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var prim Prim
			err := json.Unmarshal([]byte(tt.json), &prim)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for JSON: %s", tt.json)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for JSON %s: %v", tt.json, err)
				}
				if prim.OpCode != tt.opcode {
					t.Errorf("Expected opcode %v, got %v for JSON: %s", tt.opcode, prim.OpCode, tt.json)
				}
			}
		})
	}
}

func TestSeoulPrimitivesInOrigination(t *testing.T) {
	// Test cases that simulate the original error scenario
	// These represent Michelson code that might appear in origination operations
	tests := []struct {
		name        string
		michelson   string
		expectError bool
		description string
	}{
		{
			name: "contract with IS_IMPLICIT_ACCOUNT",
			michelson: `[
				{"prim": "parameter", "args": [{"prim": "address"}]},
				{"prim": "storage", "args": [{"prim": "bool"}]},
				{"prim": "code", "args": [[
					{"prim": "CAR"},
					{"prim": "IS_IMPLICIT_ACCOUNT"},
					{"prim": "NIL", "args": [{"prim": "operation"}]},
					{"prim": "PAIR"}
				]]}
			]`,
			expectError: false,
			description: "Contract that uses IS_IMPLICIT_ACCOUNT to check if parameter is implicit account",
		},
		{
			name: "complex contract with IS_IMPLICIT_ACCOUNT",
			michelson: `[
				{"prim": "parameter", "args": [{"prim": "pair", "args": [{"prim": "address"}, {"prim": "mutez"}]}]},
				{"prim": "storage", "args": [{"prim": "big_map", "args": [{"prim": "address"}, {"prim": "mutez"}]}]},
				{"prim": "code", "args": [[
					{"prim": "DUP"},
					{"prim": "CAR"},
					{"prim": "CAR"},
					{"prim": "IS_IMPLICIT_ACCOUNT"},
					{"prim": "IF", "args": [
						[{"prim": "PUSH", "args": [{"prim": "string"}, {"string": "implicit account"}]}],
						[{"prim": "PUSH", "args": [{"prim": "string"}, {"string": "contract account"}]}]
					]},
					{"prim": "DROP"},
					{"prim": "NIL", "args": [{"prim": "operation"}]},
					{"prim": "PAIR"}
				]]}
			]`,
			expectError: false,
			description: "Complex contract with conditional logic based on IS_IMPLICIT_ACCOUNT",
		},
		{
			name: "contract with unknown primitive",
			michelson: `[
				{"prim": "parameter", "args": [{"prim": "address"}]},
				{"prim": "storage", "args": [{"prim": "bool"}]},
				{"prim": "code", "args": [[
					{"prim": "CAR"},
					{"prim": "UNKNOWN_SEOUL_PRIMITIVE"},
					{"prim": "NIL", "args": [{"prim": "operation"}]},
					{"prim": "PAIR"}
				]]}
			]`,
			expectError: true,
			description: "Contract with unknown primitive should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var script Script
			err := json.Unmarshal([]byte(tt.michelson), &script.Code)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.description)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.description, err)
				}

				// Verify the script is valid
				if !script.IsValid() {
					t.Errorf("Script should be valid for %s", tt.description)
				}
			}
		})
	}
}

func TestSeoulPrimitivesRoundTrip(t *testing.T) {
	// Test that Seoul primitives can be marshaled and unmarshaled correctly
	tests := []struct {
		name string
		prim Prim
	}{
		{
			name: "IS_IMPLICIT_ACCOUNT nullary",
			prim: Prim{
				Type:   PrimNullary,
				OpCode: I_IS_IMPLICIT_ACCOUNT,
			},
		},
		{
			name: "IS_IMPLICIT_ACCOUNT with annotation",
			prim: Prim{
				Type:   PrimNullaryAnno,
				OpCode: I_IS_IMPLICIT_ACCOUNT,
				Anno:   []string{"%is_implicit"},
			},
		},
		{
			name: "IS_IMPLICIT_ACCOUNT in sequence",
			prim: Prim{
				Type: PrimSequence,
				Args: []Prim{
					{Type: PrimNullary, OpCode: I_PUSH},
					{Type: PrimNullary, OpCode: I_IS_IMPLICIT_ACCOUNT},
					{Type: PrimNullary, OpCode: I_IF},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.prim)
			if err != nil {
				t.Errorf("Failed to marshal prim: %v", err)
				return
			}

			// Unmarshal back
			var unmarshaled Prim
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal prim: %v", err)
				return
			}

			// Compare
			if !tt.prim.IsEqual(unmarshaled) {
				t.Errorf("Round-trip failed:\nOriginal: %s\nUnmarshaled: %s",
					tt.prim.Dump(), unmarshaled.Dump())
			}
		})
	}
}

func TestSeoulPrimitivesErrorMessages(t *testing.T) {
	// Test that error messages are helpful for debugging
	tests := []struct {
		name        string
		json        string
		expectedErr string
	}{
		{
			name:        "typo in IS_IMPLICIT_ACCOUNT",
			json:        `{"prim": "IS_IMPLICIT_ACCONT"}`,
			expectedErr: "Unknown michelson primitive IS_IMPLICIT_ACCONT",
		},
		{
			name:        "lowercase IS_IMPLICIT_ACCOUNT",
			json:        `{"prim": "is_implicit_account"}`,
			expectedErr: "Unknown michelson primitive is_implicit_account",
		},
		{
			name:        "partial primitive name",
			json:        `{"prim": "IS_IMPLICIT"}`,
			expectedErr: "Unknown michelson primitive IS_IMPLICIT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var prim Prim
			err := json.Unmarshal([]byte(tt.json), &prim)

			if err == nil {
				t.Errorf("Expected error for JSON: %s", tt.json)
				return
			}

			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

// Test performance of Seoul primitives
func BenchmarkSeoulPrimitivesUnmarshal(b *testing.B) {
	jsonData := []byte(`{"prim": "IS_IMPLICIT_ACCOUNT"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var prim Prim
		_ = json.Unmarshal(jsonData, &prim)
	}
}

func BenchmarkSeoulPrimitivesMarshal(b *testing.B) {
	prim := Prim{
		Type:   PrimNullary,
		OpCode: I_IS_IMPLICIT_ACCOUNT,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(prim)
	}
}
