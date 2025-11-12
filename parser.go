package kdly

import (
	"fmt"
	"strconv"
	"strings"
)

// Parser parses a KDL document.
type Parser struct {
	lexer   *Lexer
	current Token
	peek    Token
	errors  []error
}

// NewParser creates a new parser for the given input.
func NewParser(input string) *Parser {
	p := &Parser{
		lexer:  NewLexer(input),
		errors: make([]error, 0),
	}
	// Read two tokens to initialize current and peek
	p.nextToken()
	p.nextToken()
	return p
}

// Parse parses the entire document and returns a Document.
func (p *Parser) Parse() (*Document, error) {
	doc := NewDocument()

	for !p.currentTokenIs(TokenEOF) {
		// Skip newlines and comments at document level
		if p.currentTokenIs(TokenNewline) || p.current.IsComment() {
			p.nextToken()
			continue
		}

		// Handle slashdash comment - skip the next node
		if p.currentTokenIs(TokenSlashDash) {
			p.nextToken()
			p.skipNode()
			continue
		}

		node := p.parseNode()
		if node != nil {
			doc.Nodes = append(doc.Nodes, node)
		}
	}

	if len(p.errors) > 0 {
		return doc, fmt.Errorf("parsing errors: %v", p.errors)
	}

	return doc, nil
}

// parseNode parses a single node.
func (p *Parser) parseNode() *Node {
	// Parse optional type annotation
	var typeAnnotation *string
	if p.currentTokenIs(TokenLeftParen) {
		typeAnnotation = p.parseTypeAnnotation()
	}

	// Parse node name
	if !p.currentTokenIs(TokenIdentifier) {
		p.addError(fmt.Errorf("expected node name at %s, got %s", p.current.Position, p.current.Type))
		p.skipToNextNode()
		return nil
	}

	node := NewNode(p.current.Literal)
	node.TypeAnnotation = typeAnnotation
	p.nextToken()

	// Parse arguments and properties
	for {
		// Skip whitespace but not newlines
		if p.current.IsComment() {
			p.nextToken()
			continue
		}

		// Handle slashdash comment - skip the next value/property
		if p.currentTokenIs(TokenSlashDash) {
			p.nextToken()
			p.skipValue()
			continue
		}

		// Check for line continuation
		if p.currentTokenIs(TokenBackslash) {
			p.nextToken()
			// Skip newline after backslash
			if p.currentTokenIs(TokenNewline) {
				p.nextToken()
			}
			continue
		}

		// Check for node termination
		if p.current.IsNodeTerminator() {
			if p.currentTokenIs(TokenNewline) || p.currentTokenIs(TokenSemicolon) {
				p.nextToken()
			}
			break
		}

		// Check for children block
		if p.currentTokenIs(TokenLeftBrace) {
			p.parseChildren(node)
			// After children, node is terminated
			if p.currentTokenIs(TokenNewline) || p.currentTokenIs(TokenSemicolon) {
				p.nextToken()
			}
			break
		}

		// Try to parse property (key=value)
		if p.currentTokenIs(TokenIdentifier) && p.peekTokenIs(TokenEquals) {
			p.parseProperty(node)
			continue
		}

		// Try to parse type-annotated property or value ((type)key=value or (type)value)
		if p.currentTokenIs(TokenLeftParen) {
			// Parse the type annotation
			typeAnnotation := p.parseTypeAnnotation()
			if p.currentTokenIs(TokenIdentifier) && p.peekTokenIs(TokenEquals) {
				// It's a property with type annotation
				key := p.current.Literal
				p.nextToken() // consume identifier
				p.nextToken() // consume =
				value := p.parseValueWithoutTypeAnnotation()
				if value != nil {
					value.TypeAnnotation = typeAnnotation
					node.Properties[key] = value
				}
				continue
			} else {
				// It's an argument with type annotation
				value := p.parseValueWithoutTypeAnnotation()
				if value != nil {
					value.TypeAnnotation = typeAnnotation
					node.Arguments = append(node.Arguments, value)
				} else {
					break
				}
				continue
			}
		}

		// Otherwise, parse as an argument
		value := p.parseValue()
		if value != nil {
			node.Arguments = append(node.Arguments, value)
		} else {
			break
		}
	}

	return node
}

// parseProperty parses a property (key=value).
func (p *Parser) parseProperty(node *Node) {
	key := p.current.Literal
	p.nextToken() // consume identifier
	p.nextToken() // consume =

	value := p.parseValue()
	if value != nil {
		node.Properties[key] = value
	}
}

