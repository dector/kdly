package parser

import "fmt"

// Document represents the top-level KDL document
type Document struct {
	Nodes []Node
	// Comments []Comment
}

// Node represents a KDL node with name, arguments, properties, and children
type Node struct {
	Name       string
	Arguments  []Value
	Properties []Property
	Children   []Node
}

// Value represents any KDL value (argument or property value)
type Value struct {
	Type           ValueType // "string", "number", "boolean", "null"
	Value          string    // Raw value as string
	TypeAnnotation string    // Optional type annotation like "u8", "uuid"
}

// ValueType represents the type of a value
type ValueType string

const (
	ValueTypeString  ValueType = "string"
	ValueTypeNumber  ValueType = "number"
	ValueTypeBoolean ValueType = "boolean"
	ValueTypeNull    ValueType = "null"
)

// Property represents a key-value property on a node
type Property struct {
	Key   string
	Value Value
}

// Comment represents a comment in the document
// type Comment struct {
// 	Type    CommentType // "line", "multiline", "slashdash"
// 	Content string      // The comment text
// }

// CommentType represents the type of comment
// type CommentType string

// const (
// 	CommentTypeLine      CommentType = "line"
// 	CommentTypeMultiline CommentType = "multiline"
// 	CommentTypeSlashdash CommentType = "slashdash"
// )

// parserState represents the current state of the parser
type parserState int

const (
	stDocumentStart parserState = iota
	stNodeName
	stNodeNameQuoted
	stNodeBody        // After node name, parsing arguments/properties/children
	stArgumentValue   // Parsing an argument value
	stPropertyValue   // Parsing a property value (after key=)
	stDocumentEnd
)

// Parser represents a KDL v2 parser
type Parser struct {
	input           []rune // UTF-8 input as runes
	pos             int    // current position in input
	line            int    // current line (1-based)
	col             int    // current column (1-based)
	state           parserState
	currentPropKey  string // temporary storage for property key when parsing property value
}

// New creates a new Parser instance
func New() *Parser {
	return &Parser{}
}

// isEOF checks if we've reached the end of input
func (p *Parser) isEOF() bool {
	return p.pos >= len(p.input)
}

// peek returns the current rune without advancing
func (p *Parser) peek() rune {
	if p.isEOF() {
		return 0
	}
	return p.input[p.pos]
}

// advance moves to the next character and updates position tracking
func (p *Parser) advance() {
	if p.isEOF() {
		return
	}

	if p.input[p.pos] == '\n' {
		p.line++
		p.col = 1
	} else {
		p.col++
	}
	p.pos++
}

// peekRune returns the current rune (alias for peek for clarity)
// func (p *Parser) peekRune() rune {
// 	return p.peek()
// }

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

// isSlashdash checks if we're at a slashdash comment (/-)
func (p *Parser) isSlashdash() bool {
	if p.isEOF() {
		return false
	}

	if p.peek() != '/' {
		return false
	}

	if p.pos+1 >= len(p.input) {
		return false
	}

	return p.input[p.pos+1] == '-'
}

// skipSlashdashValue skips a single argument or property value commented with slashdash (/-)
// This can also skip children blocks
func (p *Parser) skipSlashdashValue() {
	// Skip the /- prefix
	p.advance() // skip /
	p.advance() // skip -

	// Skip whitespace between /- and the value
	p.skipInlineWhitespaceAndComments()

	ch := p.peek()

	// Check if this is a children block
	if ch == '{' {
		p.skipChildrenBlock()
		return
	}

	// Check if this is a property (identifier followed by =)
	if isIdentifierStart(ch) && ch != '(' && ch != '"' && ch != '#' {
		// Could be a property or just a bare identifier argument
		_ = p.parseIdentifier()
		p.skipInlineWhitespaceAndComments()

		if !p.isEOF() && p.peek() == '=' {
			// This is a property - skip the = and the value
			p.advance() // skip =
			p.skipInlineWhitespaceAndComments()
			// Skip the property value
			p.parseValueWithOptionalTypeAnnotation()
		}
		// Otherwise, it was just an argument (already consumed)
	} else {
		// Parse as a regular argument value (handles type annotations, strings, numbers, etc.)
		p.parseValueWithOptionalTypeAnnotation()
	}
}

