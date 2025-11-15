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
	// For now, accept any non-forbidden character
	// TODO: Need to reject patterns that look like numbers
	return !isForbiddenInBareIdentifier(r)
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
	}

	// Check for hex (0x), binary (0b), or octal (0o) prefix
	if ch == '0' && pos+1 < len(p.input) {
		next := p.input[pos+1]
		if next == 'x' || next == 'X' || next == 'b' || next == 'B' || next == 'o' || next == 'O' {
			return true
		}
	}

	// Must start with a digit
	return isDigit(ch)
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
			for !p.isEOF() && (isHexDigit(p.peek()) || p.peek() == '_') {
				p.advance()
			}
			return string(p.input[start:p.pos])
		}

		// Binary: 0b or 0B
		if next == 'b' || next == 'B' {
			p.advance() // skip '0'
			p.advance() // skip 'b'
			for !p.isEOF() && (isBinaryDigit(p.peek()) || p.peek() == '_') {
				p.advance()
			}
			return string(p.input[start:p.pos])
		}

		// Octal: 0o or 0O
		if next == 'o' || next == 'O' {
			p.advance() // skip '0'
			p.advance() // skip 'o'
			for !p.isEOF() && (isOctalDigit(p.peek()) || p.peek() == '_') {
				p.advance()
			}
			return string(p.input[start:p.pos])
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
		for !p.isEOF() && (isDigit(p.peek()) || p.peek() == '_') {
			p.advance()
		}
	}

	return string(p.input[start:p.pos])
}

