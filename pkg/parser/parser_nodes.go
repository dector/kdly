package parser

import "fmt"

func newNode(name, typeAnnotation string) Node {
	return Node{
		Name:           name,
		TypeAnnotation: typeAnnotation,
		Arguments:      make([]Value, 0),
		Properties:     make([]Property, 0),
		Children:       make([]Node, 0),
	}
}

func (p *Parser) parseAndAttachChildren(node *Node) error {
	p.advance() // skip '{'

	children, err := p.parseChildren()
	if err != nil {
		return err
	}

	node.Children = children
	return nil
}

func (p *Parser) parseIdentifierArgumentOrProperty(node *Node, allowCommentsAfterIdentifier bool, allowSlashdashValue bool) {
	savedPos := p.pos
	savedLine := p.line
	savedCol := p.col

	identifier := p.parseIdentifier()

	if allowCommentsAfterIdentifier {
		p.skipInlineWhitespaceAndComments()
	} else {
		for !p.isEOF() && (p.peek() == ' ' || p.peek() == '\t') {
			p.advance()
		}
	}

	if !p.isEOF() && p.peek() == '=' {
		p.advance() // skip '='
		p.skipInlineWhitespaceAndComments()

		var propValue Value
		if allowSlashdashValue && p.isSlashdash() {
			p.skipSlashdashValue()
			propValue = Value{Type: ValueTypeString, Value: ""}
		} else {
			propValue = p.parseValueWithOptionalTypeAnnotation()
		}

		node.addOrUpdateProperty(identifier, propValue, p.allowDuplicateProperties)
		return
	}

	p.pos = savedPos
	p.line = savedLine
	p.col = savedCol

	strValue := p.parseIdentifier()
	node.Arguments = append(node.Arguments, Value{
		Type:  ValueTypeString,
		Value: strValue,
	})
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

	// Skip optional type annotation before node name
	if p.peek() == '(' {
		p.advance() // skip '('
		p.parseTypeAnnotation()
		p.skipInlineWhitespaceAndComments()
		if !p.isEOF() && p.peek() == ')' {
			p.advance() // skip ')'
		}
		p.skipInlineWhitespaceAndComments()
	}

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

		// Check for slashdash comment - skip the entire next node
		if p.isSlashdash() {
			p.skipSlashdashNode()
			continue // Skip to next iteration to parse the next node
		}

		// Parse a child node
		var childNode Node

		// Check for optional type annotation before node name
		var nodeTypeAnnotation string
		if p.peek() == '(' {
			// If type annotations are disabled, error out
			if p.disableTypeAnnotations {
				return nil, fmt.Errorf("parse error at line %d, col %d: type annotations are disabled", p.line, p.col)
			}

			p.advance() // skip '('
			nodeTypeAnnotation = p.parseTypeAnnotation()
			p.skipInlineWhitespaceAndComments()
			if p.isEOF() || p.peek() != ')' {
				return nil, fmt.Errorf("parse error at line %d, col %d: unterminated type annotation: expected ')'", p.line, p.col)
			}
			p.advance() // skip ')'
			p.skipInlineWhitespaceAndComments()
		}

		// Check if node name is quoted or bare
		if p.peek() == '"' {
			nodeName := p.parseQuotedString()
			childNode = newNode(nodeName, nodeTypeAnnotation)
		} else if isIdentifierStart(p.peek()) {
			nodeName := p.parseIdentifier()
			childNode = newNode(nodeName, nodeTypeAnnotation)
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
				if err := p.parseAndAttachChildren(&childNode); err != nil {
					return nil, err
				}
				break
			}

			// Check for arguments and properties
			if ch == '(' || ch == '"' || ch == '#' || p.looksLikeNumber() {
				// This looks like an argument value (possibly with type annotation)
				argValue := p.parseValueWithOptionalTypeAnnotation()
				childNode.Arguments = append(childNode.Arguments, argValue)
			} else if isIdentifierStart(ch) {
				p.parseIdentifierArgumentOrProperty(&childNode, false, false)
			} else {
				// Unknown character - might be end of node
				break
			}
		}

		children = append(children, childNode)
	}
}
