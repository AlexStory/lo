package eval

import (
	"fmt"
	"os"
	"strings"

	"lo/consts"
	"lo/object"
)

var globalEvaluator *Evaluator

func SetEvaluator(e *Evaluator) {
	globalEvaluator = e
}

func GetBuiltin(name string) (object.BuiltinFunction, bool) {
	switch name {
	case "+":
		return add, true
	case "-":
		return subtract, true
	case "*":
		return multiply, true
	case "/":
		return divide, true
	case "str":
		return str, true
	case "print":
		return print, true
	case "println":
		return println, true
	case "head":
		return head, true
	case "get":
		return get, true
	case "assoc":
		return assoc, true
	case "dissoc":
		return dissoc, true
	case "first":
		return head, true
	case "rest":
		return rest, true
	case "cons":
		return cons, true
	case "count":
		return count, true
	case "map":
		return mapBuiltin, true
	case "filter":
		return filterBuiltin, true
	case "reduce":
		return reduceBuiltin, true
	case "slurp":
		return slurp, true
	case "spit":
		return spit, true
	case "env":
		return getenv, true
	case "=":
		return equals, true
	case "<":
		return lessThan, true
	case ">":
		return greaterThan, true
	case "<=":
		return lessThanEquals, true
	case ">=":
		return greaterThanEquals, true
	case "not":
		return not, true
	case "concat":
		return concat, true
	case "empty?":
		return isEmpty, true
	case "nil?":
		return isNil, true
	case "list?":
		return isList, true
	case "list":
		return list, true
	}
	return nil, false
}

func concat(args ...object.Object) object.Object {
	elements := []object.Object{}
	for _, arg := range args {
		list, ok := arg.(*object.List)
		if !ok {
			return &object.Error{Message: fmt.Sprintf("concat: expected list, got %s", arg.Type())}
		}
		elements = append(elements, list.Elements...)
	}
	return &object.List{Elements: elements}
}

func isEmpty(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("empty?: expected 1 argument, got %d", len(args))}
	}
	switch arg := args[0].(type) {
	case *object.List:
		if len(arg.Elements) == 0 {
			return &consts.TrueBool
		}
	case *object.Map:
		if len(arg.Pairs) == 0 {
			return &consts.TrueBool
		}
	case *object.String:
		if len(arg.Value) == 0 {
			return &consts.TrueBool
		}
	case *object.Nil:
		return &consts.TrueBool
	}
	return &consts.FalseBool
}

func isNil(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("nil?: expected 1 argument, got %d", len(args))}
	}
	if args[0].Type() == object.NIL_OBJ {
		return &consts.TrueBool
	}
	return &consts.FalseBool
}

func isList(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("list?: expected 1 argument, got %d", len(args))}
	}
	if args[0].Type() == object.LIST_OBJ {
		return &consts.TrueBool
	}
	return &consts.FalseBool
}

func list(args ...object.Object) object.Object {
	return &object.List{Elements: args}
}

func add(args ...object.Object) object.Object {
	var result object.Object = &object.Integer{Value: 0}

	for _, arg := range args {
		result = subAdd(result, arg)
	}
	return result
}

func subAdd(total, arg object.Object) object.Object {
	switch total := total.(type) {
	case *object.Integer:
		switch arg := arg.(type) {
		case *object.Integer:
			return &object.Integer{Value: total.Value + arg.Value}
		case *object.Float:
			return &object.Float{Value: float64(total.Value) + arg.Value}
		}
	case *object.Float:
		switch arg := arg.(type) {
		case *object.Integer:
			return &object.Float{Value: total.Value + float64(arg.Value)}
		case *object.Float:
			return &object.Float{Value: total.Value + arg.Value}
		}
	}
	return nil
}

func subtract(args ...object.Object) object.Object {
	var result object.Object = args[0]

	for _, arg := range args[1:] {
		result = sub(result, arg)
	}
	return result
}

