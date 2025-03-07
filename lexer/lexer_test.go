package lexer

import (
	"lo/token"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("math thest", func(t *testing.T) {
		input := `(+ 1 2)
(+ (+ 1 2) 3)
[1 2]`
		l := New(input, "test")

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
			expectedLine    int
			expectedColumn  int
		}{
			{token.OpenParen, "(", 1, 1},
			{token.Ident, "+", 1, 2},
			{token.Number, "1", 1, 4},
			{token.Number, "2", 1, 6},
			{token.CloseParen, ")", 1, 7},
			{token.OpenParen, "(", 2, 1},
			{token.Ident, "+", 2, 2},
			{token.OpenParen, "(", 2, 4},
			{token.Ident, "+", 2, 5},
			{token.Number, "1", 2, 7},
			{token.Number, "2", 2, 9},
			{token.CloseParen, ")", 2, 10},
			{token.Number, "3", 2, 12},
			{token.CloseParen, ")", 2, 13},
			{token.OpenBracket, "[", 3, 1},
			{token.Number, "1", 3, 2},
			{token.Number, "2", 3, 4},
			{token.CloseBracket, "]", 3, 5},
			{token.EOF, "", 3, 6},
		}

		for i, tt := range tests {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}

			if tok.Line != tt.expectedLine {
				t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, tt.expectedLine, tok.Line)
			}

			if tok.Column != tt.expectedColumn {
				t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d", i, tt.expectedColumn, tok.Column)
			}
		}
	})
}
