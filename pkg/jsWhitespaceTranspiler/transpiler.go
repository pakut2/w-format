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

func (t *Transpiler) TranspileProgram(program *ast.Program) object.Object {
	for _, statement := range program.Statements {
		t.transpile(statement, nil)
	}

	t.addInstruction(whitespace.EndProgram())

	return &object.Program{
		WhitespaceInstructions: t.instructions,
	}
}

func (t *Transpiler) transpile(node ast.Node, scopeContext *object.ScopeContext) object.Object {
	switch node := node.(type) {
	case *ast.LetStatement:
		return t.transpileLetStatement(node)
	case *ast.AssignmentStatement:
		return t.transpileAssignmentStatement(node)
	case *ast.IfStatement:
		return t.transpileIfStatement(node, scopeContext)
	case *ast.BlockStatement:
		return t.transpileBlockStatement(node, scopeContext)
	case *ast.ForStatement:
		return t.transpileForStatement(node)
	case *ast.BreakStatement:
		return t.transpileBreakStatement(node, scopeContext)
	case *ast.ContinueStatement:
		return t.transpileContinueStatement(node, scopeContext)
	case *ast.ExpressionStatement:
		return t.transpile(node.Expression, scopeContext)
	case *ast.StringLiteral:
		return t.transpileString([]byte(node.Value))
	case *ast.IntegerLiteral:
		return t.transpileInteger(node.Value)
	case *ast.Identifier:
		return t.transpileIdentifier(node)
	case *ast.CallExpression:
		function := t.transpile(node.Function, scopeContext)
		args := t.transpileExpressions(node.Arguments)

		return t.applyFunction(function, args)
	case *ast.PrefixExpression:
		right := t.transpile(node.Right, scopeContext)

		return t.transpilePrefixExpression(node, right)
	case *ast.InfixExpression:
		left := t.transpile(node.Left, scopeContext)
		right := t.transpile(node.Right, scopeContext)

		return t.transpileInfixExpression(node, left, right)
	case *ast.SuffixExpression:
		left := t.transpile(node.Left, scopeContext)

		return t.transpileSuffixExpression(node, left)
	}

	return &object.Void{}
}

func (t *Transpiler) transpileLetStatement(statement *ast.LetStatement) object.Object {
	_, ok := t.environment.Get(statement.Name.Value)
	if ok {
		panic(fmt.Sprintf("[:%d] redeclaration of %s", statement.Token.LineNumber, statement.Name.Value))
	}

	value := t.transpile(statement.Value, nil)
	t.environment.Set(statement.Name.Value, value)

	return &object.Void{}
}

func (t *Transpiler) transpileAssignmentStatement(statement *ast.AssignmentStatement) object.Object {
	previousValue, ok := t.environment.Get(statement.Name.Value)
	if !ok {
		panic(fmt.Sprintf("[:%d] %s is not defined", statement.Token.LineNumber, statement.Name.Value))
	}

	assignedValue := t.transpile(statement.Value, nil)

	if previousValue.Type() != assignedValue.Type() {
		panic(
			fmt.Sprintf("[:%d] assignment type mismatch %s = %s",
				statement.Token.LineNumber,
				previousValue.Type(),
				assignedValue.Type(),
			),
		)
	}

	if previousValue.Type() == object.INT_OBJ {
		t.retrieveFromHeapInstruction(assignedValue.(*object.Integer).HeapAddress)
		t.storeTopStackValueInHeapInstruction(previousValue.(*object.Integer).HeapAddress)
	}

	return &object.Void{}
}

func (t *Transpiler) transpileIfStatement(statement *ast.IfStatement, scopeContext *object.ScopeContext) object.Object {
	alternativeLabel := t.getEmptyLabelId()
	endIfLabel := t.getEmptyLabelId()

	conditionResult := t.transpile(statement.Condition, nil)
	if conditionResult.Type() != object.INT_OBJ {
		panic(fmt.Sprintf("[:%d] invalid if condition expression", statement.Token.LineNumber))
	}

	t.retrieveFromHeapInstruction(conditionResult.(*object.Integer).HeapAddress)
	t.addInstruction(whitespace.JumpToLabelIfZero(alternativeLabel))

	t.transpile(statement.Consequence, scopeContext)
	t.addInstruction(whitespace.JumpToLabel(endIfLabel))

	t.addInstruction(whitespace.Label(alternativeLabel))

	if statement.Alternative != nil {
		t.transpile(statement.Alternative, scopeContext)
	}

	t.addInstruction(whitespace.Label(endIfLabel))

	return &object.Void{}
}

