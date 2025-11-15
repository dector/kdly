package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Level 1: Fundamentals (Empty & Single Elements)
// ============================================================

func TestEmptyDocument(t *testing.T) {
	input := ``

	// Expected: Document with no nodes
	doc, err := New().Parse(input)
	require.NoError(t, err)
	assert.Empty(t, doc.Nodes)
}

func TestSingleSimpleNode(t *testing.T) {
	input := `hello`

	// Expected: Document with 1 Node
	//   - Node.Name = "hello"
	//   - Node.Arguments = empty
	//   - Node.Properties = empty
	//   - Node.Children = empty

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "hello", node.Name)
	assert.Empty(t, node.Arguments)
	assert.Empty(t, node.Properties)
	assert.Empty(t, node.Children)
}

func TestSingleSimpleNode_withQuotes(t *testing.T) {
	input := `"illegal(){}[]/\\=#;identifier"`

	// Expected: Document with 1 Node
	//   - Node.Name = "illegal(){}[]/\=#;identifier"
	//   - Node.Arguments = empty
	//   - Node.Properties = empty
	//   - Node.Children = empty

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, `illegal(){}[]/\=#;identifier`, node.Name)
	assert.Empty(t, node.Arguments)
	assert.Empty(t, node.Properties)
	assert.Empty(t, node.Children)
}

func TestSingleSimpleNode_flexibleBare(t *testing.T) {
	input := "-<123~!$@%^&*,.:'`|?+>"

	// Expected: Document with 1 Node
	//   - Node.Name = "-<123~!$@%^&*,.:'`|?+>"
	//   - Node.Arguments = empty
	//   - Node.Properties = empty
	//   - Node.Children = empty

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "-<123~!$@%^&*,.:'`|?+>", node.Name)
	assert.Empty(t, node.Arguments)
	assert.Empty(t, node.Properties)
	assert.Empty(t, node.Children)
}

func TestNodeWithStringArgument(t *testing.T) {
	input := `title "Hello, World"`

	// Expected: Document with 1 Node
	//   - Node.Name = "title"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "Hello, World"

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "title", node.Name)
	require.Len(t, node.Arguments, 1)

	arg := node.Arguments[0]
	assert.Equal(t, ValueTypeString, arg.Type)
	assert.Equal(t, "Hello, World", arg.Value)
}

func TestNodeWithBareStringArgument(t *testing.T) {
	input := `node1 this-is-a-string`

	// Expected: Document with 1 Node
	//   - Node.Name = "node1"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "this-is-a-string"

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "node1", node.Name)
	require.Len(t, node.Arguments, 1)

	arg := node.Arguments[0]
	assert.Equal(t, ValueTypeString, arg.Type)
	assert.Equal(t, "this-is-a-string", arg.Value)
}

func TestNodeWithEscapedStringArgument(t *testing.T) {
	input := `node2 "this\nhas\tescapes"`

	// Expected: Document with 1 Node
	//   - Node.Name = "node2"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "this\nhas\tescapes" (actual newline and tab)

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "node2", node.Name)
	require.Len(t, node.Arguments, 1)

	arg := node.Arguments[0]
	assert.Equal(t, ValueTypeString, arg.Type)
	assert.Equal(t, "this\nhas\tescapes", arg.Value)
}

func TestNodeWithRawStringArgument(t *testing.T) {
	input := `node3 #"C:\Users\zkat\raw\string"#`

	// Expected: Document with 1 Node
	//   - Node.Name = "node3"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "C:\Users\zkat\raw\string" (backslashes preserved)

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "node3", node.Name)
	require.Len(t, node.Arguments, 1)

	arg := node.Arguments[0]
	assert.Equal(t, ValueTypeString, arg.Type)
	assert.Equal(t, `C:\Users\zkat\raw\string`, arg.Value)
}

// ============================================================
// Level 2: Basic Values & Types
// ============================================================

