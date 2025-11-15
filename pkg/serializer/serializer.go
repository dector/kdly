package serializer

import (
	"strings"
	"unicode"

	"github.com/dector/kdly/pkg/parser"
)

// ToKDL converts a parsed KDL Document back to its text representation
func ToKDL(doc *parser.Document) string {
	if doc == nil {
		return ""
	}

	var sb strings.Builder
	for i, node := range doc.Nodes {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(serializeNode(node, 0))
	}
	return sb.String()
}

// serializeNode converts a node to its KDL text representation with proper indentation
func serializeNode(node parser.Node, indent int) string {
	var sb strings.Builder

	// Add indentation
	if indent > 0 {
		sb.WriteString(strings.Repeat("  ", indent))
	}

	// Add type annotation for node if present
	if node.TypeAnnotation != "" {
		sb.WriteString("(")
		sb.WriteString(node.TypeAnnotation)
		sb.WriteString(")")
	}

	// Add node name
	sb.WriteString(formatIdentifier(node.Name))

	// Add arguments
	for _, arg := range node.Arguments {
		sb.WriteString(" ")
		sb.WriteString(serializeValue(arg))
	}

	// Add properties
	for _, prop := range node.Properties {
		sb.WriteString(" ")
		sb.WriteString(serializeProperty(prop))
	}

	// Add children if present
	if len(node.Children) > 0 {
		sb.WriteString(" {\n")
		for _, child := range node.Children {
			sb.WriteString(serializeNode(child, indent+1))
			sb.WriteString("\n")
		}
		if indent > 0 {
			sb.WriteString(strings.Repeat("  ", indent))
		}
		sb.WriteString("}")
	}

	return sb.String()
}

// serializeValue converts a value to its KDL text representation
func serializeValue(val parser.Value) string {
	var sb strings.Builder

	// Add type annotation if present
	if val.TypeAnnotation != "" {
		sb.WriteString("(")
		sb.WriteString(val.TypeAnnotation)
		sb.WriteString(")")
	}

	// Format the value based on its type
	switch val.Type {
	case parser.ValueTypeString:
		sb.WriteString(formatIdentifier(val.Value))
	case parser.ValueTypeNumber:
		// Numbers are stored as-is (including hex, binary, octal, scientific notation)
		sb.WriteString(val.Value)
	case parser.ValueTypeBoolean:
		sb.WriteString("#")
		sb.WriteString(val.Value)
	case parser.ValueTypeNull:
		sb.WriteString("#null")
	default:
		// Fallback: treat as string
		sb.WriteString(formatString(val.Value))
	}

	return sb.String()
}

// serializeProperty converts a property to its KDL text representation (key=value)
func serializeProperty(prop parser.Property) string {
	var sb strings.Builder
	sb.WriteString(formatIdentifier(prop.Key))
	sb.WriteString("=")
	sb.WriteString(serializeValue(prop.Value))
	return sb.String()
}

// formatIdentifier formats a string as either a bare identifier or quoted string
func formatIdentifier(s string) string {
	if canBeBareIdentifier(s) {
		return s
	}
	return formatString(s)
}

// formatString formats a string value with appropriate quoting and escaping
func formatString(s string) string {
	// Check if we should use raw string (contains backslashes but no quotes or special chars)
	if shouldUseRawString(s) {
		return "#\"" + s + "\"#"
	}

	// Use quoted string with escaping
	return "\"" + escapeString(s) + "\""
}

// shouldUseRawString determines if a string should use raw string syntax (#"..."#)
func shouldUseRawString(s string) bool {
	// Use raw strings for paths with backslashes
	hasBackslash := strings.Contains(s, "\\")
	hasDoubleQuote := strings.Contains(s, "\"")

	// Only use raw string if it has backslashes but doesn't have double quotes
	// (since raw strings can't escape quotes)
	return hasBackslash && !hasDoubleQuote
}

// escapeString escapes special characters in a string for quoted string syntax
func escapeString(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		case '\b':
			sb.WriteString("\\b")
		case '\f':
			sb.WriteString("\\f")
		case '\\':
			sb.WriteString("\\\\")
		case '"':
			sb.WriteString("\\\"")
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// canBeBareIdentifier determines if a string can be used as a bare identifier
// Based on KDL v2 spec and parser implementation
func canBeBareIdentifier(s string) bool {
	if s == "" {
		return false
	}

	// Get first rune properly
	firstRune := []rune(s)[0]

	// First character cannot be a digit
	if unicode.IsDigit(firstRune) {
		return false
	}

	// First character cannot be '.' (decimal point for numbers like .5)
	if firstRune == '.' {
		return false
	}

	// Check all characters for forbidden chars
	for _, r := range s {
		if isForbiddenInBareIdentifier(r) {
			return false
		}
	}

	return true
}

// isForbiddenInBareIdentifier checks if a rune is forbidden in bare identifiers
// Matches the parser's logic
func isForbiddenInBareIdentifier(r rune) bool {
	// Whitespace is forbidden
	if unicode.IsSpace(r) {
		return true
	}

	// Reserved syntax characters
	switch r {
	case '[', ']', '{', '}', '(', ')', '\\', '/', '#', '"', ';', '=':
		return true
	}

	return false
}
