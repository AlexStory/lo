package eval

import (
	"fmt"
	"strings"

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
			return &object.Integer{Value: total.Value / arg.Value}
		case *object.Float:
			return &object.Float{Value: float64(total.Value) / arg.Value}
		}
	case *object.Float:
		switch arg := arg.(type) {
		case *object.Integer:
			return &object.Float{Value: total.Value / float64(arg.Value)}
		case *object.Float:
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
		fmt.Printf(arg.Inspect())
	}
	return nil
}

func println(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Printf(arg.Inspect())
	}
	fmt.Println()
	return nil
}
