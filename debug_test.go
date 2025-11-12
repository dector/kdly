package kdly

import (
	"fmt"
	"testing"
)

func TestDebugMultilineRaw(t *testing.T) {
	input := `node #"""
"""#`
	lexer := NewLexer(input)
	
	// Skip "node"
	lexer.NextToken()
	
	// Get the string token
	tok := lexer.NextToken()
	fmt.Printf("Token type: %v\n", tok.Type)
	fmt.Printf("Token literal: %q\n", tok.Literal)
	fmt.Printf("Token literal bytes: %v\n", []byte(tok.Literal))
}
