package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestBareIdentifier_RejectNumberLike_MinusDot(t *testing.T) {
	input := `-.`

	_, err := New().Parse(input)

	assert.Error(t, err, "Should reject '-.' as it looks like a number")
}

func TestBareIdentifier_RejectNumberLike_MinusDigit(t *testing.T) {
	input := `-4`

	_, err := New().Parse(input)

	// This looks like a number and should be rejected as a bare node name
	// Number-like patterns must be quoted in KDL v2
	assert.Error(t, err, "Should reject '-4' as it looks like a number")
}

func TestBareIdentifier_RejectNumberLike_PlusDigit(t *testing.T) {
	input := `+123`

	_, err := New().Parse(input)

	// This looks like a number and should be rejected as a bare node name
	// Number-like patterns must be quoted in KDL v2
	assert.Error(t, err, "Should reject '+123' as it looks like a number")
}

// ============================================================
// Malformed Numbers - Validation Tests
// ============================================================

func TestMalformedNumbers_HexWithoutDigits(t *testing.T) {
	input := `node 0x`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hexadecimal number must have at least one digit after 0x")
}

func TestMalformedNumbers_HexWithOnlyUnderscores(t *testing.T) {
	input := `node 0x___`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hexadecimal number must have at least one digit after 0x")
}

func TestMalformedNumbers_BinaryWithoutDigits(t *testing.T) {
	input := `node 0b`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "binary number must have at least one digit after 0b")
}

func TestMalformedNumbers_BinaryWithOnlyUnderscores(t *testing.T) {
	input := `node 0b___`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "binary number must have at least one digit after 0b")
}

func TestMalformedNumbers_OctalWithoutDigits(t *testing.T) {
	input := `node 0o`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "octal number must have at least one digit after 0o")
}

func TestMalformedNumbers_OctalWithOnlyUnderscores(t *testing.T) {
	input := `node 0o___`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "octal number must have at least one digit after 0o")
}

func TestMalformedNumbers_ExponentWithoutDigits(t *testing.T) {
	input := `node 1e`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exponent must have at least one digit")
}

func TestMalformedNumbers_ExponentWithSignButNoDigits(t *testing.T) {
	input := `node 1e+`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exponent must have at least one digit")
}

func TestMalformedNumbers_ExponentWithOnlyUnderscores(t *testing.T) {
	input := `node 1e___`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exponent must have at least one digit")
}

func TestMalformedNumbers_ExponentNegativeSignOnly(t *testing.T) {
	input := `node 2.5e-`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exponent must have at least one digit")
}

func TestMalformedNumbers_UppercaseHexWithoutDigits(t *testing.T) {
	input := `node 0X`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hexadecimal number must have at least one digit after 0x")
}

func TestMalformedNumbers_UppercaseBinaryWithoutDigits(t *testing.T) {
	input := `node 0B`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "binary number must have at least one digit after 0b")
}

func TestMalformedNumbers_UppercaseOctalWithoutDigits(t *testing.T) {
	input := `node 0O`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "octal number must have at least one digit after 0o")
}

func TestMalformedNumbers_SignedHexWithoutDigits(t *testing.T) {
	input := `node +0x`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hexadecimal number must have at least one digit after 0x")
}

func TestMalformedNumbers_SignedBinaryWithoutDigits(t *testing.T) {
	input := `node -0b`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "binary number must have at least one digit after 0b")
}

func TestMalformedNumbers_SignedOctalWithoutDigits(t *testing.T) {
	input := `node +0o`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "octal number must have at least one digit after 0o")
}
