package parser

import "fmt"

// parseKeyword parses a hash-prefixed keyword (#true, #false, #null, #inf, #-inf, #nan)
// Returns the keyword without the # prefix and the value type
func (p *Parser) parseKeyword() (string, ValueType) {
	if p.peek() != '#' {
		p.panicAt("expected # for keyword")
	}
	p.advance() // Skip #

	// Check for negative sign (for #-inf)
	var keyword string
	if p.peek() == '-' {
		p.advance() // Skip -
		ident := p.parseIdentifier()
		keyword = "-" + ident
	} else {
		// Parse the keyword identifier
		keyword = p.parseIdentifier()
	}

	// Determine the type based on the keyword
	switch keyword {
	case "true", "false":
		return keyword, ValueTypeBoolean
	case "null":
		return keyword, ValueTypeNull
	case "inf", "-inf", "nan":
		return keyword, ValueTypeNumber
	default:
		p.panicAt(fmt.Sprintf("unknown keyword: #%s", keyword))
		return "", ValueTypeString // unreachable
	}
}

// parseIdentifier parses an identifier from the current position
func (p *Parser) parseIdentifier() string {
	start := p.pos

	if p.isEOF() || !isIdentifierStart(p.peek()) {
		p.panicAt("expected identifier")
	}

	p.advance()

	for !p.isEOF() && isIdentifierContinue(p.peek()) {
		p.advance()
	}

	return string(p.input[start:p.pos])
}

// parseTypeAnnotation parses a type annotation token inside (...).
// Type annotations can be bare identifiers or quoted strings.
func (p *Parser) parseTypeAnnotation() string {
	if p.isEOF() {
		p.panicAt("expected type annotation")
	}

	if p.peek() == '"' {
		return p.parseQuotedString()
	}

	return p.parseIdentifier()
}

// parseQuotedString parses a quoted string from the current position
// Expects the current position to be at the opening quote (")
// Returns the unescaped string content (without quotes)
// Handles both single-line strings ("...") and multiline strings ("""...""")
func (p *Parser) parseQuotedString() string {
	if p.peek() != '"' {
		p.panicAt("expected opening quote")
	}
	p.advance() // Skip first quote

	// Check if this is a multiline string (""")
	isMultiline := false
	if !p.isEOF() && p.peek() == '"' {
		p.advance() // Skip second quote
		if !p.isEOF() && p.peek() == '"' {
			p.advance() // Skip third quote
			isMultiline = true

			// Multiline strings must be followed by a newline
			if p.isEOF() || p.peek() != '\n' {
				p.panicAt("multiline string must be followed by newline")
			}
			p.advance() // Skip newline
		} else {
			// It was just an empty string ""
			return ""
		}
	}

	if isMultiline {
		return p.parseMultilineString()
	}

	// Single-line string parsing
	var result []rune

	for !p.isEOF() {
		ch := p.peek()

		if ch == '"' {
			// End of string
			p.advance() // Skip closing quote
			return string(result)
		}

		if ch == '\\' {
			// Escape sequence
			p.advance() // Skip backslash
			if p.isEOF() {
				p.panicAt("unterminated string: EOF after backslash")
			}

			escapeChar := p.peek()
			p.advance()

			switch escapeChar {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case '/':
				result = append(result, '/')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'u':
				// Unicode escape: \u{XXXX}
				if p.peek() != '{' {
					p.panicAt("invalid unicode escape: expected '{'")
				}
				p.advance() // Skip '{'

				var hexDigits []rune
				for !p.isEOF() && p.peek() != '}' {
					hexDigits = append(hexDigits, p.peek())
					p.advance()
				}

				if p.isEOF() {
					p.panicAt("unterminated unicode escape")
				}

				p.advance() // Skip '}'

				// Parse hex digits
				var codepoint int
				for _, digit := range hexDigits {
					codepoint *= 16
					if digit >= '0' && digit <= '9' {
						codepoint += int(digit - '0')
					} else if digit >= 'a' && digit <= 'f' {
						codepoint += int(digit-'a') + 10
					} else if digit >= 'A' && digit <= 'F' {
						codepoint += int(digit-'A') + 10
					} else {
						p.panicAt(fmt.Sprintf("invalid hex digit in unicode escape: %c", digit))
					}
				}

				result = append(result, rune(codepoint))
			default:
				p.panicAt(fmt.Sprintf("invalid escape sequence: \\%c", escapeChar))
			}
		} else {
			// Regular character
			result = append(result, ch)
			p.advance()
		}
	}

	p.panicAt("unterminated string: unexpected EOF")
	return "" // unreachable
}

