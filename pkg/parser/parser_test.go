package parser

import "testing"

// ============================================================
// Level 1: Fundamentals (Empty & Single Elements)
// ============================================================

func TestEmptyDocument(t *testing.T) {
	input := ``

	// Expected: Document with no nodes
	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(doc.Nodes))
	}

	// if len(doc.Comments) != 0 {
	// 	t.Errorf("expected 0 comments, got %d", len(doc.Comments))
	// }
}

func TestSingleSimpleNode(t *testing.T) {
	input := `hello`

	// Expected: Document with 1 Node
	//   - Node.Name = "hello"
	//   - Node.Arguments = empty
	//   - Node.Properties = empty
	//   - Node.Children = empty

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "hello" {
		t.Errorf("expected node name 'hello', got %q", node.Name)
	}

	if len(node.Arguments) != 0 {
		t.Errorf("expected 0 arguments, got %d", len(node.Arguments))
	}

	if len(node.Properties) != 0 {
		t.Errorf("expected 0 properties, got %d", len(node.Properties))
	}

	if len(node.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(node.Children))
	}
}

func TestSingleSimpleNode_withQuotes(t *testing.T) {
	input := `"illegal(){}[]/\\=#;identifier"`

	// Expected: Document with 1 Node
	//   - Node.Name = "illegal(){}[]/\=#;identifier"
	//   - Node.Arguments = empty
	//   - Node.Properties = empty
	//   - Node.Children = empty

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != `illegal(){}[]/\=#;identifier` {
		t.Errorf("expected node name 'illegal(){}[]/\\=#;identifier', got %q", node.Name)
	}

	if len(node.Arguments) != 0 {
		t.Errorf("expected 0 arguments, got %d", len(node.Arguments))
	}

	if len(node.Properties) != 0 {
		t.Errorf("expected 0 properties, got %d", len(node.Properties))
	}

	if len(node.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(node.Children))
	}
}

func TestSingleSimpleNode_flexibleBare(t *testing.T) {
	input := "-<123~!$@%^&*,.:'`|?+>"

	// Expected: Document with 1 Node
	//   - Node.Name = "-<123~!$@%^&*,.:'`|?+>"
	//   - Node.Arguments = empty
	//   - Node.Properties = empty
	//   - Node.Children = empty

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "-<123~!$@%^&*,.:'`|?+>" {
		t.Errorf("expected node name 'illegal(){}[]/\\=#;identifier', got %q", node.Name)
	}

	if len(node.Arguments) != 0 {
		t.Errorf("expected 0 arguments, got %d", len(node.Arguments))
	}

	if len(node.Properties) != 0 {
		t.Errorf("expected 0 properties, got %d", len(node.Properties))
	}

	if len(node.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(node.Children))
	}
}

func TestNodeWithStringArgument(t *testing.T) {
	input := `title "Hello, World"`

	// Expected: Document with 1 Node
	//   - Node.Name = "title"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "Hello, World"

	_ = input
	t.Skip("Parser not yet implemented")
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

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestNodeWithBooleanAndNull(t *testing.T) {
	input := `flags #true #false #null`

	// Expected: Document with 1 Node
	//   - Node.Name = "flags"
	//   - Node.Arguments[0].Type = "boolean", Value = "true"
	//   - Node.Arguments[1].Type = "boolean", Value = "false"
	//   - Node.Arguments[2].Type = "null"

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestNodeWithHexNumbers(t *testing.T) {
	input := `color 0xdeadbeef`

	// Expected: Document with 1 Node
	//   - Node.Name = "color"
	//   - Node.Arguments[0].Type = "number"
	//   - Node.Arguments[0].Value = "0xdeadbeef" (or converted to decimal)

	_ = input
	t.Skip("Parser not yet implemented")
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

	_ = input
	t.Skip("Parser not yet implemented")
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
	//   - Node.Arguments[0].Value = "\n  hello\n  world\n  " (or normalized)

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestRawString(t *testing.T) {
	input := `path #"C:\path\to\file"#`

	// Expected: Document with 1 Node
	//   - Node.Name = "path"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "C:\\path\\to\\file" (backslashes preserved literally)

	_ = input
	t.Skip("Parser not yet implemented")
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

	// Expected: Document with 1 Node
	//   - Node.Name = "contents"
	//   - Node.Children[0].Name = "section"
	//   - Node.Children[0].Arguments[0].Value = "First section"
	//   - Node.Children[0].Children[0].Name = "paragraph"
	//   - Node.Children[0].Children[0].Arguments[0].Value = "Text"

	_ = input
	t.Skip("Parser not yet implemented")
}

func TestMixedConstructs(t *testing.T) {
	input := `server {
  host "localhost"
  port 8080
  // Debug mode
  debug #true
}`

	// Expected: Document with 1 Node
	//   - Node.Name = "server"
	//   - Node.Children[0].Name = "host", Arguments[0].Value = "localhost"
	//   - Node.Children[1].Name = "port", Arguments[0].Value = "8080"
	//   - Node.Children[2] or Comment between: Type = "line", Content = "Debug mode"
	//   - Node.Children[3].Name = "debug", Arguments[0].Type = "boolean", Value = "true"

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
