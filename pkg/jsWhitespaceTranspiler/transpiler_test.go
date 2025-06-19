package jsWhitespaceTranspiler

import (
	"testing"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceTranspiler/internal/ast"
	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceTranspiler/internal/token"
	"github.com/pakut2/js-whitespace/pkg/whitespace"
)

func TestEvaluator(t *testing.T) {
	mockAst := &ast.Program{
		Statements: []ast.Statement{
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
			},
		},
	}

	expectedInstructions := []whitespace.Instruction{
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.LINE_FEED, whitespace.LINE_FEED}},
	}

	transpiler := NewTranspiler()
	instructions := transpiler.Transpile(mockAst)

	for i, instruction := range instructions.Instruction() {
		currentInstruction := instruction.String()
		expectedInstruction := expectedInstructions[i].String()

		if currentInstruction != expectedInstruction {
			t.Errorf("instruction (#%d) incorrect. expected=%q, got=%q", i+1, currentInstruction, expectedInstruction)
		}
	}
}
