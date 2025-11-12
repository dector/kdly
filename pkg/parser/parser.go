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
	stDocumentEnd
)

// Parser represents a KDL v2 parser
type Parser struct {
	input []rune // UTF-8 input as runes
	pos   int    // current position in input
	line  int    // current line (1-based)
	col   int    // current column (1-based)
	state parserState
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

// isIdentifierStart checks if a rune can start an identifier
func isIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

// isIdentifierContinue checks if a rune can continue an identifier
func isIdentifierContinue(r rune) bool {
	return isIdentifierStart(r) || (r >= '0' && r <= '9') || r == '-'
}

// panicAt panics with a formatted error message including position
func (p *Parser) panicAt(message string) {
	panic(fmt.Sprintf("parse error at line %d, col %d: %s", p.line, p.col, message))
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
func (p *Parser) parseQuotedString() string {
	if p.peek() != '"' {
		p.panicAt("expected opening quote")
	}
	p.advance() // Skip opening quote

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

			// Add the node to the document
			doc.Nodes = append(doc.Nodes, node)

			// Skip any trailing whitespace
			p.skipWhitespace()

			// For now, we only support single node, so transition to end
			// Later this will be extended to handle arguments, properties, children, etc.
			p.state = stDocumentEnd

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

			// Add the node to the document
			doc.Nodes = append(doc.Nodes, node)

			// Skip any trailing whitespace
			p.skipWhitespace()

			// For now, we only support single node, so transition to end
			// Later this will be extended to handle arguments, properties, children, etc.
			p.state = stDocumentEnd

		default:
			p.panicAt(fmt.Sprintf("unexpected parser state: %d", p.state))
		}
	}

	return doc, nil
}
