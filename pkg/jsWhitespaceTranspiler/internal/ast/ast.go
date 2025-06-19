package ast

import "github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/token"

type Node interface{}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
