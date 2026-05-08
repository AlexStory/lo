package eval

import (
	"lo/object"
	"os"
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

func TestGet(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`(get {:key "value"} :key)`, "value"},
		{`(get {:key "value"} :missing)`, nil},
		{`(get {:key "value"} :missing "default")`, "default"},
		{`(get {1 2} 1)`, int64(2)},
		{`(get {1 2} 2)`, nil},
		{`(get {1 2} 2 3)`, int64(3)},
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
			case int64:
				testIntegerObject(t, evaluated, v)
			}
		}
	}
}

func TestAssoc(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		{
			input: `(assoc {:a 1} :b 2)`,
			expected: map[string]interface{}{
				":a": int64(1),
				":b": int64(2),
			},
		},
		{
			input: `(assoc {:a 1} :a 2)`,
			expected: map[string]interface{}{
				":a": int64(2),
			},
		},
		{
			input: `(assoc {:a 1} :b 2 :c 3)`,
			expected: map[string]interface{}{
				":a": int64(1),
				":b": int64(2),
				":c": int64(3),
			},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		m, ok := evaluated.(*object.Map)
		if !ok {
			t.Fatalf("expected map, got %T", evaluated)
		}

		if len(m.Pairs) != len(tt.expected) {
			t.Fatalf("expected %d pairs, got %d", len(tt.expected), len(m.Pairs))
		}

		for _, pair := range m.Pairs {
			key := pair.Key.Inspect()
			expectedVal := tt.expected[key]
			switch v := expectedVal.(type) {
			case int64:
				testIntegerObject(t, pair.Value, v)
			}
		}
	}

	// Test immutability
	input := `(def my-map {:a 1}) (assoc my-map :b 2) my-map`
	evaluated := testEval(input)
	m, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("expected map, got %T", evaluated)
	}
	if len(m.Pairs) != 1 {
		t.Errorf("original map should have 1 pair, got %d", len(m.Pairs))
	}
}

func TestDissoc(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		{
			input: `(dissoc {:a 1 :b 2} :b)`,
			expected: map[string]interface{}{
				":a": int64(1),
			},
		},
		{
			input:    `(dissoc {:a 1 :b 2} :a :b)`,
			expected: map[string]interface{}{},
		},
		{
			input: `(dissoc {:a 1} :c)`,
			expected: map[string]interface{}{
				":a": int64(1),
			},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		m, ok := evaluated.(*object.Map)
		if !ok {
			t.Fatalf("expected map, got %T", evaluated)
		}

		if len(m.Pairs) != len(tt.expected) {
			t.Fatalf("expected %d pairs, got %d for input %s", len(tt.expected), len(m.Pairs), tt.input)
		}

		for _, pair := range m.Pairs {
			key := pair.Key.Inspect()
			expectedVal := tt.expected[key]
			switch v := expectedVal.(type) {
			case int64:
				testIntegerObject(t, pair.Value, v)
			}
		}
	}

	// Test immutability
	input := `(def my-map {:a 1}) (dissoc my-map :a) my-map`
	evaluated := testEval(input)
	m, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("expected map, got %T", evaluated)
	}
	if len(m.Pairs) != 1 {
		t.Errorf("original map should have 1 pair, got %d", len(m.Pairs))
	}
}

func TestSlurpSpit(t *testing.T) {
	filename := "test_slurp_spit.txt"
	content := "hello lo language"

	// test spit
	spitInput := `(spit "` + filename + `" "` + content + `")`
	testEval(spitInput)

	// test slurp
	slurpInput := `(slurp "` + filename + `")`
	evaluated := testEval(slurpInput)
	testStringObject(t, evaluated, content)

	// cleanup
	os.Remove(filename)
}

func TestEnv(t *testing.T) {
	os.Setenv("LO_TEST_VAR", "lo_value")
	defer os.Unsetenv("LO_TEST_VAR")

	input := `(env "LO_TEST_VAR")`
	evaluated := testEval(input)
	testStringObject(t, evaluated, "lo_value")

	input = `(env "NON_EXISTENT_VAR")`
	evaluated = testEval(input)
	if _, ok := evaluated.(*object.Nil); !ok {
		t.Errorf("expected nil for non-existent env var, got=%T", evaluated)
	}
}
