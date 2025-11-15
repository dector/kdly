package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// Duplication Tests
// Tests for handling duplicate properties, node names, and arguments
// according to KDL v2 specification
// ============================================================

// ============================================================
// Duplicate Properties - Rightmost Override
// ============================================================

func TestDuplicateProperties_Simple(t *testing.T) {
	input := `node a=1 a=2`

	doc, err := New().Parse(input)
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
	assert.Equal(t, want, doc, "Rightmost property value should override")
}

func TestDuplicateProperties_ThreeDuplicates(t *testing.T) {
	input := `node key="first" key="second" key="third"`

	doc, err := New().Parse(input)
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
	assert.Equal(t, want, doc, "Rightmost property value should override (3 duplicates)")
}

func TestDuplicateProperties_DifferentTypes(t *testing.T) {
	input := `node val="string" val=42 val=#true`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "val", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Rightmost property value should override (different types)")
}

func TestDuplicateProperties_WithTypeAnnotations(t *testing.T) {
	input := `node x=(i32)10 x=(f64)20.5`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "20.5", TypeAnnotation: "f64"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Rightmost property with type annotation should override")
}

func TestDuplicateProperties_MixedWithArguments(t *testing.T) {
	input := `node "arg1" key=1 "arg2" key=2 "arg3" key=3`

	doc, err := New().Parse(input)
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
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "3"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Rightmost property should override when mixed with arguments")
}

func TestDuplicateProperties_MultipleDifferentKeys(t *testing.T) {
	input := `node a=1 b=2 a=3 b=4 c=5`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "3"}},
					{Key: "b", Value: Value{Type: ValueTypeNumber, Value: "4"}},
					{Key: "c", Value: Value{Type: ValueTypeNumber, Value: "5"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Each property key should have its rightmost value")
}

func TestDuplicateProperties_WithLineContinuation(t *testing.T) {
	input := `node \
  key=1 \
  key=2 \
  key=3`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "3"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Rightmost property should override across line continuations")
}

func TestDuplicateProperties_InChildren(t *testing.T) {
	input := `parent {
  child x=1 x=2
}`

	doc, err := New().Parse(input)
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
							{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "2"}},
						},
						Children: []Node{},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Rightmost property should override in child nodes")
}

func TestDuplicateProperties_NullOverridesValue(t *testing.T) {
	input := `node key="value" key=#null`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNull, Value: "null"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Rightmost null value should override previous value")
}

