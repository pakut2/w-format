package main

import (
	"fmt"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser"
)

func main() {
	jsFile := "console.log('Hello', 'There!')"

	lexer := jsWhitespaceParser.NewLexer(jsFile)
	parser := jsWhitespaceParser.NewParser(lexer)

	ast := parser.ParseProgram()

	//astJson, err := json.MarshalIndent(ast.Statements, "", "  ")
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//fmt.Printf("%s\n", astJson)

	evaluator := jsWhitespaceParser.NewEvaluator()
	result := evaluator.Eval(ast)

	fmt.Printf("%s", result.Instruction())
}