// parseMultilineString parses a multiline string (already past the opening """ and newline)
// Handles dedentation based on the closing quotes' indentation
func (p *Parser) parseMultilineString() string {
	var lines []string
	var currentLine []rune

	// Parse lines until we find the closing """
	for !p.isEOF() {
		ch := p.peek()

		if ch == '"' {
			// Check if this is the closing """
			if p.pos+2 < len(p.input) && p.input[p.pos+1] == '"' && p.input[p.pos+2] == '"' {
				// Found closing """
				// Save the current line first (this is the closing quotes line)
				lines = append(lines, string(currentLine))

				// Skip the closing """
				p.advance()
				p.advance()
				p.advance()

				// Calculate dedentation
				// The indentation of the last line (closing quotes line) determines what to strip
				lastLine := lines[len(lines)-1]
				dedentAmount := 0
				for dedentAmount < len(lastLine) && (lastLine[dedentAmount] == ' ' || lastLine[dedentAmount] == '\t') {
					dedentAmount++
				}

				// Apply dedentation to all lines
				var result []rune
				for i, line := range lines {
					// Strip the common indentation (but only if the line has enough characters)
					stripped := line
					if len(line) >= dedentAmount {
						stripped = line[dedentAmount:]
					}

					// For all lines except the last (which is the closing quotes line),
					// add them with their newlines
					if i < len(lines)-1 {
						result = append(result, []rune(stripped)...)
						result = append(result, '\n')
					} else {
						// Last line is the closing quotes line - only include if non-empty after stripping
						if len(stripped) > 0 {
							result = append(result, []rune(stripped)...)
						}
					}
				}

				return string(result)
			} else {
				// Just a regular quote character
				currentLine = append(currentLine, ch)
				p.advance()
			}
		} else if ch == '\n' {
			// End of line
			lines = append(lines, string(currentLine))
			currentLine = []rune{}
			p.advance()
		} else if ch == '\\' {
			// Escape sequence
			p.advance() // Skip backslash
			if p.isEOF() {
				p.panicAt("unterminated multiline string: EOF after backslash")
			}

			escapeChar := p.peek()
			p.advance()

			switch escapeChar {
			case '"':
				currentLine = append(currentLine, '"')
			case '\\':
				currentLine = append(currentLine, '\\')
			case '/':
				currentLine = append(currentLine, '/')
			case 'n':
				currentLine = append(currentLine, '\n')
			case 'r':
				currentLine = append(currentLine, '\r')
			case 't':
				currentLine = append(currentLine, '\t')
			case 'b':
				currentLine = append(currentLine, '\b')
			case 'f':
				currentLine = append(currentLine, '\f')
			default:
				p.panicAt(fmt.Sprintf("invalid escape sequence in multiline string: \\%c", escapeChar))
			}
		} else {
			// Regular character
			currentLine = append(currentLine, ch)
			p.advance()
		}
	}

	p.panicAt("unterminated multiline string: unexpected EOF")
	return "" // unreachable
}

// parseMultilineRawString parses a multiline raw string
// Expects the opening #"""<newline> to already be consumed
// hashCount is the number of # symbols before the opening """
func (p *Parser) parseMultilineRawString(hashCount int) string {
	var lines []string
	var currentLine []rune

	// Parse lines until we find the closing """#
	for !p.isEOF() {
		ch := p.peek()

		if ch == '"' {
			// Check if this is the closing """
			if p.pos+2 < len(p.input) && p.input[p.pos+1] == '"' && p.input[p.pos+2] == '"' {
				// Found potential closing """
				// Check ahead for matching # symbols before committing
				tempPos := p.pos + 3
				closingHashCount := 0
				for tempPos < len(p.input) && p.input[tempPos] == '#' {
					closingHashCount++
					tempPos++
				}

				// Check if we found the matching closing delimiter
				if closingHashCount == hashCount {
					// This is the real closing delimiter
					// Save the current line first (this is the closing quotes line)
					lines = append(lines, string(currentLine))

					// Skip the closing """
					p.advance()
					p.advance()
					p.advance()

					// Skip the # symbols
					for i := 0; i < closingHashCount; i++ {
						p.advance()
					}

					// Calculate dedentation
					// The indentation of the last line (closing quotes line) determines what to strip
					lastLine := lines[len(lines)-1]
					dedentAmount := 0
					dedentChar := rune(0)

					// Determine the common whitespace prefix from the closing line
					for dedentAmount < len(lastLine) {
						ch := rune(lastLine[dedentAmount])
						if ch != ' ' && ch != '\t' {
							break
						}
						if dedentAmount == 0 {
							dedentChar = ch
						} else if ch != dedentChar {
							// Mixed whitespace characters
							p.panicAt("multiline raw string has mixed whitespace in indentation")
						}
						dedentAmount++
					}

					// Apply dedentation to all lines
					var result []rune
					for i, line := range lines {
						if i == len(lines)-1 {
							// Last line is the closing quotes line - only include if non-empty after stripping
							stripped := line
							if len(line) >= dedentAmount {
								stripped = line[dedentAmount:]
							}
							if len(stripped) > 0 {
								result = append(result, []rune(stripped)...)
							}
						} else {
							// Regular content line
							// Verify the line has the correct indentation if it's not empty
							if len(line) > 0 {
								// Check if the line starts with the expected indentation
								if len(line) < dedentAmount {
									p.panicAt("multiline raw string line has insufficient indentation")
								}
								for j := 0; j < dedentAmount; j++ {
									if rune(line[j]) != dedentChar {
										p.panicAt("multiline raw string has inconsistent indentation characters")
									}
								}
								stripped := line[dedentAmount:]
								result = append(result, []rune(stripped)...)
							}
							result = append(result, '\n')
						}
					}

					return string(result)
				} else {
					// Not the closing delimiter, just treat """ as regular content
					currentLine = append(currentLine, '"')
					p.advance()
				}
			} else {
				// Just a regular quote character
				currentLine = append(currentLine, ch)
				p.advance()
			}
		} else if ch == '\n' {
			// End of line
			lines = append(lines, string(currentLine))
			currentLine = []rune{}
			p.advance()
		} else {
			// Regular character (no escape processing in raw strings)
			currentLine = append(currentLine, ch)
			p.advance()
		}
	}

	p.panicAt("unterminated multiline raw string: unexpected EOF")
	return "" // unreachable
}

