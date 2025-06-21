package jsWhitespaceTranspiler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/token"
)

type Lexer struct {
	input bufio.Reader

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

func (l *Lexer) peekChar() rune {
	for peekBytes := 4; peekBytes > 0; peekBytes-- {
		peekCharResult, err := l.input.Peek(peekBytes)
		if err == nil {
			char, _ := utf8.DecodeRune(peekCharResult)
			if char == utf8.RuneError {
				return 0
			}

			return char
		}
	}

	return 0
}

func (l *Lexer) peekTwoChars() (string, error) {
	peekResultBuffer, err := l.input.Peek(8)
	if err != nil && len(peekResultBuffer) == 0 {
		return "", err
	}

	char1, char1Size := utf8.DecodeRune(peekResultBuffer)

	if char1 == utf8.RuneError && char1Size == 1 {
		return "", errors.New("malformed input")
	}

	if char1Size == 0 {
		return "", errors.New("malformed input")
	}

	peekResultChar2Buffer := peekResultBuffer[char1Size:]

	char2, char2Size := utf8.DecodeRune(peekResultChar2Buffer)

	if char2 == utf8.RuneError && char2Size == 1 {
		return "", errors.New("malformed input")
	}

	if char2Size == 0 {
		return "", errors.New("malformed input")
	}

	return string([]rune{char1, char2}), nil
}

func (l *Lexer) NextToken() token.Token {
	var currentToken token.Token

	l.skipWhitespace()

	switch l.currentChar {
	case '+':
		currentToken = token.NewTokenFromChar(token.PLUS, l.currentChar, l.currentLineNumber)
	case '-':
		currentToken = token.NewTokenFromChar(token.MINUS, l.currentChar, l.currentLineNumber)
	case '*':
		currentToken = token.NewTokenFromChar(token.ASTERISK, l.currentChar, l.currentLineNumber)
	case '/':
		currentToken = token.NewTokenFromChar(token.SLASH, l.currentChar, l.currentLineNumber)
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
	case '=':
		nextChars, err := l.peekTwoChars()
		if err == nil && nextChars == "==" {
			startingCharacter := l.currentChar
			l.readChar()
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.EQUALS,
				fmt.Sprintf("%c%s", startingCharacter, nextChars),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.ASSIGN, l.currentChar, l.currentLineNumber)
		}
	case '!':
		nextChars, err := l.peekTwoChars()
		if err == nil && nextChars == "==" {
			startingCharacter := l.currentChar
			l.readChar()
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.NOT_EQUALS,
				fmt.Sprintf("%c%s", startingCharacter, nextChars),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.BANG, l.currentChar, l.currentLineNumber)
		}
	case '<':
		if l.peekChar() == '=' {
			startingCharacter := l.currentChar
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.LESS_THAN_OR_EQUAL,
				fmt.Sprintf("%c%c", startingCharacter, l.currentChar),
				l.currentLineNumber,
			)
		} else {
			currentToken = token.NewTokenFromChar(token.LESS_THAN, l.currentChar, l.currentLineNumber)
		}
	case '>':
		if l.peekChar() == '=' {
			startingCharacter := l.currentChar
			l.readChar()

			currentToken = token.NewTokenFromString(
				token.GREATER_THAN_OR_EQUAL,
				fmt.Sprintf("%c%c", startingCharacter, l.currentChar),
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
