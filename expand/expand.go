package expand

import (
	"fmt"
	"lo/ast"
	"lo/astconv"
	"lo/object"
)

// Evaluator is an interface to avoid circular dependency with eval package
type Evaluator interface {
	Eval(node ast.Node, env *object.Environment) object.Object
}

type Expander struct {
	evaluator Evaluator
}

func New(evaluator Evaluator) *Expander {
	return &Expander{evaluator: evaluator}
}

func (e *Expander) ExpandMacros(node ast.Node, env *object.Environment) (ast.Node, error) {
	if node == nil {
		return nil, nil
	}
	switch n := node.(type) {
	case *ast.Program:
		for i, expr := range n.Expressions {
			expanded, err := e.ExpandMacros(expr, env)
			if err != nil {
				return nil, err
			}
			n.Expressions[i] = expanded.(ast.Expression)

			// Special case: if it was a defmacro, we need to evaluate it
			// so that subsequent expressions can use it.
			if le, ok := expanded.(*ast.ListExpression); ok && len(le.Expressions) > 0 {
				if ident, ok := le.Expressions[0].(*ast.Identifier); ok {
					if ident.Value == "defmacro" || ident.Value == "import" {
						e.evaluator.Eval(le, env)
					}
				}
			}
		}
		return n, nil

	case *ast.ListExpression:
		return e.expandListExpression(n, env)

	case *ast.ListLiteral:
		for i, expr := range n.Expressions {
			expanded, err := e.ExpandMacros(expr, env)
			if err != nil {
				return nil, err
			}
			n.Expressions[i] = expanded.(ast.Expression)
		}
		return n, nil

	case *ast.MapLiteral:
		for i, pair := range n.Pairs {
			expandedKey, err := e.ExpandMacros(pair.Key, env)
			if err != nil {
				return nil, err
			}
			expandedVal, err := e.ExpandMacros(pair.Value, env)
			if err != nil {
				return nil, err
			}
			n.Pairs[i] = ast.MapPair{
				Key:   expandedKey.(ast.Expression),
				Value: expandedVal.(ast.Expression),
			}
		}
		return n, nil

	case *ast.LambdaLiteral:
		for i, expr := range n.Expressions {
			expanded, err := e.ExpandMacros(expr, env)
			if err != nil {
				return nil, err
			}
			n.Expressions[i] = expanded.(ast.Expression)
		}
		return n, nil

	case *ast.QuoteExpression:
		// Don't expand inside quote
		return n, nil

	case *ast.BacktickExpression:
		// Inside backtick, we might have unquote (~), which SHOULD be expanded if they contain macros.
		expanded, err := e.expandQuasiquote(n.Expr, env)
		if err != nil {
			return nil, err
		}
		n.Expr = expanded.(ast.Expression)
		return n, nil

	default:
		return n, nil
	}
}

func (e *Expander) expandQuasiquote(node ast.Expression, env *object.Environment) (ast.Expression, error) {
	switch n := node.(type) {
	case *ast.TildeExpression:
		expanded, err := e.ExpandMacros(n.Expr, env)
		if err != nil {
			return nil, err
		}
		n.Expr = expanded.(ast.Expression)
		return n, nil
	case *ast.ListExpression:
		for i, expr := range n.Expressions {
			expanded, err := e.expandQuasiquote(expr, env)
			if err != nil {
				return nil, err
			}
			n.Expressions[i] = expanded
		}
		return n, nil
	case *ast.ListLiteral:
		for i, expr := range n.Expressions {
			expanded, err := e.expandQuasiquote(expr, env)
			if err != nil {
				return nil, err
			}
			n.Expressions[i] = expanded
		}
		return n, nil
	default:
		return n, nil
	}
}

func (e *Expander) expandListExpression(le *ast.ListExpression, env *object.Environment) (ast.Node, error) {
	if len(le.Expressions) == 0 {
		return le, nil
	}

	first := le.Expressions[0]
	if ident, ok := first.(*ast.Identifier); ok {
		// If it's a macro, expand it.
		if obj, ok := env.Get(ident.Value); ok && obj.Type() == object.MACRO_OBJ {
			macro := obj.(*object.Macro)
			args := []object.Object{}
			for _, arg := range le.Expressions[1:] {
				args = append(args, astconv.AstToObject(arg))
			}

			// Expand macro
			expandedObj, err := e.applyMacro(macro, args, env)
			if err != nil {
				return nil, err
			}

			expandedAst := astconv.ObjectToAst(expandedObj)
			// Fixed point: expand the result of expansion
			return e.ExpandMacros(expandedAst, env)
		}
	}

	// Recursively expand all elements of the list
	for i, expr := range le.Expressions {
		expanded, err := e.ExpandMacros(expr, env)
		if err != nil {
			return nil, err
		}
		le.Expressions[i] = expanded.(ast.Expression)
	}

	return le, nil
}

func (e *Expander) applyMacro(macro *object.Macro, args []object.Object, env *object.Environment) (object.Object, error) {
	extendedEnv := object.NewEnclosedEnvironment(macro.Env)
	for i, param := range macro.Parameters {
		if param.Value == "&" {
			if i+1 < len(macro.Parameters) {
				restParam := macro.Parameters[i+1]
				restArgs := []object.Object{}
				if i < len(args) {
					restArgs = args[i:]
				}
				extendedEnv.Set(restParam.Value, &object.List{Elements: restArgs})
				break
			}
		}
		if i < len(args) {
			extendedEnv.Set(param.Value, args[i])
		} else {
			extendedEnv.Set(param.Value, &object.Nil{})
		}
	}

	var result object.Object
	for _, exp := range macro.Body {
		result = e.evaluator.Eval(exp, extendedEnv)
		if isError(result) {
			return nil, fmt.Errorf("%s", result.Inspect())
		}
	}

	return result, nil
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
