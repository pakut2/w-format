package jsWhitespaceParser

import (
	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser/internal/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	currentChar  byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.currentChar {
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
	case ',':
		tok = newToken(token.COMMA, l.currentChar)
	case '(':
		tok = newToken(token.LEFT_PARENTHESIS, l.currentChar)
	case ')':
		tok = newToken(token.RIGHT_PARENTHESIS, l.currentChar)
	case '"', '\'', '`':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.currentChar) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.IDENTIFIER

			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.currentChar)
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.currentChar) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_' || char == '.'
}

func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()
		if l.currentChar == '"' || l.currentChar == '\'' || l.currentChar == '`' || l.currentChar == 0 {
			break
		}
	}

	return l.input[position:l.position]
}
