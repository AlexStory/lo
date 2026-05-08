package eval

import (
	"lo/object"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(+ 1 2)", 3},
		{"(+ 1 2 3)", 6},
		{"(+ 1 2 3 4)", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(- 1 2)", -1},
		{"(- 1 2 3)", -4},
		{"(- 1 2 3 4)", -8},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(* 1 2)", 2},
		{"(* 1 2 3)", 6},
		{"(* 1 2 3 4)", 24},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(/ 1 2)", 0},
		{"(/ 1 2 3)", 0},
		{"(/ 1 2 3 4)", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestStr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`(str "hello" "world")`, "helloworld"},
		{`(str "hello" " " "world")`, "hello world"},
		{`(str "hello" " " "world" "!")`, "hello world!"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestHead(t *testing.T) {
	input := `(head [1 2 3])`
	evaluated := testEval(input)

	testIntegerObject(t, evaluated, 1)

	input = `(head [])`
	evaluated = testEval(input)

	_, ok := evaluated.(*object.Nil)
	if !ok {
		t.Errorf("Expected result to be nil type")
	}

	input = `(head 4)`
	evaluated = testEval(input)
	testError(t, evaluated)
}
