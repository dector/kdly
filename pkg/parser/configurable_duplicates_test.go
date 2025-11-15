package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// Configurable Duplicate Property Behavior Tests
// Tests for the WithAllowDuplicateProperties() configuration
// ============================================================

// ============================================================
// Default Behavior - Spec Compliant (Rightmost Wins)
// ============================================================

func TestConfigurableDefault_SpecCompliant(t *testing.T) {
	input := `node a=1 a=2`

	// Default behavior: spec-compliant (rightmost wins)
	parser := New()
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "2"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Default parser should follow KDL v2 spec (rightmost wins)")
}

func TestConfigurableExplicitFalse_SpecCompliant(t *testing.T) {
	input := `node key="first" key="second" key="third"`

	// Explicitly set to false: spec-compliant (rightmost wins)
	parser := New().WithAllowDuplicateProperties(false)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "third"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "WithAllowDuplicateProperties(false) should follow KDL v2 spec")
}

// ============================================================
// Allow Duplicates - Keep All
// ============================================================

func TestConfigurableAllowDuplicates_KeepAll(t *testing.T) {
	input := `node a=1 a=2`

	// Allow duplicates: keep all
	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "1"}},
					{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "2"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "WithAllowDuplicateProperties(true) should keep all duplicates")
}

func TestConfigurableAllowDuplicates_ThreeDuplicates(t *testing.T) {
	input := `node key="first" key="second" key="third"`

	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "first"}},
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "second"}},
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "third"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Should keep all three duplicates in order")
}

func TestConfigurableAllowDuplicates_DifferentTypes(t *testing.T) {
	input := `node val="string" val=42 val=#true val=#null`

	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "val", Value: Value{Type: ValueTypeString, Value: "string"}},
					{Key: "val", Value: Value{Type: ValueTypeNumber, Value: "42"}},
					{Key: "val", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
					{Key: "val", Value: Value{Type: ValueTypeNull, Value: "null"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Should keep all duplicates with different types")
}

func TestConfigurableAllowDuplicates_MultipleDifferentKeys(t *testing.T) {
	input := `node a=1 b=2 a=3 b=4 c=5 a=6`

	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "1"}},
					{Key: "b", Value: Value{Type: ValueTypeNumber, Value: "2"}},
					{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "3"}},
					{Key: "b", Value: Value{Type: ValueTypeNumber, Value: "4"}},
					{Key: "c", Value: Value{Type: ValueTypeNumber, Value: "5"}},
					{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "6"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Should preserve order of all properties including duplicates")
}

func TestConfigurableAllowDuplicates_WithTypeAnnotations(t *testing.T) {
	input := `node x=(i32)10 x=(f64)20.5 x=(i32)30`

	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "10", TypeAnnotation: "i32"}},
					{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "20.5", TypeAnnotation: "f64"}},
					{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "30", TypeAnnotation: "i32"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Should keep duplicates with type annotations")
}

func TestConfigurableAllowDuplicates_MixedWithArguments(t *testing.T) {
	input := `node "arg1" key=1 "arg2" key=2 "arg3" key=3`

	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg1"},
					{Type: ValueTypeString, Value: "arg2"},
					{Type: ValueTypeString, Value: "arg3"},
				},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "1"}},
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "2"}},
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "3"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Should keep all duplicate properties when mixed with arguments")
}

// ============================================================
// Global Behavior - Applies to Child Nodes
// ============================================================

func TestConfigurableAllowDuplicates_InChildren_Enabled(t *testing.T) {
	input := `parent {
  child x=1 x=2 x=3
}`

	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "parent",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name:      "child",
						Arguments: []Value{},
						Properties: []Property{
							{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "1"}},
							{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "2"}},
							{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "3"}},
						},
						Children: []Node{},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Configuration should apply to child nodes")
}

func TestConfigurableAllowDuplicates_InChildren_Disabled(t *testing.T) {
	input := `parent {
  child x=1 x=2 x=3
}`

	parser := New().WithAllowDuplicateProperties(false)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "parent",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name:      "child",
						Arguments: []Value{},
						Properties: []Property{
							{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "3"}},
						},
						Children: []Node{},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Spec-compliant behavior should apply to child nodes")
}

func TestConfigurableAllowDuplicates_NestedChildren(t *testing.T) {
	input := `grandparent {
  parent a=1 a=2 {
    child b=3 b=4 b=5
  }
}`

	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "grandparent",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name:      "parent",
						Arguments: []Value{},
						Properties: []Property{
							{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "1"}},
							{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "2"}},
						},
						Children: []Node{
							{
								Name:      "child",
								Arguments: []Value{},
								Properties: []Property{
									{Key: "b", Value: Value{Type: ValueTypeNumber, Value: "3"}},
									{Key: "b", Value: Value{Type: ValueTypeNumber, Value: "4"}},
									{Key: "b", Value: Value{Type: ValueTypeNumber, Value: "5"}},
								},
								Children: []Node{},
							},
						},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Configuration should apply to all nested levels")
}

