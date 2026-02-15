package parser

import "fmt"

func (p *Parser) validateDisallowedLiteralCodePoints() {
	line := 1
	col := 1

	for i, r := range p.input {
		if isDisallowedLiteralCodePoint(r) && !(i == 0 && r == '\ufeff') {
			panic(fmt.Sprintf("parse error at line %d, col %d: disallowed literal code point U+%04X", line, col, r))
		}

		if r == '\n' {
			line++
			col = 1
		} else {
			col++
		}
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
	p.pendingNodeType = ""
	p.validateDisallowedLiteralCodePoints()
	if len(p.input) > 0 && p.input[0] == '\ufeff' {
		p.advance()
	}

	doc = &Document{
		Nodes: make([]Node, 0),
	}

	var currentNode *Node

	// Main state machine loop
	for p.state != stDocumentEnd {
		switch p.state {
		case stDocumentStart:
			p.skipWhitespaceAndCommentsWithEscline()
			if p.isEOF() {
				// Empty document or no more nodes
				p.state = stDocumentEnd
			} else if p.isSlashdash() {
				// Slashdash comment - skip the entire next node
				p.skipSlashdashNode()
				// Stay in stDocumentStart to process the next node
			} else {
				// Check for optional type annotation before node name
				var nodeTypeAnnotation string
				if p.peek() == '(' {
					// If type annotations are disabled, error out
					if p.disableTypeAnnotations {
						p.panicAt("type annotations are disabled")
					}

					p.advance() // skip '('
					nodeTypeAnnotation = p.parseTypeAnnotation()
					p.skipInlineCommentsOnly()
					if p.isEOF() || p.peek() != ')' {
						p.panicAt("unterminated type annotation: expected ')'")
					}
					p.advance() // skip ')'
					p.skipInlineCommentsOnly()
				}

				// Transition to parsing node name based on next character
				if p.peek() == '"' || p.isRawString() {
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

				// Store the type annotation temporarily - we'll apply it when creating the node
				p.pendingNodeType = nodeTypeAnnotation
			}

		case stNodeName:
			// Parse the node name (identifier)
			nodeName := p.parseIdentifier()

			// Create a node with the parsed name and type annotation (if any)
			node := newNode(nodeName, p.pendingNodeType)

			currentNode = &node
			p.pendingNodeType = ""
			p.nodeBodySawChildBlock = false
			p.nodeBodyPendingSeparator = false

			// Transition to parsing node body
			p.state = stNodeBody

		case stNodeNameQuoted:
			// Parse the node name (quoted string or raw string)
			var nodeName string
			if p.peek() == '"' {
				nodeName = p.parseQuotedString()
			} else {
				nodeName = p.parseRawString()
			}

			// Create a node with the parsed name and type annotation (if any)
			node := newNode(nodeName, p.pendingNodeType)

			currentNode = &node
			p.pendingNodeType = ""
			p.nodeBodySawChildBlock = false
			p.nodeBodyPendingSeparator = false

			// Transition to parsing node body
			p.state = stNodeBody

		case stNodeBody:
			// After node name, check for arguments, properties, children, or end
			// Skip only spaces and tabs, not newlines (newlines terminate nodes)
			bodyStartPos := p.pos
			p.skipInlineWhitespaceAndComments()
			hadSeparator := p.pos != bodyStartPos || p.nodeBodyPendingSeparator
			p.nodeBodyPendingSeparator = false

			if p.isEOF() {
				// End of document - add current node
				if currentNode != nil {
					doc.Nodes = append(doc.Nodes, *currentNode)
					currentNode = nil
				}
				p.state = stDocumentEnd
			} else if p.isSlashdash() {
				if !hadSeparator {
					p.panicAt("expected whitespace before slashdash")
				}
				// Slashdash comment - skip the next argument or property
				if p.skipSlashdashValue() {
					p.nodeBodySawChildBlock = true
				}
				p.nodeBodyPendingSeparator = true
				// Stay in stNodeBody to continue parsing
			} else {
				ch := p.peek()

				if p.nodeBodySawChildBlock && ch != '\n' && ch != '\r' && ch != ';' && ch != '{' {
					p.panicAt("entries are not allowed after a child block")
				}

				if !hadSeparator && ch != '\n' && ch != '\r' && ch != ';' && ch != '{' {
					p.panicAt("expected whitespace before value")
				}

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
					if currentNode != nil {
						p.parseQuotedStringArgumentOrProperty(currentNode, true, true)
					}
				} else if ch == '{' {
					// Children block
					if currentNode != nil {
						if err := p.parseAndAttachChildren(currentNode); err != nil {
							return nil, err
						}
						p.nodeBodySawChildBlock = true

						p.skipInlineWhitespaceAndComments()
						for p.isSlashdash() {
							if !p.skipSlashdashValue() {
								p.panicAt("only child blocks may follow a child block")
							}
							p.skipInlineWhitespaceAndComments()
						}
						if !p.isEOF() && p.peek() != '\n' && p.peek() != '\r' && p.peek() != ';' {
							p.panicAt("expected node terminator after children block")
						}
						if !p.isEOF() && p.peek() == ';' {
							p.advance()
						}
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
					if currentNode != nil {
						p.parseIdentifierArgumentOrProperty(currentNode, true, true)
					}
				} else {
					p.panicAt(fmt.Sprintf("unexpected character in node body: '%c'", ch))
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

		default:
			p.panicAt(fmt.Sprintf("unexpected parser state: %d", p.state))
		}
	}

	return doc, nil
}
