package eval

import (
	"fmt"
	"lo/ast"
	"lo/astconv"
	"lo/consts"
	"lo/lexer"
	"lo/object"
	"lo/parser"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Evaluator struct {
	Expander interface {
		ExpandMacros(ast.Node, *object.Environment) (ast.Node, error)
	}
}

func (e *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	return Eval(node, env, e)
}

func Eval(node ast.Node, env *object.Environment, evaluator *Evaluator) object.Object {
	if node == nil {
		return &consts.Nil
	}
	if ow, ok := node.(*astconv.ObjectWrapper); ok {
		return ow.Obj
	}
	switch node := node.(type) {
	case *ast.Program:
		result := evalProgram(node.Expressions, env, evaluator)
		return result

	case *ast.ListExpression:
		return evalList(node, env, evaluator)

	case *ast.IntLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.ListLiteral:
		return evalListLiteral(node, env, evaluator)
	case *ast.MapLiteral:
		return evalMapLiteral(node, env, evaluator)
	case *ast.LambdaLiteral:
		return evalLambdaLiteral(node, env, evaluator)
	case *ast.Keyword:
		return &object.Keyword{Value: node.Value}
	case *ast.QuoteExpression:
		return evalQuote(node, env, evaluator)
	case *ast.BacktickExpression:
		return evalBacktick(node, env, evaluator)
	}

	return nil
}

func evalProgram(exps []ast.Expression, env *object.Environment, evaluator *Evaluator) object.Object {
	var result object.Object

	for _, exp := range exps {
		result = Eval(exp, env, evaluator)
	}

	return result
}

func evalList(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	if len(le.Expressions) == 0 {
		return &object.Error{Message: "empty list"}
	}

	var f object.Object

	first := le.Expressions[0]
	if ident, ok := first.(*ast.Identifier); ok {
		switch ident.Value {
		case "def":
			return doDef(le, env, evaluator)
		case "defn":
			return evalDefn(le, env, evaluator)
		case "defmacro":
			return evalDefMacro(le, env, evaluator)
		case "\\":
			return evalLambda(le, env, evaluator)
		case "if":
			return evalIf(le, env, evaluator)
		case "when":
			return evalWhen(le, env, evaluator)
		case "let":
			return evalLet(le, env, evaluator)
		case "and":
			return evalAnd(le, env, evaluator)
		case "or":
			return evalOr(le, env, evaluator)
		case "quote":
			return evalQuoteSpecialForm(le, env, evaluator)
		case "import":
			return evalImport(le, env, evaluator)
		default:
			f = evalIdentifier(ident, env)
		}
	} else {
		f = Eval(first, env, evaluator)
	}

	if f == nil {
		return &object.Error{Message: "function/macro not found"}
	}

	if f.Type() != object.FUNCTION_OBJ && f.Type() != object.BUILTIN_OBJ && f.Type() != object.KEYWORD_OBJ {
		return &object.Error{Message: fmt.Sprintf("first element is not a function, got %s", f.Type())}
	}

	args := []object.Object{}
	for _, arg := range le.Expressions[1:] {
		evaluated := Eval(arg, env, evaluator)
		if isError(evaluated) {
			return evaluated
		}
		args = append(args, evaluated)
	}

	if f.Type() == object.KEYWORD_OBJ {
		return applyKeywordAsFunction(f.(*object.Keyword), args)
	}

	return applyFunction(f, args, env, evaluator)
}

func applyKeywordAsFunction(k *object.Keyword, args []object.Object) object.Object {
	if len(args) != 1 && len(args) != 2 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to keyword getter: expected 1 or 2, got %d.", len(args)),
		}
	}

	m, ok := args[0].(*object.Map)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: first argument to keyword getter should be a map, got %s.", args[0].Type()),
		}
	}

	hashed := k.HashKey()
	pair, ok := m.Pairs[hashed]
	if !ok {
		if len(args) == 2 {
			return args[1]
		}
		return &consts.Nil
	}

	return pair.Value
}

