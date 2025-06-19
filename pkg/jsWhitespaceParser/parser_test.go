package jsWhitespaceParser

import (
	"reflect"
	"testing"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser/internal/ast"
	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser/internal/token"
)

func TestParser(t *testing.T) {
	input := "console.log('Hello', 'There!');"

	expectedAst :=
		&ast.ExpressionStatement{
			Expression: &ast.CallExpression{
				Function: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "console.log",
					},
					Value: "console.log",
				},
				Arguments: []ast.Expression{
					&ast.StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "Hello",
						},
						Value: "Hello",
					},
					&ast.StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "There!",
						},
						Value: "There!",
					},
				},
			},
		}

	lexer := NewLexer(input)
	parser := NewParser(lexer)

	parsedAst := parser.ParseProgram()

	if len(parsedAst.Statements) != 1 {
		t.Fatalf("invalid number of statements. expected=1, got=%d", len(parsedAst.Statements))
	}

	if !reflect.DeepEqual(parsedAst.Statements[0], expectedAst) {
		t.Fatalf("invalid ast. expected=%v, got=%v", expectedAst, parsedAst.Statements[0])
	}
}
