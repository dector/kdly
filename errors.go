package kdly

import "fmt"

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Position Position
	Message  string
}

// Error implements the error interface.
func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at %s: %s", e.Position, e.Message)
}

// NewParseError creates a new ParseError.
func NewParseError(pos Position, message string) *ParseError {
	return &ParseError{
		Position: pos,
		Message:  message,
	}
}

// LexError represents an error that occurred during lexing.
type LexError struct {
	Position Position
	Message  string
}

// Error implements the error interface.
func (e *LexError) Error() string {
	return fmt.Sprintf("lex error at %s: %s", e.Position, e.Message)
}

// NewLexError creates a new LexError.
func NewLexError(pos Position, message string) *LexError {
	return &LexError{
		Position: pos,
		Message:  message,
	}
}

// ErrorList is a list of errors accumulated during parsing.
type ErrorList []error

// Error implements the error interface.
func (el ErrorList) Error() string {
	if len(el) == 0 {
		return "no errors"
	}
	if len(el) == 1 {
		return el[0].Error()
	}
	result := fmt.Sprintf("%d errors occurred:\n", len(el))
	for i, err := range el {
		result += fmt.Sprintf("  %d. %s\n", i+1, err.Error())
	}
	return result
}

// Add adds an error to the list.
func (el *ErrorList) Add(err error) {
	*el = append(*el, err)
}

// HasErrors returns true if the list contains any errors.
func (el ErrorList) HasErrors() bool {
	return len(el) > 0
}
