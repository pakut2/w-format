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
						&ast.Identifier{
							Token: token.Token{Type: token.IDENTIFIER, Literal: "text", LineNumber: 7},
							Value: "text",
						},
						&ast.Identifier{
							Token: token.Token{Type: token.IDENTIFIER, Literal: "number1", LineNumber: 7},
							Value: "number1",
						},
						&ast.Identifier{
							Token: token.Token{Type: token.IDENTIFIER, Literal: "number2", LineNumber: 7},
							Value: "number2",
						},
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
					Token: token.Token{Type: token.GREATER_THAN, Literal: ">", LineNumber: 9},
					Left: &ast.InfixExpression{
						Token: token.Token{Type: token.SLASH, Literal: "/", LineNumber: 9},
						Left: &ast.InfixExpression{
							Token: token.Token{Type: token.PLUS, Literal: "+", LineNumber: 9},
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENTIFIER, Literal: "number1", LineNumber: 9},
								Value: "number1",
							},
							Operator: "+",
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "2", LineNumber: 9},
								Value: 2,
							},
						},
						Operator: "/",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "2", LineNumber: 9},
							Value: 2,
						},
					},
					Operator: ">",
					Right: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "1000", LineNumber: 9},
						Value: 1000,
					},
				},
			},
			&ast.ExpressionStatement{
				Expression: &ast.InfixExpression{
					Token: token.Token{Type: token.EQUALS, Literal: "===", LineNumber: 10},
					Left: &ast.Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "expression", LineNumber: 10},
						Value: "expression",
					},
					Operator: "===",
					Right: &ast.IntegerLiteral{
						Token: token.Token{Type: token.TRUE, Literal: "true", LineNumber: 10},
						Value: 1,
					},
				},
			},
			&ast.IfStatement{
				Token: token.Token{Type: token.IF, Literal: "if", LineNumber: 12},
				Condition: &ast.IntegerLiteral{
					Token: token.Token{Type: token.FALSE, Literal: "false", LineNumber: 12},
					Value: 0,
				},
				Consequence: &ast.BlockStatement{
					Token: token.Token{Type: token.LEFT_BRACE, Literal: "{", LineNumber: 12},
					Statements: []ast.Statement{
						&ast.AssignmentStatement{
							Token: token.Token{Type: token.IDENTIFIER, Literal: "expression", LineNumber: 13},
							Name: &ast.Identifier{
								Token: token.Token{Type: token.IDENTIFIER, Literal: "expression", LineNumber: 13},
								Value: "expression",
							},
							Value: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "1", LineNumber: 13},
								Value: 1,
							},
						},
					},
				},
				Alternative: &ast.BlockStatement{
					Token: token.Token{Type: token.LEFT_BRACE, Literal: "{", LineNumber: 14},
					Statements: []ast.Statement{
						&ast.AssignmentStatement{
							Token: token.Token{Type: token.IDENTIFIER, Literal: "expression", LineNumber: 15},
							Name: &ast.Identifier{
								Token: token.Token{Type: token.IDENTIFIER, Literal: "expression", LineNumber: 15},
								Value: "expression",
							},
							Value: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "2", LineNumber: 15},
								Value: 2,
							},
						},
					},
				},
			},
			&ast.ForStatement{
				Token: token.Token{Type: token.FOR, Literal: "for", LineNumber: 18},
				Declaration: &ast.LetStatement{
					Token: token.Token{Type: token.LET, Literal: "let", LineNumber: 18},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "i", LineNumber: 18},
						Value: "i",
					},
					Value: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "0", LineNumber: 18},
						Value: 0,
					},
				},
				Boundary: &ast.InfixExpression{
					Token: token.Token{Type: token.LESS_THAN, Literal: "<", LineNumber: 18},
					Left: &ast.Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "i", LineNumber: 18},
						Value: "i",
					},
					Operator: "<",
					Right: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "10", LineNumber: 18},
						Value: 10,
					},
				},
				Increment: &ast.SuffixExpression{
					Token: token.Token{Type: token.INCREMENT, Literal: "++", LineNumber: 18},
					Left: &ast.Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "i", LineNumber: 18},
						Value: "i",
					},
					Operator: "++",
				},
				Body: &ast.BlockStatement{
					Token: token.Token{Type: token.LEFT_BRACE, Literal: "{", LineNumber: 18},
					Statements: []ast.Statement{
						&ast.IfStatement{
							Token: token.Token{Type: token.IF, Literal: "if", LineNumber: 19},
							Condition: &ast.InfixExpression{
								Token: token.Token{Type: token.OR, Literal: "||", LineNumber: 19},
								Left: &ast.InfixExpression{
									Token: token.Token{Type: token.EQUALS, Literal: "===", LineNumber: 19},
									Left: &ast.InfixExpression{
										Token: token.Token{Type: token.PERCENT, Literal: "%", LineNumber: 19},
										Left: &ast.Identifier{
											Token: token.Token{Type: token.IDENTIFIER, Literal: "i", LineNumber: 19},
											Value: "i",
										},
										Operator: "%",
										Right: &ast.IntegerLiteral{
											Token: token.Token{Type: token.INT, Literal: "2", LineNumber: 19},
											Value: 2,
										},
									},
									Operator: "===",
									Right: &ast.IntegerLiteral{
										Token: token.Token{Type: token.INT, Literal: "0", LineNumber: 19},
										Value: 0,
									},
								},
								Operator: "||",
								Right: &ast.InfixExpression{
									Token: token.Token{Type: token.EQUALS, Literal: "===", LineNumber: 19},
									Left: &ast.Identifier{
										Token: token.Token{Type: token.IDENTIFIER, Literal: "i", LineNumber: 19},
										Value: "i",
									},
									Operator: "===",
									Right: &ast.IntegerLiteral{
										Token: token.Token{Type: token.INT, Literal: "8", LineNumber: 19},
										Value: 8,
									},
								},
							},
							Consequence: &ast.BlockStatement{
								Token: token.Token{Type: token.LEFT_BRACE, Literal: "{", LineNumber: 19},
								Statements: []ast.Statement{
									&ast.ContinueStatement{
										Token: token.Token{Type: token.CONTINUE, Literal: "continue", LineNumber: 20},
									},
								},
							},
							Alternative: nil,
						},
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
		parsedAstJson, _ := json.MarshalIndent(parsedAst.Statements, "", "  ")

		t.Fatalf("invalid ast. expected=%s, got=%s", expectedAstJson, parsedAstJson)
	}
}
