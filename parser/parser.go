package parser

import (
	"fmt"
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
	case token.OpenBrace:
		return p.parseMapLiteral()
	case token.HashParen:
		return p.parseLambdaLiteral()
	case token.Quote:
		return p.parseQuoteExpression()
	case token.Backtick:
		return p.parseBacktickExpression()
	case token.Tilde:
		return p.parseTildeExpression()
	case token.TildeAt:
		return p.parseTildeAtExpression()
	case token.Keyword:
		return &ast.Keyword{Token: p.curToken, Value: p.curToken.Literal}
	default:
		return nil
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseNumber() ast.Expression {
	if strings.Contains(p.curToken.Literal, ".") {
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			fmt.Println(value)
			fmt.Println(err)
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

func (p *Parser) parseMapLiteral() *ast.MapLiteral {
	mapLiteral := &ast.MapLiteral{Token: p.curToken}
	mapLiteral.Pairs = []ast.MapPair{}

	p.nextToken() // Skip '{'

	for !p.curTokenIs(token.CloseBrace) && !p.curTokenIs(token.EOF) {
		key := p.parseExpression()
		if key == nil {
			return nil
		}

		p.nextToken()
		value := p.parseExpression()
		if value == nil {
			return nil
		}

		mapLiteral.Pairs = append(mapLiteral.Pairs, ast.MapPair{Key: key, Value: value})
		p.nextToken()
	}

	return mapLiteral
}

func (p *Parser) parseLambdaLiteral() *ast.LambdaLiteral {
	ll := &ast.LambdaLiteral{Token: p.curToken}
	ll.Expressions = []ast.Expression{}

	for p.peekToken.Type != token.CloseParen && p.peekToken.Type != token.EOF {
		p.nextToken()
		expr := p.parseExpression()
		if expr != nil {
			ll.Expressions = append(ll.Expressions, expr)
		}
	}
	p.nextToken()

	return ll
}

func (p *Parser) parseQuoteExpression() *ast.QuoteExpression {
	qe := &ast.QuoteExpression{Token: p.curToken}
	p.nextToken()
	qe.Expr = p.parseExpression()
	return qe
}

func (p *Parser) parseBacktickExpression() *ast.BacktickExpression {
	be := &ast.BacktickExpression{Token: p.curToken}
	p.nextToken()
	be.Expr = p.parseExpression()
	return be
}

func (p *Parser) parseTildeExpression() *ast.TildeExpression {
	te := &ast.TildeExpression{Token: p.curToken}
	p.nextToken()
	te.Expr = p.parseExpression()
	return te
}

func (p *Parser) parseTildeAtExpression() *ast.TildeAtExpression {
	tae := &ast.TildeAtExpression{Token: p.curToken}
	p.nextToken()
	tae.Expr = p.parseExpression()
	return tae
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
