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
	input := `
console.log('Hello', 42);

let text = "value";
let number1 = 1337;
let number2 = number1;
console.log(text, number1, number2);

let expression = (number1 + 2) / 2 > 1000;
expression === true;
`
	expectedAst := &ast.Program{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Expression: &ast.CallExpression{
					Function: &ast.Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "console.log", LineNumber: 2},
						Value: "console.log",
					},
					Arguments: []ast.Expression{
						&ast.StringLiteral{
							Token: token.Token{Type: token.STRING, Literal: "Hello", LineNumber: 2},
							Value: "Hello",
						},
						&ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "42", LineNumber: 2},
							Value: 42,
						},
					},
				},
			},
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let", LineNumber: 4},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "text", LineNumber: 4},
					Value: "text",
				},
				Value: &ast.StringLiteral{
					Token: token.Token{Type: token.STRING, Literal: "value", LineNumber: 4},
					Value: "value",
				},
			},
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let", LineNumber: 5},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "number1", LineNumber: 5},
					Value: "number1",
				},
				Value: &ast.IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "1337", LineNumber: 5},
					Value: 1337,
				},
			},
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let", LineNumber: 6},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "number2", LineNumber: 6},
					Value: "number2",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "number1", LineNumber: 6},
					Value: "number1",
				},
			},
			&ast.ExpressionStatement{
				Expression: &ast.CallExpression{
					Function: &ast.Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "console.log", LineNumber: 7},
						Value: "console.log",
					},
					Arguments: []ast.Expression{
						&ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Literal: "text", LineNumber: 7}, Value: "text"},
						&ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Literal: "number1", LineNumber: 7}, Value: "number1"},
						&ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Literal: "number2", LineNumber: 7}, Value: "number2"},
					},
				},
			},
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let", LineNumber: 9},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "expression", LineNumber: 9},
					Value: "expression",
				},
				Value: &ast.InfixExpression{
					Token:    token.Token{Type: token.TokenType(">"), Literal: ">", LineNumber: 9},
					Operator: ">",
					Left: &ast.InfixExpression{
						Token:    token.Token{Type: token.TokenType("/"), Literal: "/", LineNumber: 9},
						Operator: "/",
						Left: &ast.InfixExpression{
							Token:    token.Token{Type: token.TokenType("+"), Literal: "+", LineNumber: 9},
							Operator: "+",
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENTIFIER, Literal: "number1", LineNumber: 9},
								Value: "number1",
							},
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "2", LineNumber: 9},
								Value: 2,
							},
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "2", LineNumber: 9},
							Value: 2,
						},
					},
					Right: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "1000", LineNumber: 9},
						Value: 1000,
					},
				},
			},
			&ast.ExpressionStatement{
				Expression: &ast.InfixExpression{
					Token:    token.Token{Type: token.TokenType("==="), Literal: "===", LineNumber: 10},
					Operator: "===",
					Left: &ast.Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "expression", LineNumber: 10},
						Value: "expression",
					},
					Right: &ast.IntegerLiteral{
						Token: token.Token{Type: token.TRUE, Literal: "true", LineNumber: 10},
						Value: 1,
					},
				},
			},
		},
	}

	lexer := NewLexer(strings.NewReader(input))
	parser := NewParser(lexer)

	parsedAst := parser.ParseProgram()

	if !reflect.DeepEqual(parsedAst, expectedAst) {
		expectedAstJson, _ := json.MarshalIndent(expectedAst, "", "  ")
		parsedAstJson, _ := json.MarshalIndent(parsedAst.Statements[0], "", "  ")

		t.Fatalf("invalid ast. expected=%s, got=%s", expectedAstJson, parsedAstJson)
	}
}
