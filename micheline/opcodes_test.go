// Copyright (c) 2025 TriliTech Ltd.
// Author: tzstats@trili.tech

package micheline

import (
	"testing"
)

func TestParseOpCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected OpCode
		wantErr  bool
	}{
		// Positive cases - existing opcodes
		{
			name:     "IMPLICIT_ACCOUNT",
			input:    "IMPLICIT_ACCOUNT",
			expected: I_IMPLICIT_ACCOUNT,
			wantErr:  false,
		},
		{
			name:     "PACK",
			input:    "PACK",
			expected: I_PACK,
			wantErr:  false,
		},
		{
			name:     "ADD",
			input:    "ADD",
			expected: I_ADD,
			wantErr:  false,
		},
		{
			name:     "INDEX_ADDRESS",
			input:    "INDEX_ADDRESS",
			expected: I_INDEX_ADDRESS,
			wantErr:  false,
		},
		{
			name:     "GET_ADDRESS_INDEX",
			input:    "GET_ADDRESS_INDEX",
			expected: I_GET_ADDRESS_INDEX,
			wantErr:  false,
		},
		// Positive case - new Seoul primitive
		{
			name:     "IS_IMPLICIT_ACCOUNT",
			input:    "IS_IMPLICIT_ACCOUNT",
			expected: I_IS_IMPLICIT_ACCOUNT,
			wantErr:  false,
		},
		// Negative cases - invalid opcodes
		{
			name:     "empty string",
			input:    "",
			expected: 255,
			wantErr:  true,
		},
		{
			name:     "unknown opcode",
			input:    "UNKNOWN_OPCODE",
			expected: 255,
			wantErr:  true,
		},
		{
			name:     "case sensitive - lowercase",
			input:    "is_implicit_account",
			expected: 255,
			wantErr:  true,
		},
		{
			name:     "invalid characters",
			input:    "IS_IMPLICIT_ACCOUNT!",
			expected: 255,
			wantErr:  true,
		},
		{
			name:     "partial match",
			input:    "IS_IMPLICIT",
			expected: 255,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseOpCode(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseOpCode(%q) expected error, got nil", tt.input)
				}
				if result != tt.expected {
					t.Errorf("ParseOpCode(%q) = %v, expected %v on error", tt.input, result, tt.expected)
				}
			} else {
				if err != nil {
					t.Errorf("ParseOpCode(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseOpCode(%q) = %v, expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestOpCodeString(t *testing.T) {
	tests := []struct {
		name     string
		opcode   OpCode
		expected string
	}{
		// Positive cases
		{
			name:     "IS_IMPLICIT_ACCOUNT string representation",
			opcode:   I_IS_IMPLICIT_ACCOUNT,
			expected: "IS_IMPLICIT_ACCOUNT",
		},
		{
			name:     "IMPLICIT_ACCOUNT string representation",
			opcode:   I_IMPLICIT_ACCOUNT,
			expected: "IMPLICIT_ACCOUNT",
		},
		{
			name:     "PACK string representation",
			opcode:   I_PACK,
			expected: "PACK",
		},
		{
			name:     "INDEX_ADDRESS string representation",
			opcode:   I_INDEX_ADDRESS,
			expected: "INDEX_ADDRESS",
		},
		{
			name:     "GET_ADDRESS_INDEX string representation",
			opcode:   I_GET_ADDRESS_INDEX,
			expected: "GET_ADDRESS_INDEX",
		},
		// Negative case - invalid opcode
		{
			name:     "invalid opcode",
			opcode:   OpCode(255),
			expected: "Unknown michelson opcode 0xff",
		},
		{
			name:     "out of range opcode",
			opcode:   OpCode(200),
			expected: "Unknown michelson opcode 0xc8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.opcode.String()
			if result != tt.expected {
				t.Errorf("OpCode(%d).String() = %q, expected %q", tt.opcode, result, tt.expected)
			}
		})
	}
}

func TestOpCodeIsValid(t *testing.T) {
	tests := []struct {
		name     string
		opcode   OpCode
		expected bool
	}{
		// Positive cases - valid opcodes
		{
			name:     "IS_IMPLICIT_ACCOUNT is valid",
			opcode:   I_IS_IMPLICIT_ACCOUNT,
			expected: true,
		},
		{
			name:     "IMPLICIT_ACCOUNT is valid",
			opcode:   I_IMPLICIT_ACCOUNT,
			expected: true,
		},
		{
			name:     "first opcode K_PARAMETER is valid",
			opcode:   K_PARAMETER,
			expected: true,
		},
		{
			name:     "D_TICKET is valid",
			opcode:   D_TICKET,
			expected: true,
		},
		{
			name:     "I_NAT is valid",
			opcode:   I_NAT,
			expected: true,
		},
		{
			name:     "INDEX_ADDRESS is valid",
			opcode:   I_INDEX_ADDRESS,
			expected: true,
		},
		{
			name:     "GET_ADDRESS_INDEX is valid",
			opcode:   I_GET_ADDRESS_INDEX,
			expected: true,
		},
		// Negative cases - invalid opcodes
		{
			name:     "opcode beyond GET_ADDRESS_INDEX is invalid",
			opcode:   I_GET_ADDRESS_INDEX + 1,
			expected: false,
		},
		{
			name:     "high value opcode is invalid",
			opcode:   OpCode(200),
			expected: false,
		},
		{
			name:     "max byte value is invalid",
			opcode:   OpCode(255),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.opcode.IsValid()
			if result != tt.expected {
				t.Errorf("OpCode(%d).IsValid() = %v, expected %v", tt.opcode, result, tt.expected)
			}
		})
	}
}

func TestOpCodeByte(t *testing.T) {
	tests := []struct {
		name     string
		opcode   OpCode
		expected byte
	}{
		{
			name:     "IS_IMPLICIT_ACCOUNT byte value",
			opcode:   I_IS_IMPLICIT_ACCOUNT,
			expected: 0x9E, // 158 decimal
		},
		{
			name:     "D_TICKET byte value",
			opcode:   D_TICKET,
			expected: 0x9D, // 157 decimal
		},
		{
			name:     "I_NAT byte value",
			opcode:   I_NAT,
			expected: 0x9C, // 156 decimal
		},
		{
			name:     "K_PARAMETER byte value",
			opcode:   K_PARAMETER,
			expected: 0x00, // 0 decimal
		},
		{
			name:     "INDEX_ADDRESS byte value",
			opcode:   I_INDEX_ADDRESS,
			expected: 0x9F, // 159 decimal
		},
		{
			name:     "GET_ADDRESS_INDEX byte value",
			opcode:   I_GET_ADDRESS_INDEX,
			expected: 0xA0, // 160 decimal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.opcode.Byte()
			if result != tt.expected {
				t.Errorf("OpCode(%d).Byte() = 0x%02x, expected 0x%02x", tt.opcode, result, tt.expected)
			}
		})
	}
}

