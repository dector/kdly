package parser

import "fmt"

func (p *Parser) consumeEscapedWhitespace(first rune) bool {
	consumedNewline := false

	if first == '\r' && !p.isEOF() && p.peek() == '\n' {
		p.advance()
		consumedNewline = true
	}

	if first == '\n' || first == '\r' {
		consumedNewline = true
	}

	for !p.isEOF() {
		ch := p.peek()
		if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			break
		}
		if ch == '\n' || ch == '\r' {
			consumedNewline = true
		}
		p.advance()
	}

	return consumedNewline
}

func (p *Parser) parseUnicodeEscape() rune {
	if p.peek() != '{' {
		p.panicAt("invalid unicode escape: expected '{'")
	}
	p.advance() // Skip '{'

	hexDigits := make([]rune, 0, 6)
	for !p.isEOF() && p.peek() != '}' {
		if len(hexDigits) >= 6 {
			p.panicAt("invalid unicode escape: codepoint must be 1-6 hex digits")
		}
		hexDigits = append(hexDigits, p.peek())
		p.advance()
	}

	if p.isEOF() {
		p.panicAt("unterminated unicode escape")
	}

	p.advance() // Skip '}'

	if len(hexDigits) == 0 {
		p.panicAt("invalid unicode escape: empty codepoint")
	}

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

	if codepoint > 0x10ffff {
		p.panicAt("invalid unicode escape: codepoint out of range")
	}

	if codepoint >= 0xd800 && codepoint <= 0xdfff {
		p.panicAt("invalid unicode escape: surrogate codepoint")
	}

	return rune(codepoint)
}

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

		if ch == '\n' || ch == '\r' {
			p.panicAt("single-line string cannot contain newline")
		}

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
			case 's':
				result = append(result, ' ')
			case 'u':
				result = append(result, p.parseUnicodeEscape())
			case ' ', '\t', '\n', '\r':
				p.consumeEscapedWhitespace(escapeChar)
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
	type multilineLine struct {
		text          []rune
		literalPrefix int
	}

	var lines []multilineLine
	var currentLine []rune
	literalPrefix := 0
	prefixActive := true
	currentLineHasWhitespaceEscape := false
	sawWhitespaceEscape := false

	// Parse lines until we find the closing """
	for !p.isEOF() {
		ch := p.peek()

		if ch == '"' {
			// Check if this is the closing """
			if p.pos+2 < len(p.input) && p.input[p.pos+1] == '"' && p.input[p.pos+2] == '"' {
				closingLine := multilineLine{
					text:          append([]rune(nil), currentLine...),
					literalPrefix: literalPrefix,
				}
				lines = append(lines, closingLine)
				closingLineHasWhitespaceEscape := currentLineHasWhitespaceEscape

				// Skip the closing """
				p.advance()
				p.advance()
				p.advance()

				// Dedent prefix comes from literal leading whitespace on the closing line.
				dedentPrefixLen := 0
				for dedentPrefixLen < len(closingLine.text) && isNonNewlineWhitespace(closingLine.text[dedentPrefixLen]) {
					dedentPrefixLen++
				}
				dedentPrefix := append([]rune(nil), closingLine.text[:dedentPrefixLen]...)

				applyDedent := func(line multilineLine) []rune {
					if len(line.text) == 0 {
						return []rune{}
					}
					if len(line.text) < len(dedentPrefix) {
						p.panicAt("multiline string line has insufficient indentation")
					}
					if line.literalPrefix < len(dedentPrefix) {
						p.panicAt("multiline string has non-literal indentation prefix")
					}
					for i := 0; i < len(dedentPrefix); i++ {
						if line.text[i] != dedentPrefix[i] {
							p.panicAt("multiline string has inconsistent indentation characters")
						}
					}
					return line.text[len(dedentPrefix):]
				}

				var result []rune
				for i, line := range lines {
					if i < len(lines)-1 {
						stripped := applyDedent(line)
						result = append(result, stripped...)
						if i < len(lines)-2 {
							result = append(result, '\n')
						} else if !closingLineHasWhitespaceEscape && !sawWhitespaceEscape {
							result = append(result, '\n')
						}
					} else {
						stripped := applyDedent(line)
						if len(stripped) > 0 {
							result = append(result, stripped...)
						}
					}
				}

				return string(result)
			} else {
				// Just a regular quote character
				currentLine = append(currentLine, ch)
				prefixActive = false
				p.advance()
			}
		} else if ch == '\n' {
			// End of line
			lines = append(lines, multilineLine{text: append([]rune(nil), currentLine...), literalPrefix: literalPrefix})
			currentLine = []rune{}
			literalPrefix = 0
			prefixActive = true
			currentLineHasWhitespaceEscape = false
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
				prefixActive = false
			case '\\':
				currentLine = append(currentLine, '\\')
				prefixActive = false
			case 'n':
				currentLine = append(currentLine, '\n')
				prefixActive = false
			case 'r':
				currentLine = append(currentLine, '\r')
				prefixActive = false
			case 't':
				currentLine = append(currentLine, '\t')
				prefixActive = false
			case 'b':
				currentLine = append(currentLine, '\b')
				prefixActive = false
			case 'f':
				currentLine = append(currentLine, '\f')
				prefixActive = false
			case 's':
				currentLine = append(currentLine, ' ')
				prefixActive = false
			case 'u':
				r := p.parseUnicodeEscape()
				currentLine = append(currentLine, r)
				if prefixActive {
					if isNonNewlineWhitespace(r) {
						prefixActive = false
					} else {
						prefixActive = false
					}
				}
			case ' ', '\t', '\n', '\r':
				consumedNewline := p.consumeEscapedWhitespace(escapeChar)
				currentLineHasWhitespaceEscape = true
				sawWhitespaceEscape = true
				prefixActive = false
				if consumedNewline && p.pos+2 < len(p.input) && p.input[p.pos] == '"' && p.input[p.pos+1] == '"' && p.input[p.pos+2] == '"' {
					for _, r := range currentLine {
						if !isNonNewlineWhitespace(r) {
							p.panicAt("multiline string cannot end with escaped trailing whitespace")
						}
					}
				}
			default:
				p.panicAt(fmt.Sprintf("invalid escape sequence in multiline string: \\%c", escapeChar))
			}
		} else {
			// Regular character
			currentLine = append(currentLine, ch)
			if prefixActive {
				if isNonNewlineWhitespace(ch) {
					literalPrefix++
				} else {
					prefixActive = false
				}
			}
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

					lastLine := []rune(lines[len(lines)-1])
					dedentPrefixLen := 0
					for dedentPrefixLen < len(lastLine) && isNonNewlineWhitespace(lastLine[dedentPrefixLen]) {
						dedentPrefixLen++
					}
					dedentPrefix := append([]rune(nil), lastLine[:dedentPrefixLen]...)

					applyDedent := func(line string) []rune {
						runes := []rune(line)
						if len(runes) == 0 {
							return []rune{}
						}
						if len(runes) < len(dedentPrefix) {
							p.panicAt("multiline raw string line has insufficient indentation")
						}
						for i := 0; i < len(dedentPrefix); i++ {
							if runes[i] != dedentPrefix[i] {
								p.panicAt("multiline raw string has inconsistent indentation characters")
							}
						}
						return runes[len(dedentPrefix):]
					}

					var result []rune
					for i, line := range lines {
						if i == len(lines)-1 {
							stripped := applyDedent(line)
							if len(stripped) > 0 {
								result = append(result, stripped...)
							}
						} else {
							stripped := applyDedent(line)
							result = append(result, stripped...)
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

	// Check if this is a multiline raw string (""").
	isMultiline := false
	if p.pos+1 < len(p.input) && p.peek() == '"' && p.input[p.pos+1] == '"' {
		p.advance() // Skip second quote
		p.advance() // Skip third quote
		isMultiline = true

		// Multiline raw strings must be followed by a newline.
		if p.isEOF() || p.peek() != '\n' {
			p.panicAt("multiline raw string must be followed by newline")
		}
		p.advance() // Skip newline
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
