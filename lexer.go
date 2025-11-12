package kdly

import (
	"fmt"
	"strings"
	"unicode"
)

// Lexer tokenizes a KDL document.
type Lexer struct {
	input    string
	position int  // current position in input
	line     int  // current line (1-based)
	column   int  // current column (1-based)
	ch       rune // current character
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar() // Initialize first character
	return l
}

// readChar reads the next character and advances position.
func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = rune(l.input[l.position])
	}
	l.position++
	l.column++
}

// peekChar returns the next character without advancing.
func (l *Lexer) peekChar() rune {
	if l.position >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position])
}

// peekCharAt returns the character at offset positions ahead without advancing.
func (l *Lexer) peekCharAt(offset int) rune {
	pos := l.position + offset - 1
	if pos >= len(l.input) {
		return 0
	}
	return rune(l.input[pos])
}

// currentPosition returns the current position in the source.
func (l *Lexer) currentPosition() Position {
	return Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.position - 1,
	}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	pos := l.currentPosition()

	switch l.ch {
	case 0:
		return NewToken(TokenEOF, "", pos)
	case '\n':
		l.readChar()
		tok := NewToken(TokenNewline, "\n", pos)
		l.line++
		l.column = 0
		return tok
	case '\r':
		// Handle \r\n as a single newline
		l.readChar()
		if l.ch == '\n' {
			l.readChar()
		}
		tok := NewToken(TokenNewline, "\n", pos)
		l.line++
		l.column = 0
		return tok
	case '=':
		l.readChar()
		return NewToken(TokenEquals, "=", pos)
	case '{':
		l.readChar()
		return NewToken(TokenLeftBrace, "{", pos)
	case '}':
		l.readChar()
		return NewToken(TokenRightBrace, "}", pos)
	case '(':
		l.readChar()
		return NewToken(TokenLeftParen, "(", pos)
	case ')':
		l.readChar()
		return NewToken(TokenRightParen, ")", pos)
	case ';':
		l.readChar()
		return NewToken(TokenSemicolon, ";", pos)
	case '\\':
		l.readChar()
		return NewToken(TokenBackslash, "\\", pos)
	case '/':
		// Could be //, /*, or /-
		next := l.peekChar()
		if next == '/' {
			return l.readLineComment(pos)
		} else if next == '*' {
			return l.readBlockComment(pos)
		} else if next == '-' {
			l.readChar() // consume '/'
			l.readChar() // consume '-'
			return NewToken(TokenSlashDash, "/-", pos)
		}
		// Otherwise it's an identifier starting with /
		return l.readIdentifier(pos)
	case '"':
		// Check for multiline string (""")
		if l.peekChar() == '"' && l.peekCharAt(2) == '"' {
			return l.readMultilineString(pos)
		}
		return l.readQuotedString(pos)
	case '#':
		// Could be #true, #false, #null, or raw string #"..."#
		next := l.peekChar()
		if next == '"' {
			return l.readRawString(pos)
		}
		// Check for keywords
		return l.readHashKeyword(pos)
	default:
		// Number or identifier
		if isDigit(l.ch) || (l.ch == '-' || l.ch == '+') && isDigit(l.peekChar()) {
			return l.readNumber(pos)
		}
		if isIdentifierStart(l.ch) {
			return l.readIdentifier(pos)
		}
		// Unknown character
		ch := l.ch
		l.readChar()
		return NewToken(TokenError, fmt.Sprintf("unexpected character: %q", ch), pos)
	}
}

// skipWhitespace skips whitespace characters except newlines.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\u00A0' || l.ch == '\uFEFF' {
		l.readChar()
	}
}

// readLineComment reads a line comment (// ...).
func (l *Lexer) readLineComment(pos Position) Token {
	start := l.position - 1
	l.readChar() // consume first /
	l.readChar() // consume second /

	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}

	return NewToken(TokenLineComment, l.input[start:l.position-1], pos)
}

// readBlockComment reads a block comment (/* ... */) with nesting support.
func (l *Lexer) readBlockComment(pos Position) Token {
	start := l.position - 1
	l.readChar() // consume /
	l.readChar() // consume *

	depth := 1
	for depth > 0 && l.ch != 0 {
		if l.ch == '/' && l.peekChar() == '*' {
			l.readChar()
			l.readChar()
			depth++
		} else if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			depth--
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
		}
	}

	return NewToken(TokenBlockComment, l.input[start:l.position-1], pos)
}

// readQuotedString reads a quoted string ("...").
func (l *Lexer) readQuotedString(pos Position) Token {
	start := l.position - 1
	l.readChar() // consume opening "

	var result strings.Builder
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			// Handle escape sequences
			switch l.ch {
			case 'n':
				result.WriteRune('\n')
			case 't':
				result.WriteRune('\t')
			case 'r':
				result.WriteRune('\r')
			case '\\':
				result.WriteRune('\\')
			case '"':
				result.WriteRune('"')
			case '/':
				result.WriteRune('/')
			case 'b':
				result.WriteRune('\b')
			case 'f':
				result.WriteRune('\f')
			default:
				// For now, include the escape sequence as-is
				result.WriteRune('\\')
				result.WriteRune(l.ch)
			}
			l.readChar()
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			result.WriteRune(l.ch)
			l.readChar()
		}
	}

	if l.ch == '"' {
		l.readChar() // consume closing "
	}

	return NewToken(TokenString, l.input[start:l.position-1], pos)
}