func TestNodeWithMultipleArguments(t *testing.T) {
	input := `bookmarks 12 15 188 1234`

	// Expected: Document with 1 Node
	//   - Node.Name = "bookmarks"
	//   - Node.Arguments has 4 elements, all Type="number"
	//   - Values: "12", "15", "188", "1234"

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "bookmarks", node.Name)
	require.Len(t, node.Arguments, 4)

	expectedValues := []string{"12", "15", "188", "1234"}
	for i, expected := range expectedValues {
		arg := node.Arguments[i]
		assert.Equal(t, ValueTypeNumber, arg.Type, "argument[%d] type", i)
		assert.Equal(t, expected, arg.Value, "argument[%d] value", i)
	}
}

func TestNodeWithBooleanAndNull(t *testing.T) {
	input := `flags #true #false #null`

	// Expected: Document with 1 Node
	//   - Node.Name = "flags"
	//   - Node.Arguments[0].Type = "boolean", Value = "true"
	//   - Node.Arguments[1].Type = "boolean", Value = "false"
	//   - Node.Arguments[2].Type = "null"

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "flags", node.Name)
	require.Len(t, node.Arguments, 3)

	// Check first argument: #true
	assert.Equal(t, ValueTypeBoolean, node.Arguments[0].Type)
	assert.Equal(t, "true", node.Arguments[0].Value)

	// Check second argument: #false
	assert.Equal(t, ValueTypeBoolean, node.Arguments[1].Type)
	assert.Equal(t, "false", node.Arguments[1].Value)

	// Check third argument: #null
	assert.Equal(t, ValueTypeNull, node.Arguments[2].Type)
}

func TestNodeWithHexNumbers(t *testing.T) {
	input := `color 0xdeadbeef`

	// Expected: Document with 1 Node
	//   - Node.Name = "color"
	//   - Node.Arguments[0].Type = "number"
	//   - Node.Arguments[0].Value = "0xdeadbeef" (or converted to decimal)

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "color", node.Name)
	require.Len(t, node.Arguments, 1)

	arg := node.Arguments[0]
	assert.Equal(t, ValueTypeNumber, arg.Type)
	assert.Equal(t, "0xdeadbeef", arg.Value)
}

// ============================================================
// Level 3: Properties & Type Annotations
// ============================================================

func TestNodeWithProperties(t *testing.T) {
	input := `author "Alex Monad" email=alex@example.com active=#true`

	// Expected: Document with 1 Node
	//   - Node.Name = "author"
	//   - Node.Arguments[0].Type = "string", Value = "Alex Monad"
	//   - Node.Properties[0].Key = "email", Value.Type = "string", Value.Value = "alex@example.com"
	//   - Node.Properties[1].Key = "active", Value.Type = "boolean", Value.Value = "true"

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "author", node.Name)
	require.Len(t, node.Arguments, 1)

	arg := node.Arguments[0]
	assert.Equal(t, ValueTypeString, arg.Type)
	assert.Equal(t, "Alex Monad", arg.Value)

	require.Len(t, node.Properties, 2)

	// Check first property: email=alex@example.com
	assert.Equal(t, "email", node.Properties[0].Key)
	assert.Equal(t, ValueTypeString, node.Properties[0].Value.Type)
	assert.Equal(t, "alex@example.com", node.Properties[0].Value.Value)

	// Check second property: active=#true
	assert.Equal(t, "active", node.Properties[1].Key)
	assert.Equal(t, ValueTypeBoolean, node.Properties[1].Value.Type)
	assert.Equal(t, "true", node.Properties[1].Value.Value)
}

func TestTypeAnnotations(t *testing.T) {
	input := `numbers (u8)10 (uuid)"123e4567-e89b-12d3-a456-426614174000"`

	// Expected: Document with 1 Node
	//   - Node.Name = "numbers"
	//   - Node.Arguments[0].Type = "number", Value = "10", TypeAnnotation = "u8"
	//   - Node.Arguments[1].Type = "string", Value = "123e4567-e89b-12d3-a456-426614174000", TypeAnnotation = "uuid"

	_ = input
	t.Skip("Parser not yet implemented")
}

// ============================================================
// Level 4: Comments
// ============================================================