func (t *Transpiler) transpileBlockStatement(block *ast.BlockStatement, scopeContext *object.ScopeContext) object.Object {
	for _, statement := range block.Statements {
		t.transpile(statement, scopeContext)
	}

	return &object.Void{}
}

func (t *Transpiler) transpileForStatement(statement *ast.ForStatement) object.Object {
	loopControlLabel := t.getEmptyLabelId()
	loopBodyLabel := t.getEmptyLabelId()
	loopEndLabel := t.getEmptyLabelId()

	t.transpile(statement.Declaration, nil)

	initialConditionResult := t.transpile(statement.Boundary, nil)
	if initialConditionResult.Type() != object.INT_OBJ {
		panic(fmt.Sprintf("[:%d] invalid loop boundary condition expression", statement.Token.LineNumber))
	}

	t.retrieveFromHeapInstruction(initialConditionResult.(*object.Integer).HeapAddress)
	t.addInstruction(whitespace.JumpToLabelIfZero(loopEndLabel))

	t.addInstruction(whitespace.JumpToLabel(loopBodyLabel))

	t.addInstruction(whitespace.Label(loopControlLabel))

	incrementResult := t.transpile(statement.Increment, nil)
	if incrementResult.Type() != object.INT_OBJ {
		panic(fmt.Sprintf("[:%d] invalid loop increment expression", statement.Token.LineNumber))
	}

	previousIteratorValue, _ := t.environment.Get(statement.Declaration.Name.Value)
	previousIteratorValueHeapAddress := previousIteratorValue.(*object.Integer).HeapAddress

	t.retrieveFromHeapInstruction(incrementResult.(*object.Integer).HeapAddress)
	t.storeTopStackValueInHeapInstruction(previousIteratorValueHeapAddress)

	conditionResult := t.transpile(statement.Boundary, nil)
	if conditionResult.Type() != object.INT_OBJ {
		panic(fmt.Sprintf("[:%d] invalid loop boundary condition expression", statement.Token.LineNumber))
	}

	t.retrieveFromHeapInstruction(conditionResult.(*object.Integer).HeapAddress)
	t.addInstruction(whitespace.JumpToLabelIfZero(loopEndLabel))

	t.addInstruction(whitespace.Label(loopBodyLabel))

	t.transpile(statement.Body, &object.ScopeContext{
		For: &object.ForContext{
			ControlLabelId: loopControlLabel,
			EndLabelId:     loopEndLabel,
		},
	})

	t.addInstruction(whitespace.JumpToLabel(loopControlLabel))

	t.addInstruction(whitespace.Label(loopEndLabel))

	return &object.Void{}
}

func (t *Transpiler) transpileBreakStatement(statement *ast.BreakStatement, scopeContext *object.ScopeContext) object.Object {
	if scopeContext == nil {
		panic(fmt.Sprintf("[:%d] cannot determine break target", statement.Token.LineNumber))
	}

	t.addInstruction(whitespace.JumpToLabel(scopeContext.For.EndLabelId))

	return &object.Void{}
}

func (t *Transpiler) transpileContinueStatement(statement *ast.ContinueStatement, scopeContext *object.ScopeContext) object.Object {
	if scopeContext == nil {
		panic(fmt.Sprintf("[:%d] cannot determine continue target", statement.Token.LineNumber))
	}

	t.addInstruction(whitespace.JumpToLabel(scopeContext.For.ControlLabelId))

	return &object.Void{}
}