// readRawString reads a raw string (#"..."#, ##"..."##, etc.).
func (l *Lexer) readRawString(pos Position) Token {
	start := l.position - 1
	l.readChar() // consume #

	// Count hashes
	hashCount := 1
	for l.ch == '#' {
		hashCount++
		l.readChar()
	}

	if l.ch != '"' {
		return NewToken(TokenError, "expected '\"' after '#' in raw string", pos)
	}
	l.readChar() // consume opening "

	// Read until we find closing "###...
	for l.ch != 0 {
		if l.ch == '"' {
			// Check if followed by the right number of hashes
			matched := true
			for i := 0; i < hashCount; i++ {
				if l.peekCharAt(i+1) != '#' {
					matched = false
					break
				}
			}
			if matched {
				l.readChar() // consume "
				for i := 0; i < hashCount; i++ {
					l.readChar() // consume #
				}
				break
			}
		}
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}

	return NewToken(TokenString, l.input[start:l.position-1], pos)
}

// readMultilineString reads a multiline string ("""...""").
func (l *Lexer) readMultilineString(pos Position) Token {
	start := l.position - 1
	l.readChar() // consume first "
	l.readChar() // consume second "
	l.readChar() // consume third "

	// Read until we find closing """
	for l.ch != 0 {
		if l.ch == '"' && l.peekChar() == '"' && l.peekCharAt(2) == '"' {
			l.readChar() // consume first "
			l.readChar() // consume second "
			l.readChar() // consume third "
			break
		}
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}

	return NewToken(TokenString, l.input[start:l.position-1], pos)
}

// readHashKeyword reads #true, #false, or #null.
func (l *Lexer) readHashKeyword(pos Position) Token {
	start := l.position - 1
	l.readChar() // consume #

	for isIdentifierChar(l.ch) {
		l.readChar()
	}

	literal := l.input[start : l.position-1]
	switch literal {
	case "#true":
		return NewToken(TokenTrue, literal, pos)
	case "#false":
		return NewToken(TokenFalse, literal, pos)
	case "#null":
		return NewToken(TokenNull, literal, pos)
	default:
		return NewToken(TokenError, fmt.Sprintf("unknown keyword: %s", literal), pos)
	}
}

// readNumber reads a number (decimal, hex, octal, binary).
func (l *Lexer) readNumber(pos Position) Token {
	start := l.position - 1

	// Handle sign
	if l.ch == '-' || l.ch == '+' {
		l.readChar()
	}

	// Check for hex, octal, binary
	if l.ch == '0' {
		next := l.peekChar()
		if next == 'x' || next == 'X' {
			l.readChar() // consume 0
			l.readChar() // consume x
			for isHexDigit(l.ch) || l.ch == '_' {
				l.readChar()
			}
			return NewToken(TokenNumber, l.input[start:l.position-1], pos)
		} else if next == 'o' || next == 'O' {
			l.readChar() // consume 0
			l.readChar() // consume o
			for isOctalDigit(l.ch) || l.ch == '_' {
				l.readChar()
			}
			return NewToken(TokenNumber, l.input[start:l.position-1], pos)
		} else if next == 'b' || next == 'B' {
			l.readChar() // consume 0
			l.readChar() // consume b
			for isBinaryDigit(l.ch) || l.ch == '_' {
				l.readChar()
			}
			return NewToken(TokenNumber, l.input[start:l.position-1], pos)
		}
	}

	// Decimal number
	for isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume .
		for isDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}
	}

	// Check for exponent
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}
	}

	return NewToken(TokenNumber, l.input[start:l.position-1], pos)
}

// readIdentifier reads an identifier or unquoted string.
func (l *Lexer) readIdentifier(pos Position) Token {
	start := l.position - 1

	for isIdentifierChar(l.ch) {
		l.readChar()
	}

	return NewToken(TokenIdentifier, l.input[start:l.position-1], pos)
}

// Helper functions for character classification

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch rune) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isOctalDigit(ch rune) bool {
	return ch >= '0' && ch <= '7'
}

func isBinaryDigit(ch rune) bool {
	return ch == '0' || ch == '1'
}

func isIdentifierStart(ch rune) bool {
	// KDL allows most unicode characters in identifiers except special ones
	if ch == 0 {
		return false
	}
	// Disallowed characters
	disallowed := "\\(){}[]/;\"#=\n\r\t "
	if strings.ContainsRune(disallowed, ch) {
		return false
	}
	return true
}

func isIdentifierChar(ch rune) bool {
	if ch == 0 {
		return false
	}
	// Disallowed characters
	disallowed := "\\(){}[]/;\"#=\n\r\t "
	if strings.ContainsRune(disallowed, ch) {
		return false
	}
	return unicode.IsPrint(ch)
}
