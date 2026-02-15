package parser

// Document represents the top-level KDL document
type Document struct {
	Nodes []Node
	// Comments []Comment
}

// Node represents a KDL node with name, arguments, properties, and children
type Node struct {
	Name           string
	TypeAnnotation string // Optional type annotation like "author", "contributor"
	Arguments      []Value
	Properties     []Property
	Children       []Node
}

// Value represents any KDL value (argument or property value)
type Value struct {
	Type           ValueType  // "string", "number", "boolean", "null"
	Value          string     // String value; numeric literals are normalized to decimal text
	TypeAnnotation string     // Optional type annotation like "u8", "uuid"
	OriginalBase   NumberBase // Original numeric base for number literals
}

// ValueType represents the type of a value
type ValueType string

const (
	ValueTypeString  ValueType = "string"
	ValueTypeNumber  ValueType = "number"
	ValueTypeBoolean ValueType = "boolean"
	ValueTypeNull    ValueType = "null"
)

// NumberBase represents the original base of a parsed numeric literal.
type NumberBase int

const (
	NumberBaseDecimal NumberBase = iota
	NumberBaseHexadecimal
	NumberBaseOctal
	NumberBaseBinary
)

// Property represents a key-value property on a node
type Property struct {
	Key   string
	Value Value
}

// addOrUpdateProperty adds a property to a node.
// If allowDuplicates is true, all properties are kept (even with duplicate keys).
// If allowDuplicates is false, follows KDL v2 spec: "rightmost values override duplicates".
func (n *Node) addOrUpdateProperty(key string, value Value, allowDuplicates bool) {
	if allowDuplicates {
		// Keep all duplicates - just append
		n.Properties = append(n.Properties, Property{Key: key, Value: value})
		return
	}

	// Spec-compliant behavior: look for an existing property with the same key
	for i := range n.Properties {
		if n.Properties[i].Key == key {
			// Found duplicate - update the value (rightmost wins)
			n.Properties[i].Value = value
			return
		}
	}
	// No duplicate found - append new property
	n.Properties = append(n.Properties, Property{Key: key, Value: value})
}

func (d *Document) NodesByName(name string) []*Node {
	var nodes []*Node
	for _, node := range d.Nodes {
		if node.Name == name {
			nodes = append(nodes, &node)
		}
	}
	return nodes
}

func (d *Document) NodeFirstByName(name string) *Node {
	for _, node := range d.Nodes {
		if node.Name == name {
			return &node
		}
	}
	return nil
}
