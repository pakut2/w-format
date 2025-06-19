package jsWhitespaceTranspiler

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/ast"
	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/token"
)

func TestParser(t *testing.T) {
	input := "console.log('Hello', 'There!');"

	expectedAst :=
		&ast.ExpressionStatement{
			Expression: &ast.CallExpression{
				Function: &ast.Identifier{
					Token: token.Token{
						Type:       token.IDENTIFIER,
						Literal:    "console.log",
						LineNumber: 1,
					},
					Value: "console.log",
				},
				Arguments: []ast.Expression{
					&ast.StringLiteral{
						Token: token.Token{
							Type:       token.STRING,
							Literal:    "Hello",
							LineNumber: 1,
						},
						Value: "Hello",
					},
					&ast.StringLiteral{
						Token: token.Token{
							Type:       token.STRING,
							Literal:    "There!",
							LineNumber: 1,
						},
						Value: "There!",
					},
				},
			},
		}

	lexer := NewLexer(strings.NewReader(input))
	parser := NewParser(lexer)

	parsedAst := parser.ParseProgram()

	if len(parsedAst.Statements) != 1 {
		t.Fatalf("invalid number of statements. expected=1, got=%d", len(parsedAst.Statements))
	}

	if !reflect.DeepEqual(parsedAst.Statements[0], expectedAst) {
		expectedAstJson, _ := json.MarshalIndent(expectedAst, "", "  ")
		parsedAstJson, _ := json.MarshalIndent(parsedAst.Statements[0], "", "  ")

		t.Fatalf("invalid ast. expected=%s, got=%s", expectedAstJson, parsedAstJson)
	}
}
