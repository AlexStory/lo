package eval

import (
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
