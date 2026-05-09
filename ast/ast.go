package ast

import "lo/token"

// Node is the interface for all AST nodes
type Node interface {
	TokenLiteral() string
}

// Expression is the interface for all expression nodes
type Expression interface {
	Node
	ExpressionNode()
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

func (le *ListExpression) ExpressionNode()      {}
func (le *ListExpression) TokenLiteral() string { return le.Token.Literal }

// Identifier represents an identifier node
type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) ExpressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// IntLiteral represents an integer literal node
type IntLiteral struct {
	Token token.Token // The token.INT token
	Value int64
}

func (il *IntLiteral) ExpressionNode()      {}
func (il *IntLiteral) TokenLiteral() string { return il.Token.Literal }

// FloatLiteral represents a float literal node
type FloatLiteral struct {
	Token token.Token // The token.FLOAT token
	Value float64
}

func (fl *FloatLiteral) ExpressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }

// ListLiteral represents a list literal node
type ListLiteral struct {
	Token       token.Token // The [ token
	Expressions []Expression
}

func (ll *ListLiteral) ExpressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }

// StringLiteral represents a string literal node
type StringLiteral struct {
	Token token.Token // The token.STRING token
	Value string
}

func (sl *StringLiteral) ExpressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// Keyword represents a keyword literal node
type Keyword struct {
	Token token.Token // The token.KEYWORD token
	Value string
}

func (k *Keyword) ExpressionNode()      {}
func (k *Keyword) TokenLiteral() string { return k.Token.Literal }

// MapLiteral represents a map literal node
type MapLiteral struct {
	Token token.Token // The { token
	Pairs []MapPair
}

type MapPair struct {
	Key   Expression
	Value Expression
}

func (ml *MapLiteral) ExpressionNode()      {}
func (ml *MapLiteral) TokenLiteral() string { return ml.Token.Literal }

// LambdaLiteral represents a #(...) lambda literal
type LambdaLiteral struct {
	Token       token.Token // The #( token
	Expressions []Expression
}

func (ll *LambdaLiteral) ExpressionNode()      {}
func (ll *LambdaLiteral) TokenLiteral() string { return ll.Token.Literal }

// QuoteExpression represents a 'expression
type QuoteExpression struct {
	Token token.Token // The ' token
	Expr  Expression
}

func (qe *QuoteExpression) ExpressionNode()      {}
func (qe *QuoteExpression) TokenLiteral() string { return qe.Token.Literal }

// BacktickExpression represents a `expression
type BacktickExpression struct {
	Token token.Token // The ` token
	Expr  Expression
}

func (be *BacktickExpression) ExpressionNode()      {}
func (be *BacktickExpression) TokenLiteral() string { return be.Token.Literal }

// TildeExpression represents a ~expression
type TildeExpression struct {
	Token token.Token // The ~ token
	Expr  Expression
}

func (te *TildeExpression) ExpressionNode()      {}
func (te *TildeExpression) TokenLiteral() string { return te.Token.Literal }

// TildeAtExpression represents a ~@expression
type TildeAtExpression struct {
	Token token.Token // The ~@ token
	Expr  Expression
}

func (tae *TildeAtExpression) ExpressionNode()      {}
func (tae *TildeAtExpression) TokenLiteral() string { return tae.Token.Literal }