func TestOpCodeMarshalText(t *testing.T) {
	tests := []struct {
		name     string
		opcode   OpCode
		expected string
		wantErr  bool
	}{
		{
			name:     "IS_IMPLICIT_ACCOUNT marshal text",
			opcode:   I_IS_IMPLICIT_ACCOUNT,
			expected: "IS_IMPLICIT_ACCOUNT",
			wantErr:  false,
		},
		{
			name:     "IMPLICIT_ACCOUNT marshal text",
			opcode:   I_IMPLICIT_ACCOUNT,
			expected: "IMPLICIT_ACCOUNT",
			wantErr:  false,
		},
		{
			name:     "INDEX_ADDRESS marshal text",
			opcode:   I_INDEX_ADDRESS,
			expected: "INDEX_ADDRESS",
			wantErr:  false,
		},
		{
			name:     "GET_ADDRESS_INDEX marshal text",
			opcode:   I_GET_ADDRESS_INDEX,
			expected: "GET_ADDRESS_INDEX",
			wantErr:  false,
		},
		{
			name:     "invalid opcode marshal text",
			opcode:   OpCode(255),
			expected: "Unknown michelson opcode 0xff",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.opcode.MarshalText()

			if tt.wantErr && err == nil {
				t.Errorf("OpCode(%d).MarshalText() expected error, got nil", tt.opcode)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("OpCode(%d).MarshalText() unexpected error: %v", tt.opcode, err)
			}

			if string(result) != tt.expected {
				t.Errorf("OpCode(%d).MarshalText() = %q, expected %q", tt.opcode, string(result), tt.expected)
			}
		})
	}
}

func TestSeoulProtocolPrimitives(t *testing.T) {
	// Test that Seoul protocol primitives are properly defined
	t.Run("Seoul primitives exist", func(t *testing.T) {
		// Test IS_IMPLICIT_ACCOUNT exists and has correct properties
		if !I_IS_IMPLICIT_ACCOUNT.IsValid() {
			t.Error("I_IS_IMPLICIT_ACCOUNT should be valid")
		}

		if I_IS_IMPLICIT_ACCOUNT.String() != "IS_IMPLICIT_ACCOUNT" {
			t.Errorf("I_IS_IMPLICIT_ACCOUNT.String() = %q, expected %q",
				I_IS_IMPLICIT_ACCOUNT.String(), "IS_IMPLICIT_ACCOUNT")
		}

		if I_IS_IMPLICIT_ACCOUNT.Byte() != 0x9E {
			t.Errorf("I_IS_IMPLICIT_ACCOUNT.Byte() = 0x%02x, expected 0x9E",
				I_IS_IMPLICIT_ACCOUNT.Byte())
		}
	})

	t.Run("Seoul primitives round-trip", func(t *testing.T) {
		// Test that we can parse the string representation back to the opcode
		parsed, err := ParseOpCode("IS_IMPLICIT_ACCOUNT")
		if err != nil {
			t.Errorf("ParseOpCode(IS_IMPLICIT_ACCOUNT) error: %v", err)
		}
		if parsed != I_IS_IMPLICIT_ACCOUNT {
			t.Errorf("ParseOpCode(IS_IMPLICIT_ACCOUNT) = %v, expected %v",
				parsed, I_IS_IMPLICIT_ACCOUNT)
		}
	})
}

// Benchmark tests for performance validation
func BenchmarkParseOpCode(b *testing.B) {
	b.Run("IS_IMPLICIT_ACCOUNT", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = ParseOpCode("IS_IMPLICIT_ACCOUNT")
		}
	})

	b.Run("PACK", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = ParseOpCode("PACK")
		}
	})

	b.Run("invalid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = ParseOpCode("INVALID_OPCODE")
		}
	})
}

func BenchmarkOpCodeString(b *testing.B) {
	b.Run("IS_IMPLICIT_ACCOUNT", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = I_IS_IMPLICIT_ACCOUNT.String()
		}
	})

	b.Run("PACK", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = I_PACK.String()
		}
	})
}
