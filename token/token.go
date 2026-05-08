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
	OpenBrace    TokenType = "LBRACE"
	CloseBrace   TokenType = "RBRACE"
	Quote        TokenType = "QUOTE"
	Keyword      TokenType = "KEYWORD"
	HashParen    TokenType = "#("
)

type Token struct {
	Type     TokenType
	Line     int
	Column   int
	Filename string
	Literal  string
}