// parseValue parses a value (with optional type annotation).
func (p *Parser) parseValue() *Value {
	// Check for type annotation
	var typeAnnotation *string
	if p.currentTokenIs(TokenLeftParen) {
		typeAnnotation = p.parseTypeAnnotation()
	}

	var value *Value

	switch p.current.Type {
	case TokenString:
		value = p.parseStringValue()
	case TokenNumber:
		value = p.parseNumberValue()
	case TokenTrue:
		value = NewBooleanValue(true)
		p.nextToken()
	case TokenFalse:
		value = NewBooleanValue(false)
		p.nextToken()
	case TokenNull:
		value = NewNullValue()
		p.nextToken()
	case TokenIdentifier:
		// Unquoted string (identifier)
		value = NewStringValue(p.current.Literal)
		p.nextToken()
	default:
		p.addError(fmt.Errorf("expected value at %s, got %s", p.current.Position, p.current.Type))
		return nil
	}

	if value != nil && typeAnnotation != nil {
		value.TypeAnnotation = typeAnnotation
	}

	return value
}

// parseValueWithoutTypeAnnotation parses a value without checking for a type annotation first.
// This is used when the type annotation has already been parsed.
func (p *Parser) parseValueWithoutTypeAnnotation() *Value {
	var value *Value

	switch p.current.Type {
	case TokenString:
		value = p.parseStringValue()
	case TokenNumber:
		value = p.parseNumberValue()
	case TokenTrue:
		value = NewBooleanValue(true)
		p.nextToken()
	case TokenFalse:
		value = NewBooleanValue(false)
		p.nextToken()
	case TokenNull:
		value = NewNullValue()
		p.nextToken()
	case TokenIdentifier:
		// Unquoted string (identifier)
		value = NewStringValue(p.current.Literal)
		p.nextToken()
	default:
		p.addError(fmt.Errorf("expected value at %s, got %s", p.current.Position, p.current.Type))
		return nil
	}

	return value
}

// parseStringValue parses a string value.
func (p *Parser) parseStringValue() *Value {
	literal := p.current.Literal
	var stringValue string

	// Extract the actual string content based on the format
	if strings.HasPrefix(literal, "\"\"\"") {
		// Multiline string - remove """ at start and end, handle dedentation
		content := literal[3 : len(literal)-3]
		stringValue = dedentMultilineString(content)
	} else if strings.HasPrefix(literal, "#") {
		// Raw string - remove #"...# wrapper
		hashCount := 0
		for i := 0; i < len(literal) && literal[i] == '#'; i++ {
			hashCount++
		}
		start := hashCount + 1 // skip hashes and opening "
		end := len(literal) - hashCount - 1 // skip closing " and hashes
		stringValue = literal[start:end]
	} else if strings.HasPrefix(literal, "\"") {
		// Regular quoted string - already processed escape sequences in lexer
		// For now, just strip quotes (we should actually process escapes here)
		stringValue = literal[1 : len(literal)-1]
	} else {
		// Unquoted string
		stringValue = literal
	}

	p.nextToken()
	return NewStringValue(stringValue)
}

// parseNumberValue parses a number value.
func (p *Parser) parseNumberValue() *Value {
	literal := p.current.Literal
	cleanLiteral := strings.ReplaceAll(literal, "_", "")

	var numValue interface{}
	var err error

	// Determine number format
	if strings.HasPrefix(cleanLiteral, "0x") || strings.HasPrefix(cleanLiteral, "0X") {
		// Hexadecimal
		numValue, err = strconv.ParseInt(cleanLiteral[2:], 16, 64)
	} else if strings.HasPrefix(cleanLiteral, "0o") || strings.HasPrefix(cleanLiteral, "0O") {
		// Octal
		numValue, err = strconv.ParseInt(cleanLiteral[2:], 8, 64)
	} else if strings.HasPrefix(cleanLiteral, "0b") || strings.HasPrefix(cleanLiteral, "0B") {
		// Binary
		numValue, err = strconv.ParseInt(cleanLiteral[2:], 2, 64)
	} else if strings.Contains(cleanLiteral, ".") || strings.ContainsAny(cleanLiteral, "eE") {
		// Float
		numValue, err = strconv.ParseFloat(cleanLiteral, 64)
	} else {
		// Integer
		numValue, err = strconv.ParseInt(cleanLiteral, 10, 64)
	}

	if err != nil {
		p.addError(fmt.Errorf("failed to parse number %q at %s: %v", literal, p.current.Position, err))
		p.nextToken()
		return nil
	}

	p.nextToken()
	return NewNumberValue(literal, numValue)
}