// skipSlashdashNode skips a node that's commented out with slashdash (/-)
func (p *Parser) skipSlashdashNode() {
	// Skip the /- prefix
	p.advance() // skip /
	p.advance() // skip -

	// Skip whitespace between /- and the node
	p.skipInlineWhitespaceAndComments()

	// Skip the node name (quoted or bare identifier)
	if p.peek() == '"' {
		// Skip quoted node name
		p.parseQuotedString()
	} else {
		// Skip bare identifier node name
		p.parseIdentifier()
	}

	// Skip the rest of the node (arguments, properties, children)
	// We need to skip until we hit a node terminator or children block
	for !p.isEOF() {
		p.skipInlineWhitespaceAndComments()

		if p.isEOF() {
			break
		}

		ch := p.peek()

		// Node terminators
		if ch == '\n' || ch == '\r' || ch == ';' {
			p.advance()
			break
		}

		// Children block - need to skip the entire block
		if ch == '{' {
			p.skipChildrenBlock()
			break
		}

		// Skip any value (argument or property)
		// First, try to parse the identifier for the property key
		if ch == '(' || ch == '"' || ch == '#' || p.looksLikeNumber() || isIdentifierStart(ch) {
			// Check if this is a property (identifier followed by =)
			// Try to read what looks like a property key
			if isIdentifierStart(ch) && ch != '(' && ch != '"' && ch != '#' {
				_ = p.parseIdentifier()
				p.skipInlineWhitespaceAndComments()

				if !p.isEOF() && p.peek() == '=' {
					// This is a property
					p.advance() // skip =
					p.skipInlineWhitespaceAndComments()
					// Skip the property value
					p.parseValueWithOptionalTypeAnnotation()
				} else {
					// Not a property, restore and parse as argument
					// We already consumed it as an identifier, which is fine
					// It was an argument value
				}
			} else {
				// Parse as a regular argument value
				p.parseValueWithOptionalTypeAnnotation()
			}
		} else {
			// Unknown character, skip it
			p.advance()
		}
	}
}

