package kdly

import "fmt"

// TokenType represents the type of a token in the KDL document.
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenError

	// Literals
	TokenIdentifier    // node-name, property-key
	TokenString        // "quoted string", #"raw string"#, """multiline"""
	TokenNumber        // 123, 1.5, 0xDEAD, 0o755, 0b101
	TokenTrue          // #true
	TokenFalse         // #false
	TokenNull          // #null

	// Delimiters
	TokenEquals        // =
	TokenLeftBrace     // {
	TokenRightBrace    // }
	TokenLeftParen     // (
	TokenRightParen    // )
	TokenSemicolon     // ;
	TokenBackslash     // \ (line continuation)
	TokenNewline       // \n
	TokenSlashDash     // /-

	// Comments
	TokenLineComment   // // comment
	TokenBlockComment  // /* comment */
)

// String returns the string representation of a TokenType.
func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return "ERROR"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenString:
		return "STRING"
	case TokenNumber:
		return "NUMBER"
	case TokenTrue:
		return "TRUE"
	case TokenFalse:
		return "FALSE"
	case TokenNull:
		return "NULL"
	case TokenEquals:
		return "="
	case TokenLeftBrace:
		return "{"
	case TokenRightBrace:
		return "}"
	case TokenLeftParen:
		return "("
	case TokenRightParen:
		return ")"
	case TokenSemicolon:
		return ";"
	case TokenBackslash:
		return "\\"
	case TokenNewline:
		return "NEWLINE"
	case TokenSlashDash:
		return "/-"
	case TokenLineComment:
		return "LINE_COMMENT"
	case TokenBlockComment:
		return "BLOCK_COMMENT"
	default:
		return "UNKNOWN"
	}
}

// Position represents a position in the source code.
type Position struct {
	Line   int // 1-based line number
	Column int // 1-based column number
	Offset int // 0-based byte offset
}

// String returns a string representation of the position.
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Token represents a lexical token from the KDL document.
type Token struct {
	Type     TokenType
	Literal  string   // The raw text of the token
	Position Position // Position in the source
}

// NewToken creates a new token with the given type, literal, and position.
func NewToken(tokenType TokenType, literal string, pos Position) Token {
	return Token{
		Type:     tokenType,
		Literal:  literal,
		Position: pos,
	}
}

// String returns a string representation of the token.
func (t Token) String() string {
	if t.Type == TokenNewline {
		return fmt.Sprintf("[%s at %s]", t.Type, t.Position)
	}
	return fmt.Sprintf("[%s %q at %s]", t.Type, t.Literal, t.Position)
}

// IsWhitespace returns true if the token is whitespace (newline).
func (t Token) IsWhitespace() bool {
	return t.Type == TokenNewline
}

// IsComment returns true if the token is a comment.
func (t Token) IsComment() bool {
	return t.Type == TokenLineComment || t.Type == TokenBlockComment
}

// IsValue returns true if the token represents a value (string, number, boolean, null).
func (t Token) IsValue() bool {
	switch t.Type {
	case TokenString, TokenNumber, TokenTrue, TokenFalse, TokenNull:
		return true
	default:
		return false
	}
}

// IsNodeTerminator returns true if the token can terminate a node (newline, semicolon, EOF).
func (t Token) IsNodeTerminator() bool {
	return t.Type == TokenNewline || t.Type == TokenSemicolon || t.Type == TokenEOF
}