// parseKeyword parses a hash-prefixed keyword (#true, #false, #null)
// Returns the keyword without the # prefix and the value type
func (p *Parser) parseKeyword() (string, ValueType) {
	if p.peek() != '#' {
		p.panicAt("expected # for keyword")
	}
	p.advance() // Skip #

	// Parse the keyword identifier
	keyword := p.parseIdentifier()

	// Determine the type based on the keyword
	switch keyword {
	case "true", "false":
		return keyword, ValueTypeBoolean
	case "null":
		return keyword, ValueTypeNull
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

// parseRawString parses a raw string from the current position
// Expects the current position to be at the # before the opening quote
// Raw strings are in the form #"..."# and don't process escape sequences
// Returns the literal string content (without delimiters)
func (p *Parser) parseRawString() string {
	if p.peek() != '#' {
		p.panicAt("expected # for raw string")
	}
	p.advance() // Skip #

	if p.peek() != '"' {
		p.panicAt("expected \" after # for raw string")
	}
	p.advance() // Skip opening quote

	var result []rune

	for !p.isEOF() {
		ch := p.peek()

		// Look for closing "#
		if ch == '"' {
			p.advance()
			if !p.isEOF() && p.peek() == '#' {
				p.advance() // Skip closing #
				return string(result)
			} else {
				// Just a quote in the middle, not the end
				result = append(result, '"')
			}
		} else {
			result = append(result, ch)
			p.advance()
		}
	}

	p.panicAt("unterminated raw string: unexpected EOF")
	return "" // unreachable
}

// parseChildren parses a block of child nodes enclosed in {}
// Assumes the opening { has already been consumed
// Returns when the closing } is encountered
func (p *Parser) parseChildren() ([]Node, error) {
	children := make([]Node, 0)

	for {
		p.skipWhitespace()

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
			for !p.isEOF() && (p.peek() == ' ' || p.peek() == '\t') {
				p.advance()
			}

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
			if ch == '"' {
				// Quoted string argument
				strValue := p.parseQuotedString()
				childNode.Arguments = append(childNode.Arguments, Value{
					Type:  ValueTypeString,
					Value: strValue,
				})
			} else if ch == '#' {
				// Could be raw string or keyword
				if p.pos+1 < len(p.input) && p.input[p.pos+1] == '"' {
					// Raw string
					strValue := p.parseRawString()
					childNode.Arguments = append(childNode.Arguments, Value{
						Type:  ValueTypeString,
						Value: strValue,
					})
				} else {
					// Keyword (#true, #false, #null)
					keyword, valueType := p.parseKeyword()
					childNode.Arguments = append(childNode.Arguments, Value{
						Type:  valueType,
						Value: keyword,
					})
				}
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
					for !p.isEOF() && (p.peek() == ' ' || p.peek() == '\t') {
						p.advance()
					}

					// Parse property value
					propValue := p.parseValue()
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
			} else if p.looksLikeNumber() {
				// Numeric literal
				numValue := p.parseNumber()
				childNode.Arguments = append(childNode.Arguments, Value{
					Type:  ValueTypeNumber,
					Value: numValue,
				})
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
	} else if ch == '#' && p.pos+1 < len(p.input) && p.input[p.pos+1] == '"' {
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
func (p *Parser) Parse(input string) (*Document, error) {
	// Initialize parser state
	p.input = []rune(input)
	p.pos = 0
	p.line = 1
	p.col = 1
	p.state = stDocumentStart

	doc := &Document{
		Nodes: make([]Node, 0),
	}

	var currentNode *Node

	// Main state machine loop
	for p.state != stDocumentEnd {
		switch p.state {
		case stDocumentStart:
			p.skipWhitespace()
			if p.isEOF() {
				// Empty document or no more nodes
				p.state = stDocumentEnd
			} else {
				// Transition to parsing node name based on next character
				if p.peek() == '"' {
					p.state = stNodeNameQuoted
				} else {
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
			p.skipWhitespace()

			if p.isEOF() {
				// End of document - add current node
				if currentNode != nil {
					doc.Nodes = append(doc.Nodes, *currentNode)
					currentNode = nil
				}
				p.state = stDocumentEnd
			} else {
				ch := p.peek()

				// Check what comes next
				if ch == '"' {
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
				} else if isIdentifierStart(ch) {
					// Could be a bare identifier argument or property key
					// Look ahead to see if there's an = sign after the identifier
					savedPos := p.pos
					savedLine := p.line
					savedCol := p.col

					identifier := p.parseIdentifier()
					p.skipWhitespace()

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
			// Parse an argument value
			ch := p.peek()

			var argValue Value

			if ch == '"' {
				// Quoted string
				strValue := p.parseQuotedString()
				argValue = Value{
					Type:  ValueTypeString,
					Value: strValue,
				}
			} else if ch == '#' && p.pos+1 < len(p.input) && p.input[p.pos+1] == '"' {
				// Raw string
				strValue := p.parseRawString()
				argValue = Value{
					Type:  ValueTypeString,
					Value: strValue,
				}
			} else if ch == '#' {
				// Keyword (#true, #false, #null)
				keyword, valueType := p.parseKeyword()
				argValue = Value{
					Type:  valueType,
					Value: keyword,
				}
			} else if p.looksLikeNumber() {
				// Numeric literal
				numValue := p.parseNumber()
				argValue = Value{
					Type:  ValueTypeNumber,
					Value: numValue,
				}
			} else if isIdentifierStart(ch) {
				// Bare identifier string
				strValue := p.parseIdentifier()
				argValue = Value{
					Type:  ValueTypeString,
					Value: strValue,
				}
			} else {
				p.panicAt("unexpected character in argument value")
			}

			// Add argument to current node
			if currentNode != nil {
				currentNode.Arguments = append(currentNode.Arguments, argValue)
			}

			// Transition back to node body to check for more arguments
			p.state = stNodeBody

		case stPropertyValue:
			// Parse a property value (after key=)
			p.skipWhitespace()
			ch := p.peek()

			var propValue Value

			if ch == '"' {
				// Quoted string
				strValue := p.parseQuotedString()
				propValue = Value{
					Type:  ValueTypeString,
					Value: strValue,
				}
			} else if ch == '#' && p.pos+1 < len(p.input) && p.input[p.pos+1] == '"' {
				// Raw string
				strValue := p.parseRawString()
				propValue = Value{
					Type:  ValueTypeString,
					Value: strValue,
				}
			} else if ch == '#' {
				// Keyword (#true, #false, #null)
				keyword, valueType := p.parseKeyword()
				propValue = Value{
					Type:  valueType,
					Value: keyword,
				}
			} else if p.looksLikeNumber() {
				// Numeric literal
				numValue := p.parseNumber()
				propValue = Value{
					Type:  ValueTypeNumber,
					Value: numValue,
				}
			} else if isIdentifierStart(ch) {
				// Bare identifier string
				strValue := p.parseIdentifier()
				propValue = Value{
					Type:  ValueTypeString,
					Value: strValue,
				}
			} else {
				p.panicAt("unexpected character in property value")
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
