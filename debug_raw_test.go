package kdly

import (
	"fmt"
	"testing"
)

func TestDebugRawString(t *testing.T) {
	input := `node (type)#"str"#`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	fmt.Printf("Node name: %q\n", node.Name)
	fmt.Printf("Node type annotation: %v\n", node.TypeAnnotation)
	fmt.Printf("Num args: %d\n", len(node.Arguments))

	if len(node.Arguments) > 0 {
		arg := node.Arguments[0]
		fmt.Printf("Arg type: %v\n", arg.Type)
		fmt.Printf("Arg value: %q\n", arg.Value)
		fmt.Printf("Arg type annotation: %v\n", arg.TypeAnnotation)
	}
}
