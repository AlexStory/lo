package astconv

import (
	"fmt"
	"lo/ast"
	"lo/consts"
	"lo/object"
)

// ObjectWrapper is used to inject already-evaluated objects back into the AST
type ObjectWrapper struct {
	Obj object.Object
}

func (ow *ObjectWrapper) TokenLiteral() string { return ow.Obj.Inspect() }
func (ow *ObjectWrapper) ExpressionNode()      {}

func AstToObject(node ast.Node) object.Object {
	if node == nil {
		return &consts.Nil
	}
	if ow, ok := node.(*ObjectWrapper); ok {
		return ow.Obj
	}
	switch node := node.(type) {
	case *ast.IntLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Identifier:
		return &object.Symbol{Value: node.Value}
	case *ast.ListExpression:
		elements := []object.Object{}
		for _, e := range node.Expressions {
			elements = append(elements, AstToObject(e))
		}
		return &object.List{Elements: elements}
	case *ast.ListLiteral:
		elements := []object.Object{}
		elements = append(elements, &object.Symbol{Value: "vec*"})
		for _, e := range node.Expressions {
			elements = append(elements, AstToObject(e))
		}
		return &object.List{Elements: elements}
	case *ast.MapLiteral:
		pairs := make(map[object.HashKey]object.MapPair)
		for _, p := range node.Pairs {
			key := AstToObject(p.Key)
			val := AstToObject(p.Value)
			if hashable, ok := key.(object.Hashable); ok {
				pairs[hashable.HashKey()] = object.MapPair{Key: key, Value: val}
			}
		}
		return &object.Map{Pairs: pairs}
	case *ast.LambdaLiteral:
		elements := []object.Object{}
		elements = append(elements, &object.Symbol{Value: "fn*"})
		for _, e := range node.Expressions {
			elements = append(elements, AstToObject(e))
		}
		return &object.List{Elements: elements}
	case *ast.Keyword:
		return &object.Keyword{Value: node.Value}
	case *ast.QuoteExpression:
		return &object.List{Elements: []object.Object{
			&object.Symbol{Value: "quote"},
			AstToObject(node.Expr),
		}}
	case *ast.BacktickExpression:
		return &object.List{Elements: []object.Object{
			&object.Symbol{Value: "backtick"},
			AstToObject(node.Expr),
		}}
	case *ast.TildeExpression:
		return &object.List{Elements: []object.Object{
			&object.Symbol{Value: "unquote"},
			AstToObject(node.Expr),
		}}
	case *ast.TildeAtExpression:
		return &object.List{Elements: []object.Object{
			&object.Symbol{Value: "unquote-splicing"},
			AstToObject(node.Expr),
		}}
	}
	return &consts.Nil
}

func ObjectToAst(obj object.Object) ast.Expression {
	switch obj := obj.(type) {
	case *object.Integer:
		return &ast.IntLiteral{Value: obj.Value}
	case *object.Float:
		return &ast.FloatLiteral{Value: obj.Value}
	case *object.String:
		return &ast.StringLiteral{Value: obj.Value}
	case *object.Symbol:
		return &ast.Identifier{Value: obj.Value}
	case *object.List:
		if len(obj.Elements) > 0 {
			if sym, ok := obj.Elements[0].(*object.Symbol); ok {
				if sym.Value == "fn*" {
					exprs := []ast.Expression{}
					for _, e := range obj.Elements[1:] {
						exprs = append(exprs, ObjectToAst(e))
					}
					return &ast.LambdaLiteral{Expressions: exprs}
				}
				if sym.Value == "vec*" {
					exprs := []ast.Expression{}
					for _, e := range obj.Elements[1:] {
						exprs = append(exprs, ObjectToAst(e))
					}
					return &ast.ListLiteral{Expressions: exprs}
				}
			}
		}
		exprs := []ast.Expression{}
		for _, e := range obj.Elements {
			exprs = append(exprs, ObjectToAst(e))
		}
		return &ast.ListExpression{Expressions: exprs}
	case *object.Map:
		pairs := []ast.MapPair{}
		for _, p := range obj.Pairs {
			pairs = append(pairs, ast.MapPair{
				Key:   ObjectToAst(p.Key),
				Value: ObjectToAst(p.Value),
			})
		}
		return &ast.MapLiteral{Pairs: pairs}
	case *object.Keyword:
		return &ast.Keyword{Value: obj.Value}
	case *object.Boolean:
		return &ast.Identifier{Value: fmt.Sprintf("%t", obj.Value)}
	case *object.Nil:
		return &ast.Identifier{Value: "nil"}
	}
	return &ObjectWrapper{Obj: obj}
}
