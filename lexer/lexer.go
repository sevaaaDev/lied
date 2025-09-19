package lexer

import (
	"unicode"
)

type TokenType int

type Token struct {
	Type  TokenType
	Value []byte
}

const (
	TokCmd TokenType = iota
	TokDigits
)

func Tokenize(input []byte) []Token {
	i := 0
	var tokens []Token
	for i < len(input) {
		var buf []byte
		switch {
		case unicode.IsDigit(rune(input[i])):
			for i < len(input) && unicode.IsDigit(rune(input[i])) {
				buf = append(buf, input[i])
				i++
			}
			tokens = append(tokens, Token{Type: TokDigits, Value: buf})
		case input[i] == 'p':
			tokens = append(tokens, Token{Type: TokCmd, Value: []byte{input[i]}})
			i++
		case input[i] == 'q':
			tokens = append(tokens, Token{Type: TokCmd, Value: []byte{input[i]}})
			i++
		case input[i] == 'd':
			tokens = append(tokens, Token{Type: TokCmd, Value: []byte{input[i]}})
			i++
		default:
			i++
		}
	}
	return tokens
}