// parseTypeAnnotation parses a type annotation (type).
func (p *Parser) parseTypeAnnotation() *string {
	if !p.currentTokenIs(TokenLeftParen) {
		return nil
	}
	p.nextToken() // consume (

	if !p.currentTokenIs(TokenIdentifier) {
		p.addError(fmt.Errorf("expected type name at %s", p.current.Position))
		return nil
	}

	typeName := p.current.Literal
	p.nextToken() // consume type name

	if !p.currentTokenIs(TokenRightParen) {
		p.addError(fmt.Errorf("expected ')' at %s", p.current.Position))
		return nil
	}
	p.nextToken() // consume )

	return &typeName
}

// parseChildren parses a children block { ... }.
func (p *Parser) parseChildren(parent *Node) {
	if !p.currentTokenIs(TokenLeftBrace) {
		return
	}
	p.nextToken() // consume {

	for !p.currentTokenIs(TokenRightBrace) && !p.currentTokenIs(TokenEOF) {
		// Skip newlines and comments
		if p.currentTokenIs(TokenNewline) || p.current.IsComment() {
			p.nextToken()
			continue
		}

		// Handle slashdash comment - skip the next node
		if p.currentTokenIs(TokenSlashDash) {
			p.nextToken()
			p.skipNode()
			continue
		}

		child := p.parseNode()
		if child != nil {
			parent.Children = append(parent.Children, child)
		}
	}

	if p.currentTokenIs(TokenRightBrace) {
		p.nextToken() // consume }
	} else {
		p.addError(fmt.Errorf("expected '}' at %s", p.current.Position))
	}
}

// skipNode skips the current node (used for slashdash comments).
func (p *Parser) skipNode() {
	// Skip type annotation if present
	if p.currentTokenIs(TokenLeftParen) {
		p.parseTypeAnnotation()
	}

	// Skip node name
	if p.currentTokenIs(TokenIdentifier) {
		p.nextToken()
	}

	// Skip arguments and properties
	for !p.current.IsNodeTerminator() && !p.currentTokenIs(TokenLeftBrace) && !p.currentTokenIs(TokenEOF) {
		p.skipValue()
	}

	// Skip children if present
	if p.currentTokenIs(TokenLeftBrace) {
		depth := 1
		p.nextToken()
		for depth > 0 && !p.currentTokenIs(TokenEOF) {
			if p.currentTokenIs(TokenLeftBrace) {
				depth++
			} else if p.currentTokenIs(TokenRightBrace) {
				depth--
			}
			p.nextToken()
		}
	}

	// Skip terminator
	if p.currentTokenIs(TokenNewline) || p.currentTokenIs(TokenSemicolon) {
		p.nextToken()
	}
}

// skipValue skips a single value.
func (p *Parser) skipValue() {
	// Skip type annotation if present
	if p.currentTokenIs(TokenLeftParen) {
		p.parseTypeAnnotation()
	}

	// Skip property key= if present
	if p.currentTokenIs(TokenIdentifier) && p.peekTokenIs(TokenEquals) {
		p.nextToken() // skip identifier
		p.nextToken() // skip =
	}

	// Skip the value
	if p.current.IsValue() || p.currentTokenIs(TokenIdentifier) {
		p.nextToken()
	}
}

// skipToNextNode skips tokens until the next node or EOF.
func (p *Parser) skipToNextNode() {
	for !p.current.IsNodeTerminator() && !p.currentTokenIs(TokenEOF) {
		p.nextToken()
	}
	if p.currentTokenIs(TokenNewline) || p.currentTokenIs(TokenSemicolon) {
		p.nextToken()
	}
}

// Token management

func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.current.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peek.Type == t
}

// Error handling

func (p *Parser) addError(err error) {
	p.errors = append(p.errors, err)
}

// dedentMultilineString removes common leading whitespace from multiline strings.
func dedentMultilineString(s string) string {
	lines := strings.Split(s, "\n")

	// Remove first line if empty
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}

	// Remove last line if empty
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) == 0 {
		return ""
	}

	// Find minimum indentation
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := 0
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				indent++
			} else {
				break
			}
		}
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	// Remove common indentation
	if minIndent > 0 {
		for i, line := range lines {
			if len(line) >= minIndent {
				lines[i] = line[minIndent:]
			}
		}
	}

	return strings.Join(lines, "\n")
}
