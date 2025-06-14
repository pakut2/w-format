package jsWhitespaceParser

import (
	"fmt"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser/internal/ast"
	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceParser/internal/object"
)

const (
	TAB       = '\t'
	LINE_FEED = '\n'
	SPACE     = ' '
	//TAB       = 'T'
	//LINE_FEED = 'L'
	//SPACE     = 'S'
)

type Evaluator struct {
	currentHeapAddress byte
	builtInFunctions   map[string]*object.BuildIn
}

func NewEvaluator() *Evaluator {
	e := &Evaluator{
		currentHeapAddress: 0,
	}

	e.builtInFunctions = make(map[string]*object.BuildIn)
	e.registerBuildInFunction("console.log", &object.BuildIn{Function: e.consoleLogBuiltInFunction})

	return e
}

func (e *Evaluator) registerBuildInFunction(functionName string, function *object.BuildIn) {
	e.builtInFunctions[functionName] = function
}

func (e *Evaluator) consoleLogBuiltInFunction(args ...object.Object) object.Object {
	var instruction string

	for i, arg := range args {
		switch arg := arg.(type) {
		case *object.String:
			instruction = fmt.Sprintf("%s%s", instruction, arg.InstructionBody)

			for _, char := range arg.Chars {
				printCharInstruction := fmt.Sprintf("%s%s",
					e.retrieveFromHeapInstruction(char.HeapAddress),
					e.printTopStackCharInstruction(),
				)

				instruction = fmt.Sprintf("%s%s", instruction, printCharInstruction)
			}

			if i != len(args)-1 {
				printSpaceInstruction := fmt.Sprintf("%s%s",
					e.pushNumberLiteralToStackInstruction(' '),
					e.printTopStackCharInstruction(),
				)

				instruction = fmt.Sprintf("%s%s", instruction, printSpaceInstruction)
			}
		default:
			panic(fmt.Sprintf("argument %s not supported", arg.Type()))
		}
	}

	printNewLineInstruction := fmt.Sprintf("%s%s",
		e.pushNumberLiteralToStackInstruction('\n'),
		e.printTopStackCharInstruction(),
	)

	instruction = fmt.Sprintf("%s%s", instruction, printNewLineInstruction)

	return &object.Void{InstructionBody: instruction}
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
	var instructions string

	for _, statement := range program.Statements {
		instructions += e.Eval(statement).Instruction()
	}

	return &object.Program{
		InstructionBody: fmt.Sprintf("%s%c%c%c", instructions, LINE_FEED, LINE_FEED, LINE_FEED),
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
		instruction := e.storeInHeapInstruction(heapAddress, c)

		result.Chars = append(result.Chars, object.Char{HeapAddress: e.currentHeapAddress})
		result.InstructionBody = fmt.Sprintf("%s%s", result.InstructionBody, instruction)
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
	case *object.BuildIn:
		return function.Function(args...)
	default:
		panic(fmt.Sprintf("not a function: %s", function.Type()))
	}
}

func (e *Evaluator) storeInHeapInstruction(heapAddress byte, value byte) string {
	instruction := fmt.Sprintf(
		"%s%s",
		e.pushNumberLiteralToStackInstruction(heapAddress),
		e.pushNumberLiteralToStackInstruction(value),
	)

	return fmt.Sprintf("%s%c%c%c", instruction, TAB, TAB, SPACE)
}

func (e *Evaluator) retrieveFromHeapInstruction(heapAddress byte) string {
	instruction := e.pushNumberLiteralToStackInstruction(heapAddress)

	return fmt.Sprintf("%s%c%c%c", instruction, TAB, TAB, TAB)
}

func (e *Evaluator) pushNumberLiteralToStackInstruction(value byte) string {
	return fmt.Sprintf("%c%c%s", SPACE, SPACE, e.numberLiteralInstruction(value))
}

//func (e *Evaluator) popFromStackInstruction() string {
//	return fmt.Sprintf("%c%c%c", SPACE, LINE_FEED, LINE_FEED)
//}

func (e *Evaluator) numberLiteralInstruction(value byte) string {
	var instruction string

	charBinary := fmt.Sprintf("%s%.8b", instruction, value)

	for _, bit := range charBinary {
		if bit == '1' {
			instruction = fmt.Sprintf("%s%c", instruction, TAB)

			continue
		}

		instruction = fmt.Sprintf("%s%c", instruction, SPACE)
	}

	return fmt.Sprintf("%c%s%c", SPACE, instruction, LINE_FEED)
}

func (e *Evaluator) printTopStackCharInstruction() string {
	return fmt.Sprintf("%c%c%c%c", TAB, LINE_FEED, SPACE, SPACE)
}
