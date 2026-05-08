package eval

import (
	"fmt"
	"os"
	"strings"

	"lo/consts"
	"lo/object"
)

var builtinFunctions = map[string]object.BuiltinFunction{
	"+":       add,
	"-":       subtract,
	"*":       multiply,
	"/":       divide,
	"str":     str,
	"print":   print,
	"println": println,
	"head":    head,
	"get":     get,
	"slurp":   slurp,
	"spit":    spit,
	"env":     getenv,
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
			Message: "Incorrect number of arguments to add: expected 1 argument.",
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

func get(args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{
			Message: fmt.Sprintf("Incorrect number of arguments to get: expected 2, got %d.", len(args)),
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
		return &consts.Nil
	}

	return pair.Value
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
