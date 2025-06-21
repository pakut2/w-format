package jsWhitespaceTranspiler

import (
	"fmt"

	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/ast"
	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/object"
	"github.com/pakut2/w-format/pkg/whitespace"
)

type Transpiler struct {
	instructions []whitespace.Instruction

	currentHeapAddress int64
	currentLabelId     int64

	environment      *object.Environment
	builtInFunctions map[string]*object.BuiltIn
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

func (t *Transpiler) getEmptyHeapAddress() int64 {
	t.currentHeapAddress++

	return t.currentHeapAddress
}

func (t *Transpiler) getEmptyLabelId() int64 {
	t.currentLabelId++

	return t.currentLabelId
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
	case *ast.BlockStatement:
		return t.transpileBlockStatement(node)
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
	case *ast.PrefixExpression:
		rightExpression := t.Transpile(node.Right)

		return t.transpilePrefixExpression(node.Operator, rightExpression)
	case *ast.InfixExpression:
		leftExpression := t.Transpile(node.Left)
		rightExpression := t.Transpile(node.Right)

		return t.transpileInfixExpression(node.Operator, leftExpression, rightExpression)
	case *ast.IfExpression:
		return t.transpileIfExpression(node)
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

func (t *Transpiler) transpileBlockStatement(block *ast.BlockStatement) object.Object {
	for _, statement := range block.Statements {
		t.Transpile(statement)
	}

	return &object.Void{}
}

func (t *Transpiler) transpileIdentifier(identifier *ast.Identifier) object.Object {
	if val, ok := t.environment.Get(identifier.Value); ok {
		return val
	}

	if buildInFunction, ok := t.builtInFunctions[identifier.Value]; ok {
		return buildInFunction
	}

	panic(fmt.Sprintf("[:%d] identifier %s is not defined", identifier.Token.LineNumber, identifier.Value))
}

func (t *Transpiler) transpileString(value []byte) object.Object {
	var stringObject object.String

	for _, c := range value {
		heapAddress := t.getEmptyHeapAddress()
		t.storeValueInHeapInstruction(heapAddress, int64(c))

		stringObject.Chars = append(stringObject.Chars, object.Char{HeapAddress: t.currentHeapAddress})
	}

	return &stringObject
}

func (t *Transpiler) transpileInteger(value int64) object.Object {
	var integerObject object.Integer

	integerObject.HeapAddress = t.getEmptyHeapAddress()
	t.storeValueInHeapInstruction(integerObject.HeapAddress, value)

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

func (t *Transpiler) transpilePrefixExpression(operator string, rightExpression object.Object) object.Object {
	switch operator {
	case ast.SUBTRACTION:
		return t.transpileMinusPrefixOperatorExpression(rightExpression)
	case ast.NEGATION:
		return t.transpileNegationPrefixOperatorExpression(rightExpression)
	default:
		panic(fmt.Sprintf("unknown operator %s%s", operator, rightExpression.Type()))
	}
}

func (t *Transpiler) transpileMinusPrefixOperatorExpression(rightExpression object.Object) object.Object {
	if rightExpression.Type() != object.INT_OBJ {
		panic(fmt.Sprintf("unsupported opposition target %q", rightExpression.Type()))
	}

	rightInteger := rightExpression.(*object.Integer)

	t.literalMultiplicationInstruction(rightInteger.HeapAddress, -1)

	resultHeapAddress := t.getEmptyHeapAddress()
	t.storeTopStackValueInHeapInstruction(resultHeapAddress)

	rightInteger.HeapAddress = resultHeapAddress

	return rightInteger
}

func (t *Transpiler) transpileNegationPrefixOperatorExpression(rightExpression object.Object) object.Object {
	if rightExpression.Type() != object.INT_OBJ {
		panic(fmt.Sprintf("unsupported opposition target %q", rightExpression.Type()))
	}

	rightInteger := rightExpression.(*object.Integer)

	comparatorHeapAddress := t.getEmptyHeapAddress()
	t.storeValueInHeapInstruction(comparatorHeapAddress, 0)

	t.integerBooleanInstruction(ast.EQUALS, rightInteger.HeapAddress, comparatorHeapAddress)

	resultHeapAddress := t.getEmptyHeapAddress()
	t.storeTopStackValueInHeapInstruction(resultHeapAddress)

	rightInteger.HeapAddress = resultHeapAddress

	return rightInteger
}

func (t *Transpiler) transpileInfixExpression(operator string, leftExpression, rightExpression object.Object) object.Object {
	switch {
	case leftExpression.Type() == object.INT_OBJ && rightExpression.Type() == object.INT_OBJ:
		return t.transpileIntegerInfixExpression(operator, leftExpression, rightExpression)
	case leftExpression.Type() == object.STRING_OBJ && rightExpression.Type() == object.STRING_OBJ:
		panic("unsupported")
	case leftExpression.Type() != rightExpression.Type():
		panic(fmt.Sprintf("type mismatch %s %s %s", leftExpression.Type(), operator, rightExpression.Type()))
	default:
		panic(fmt.Sprintf("unknown operator %s %s %s", leftExpression.Type(), operator, rightExpression.Type()))
	}
}

func (t *Transpiler) transpileIntegerInfixExpression(operator string, leftExpression, rightExpression object.Object) object.Object {
	leftHeapAddress := leftExpression.(*object.Integer).HeapAddress
	rightHeapAddress := rightExpression.(*object.Integer).HeapAddress

	switch operator {
	case ast.ADDITION:
		t.additionInstruction(leftHeapAddress, rightHeapAddress)
	case ast.SUBTRACTION:
		t.subtractionInstruction(leftHeapAddress, rightHeapAddress)
	case ast.MULTIPLICATION:
		t.multiplicationInstruction(leftHeapAddress, rightHeapAddress)
	case ast.DIVISION:
		t.divisionInstruction(leftHeapAddress, rightHeapAddress)
	case ast.MODULO:
		t.moduloInstruction(leftHeapAddress, rightHeapAddress)
	case ast.EQUALS, ast.NOT_EQUALS, ast.LESS_THAN, ast.LESS_THAN_OR_EQUAL, ast.GREATER_THAN, ast.GREATER_THAN_OR_EQUAL:
		t.integerBooleanInstruction(operator, leftHeapAddress, rightHeapAddress)
	default:
		panic(fmt.Sprintf("unknown operator %s %s %s", leftExpression.Type(), operator, rightExpression.Type()))
	}

	resultHeapAddress := t.getEmptyHeapAddress()
	t.storeTopStackValueInHeapInstruction(resultHeapAddress)

	return &object.Integer{HeapAddress: resultHeapAddress}
}

func (t *Transpiler) transpileIfExpression(ifExpression *ast.IfExpression) object.Object {
	alternativeLabel := t.getEmptyLabelId()
	endIfLabel := t.getEmptyLabelId()

	conditionResult := t.Transpile(ifExpression.Condition)

	conditionResultLiteral, ok := conditionResult.(*object.Integer)
	if !ok {
		panic("invalid if condition expression")
	}

	t.retrieveFromHeapInstruction(conditionResultLiteral.HeapAddress)
	t.addInstruction(whitespace.JumpToLabelIfZero(alternativeLabel))

	t.Transpile(ifExpression.Consequence)
	t.addInstruction(whitespace.JumpToLabel(endIfLabel))

	t.addInstruction(whitespace.Label(alternativeLabel))

	if ifExpression.Alternative != nil {
		t.Transpile(ifExpression.Alternative)
	}

	t.addInstruction(whitespace.Label(endIfLabel))

	return &object.Void{}
}

func (t *Transpiler) storeTopStackValueInHeapInstruction(heapAddress int64) {
	t.pushNumberLiteralToStackInstruction(heapAddress)
	t.addInstruction(whitespace.SwapTwoTopStackItems())
	t.addInstruction(whitespace.StoreInHeap())
}

func (t *Transpiler) storeValueInHeapInstruction(heapAddress int64, value int64) {
	t.pushNumberLiteralToStackInstruction(heapAddress)
	t.pushNumberLiteralToStackInstruction(value)
	t.addInstruction(whitespace.StoreInHeap())
}

func (t *Transpiler) retrieveFromHeapInstruction(heapAddress int64) {
	t.pushNumberLiteralToStackInstruction(heapAddress)
	t.addInstruction(whitespace.RetrieveFromHeap())
}

func (t *Transpiler) retrieveMultipleFromHeapInstruction(heapAddress1 int64, heapAddress2 int64) {
	t.retrieveFromHeapInstruction(heapAddress1)
	t.retrieveFromHeapInstruction(heapAddress2)
	t.addInstruction(whitespace.LiftStackItem(1))
	t.addInstruction(whitespace.SwapTwoTopStackItems())
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

func (t *Transpiler) additionInstruction(heapAddress1 int64, heapAddress2 int64) {
	t.retrieveMultipleFromHeapInstruction(heapAddress1, heapAddress2)
	t.addInstruction(whitespace.Add())
}

func (t *Transpiler) subtractionInstruction(heapAddress1 int64, heapAddress2 int64) {
	t.retrieveMultipleFromHeapInstruction(heapAddress1, heapAddress2)
	t.addInstruction(whitespace.Subtract())
}

func (t *Transpiler) multiplicationInstruction(heapAddress1 int64, heapAddress2 int64) {
	t.retrieveMultipleFromHeapInstruction(heapAddress1, heapAddress2)
	t.addInstruction(whitespace.Multiply())
}

func (t *Transpiler) literalMultiplicationInstruction(heapAddress int64, value int64) {
	t.retrieveFromHeapInstruction(heapAddress)
	t.pushNumberLiteralToStackInstruction(value)
	t.addInstruction(whitespace.Multiply())
}

func (t *Transpiler) literalStackMultiplicationInstruction(value int64) {
	t.pushNumberLiteralToStackInstruction(value)
	t.addInstruction(whitespace.Multiply())
}

func (t *Transpiler) divisionInstruction(heapAddress1 int64, heapAddress2 int64) {
	t.retrieveMultipleFromHeapInstruction(heapAddress1, heapAddress2)
	t.addInstruction(whitespace.Divide())
}

func (t *Transpiler) moduloInstruction(heapAddress1 int64, heapAddress2 int64) {
	t.retrieveMultipleFromHeapInstruction(heapAddress1, heapAddress2)
	t.addInstruction(whitespace.Mod())
}

func (t *Transpiler) integerBooleanInstruction(operator string, leftHeapAddress, rightHeapAddress int64) {
	consequenceLabel := t.getEmptyLabelId()
	endIfLabel := t.getEmptyLabelId()

	if operator == ast.EQUALS || operator == ast.NOT_EQUALS || operator == ast.GREATER_THAN || operator == ast.LESS_THAN_OR_EQUAL {
		t.subtractionInstruction(leftHeapAddress, rightHeapAddress)
		t.addInstruction(whitespace.JumpToLabelIfZero(consequenceLabel))
	}

	if operator == ast.LESS_THAN || operator == ast.GREATER_THAN || operator == ast.LESS_THAN_OR_EQUAL || operator == ast.GREATER_THAN_OR_EQUAL {
		t.subtractionInstruction(leftHeapAddress, rightHeapAddress)
		t.addInstruction(whitespace.JumpToLabelIfNegative(consequenceLabel))
	}

	switch operator {
	case ast.EQUALS, ast.LESS_THAN, ast.LESS_THAN_OR_EQUAL:
		t.pushNumberLiteralToStackInstruction(whitespace.FALSE)
	case ast.NOT_EQUALS, ast.GREATER_THAN, ast.GREATER_THAN_OR_EQUAL:
		t.pushNumberLiteralToStackInstruction(whitespace.TRUE)
	default:
		panic(fmt.Sprintf("unknown operator %q", operator))
	}

	t.addInstruction(whitespace.JumpToLabel(endIfLabel))

	t.addInstruction(whitespace.Label(consequenceLabel))

	switch operator {
	case ast.EQUALS, ast.LESS_THAN, ast.LESS_THAN_OR_EQUAL:
		t.pushNumberLiteralToStackInstruction(whitespace.TRUE)
	case ast.NOT_EQUALS, ast.GREATER_THAN, ast.GREATER_THAN_OR_EQUAL:
		t.pushNumberLiteralToStackInstruction(whitespace.FALSE)
	default:
		panic(fmt.Sprintf("unknown operator %q", operator))
	}

	t.addInstruction(whitespace.Label(endIfLabel))
}
