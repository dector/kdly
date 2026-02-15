package kdly

import (
	"github.com/dector/kdly/pkg/parser"
	"github.com/dector/kdly/pkg/serializer"
)

// Document represents the top-level KDL document
type Document struct {
	*parser.Document
}

// Node represents a KDL node with name, arguments, properties, and children
type Node = parser.Node

// Value represents any KDL value (argument or property value)
type Value = parser.Value

// ValueType represents the type of a value
type ValueType = parser.ValueType

// NumberBase represents the original base of a parsed numeric literal.
type NumberBase = parser.NumberBase

const (
	ValueTypeString  = parser.ValueTypeString
	ValueTypeNumber  = parser.ValueTypeNumber
	ValueTypeBoolean = parser.ValueTypeBoolean
	ValueTypeNull    = parser.ValueTypeNull

	NumberBaseDecimal     = parser.NumberBaseDecimal
	NumberBaseHexadecimal = parser.NumberBaseHexadecimal
	NumberBaseOctal       = parser.NumberBaseOctal
	NumberBaseBinary      = parser.NumberBaseBinary
)

// Property represents a key-value property on a node
type Property = parser.Property

// Parser represents a KDL v2 parser.
type Parser struct {
	*parser.Parser
}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{Parser: parser.New()}
}

// Parse parses the given input string and returns a Document.
func Parse(input string) (*Document, error) {
	doc, err := parser.New().Parse(input)
	if err != nil {
		return nil, err
	}
	return &Document{Document: doc}, nil
}

// WithAllowDuplicateProperties configures the parser to keep all duplicate properties
// instead of following the KDL v2 spec behavior (rightmost wins).
// When set to true, all properties with duplicate keys will be preserved in the Properties slice.
// When set to false (default), duplicate property keys will follow KDL v2 spec: rightmost value overrides.
func (p *Parser) WithAllowDuplicateProperties() *Parser {
	p.Parser.WithAllowDuplicateProperties(true)
	return p
}

// WithNoTypeAnnotations configures the parser to reject type annotations.
// When set to true, any type annotation (e.g., (type)node or (type)value) will cause a parse error.
// When set to false (default), type annotations are allowed per KDL v2 spec.
// This is useful for enforcing explicit types in configuration files.
func (p *Parser) WithNoTypeAnnotations() *Parser {
	p.Parser.WithDisableTypeAnnotations(true)
	return p
}

// NodesByName returns all nodes with the given name from the document.
func (d *Document) NodesByName(name string) []*Node {
	return d.Document.NodesByName(name)
}

// NodeFirstByName returns the first node with the given name, or nil if not found.
func (d *Document) NodeFirstByName(name string) *Node {
	return d.Document.NodeFirstByName(name)
}

// ToKDL converts a parsed KDL Document back to its text representation.
// This function serializes the document structure into valid KDL v2 format.
func ToKDL(doc *Document) string {
	return serializer.ToKDL(doc.Document)
}
