package jsWhitespaceParser

import (
	"fmt"

	"github.com/pakut2/js-whitespace/pkg/whitespace"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser/internal/ast"
	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser/internal/object"
)

type Evaluator struct {
	instructions       []whitespace.Instruction
	currentHeapAddress byte
	builtInFunctions   map[string]*object.BuiltIn
}

func NewEvaluator() *Evaluator {
	e := &Evaluator{
		currentHeapAddress: 0,
	}

	e.builtInFunctions = make(map[string]*object.BuiltIn)
	e.registerBuildInFunction("console.log", &object.BuiltIn{Function: e.consoleLogBuiltInFunction})

	return e
}

func (e *Evaluator) registerBuildInFunction(functionName string, function *object.BuiltIn) {
	e.builtInFunctions[functionName] = function
}

func (e *Evaluator) addInstruction(instruction whitespace.Instruction) {
	e.instructions = append(e.instructions, instruction)
}

func (e *Evaluator) consoleLogBuiltInFunction(args ...object.Object) object.Object {
	for i, arg := range args {
		switch arg := arg.(type) {
		case *object.String:
			for _, char := range arg.Chars {
				e.retrieveFromHeapInstruction(char.HeapAddress)
				e.printTopStackCharInstruction()
			}

			if i != len(args)-1 {
				e.pushNumberLiteralToStackInstruction(' ')
				e.printTopStackCharInstruction()
			}
		default:
			panic(fmt.Sprintf("argument %s not supported", arg.Type()))
		}
	}

	e.pushNumberLiteralToStackInstruction('\n')
	e.printTopStackCharInstruction()

	return &object.Void{}
}

func (e *Evaluator) getCurrentHeapAddressWithIncrement() byte {
	e.currentHeapAddress++

	return e.currentHeapAddress
}

func (e *Evaluator) Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.StringLiteral:
		return e.evalString([]byte(node.Value))
	case *ast.Identifier:
		return e.evalIdentifier(node)
	case *ast.CallExpression:
		function := e.Eval(node.Function)
		args := e.evalExpressions(node.Arguments)

		return e.applyFunction(function, args)
	}

	return nil
}

func (e *Evaluator) evalProgram(program *ast.Program) object.Object {
	for _, statement := range program.Statements {
		e.Eval(statement)
	}

	e.addInstruction(
		whitespace.Instruction{
			Body: []whitespace.Token{whitespace.LINE_FEED, whitespace.LINE_FEED, whitespace.LINE_FEED},
		},
	)

	return &object.Program{
		Instructions: e.instructions,
	}
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier) object.Object {
	if buildInFunction, ok := e.builtInFunctions[node.Value]; ok {
		return buildInFunction
	}

	panic(fmt.Sprintf("identifier %s not implemeted", node.Value))
}

func (e *Evaluator) evalString(value []byte) object.Object {
	var result object.String

	for _, c := range value {
		heapAddress := e.getCurrentHeapAddressWithIncrement()
		e.storeInHeapInstruction(heapAddress, c)

		result.Chars = append(result.Chars, object.Char{HeapAddress: e.currentHeapAddress})
	}

	return &result
}

func (e *Evaluator) evalExpressions(expressions []ast.Expression) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := e.Eval(expression)

		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) applyFunction(function object.Object, args []object.Object) object.Object {
	switch function := function.(type) {
	case *object.BuiltIn:
		return function.Function(args...)
	default:
		panic(fmt.Sprintf("not a function: %s", function.Type()))
	}
}

func (e *Evaluator) storeInHeapInstruction(heapAddress byte, value byte) {
	e.pushNumberLiteralToStackInstruction(heapAddress)
	e.pushNumberLiteralToStackInstruction(value)
	e.addInstruction(
		whitespace.Instruction{
			Body: []whitespace.Token{whitespace.TAB, whitespace.TAB, whitespace.SPACE},
		},
	)
}

func (e *Evaluator) retrieveFromHeapInstruction(heapAddress byte) {
	e.pushNumberLiteralToStackInstruction(heapAddress)
	e.addInstruction(
		whitespace.Instruction{
			Body: []whitespace.Token{
				whitespace.TAB, whitespace.TAB, whitespace.TAB,
			},
		},
	)
}

func (e *Evaluator) pushNumberLiteralToStackInstruction(value byte) {
	e.addInstruction(
		whitespace.Instruction{
			Body: append(
				[]whitespace.Token{whitespace.SPACE, whitespace.SPACE},
				e.prepareNumberLiteralInstruction(value).Body...,
			),
		},
	)

}

//func (e *Evaluator) popFromStackInstruction() string {
//	return fmt.Sprintf("%c%c%c", whitespace.SPACE, whitespace.LINE_FEED, whitespace.LINE_FEED)
//}

func (e *Evaluator) prepareNumberLiteralInstruction(value byte) whitespace.Instruction {
	instruction := whitespace.Instruction{Body: []whitespace.Token{whitespace.SPACE}}

	charBinary := fmt.Sprintf("%s%.8b", instruction, value)

	for _, bit := range charBinary {
		if bit == '1' {
			instruction.Body = append(instruction.Body, whitespace.TAB)

			continue
		}

		instruction.Body = append(instruction.Body, whitespace.SPACE)
	}

	instruction.Body = append(instruction.Body, whitespace.LINE_FEED)

	return instruction
}

func (e *Evaluator) printTopStackCharInstruction() {
	e.addInstruction(whitespace.Instruction{
		Body: []whitespace.Token{whitespace.TAB, whitespace.LINE_FEED, whitespace.SPACE, whitespace.SPACE},
	})
}
