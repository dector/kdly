package parser

import "fmt"

// parserState represents the current state of the parser
type parserState int

const (
	stDocumentStart parserState = iota
	stNodeName
	stNodeNameQuoted
	stNodeBody      // After node name, parsing arguments/properties/children
	stArgumentValue // Parsing an argument value
	stDocumentEnd
)

// Parser represents a KDL v2 parser
type Parser struct {
	input                    []rune // UTF-8 input as runes
	pos                      int    // current position in input
	line                     int    // current line (1-based)
	col                      int    // current column (1-based)
	state                    parserState
	pendingNodeType          string // temporary storage for node type annotation before node creation
	allowDuplicateProperties bool   // if true, keeps all duplicate properties; if false (default), rightmost wins per KDL v2 spec
	disableTypeAnnotations   bool   // if true, type annotations are not allowed and will cause parse errors
}

// New creates a new Parser instance
func New() *Parser {
	return &Parser{}
}

// WithAllowDuplicateProperties configures the parser to keep all duplicate properties
// instead of following the KDL v2 spec behavior (rightmost wins).
// When set to true, all properties with duplicate keys will be preserved in the Properties slice.
// When set to false (default), duplicate property keys will follow KDL v2 spec: rightmost value overrides.
func (p *Parser) WithAllowDuplicateProperties(allow bool) *Parser {
	p.allowDuplicateProperties = allow
	return p
}

// WithDisableTypeAnnotations configures the parser to reject type annotations.
// When set to true, any type annotation (e.g., (type)node or (type)value) will cause a parse error.
// When set to false (default), type annotations are allowed per KDL v2 spec.
// This is useful for enforcing explicit types in configuration files.
func (p *Parser) WithDisableTypeAnnotations(disable bool) *Parser {
	p.disableTypeAnnotations = disable
	return p
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

// panicAt panics with a formatted error message including position
func (p *Parser) panicAt(message string) {
	panic(fmt.Sprintf("parse error at line %d, col %d: %s", p.line, p.col, message))
}
