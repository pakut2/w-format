package jsWhitespaceTranspiler

import (
	"strings"
	"testing"

	"github.com/pakut2/w-format/pkg/whitespace"
)

func TestPrint(t *testing.T) {
	input := `
console.log('Hello', 42);

let text = "value";
let number1 = 1337;
let number2 = number1;
console.log(text, number1, number2);

let expression = (number1 + 2) / 2 > 1000;
expression === true;
`

	expectedInstructions := []whitespace.Instruction{
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.LINE_FEED, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.LINE_FEED, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.LINE_FEED, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.LINE_FEED, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.LINE_FEED, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.LINE_FEED, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.TAB, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.SPACE, whitespace.LINE_FEED, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.LINE_FEED, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.LINE_FEED, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.SPACE, whitespace.LINE_FEED, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.SPACE, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.TAB, whitespace.SPACE, whitespace.SPACE, whitespace.LINE_FEED}},
		{Body: []whitespace.Token{whitespace.SPACE, whitespace.LINE_FEED, whitespace.TAB}},
		{Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE}},
		{Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.LINE_FEED, whitespace.LINE_FEED}},
	}

	lexer := NewLexer(strings.NewReader(input))
	parsedAst := NewParser(lexer).ParseProgram()

	whitespaceProgram := NewTranspiler().Transpile(parsedAst)

	for i, instruction := range whitespaceProgram.Instructions() {
		currentInstruction := instruction.String()
		expectedInstruction := expectedInstructions[i].String()

		if currentInstruction != expectedInstruction {
			t.Errorf("instruction (#%d) incorrect. expected=%q, got=%q", i+1, expectedInstruction, currentInstruction)
		}
	}
}
