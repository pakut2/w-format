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

let expression = (number1 + 2) / 2 > 1000;
expression === true;

if (false) {
	expression = 1;
} else {
	expression = 2;
}

for (let i = 0; i < 10; i++) { 
	if (i % 2 === 0 || i === 8) {
		continue;
	}
}
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

		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "expression"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.LEFT_PARENTHESIS, Literal: "("},
		{Type: token.IDENTIFIER, Literal: "number1"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.INT, Literal: "2"},
		{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
		{Type: token.SLASH, Literal: "/"},
		{Type: token.INT, Literal: "2"},
		{Type: token.GREATER_THAN, Literal: ">"},
		{Type: token.INT, Literal: "1000"},
		{Type: token.SEMICOLON, Literal: ";"},

		{Type: token.IDENTIFIER, Literal: "expression"},
		{Type: token.EQUALS, Literal: "==="},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.SEMICOLON, Literal: ";"},

		{Type: token.IF, Literal: "if"},
		{Type: token.LEFT_PARENTHESIS, Literal: "("},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
		{Type: token.LEFT_BRACE, Literal: "{"},
		{Type: token.IDENTIFIER, Literal: "expression"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "1"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RIGHT_BRACE, Literal: "}"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.LEFT_BRACE, Literal: "{"},
		{Type: token.IDENTIFIER, Literal: "expression"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "2"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RIGHT_BRACE, Literal: "}"},

		{Type: token.FOR, Literal: "for"},
		{Type: token.LEFT_PARENTHESIS, Literal: "("},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "i"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "0"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IDENTIFIER, Literal: "i"},
		{Type: token.LESS_THAN, Literal: "<"},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IDENTIFIER, Literal: "i"},
		{Type: token.INCREMENT, Literal: "++"},
		{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
		{Type: token.LEFT_BRACE, Literal: "{"},
		{Type: token.IF, Literal: "if"},
		{Type: token.LEFT_PARENTHESIS, Literal: "("},
		{Type: token.IDENTIFIER, Literal: "i"},
		{Type: token.PERCENT, Literal: "%"},
		{Type: token.INT, Literal: "2"},
		{Type: token.EQUALS, Literal: "==="},
		{Type: token.INT, Literal: "0"},
		{Type: token.OR, Literal: "||"},
		{Type: token.IDENTIFIER, Literal: "i"},
		{Type: token.EQUALS, Literal: "==="},
		{Type: token.INT, Literal: "8"},
		{Type: token.RIGHT_PARENTHESIS, Literal: ")"},
		{Type: token.LEFT_BRACE, Literal: "{"},
		{Type: token.CONTINUE, Literal: "continue"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RIGHT_BRACE, Literal: "}"},
		{Type: token.RIGHT_BRACE, Literal: "}"},

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
