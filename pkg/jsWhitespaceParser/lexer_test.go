package jsWhitespaceParser

import (
	"strings"
	"testing"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser/internal/token"
)

func TestLexer(t *testing.T) {
	input := "console.log('Hello', 'There!');"

	expectedTokens := []token.Token{
		{Type: token.IDENTIFIER, Literal: "console.log"},
		{Type: token.LEFT_PARENTHESIS, Literal: "("},
		{Type: token.STRING, Literal: "Hello"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.STRING, Literal: "There!"},
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