// skipChildrenBlock skips an entire children block {...}
func (p *Parser) skipChildrenBlock() {
	if p.peek() != '{' {
		return
	}

	p.advance() // skip opening {

	depth := 1
	for !p.isEOF() && depth > 0 {
		p.skipWhitespaceAndComments()

		if p.isEOF() {
			break
		}

		ch := p.peek()

		if ch == '{' {
			p.advance()
			depth++
		} else if ch == '}' {
			p.advance()
			depth--
		} else if ch == '"' {
			// Skip quoted strings
			p.parseQuotedString()
		} else if p.isRawString() {
			// Skip raw strings
			p.parseRawString()
		} else {
			// Skip any other character
			p.advance()
		}
	}
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

// skipInlineWhitespaceAndComments skips spaces, tabs, and comments but NOT newlines
// However, line continuation (backslash before newline) allows newlines to be treated as whitespace
func (p *Parser) skipInlineWhitespaceAndComments() {
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
					continue    // Continue skipping whitespace after the continuation
				} else if nextCh == '\r' {
					// Handle \r\n or just \r
					p.advance() // skip \
					p.advance() // skip \r
					if !p.isEOF() && p.peek() == '\n' {
						p.advance() // skip \n
					}
					continue // Continue skipping whitespace after the continuation
				}
			}
		}

		// Try to skip comments
		p.skipComments()
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
// - Patterns that look like numbers (digit, +/- followed by digit, etc.)
func isIdentifierStart(r rune) bool {
	// Reject forbidden characters
	if isForbiddenInBareIdentifier(r) {
		return false
	}

	// Reject digits - numbers can start with digits
	if isDigit(r) {
		return false
	}

	// Reject decimal point - numbers like .5, +.5, -.5 start with '.'
	// Note: This is conservative since '.' could be valid in some contexts,
	// but according to KDL v2 spec, patterns like .5 are numeric literals
	if r == '.' {
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

// panicAt panics with a formatted error message including position
func (p *Parser) panicAt(message string) {
	panic(fmt.Sprintf("parse error at line %d, col %d: %s", p.line, p.col, message))
}

// looksLikeNumber checks if the current position starts a numeric literal
// Returns true for: integers (123), floats (1.23), hex (0x1f), binary (0b101), octal (0o77)
// Also handles signs (+123, -456) and scientific notation (1e10, 1.5e-3)
// Also handles decimal-only numbers like .5, +.5, -.5
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

// parseNumber parses a numeric literal from the current position
// Returns the raw numeric string
func (p *Parser) parseNumber() string {
	start := p.pos

	// Handle optional sign
	if p.peek() == '+' || p.peek() == '-' {
		p.advance()
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
						codepoint += int(digit - 'a' + 10)
					} else if digit >= 'A' && digit <= 'F' {
						codepoint += int(digit - 'A' + 10)
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

	// The first line is always empty (we already consumed the newline after """)
	// We need to preserve this as a leading newline in the result
	lines = append(lines, "")

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
	p.advance() // Skip opening quote

	var result []rune

	for !p.isEOF() {
		ch := p.peek()

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

// parseValueWithOptionalTypeAnnotation parses an optional type annotation followed by a value
// Returns the parsed Value with TypeAnnotation field set if present
func (p *Parser) parseValueWithOptionalTypeAnnotation() Value {
	var typeAnnotation string

	// Check for type annotation: (type)
	if p.peek() == '(' {
		p.advance() // skip '('

		// Parse the type annotation (identifier)
		typeAnnotation = p.parseIdentifier()

		// Expect closing ')'
		if p.peek() != ')' {
			p.panicAt("expected ')' after type annotation")
		}
		p.advance() // skip ')'
	}

	// Now parse the actual value
	ch := p.peek()
	var value Value

	if ch == '"' {
		// Quoted string
		strValue := p.parseQuotedString()
		value = Value{
			Type:           ValueTypeString,
			Value:          strValue,
			TypeAnnotation: typeAnnotation,
		}
	} else if p.isRawString() {
		// Raw string
		strValue := p.parseRawString()
		value = Value{
			Type:           ValueTypeString,
			Value:          strValue,
			TypeAnnotation: typeAnnotation,
		}
	} else if ch == '#' {
		// Keyword (#true, #false, #null)
		keyword, valueType := p.parseKeyword()
		value = Value{
			Type:           valueType,
			Value:          keyword,
			TypeAnnotation: typeAnnotation,
		}
	} else if p.looksLikeNumber() {
		// Numeric literal
		numValue := p.parseNumber()
		value = Value{
			Type:           ValueTypeNumber,
			Value:          numValue,
			TypeAnnotation: typeAnnotation,
		}
	} else if isIdentifierStart(ch) {
		// Bare identifier string
		strValue := p.parseIdentifier()
		value = Value{
			Type:           ValueTypeString,
			Value:          strValue,
			TypeAnnotation: typeAnnotation,
		}
	} else {
		p.panicAt("unexpected character in value")
	}

	return value
}

// parseChildren parses a block of child nodes enclosed in {}
// Assumes the opening { has already been consumed
// Returns when the closing } is encountered
func (p *Parser) parseChildren() ([]Node, error) {
	children := make([]Node, 0)

	for {
		p.skipWhitespaceAndComments()

		// Check for closing brace
		if p.isEOF() {
			return nil, fmt.Errorf("unexpected EOF while parsing children block")
		}

		if p.peek() == '}' {
			p.advance() // consume '}'
			return children, nil
		}

		// Parse a child node
		var childNode Node

		// Check if node name is quoted or bare
		if p.peek() == '"' {
			nodeName := p.parseQuotedString()
			childNode = Node{
				Name:       nodeName,
				Arguments:  make([]Value, 0),
				Properties: make([]Property, 0),
				Children:   make([]Node, 0),
			}
		} else if isIdentifierStart(p.peek()) {
			nodeName := p.parseIdentifier()
			childNode = Node{
				Name:       nodeName,
				Arguments:  make([]Value, 0),
				Properties: make([]Property, 0),
				Children:   make([]Node, 0),
			}
		} else {
			return nil, fmt.Errorf("expected node name in children block at line %d, col %d", p.line, p.col)
		}

		// Parse node body (arguments, properties, children)
		for {
			// Skip only spaces and tabs, not newlines (newlines terminate nodes)
			p.skipInlineWhitespaceAndComments()

			if p.isEOF() || p.peek() == '}' {
				// End of this child node
				break
			}

			ch := p.peek()

			// Check for newline or semicolon (node terminators)
			if ch == '\n' || ch == ';' {
				p.advance()
				break
			}

			// Check for children block
			if ch == '{' {
				p.advance() // skip '{'

				grandchildren, err := p.parseChildren()
				if err != nil {
					return nil, err
				}

				childNode.Children = grandchildren
				break
			}

			// Check for arguments and properties
			if ch == '(' || ch == '"' || ch == '#' || p.looksLikeNumber() {
				// This looks like an argument value (possibly with type annotation)
				argValue := p.parseValueWithOptionalTypeAnnotation()
				childNode.Arguments = append(childNode.Arguments, argValue)
			} else if isIdentifierStart(ch) {
				// Could be a bare identifier argument or property key
				savedPos := p.pos
				savedLine := p.line
				savedCol := p.col

				identifier := p.parseIdentifier()

				// Skip only spaces and tabs after identifier
				for !p.isEOF() && (p.peek() == ' ' || p.peek() == '\t') {
					p.advance()
				}

				if !p.isEOF() && p.peek() == '=' {
					// This is a property: key=value
					p.advance() // skip '='

					// Skip spaces after =
					p.skipInlineWhitespaceAndComments()

					// Parse property value (with optional type annotation)
					propValue := p.parseValueWithOptionalTypeAnnotation()
					childNode.Properties = append(childNode.Properties, Property{
						Key:   identifier,
						Value: propValue,
					})
				} else {
					// This is a bare identifier argument - restore position and parse as argument
					p.pos = savedPos
					p.line = savedLine
					p.col = savedCol

					strValue := p.parseIdentifier()
					childNode.Arguments = append(childNode.Arguments, Value{
						Type:  ValueTypeString,
						Value: strValue,
					})
				}
			} else {
				// Unknown character - might be end of node
				break
			}
		}

		children = append(children, childNode)
	}
}

// parseValue parses a single value (used for arguments and properties)
func (p *Parser) parseValue() Value {
	ch := p.peek()

	if ch == '"' {
		// Quoted string
		strValue := p.parseQuotedString()
		return Value{
			Type:  ValueTypeString,
			Value: strValue,
		}
	} else if p.isRawString() {
		// Raw string
		strValue := p.parseRawString()
		return Value{
			Type:  ValueTypeString,
			Value: strValue,
		}
	} else if ch == '#' {
		// Keyword (#true, #false, #null)
		keyword, valueType := p.parseKeyword()
		return Value{
			Type:  valueType,
			Value: keyword,
		}
	} else if p.looksLikeNumber() {
		// Numeric literal
		numValue := p.parseNumber()
		return Value{
			Type:  ValueTypeNumber,
			Value: numValue,
		}
	} else if isIdentifierStart(ch) {
		// Bare identifier string
		strValue := p.parseIdentifier()
		return Value{
			Type:  ValueTypeString,
			Value: strValue,
		}
	} else {
		p.panicAt("unexpected character in value")
		return Value{} // unreachable
	}
}

// Parse parses a KDL document from the provided string
func (p *Parser) Parse(input string) (doc *Document, err error) {
	// Recover from panics and convert to error
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(string); ok {
				err = fmt.Errorf("%s", e)
			} else if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("parse error: %v", r)
			}
		}
	}()

	// Initialize parser state
	p.input = []rune(input)
	p.pos = 0
	p.line = 1
	p.col = 1
	p.state = stDocumentStart

	doc = &Document{
		Nodes: make([]Node, 0),
	}

	var currentNode *Node

	// Main state machine loop
	for p.state != stDocumentEnd {
		switch p.state {
		case stDocumentStart:
			p.skipWhitespaceAndComments()
			if p.isEOF() {
				// Empty document or no more nodes
				p.state = stDocumentEnd
			} else if p.isSlashdash() {
				// Slashdash comment - skip the entire next node
				p.skipSlashdashNode()
				// Stay in stDocumentStart to process the next node
			} else {
				// Transition to parsing node name based on next character
				if p.peek() == '"' {
					p.state = stNodeNameQuoted
				} else {
					// Check if it looks like a number - if so, it's an error
					// because node names can't be numbers
					if p.looksLikeNumber() {
						p.panicAt("node name cannot be a number")
					}
					// Check if the character can start a valid identifier
					if !isIdentifierStart(p.peek()) {
						p.panicAt(fmt.Sprintf("invalid character for node name: '%c'", p.peek()))
					}
					p.state = stNodeName
				}
			}

		case stNodeName:
			// Parse the node name (identifier)
			nodeName := p.parseIdentifier()

			// Create a node with the parsed name
			node := Node{
				Name:       nodeName,
				Arguments:  make([]Value, 0),
				Properties: make([]Property, 0),
				Children:   make([]Node, 0),
			}

			currentNode = &node

			// Transition to parsing node body
			p.state = stNodeBody

		case stNodeNameQuoted:
			// Parse the node name (quoted string)
			nodeName := p.parseQuotedString()

			// Create a node with the parsed name
			node := Node{
				Name:       nodeName,
				Arguments:  make([]Value, 0),
				Properties: make([]Property, 0),
				Children:   make([]Node, 0),
			}

			currentNode = &node

			// Transition to parsing node body
			p.state = stNodeBody

		case stNodeBody:
			// After node name, check for arguments, properties, children, or end
			// Skip only spaces and tabs, not newlines (newlines terminate nodes)
			p.skipInlineWhitespaceAndComments()

			if p.isEOF() {
				// End of document - add current node
				if currentNode != nil {
					doc.Nodes = append(doc.Nodes, *currentNode)
					currentNode = nil
				}
				p.state = stDocumentEnd
			} else if p.isSlashdash() {
				// Slashdash comment - skip the next argument or property
				p.skipSlashdashValue()
				// Stay in stNodeBody to continue parsing
			} else {
				ch := p.peek()

				// Check for newline or semicolon (node terminators)
				if ch == '\n' || ch == '\r' || ch == ';' {
					// Node is complete
					if currentNode != nil {
						doc.Nodes = append(doc.Nodes, *currentNode)
						currentNode = nil
					}
					p.advance() // consume terminator
					p.state = stDocumentStart
				} else if ch == '(' {
					// Type annotation followed by value
					p.state = stArgumentValue
				} else if ch == '"' {
					// Quoted string argument
					p.state = stArgumentValue
				} else if ch == '{' {
					// Children block
					p.advance() // skip '{'

					// Parse children nodes recursively
					children, err := p.parseChildren()
					if err != nil {
						return nil, err
					}

					if currentNode != nil {
						currentNode.Children = children
					}

					// After parsing children, we've consumed the closing '}'
					// Add the current node to the document and move on
					if currentNode != nil {
						doc.Nodes = append(doc.Nodes, *currentNode)
						currentNode = nil
					}
					p.state = stDocumentStart
				} else if ch == '#' {
					// Could be raw string or keyword
					// Peek ahead to see if it's #"
					if p.pos+1 < len(p.input) && p.input[p.pos+1] == '"' {
						// Raw string
						p.state = stArgumentValue
					} else {
						// Keyword (#true, #false, #null)
						p.state = stArgumentValue
					}
				} else if p.looksLikeNumber() {
					// Numeric literal - check this BEFORE isIdentifierStart
					p.state = stArgumentValue
				} else if isIdentifierStart(ch) {
					// Could be a bare identifier argument or property key
					// Look ahead to see if there's an = sign after the identifier
					savedPos := p.pos
					savedLine := p.line
					savedCol := p.col

					identifier := p.parseIdentifier()

					// Skip only spaces and tabs after identifier
					p.skipInlineWhitespaceAndComments()

					if !p.isEOF() && p.peek() == '=' {
						// This is a property: key=value
						p.currentPropKey = identifier
						p.advance() // skip '='
						p.state = stPropertyValue
					} else {
						// This is a bare identifier argument - restore position and parse as argument
						p.pos = savedPos
						p.line = savedLine
						p.col = savedCol
						p.state = stArgumentValue
					}
				} else {
					// End of node
					if currentNode != nil {
						doc.Nodes = append(doc.Nodes, *currentNode)
						currentNode = nil
					}
					p.state = stDocumentEnd
				}
			}

		case stArgumentValue:
			// Parse an argument value (with optional type annotation)
			argValue := p.parseValueWithOptionalTypeAnnotation()

			// Add argument to current node
			if currentNode != nil {
				currentNode.Arguments = append(currentNode.Arguments, argValue)
			}

			// Transition back to node body to check for more arguments
			p.state = stNodeBody

		case stPropertyValue:
			// Parse a property value (after key=)
			p.skipInlineWhitespaceAndComments()

			var propValue Value

			// Check for slashdash commenting out the value
			if p.isSlashdash() {
				// Skip the slashdash and the value it comments out
				p.skipSlashdashValue()
				// Use an empty string value as placeholder
				propValue = Value{
					Type:  ValueTypeString,
					Value: "",
				}
			} else {
				// Parse property value (with optional type annotation)
				propValue = p.parseValueWithOptionalTypeAnnotation()
			}

			// Add property to current node
			if currentNode != nil {
				prop := Property{
					Key:   p.currentPropKey,
					Value: propValue,
				}
				currentNode.Properties = append(currentNode.Properties, prop)
			}

			// Clear the temporary property key
			p.currentPropKey = ""

			// Transition back to node body to check for more arguments/properties
			p.state = stNodeBody

		default:
			p.panicAt(fmt.Sprintf("unexpected parser state: %d", p.state))
		}
	}

	return doc, nil
}
