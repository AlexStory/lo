package parser

import (
	"fmt"
	"lo/ast"
	"lo/lexer"
	"testing"
)

func TestSimpleParse(t *testing.T) {
	input := `(+ 1 2)`
	l := lexer.New(input, "test")
	p := New(l)

	program := p.Parse()

	if len(program.Expressions) != 1 {
		t.Fatalf("program.Expressions does not contain 1 expression. got=%d", len(program.Expressions))
	}

	expr := program.Expressions[0]

	listExpr, ok := expr.(*ast.ListExpression)
	if !ok {
		t.Fatalf("expr not *ast.ListExpression. got=%T", expr)
	}

	if listExpr.TokenLiteral() != "(" {
		t.Fatalf("listExpr.TokenLiteral not (. got=%q", listExpr.TokenLiteral())
	}

	if len(listExpr.Expressions) != 3 {
		t.Fatalf("listExpr.Expressions does not contain 3 expressions. got=%d", len(listExpr.Expressions))
	}

	testIdent(t, listExpr.Expressions[0], "+")
	testIntLiteral(t, listExpr.Expressions[1], 1)
	testIntLiteral(t, listExpr.Expressions[2], 2)
}

func TestNestedParse(t *testing.T) {
	input := "(+ (+ 1 2) 3)"
	l := lexer.New(input, "test")
	p := New(l)

	program := p.Parse()

	if len(program.Expressions) != 1 {
		for _, e := range program.Expressions {
			fmt.Printf("%+v\n", e)
		}
		t.Fatalf("program.Expressions does not contain 1 expression. got=%d", len(program.Expressions))
	}

	expr := program.Expressions[0]

	listExpr, ok := expr.(*ast.ListExpression)
	if !ok {
		t.Fatalf("expr not *ast.ListExpression. got=%T", expr)
	}

	if listExpr.TokenLiteral() != "(" {
		t.Fatalf("listExpr.TokenLiteral not (. got=%q", listExpr.TokenLiteral())
	}

	if len(listExpr.Expressions) != 3 {
		t.Fatalf("listExpr.Expressions does not contain 3 expressions. got=%d", len(listExpr.Expressions))
	}
}

func TestListLiteralParse(t *testing.T) {
	input := "[1 2]"
	l := lexer.New(input, "test")
	p := New(l)

	program := p.Parse()

	if len(program.Expressions) != 1 {
		t.Fatalf("program.Expressions does not contain 1 expression. got=%d", len(program.Expressions))
	}

	expr := program.Expressions[0]

	listExpr, ok := expr.(*ast.ListLiteral)
	if !ok {
		t.Fatalf("expr not *ast.ListExpression. got=%T", expr)
	}

	if listExpr.TokenLiteral() != "[" {
		t.Fatalf("listExpr.TokenLiteral not [. got=%q", listExpr.TokenLiteral())
	}

	if len(listExpr.Expressions) != 2 {
		t.Fatalf("listExpr.Expressions does not contain 2 expressions. got=%d", len(listExpr.Expressions))
	}

	testIntLiteral(t, listExpr.Expressions[0], 1)
	testIntLiteral(t, listExpr.Expressions[1], 2)
}

func TestDefnParse(t *testing.T) {
	input := "(defn add [a b] (+ a b))"
	l := lexer.New(input, "test")
	p := New(l)

	program := p.Parse()

	if len(program.Expressions) != 1 {
		for _, e := range program.Expressions {
			fmt.Printf("%+v\n", e)
		}
		t.Fatalf("program.Expressions does not contain 1 expression. got=%d", len(program.Expressions))
	}

	expr := program.Expressions[0]

	defnExpr, ok := expr.(*ast.ListExpression)

	if !ok {
		t.Fatalf("expr not *ast.ListExpression. got=%T", expr)
	}

	if defnExpr.TokenLiteral() != "(" {
		t.Fatalf("defnExpr.TokenLiteral not (. got=%q", defnExpr.TokenLiteral())
	}

	if len(defnExpr.Expressions) != 4 {
		for _, e := range defnExpr.Expressions {
			fmt.Printf("%+v\n", e)
		}

		t.Fatalf("defnExpr.Expressions does not contain 4 expressions. got=%d", len(defnExpr.Expressions))
	}

}

// Helpers
func testIdent(t *testing.T, expr ast.Expression, value string) {
	t.Helper()

	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("expr not *ast.Identifier. got=%T", expr)
	}

	if ident.Value != value {
		t.Fatalf("ident.Value not %s. got=%s", value, ident.Value)
	}

	if ident.TokenLiteral() != value {
		t.Fatalf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
	}
}

func testIntLiteral(t *testing.T, expr ast.Expression, value int64) {
	t.Helper()

	intLiteral, ok := expr.(*ast.IntLiteral)
	if !ok {
		t.Fatalf("expr not *ast.IntLiteral. got=%T", expr)
	}

	if intLiteral.Value != value {
		t.Fatalf("intLiteral.Value not %d. got=%d", value, intLiteral.Value)
	}

	if intLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("intLiteral.TokenLiteral not 1. got=%s", intLiteral.TokenLiteral())
	}
}

func testFloatLiteral(t *testing.T, expr ast.Expression, value float64) {
	t.Helper()

	floatLiteral, ok := expr.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expr not *ast.FloatLiteral. got=%T", expr)
	}

	if floatLiteral.Value != value {
		t.Fatalf("floatLiteral.Value not %f. got=%f", value, floatLiteral.Value)
	}
}
