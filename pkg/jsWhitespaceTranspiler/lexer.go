package jsWhitespaceTranspiler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceTranspiler/internal/token"
)

type Lexer struct {
	input             bufio.Reader
	currentChar       rune
	currentLineNumber int
}

func NewLexer(input io.Reader) *Lexer {
	l := &Lexer{input: *bufio.NewReader(input), currentLineNumber: 1}
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
			panic(fmt.Sprintf("input processing error: %v", err))
		}
	}

	l.currentChar = char
}

func (l *Lexer) NextToken() token.Token {
	var currentToken token.Token

	l.skipWhitespace()

	switch l.currentChar {
	case ';':
		currentToken = token.NewToken(token.SEMICOLON, l.currentChar, l.currentLineNumber)
	case ',':
		currentToken = token.NewToken(token.COMMA, l.currentChar, l.currentLineNumber)
	case '(':
		currentToken = token.NewToken(token.LEFT_PARENTHESIS, l.currentChar, l.currentLineNumber)
	case ')':
		currentToken = token.NewToken(token.RIGHT_PARENTHESIS, l.currentChar, l.currentLineNumber)
	case '"', '\'', '`':
		currentToken.Type = token.STRING
		currentToken.Literal = l.readString()
		currentToken.LineNumber = l.currentLineNumber
	case 0:
		currentToken.Literal = ""
		currentToken.Type = token.EOF
		currentToken.LineNumber = l.currentLineNumber
	default:
		if isLetter(l.currentChar) {
			currentToken.Literal = l.readIdentifier()
			currentToken.Type = token.IDENTIFIER
			currentToken.LineNumber = l.currentLineNumber

			return currentToken
		} else {
			currentToken = token.NewToken(token.ILLEGAL, l.currentChar, l.currentLineNumber)
		}
	}

	l.readChar()

	return currentToken
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		if l.currentChar == '\n' {
			l.currentLineNumber++
		}

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
