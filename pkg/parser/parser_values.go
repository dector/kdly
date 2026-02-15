package parser

// parseValueWithOptionalTypeAnnotation parses an optional type annotation followed by a value
// Returns the parsed Value with TypeAnnotation field set if present
func (p *Parser) parseValueWithOptionalTypeAnnotation() Value {
	var typeAnnotation string

	// Check for type annotation: (type)
	if p.peek() == '(' {
		// If type annotations are disabled, error out
		if p.disableTypeAnnotations {
			p.panicAt("type annotations are disabled")
		}

		p.advance() // skip '('

		// Parse the type annotation token
		typeAnnotation = p.parseTypeAnnotation()
		p.skipInlineCommentsOnly()

		// Expect closing ')'
		if p.isEOF() || p.peek() != ')' {
			p.panicAt("expected ')' after type annotation")
		}
		p.advance() // skip ')'
		p.skipInlineCommentsOnly()
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
	} else if p.isSignedDotBareIdentifier() {
		// +. and -. without trailing digit are bare identifiers
		strValue := p.parseIdentifier()
		value = Value{
			Type:           ValueTypeString,
			Value:          strValue,
			TypeAnnotation: typeAnnotation,
		}
	} else if p.looksLikeNumber() {
		// Numeric literal
		numValue := p.parseNumber()
		p.validateNumberBoundary()
		numValue, originalBase := normalizeNumberValue(numValue)
		value = Value{
			Type:           ValueTypeNumber,
			Value:          numValue,
			TypeAnnotation: typeAnnotation,
			OriginalBase:   originalBase,
		}
	} else if isIdentifierStart(ch) {
		// Bare identifier string
		strValue := p.parseIdentifier()
		if isSpecialFloatIdentifier(strValue) {
			p.panicAt("special floating-point values must be written as #inf, #-inf, or #nan")
		}
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

func isSpecialFloatIdentifier(s string) bool {
	return s == "inf" || s == "-inf" || s == "nan"
}
