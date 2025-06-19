package jsWhitespaceTranspiler

import (
	"fmt"

	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceTranspiler/internal/ast"
	"github.com/pakut2/js-whitespace/pkg/jsWhitespaceTranspiler/internal/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

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

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.LEFT_PARENTHESIS, p.parseCallExpression)

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
		statement := p.parseExpressionStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{}
	statement.Expression = p.parseExpression()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpression() ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		panic(fmt.Sprintf("[:%d] no prefix parse function for %s found", p.currentToken.LineNumber, p.currentToken.Type))

		return nil
	}

	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		p.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
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
	args = append(args, p.parseExpression())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		args = append(args, p.parseExpression())
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
