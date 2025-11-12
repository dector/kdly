package kdly

import (
	"testing"
)

func TestParseSimpleNode(t *testing.T) {
	input := `node`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(doc.Nodes))
	}

	if doc.Nodes[0].Name != "node" {
		t.Errorf("Expected node name 'node', got %q", doc.Nodes[0].Name)
	}
}

func TestParseNodeWithArguments(t *testing.T) {
	input := `node "arg1" 42 #true #null`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if len(node.Arguments) != 4 {
		t.Fatalf("Expected 4 arguments, got %d", len(node.Arguments))
	}

	// Check string argument
	if node.Arguments[0].Type != ValueTypeString {
		t.Errorf("Expected string type, got %s", node.Arguments[0].Type)
	}
	if node.Arguments[0].Value != "arg1" {
		t.Errorf("Expected 'arg1', got %v", node.Arguments[0].Value)
	}

	// Check number argument
	if node.Arguments[1].Type != ValueTypeNumber {
		t.Errorf("Expected number type, got %s", node.Arguments[1].Type)
	}
	if node.Arguments[1].Value != int64(42) {
		t.Errorf("Expected 42, got %v", node.Arguments[1].Value)
	}

	// Check boolean argument
	if node.Arguments[2].Type != ValueTypeBoolean {
		t.Errorf("Expected boolean type, got %s", node.Arguments[2].Type)
	}
	if node.Arguments[2].Value != true {
		t.Errorf("Expected true, got %v", node.Arguments[2].Value)
	}

	// Check null argument
	if node.Arguments[3].Type != ValueTypeNull {
		t.Errorf("Expected null type, got %s", node.Arguments[3].Type)
	}
	if node.Arguments[3].Value != nil {
		t.Errorf("Expected nil, got %v", node.Arguments[3].Value)
	}
}

func TestParseNodeWithProperties(t *testing.T) {
	input := `node key1="value1" key2=42`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	node := doc.Nodes[0]
	if len(node.Properties) != 2 {
		t.Fatalf("Expected 2 properties, got %d", len(node.Properties))
	}

	if node.Properties["key1"].Value != "value1" {
		t.Errorf("Expected 'value1', got %v", node.Properties["key1"].Value)
	}

	if node.Properties["key2"].Value != int64(42) {
		t.Errorf("Expected 42, got %v", node.Properties["key2"].Value)
	}
}

func TestParseNodeWithChildren(t *testing.T) {
	input := `parent {
	child1
	child2 "arg"
}`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	parent := doc.Nodes[0]
	if len(parent.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(parent.Children))
	}

	if parent.Children[0].Name != "child1" {
		t.Errorf("Expected 'child1', got %q", parent.Children[0].Name)
	}

	if parent.Children[1].Name != "child2" {
		t.Errorf("Expected 'child2', got %q", parent.Children[1].Name)
	}

	if len(parent.Children[1].Arguments) != 1 {
		t.Errorf("Expected 1 argument for child2, got %d", len(parent.Children[1].Arguments))
	}
}

func TestParseMultipleNodes(t *testing.T) {
	input := `node1
node2
node3`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 3 {
		t.Fatalf("Expected 3 nodes, got %d", len(doc.Nodes))
	}
}

func TestParseNodeWithSemicolon(t *testing.T) {
	input := `node1; node2; node3`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 3 {
		t.Fatalf("Expected 3 nodes, got %d", len(doc.Nodes))
	}
}

func TestParseNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`node 42`, int64(42)},
		{`node -42`, int64(-42)},
		{`node 3.14`, float64(3.14)},
		{`node 1.5e10`, float64(1.5e10)},
		{`node 0xDEAD`, int64(0xDEAD)},
		{`node 0o755`, int64(0755)},
		{`node 0b1010`, int64(0b1010)},
	}

	for _, tt := range tests {
		parser := NewParser(tt.input)
		doc, err := parser.Parse()

		if err != nil {
			t.Fatalf("Parse error for %q: %v", tt.input, err)
		}

		if len(doc.Nodes[0].Arguments) != 1 {
			t.Fatalf("Expected 1 argument for %q", tt.input)
		}

		arg := doc.Nodes[0].Arguments[0]
		if arg.Value != tt.expected {
			t.Errorf("For %q: expected %v, got %v", tt.input, tt.expected, arg.Value)
		}
	}
}

