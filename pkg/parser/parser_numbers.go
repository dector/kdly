package parser

import (
	"fmt"
	"math/big"
	"strings"
)

// looksLikeNumber checks if the current position starts a numeric literal
// Returns true for: integers (123), floats (1.23), hex (0x1f), binary (0b101), octal (0o77)
// Also handles signs (+123, -456) and scientific notation (1e10, 1.5e-3)
// Also recognizes decimal-only forms like .5, +.5, -.5 so they can
// be validated and rejected by parseNumber.
func (p *Parser) looksLikeNumber() bool {
	pos := p.pos
	if pos >= len(p.input) {
		return false
	}

	ch := p.input[pos]

	// Check for sign
	if ch == '+' || ch == '-' {
		pos++
		if pos >= len(p.input) {
			return false
		}
		ch = p.input[pos]

		// +. or -. patterns (even without following digit) look like numbers
		// and must be quoted according to KDL v2 spec
		if ch == '.' {
			return true
		}
	}

	// Check for hex (0x), binary (0b), or octal (0o) prefix
	if ch == '0' && pos+1 < len(p.input) {
		next := p.input[pos+1]
		if next == 'x' || next == 'X' || next == 'b' || next == 'B' || next == 'o' || next == 'O' {
			return true
		}
	}

	// Can start with a digit or a decimal point (for numbers like .5, +.5, -.5)
	if isDigit(ch) {
		return true
	}

	// Check for decimal point followed by a digit (.5, +.5, -.5)
	if ch == '.' && pos+1 < len(p.input) && isDigit(p.input[pos+1]) {
		return true
	}

	return false
}

// isSignedDotBareIdentifier checks for +. or -. tokens that are not followed
// by a digit, which should be parsed as bare identifiers rather than numbers.
func (p *Parser) isSignedDotBareIdentifier() bool {
	if p.isEOF() {
		return false
	}

	ch := p.peek()
	if ch != '+' && ch != '-' {
		return false
	}

	if p.pos+1 >= len(p.input) || p.input[p.pos+1] != '.' {
		return false
	}

	if p.pos+2 < len(p.input) && isDigit(p.input[p.pos+2]) {
		return false
	}

	return true
}

// isNumberBoundary checks whether a rune can legally follow a numeric literal.
func isNumberBoundary(r rune) bool {
	if isWhitespace(r) {
		return true
	}

	switch r {
	case ';', '{', '}', '/', '\\':
		return true
	}

	return false
}

// validateNumberBoundary ensures a parsed numeric literal is properly delimited.
func (p *Parser) validateNumberBoundary() {
	if p.isEOF() {
		return
	}

	if !isNumberBoundary(p.peek()) {
		p.panicAt(fmt.Sprintf("invalid character '%c' after numeric literal", p.peek()))
	}
}