// ============================================================
// Builder Pattern Chaining
// ============================================================

func TestConfigurableBuilderPattern_Chaining(t *testing.T) {
	input := `node key=1 key=2`

	// Test that builder pattern returns parser instance for chaining
	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(doc.Nodes[0].Properties), "Builder pattern should work")
}

func TestConfigurableBuilderPattern_MultipleCalls(t *testing.T) {
	input := `node a=1 a=2`

	// Test toggling the setting
	parser := New().WithAllowDuplicateProperties(true).WithAllowDuplicateProperties(false)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "2"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Last setting should win when called multiple times")
}

// ============================================================
// Complex Scenarios
// ============================================================

func TestConfigurableAllowDuplicates_ComplexDocument(t *testing.T) {
	input := `node1 a=1 b=2 a=3
node2 x=10 x=20 {
  child y=100 y=200
}
node3 p="first" p="second"`

	parser := New().WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(doc.Nodes), "Should have 3 nodes")

	// node1 should have 3 properties (a, b, a)
	assert.Equal(t, 3, len(doc.Nodes[0].Properties), "node1 should have 3 properties")
	assert.Equal(t, "a", doc.Nodes[0].Properties[0].Key)
	assert.Equal(t, "1", doc.Nodes[0].Properties[0].Value.Value)
	assert.Equal(t, "b", doc.Nodes[0].Properties[1].Key)
	assert.Equal(t, "2", doc.Nodes[0].Properties[1].Value.Value)
	assert.Equal(t, "a", doc.Nodes[0].Properties[2].Key)
	assert.Equal(t, "3", doc.Nodes[0].Properties[2].Value.Value)

	// node2 should have 2 properties (x, x)
	assert.Equal(t, 2, len(doc.Nodes[1].Properties), "node2 should have 2 properties")
	assert.Equal(t, "x", doc.Nodes[1].Properties[0].Key)
	assert.Equal(t, "10", doc.Nodes[1].Properties[0].Value.Value)
	assert.Equal(t, "x", doc.Nodes[1].Properties[1].Key)
	assert.Equal(t, "20", doc.Nodes[1].Properties[1].Value.Value)

	// node2's child should have 2 properties (y, y)
	assert.Equal(t, 1, len(doc.Nodes[1].Children), "node2 should have 1 child")
	assert.Equal(t, 2, len(doc.Nodes[1].Children[0].Properties), "child should have 2 properties")

	// node3 should have 2 properties (p, p)
	assert.Equal(t, 2, len(doc.Nodes[2].Properties), "node3 should have 2 properties")
}

func TestConfigurableDefault_ComplexDocument(t *testing.T) {
	input := `node1 a=1 b=2 a=3
node2 x=10 x=20 {
  child y=100 y=200
}`

	parser := New() // Default: spec-compliant
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(doc.Nodes), "Should have 2 nodes")

	// node1 should have 2 unique properties (b=2, a=3 - rightmost a wins)
	assert.Equal(t, 2, len(doc.Nodes[0].Properties), "node1 should have 2 unique properties")
	foundA := false
	foundB := false
	for _, prop := range doc.Nodes[0].Properties {
		if prop.Key == "a" {
			assert.Equal(t, "3", prop.Value.Value, "Rightmost 'a' should win")
			foundA = true
		}
		if prop.Key == "b" {
			assert.Equal(t, "2", prop.Value.Value)
			foundB = true
		}
	}
	assert.True(t, foundA && foundB, "Should have both 'a' and 'b' properties")

	// node2 should have 1 property (x=20 - rightmost wins)
	assert.Equal(t, 1, len(doc.Nodes[1].Properties), "node2 should have 1 property")
	assert.Equal(t, "x", doc.Nodes[1].Properties[0].Key)
	assert.Equal(t, "20", doc.Nodes[1].Properties[0].Value.Value)

	// node2's child should have 1 property (y=200 - rightmost wins)
	assert.Equal(t, 1, len(doc.Nodes[1].Children[0].Properties), "child should have 1 property")
	assert.Equal(t, "y", doc.Nodes[1].Children[0].Properties[0].Key)
	assert.Equal(t, "200", doc.Nodes[1].Children[0].Properties[0].Value.Value)
}
