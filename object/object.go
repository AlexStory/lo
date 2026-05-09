package object

import (
	"fmt"
	"hash/fnv"
	"lo/ast"
	"strconv"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ  ObjectType = "INTEGER"
	FLOAT_OBJ    ObjectType = "FLOAT"
	BOOLEAN_OBJ  ObjectType = "BOOLEAN"
	ERROR_OBJ    ObjectType = "ERROR"
	FUNCTION_OBJ ObjectType = "FUNCTION"
	BUILTIN_OBJ  ObjectType = "BUILTIN"
	LIST_OBJ     ObjectType = "LIST"
	STRING_OBJ   ObjectType = "STRING"
	KEYWORD_OBJ  ObjectType = "KEYWORD"
	MAP_OBJ      ObjectType = "MAP"
	NIL_OBJ      ObjectType = "Nil"
	SYMBOL_OBJ   ObjectType = "SYMBOL"
	MACRO_OBJ    ObjectType = "MACRO"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer represents an integer object
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Float represents a float object
type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return strconv.FormatFloat(f.Value, 'g', -1, 64) }

// Boolean represents a boolean object
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// Error represents an error object
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Function represents a user-defined function object
type Function struct {
	Name       string
	Parameters []*ast.Identifier
	Body       []ast.Expression
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	return fmt.Sprintf("(fn %s)", f.Name)
}

// BuiltinFunction represents a built-in function object
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// List represents a list object
type List struct {
	Elements []Object
}

func (l *List) Type() ObjectType { return LIST_OBJ }
func (l *List) Inspect() string {
	var out strings.Builder
	out.WriteString("[")
	for i, elem := range l.Elements {
		out.WriteString(elem.Inspect())
		if i < len(l.Elements)-1 {
			out.WriteString(" ")
		}
	}
	out.WriteString("]")
	return out.String()
}

// String returns the string representation of the object
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Keyword represents a keyword object
type Keyword struct {
	Value string
}

func (k *Keyword) Type() ObjectType { return KEYWORD_OBJ }
func (k *Keyword) Inspect() string  { return k.Value }
func (k *Keyword) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(k.Value))
	return HashKey{Type: k.Type(), Value: h.Sum64()}
}

type MapPair struct {
	Key   Object
	Value Object
}

// Map represents a map object
type Map struct {
	Pairs map[HashKey]MapPair
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	var out strings.Builder
	pairs := []string{}
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// Nil represents a lack of value
type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Inspect() string  { return "nil" }

// Symbol represents a symbol object
type Symbol struct {
	Value string
}

func (s *Symbol) Type() ObjectType { return SYMBOL_OBJ }
func (s *Symbol) Inspect() string  { return s.Value }

// Macro represents a macro object
type Macro struct {
	Parameters []*ast.Identifier
	Body       []ast.Expression
	Env        *Environment
}

func (m *Macro) Type() ObjectType { return MACRO_OBJ }
func (m *Macro) Inspect() string  { return "macro" }
