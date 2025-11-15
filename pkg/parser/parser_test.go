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

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "title" {
		t.Errorf("expected node name 'title', got %q", node.Name)
	}

	if len(node.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(node.Arguments))
	}

	arg := node.Arguments[0]
	if arg.Type != ValueTypeString {
		t.Errorf("expected argument type 'string', got %q", arg.Type)
	}

	if arg.Value != "Hello, World" {
		t.Errorf("expected argument value 'Hello, World', got %q", arg.Value)
	}
}

func TestNodeWithBareStringArgument(t *testing.T) {
	input := `node1 this-is-a-string`

	// Expected: Document with 1 Node
	//   - Node.Name = "node1"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "this-is-a-string"

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "node1" {
		t.Errorf("expected node name 'node1', got %q", node.Name)
	}

	if len(node.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(node.Arguments))
	}

	arg := node.Arguments[0]
	if arg.Type != ValueTypeString {
		t.Errorf("expected argument type 'string', got %q", arg.Type)
	}

	if arg.Value != "this-is-a-string" {
		t.Errorf("expected argument value 'this-is-a-string', got %q", arg.Value)
	}
}

func TestNodeWithEscapedStringArgument(t *testing.T) {
	input := `node2 "this\nhas\tescapes"`

	// Expected: Document with 1 Node
	//   - Node.Name = "node2"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "this\nhas\tescapes" (actual newline and tab)

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "node2" {
		t.Errorf("expected node name 'node2', got %q", node.Name)
	}

	if len(node.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(node.Arguments))
	}

	arg := node.Arguments[0]
	if arg.Type != ValueTypeString {
		t.Errorf("expected argument type 'string', got %q", arg.Type)
	}

	expected := "this\nhas\tescapes"
	if arg.Value != expected {
		t.Errorf("expected argument value %q, got %q", expected, arg.Value)
	}
}

func TestNodeWithRawStringArgument(t *testing.T) {
	input := `node3 #"C:\Users\zkat\raw\string"#`

	// Expected: Document with 1 Node
	//   - Node.Name = "node3"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "C:\Users\zkat\raw\string" (backslashes preserved)

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "node3" {
		t.Errorf("expected node name 'node3', got %q", node.Name)
	}

	if len(node.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(node.Arguments))
	}

	arg := node.Arguments[0]
	if arg.Type != ValueTypeString {
		t.Errorf("expected argument type 'string', got %q", arg.Type)
	}

	expected := `C:\Users\zkat\raw\string`
	if arg.Value != expected {
		t.Errorf("expected argument value %q, got %q", expected, arg.Value)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "bookmarks" {
		t.Errorf("expected node name 'bookmarks', got %q", node.Name)
	}

	if len(node.Arguments) != 4 {
		t.Fatalf("expected 4 arguments, got %d", len(node.Arguments))
	}

	expectedValues := []string{"12", "15", "188", "1234"}
	for i, expected := range expectedValues {
		arg := node.Arguments[i]
		if arg.Type != ValueTypeNumber {
			t.Errorf("expected argument[%d] type 'number', got %q", i, arg.Type)
		}
		if arg.Value != expected {
			t.Errorf("expected argument[%d] value %q, got %q", i, expected, arg.Value)
		}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "flags" {
		t.Errorf("expected node name 'flags', got %q", node.Name)
	}

	if len(node.Arguments) != 3 {
		t.Fatalf("expected 3 arguments, got %d", len(node.Arguments))
	}

	// Check first argument: #true
	arg0 := node.Arguments[0]
	if arg0.Type != ValueTypeBoolean {
		t.Errorf("expected argument[0] type 'boolean', got %q", arg0.Type)
	}
	if arg0.Value != "true" {
		t.Errorf("expected argument[0] value 'true', got %q", arg0.Value)
	}

	// Check second argument: #false
	arg1 := node.Arguments[1]
	if arg1.Type != ValueTypeBoolean {
		t.Errorf("expected argument[1] type 'boolean', got %q", arg1.Type)
	}
	if arg1.Value != "false" {
		t.Errorf("expected argument[1] value 'false', got %q", arg1.Value)
	}

	// Check third argument: #null
	arg2 := node.Arguments[2]
	if arg2.Type != ValueTypeNull {
		t.Errorf("expected argument[2] type 'null', got %q", arg2.Type)
	}
}

func TestNodeWithHexNumbers(t *testing.T) {
	input := `color 0xdeadbeef`

	// Expected: Document with 1 Node
	//   - Node.Name = "color"
	//   - Node.Arguments[0].Type = "number"
	//   - Node.Arguments[0].Value = "0xdeadbeef" (or converted to decimal)

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "color" {
		t.Errorf("expected node name 'color', got %q", node.Name)
	}

	if len(node.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(node.Arguments))
	}

	arg := node.Arguments[0]
	if arg.Type != ValueTypeNumber {
		t.Errorf("expected argument type 'number', got %q", arg.Type)
	}

	if arg.Value != "0xdeadbeef" {
		t.Errorf("expected argument value '0xdeadbeef', got %q", arg.Value)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "author" {
		t.Errorf("expected node name 'author', got %q", node.Name)
	}

	if len(node.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(node.Arguments))
	}

	arg := node.Arguments[0]
	if arg.Type != ValueTypeString {
		t.Errorf("expected argument type 'string', got %q", arg.Type)
	}
	if arg.Value != "Alex Monad" {
		t.Errorf("expected argument value 'Alex Monad', got %q", arg.Value)
	}

	if len(node.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(node.Properties))
	}

	// Check first property: email=alex@example.com
	prop0 := node.Properties[0]
	if prop0.Key != "email" {
		t.Errorf("expected property[0] key 'email', got %q", prop0.Key)
	}
	if prop0.Value.Type != ValueTypeString {
		t.Errorf("expected property[0] value type 'string', got %q", prop0.Value.Type)
	}
	if prop0.Value.Value != "alex@example.com" {
		t.Errorf("expected property[0] value 'alex@example.com', got %q", prop0.Value.Value)
	}

	// Check second property: active=#true
	prop1 := node.Properties[1]
	if prop1.Key != "active" {
		t.Errorf("expected property[1] key 'active', got %q", prop1.Key)
	}
	if prop1.Value.Type != ValueTypeBoolean {
		t.Errorf("expected property[1] value type 'boolean', got %q", prop1.Value.Type)
	}
	if prop1.Value.Value != "true" {
		t.Errorf("expected property[1] value 'true', got %q", prop1.Value.Value)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "message" {
		t.Errorf("expected node name 'message', got %q", node.Name)
	}

	if len(node.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(node.Arguments))
	}

	arg := node.Arguments[0]
	if arg.Type != ValueTypeString {
		t.Errorf("expected argument type 'string', got %q", arg.Type)
	}

	// The common indentation (2 spaces) should be stripped based on the closing quotes' indentation
	expected := "\nhello\nworld\n"
	if arg.Value != expected {
		t.Errorf("expected argument value %q, got %q", expected, arg.Value)
	}
}

