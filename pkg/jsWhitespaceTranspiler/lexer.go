package jsWhitespaceTranspiler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/token"
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
		currentToken.Type = token.EOF
		currentToken.Literal = ""
		currentToken.LineNumber = l.currentLineNumber
	default:
		if l.isLetter() {
			currentToken.Type = token.IDENTIFIER
			currentToken.Literal = l.readIdentifier()
			currentToken.LineNumber = l.currentLineNumber

			return currentToken
		} else if l.isDigit() {
			currentToken.Type = token.INT
			currentToken.Literal = l.readNumber()
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

func (l *Lexer) isLetter() bool {
	return 'a' <= l.currentChar && l.currentChar <= 'z' || 'A' <= l.currentChar && l.currentChar <= 'Z' || l.currentChar == '_' || l.currentChar == '.'
}

func (l *Lexer) readIdentifier() string {
	var identifier string

	for l.isLetter() {
		identifier = fmt.Sprintf("%s%c", identifier, l.currentChar)

		l.readChar()
	}

	return identifier
}

func (l *Lexer) isDigit() bool {
	return '0' <= l.currentChar && l.currentChar <= '9'
}

func (l *Lexer) readNumber() string {
	var numberLiteral string

	for l.isDigit() {
		numberLiteral = fmt.Sprintf("%s%c", numberLiteral, l.currentChar)

		l.readChar()
	}

	return numberLiteral
}
