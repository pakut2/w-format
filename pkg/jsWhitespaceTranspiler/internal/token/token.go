package token

type TokenType string

type Token struct {
	Type       TokenType
	Literal    string
	LineNumber int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"

	COMMA             = ","
	SEMICOLON         = ";"
	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
)

func NewToken(tokenType TokenType, char rune, lineNumber int) Token {
	return Token{Type: tokenType, Literal: string(char), LineNumber: lineNumber}
}
