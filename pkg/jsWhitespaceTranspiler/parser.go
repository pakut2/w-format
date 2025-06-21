package jsWhitespaceTranspiler

import (
	"fmt"
	"strconv"

	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/ast"
	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler/internal/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESS_GREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQUALS:                EQUALS,
	token.NOT_EQUALS:            EQUALS,
	token.LESS_THAN:             LESS_GREATER,
	token.LESS_THAN_OR_EQUAL:    LESS_GREATER,
	token.GREATER_THAN:          LESS_GREATER,
	token.GREATER_THAN_OR_EQUAL: LESS_GREATER,
	token.PLUS:                  SUM,
	token.MINUS:                 SUM,
	token.SLASH:                 PRODUCT,
	token.ASTERISK:              PRODUCT,
	token.PERCENT:               PRODUCT,
	token.LEFT_PARENTHESIS:      CALL,
}

type Parser struct {
	l      *Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LEFT_PARENTHESIS, p.parseGroupedExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.LEFT_PARENTHESIS, p.parseCallExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.EQUALS, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUALS, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN_OR_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN_OR_EQUAL, p.parseInfixExpression)

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(TokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[TokenType] = fn
}

func (p *Parser) registerInfix(TokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[TokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		panic(fmt.Sprintf("[:%d] invalid declaration statement, identifier must follow %q", p.currentToken.LineNumber, p.currentToken.Type))
	}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		panic(fmt.Sprintf("[:%d] invalid declaration statement, assignment must follow %q", p.currentToken.LineNumber, p.currentToken.Type))
	}

	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	expressionLineNumber := p.currentToken.LineNumber

	for !p.currentTokenIs(token.SEMICOLON) {
		// TODO check end of line
		if p.currentTokenIs(token.EOF) {
			panic(fmt.Sprintf("[:%d] missing semicolon at the end of let statement", expressionLineNumber))
		}

		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{}
	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		panic(fmt.Sprintf("[:%d] no prefix parse function for %s found", p.currentToken.LineNumber, p.currentToken.Type))

		return nil
	}

	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		p.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()

	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		panic(fmt.Sprintf("[:%d] invalid expression grouping", p.currentToken.LineNumber))
	}

	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Function: function}
	expression.Arguments = p.parseCallArguments()

	return expression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(token.RIGHT_PARENTHESIS) {
		p.nextToken()

		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	return args
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	integerLiteral := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		panic(fmt.Sprintf("[:%d] cannot parse %q as integer", p.currentToken.LineNumber, p.currentToken.Literal))
	}

	integerLiteral.Value = value

	return integerLiteral
}

func (p *Parser) parseBoolean() ast.Expression {
	var integerBooleanValue int64

	if p.currentTokenIs(token.TRUE) {
		integerBooleanValue = 1
	} else {
		integerBooleanValue = 0
	}

	return &ast.IntegerLiteral{Token: p.currentToken, Value: integerBooleanValue}
}

func (p *Parser) currentTokenIs(expectedCurrentToken token.TokenType) bool {
	return p.currentToken.Type == expectedCurrentToken
}

func (p *Parser) peekTokenIs(expectedPeekToken token.TokenType) bool {
	return p.peekToken.Type == expectedPeekToken
}

func (p *Parser) expectPeek(expectedPeekToken token.TokenType) bool {
	if p.peekTokenIs(expectedPeekToken) {
		p.nextToken()

		return true
	}

	panic(
		fmt.Sprintf(
			"[:%d] expected next token to be %s, got %s instead",
			p.currentToken.LineNumber,
			expectedPeekToken,
			p.peekToken.Type,
		),
	)
}

func (p *Parser) currentPrecedence() int {
	if precedence, ok := precedences[p.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}