func TestRawString(t *testing.T) {
	input := `path #"C:\path\to\file"#`

	// Expected: Document with 1 Node
	//   - Node.Name = "path"
	//   - Node.Arguments[0].Type = "string"
	//   - Node.Arguments[0].Value = "C:\path\to\file" (backslashes preserved literally, no escape processing)

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "path" {
		t.Errorf("expected node name 'path', got %q", node.Name)
	}

	if len(node.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(node.Arguments))
	}

	arg := node.Arguments[0]
	if arg.Type != ValueTypeString {
		t.Errorf("expected argument type 'string', got %q", arg.Type)
	}

	// Raw strings preserve backslashes literally without escape processing
	expected := `C:\path\to\file`
	if arg.Value != expected {
		t.Errorf("expected argument value %q, got %q", expected, arg.Value)
	}
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

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "contents" {
		t.Errorf("expected node name 'contents', got %q", node.Name)
	}

	if len(node.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(node.Children))
	}

	section := node.Children[0]
	if section.Name != "section" {
		t.Errorf("expected child node name 'section', got %q", section.Name)
	}

	if len(section.Arguments) != 1 {
		t.Fatalf("expected 1 argument on section, got %d", len(section.Arguments))
	}

	if section.Arguments[0].Value != "First section" {
		t.Errorf("expected section argument value 'First section', got %q", section.Arguments[0].Value)
	}

	if len(section.Children) != 1 {
		t.Fatalf("expected 1 child on section, got %d", len(section.Children))
	}

	paragraph := section.Children[0]
	if paragraph.Name != "paragraph" {
		t.Errorf("expected child node name 'paragraph', got %q", paragraph.Name)
	}

	if len(paragraph.Arguments) != 1 {
		t.Fatalf("expected 1 argument on paragraph, got %d", len(paragraph.Arguments))
	}

	if paragraph.Arguments[0].Value != "Text" {
		t.Errorf("expected paragraph argument value 'Text', got %q", paragraph.Arguments[0].Value)
	}
}

func TestMixedConstructs(t *testing.T) {
	input := `server {
  host "localhost"
  port 8080
  debug #true
}`

	// Expected: Document with 1 Node
	//   - Node.Name = "server"
	//   - Node.Children[0].Name = "host", Arguments[0].Value = "localhost"
	//   - Node.Children[1].Name = "port", Arguments[0].Value = "8080"
	//   - Node.Children[2].Name = "debug", Arguments[0].Type = "boolean", Value = "true"
	// Note: Ignoring comment line for now as per instructions

	doc, err := New().Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if node.Name != "server" {
		t.Errorf("expected node name 'server', got %q", node.Name)
	}

	if len(node.Children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(node.Children))
	}

	// Check first child: host "localhost"
	host := node.Children[0]
	if host.Name != "host" {
		t.Errorf("expected child[0] name 'host', got %q", host.Name)
	}
	if len(host.Arguments) != 1 {
		t.Fatalf("expected 1 argument on host, got %d", len(host.Arguments))
	}
	if host.Arguments[0].Value != "localhost" {
		t.Errorf("expected host argument value 'localhost', got %q", host.Arguments[0].Value)
	}

	// Check second child: port 8080
	port := node.Children[1]
	if port.Name != "port" {
		t.Errorf("expected child[1] name 'port', got %q", port.Name)
	}
	if len(port.Arguments) != 1 {
		t.Fatalf("expected 1 argument on port, got %d", len(port.Arguments))
	}
	if port.Arguments[0].Value != "8080" {
		t.Errorf("expected port argument value '8080', got %q", port.Arguments[0].Value)
	}

	// Check third child: debug #true
	debug := node.Children[2]
	if debug.Name != "debug" {
		t.Errorf("expected child[2] name 'debug', got %q", debug.Name)
	}
	if len(debug.Arguments) != 1 {
		t.Fatalf("expected 1 argument on debug, got %d", len(debug.Arguments))
	}
	if debug.Arguments[0].Type != ValueTypeBoolean {
		t.Errorf("expected debug argument type 'boolean', got %q", debug.Arguments[0].Type)
	}
	if debug.Arguments[0].Value != "true" {
		t.Errorf("expected debug argument value 'true', got %q", debug.Arguments[0].Value)
	}
}

// ============================================================
// Legacy Test (Keep for compatibility)
// ============================================================

func TestGoldenDataLocal(t *testing.T) {
	t.Run("sample test", func(t *testing.T) {
		// This test passes
	})
}
