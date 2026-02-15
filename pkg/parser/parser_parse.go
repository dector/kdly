package parser

import "fmt"

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
				// Check for optional type annotation before node name
				var nodeTypeAnnotation string
				if p.peek() == '(' {
					// If type annotations are disabled, error out
					if p.disableTypeAnnotations {
						p.panicAt("type annotations are disabled")
					}

					p.advance() // skip '('
					nodeTypeAnnotation = p.parseTypeAnnotation()
					p.skipInlineWhitespaceAndComments()
					if p.isEOF() || p.peek() != ')' {
						p.panicAt("unterminated type annotation: expected ')'")
					}
					p.advance() // skip ')'
					p.skipInlineWhitespaceAndComments()
				}

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

			// Transition to parsing node body
			p.state = stNodeBody

		case stNodeNameQuoted:
			// Parse the node name (quoted string)
			nodeName := p.parseQuotedString()

			// Create a node with the parsed name and type annotation (if any)
			node := newNode(nodeName, p.pendingNodeType)

			currentNode = &node
			p.pendingNodeType = ""

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
					if currentNode != nil {
						if err := p.parseAndAttachChildren(currentNode); err != nil {
							return nil, err
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

		default:
			p.panicAt(fmt.Sprintf("unexpected parser state: %d", p.state))
		}
	}

	return doc, nil
}
