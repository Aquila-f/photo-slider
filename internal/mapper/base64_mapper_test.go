package mapper

import (
	"testing"
)

func TestBase64Mapper_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", "src1/vacation/2023"},
		{"with underscores", "src1/vacation_2023"},
		{"slash and underscore collision", "src1/a_b"},
		{"deep path", "source/nested/deep/album"},
		{"empty", ""},
	}

	m := NewBase64Mapper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := m.Encode(tt.input)
			decoded, err := m.Decode(encoded)
			if err != nil {
				t.Fatalf("Decode(%q) error: %v", encoded, err)
			}
			if decoded != tt.input {
				t.Errorf("round-trip failed: got %q, want %q", decoded, tt.input)
			}
		})
	}
}

func TestBase64Mapper_NoCollision(t *testing.T) {
	m := NewBase64Mapper()

	// These two inputs would collide under SlashMapper
	a := m.Encode("src1/vacation_2023")
	b := m.Encode("src1/vacation/2023")

	if a == b {
		t.Errorf("collision: both %q and %q encode to %q", "src1/vacation_2023", "src1/vacation/2023", a)
	}
}

func TestBase64Mapper_DecodeInvalid(t *testing.T) {
	m := NewBase64Mapper()
	_, err := m.Decode("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64 input, got nil")
	}
}