func (t *Transpiler) transpileIdentifier(identifier *ast.Identifier) object.Object {
	if val, ok := t.environment.Get(identifier.Value); ok {
		return val
	}

	if buildInFunction, ok := t.builtInFunctions[identifier.Value]; ok {
		return buildInFunction
	}

	panic(fmt.Sprintf("[:%d] %s is not defined", identifier.Token.LineNumber, identifier.Value))
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
		transpiledExpression := t.transpile(expression, nil)

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

func (t *Transpiler) transpilePrefixExpression(expression *ast.PrefixExpression, right object.Object) object.Object {
	switch expression.Operator {
	case ast.SUBTRACTION:
		return t.transpileMinusPrefixOperatorExpression(right)
	case ast.NEGATION:
		return t.transpileNegationPrefixOperatorExpression(right)
	default:
		panic(fmt.Sprintf("[:%d] unknown operator %s%s", expression.Token.LineNumber, expression.Operator, right.Type()))
	}
}

func (t *Transpiler) transpileMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INT_OBJ {
		panic(fmt.Sprintf("unsupported opposition target %q", right.Type()))
	}

	rightInteger := right.(*object.Integer)

	t.literalMultiplicationInstruction(rightInteger.HeapAddress, -1)

	resultHeapAddress := t.getEmptyHeapAddress()
	t.storeTopStackValueInHeapInstruction(resultHeapAddress)

	rightInteger.HeapAddress = resultHeapAddress

	return rightInteger
}

func (t *Transpiler) transpileNegationPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INT_OBJ {
		panic(fmt.Sprintf("unsupported opposition target %q", right.Type()))
	}

	rightInteger := right.(*object.Integer)

	comparatorHeapAddress := t.getEmptyHeapAddress()
	t.storeValueInHeapInstruction(comparatorHeapAddress, 0)

	t.integerComparisonInstruction(ast.EQUALS, rightInteger.HeapAddress, comparatorHeapAddress)

	resultHeapAddress := t.getEmptyHeapAddress()
	t.storeTopStackValueInHeapInstruction(resultHeapAddress)

	rightInteger.HeapAddress = resultHeapAddress

	return rightInteger
}

func (t *Transpiler) transpileInfixExpression(expression *ast.InfixExpression, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ:
		return t.transpileIntegerInfixExpression(expression, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return t.transpileStringInfixExpression(expression, left, right)
	case left.Type() != right.Type():
		panic(
			fmt.Sprintf("[:%d] type mismatch %s %s %s",
				expression.Token.LineNumber,
				left.Type(),
				expression.Operator,
				right.Type(),
			),
		)
	default:
		panic(
			fmt.Sprintf("[:%d] unknown operator %s %s %s",
				expression.Token.LineNumber,
				left.Type(),
				expression.Operator,
				right.Type(),
			),
		)
	}
}

func (t *Transpiler) transpileIntegerInfixExpression(expression *ast.InfixExpression, left, right object.Object) object.Object {
	leftHeapAddress := left.(*object.Integer).HeapAddress
	rightHeapAddress := right.(*object.Integer).HeapAddress

	switch expression.Operator {
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
		t.integerComparisonInstruction(expression.Operator, leftHeapAddress, rightHeapAddress)
	case ast.AND:
		t.integerAndInstruction(leftHeapAddress, rightHeapAddress)
	case ast.OR:
		t.integerOrInstruction(leftHeapAddress, rightHeapAddress)
	default:
		panic(
			fmt.Sprintf("[:%d] unknown operator %s %s %s",
				expression.Token.LineNumber,
				left.Type(),
				expression.Operator,
				right.Type(),
			),
		)
	}

	resultHeapAddress := t.getEmptyHeapAddress()
	t.storeTopStackValueInHeapInstruction(resultHeapAddress)

	return &object.Integer{HeapAddress: resultHeapAddress}
}

func (t *Transpiler) transpileStringInfixExpression(expression *ast.InfixExpression, left, right object.Object) object.Object {
	leftChars := left.(*object.String).Chars
	rightChars := right.(*object.String).Chars

	switch expression.Operator {
	case ast.ADDITION:
		return &object.String{Chars: append(leftChars, rightChars...)} // This is naive and won't work for runtime assignment statements without dynamic memory allocation
	default:
		panic(
			fmt.Sprintf("[:%d] unknown operator %s %s %s",
				expression.Token.LineNumber,
				left.Type(),
				expression.Operator,
				right.Type(),
			),
		)
	}
}

