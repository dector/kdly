package parser

// Document represents the top-level KDL document
type Document struct {
	Nodes    []Node
	Comments []Comment
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
type Comment struct {
	Type    CommentType // "line", "multiline", "slashdash"
	Content string      // The comment text
}

// CommentType represents the type of comment
type CommentType string

const (
	CommentTypeLine      CommentType = "line"
	CommentTypeMultiline CommentType = "multiline"
	CommentTypeSlashdash CommentType = "slashdash"
)

// Parser represents a KDL v2 parser
type Parser struct {
	// TODO: Add parser fields
}

// New creates a new Parser instance
func New() *Parser {
	// TODO: Initialize parser
	return &Parser{}
}

// Parse parses a KDL document from the provided string
func (p *Parser) Parse(input string) (*Document, error) {
	// TODO: Implement parsing logic
	return nil, nil
}
