package jsWhitespaceTranspiler

import (
	"fmt"

	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/ast"
	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/object"
	"github.com/pakut2/w-format/pkg/whitespace"
)

type Transpiler struct {
	instructions       []whitespace.Instruction
	currentHeapAddress int64
	environment        *object.Environment
	builtInFunctions   map[string]*object.BuiltIn
}

func NewTranspiler() *Transpiler {
	t := &Transpiler{
		currentHeapAddress: 0,
		environment:        object.NewEnvironment(),
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
		case *object.Integer:
			t.retrieveFromHeapInstruction(arg.HeapAddress)
			t.printTopStackIntegerInstruction()

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

func (t *Transpiler) getCurrentHeapAddressWithIncrement() int64 {
	t.currentHeapAddress++

	return t.currentHeapAddress
}

func (t *Transpiler) Transpile(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return t.transpileProgram(node)
	case *ast.LetStatement:
		value := t.Transpile(node.Value)
		t.environment.Set(node.Name.Value, value)

		return &object.Void{}
	case *ast.ExpressionStatement:
		return t.Transpile(node.Expression)
	case *ast.StringLiteral:
		return t.transpileString([]byte(node.Value))
	case *ast.IntegerLiteral:
		return t.transpileInteger(node.Value)
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
		WhitespaceInstructions: t.instructions,
	}
}

func (t *Transpiler) transpileIdentifier(identifier *ast.Identifier) object.Object {
	if val, ok := t.environment.Get(identifier.Value); ok {
		return val
	}

	if buildInFunction, ok := t.builtInFunctions[identifier.Value]; ok {
		return buildInFunction
	}

	panic(fmt.Sprintf("[:%d] identifier %s not implemeted", identifier.Token.LineNumber, identifier.Value))
}

func (t *Transpiler) transpileString(value []byte) object.Object {
	var stringObject object.String

	for _, c := range value {
		heapAddress := t.getCurrentHeapAddressWithIncrement()
		t.storeInHeapInstruction(heapAddress, int64(c))

		stringObject.Chars = append(stringObject.Chars, object.Char{HeapAddress: t.currentHeapAddress})
	}

	return &stringObject
}

func (t *Transpiler) transpileInteger(value int64) object.Object {
	var integerObject object.Integer

	integerObject.HeapAddress = t.getCurrentHeapAddressWithIncrement()
	t.storeInHeapInstruction(integerObject.HeapAddress, value)

	return &integerObject
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

func (t *Transpiler) storeInHeapInstruction(heapAddress int64, value int64) {
	t.pushNumberLiteralToStackInstruction(heapAddress)
	t.pushNumberLiteralToStackInstruction(value)
	t.addInstruction(whitespace.StoreInHeap())
}

func (t *Transpiler) retrieveFromHeapInstruction(heapAddress int64) {
	t.pushNumberLiteralToStackInstruction(heapAddress)
	t.addInstruction(whitespace.RetrieveFromHeap())
}

func (t *Transpiler) pushNumberLiteralToStackInstruction(value int64) {
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
	t.addInstruction(whitespace.PrintTopStackChar())
}

func (t *Transpiler) printTopStackIntegerInstruction() {
	t.addInstruction(whitespace.PrintTopStackInteger())
}
