package jsWhitespaceTranspiler

import (
	"fmt"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceTranspiler/internal/ast"
	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceTranspiler/internal/object"
	"github.com/pakut2/js-whitespace/pkg/whitespace"
)

type Transpiler struct {
	instructions       []whitespace.Instruction
	currentHeapAddress byte
	builtInFunctions   map[string]*object.BuiltIn
}

func NewTranspiler() *Transpiler {
	t := &Transpiler{
		currentHeapAddress: 0,
	}

	t.builtInFunctions = make(map[string]*object.BuiltIn)
	t.registerBuildInFunction("console.log", &object.BuiltIn{Function: t.consoleLogBuiltInFunction})

	return t
}

func (t *Transpiler) registerBuildInFunction(functionName string, function *object.BuiltIn) {
	t.builtInFunctions[functionName] = function
}

func (t *Transpiler) addInstruction(instruction whitespace.Instruction) {
	t.instructions = append(t.instructions, instruction)
}

func (t *Transpiler) consoleLogBuiltInFunction(args ...object.Object) object.Object {
	for i, arg := range args {
		switch arg := arg.(type) {
		case *object.String:
			for _, char := range arg.Chars {
				t.retrieveFromHeapInstruction(char.HeapAddress)
				t.printTopStackCharInstruction()
			}

			if i != len(args)-1 {
				t.pushNumberLiteralToStackInstruction(' ')
				t.printTopStackCharInstruction()
			}
		default:
			panic(fmt.Sprintf("argument %s not supported", arg.Type()))
		}
	}

	t.pushNumberLiteralToStackInstruction('\n')
	t.printTopStackCharInstruction()

	return &object.Void{}
}

func (t *Transpiler) getCurrentHeapAddressWithIncrement() byte {
	t.currentHeapAddress++

	return t.currentHeapAddress
}

func (t *Transpiler) Transpile(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return t.transpileProgram(node)
	case *ast.ExpressionStatement:
		return t.Transpile(node.Expression)
	case *ast.StringLiteral:
		return t.transpileString([]byte(node.Value))
	case *ast.Identifier:
		return t.transpileIdentifier(node)
	case *ast.CallExpression:
		function := t.Transpile(node.Function)
		args := t.transpileExpressions(node.Arguments)

		return t.applyFunction(function, args)
	}

	return nil
}

func (t *Transpiler) transpileProgram(program *ast.Program) object.Object {
	for _, statement := range program.Statements {
		t.Transpile(statement)
	}

	t.addInstruction(whitespace.EndProgram())

	return &object.Program{
		Instructions: t.instructions,
	}
}

func (t *Transpiler) transpileIdentifier(identifier *ast.Identifier) object.Object {
	if buildInFunction, ok := t.builtInFunctions[identifier.Value]; ok {
		return buildInFunction
	}

	panic(fmt.Sprintf("[:%d] identifier %s not implemeted", identifier.Token.LineNumber, identifier.Value))
}

func (t *Transpiler) transpileString(value []byte) object.Object {
	var result object.String

	for _, c := range value {
		heapAddress := t.getCurrentHeapAddressWithIncrement()
		t.storeInHeapInstruction(heapAddress, c)

		result.Chars = append(result.Chars, object.Char{HeapAddress: t.currentHeapAddress})
	}

	return &result
}

func (t *Transpiler) transpileExpressions(expressions []ast.Expression) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		transpiledExpression := t.Transpile(expression)

		result = append(result, transpiledExpression)
	}

	return result
}

func (t *Transpiler) applyFunction(function object.Object, args []object.Object) object.Object {
	switch function := function.(type) {
	case *object.BuiltIn:
		return function.Function(args...)
	default:
		panic(fmt.Sprintf("%s is not a function", function.Type()))
	}
}

func (t *Transpiler) storeInHeapInstruction(heapAddress byte, value byte) {
	t.pushNumberLiteralToStackInstruction(heapAddress)
	t.pushNumberLiteralToStackInstruction(value)
	t.addInstruction(whitespace.StoreInHeap())
}

func (t *Transpiler) retrieveFromHeapInstruction(heapAddress byte) {
	t.pushNumberLiteralToStackInstruction(heapAddress)
	t.addInstruction(whitespace.RetrieveFromHeap())
}

func (t *Transpiler) pushNumberLiteralToStackInstruction(value byte) {
	t.addInstruction(
		whitespace.Instruction{
			Body: append(
				whitespace.PushToStack().Body,
				whitespace.NumberLiteral(value).Body...,
			),
		},
	)
}

func (t *Transpiler) printTopStackCharInstruction() {
	t.addInstruction(whitespace.PrintTopStack())
}
