package lexer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenize(t *testing.T) {
	cases := []struct {
		Input []byte
		Want  []Token
	}{
		{
			Input: []byte("p"),
			Want:  []Token{{Type: TokCmd, Value: []byte("p")}},
		},
		{
			Input: []byte("12p"),
			Want:  []Token{{Type: TokDigits, Value: []byte("12")}, {Type: TokCmd, Value: []byte("p")}},
		},
		{
			Input: []byte("12pl"),
			Want: []Token{
				{Type: TokDigits, Value: []byte("12")},
				{Type: TokCmd, Value: []byte("p")},
				{Type: TokArg, Value: []byte("l")}},
		},
		{
			Input: []byte("w filename.txt"),
			Want: []Token{
				{Type: TokCmd, Value: []byte("w")},
				{Type: TokArg, Value: []byte(" filename.txt")}},
		},
	}
	for _, c := range cases {
		got, _ := Tokenize(c.Input)
		assert.Equal(t, got, c.Want)
	}
}