// parseNumber parses a numeric literal from the current position
// Returns the raw numeric string
func (p *Parser) parseNumber() string {
	start := p.pos

	// Handle optional sign
	hadSign := false
	if p.peek() == '+' || p.peek() == '-' {
		hadSign = true
		p.advance()
	}

	if !p.isEOF() && p.peek() == '.' {
		if hadSign && p.pos+1 < len(p.input) && isDigit(p.input[p.pos+1]) {
			// Signed fractional forms like -.5 and +.5 are allowed.
		} else {
			p.panicAt("numeric literal must have at least one digit before decimal point")
		}
	}

	// Check for special bases
	if p.peek() == '0' && p.pos+1 < len(p.input) {
		next := p.input[p.pos+1]

		// Hexadecimal: 0x or 0X
		if next == 'x' || next == 'X' {
			p.advance() // skip '0'
			p.advance() // skip 'x'
			digitStart := p.pos
			for !p.isEOF() && (isHexDigit(p.peek()) || p.peek() == '_') {
				p.advance()
			}
			numStr := string(p.input[start:p.pos])
			// Validate: must have at least one hex digit after 0x
			if !hasDigitsAfterPrefix(p.input[digitStart:p.pos], isHexDigit) {
				p.panicAt("hexadecimal number must have at least one digit after 0x")
			}
			if p.pos > digitStart && p.input[digitStart] == '_' {
				p.panicAt("hexadecimal number cannot start with underscore after 0x")
			}
			return numStr
		}

		// Binary: 0b or 0B
		if next == 'b' || next == 'B' {
			p.advance() // skip '0'
			p.advance() // skip 'b'
			digitStart := p.pos
			for !p.isEOF() && (isBinaryDigit(p.peek()) || p.peek() == '_') {
				p.advance()
			}
			numStr := string(p.input[start:p.pos])
			// Validate: must have at least one binary digit after 0b
			if !hasDigitsAfterPrefix(p.input[digitStart:p.pos], isBinaryDigit) {
				p.panicAt("binary number must have at least one digit after 0b")
			}
			if p.pos > digitStart && p.input[digitStart] == '_' {
				p.panicAt("binary number cannot start with underscore after 0b")
			}
			return numStr
		}

		// Octal: 0o or 0O
		if next == 'o' || next == 'O' {
			p.advance() // skip '0'
			p.advance() // skip 'o'
			digitStart := p.pos
			for !p.isEOF() && (isOctalDigit(p.peek()) || p.peek() == '_') {
				p.advance()
			}
			numStr := string(p.input[start:p.pos])
			// Validate: must have at least one octal digit after 0o
			if !hasDigitsAfterPrefix(p.input[digitStart:p.pos], isOctalDigit) {
				p.panicAt("octal number must have at least one digit after 0o")
			}
			if p.pos > digitStart && p.input[digitStart] == '_' {
				p.panicAt("octal number cannot start with underscore after 0o")
			}
			return numStr
		}
	}

	// Parse integer part (decimal)
	for !p.isEOF() && (isDigit(p.peek()) || p.peek() == '_') {
		p.advance()
	}

	// Check for decimal point (float)
	if !p.isEOF() && p.peek() == '.' {
		// Make sure it's not a trailing dot (like "123.")
		if p.pos+1 < len(p.input) && isDigit(p.input[p.pos+1]) {
			p.advance() // skip '.'
			for !p.isEOF() && (isDigit(p.peek()) || p.peek() == '_') {
				p.advance()
			}
		}
	}

	// Check for exponent (scientific notation)
	if !p.isEOF() && (p.peek() == 'e' || p.peek() == 'E') {
		p.advance() // skip 'e'

		// Handle optional sign in exponent
		if !p.isEOF() && (p.peek() == '+' || p.peek() == '-') {
			p.advance()
		}

		// Parse exponent digits
		exponentStart := p.pos
		for !p.isEOF() && (isDigit(p.peek()) || p.peek() == '_') {
			p.advance()
		}
		// Validate: exponent must have at least one digit
		if !hasDigitsAfterPrefix(p.input[exponentStart:p.pos], isDigit) {
			p.panicAt("exponent must have at least one digit")
		}
	}

	return string(p.input[start:p.pos])
}

// normalizeNumberValue canonicalizes numeric literals to decimal string form
// and returns the original parsed base.
func normalizeNumberValue(num string) (string, NumberBase) {
	if num == "" {
		return num, NumberBaseDecimal
	}

	sign := ""
	body := num
	if body[0] == '+' || body[0] == '-' {
		sign = body[:1]
		body = body[1:]
		if body == "" {
			return num, NumberBaseDecimal
		}
	}

	base := NumberBaseDecimal
	radix := 10
	if len(body) >= 2 && body[0] == '0' {
		switch body[1] {
		case 'x', 'X':
			base = NumberBaseHexadecimal
			radix = 16
		case 'o', 'O':
			base = NumberBaseOctal
			radix = 8
		case 'b', 'B':
			base = NumberBaseBinary
			radix = 2
		}
	}

	if radix == 10 {
		return num, NumberBaseDecimal
	}

	digits := strings.ReplaceAll(body[2:], "_", "")
	if digits == "" {
		return num, base
	}

	parsed := new(big.Int)
	if _, ok := parsed.SetString(digits, radix); !ok {
		return num, base
	}

	if sign == "-" {
		parsed.Neg(parsed)
	}

	return parsed.String(), base
}

// hasDigitsAfterPrefix checks if the given rune slice contains at least one valid digit
// (excluding underscores) according to the provided digit validator function
func hasDigitsAfterPrefix(runes []rune, isValidDigit func(rune) bool) bool {
	for _, r := range runes {
		if r != '_' && isValidDigit(r) {
			return true
		}
	}
	return false
}
