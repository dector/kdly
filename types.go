package kdly

// Document represents a complete KDL document containing a list of nodes.
type Document struct {
	Nodes []*Node
}

// Node represents a KDL node with a name, optional arguments, properties, and children.
type Node struct {
	Name       string
	Arguments  []*Value
	Properties map[string]*Value
	Children   []*Node
	// Type annotation for the node itself (e.g., (date)node-name)
	TypeAnnotation *string
}

// Value represents a KDL value which can be a String, Number, Boolean, or Null.
type Value struct {
	Type ValueType
	// Raw string representation
	Raw string
	// Parsed value (one of: string, int64, float64, bool, nil)
	Value interface{}
	// Optional type annotation (e.g., (f32)1.5)
	TypeAnnotation *string
}

// ValueType represents the type of a KDL value.
type ValueType int

const (
	ValueTypeString ValueType = iota
	ValueTypeNumber
	ValueTypeBoolean
	ValueTypeNull
)

// String returns the string representation of a ValueType.
func (vt ValueType) String() string {
	switch vt {
	case ValueTypeString:
		return "String"
	case ValueTypeNumber:
		return "Number"
	case ValueTypeBoolean:
		return "Boolean"
	case ValueTypeNull:
		return "Null"
	default:
		return "Unknown"
	}
}

// NewDocument creates a new empty KDL document.
func NewDocument() *Document {
	return &Document{
		Nodes: make([]*Node, 0),
	}
}

// NewNode creates a new KDL node with the given name.
func NewNode(name string) *Node {
	return &Node{
		Name:       name,
		Arguments:  make([]*Value, 0),
		Properties: make(map[string]*Value),
		Children:   make([]*Node, 0),
	}
}

// NewStringValue creates a new string Value.
func NewStringValue(s string) *Value {
	return &Value{
		Type:  ValueTypeString,
		Raw:   s,
		Value: s,
	}
}

// NewNumberValue creates a new number Value.
// The value parameter should be int64 or float64.
func NewNumberValue(raw string, value interface{}) *Value {
	return &Value{
		Type:  ValueTypeNumber,
		Raw:   raw,
		Value: value,
	}
}

// NewBooleanValue creates a new boolean Value.
func NewBooleanValue(b bool) *Value {
	raw := "#false"
	if b {
		raw = "#true"
	}
	return &Value{
		Type:  ValueTypeBoolean,
		Raw:   raw,
		Value: b,
	}
}

// NewNullValue creates a new null Value.
func NewNullValue() *Value {
	return &Value{
		Type:  ValueTypeNull,
		Raw:   "#null",
		Value: nil,
	}
}

// AddArgument adds an argument to the node.
func (n *Node) AddArgument(value *Value) {
	n.Arguments = append(n.Arguments, value)
}

// AddProperty adds a property to the node.
func (n *Node) AddProperty(key string, value *Value) {
	n.Properties[key] = value
}

// AddChild adds a child node.
func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

// HasChildren returns true if the node has children.
func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}