func TestParseTypeAnnotation(t *testing.T) {
	input := `node (type)"value" key=(int)42`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	node := doc.Nodes[0]

	// Check type annotation on argument
	if node.Arguments[0].TypeAnnotation == nil {
		t.Error("Expected type annotation on argument")
	} else if *node.Arguments[0].TypeAnnotation != "type" {
		t.Errorf("Expected 'type', got %q", *node.Arguments[0].TypeAnnotation)
	}

	// Check type annotation on property
	if node.Properties["key"].TypeAnnotation == nil {
		t.Error("Expected type annotation on property")
	} else if *node.Properties["key"].TypeAnnotation != "int" {
		t.Errorf("Expected 'int', got %q", *node.Properties["key"].TypeAnnotation)
	}
}

func TestParseLineComment(t *testing.T) {
	input := `// This is a comment
node`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(doc.Nodes))
	}

	if doc.Nodes[0].Name != "node" {
		t.Errorf("Expected 'node', got %q", doc.Nodes[0].Name)
	}
}

func TestParseBlockComment(t *testing.T) {
	input := `/* This is a
	multiline comment */
node`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(doc.Nodes))
	}
}

func TestParseSlashdashComment(t *testing.T) {
	input := `/-node1
node2
node3 /-"arg1" "arg2"`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// node1 should be skipped
	if len(doc.Nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(doc.Nodes))
	}

	if doc.Nodes[0].Name != "node2" {
		t.Errorf("Expected 'node2', got %q", doc.Nodes[0].Name)
	}

	// node3 should have only 1 argument (arg1 is commented out)
	if len(doc.Nodes[1].Arguments) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(doc.Nodes[1].Arguments))
	}

	if doc.Nodes[1].Arguments[0].Value != "arg2" {
		t.Errorf("Expected 'arg2', got %v", doc.Nodes[1].Arguments[0].Value)
	}
}

func TestParseLineContinuation(t *testing.T) {
	input := `node "arg1" \
	"arg2" \
	"arg3"`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	node := doc.Nodes[0]
	if len(node.Arguments) != 3 {
		t.Fatalf("Expected 3 arguments, got %d", len(node.Arguments))
	}
}

func TestParseMixedArgumentsAndProperties(t *testing.T) {
	input := `node "arg1" key1="value1" "arg2" key2=42`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	node := doc.Nodes[0]
	if len(node.Arguments) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(node.Arguments))
	}

	if len(node.Properties) != 2 {
		t.Fatalf("Expected 2 properties, got %d", len(node.Properties))
	}
}

func TestParseRawString(t *testing.T) {
	input := `node #"C:\path\to\file"#`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	node := doc.Nodes[0]
	if len(node.Arguments) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(node.Arguments))
	}

	expected := `C:\path\to\file`
	if node.Arguments[0].Value != expected {
		t.Errorf("Expected %q, got %q", expected, node.Arguments[0].Value)
	}
}

func TestLexerTokenTypes(t *testing.T) {
	input := `node { } = ; // comment
"string" 123 #true #false #null`

	lexer := NewLexer(input)
	tokens := []Token{}

	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}

	expectedTypes := []TokenType{
		TokenIdentifier,    // node
		TokenLeftBrace,     // {
		TokenRightBrace,    // }
		TokenEquals,        // =
		TokenSemicolon,     // ;
		TokenLineComment,   // // comment
		TokenNewline,       // newline
		TokenString,        // "string"
		TokenNumber,        // 123
		TokenTrue,          // #true
		TokenFalse,         // #false
		TokenNull,          // #null
		TokenEOF,           // EOF
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expected := range expectedTypes {
		if tokens[i].Type != expected {
			t.Errorf("Token %d: expected %s, got %s", i, expected, tokens[i].Type)
		}
	}
}
