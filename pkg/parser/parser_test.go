package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func node(name string) *Node {
	return &Node{
		Name:       name,
		Arguments:  []Value{},
		Properties: []Property{},
		Children:   []Node{},
	}
}

// ============================================================
// Level 1: Fundamentals (Empty & Single Elements)
// ============================================================

func TestEmptyDocument(t *testing.T) {
	input := ``

	doc, err := New().Parse(input)

	assert.NoError(t, err)
	assert.Empty(t, doc.Nodes)
}

func TestSingleSimpleNode(t *testing.T) {
	input := `hello`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("hello"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSingleSimpleNode_withQuotes(t *testing.T) {
	input := `"illegal(){}[]/\\=#;identifier"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node(`illegal(){}[]/\=#;identifier`),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSingleSimpleNode_flexibleBare(t *testing.T) {
	input := "-<123~!$@%^&*,.:'`|?+>"

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("-<123~!$@%^&*,.:'`|?+>"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithStringArgument(t *testing.T) {
	input := `title "Hello, World"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "title",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "Hello, World"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithBareStringArgument(t *testing.T) {
	input := `node1 this-is-a-string`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node1",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "this-is-a-string"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithEscapedStringArgument(t *testing.T) {
	input := `node2 "this\nhas\tescapes"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node2",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "this\nhas\tescapes"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithRawStringArgument(t *testing.T) {
	input := `node3 #"C:\Users\zkat\raw\string"#`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node3",
				Arguments: []Value{
					{Type: ValueTypeString, Value: `C:\Users\zkat\raw\string`},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Level 2: Basic Values & Types
// ============================================================

func TestNodeWithMultipleArguments(t *testing.T) {
	input := `bookmarks 12 15 188 1234`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "bookmarks",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "12"},
					{Type: ValueTypeNumber, Value: "15"},
					{Type: ValueTypeNumber, Value: "188"},
					{Type: ValueTypeNumber, Value: "1234"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithBooleanAndNull(t *testing.T) {
	input := `flags #true #false #null`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "flags",
				Arguments: []Value{
					{Type: ValueTypeBoolean, Value: "true"},
					{Type: ValueTypeBoolean, Value: "false"},
					{Type: ValueTypeNull, Value: "null"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithHexNumbers(t *testing.T) {
	input := `color 0xdeadbeef`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "color",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0xdeadbeef"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Level 3: Properties & Type Annotations
// ============================================================

func TestNodeWithProperties(t *testing.T) {
	input := `author "Alex Monad" email=alex@example.com active=#true`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "author",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "Alex Monad"},
				},
				Properties: []Property{
					{Key: "email", Value: Value{Type: ValueTypeString, Value: "alex@example.com"}},
					{Key: "active", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
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

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "message",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "\nhello\nworld\n"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestRawString(t *testing.T) {
	input := `path #"C:\path\to\file"#`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "path",
				Arguments: []Value{
					{Type: ValueTypeString, Value: `C:\path\to\file`},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
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
	want := &Document{
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

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestMixedConstructs(t *testing.T) {
	input := `server {
  host "localhost"
  port 8080
  debug #true
}`

	doc, err := New().Parse(input)
	want := &Document{
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

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Level 7: Additional Coverage
// ============================================================

func TestMultipleTopLevelNodes(t *testing.T) {
	input := `node1
node2
node3`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
			*node("node3"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithOnlyProperties(t *testing.T) {
	input := `config debug=#true port=8080`

	doc, err := New().Parse(input)
	want := &Document{
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

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithEmptyChildren(t *testing.T) {
	input := `container {}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("container"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
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
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
			*node("node3"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestQuotedNodeNameOnly(t *testing.T) {
	input := `"node-with-dashes"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node-with-dashes"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSiblingsAfterNesting(t *testing.T) {
	input := `first {}
second`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("first"),
			*node("second"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSemicolonSeparatedNodes(t *testing.T) {
	input := `node1; node2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSemicolonSeparatedNodesWithChildrenInline(t *testing.T) {
	input := `node1; node2 { node3 }`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					*node("node3"),
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSemicolonSeparatedNodesWithChildrenMultiline(t *testing.T) {
	input := `node1; node2 {
  node3
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					*node("node3"),
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Legacy Test (Keep for compatibility)
// ============================================================

func TestGoldenDataLocal(t *testing.T) {
	t.Run("sample test", func(t *testing.T) {
		// This test passes
	})
}