func TestDuplicateProperties_SpecialNumbers(t *testing.T) {
	input := `node val=1.5 val=#inf val=#-inf val=#nan`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "val", Value: Value{Type: ValueTypeNumber, Value: "nan"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Rightmost special number should override")
}

func TestDuplicateProperties_BareIdentifiers(t *testing.T) {
	input := `node key=first-value key=second-value key=third-value`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "third-value"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Rightmost bare identifier value should override")
}

func TestDuplicateProperties_RawStrings(t *testing.T) {
	input := `node path=#"C:\first"# path=#"C:\second"#`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "path", Value: Value{Type: ValueTypeString, Value: `C:\second`}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Rightmost raw string value should override")
}

// ============================================================
// Duplicate Node Names - Should Be Allowed
// ============================================================

func TestDuplicateNodeNames_TwoNodes(t *testing.T) {
	input := `node
node`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node"),
			*node("node"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple nodes with same name should be allowed")
}

func TestDuplicateNodeNames_ThreeNodesWithDifferentValues(t *testing.T) {
	input := `item "first"
item "second"
item "third"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "item",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "first"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name: "item",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "second"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name: "item",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "third"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple nodes with same name but different arguments should be allowed")
}

func TestDuplicateNodeNames_WithProperties(t *testing.T) {
	input := `config debug=#true
config port=8080
config host="localhost"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "config",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "debug", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
				},
				Children: []Node{},
			},
			{
				Name:      "config",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "port", Value: Value{Type: ValueTypeNumber, Value: "8080"}},
				},
				Children: []Node{},
			},
			{
				Name:      "config",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "host", Value: Value{Type: ValueTypeString, Value: "localhost"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple nodes with same name and different properties should be allowed")
}

func TestDuplicateNodeNames_InChildren(t *testing.T) {
	input := `parent {
  child "first"
  child "second"
  child "third"
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "parent",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "child",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "first"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "child",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "second"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "child",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "third"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple child nodes with same name should be allowed")
}

func TestDuplicateNodeNames_MixedWithUniqueNodes(t *testing.T) {
	input := `node1
node2
node1
node3
node2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
			*node("node1"),
			*node("node3"),
			*node("node2"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Duplicate node names mixed with unique names should preserve order")
}

// ============================================================
// Duplicate Arguments - Should Be Allowed
// ============================================================

func TestDuplicateArguments_SameValue(t *testing.T) {
	input := `node "value" "value" "value"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "value"},
					{Type: ValueTypeString, Value: "value"},
					{Type: ValueTypeString, Value: "value"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple arguments with same value should be allowed")
}

func TestDuplicateArguments_SameNumber(t *testing.T) {
	input := `node 42 42 42`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "42"},
					{Type: ValueTypeNumber, Value: "42"},
					{Type: ValueTypeNumber, Value: "42"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple arguments with same number should be allowed")
}

func TestDuplicateArguments_SameBareIdentifier(t *testing.T) {
	input := `node foo foo foo`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "foo"},
					{Type: ValueTypeString, Value: "foo"},
					{Type: ValueTypeString, Value: "foo"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple arguments with same bare identifier should be allowed")
}

func TestDuplicateArguments_SameBoolean(t *testing.T) {
	input := `node #true #true #false #false`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeBoolean, Value: "true"},
					{Type: ValueTypeBoolean, Value: "true"},
					{Type: ValueTypeBoolean, Value: "false"},
					{Type: ValueTypeBoolean, Value: "false"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple arguments with same boolean values should be allowed")
}

func TestDuplicateArguments_SameNull(t *testing.T) {
	input := `node #null #null #null`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNull, Value: "null"},
					{Type: ValueTypeNull, Value: "null"},
					{Type: ValueTypeNull, Value: "null"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple null arguments should be allowed")
}

func TestDuplicateArguments_WithTypeAnnotations(t *testing.T) {
	input := `node (i32)10 (i32)10 (i32)10`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "10", TypeAnnotation: "i32"},
					{Type: ValueTypeNumber, Value: "10", TypeAnnotation: "i32"},
					{Type: ValueTypeNumber, Value: "10", TypeAnnotation: "i32"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Multiple arguments with same value and type annotation should be allowed")
}

func TestDuplicateArguments_MixedWithProperties(t *testing.T) {
	input := `node "arg" key=1 "arg" key=2 "arg"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
					{Type: ValueTypeString, Value: "arg"},
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "2"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Duplicate arguments should be preserved while duplicate properties are overridden")
}

// ============================================================
// Complex Duplication Scenarios
// ============================================================

func TestComplexDuplication_AllTypes(t *testing.T) {
	input := `node1 "arg1" a=1 "arg2" a=2
node1 "arg3" b=3
node2 x=10 x=20 x=30`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node1",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg1"},
					{Type: ValueTypeString, Value: "arg2"},
				},
				Properties: []Property{
					{Key: "a", Value: Value{Type: ValueTypeNumber, Value: "2"}},
				},
				Children: []Node{},
			},
			{
				Name: "node1",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg3"},
				},
				Properties: []Property{
					{Key: "b", Value: Value{Type: ValueTypeNumber, Value: "3"}},
				},
				Children: []Node{},
			},
			{
				Name:      "node2",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "x", Value: Value{Type: ValueTypeNumber, Value: "30"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Complex scenario with duplicate nodes, arguments, and properties")
}

func TestComplexDuplication_NestedChildren(t *testing.T) {
	input := `parent {
  child key=1 key=2
  child key=3 key=4
}`

	doc, err := New().Parse(input)
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
							{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "2"}},
						},
						Children: []Node{},
					},
					{
						Name:      "child",
						Arguments: []Value{},
						Properties: []Property{
							{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "4"}},
						},
						Children: []Node{},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Duplicate child node names with overridden properties")
}

func TestComplexDuplication_WithSlashdash(t *testing.T) {
	input := `node key=1 /-key=999 key=2 key=3`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "3"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Slashdash property should be ignored, rightmost non-slashdash should win")
}