func sub(total, arg object.Object) object.Object {
	switch total := total.(type) {
	case *object.Integer:
		switch arg := arg.(type) {
		case *object.Integer:
			return &object.Integer{Value: total.Value - arg.Value}
		case *object.Float:
			return &object.Float{Value: float64(total.Value) - arg.Value}
		}
	case *object.Float:
		switch arg := arg.(type) {
		case *object.Integer:
			return &object.Float{Value: total.Value - float64(arg.Value)}
		case *object.Float:
			return &object.Float{Value: total.Value - arg.Value}
		}
	}
	return nil
}

func multiply(args ...object.Object) object.Object {
	var result object.Object = &object.Integer{Value: 1}

	for _, arg := range args {
		result = mul(result, arg)
	}
	return result
}

func mul(total, arg object.Object) object.Object {
	switch total := total.(type) {
	case *object.Integer:
		switch arg := arg.(type) {
		case *object.Integer:
			return &object.Integer{Value: total.Value * arg.Value}
		case *object.Float:
			return &object.Float{Value: float64(total.Value) * arg.Value}
		}
	case *object.Float:
		switch arg := arg.(type) {
		case *object.Integer:
			return &object.Float{Value: total.Value * float64(arg.Value)}
		case *object.Float:
			return &object.Float{Value: total.Value * arg.Value}
		}
	}
	return nil
}

func divide(args ...object.Object) object.Object {
	var result object.Object = args[0]

	for _, arg := range args[1:] {
		result = div(result, arg)
	}
	return result
}

func div(total, arg object.Object) object.Object {
	switch total := total.(type) {
	case *object.Integer:
		switch arg := arg.(type) {
		case *object.Integer:
			if arg.Value == 0 {
				return &object.Integer{Value: 0}
			}
			return &object.Integer{Value: total.Value / arg.Value}
		case *object.Float:
			if arg.Value == 0.0 {
				return &object.Float{Value: 0.0}
			}
			return &object.Float{Value: float64(total.Value) / arg.Value}
		}
	case *object.Float:
		switch arg := arg.(type) {
		case *object.Integer:
			if arg.Value == 0 {
				return &object.Float{Value: 0.0}
			}
			return &object.Float{Value: total.Value / float64(arg.Value)}
		case *object.Float:
			if arg.Value == 0.0 {
				return &object.Float{Value: 0.0}
			}
			return &object.Float{Value: total.Value / arg.Value}
		}
	}
	return nil
}

func str(args ...object.Object) object.Object {
	var s strings.Builder

	for _, arg := range args {
		s.WriteString(arg.Inspect())
	}
	return &object.String{Value: s.String()}
}

func print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Printf("%s", arg.Inspect())
	}
	return nil
}

func println(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Printf("%s", arg.Inspect())
	}
	fmt.Println()
	return nil
}

func head(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{
			Message: "Incorrect number of arguments to head: expected 1 argument.",
		}
	}

	switch arg := args[0].(type) {
	case *object.List:
		if len(arg.Elements) < 1 {
			return &consts.Nil
		} else {
			return arg.Elements[0]
		}
	default:
		return &object.Error{
			Message: "Argument Error: head should be called with a list.",
		}
	}
}

func rest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{
			Message: "Incorrect number of arguments to rest: expected 1 argument.",
		}
	}

	switch arg := args[0].(type) {
	case *object.List:
		if len(arg.Elements) < 1 {
			return &object.List{Elements: []object.Object{}}
		} else {
			return &object.List{Elements: arg.Elements[1:]}
		}
	default:
		return &object.Error{
			Message: "Argument Error: rest should be called with a list.",
		}
	}
}

func cons(args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to cons: expected 2, got %d.", len(args)),
		}
	}

	first := args[0]
	rest, ok := args[1].(*object.List)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: second argument to cons should be a list, got %s.", args[1].Type()),
		}
	}

	newElements := append([]object.Object{first}, rest.Elements...)
	return &object.List{Elements: newElements}
}

func count(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to count: expected 1, got %d.", len(args)),
		}
	}

	switch arg := args[0].(type) {
	case *object.List:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.Map:
		return &object.Integer{Value: int64(len(arg.Pairs))}
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: count not supported for %s.", args[0].Type()),
		}
	}
}