func applyFunction(fn object.Object, args []object.Object, env *object.Environment, evaluator *Evaluator) object.Object {
	if fn.Type() != object.BUILTIN_OBJ && fn.Type() != object.FUNCTION_OBJ {
		return &object.Error{Message: fmt.Sprintf("not a function, got %s", fn.Type())}
	}

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := object.NewEnclosedEnvironment(fn.Env)
		for i, param := range fn.Parameters {
			if param.Value == "&" {
				if i+1 < len(fn.Parameters) {
					restParam := fn.Parameters[i+1]
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
				extendedEnv.Set(param.Value, &consts.Nil)
			}
		}

		var result object.Object
		for _, exp := range fn.Body {
			result = Eval(exp, extendedEnv, evaluator)
		}
		return result
	case *object.Builtin:
		return fn.Fn(args...)
	}
	return nil
}

func evalListLiteral(ll *ast.ListLiteral, env *object.Environment, evaluator *Evaluator) object.Object {
	elements := []object.Object{}

	for _, exp := range ll.Expressions {
		evaluated := Eval(exp, env, evaluator)
		elements = append(elements, evaluated)
	}
	return &object.List{Elements: elements}
}

func evalMapLiteral(ml *ast.MapLiteral, env *object.Environment, evaluator *Evaluator) object.Object {
	pairs := make(map[object.HashKey]object.MapPair)

	for _, pair := range ml.Pairs {
		key := Eval(pair.Key, env, evaluator)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return &object.Error{Message: fmt.Sprintf("unusable as hash key: %s", key.Type())}
		}

		value := Eval(pair.Value, env, evaluator)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.MapPair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {

	name := ident.Value
	if name == "%" {
		name = "%1"
	}

	switch name {
	case "true":
		return &consts.TrueBool
	case "false":
		return &consts.FalseBool
	case "nil":
		return &consts.Nil
	}

	if b, ok := GetBuiltin(name); ok {
		return &object.Builtin{Fn: b}
	}

	val, ok := env.Get(name)
	if !ok {
		return &object.Error{Message: "identifier not found: " + name}
	}

	return val
}

func doDef(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	if len(le.Expressions) != 3 {
		return &object.Error{Message: "wrong number of arguments to def, got " + fmt.Sprint(len(le.Expressions)-1) + ", expected 2"}
	}

	ident, ok := le.Expressions[1].(*ast.Identifier)
	if !ok {
		return &object.Error{Message: "first argument to def must be an identifier"}
	}

	val := Eval(le.Expressions[2], env, evaluator)
	env.Set(ident.Value, val)
	return val
}

func evalDefn(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
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

	body := le.Expressions[3:]
	fn := &object.Function{Name: ident.Value, Parameters: params, Body: body, Env: env}
	env.Set(ident.Value, fn)
	return fn
}

func evalDefMacro(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	if len(le.Expressions) < 4 {
		return &object.Error{Message: "wrong number of arguments to defmacro, got " + fmt.Sprint(len(le.Expressions)-1) + ", expected 3"}
	}

	ident, ok := le.Expressions[1].(*ast.Identifier)
	if !ok {
		return &object.Error{Message: "first argument to defmacro must be an identifier"}
	}

	paramsExpr, ok := le.Expressions[2].(*ast.ListLiteral)
	if !ok {
		return &object.Error{Message: "second argument to defmacro must be a list of identifiers"}
	}

	params := []*ast.Identifier{}
	for _, p := range paramsExpr.Expressions {
		param, ok := p.(*ast.Identifier)
		if !ok {
			return &object.Error{Message: "parameters to defmacro must be identifiers"}
		}
		params = append(params, param)
	}

	body := le.Expressions[3:]
	macro := &object.Macro{Parameters: params, Body: body, Env: env}
	env.Set(ident.Value, macro)
	return macro
}

func evalLambda(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
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

	body := le.Expressions[2:]
	return &object.Function{Name: "lambda", Parameters: params, Body: body, Env: env}
}

func evalLambdaLiteral(ll *ast.LambdaLiteral, env *object.Environment, evaluator *Evaluator) object.Object {
	maxArg := 0

	// Find all %n and %& in the expressions
	var findArgs func(node ast.Node)
	findArgs = func(node ast.Node) {
		switch n := node.(type) {
		case *ast.Identifier:
			if n.Value == "%" {
				if maxArg < 1 {
					maxArg = 1
				}
			} else if len(n.Value) > 1 && n.Value[0] == '%' {
				if n.Value == "%&" {
					// TODO: handle variadic if needed
				} else {
					argNum, err := strconv.Atoi(n.Value[1:])
					if err == nil {
						if argNum > maxArg {
							maxArg = argNum
						}
					}
				}
			}
		case *ast.ListExpression:
			for _, expr := range n.Expressions {
				findArgs(expr)
			}
		case *ast.ListLiteral:
			for _, expr := range n.Expressions {
				findArgs(expr)
			}
		case *ast.MapLiteral:
			for _, pair := range n.Pairs {
				findArgs(pair.Key)
				findArgs(pair.Value)
			}
		case *ast.LambdaLiteral:
			// Don't descend into nested lambdas for argument discovery of the outer one
		}
	}

	for _, expr := range ll.Expressions {
		findArgs(expr)
	}

	params := make([]*ast.Identifier, maxArg)
	for i := 0; i < maxArg; i++ {
		params[i] = &ast.Identifier{Value: fmt.Sprintf("%%%d", i+1)}
	}

	return &object.Function{
		Name:       "lambda",
		Parameters: params,
		Body:       []ast.Expression{&ast.ListExpression{Token: ll.Token, Expressions: ll.Expressions}},
		Env:        env,
	}
}

func evalIf(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {

	if len(le.Expressions) != 4 {
		return &object.Error{Message: "wrong number of arguments to if, got " + fmt.Sprint(len(le.Expressions)-1) + ", expected 3"}
	}

	cond := Eval(le.Expressions[1], env, evaluator)
	if !isFalsey(cond) {
		return Eval(le.Expressions[2], env, evaluator)
	}
	return Eval(le.Expressions[3], env, evaluator)
}

func evalWhen(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	if len(le.Expressions) < 3 {
		return &object.Error{Message: "wrong number of arguments to when, expected at least 2"}
	}

	cond := Eval(le.Expressions[1], env, evaluator)
	if !isFalsey(cond) {
		var result object.Object
		for _, exp := range le.Expressions[2:] {
			result = Eval(exp, env, evaluator)
		}
		return result
	}
	return &consts.Nil
}

func evalLet(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	if len(le.Expressions) < 3 {
		return &object.Error{Message: "wrong number of arguments to let, expected at least 2"}
	}

	bindings, ok := le.Expressions[1].(*ast.ListLiteral)
	if !ok {
		return &object.Error{Message: "first argument to let must be a list of bindings"}
	}

	if len(bindings.Expressions)%2 != 0 {
		return &object.Error{Message: "bindings must be a list of even number of elements"}
	}

	enclosedEnv := object.NewEnclosedEnvironment(env)
	for i := 0; i < len(bindings.Expressions); i += 2 {
		ident, ok := bindings.Expressions[i].(*ast.Identifier)
		if !ok {
			return &object.Error{Message: "binding keys must be identifiers"}
		}
		val := Eval(bindings.Expressions[i+1], enclosedEnv, evaluator)
		if isError(val) {
			return val
		}
		enclosedEnv.Set(ident.Value, val)
	}

	var result object.Object
	for _, exp := range le.Expressions[2:] {
		result = Eval(exp, enclosedEnv, evaluator)
	}
	return result
}

type objectWrapper struct {
	obj object.Object
}

func (ow *objectWrapper) TokenLiteral() string { return "" }
func (ow *objectWrapper) ExpressionNode()      {}

func evalAnd(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	var result object.Object = &consts.TrueBool
	for _, exp := range le.Expressions[1:] {
		result = Eval(exp, env, evaluator)
		if isFalsey(result) {
			return result
		}
	}
	return result
}

func evalOr(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	var result object.Object = &consts.FalseBool
	for _, exp := range le.Expressions[1:] {
		result = Eval(exp, env, evaluator)
		if !isFalsey(result) {
			return result
		}
	}
	return result
}

func evalQuote(qe *ast.QuoteExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	return astconv.AstToObject(qe.Expr)
}

func evalQuoteSpecialForm(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	if len(le.Expressions) != 2 {
		return &object.Error{Message: fmt.Sprintf("quote: expected 1 argument, got %d", len(le.Expressions)-1)}
	}
	return astconv.AstToObject(le.Expressions[1])
}

func evalBacktick(be *ast.BacktickExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	return evalQuasiquote(be.Expr, env, evaluator)
}

func evalQuasiquote(node ast.Expression, env *object.Environment, evaluator *Evaluator) object.Object {
	switch n := node.(type) {
	case *ast.TildeExpression:
		return Eval(n.Expr, env, evaluator)
	case *ast.ListExpression:
		newElements := []object.Object{}
		for _, e := range n.Expressions {
			if tae, ok := e.(*ast.TildeAtExpression); ok {
				expanded := Eval(tae.Expr, env, evaluator)
				if list, ok := expanded.(*object.List); ok {
					newElements = append(newElements, list.Elements...)
				} else {
					return &object.Error{Message: "unquote-splicing: expected list"}
				}
			} else {
				newElements = append(newElements, evalQuasiquote(e, env, evaluator))
			}
		}
		return &object.List{Elements: newElements}
	default:
		return astconv.AstToObject(node)
	}
}

func evalImport(le *ast.ListExpression, env *object.Environment, evaluator *Evaluator) object.Object {
	if len(le.Expressions) != 2 {
		return &object.Error{Message: fmt.Sprintf("import: expected 1 argument, got %d", len(le.Expressions)-1)}
	}

	var path string
	var namespace string

	switch arg := le.Expressions[1].(type) {
	case *ast.StringLiteral:
		path = arg.Value
		// Namespace is the filename without extension
		namespace = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	case *ast.Identifier:
		// Project-style import: (import math) looks for src/math.lo or src/math/main.lo
		name := arg.Value
		namespace = name
		path = filepath.Join("src", name+".lo")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = filepath.Join("src", name, "main.lo")
		}
	default:
		return &object.Error{Message: "import: argument must be a string or identifier"}
	}

	// Resolve relative path if it's a string literal and looks like a relative path
	if _, ok := le.Expressions[1].(*ast.StringLiteral); ok {
		if !filepath.IsAbs(path) {
			// Get current file path from token
			currentFile := le.Token.Filename
			if currentFile != "" && currentFile != "repl" && currentFile != "main" {
				dir := filepath.Dir(currentFile)
				path = filepath.Join(dir, path)
			}
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return &object.Error{Message: fmt.Sprintf("import: failed to read file %s: %s", path, err)}
	}

	l := lexer.New(string(content), path)
	p := parser.New(l)
	program := p.Parse()
	if len(p.Errors) != 0 {
		return &object.Error{Message: fmt.Sprintf("import: parser errors in %s", path)}
	}

	// We need to expand macros in the imported file too!
	// But we don't have access to the expander here easily without creating a circular dependency
	// or passing it in.
	// For now, let's try to just evaluate it in a fresh environment.
	moduleEnv := object.NewEnvironment()

	// We need an expander. main.go has one.
	// This is a bit tricky because of the package structure.
	// But eval is a dependency of expand, not the other way around.
	// Wait, we can't use expand here.
	// Actually, ExpandMacros is what we should call if we want macros to work in modules.

	// Let's see how main.go handles it.
	// For now, let's just evaluate it and see.
	// If the module has macros, they won't be expanded in the module itself unless we expand it.

	// To avoid circular dependency, we might need a way to pass the expander to Eval
	// or have a global registry.
	if evaluator.Expander != nil {
		expanded, err := evaluator.Expander.ExpandMacros(program, moduleEnv)
		if err == nil {
			program = expanded.(*ast.Program)
		}
	}

	_ = Eval(program, moduleEnv, evaluator)

	// Export symbols to the current environment with namespace prefix
	for _, key := range moduleEnv.GetKeys() {
		val, _ := moduleEnv.Get(key)
		env.Set(namespace+"/"+key, val)
	}

	return &consts.Nil
}

func isFalsey(obj object.Object) bool {
	if obj == nil {
		return true
	}
	if b, ok := obj.(*object.Boolean); ok {
		return !b.Value
	}
	if obj.Type() == object.NIL_OBJ {
		return true
	}
	return false
}