func (t *Transpiler) transpileSuffixExpression(expression *ast.SuffixExpression, operand object.Object) object.Object {
	if operand.Type() != object.INT_OBJ {
		panic(
			fmt.Sprintf("[:%d] unsupported %s target %q",
				expression.Token.LineNumber,
				expression.Operator,
				operand.Type()),
		)
	}

	switch expression.Operator {
	case ast.INCREMENT:
		t.literalAdditionInstruction(operand.(*object.Integer).HeapAddress, 1)
	case ast.DECREMENT:
		t.literalSubtractionInstruction(operand.(*object.Integer).HeapAddress, 1)
	default:
		panic(
			fmt.Sprintf("[:%d] unknown operator %s %s",
				expression.Token.LineNumber,
				operand.Type(),
				expression.Operator,
			),
		)
	}

	resultHeapAddress := t.getEmptyHeapAddress()
	t.storeTopStackValueInHeapInstruction(resultHeapAddress)

	return &object.Integer{HeapAddress: resultHeapAddress}
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

func (t *Transpiler) literalAdditionInstruction(heapAddress, value int64) {
	t.retrieveFromHeapInstruction(heapAddress)
	t.pushNumberLiteralToStackInstruction(value)
	t.addInstruction(whitespace.Add())
}

func (t *Transpiler) subtractionInstruction(heapAddress1, heapAddress2 int64) {
	t.retrieveMultipleFromHeapInstruction(heapAddress1, heapAddress2)
	t.addInstruction(whitespace.Subtract())
}

func (t *Transpiler) literalSubtractionInstruction(heapAddress, value int64) {
	t.retrieveFromHeapInstruction(heapAddress)
	t.pushNumberLiteralToStackInstruction(value)
	t.addInstruction(whitespace.Subtract())
}

func (t *Transpiler) multiplicationInstruction(heapAddress1, heapAddress2 int64) {
	t.retrieveMultipleFromHeapInstruction(heapAddress1, heapAddress2)
	t.addInstruction(whitespace.Multiply())
}

func (t *Transpiler) literalMultiplicationInstruction(heapAddress, value int64) {
	t.retrieveFromHeapInstruction(heapAddress)
	t.pushNumberLiteralToStackInstruction(value)
	t.addInstruction(whitespace.Multiply())
}

func (t *Transpiler) literalStackMultiplicationInstruction(value int64) {
	t.pushNumberLiteralToStackInstruction(value)
	t.addInstruction(whitespace.Multiply())
}

func (t *Transpiler) divisionInstruction(heapAddress1, heapAddress2 int64) {
	t.retrieveMultipleFromHeapInstruction(heapAddress1, heapAddress2)
	t.addInstruction(whitespace.Divide())
}

func (t *Transpiler) moduloInstruction(heapAddress1, heapAddress2 int64) {
	t.retrieveMultipleFromHeapInstruction(heapAddress1, heapAddress2)
	t.addInstruction(whitespace.Mod())
}

func (t *Transpiler) integerComparisonInstruction(operator string, leftHeapAddress, rightHeapAddress int64) {
	matchLabel := t.getEmptyLabelId()
	endComparisonLabel := t.getEmptyLabelId()

	if operator == ast.EQUALS || operator == ast.NOT_EQUALS || operator == ast.GREATER_THAN || operator == ast.LESS_THAN_OR_EQUAL {
		t.subtractionInstruction(leftHeapAddress, rightHeapAddress)
		t.addInstruction(whitespace.JumpToLabelIfZero(matchLabel))
	}

	if operator == ast.LESS_THAN || operator == ast.GREATER_THAN || operator == ast.LESS_THAN_OR_EQUAL || operator == ast.GREATER_THAN_OR_EQUAL {
		t.subtractionInstruction(leftHeapAddress, rightHeapAddress)
		t.addInstruction(whitespace.JumpToLabelIfNegative(matchLabel))
	}

	switch operator {
	case ast.EQUALS, ast.LESS_THAN, ast.LESS_THAN_OR_EQUAL:
		t.pushNumberLiteralToStackInstruction(whitespace.FALSE)
	case ast.NOT_EQUALS, ast.GREATER_THAN, ast.GREATER_THAN_OR_EQUAL:
		t.pushNumberLiteralToStackInstruction(whitespace.TRUE)
	default:
		panic(fmt.Sprintf("unknown operator %q", operator))
	}

	t.addInstruction(whitespace.JumpToLabel(endComparisonLabel))

	t.addInstruction(whitespace.Label(matchLabel))

	switch operator {
	case ast.EQUALS, ast.LESS_THAN, ast.LESS_THAN_OR_EQUAL:
		t.pushNumberLiteralToStackInstruction(whitespace.TRUE)
	case ast.NOT_EQUALS, ast.GREATER_THAN, ast.GREATER_THAN_OR_EQUAL:
		t.pushNumberLiteralToStackInstruction(whitespace.FALSE)
	default:
		panic(fmt.Sprintf("unknown operator %q", operator))
	}

	t.addInstruction(whitespace.Label(endComparisonLabel))
}

func (t *Transpiler) integerAndInstruction(leftHeapAddress, rightHeapAddress int64) {
	firstMatchLabel := t.getEmptyLabelId()
	secondMatchLabel := t.getEmptyLabelId()
	endComparisonLabel := t.getEmptyLabelId()

	comparatorHeapAddress := t.getEmptyHeapAddress()
	t.storeValueInHeapInstruction(comparatorHeapAddress, 0)

	t.integerComparisonInstruction(ast.EQUALS, leftHeapAddress, comparatorHeapAddress)
	t.addInstruction(whitespace.JumpToLabelIfZero(firstMatchLabel))

	t.pushNumberLiteralToStackInstruction(whitespace.FALSE)
	t.addInstruction(whitespace.JumpToLabel(endComparisonLabel))

	t.addInstruction(whitespace.Label(firstMatchLabel))

	firstComparisonResultHeapAddress := t.getEmptyHeapAddress()
	t.storeTopStackValueInHeapInstruction(firstComparisonResultHeapAddress)

	t.integerComparisonInstruction(ast.EQUALS, rightHeapAddress, comparatorHeapAddress)
	t.addInstruction(whitespace.JumpToLabelIfZero(secondMatchLabel))

	t.pushNumberLiteralToStackInstruction(whitespace.FALSE)
	t.addInstruction(whitespace.JumpToLabel(endComparisonLabel))

	t.addInstruction(whitespace.Label(secondMatchLabel))
	t.pushNumberLiteralToStackInstruction(whitespace.TRUE)

	t.addInstruction(whitespace.Label(endComparisonLabel))
}

func (t *Transpiler) integerOrInstruction(leftHeapAddress, rightHeapAddress int64) {
	matchLabel := t.getEmptyLabelId()
	endComparisonLabel := t.getEmptyLabelId()

	comparatorHeapAddress := t.getEmptyHeapAddress()
	t.storeValueInHeapInstruction(comparatorHeapAddress, 0)

	t.integerComparisonInstruction(ast.EQUALS, leftHeapAddress, comparatorHeapAddress)
	t.addInstruction(whitespace.JumpToLabelIfZero(matchLabel))

	firstComparisonResultHeapAddress := t.getEmptyHeapAddress()
	t.storeTopStackValueInHeapInstruction(firstComparisonResultHeapAddress)

	t.integerComparisonInstruction(ast.EQUALS, rightHeapAddress, comparatorHeapAddress)
	t.addInstruction(whitespace.JumpToLabelIfZero(matchLabel))

	t.pushNumberLiteralToStackInstruction(whitespace.FALSE)
	t.addInstruction(whitespace.JumpToLabel(endComparisonLabel))

	t.addInstruction(whitespace.Label(matchLabel))
	t.pushNumberLiteralToStackInstruction(whitespace.TRUE)

	t.addInstruction(whitespace.Label(endComparisonLabel))
}
