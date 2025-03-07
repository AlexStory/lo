package lexer

import (
	"lo/token"
	"unicode/utf8"
)

type Lexer struct {
	filename     string
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	column       int
}

func New(input string, name string) *Lexer {
	l := &Lexer{input: input, filename: name, line: 1}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	tok.Filename = l.filename
	for {
		l.skipWhitespace()
		l.skipComments()
		if l.ch != ' ' && l.ch != ';' {
			break
		}
	}

	switch l.ch {
	case '(':
		tok = newToken(token.OpenParen, l, string(l.ch))
	case ')':
		tok = newToken(token.CloseParen, l, string(l.ch))
	case '[':
		tok = newToken(token.OpenBracket, l, string(l.ch))
	case ']':
		tok = newToken(token.CloseBracket, l, string(l.ch))
	case 0:
		tok = newToken(token.EOF, l, "")
	default:
		peekChar, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
		if isStartingDigit(l.ch, peekChar) {
			tok.Column = l.column
			tok.Line = l.line
			tok.Type = token.Number
			value := readDigit(l)
			tok.Literal = value
			return tok
		}
		tok.Column = l.column
		tok.Line = l.line
		tok.Type = token.Ident
		value := readIdentifier(l)
		tok.Literal = value
		return tok
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	var charLength int

	isNewLine := l.ch == '\n'
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.position = l.readPosition
		l.column += 1
		return
	} else {
		l.ch, charLength = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}

	l.position = l.readPosition
	l.readPosition += charLength
	if isNewLine {
		l.line++
		l.column = 1
	} else {
		l.column += charLength
	}
}

func isDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ch == '_' || ch == '.'
}

func isStartingDigit(ch, next rune) bool {
	return ('0' <= ch && ch <= '9') || (ch == '-' && isDigit(next))
}

func readDigit(l *Lexer) string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func readIdentifier(l *Lexer) string {
	position := l.position
	for !isWhitespace(l.ch) && l.ch != ')' && l.ch != '(' && l.ch != 0 && l.ch != '[' && l.ch != ']' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	if l.ch == ';' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}
}

func newToken(tokenType token.TokenType, l *Lexer, literal string) token.Token {
	return token.Token{Type: tokenType, Line: l.line, Column: l.column, Filename: l.filename, Literal: literal}
}