// isRawString checks if the current position is at the start of a raw string
// Raw strings start with one or more # followed by a quote: #", ##", ###", etc.
func (p *Parser) isRawString() bool {
	if p.peek() != '#' {
		return false
	}

	// Look ahead to find a quote after one or more #
	pos := p.pos
	for pos < len(p.input) && p.input[pos] == '#' {
		pos++
	}

	// Check if we found a quote after the # symbols
	return pos < len(p.input) && p.input[pos] == '"'
}

// parseRawString parses a raw string from the current position
// Expects the current position to be at the first # before the opening quote
// Raw strings are in the form #"..."#, ##"..."##, ###"..."###, etc.
// Multiline raw strings are in the form #"""..."""#, ##"""..."""##, etc.
// The number of # symbols before and after must match
// Returns the literal string content (without delimiters)
func (p *Parser) parseRawString() string {
	if p.peek() != '#' {
		p.panicAt("expected # for raw string")
	}

	// Count the number of # symbols at the start
	hashCount := 0
	for !p.isEOF() && p.peek() == '#' {
		hashCount++
		p.advance()
	}

	if p.peek() != '"' {
		p.panicAt("expected \" after # for raw string")
	}
	p.advance() // Skip first quote

	// Check if this is a multiline raw string (""")
	isMultiline := false
	if !p.isEOF() && p.peek() == '"' {
		p.advance() // Skip second quote
		if !p.isEOF() && p.peek() == '"' {
			p.advance() // Skip third quote
			isMultiline = true

			// Multiline raw strings must be followed by a newline
			if p.isEOF() || p.peek() != '\n' {
				p.panicAt("multiline raw string must be followed by newline")
			}
			p.advance() // Skip newline
		} else {
			// It was just an empty raw string #""#
			// Check for closing # symbols
			closingHashCount := 0
			for !p.isEOF() && p.peek() == '#' {
				closingHashCount++
				p.advance()
			}
			if closingHashCount == hashCount {
				return ""
			}
			p.panicAt("mismatched # count in raw string")
		}
	}

	if isMultiline {
		return p.parseMultilineRawString(hashCount)
	}

	// Single-line raw string parsing
	var result []rune

	for !p.isEOF() {
		ch := p.peek()

		// Single-line raw strings cannot contain newlines
		if ch == '\n' {
			p.panicAt("single-line raw string cannot contain newline")
		}

		// Look for closing quote followed by matching # count
		if ch == '"' {
			p.advance() // Skip the quote

			// Count # symbols after the quote
			closingHashCount := 0
			for !p.isEOF() && p.peek() == '#' {
				closingHashCount++
				p.advance()
			}

			// Check if we found the matching closing delimiter
			if closingHashCount == hashCount {
				return string(result)
			}

			// Not the closing delimiter, restore position and add to result
			// Add the quote and all the # symbols we consumed
			result = append(result, '"')
			for i := 0; i < closingHashCount; i++ {
				result = append(result, '#')
			}
			// Don't restore position - we've already advanced past the # symbols
		} else {
			result = append(result, ch)
			p.advance()
		}
	}

	p.panicAt("unterminated raw string: unexpected EOF")
	return "" // unreachable
}