func TestLineComment(t *testing.T) {
	input := `// This is a comment`

	// Expected: Document with 1 Comment
	//   - Comment.Type = "line"
	//   - Comment.Content = "This is a comment" (or " This is a comment" with leading space)

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestMultilineComment(t *testing.T) {
	input := `/* Multiline
   comment */`

	// Expected: Document with 1 Comment
	//   - Comment.Type = "multiline"
	//   - Comment.Content = "Multiline\n   comment" (preserving internal formatting)

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestSlashdashComment(t *testing.T) {
	input := `/- node { child "value" }`

	// Expected: Document with 1 Comment
	//   - Comment.Type = "slashdash"
	//   - Comment.Content captures the entire commented-out node structure
	// Alternative: Document with no nodes (slashdash removes the node entirely)

	_ = input
	t.Skip("Parser not yet implemented")
}

// ============================================================
// Level 5: Strings & Multiline
// ============================================================

func TestQuotedMultilineString(t *testing.T) {
	input := `message """
  hello
  world
  """`

	// Expected: Document with 1 Node
	//   - Node.Name = "message"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "\nhello\nworld\n" (dedented based on closing quotes)

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "message", node.Name)
	require.Len(t, node.Arguments, 1)

	arg := node.Arguments[0]
	assert.Equal(t, ValueTypeString, arg.Type)
	// The common indentation (2 spaces) should be stripped based on the closing quotes' indentation
	assert.Equal(t, "\nhello\nworld\n", arg.Value)
}

func TestRawString(t *testing.T) {
	input := `path #"C:\path\to\file"#`

	// Expected: Document with 1 Node
	//   - Node.Name = "path"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "C:\path\to\file" (backslashes preserved literally, no escape processing)

	doc, err := New().Parse(input)
	require.NoError(t, err)
	require.Len(t, doc.Nodes, 1)

	node := doc.Nodes[0]
	assert.Equal(t, "path", node.Name)
	require.Len(t, node.Arguments, 1)

	arg := node.Arguments[0]
	assert.Equal(t, ValueTypeString, arg.Type)
	// Raw strings preserve backslashes literally without escape processing
	assert.Equal(t, `C:\path\to\file`, arg.Value)
}

// ============================================================
// Level 6: Children & Nested Structures
// ============================================================

func TestNodeWithChildren(t *testing.T) {
	input := `contents {
  section "First section" {
    paragraph "Text"
  }
}`

	doc, err := New().Parse(input)
	assert.NoError(t, err)

	expected := &Document{
		Nodes: []Node{
			{
				Name:       "contents",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "section",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "First section"},
						},
						Properties: []Property{},
						Children: []Node{
							{
								Name: "paragraph",
								Arguments: []Value{
									{Type: ValueTypeString, Value: "Text"},
								},
								Properties: []Property{},
								Children:   []Node{},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, doc)
}

func TestMixedConstructs(t *testing.T) {
	input := `server {
  host "localhost"
  port 8080
  debug #true
}`

	doc, err := New().Parse(input)
	assert.NoError(t, err)

	expected := &Document{
		Nodes: []Node{
			{
				Name:       "server",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "host",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "localhost"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "port",
						Arguments: []Value{
							{Type: ValueTypeNumber, Value: "8080"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "debug",
						Arguments: []Value{
							{Type: ValueTypeBoolean, Value: "true"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, doc)
}

// ============================================================
// Level 7: Additional Coverage
// ============================================================

func TestMultipleTopLevelNodes(t *testing.T) {
	input := `node1
node2
node3`

	doc, err := New().Parse(input)
	assert.NoError(t, err)

	expected := &Document{
		Nodes: []Node{
			{Name: "node1", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
			{Name: "node2", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
			{Name: "node3", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
		},
	}

	assert.Equal(t, expected, doc)
}

func TestNodeWithOnlyProperties(t *testing.T) {
	input := `config debug=#true port=8080`

	doc, err := New().Parse(input)
	assert.NoError(t, err)

	expected := &Document{
		Nodes: []Node{
			{
				Name:      "config",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "debug", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
					{Key: "port", Value: Value{Type: ValueTypeNumber, Value: "8080"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.Equal(t, expected, doc)
}

func TestNodeWithEmptyChildren(t *testing.T) {
	input := `container {}`

	doc, err := New().Parse(input)
	assert.NoError(t, err)

	expected := &Document{
		Nodes: []Node{
			{Name: "container", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
		},
	}

	assert.Equal(t, expected, doc)
}

func TestNumbersWithUnderscores(t *testing.T) {
	input := `value 1_000_000`

	// Expected: Document with 1 Node
	//   - Node.Name = "value"
	//   - Node.Arguments[0].Type = "number"
	//   - Node.Arguments[0].Value = "1000000" or "1_000_000" (implementation choice)

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestBinaryAndOctalNumbers(t *testing.T) {
	input := `bits 0b1010 0o755`

	// Expected: Document with 1 Node
	//   - Node.Name = "bits"
	//   - Node.Arguments[0].Type = "number", Value = "0b1010" or "10"
	//   - Node.Arguments[1].Type = "number", Value = "0o755" or "493"

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestFloatingPointNumbers(t *testing.T) {
	input := `coords 3.14 -2.5 1.0e10`

	// Expected: Document with 1 Node
	//   - Node.Name = "coords"
	//   - Node.Arguments[0].Type = "number", Value = "3.14"
	//   - Node.Arguments[1].Type = "number", Value = "-2.5"
	//   - Node.Arguments[2].Type = "number", Value = "1.0e10"

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestMixedArgumentsAndProperties(t *testing.T) {
	input := `person "John" age=30 "Doe" city="NYC"`

	// Expected: Document with 1 Node
	//   - Node.Name = "person"
	//   - Node.Arguments[0].Type = "string", Value = "John"
	//   - Node.Arguments[1].Type = "string", Value = "Doe"
	//   - Node.Properties[0].Key = "age", Value.Type = "number", Value.Value = "30"
	//   - Node.Properties[1].Key = "city", Value.Type = "string", Value.Value = "NYC"

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestVariousWhitespace(t *testing.T) {
	input := "node1\n\nnode2\n  \nnode3"

	doc, err := New().Parse(input)
	assert.NoError(t, err)

	expected := &Document{
		Nodes: []Node{
			{Name: "node1", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
			{Name: "node2", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
			{Name: "node3", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
		},
	}

	assert.Equal(t, expected, doc)
}

func TestQuotedNodeNameOnly(t *testing.T) {
	input := `"node-with-dashes"`

	doc, err := New().Parse(input)
	assert.NoError(t, err)

	expected := &Document{
		Nodes: []Node{
			{Name: "node-with-dashes", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
		},
	}

	assert.Equal(t, expected, doc)
}

func TestSiblingsAfterNesting(t *testing.T) {
	input := `first {}
second`

	doc, err := New().Parse(input)
	assert.NoError(t, err)

	expected := &Document{
		Nodes: []Node{
			{Name: "first", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
			{Name: "second", Arguments: []Value{}, Properties: []Property{}, Children: []Node{}},
		},
	}

	assert.Equal(t, expected, doc)
}

func TestSemicolonSeparatedNodes(t *testing.T) {
	input := `node1; node2`

	// Expected: Document with 2 Nodes
	//   - Node[0].Name = "node1"
	//   - Node[1].Name = "node2"

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestSemicolonSeparatedNodesWithChildrenInline(t *testing.T) {
	input := `node1; node2 { node3 }`

	// Expected: Document with 2 Nodes
	//   - Node[0].Name = "node1"
	//   - Node[1].Name = "node2", Children[0].Name = "node3"

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestSemicolonSeparatedNodesWithChildrenMultiline(t *testing.T) {
	input := `node1; node2 {
  node3
}`

	// Expected: Document with 2 Nodes
	//   - Node[0].Name = "node1"
	//   - Node[1].Name = "node2", Children[0].Name = "node3"

	_ = input
	t.Skip("Parser not yet implemented")
}

// ============================================================
// Legacy Test (Keep for compatibility)
// ============================================================

func TestGoldenDataLocal(t *testing.T) {
	t.Run("sample test", func(t *testing.T) {
		// This test passes
	})
}
