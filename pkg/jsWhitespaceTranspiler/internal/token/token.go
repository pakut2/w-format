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
	INT        = "INT"

	ASSIGN                = "="
	PLUS                  = "+"
	MINUS                 = "-"
	BANG                  = "!"
	ASTERISK              = "*"
	SLASH                 = "/"
	PERCENT               = "%"
	EQUALS                = "==="
	NOT_EQUALS            = "!=="
	LESS_THAN             = "<"
	LESS_THAN_OR_EQUAL    = "<="
	GREATER_THAN          = ">"
	GREATER_THAN_OR_EQUAL = ">="
	INCREMENT             = "++"
	DECREMENT             = "--"

	COMMA             = ","
	SEMICOLON         = ";"
	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
	LEFT_BRACE        = "{"
	RIGHT_BRACE       = "}"

	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

var keywords = map[string]TokenType{
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}

	return IDENTIFIER
}

func NewTokenFromChar(tokenType TokenType, char rune, lineNumber int) Token {
	return Token{Type: tokenType, Literal: string(char), LineNumber: lineNumber}
}

func NewTokenFromString(tokenType TokenType, literal string, lineNumber int) Token {
	return Token{Type: tokenType, Literal: literal, LineNumber: lineNumber}
}
