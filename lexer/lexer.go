package lexer

import (
	"fmt"
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
	TokComma
	TokSymbol
	TokArg
)

func peek(input []byte, want byte, currentIndex int) bool {
	if currentIndex >= len(input) {
		return false
	}
	if input[currentIndex] != want {
		return false
	}
	return true
}

// TODO: make peek func like parser
func Tokenize(input []byte) ([]Token, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty input")
	}
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
		case input[i] == 't':
			fallthrough
		case input[i] == 'w':
			tokens = append(tokens, Token{Type: TokCmd, Value: []byte{input[i]}})
			i++
			for i < len(input) {
				buf = append(buf, input[i])
				i++
			}
			if len(buf) != 0 {
				tokens = append(tokens, Token{Type: TokArg, Value: buf})
			}
		case input[i] == 's':
			tokens = append(tokens, Token{Type: TokCmd, Value: []byte{input[i]}})
			i++
			// s/re/repl
			// s/re/
			// s/re
			// s
			for peek(input, '/', i) {
				i++
				buf = []byte{}
				for i < len(input) && !peek(input, '/', i) {
					buf = append(buf, input[i])
					i++
				}
				tokens = append(tokens, Token{Type: TokArg, Value: buf})
			}
		case input[i] == 'p':
			fallthrough
		case input[i] == 'q':
			fallthrough
		case input[i] == 'd':
			tokens = append(tokens, Token{Type: TokCmd, Value: []byte{input[i]}})
			i++
		case input[i] == ',':
			tokens = append(tokens, Token{Type: TokComma})
			i++
		case input[i] == '$':
			tokens = append(tokens, Token{Type: TokSymbol, Value: []byte{input[i]}})
			i++
		case input[i] == '.':
			tokens = append(tokens, Token{Type: TokSymbol, Value: []byte{input[i]}})
			i++
		default:
			return nil, fmt.Errorf("invalid input '%s'", string(input[i]))
		}
	}
	return tokens, nil
}
