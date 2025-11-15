package parser

import "testing"

// ============================================================
// Tests for forbidden bare identifiers
// These tests verify that the parser correctly rejects
// bare identifiers that are not allowed by the KDL v2 spec
// ============================================================

func TestForbiddenBareIdentifier_StartsWithSemicolon(t *testing.T) {
	input := `;node`

	// Expected: Parse error
	// Semicolon is a forbidden character in bare identifiers

	_, err := New().Parse(input)
	if err == nil {
		t.Errorf("expected parse error for identifier starting with ';'")
	}
}

func TestForbiddenBareIdentifier_StartsWithEquals(t *testing.T) {
	input := `=node`

	// Expected: Parse error
	// Equals sign is a forbidden character in bare identifiers

	_, err := New().Parse(input)
	if err == nil {
		t.Errorf("expected parse error for identifier starting with '='")
	}
}

func TestForbiddenBareIdentifier_StartsWithMinusDot(t *testing.T) {
	input := `-.`

	// Expected: Parse error
	// "-." looks like the start of a number (e.g., "-.5" would be -0.5)
	// and should not be allowed as a bare identifier

	_, err := New().Parse(input)
	if err == nil {
		t.Errorf("expected parse error for identifier '-.'")
	}
}

func TestForbiddenBareIdentifier_StartsWithMinusDigit(t *testing.T) {
	input := `-4`

	// Expected: Parse error
	// "-4" looks like a number and should not be allowed as a bare identifier
	// (It would be parsed as a numeric value argument instead)

	_, err := New().Parse(input)
	if err == nil {
		t.Errorf("expected parse error for identifier '-4'")
	}
}
