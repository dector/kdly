package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// Configurable Type Annotation Behavior Tests
// Tests for the WithDisableTypeAnnotations() configuration
// ============================================================

// ============================================================
// Default Behavior - Type Annotations Allowed (KDL v2 Spec)
// ============================================================

func TestTypeAnnotationsDefault_Allowed(t *testing.T) {
	input := `(contributor)person name="Foo McBar"`

	// Default behavior: type annotations are allowed
	parser := New()
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "person",
				TypeAnnotation: "contributor",
				Arguments:      []Value{},
				Properties: []Property{
					{Key: "name", Value: Value{Type: ValueTypeString, Value: "Foo McBar"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Default parser should allow type annotations per KDL v2 spec")
}

func TestTypeAnnotationsExplicitFalse_Allowed(t *testing.T) {
	input := `node (i32)42`

	// Explicitly set to false: type annotations are allowed
	parser := New().WithDisableTypeAnnotations(false)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "42", TypeAnnotation: "i32"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "WithDisableTypeAnnotations(false) should allow type annotations")
}

func TestTypeAnnotationsDefault_AllowedWithProperties(t *testing.T) {
	input := `node count=(i32)42 id=(uuid)"550e8400-e29b-41d4-a716-446655440000"`

	// Default behavior: type annotations are allowed
	parser := New()
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "",
				Arguments:      []Value{},
				Properties: []Property{
					{Key: "count", Value: Value{Type: ValueTypeNumber, Value: "42", TypeAnnotation: "i32"}},
					{Key: "id", Value: Value{Type: ValueTypeString, Value: "550e8400-e29b-41d4-a716-446655440000", TypeAnnotation: "uuid"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Default parser should allow type annotations on property values")
}

func TestTypeAnnotationsDefault_MixedArgumentsAndPropertiesWithTypes(t *testing.T) {
	input := `node (i32)10 (i32)20 key1=(uuid)"abc" key2=(f64)3.14`

	// Default behavior: type annotations are allowed
	parser := New()
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "10", TypeAnnotation: "i32"},
					{Type: ValueTypeNumber, Value: "20", TypeAnnotation: "i32"},
				},
				Properties: []Property{
					{Key: "key1", Value: Value{Type: ValueTypeString, Value: "abc", TypeAnnotation: "uuid"}},
					{Key: "key2", Value: Value{Type: ValueTypeNumber, Value: "3.14", TypeAnnotation: "f64"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Default parser should allow type annotations on both arguments and properties")
}

// ============================================================
// Disabled Type Annotations - Errors
// ============================================================

func TestTypeAnnotationsDisabled_NodeTypeAnnotation_Error(t *testing.T) {
	input := `(contributor)person name="Foo McBar"`

	// Disable type annotations
	parser := New().WithDisableTypeAnnotations(true)
	_, err := parser.Parse(input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type annotations are disabled")
}

func TestTypeAnnotationsDisabled_ValueTypeAnnotation_Error(t *testing.T) {
	input := `node (i32)42`

	parser := New().WithDisableTypeAnnotations(true)
	_, err := parser.Parse(input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type annotations are disabled")
}

func TestTypeAnnotationsDisabled_PropertyValueTypeAnnotation_Error(t *testing.T) {
	input := `node key=(uuid)"550e8400-e29b-41d4-a716-446655440000"`

	parser := New().WithDisableTypeAnnotations(true)
	_, err := parser.Parse(input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type annotations are disabled")
}

func TestTypeAnnotationsDisabled_PropertyValueTypeAnnotation_Number_Error(t *testing.T) {
	input := `node count=(i32)42 name="test"`

	parser := New().WithDisableTypeAnnotations(true)
	_, err := parser.Parse(input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type annotations are disabled")
}

func TestTypeAnnotationsDisabled_MultiplePropertiesWithTypes_Error(t *testing.T) {
	input := `node id=(uuid)"123" count=(i32)42`

	parser := New().WithDisableTypeAnnotations(true)
	_, err := parser.Parse(input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type annotations are disabled")
}

func TestTypeAnnotationsDisabled_ChildNodeTypeAnnotation_Error(t *testing.T) {
	input := `parent {
  (child-type)child x=1
}`

	parser := New().WithDisableTypeAnnotations(true)
	_, err := parser.Parse(input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type annotations are disabled")
}

func TestTypeAnnotationsDisabled_MultipleAnnotations_FirstError(t *testing.T) {
	input := `(type1)node1 x=1
(type2)node2 x=2`

	parser := New().WithDisableTypeAnnotations(true)
	_, err := parser.Parse(input)

	// Should fail on the first type annotation
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type annotations are disabled")
}

// ============================================================
// Without Type Annotations - Should Work
// ============================================================

func TestTypeAnnotationsDisabled_NoAnnotations_Success(t *testing.T) {
	input := `node "arg1" key="value"`

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg1"},
				},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "value"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Should parse successfully without type annotations")
}

func TestTypeAnnotationsDisabled_PropertiesWithoutTypes_Success(t *testing.T) {
	input := `node count=42 name="test" enabled=#true`

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "",
				Arguments:      []Value{},
				Properties: []Property{
					{Key: "count", Value: Value{Type: ValueTypeNumber, Value: "42"}},
					{Key: "name", Value: Value{Type: ValueTypeString, Value: "test"}},
					{Key: "enabled", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Should parse properties without type annotations successfully")
}

func TestTypeAnnotationsDisabled_MixedArgumentsAndPropertiesNoTypes_Success(t *testing.T) {
	input := `node 10 20 key1="abc" key2=3.14`

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "10"},
					{Type: ValueTypeNumber, Value: "20"},
				},
				Properties: []Property{
					{Key: "key1", Value: Value{Type: ValueTypeString, Value: "abc"}},
					{Key: "key2", Value: Value{Type: ValueTypeNumber, Value: "3.14"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc, "Should parse arguments and properties without type annotations successfully")
}

func TestTypeAnnotationsDisabled_ComplexDocument_NoAnnotations_Success(t *testing.T) {
	input := `node1 a=1 b=2
node2 x=10 {
  child y=100
}
node3 p="first"`

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(doc.Nodes), "Should parse 3 nodes successfully")
	assert.Equal(t, "", doc.Nodes[0].TypeAnnotation, "node1 should have no type annotation")
	assert.Equal(t, "", doc.Nodes[1].TypeAnnotation, "node2 should have no type annotation")
	assert.Equal(t, "", doc.Nodes[2].TypeAnnotation, "node3 should have no type annotation")
	assert.Equal(t, "", doc.Nodes[1].Children[0].TypeAnnotation, "child should have no type annotation")
}

// ============================================================
// Parentheses in Other Contexts - Should Not Be Affected
// ============================================================

func TestTypeAnnotationsDisabled_ParenthesesInStrings_Success(t *testing.T) {
	input := `node key="(not a type)"`

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, "(not a type)", doc.Nodes[0].Properties[0].Value.Value)
}

// ============================================================
// Builder Pattern Chaining
// ============================================================

func TestTypeAnnotationsBuilderPattern_Chaining(t *testing.T) {
	input := `node key="value"`

	// Test that builder pattern returns parser instance for chaining
	parser := New().WithDisableTypeAnnotations(true).WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(doc.Nodes), "Builder pattern should work with multiple options")
}

func TestTypeAnnotationsBuilderPattern_MultipleCalls(t *testing.T) {
	input := `node 42`

	// Test toggling the setting - last one should win
	parser := New().WithDisableTypeAnnotations(true).WithDisableTypeAnnotations(false)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(doc.Nodes), "Last setting should win when called multiple times")
}

func TestTypeAnnotationsBuilderPattern_EnableThenDisable(t *testing.T) {
	input := `(type)node 42`

	// Disable, then enable (should allow type annotations)
	parser := New().WithDisableTypeAnnotations(true).WithDisableTypeAnnotations(false)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, "type", doc.Nodes[0].TypeAnnotation, "Should allow type annotations when re-enabled")
}

// ============================================================
// Combined with Other Features
// ============================================================

func TestTypeAnnotationsDisabled_WithDuplicateProperties(t *testing.T) {
	input := `node a=1 a=2 a=3`

	parser := New().WithDisableTypeAnnotations(true).WithAllowDuplicateProperties(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(doc.Nodes[0].Properties), "Should keep all duplicate properties")
}

func TestTypeAnnotationsDisabled_WithComments(t *testing.T) {
	input := `// This is a comment
node x=1 // inline comment
/* multiline
   comment */
node2 y=2`

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(doc.Nodes), "Should parse nodes with comments")
}

func TestTypeAnnotationsDisabled_WithSlashdash(t *testing.T) {
	input := `/- node1 x=1
node2 x=2`

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(doc.Nodes), "Should skip slashdash commented node")
	assert.Equal(t, "node2", doc.Nodes[0].Name)
}

// ============================================================
// Edge Cases
// ============================================================

func TestTypeAnnotationsDisabled_EmptyDocument(t *testing.T) {
	input := ``

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(doc.Nodes), "Empty document should parse successfully")
}

func TestTypeAnnotationsDisabled_OnlyWhitespace(t *testing.T) {
	input := `

	`

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(doc.Nodes), "Whitespace-only document should parse successfully")
}

func TestTypeAnnotationsDisabled_OnlyComments(t *testing.T) {
	input := `// comment 1
/* comment 2 */
// comment 3`

	parser := New().WithDisableTypeAnnotations(true)
	doc, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(doc.Nodes), "Comment-only document should parse successfully")
}
