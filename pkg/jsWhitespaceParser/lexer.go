package jsWhitespaceParser

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser/internal/token"
)

type Lexer struct {
	input       bufio.Reader
	currentChar rune
}

func NewLexer(input io.Reader) *Lexer {
	l := &Lexer{input: *bufio.NewReader(input)}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	char, _, err := l.input.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.currentChar = 0

			return
		} else {
			panic(err)
		}
	}

	l.currentChar = char
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
	var identifier string

	for isLetter(l.currentChar) {
		identifier = fmt.Sprintf("%s%c", identifier, l.currentChar)

		l.readChar()
	}

	return identifier
}

func newToken(tokenType token.TokenType, char rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_' || char == '.'
}

func (l *Lexer) readString() string {
	var stringLiteral string

	for {
		l.readChar()

		if l.currentChar == '"' || l.currentChar == '\'' || l.currentChar == '`' || l.currentChar == 0 {
			break
		}

		stringLiteral = fmt.Sprintf("%s%c", stringLiteral, l.currentChar)
	}

	return stringLiteral
}
