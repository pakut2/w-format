package ast

import "github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/token"

const (
	ADDITION              = token.PLUS
	SUBTRACTION           = token.MINUS
	MULTIPLICATION        = token.ASTERISK
	DIVISION              = token.SLASH
	MODULO                = token.PERCENT
	NEGATION              = token.BANG
	EQUALS                = token.EQUALS
	NOT_EQUALS            = token.NOT_EQUALS
	LESS_THAN             = token.LESS_THAN
	LESS_THAN_OR_EQUAL    = token.LESS_THAN_OR_EQUAL
	GREATER_THAN          = token.GREATER_THAN
	GREATER_THAN_OR_EQUAL = token.GREATER_THAN_OR_EQUAL
)

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

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}

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

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode() {}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}