func mapBuiltin(args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Message: fmt.Sprintf("map expects 2 arguments, got %d", len(args))}
	}

	fn := args[0]
	list, ok := args[1].(*object.List)
	if !ok {
		return &object.Error{Message: fmt.Sprintf("second argument to map must be a list, got %s", args[1].Type())}
	}

	newElements := make([]object.Object, len(list.Elements))
	for i, el := range list.Elements {
		res := applyFunction(fn, []object.Object{el}, nil, globalEvaluator)
		if isError(res) {
			return res
		}
		newElements[i] = res
	}

	return &object.List{Elements: newElements}
}

func filterBuiltin(args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Message: fmt.Sprintf("filter expects 2 arguments, got %d", len(args))}
	}

	fn := args[0]
	list, ok := args[1].(*object.List)
	if !ok {
		return &object.Error{Message: fmt.Sprintf("second argument to filter must be a list, got %s", args[1].Type())}
	}

	newElements := []object.Object{}
	for _, el := range list.Elements {
		res := applyFunction(fn, []object.Object{el}, nil, globalEvaluator)
		if isError(res) {
			return res
		}
		if res != &consts.FalseBool && res != &consts.Nil {
			newElements = append(newElements, el)
		}
	}

	return &object.List{Elements: newElements}
}

func reduceBuiltin(args ...object.Object) object.Object {
	if len(args) != 2 && len(args) != 3 {
		return &object.Error{Message: fmt.Sprintf("reduce expects 2 or 3 arguments, got %d", len(args))}
	}

	fn := args[0]
	var initial object.Object
	var list *object.List
	var ok bool

	if len(args) == 3 {
		initial = args[1]
		list, ok = args[2].(*object.List)
	} else {
		list, ok = args[1].(*object.List)
		if ok && len(list.Elements) > 0 {
			initial = list.Elements[0]
			list = &object.List{Elements: list.Elements[1:]}
		} else if ok {
			return &object.Error{Message: "reduce on empty list with no initial value"}
		}
	}

	if !ok {
		return &object.Error{Message: "last argument to reduce must be a list"}
	}

	result := initial
	for _, el := range list.Elements {
		result = applyFunction(fn, []object.Object{result, el}, nil, globalEvaluator)
		if isError(result) {
			return result
		}
	}

	return result
}

func get(args ...object.Object) object.Object {
	if len(args) != 2 && len(args) != 3 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to get: expected 2 or 3, got %d.", len(args)),
		}
	}

	m, ok := args[0].(*object.Map)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: first argument to get should be a map, got %s.", args[0].Type()),
		}
	}

	key, ok := args[1].(object.Hashable)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: second argument to get should be a hashable object, got %s.", args[1].Type()),
		}
	}

	hashed := key.HashKey()
	pair, ok := m.Pairs[hashed]
	if !ok {
		if len(args) == 3 {
			return args[2]
		}
		return &consts.Nil
	}

	return pair.Value
}

func assoc(args ...object.Object) object.Object {
	if len(args) < 3 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to assoc: expected at least 3, got %d.", len(args)),
		}
	}

	if (len(args)-1)%2 != 0 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to assoc: expected even number of key-value pairs, got %d.", len(args)-1),
		}
	}

	m, ok := args[0].(*object.Map)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: first argument to assoc should be a map, got %s.", args[0].Type()),
		}
	}

	// Create a new map to maintain immutability
	newPairs := make(map[object.HashKey]object.MapPair)
	for k, v := range m.Pairs {
		newPairs[k] = v
	}

	for i := 1; i < len(args); i += 2 {
		key, ok := args[i].(object.Hashable)
		if !ok {
			return &object.Error{
				Message: fmt.Sprintf("Argument Error: key in assoc should be a hashable object, got %s.", args[i].Type()),
			}
		}

		value := args[i+1]
		hashed := key.HashKey()
		newPairs[hashed] = object.MapPair{Key: args[i], Value: value}
	}

	return &object.Map{Pairs: newPairs}
}

