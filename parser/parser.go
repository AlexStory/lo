package parser

import (
	"lo/ast"
	"lo/lexer"
	"lo/token"
	"strconv"
	"strings"
)

type ParseError struct {
	Msg    string
	Line   int
	Column int
}

type Parser struct {
	l      *lexer.Lexer
	Errors []ParseError

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Expressions = []ast.Expression{}

	for p.curToken.Type != token.EOF {
		expr := p.parseExpression()
		if expr != nil {
			program.Expressions = append(program.Expressions, expr)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.curToken.Type {
	case token.Ident:
		return p.parseIdentifier()
	case token.Number:
		return p.parseNumber()
	case token.String:
		return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case token.OpenParen:
		return p.parseList()
	case token.OpenBracket:
		return p.parseListLiteral()
	default:
		return nil
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseNumber() ast.Expression {
	if strings.Contains(p.curToken.Literal, ".") {
		value, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			p.Errors = append(p.Errors, ParseError{Msg: "Could not parse float", Line: p.curToken.Line, Column: p.curToken.Column})
			return nil
		}
		return &ast.FloatLiteral{Token: p.curToken, Value: float64(value)}
	} else {
		value, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			p.Errors = append(p.Errors, ParseError{Msg: "Could not parse int", Line: p.curToken.Line, Column: p.curToken.Column})
			return nil
		}
		return &ast.IntLiteral{Token: p.curToken, Value: int64(value)}
	}
}

func (p *Parser) parseList() *ast.ListExpression {
	list := &ast.ListExpression{Token: p.curToken}
	list.Expressions = []ast.Expression{}

	for p.peekToken.Type != token.CloseParen && p.peekToken.Type != token.EOF {
		p.nextToken()
		expr := p.parseExpression()
		if expr != nil {
			list.Expressions = append(list.Expressions, expr)
		}
	}
	p.nextToken()

	return list
}

func (p *Parser) parseListLiteral() *ast.ListLiteral {
	list := &ast.ListLiteral{Token: p.curToken}
	list.Expressions = []ast.Expression{}

	p.nextToken() // Skip '['

	for !p.curTokenIs(token.CloseBracket) && !p.curTokenIs(token.EOF) {
		expr := p.parseExpression()
		if expr != nil {
			list.Expressions = append(list.Expressions, expr)
		}
		p.nextToken()
	}

	return list
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
