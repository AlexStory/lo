package eval

import (
	"lo/expand"
	"lo/lexer"
	"lo/object"
	"lo/parser"
	"os"
	"testing"
)

var testExpander = expand.New(&Evaluator{})

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

func TestKeywordEval(t *testing.T) {
	input := ":key"
	evaluated := testEval(input)

	keyword, ok := evaluated.(*object.Keyword)
	if !ok {
		t.Fatalf("object is not Keyword. got=%T (%+v)", evaluated, evaluated)
	}

	if keyword.Value != ":key" {
		t.Errorf("keyword has wrong value. got=%s, want=%s", keyword.Value, ":key")
	}
}

func TestMapEval(t *testing.T) {
	input := `{:key "value" 1 2}`
	evaluated := testEval(input)

	resultMap, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("object is not Map. got=%T (%+v)", evaluated, evaluated)
	}

	if len(resultMap.Pairs) != 2 {
		t.Errorf("map has wrong number of pairs. got=%d", len(resultMap.Pairs))
	}

	expected := map[string]interface{}{
		":key": "value",
		"1":    int64(2),
	}

	for _, pair := range resultMap.Pairs {
		key := pair.Key.Inspect()
		expectedValue, ok := expected[key]
		if !ok {
			t.Errorf("unexpected key: %s", key)
			continue
		}

		switch v := expectedValue.(type) {
		case string:
			testStringObject(t, pair.Value, v)
		case int64:
			testIntegerObject(t, pair.Value, v)
		}
	}
}

func TestKeywordAsGetter(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`(:key {:key "value"})`, "value"},
		{`(:missing {:key "value"})`, nil},
		{`(:missing {:key "value"} "default")`, "default"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if tt.expected == nil {
			if _, ok := evaluated.(*object.Nil); !ok {
				t.Errorf("expected nil, got=%T (%+v) for input %s", evaluated, evaluated, tt.input)
			}
		} else {
			switch v := tt.expected.(type) {
			case string:
				testStringObject(t, evaluated, v)
			}
		}
	}
}

func TestLambdaLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(#(+ % 1) 5)", 6},
		{"(#(+ %1 1) 5)", 6},
		{"(#(+ %1 %2) 5 10)", 15},
		{"(map #(+ % 1) [1 2 3])", 0}, // Placeholder for list comparison
	}

	for i, tt := range tests {
		if i == 3 {
			// Special case for list
			evaluated := testEval(tt.input)
			list, ok := evaluated.(*object.List)
			if !ok {
				t.Fatalf("expected list, got %T", evaluated)
			}
			if len(list.Elements) != 3 {
				t.Fatalf("expected 3 elements, got %d", len(list.Elements))
			}
			testIntegerObject(t, list.Elements[0], 2)
			testIntegerObject(t, list.Elements[1], 3)
			testIntegerObject(t, list.Elements[2], 4)
			continue
		}
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestLet(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(let [x 1] x)", 1},
		{"(let [x 1 y 2] (+ x y))", 3},
		{"(let [x 1] (let [x 2] x))", 2},
		{"(let [x 1] (let [y x] y))", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"(defmacro unless [condition body] (list 'if condition nil body)) (unless false 10)", int64(10)},
		{"(defmacro unless [condition body] (list 'if condition nil body)) (unless true 10)", nil},
		{"(def x 5) (-> x (+ 1) (* 2))", int64(12)},
		{"(->> [1 2 3] (map #(+ % 1)) (filter #(> % 2)))", []int64{3, 4}},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		if tt.expected == nil {
			if _, ok := evaluated.(*object.Nil); !ok {
				t.Errorf("tests[%d] - expected nil, got=%T (%+v) for input %s", i, evaluated, evaluated, tt.input)
			}
		} else {
			switch v := tt.expected.(type) {
			case int64:
				testIntegerObject(t, evaluated, v)
			case []int64:
				list, ok := evaluated.(*object.List)
				if !ok {
					t.Fatalf("tests[%d] - expected list, got %T", i, evaluated)
				}
				if len(list.Elements) != len(v) {
					t.Fatalf("tests[%d] - expected length %d, got %d", i, len(v), len(list.Elements))
				}
				for j, val := range v {
					testIntegerObject(t, list.Elements[j], val)
				}
			}
		}
	}
}

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'a", "a"},
		{"'(1 2 3)", "[1 2 3]"},
		{"(quote (1 2 3))", "[1 2 3]"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated.Inspect() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, evaluated.Inspect())
		}
	}
}

func TestWhen(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"(when true 1)", int64(1)},
		{"(when false 1)", nil},
		{"(when true 1 2)", int64(2)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if tt.expected == nil {
			if _, ok := evaluated.(*object.Nil); !ok {
				t.Errorf("expected nil, got=%T", evaluated)
			}
		} else {
			testIntegerObject(t, evaluated, tt.expected.(int64))
		}
	}
}

func TestThreadingMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(-> 5 (+ 1) (* 2))", 12},
		{"(-> 5 (+ 1) *)", 6}, // one-arg threading
		{"(->> 5 (+ 1) (* 2))", 12},
		{"(->> 5 (- 10))", 5}, // (->> 5 (- 10)) => (- 10 5)
		{"(-> 5 (- 10))", -5}, // (-> 5 (- 10)) => (- 5 10)
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// Helpers

func testEval(input string) object.Object {
	l := lexer.New(input, "test")
	p := parser.New(l)
	program := p.Parse()

	env := object.NewEnvironment()
	loadStdlibForTest(env)

	expanded, err := testExpander.ExpandMacros(program, env)
	if err != nil {
		return &object.Error{Message: err.Error()}
	}

	evaluator := &Evaluator{Expander: testExpander}
	return evaluator.Eval(expanded, env)
}

func loadStdlibForTest(env *object.Environment) {
	// For tests, we'll continue to read the file from the disk to avoid
	// needing to export the embedded FS from main or duplicating it here.
	// This keeps the test logic simple and independent of the main package's embedding.
	paths := []string{"../stdlib.lo", "./stdlib.lo"}
	var contents []byte
	var err error
	for _, p := range paths {
		contents, err = os.ReadFile(p)
		if err == nil {
			l := lexer.New(string(contents), p)
			parser := parser.New(l)
			program := parser.Parse()
			expanded, err := testExpander.ExpandMacros(program, env)
			if err == nil {
				evaluator := &Evaluator{Expander: testExpander}
				evaluator.Eval(expanded, env)
			}
			return
		}
	}
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	t.Helper()

	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
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

func testError(t *testing.T, obj object.Object) bool {
	t.Helper()
	_, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object is not an error. got-%T (%+v)", obj, obj)
		return false
	}

	return true
}