func dissoc(args ...object.Object) object.Object {
	if len(args) < 2 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to dissoc: expected at least 2, got %d.", len(args)),
		}
	}

	m, ok := args[0].(*object.Map)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: first argument to dissoc should be a map, got %s.", args[0].Type()),
		}
	}

	// Create a new map to maintain immutability
	newPairs := make(map[object.HashKey]object.MapPair)
	for k, v := range m.Pairs {
		newPairs[k] = v
	}

	for i := 1; i < len(args); i++ {
		key, ok := args[i].(object.Hashable)
		if !ok {
			return &object.Error{
				Message: fmt.Sprintf("Argument Error: key in dissoc should be a hashable object, got %s.", args[i].Type()),
			}
		}

		hashed := key.HashKey()
		delete(newPairs, hashed)
	}

	return &object.Map{Pairs: newPairs}
}

func slurp(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to slurp: expected 1, got %d.", len(args)),
		}
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: argument to slurp should be a string (file path), got %s.", args[0].Type()),
		}
	}

	content, err := os.ReadFile(path.Value)
	if err != nil {
		return &object.Error{
			Message: fmt.Sprintf("IO Error: failed to read file %s: %s", path.Value, err),
		}
	}

	return &object.String{Value: string(content)}
}

func spit(args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to spit: expected 2, got %d.", len(args)),
		}
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: first argument to spit should be a string (file path), got %s.", args[0].Type()),
		}
	}

	content, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: second argument to spit should be a string, got %s.", args[1].Type()),
		}
	}

	err := os.WriteFile(path.Value, []byte(content.Value), 0644)
	if err != nil {
		return &object.Error{
			Message: fmt.Sprintf("IO Error: failed to write file %s: %s", path.Value, err),
		}
	}

	return &consts.Nil
}

func getenv(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to env: expected 1, got %d.", len(args)),
		}
	}

	name, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Argument Error: argument to env should be a string (env var name), got %s.", args[0].Type()),
		}
	}

	val := os.Getenv(name.Value)
	if val == "" {
		return &consts.Nil
	}

	return &object.String{Value: val}
}
func equals(args ...object.Object) object.Object {
	if len(args) < 2 {
		return &object.Error{Message: fmt.Sprintf("= expects at least 2 arguments, got %d", len(args))}
	}

	first := args[0]
	for _, arg := range args[1:] {
		if !objectsEqual(first, arg) {
			return &consts.FalseBool
		}
	}
	return &consts.TrueBool
}

func objectsEqual(left, right object.Object) bool {
	if left.Type() != right.Type() {
		return false
	}

	switch l := left.(type) {
	case *object.Integer:
		r := right.(*object.Integer)
		return l.Value == r.Value
	case *object.Float:
		r := right.(*object.Float)
		return l.Value == r.Value
	case *object.Boolean:
		r := right.(*object.Boolean)
		return l.Value == r.Value
	case *object.String:
		r := right.(*object.String)
		return l.Value == r.Value
	case *object.Keyword:
		r := right.(*object.Keyword)
		return l.Value == r.Value
	case *object.Nil:
		return true
	}
	if left.Type() == object.NIL_OBJ && right.Type() == object.NIL_OBJ {
		return true
	}
	return false
}

func lessThan(args ...object.Object) object.Object {
	return compare(args, func(l, r float64) bool { return l < r })
}

func greaterThan(args ...object.Object) object.Object {
	return compare(args, func(l, r float64) bool { return l > r })
}

func lessThanEquals(args ...object.Object) object.Object {
	return compare(args, func(l, r float64) bool { return l <= r })
}

func greaterThanEquals(args ...object.Object) object.Object {
	return compare(args, func(l, r float64) bool { return l >= r })
}

func compare(args []object.Object, op func(float64, float64) bool) object.Object {
	if len(args) < 2 {
		return &object.Error{Message: "Comparison operators expect at least 2 arguments"}
	}

	for i := 0; i < len(args)-1; i++ {
		lVal, okL := getFloatValue(args[i])
		rVal, okR := getFloatValue(args[i+1])

		if !okL || !okR {
			return &object.Error{Message: "Comparison operators expect numeric arguments"}
		}

		if !op(lVal, rVal) {
			return &consts.FalseBool
		}
	}
	return &consts.TrueBool
}

func getFloatValue(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	}
	return 0, false
}

func not(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("not expects 1 argument, got %d", len(args))}
	}

	if isFalsey(args[0]) {
		return &consts.TrueBool
	}
	return &consts.FalseBool
}
