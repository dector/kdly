package parser

// isWhitespace checks if a rune is whitespace (space, tab, newline, carriage return)
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// skipWhitespace skips all whitespace characters
func (p *Parser) skipWhitespace() {
	for !p.isEOF() && isWhitespace(p.peek()) {
		p.advance()
	}
}

// skipComments skips line comments (//) and multiline comments (/* */)
// Note: This does NOT handle slashdash (/-) comments, which are handled separately
func (p *Parser) skipComments() bool {
	if p.isEOF() {
		return false
	}

	ch := p.peek()
	if ch != '/' {
		return false
	}

	// Need to check the next character
	if p.pos+1 >= len(p.input) {
		return false
	}

	nextCh := p.input[p.pos+1]

	if nextCh == '/' {
		// Line comment - skip until end of line or EOF
		p.advance() // skip first /
		p.advance() // skip second /
		for !p.isEOF() && p.peek() != '\n' {
			p.advance()
		}
		// Optionally skip the newline itself
		if !p.isEOF() && p.peek() == '\n' {
			p.advance()
		}
		return true
	} else if nextCh == '*' {
		// Multiline comment - skip until matching */
		// KDL v2 supports nested multiline comments
		p.advance() // skip /
		p.advance() // skip *

		depth := 1
		for !p.isEOF() && depth > 0 {
			if p.peek() == '/' && p.pos+1 < len(p.input) && p.input[p.pos+1] == '*' {
				// Found nested comment start
				p.advance() // skip /
				p.advance() // skip *
				depth++
			} else if p.peek() == '*' && p.pos+1 < len(p.input) && p.input[p.pos+1] == '/' {
				// Found comment end
				p.advance() // skip *
				p.advance() // skip /
				depth--
			} else {
				p.advance()
			}
		}
		return true
	}

	return false
}

// skipWhitespaceAndComments skips both whitespace and comments
func (p *Parser) skipWhitespaceAndComments() {
	for {
		startPos := p.pos
		p.skipWhitespace()
		p.skipComments()
		// If position didn't change, we're done
		if p.pos == startPos {
			break
		}
	}
}

// skipInlineComments skips comments in inline contexts.
// For line comments (//), it stops before the newline so callers can still
// observe newline as a node terminator.
func (p *Parser) skipInlineComments() bool {
	if p.isEOF() || p.peek() != '/' || p.pos+1 >= len(p.input) {
		return false
	}

	nextCh := p.input[p.pos+1]

	if nextCh == '/' {
		p.advance() // skip first /
		p.advance() // skip second /
		for !p.isEOF() && p.peek() != '\n' && p.peek() != '\r' {
			p.advance()
		}
		return true
	}

	if nextCh == '*' {
		p.advance() // skip /
		p.advance() // skip *

		depth := 1
		for !p.isEOF() && depth > 0 {
			if p.peek() == '/' && p.pos+1 < len(p.input) && p.input[p.pos+1] == '*' {
				p.advance() // skip /
				p.advance() // skip *
				depth++
			} else if p.peek() == '*' && p.pos+1 < len(p.input) && p.input[p.pos+1] == '/' {
				p.advance() // skip *
				p.advance() // skip /
				depth--
			} else {
				p.advance()
			}
		}
		return true
	}

	return false
}

// skipInlineWhitespaceAndComments skips spaces, tabs, and comments but NOT newlines
// However, line continuation (backslash before newline) allows newlines to be treated as whitespace
func (p *Parser) skipInlineWhitespaceAndComments() {
	allowCommentNewline := false

	for {
		startPos := p.pos
		// Skip spaces and tabs only (not newlines)
		for !p.isEOF() && (p.peek() == ' ' || p.peek() == '\t') {
			p.advance()
		}

		// Check for line continuation: backslash followed by newline
		if !p.isEOF() && p.peek() == '\\' {
			// Look ahead to see if there's a newline
			if p.pos+1 < len(p.input) {
				nextCh := p.input[p.pos+1]
				if nextCh == '\n' {
					// Line continuation - skip backslash and newline
					p.advance() // skip \
					p.advance() // skip \n
					allowCommentNewline = true
					continue // Continue skipping whitespace after the continuation
				} else if nextCh == '\r' {
					// Handle \r\n or just \r
					p.advance() // skip \
					p.advance() // skip \r
					if !p.isEOF() && p.peek() == '\n' {
						p.advance() // skip \n
					}
					allowCommentNewline = true
					continue // Continue skipping whitespace after the continuation
				}
			}
		}

		// Try to skip comments
		if p.skipInlineComments() {
			if allowCommentNewline && !p.isEOF() {
				if p.peek() == '\n' {
					p.advance()
					continue
				}
				if p.peek() == '\r' {
					p.advance()
					if !p.isEOF() && p.peek() == '\n' {
						p.advance()
					}
					continue
				}
			}
		}

		// If position didn't change, we're done
		if p.pos == startPos {
			break
		}
	}
}

// isForbiddenInBareIdentifier checks if a rune is forbidden in bare identifiers
// According to KDL v2 spec, bare identifiers cannot contain:
// - Whitespace
// - Reserved syntax characters: []{}()\/#";=
func isForbiddenInBareIdentifier(r rune) bool {
	if isWhitespace(r) {
		return true
	}
	// Check reserved syntax characters
	switch r {
	case '[', ']', '{', '}', '(', ')', '\\', '/', '#', '"', ';', '=':
		return true
	}
	return false
}

// isIdentifierStart checks if a rune can start an identifier
// In KDL v2, bare identifiers can start with almost any character except:
// - Whitespace and forbidden characters (checked by isForbiddenInBareIdentifier)
// - Digits (number parsing has precedence for numeric-looking tokens)
func isIdentifierStart(r rune) bool {
	// Reject forbidden characters
	if isForbiddenInBareIdentifier(r) {
		return false
	}

	// Reject digits - numbers can start with digits
	if isDigit(r) {
		return false
	}

	// Note: We cannot reject '+' and '-' here because they're only number-like
	// when followed by a digit or '.', which requires lookahead.
	// Use looksLikeNumber() for full number pattern validation.

	return true
}

// isIdentifierContinue checks if a rune can continue an identifier
// In KDL v2, the same rules apply for continuation as for start
func isIdentifierContinue(r rune) bool {
	return !isForbiddenInBareIdentifier(r)
}

// isDigit checks if a rune is a decimal digit (0-9)
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// isHexDigit checks if a rune is a hexadecimal digit (0-9, a-f, A-F)
func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

// isBinaryDigit checks if a rune is a binary digit (0-1)
func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1'
}

// isOctalDigit checks if a rune is an octal digit (0-7)
func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}
