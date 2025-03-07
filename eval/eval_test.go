package eval

import (
	"lo/lexer"
	"lo/object"
	"lo/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestMath(t *testing.T) {
	input := "(+ 1 2)"
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 3)
}

func TestNestedAddition(t *testing.T) {
	input := "(+ (+ 1 2) 3)"
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 6)
}

func TestMultipleAddition(t *testing.T) {
	input := "(+ 1 2 3 4)"
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func TestList(t *testing.T) {
	input := "[1 2 3 4]"
	evaluated := testEval(input)
	arr, ok := evaluated.(*object.List)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(arr.Elements) != 4 {
		t.Errorf("array has wrong number of elements. got=%d", len(arr.Elements))
	}

	for i, el := range arr.Elements {
		testIntegerObject(t, el, int64(i+1))
	}
}

func TestDef(t *testing.T) {
	input := "(def x 5)"
	evaluated := testEval(input)

	val, ok := evaluated.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", evaluated, evaluated)
	}

	if val.Value != 5 {
		t.Errorf("object has wrong value. got=%d, want=%d", val.Value, 5)
	}
}

// Helpers

func testEval(input string) object.Object {
	l := lexer.New(input, "test")
	p := parser.New(l)
	program := p.Parse()

	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	t.Helper()

	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	t.Helper()

	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}

	return true
}
