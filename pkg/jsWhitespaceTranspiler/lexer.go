package jsWhitespaceTranspiler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pakut2/w-format/internal/utilities"
	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/token"
)

type Lexer struct {
	input bufio.Reader

	previousChar      rune
	currentChar       rune
	currentLineNumber int
}

func NewLexer(input io.Reader) *Lexer {
	l := &Lexer{input: *bufio.NewReader(input), currentLineNumber: 1}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	l.previousChar = l.currentChar
	l.currentChar = utilities.ReadRune(&l.input)
}

func (l *Lexer) NextToken() token.Token {
	var currentToken token.Token

	l.skipWhitespace()

	switch l.currentChar {
	case '*':
		currentToken = token.NewTokenFromChar(token.ASTERISK, l.currentChar, l.currentLineNumber)
	case '%':
		currentToken = token.NewTokenFromChar(token.PERCENT, l.currentChar, l.currentLineNumber)
	case ';':
		currentToken = token.NewTokenFromChar(token.SEMICOLON, l.currentChar, l.currentLineNumber)
	case ',':
		currentToken = token.NewTokenFromChar(token.COMMA, l.currentChar, l.currentLineNumber)
	case '(':
		currentToken = token.NewTokenFromChar(token.LEFT_PARENTHESIS, l.currentChar, l.currentLineNumber)
	case ')':
		currentToken = token.NewTokenFromChar(token.RIGHT_PARENTHESIS, l.currentChar, l.currentLineNumber)
	case '{':
		currentToken = token.NewTokenFromChar(token.LEFT_BRACE, l.currentChar, l.currentLineNumber)
	case '}':
		currentToken = token.NewTokenFromChar(token.RIGHT_BRACE, l.currentChar, l.currentLineNumber)
	case '"', '\'', '`':
		currentToken = token.NewTokenFromString(token.STRING, l.readString(), l.currentLineNumber)
	case '&':
		if utilities.PeekRune(l.input) == '&' {
			startingChar := l.currentChar
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.AND,
				fmt.Sprintf("%c%c", startingChar, l.currentChar),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.ILLEGAL, l.currentChar, l.currentLineNumber)
		}
	case '|':
		if utilities.PeekRune(l.input) == '|' {
			startingChar := l.currentChar
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.OR,
				fmt.Sprintf("%c%c", startingChar, l.currentChar),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.ILLEGAL, l.currentChar, l.currentLineNumber)
		}
	case '+':
		if utilities.PeekRune(l.input) == '+' {
			startingChar := l.currentChar
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.INCREMENT,
				fmt.Sprintf("%c%c", startingChar, l.currentChar),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.PLUS, l.currentChar, l.currentLineNumber)
		}
	case '-':
		if utilities.PeekRune(l.input) == '-' {
			startingChar := l.currentChar
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.DECREMENT,
				fmt.Sprintf("%c%c", startingChar, l.currentChar),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.MINUS, l.currentChar, l.currentLineNumber)
		}
	case '/':
		nextChar := utilities.PeekRune(l.input)
		if nextChar == '/' || nextChar == '*' {
			panic("comments are a violation of DRY")
		}

		currentToken = token.NewTokenFromChar(token.SLASH, l.currentChar, l.currentLineNumber)
	case '=':
		nextChars, err := utilities.PeekTwoRunes(l.input)
		if err == nil && nextChars == "==" {
			startingChar := l.currentChar
			l.readChar()
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.EQUALS,
				fmt.Sprintf("%c%s", startingChar, nextChars),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.ASSIGN, l.currentChar, l.currentLineNumber)
		}
	case '!':
		nextChars, err := utilities.PeekTwoRunes(l.input)
		if err == nil && nextChars == "==" {
			startingChar := l.currentChar
			l.readChar()
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.NOT_EQUALS,
				fmt.Sprintf("%c%s", startingChar, nextChars),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.BANG, l.currentChar, l.currentLineNumber)
		}
	case '<':
		if utilities.PeekRune(l.input) == '=' {
			startingChar := l.currentChar
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.LESS_THAN_OR_EQUAL,
				fmt.Sprintf("%c%c", startingChar, l.currentChar),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.LESS_THAN, l.currentChar, l.currentLineNumber)
		}
	case '>':
		if utilities.PeekRune(l.input) == '=' {
			startingChar := l.currentChar
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.GREATER_THAN_OR_EQUAL,
				fmt.Sprintf("%c%c", startingChar, l.currentChar),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.GREATER_THAN, l.currentChar, l.currentLineNumber)
		}
	case 0:
		currentToken = token.NewTokenFromString(token.EOF, "", l.currentLineNumber)
	default:
		if l.isLetter() {
			currentToken.Literal = l.readIdentifier()
			currentToken.Type = token.LookupIdentifier(currentToken.Literal)
			currentToken.LineNumber = l.currentLineNumber

			return currentToken
		} else if l.isDigit() {
			return token.NewTokenFromString(token.INT, l.readNumber(), l.currentLineNumber)
		} else {
			currentToken = token.NewTokenFromChar(token.ILLEGAL, l.currentChar, l.currentLineNumber)
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
	startingQuote := l.currentChar

	var stringLiteral string

	for {
		l.readChar()

		if l.currentChar == 0 || (l.currentChar == startingQuote && l.previousChar != '\\') {
			break
		}

		if l.currentChar != '\\' {
			stringLiteral = fmt.Sprintf("%s%c", stringLiteral, l.currentChar)
		}
	}

	return stringLiteral
}

func (l *Lexer) isLetter() bool {
	return 'a' <= l.currentChar && l.currentChar <= 'z' || 'A' <= l.currentChar && l.currentChar <= 'Z' || l.currentChar == '_' || l.currentChar == '.'
}

func (l *Lexer) readIdentifier() string {
	var identifier string

	for l.isLetter() || l.isDigit() {
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
