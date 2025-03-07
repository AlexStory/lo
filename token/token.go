package token

type TokenType string

const (
	Ident        TokenType = "IDENT"
	Illegal      TokenType = "ILLEGAL"
	EOF          TokenType = "EOF"
	Number       TokenType = "NUMBER"
	String       TokenType = "STRING"
	OpenParen    TokenType = "LPAREN"
	CloseParen   TokenType = "RPAREN"
	OpenBracket  TokenType = "LBRACKET"
	CloseBracket TokenType = "RBRACKET"
	Quote        TokenType = "QUOTE"
)

type Token struct {
	Type     TokenType
	Line     int
	Column   int
	Filename string
	Literal  string
}
