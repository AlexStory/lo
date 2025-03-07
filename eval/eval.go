package eval

import (
	"fmt"
	"lo/ast"
	"lo/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		result := evalProgram(node.Expressions, env)
		return result

	case *ast.ListExpression:
		return evalList(node, env)

	case *ast.IntLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.ListLiteral:
		return evalListLiteral(node, env)
	}

	return nil
}

func evalProgram(exps []ast.Expression, env *object.Environment) object.Object {
	var result object.Object

	for _, exp := range exps {
		result = Eval(exp, env)
	}

	return result
}

func evalList(le *ast.ListExpression, env *object.Environment) object.Object {
	if len(le.Expressions) == 0 {
		return &object.Error{Message: "empty list"}
	}

	var f object.Object

	first := le.Expressions[0]
	if ident, ok := first.(*ast.Identifier); ok {
		switch ident.Value {
		case "def":
			return doDef(le, env)
		case "defn":
			return evalDefn(le, env)
		case "\\":
			return evalLambda(le, env)
		default:
			f = evalIdentifier(ident, env)
		}
	} else {
		f = Eval(first, env)
	}

	if f.Type() != object.FUNCTION_OBJ && f.Type() != object.BUILTIN_OBJ {

		return &object.Error{Message: "first element is not a function"}
	}

	args := []object.Object{}
	for _, arg := range le.Expressions[1:] {
		args = append(args, Eval(arg, env))
	}

	return applyFunction(f, args, env)
}

func applyFunction(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	if fn.Type() != object.BUILTIN_OBJ && fn.Type() != object.FUNCTION_OBJ {
		return &object.Error{Message: fmt.Sprintf("not a function, got %s", fn.Type())}
	}

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := object.NewEnclosedEnvironment(fn.Env)
		for i, param := range fn.Parameters {
			extendedEnv.Set(param.Value, args[i])
		}

		return Eval(fn.Body, extendedEnv)
	case *object.Builtin:
		return fn.Fn(args...)
	}
	return nil
}

func evalListLiteral(ll *ast.ListLiteral, env *object.Environment) object.Object {
	elements := []object.Object{}

	for _, exp := range ll.Expressions {
		evaluated := Eval(exp, env)
		elements = append(elements, evaluated)
	}
	return &object.List{Elements: elements}
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {
	if b, ok := builtinFunctions[ident.Value]; ok {
		return &object.Builtin{Fn: b}
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return &object.Error{Message: "identifier not found: " + ident.Value}
	}

	return val
}

func doDef(le *ast.ListExpression, env *object.Environment) object.Object {
	if len(le.Expressions) != 3 {
		return &object.Error{Message: "wrong number of arguments to def, got " + fmt.Sprint(len(le.Expressions)-1) + ", expected 2"}
	}

	ident, ok := le.Expressions[1].(*ast.Identifier)
	if !ok {
		return &object.Error{Message: "first argument to def must be an identifier"}
	}

	val := Eval(le.Expressions[2], env)
	env.Set(ident.Value, val)
	return val
}

func evalDefn(le *ast.ListExpression, env *object.Environment) object.Object {
	if len(le.Expressions) < 4 {
		return &object.Error{Message: "wrong number of arguments to defn, got " + fmt.Sprint(len(le.Expressions)-1) + ", expected 3"}
	}

	ident, ok := le.Expressions[1].(*ast.Identifier)
	if !ok {
		return &object.Error{Message: "first argument to defn must be an identifier"}
	}

	paramsExpr, ok := le.Expressions[2].(*ast.ListLiteral)
	if !ok {
		return &object.Error{Message: "second argument to defn must be a list of identifiers"}
	}

	params := []*ast.Identifier{}
	for _, p := range paramsExpr.Expressions {
		param, ok := p.(*ast.Identifier)
		if !ok {
			return &object.Error{Message: "parameters to defn must be identifiers"}
		}
		params = append(params, param)
	}

	body := le.Expressions[3]
	fn := &object.Function{Name: ident.Value, Parameters: params, Body: body, Env: env}
	env.Set(ident.Value, fn)
	return fn
}

func evalLambda(le *ast.ListExpression, env *object.Environment) object.Object {
	if len(le.Expressions) < 3 {
		return &object.Error{Message: "wrong number of arguments to lambda, got " + fmt.Sprint(len(le.Expressions)-1) + ", expected 2"}
	}

	paramsExpr, ok := le.Expressions[1].(*ast.ListLiteral)
	if !ok {
		return &object.Error{Message: "first argument to lambda must be a list of identifiers"}
	}

	params := []*ast.Identifier{}
	for _, p := range paramsExpr.Expressions {
		param, ok := p.(*ast.Identifier)
		if !ok {
			return &object.Error{Message: "parameters to lambda must be identifiers"}
		}
		params = append(params, param)
	}

	body := le.Expressions[2]
	return &object.Function{Name: "lambda", Parameters: params, Body: body, Env: env}
}
