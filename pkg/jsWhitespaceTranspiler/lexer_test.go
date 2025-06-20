package jsWhitespaceTranspiler

import (
	"strings"
	"testing"

	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/token"
)

func TestLexer(t *testing.T) {
	input := `
console.log('Hello', 42);

let text = "value";
let number1 = 1337;
let number2 = number1;
console.log(text, number1, number2);
`

	expectedTokens := []token.Token{
		{Type: token.IDENTIFIER, Literal: "console.log"},
		{Type: token.LEFT_PARENTHESIS, Literal: "("},
		{Type: token.STRING, Literal: "Hello"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.INT, Literal: "42"},
		{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},

		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "text"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.STRING, Literal: "value"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "number1"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "1337"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "number2"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.IDENTIFIER, Literal: "number1"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IDENTIFIER, Literal: "console.log"},
		{Type: token.LEFT_PARENTHESIS, Literal: "("},
		{Type: token.IDENTIFIER, Literal: "text"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENTIFIER, Literal: "number1"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENTIFIER, Literal: "number2"},
		{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	}

	lexer := NewLexer(strings.NewReader(input))

	for i, expectedToken := range expectedTokens {
		parsedToken := lexer.NextToken()

		if parsedToken.Type != expectedToken.Type {
			t.Errorf("token type (#%d) incorrect. expected=%q, got=%q", i+1, expectedToken.Type, parsedToken.Type)
		}

		if parsedToken.Literal != expectedToken.Literal {
			t.Errorf("token literal (#%d) incorrect. expected=%q, got=%q", i+1, expectedToken.Literal, parsedToken.Literal)
		}
	}
}
