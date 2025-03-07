package ast

import "lo/token"

// Node is the interface for all AST nodes
type Node interface {
	TokenLiteral() string
}

// Expression is the interface for all expression nodes
type Expression interface {
	Node
	expressionNode()
}

// Program represents the entire program
type Program struct {
	Expressions []Expression
}

func (p *Program) TokenLiteral() string {
	if len(p.Expressions) > 0 {
		return p.Expressions[0].TokenLiteral()
	}
	return ""
}

// ListExpression represents a list of expressions
type ListExpression struct {
	Token       token.Token // The '(' token
	Expressions []Expression
}

func (le *ListExpression) expressionNode()      {}
func (le *ListExpression) TokenLiteral() string { return le.Token.Literal }

// Identifier represents an identifier node
type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// IntLiteral represents an integer literal node
type IntLiteral struct {
	Token token.Token // The token.INT token
	Value int64
}

func (il *IntLiteral) expressionNode()      {}
func (il *IntLiteral) TokenLiteral() string { return il.Token.Literal }

// FloatLiteral represents a float literal node
type FloatLiteral struct {
	Token token.Token // The token.FLOAT token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }

// ListLiteral represents a list literal node
type ListLiteral struct {
	Token       token.Token // The [ token
	Expressions []Expression
}

func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }
