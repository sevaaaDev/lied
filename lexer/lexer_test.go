package lexer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
			Input: []byte("w filename.txt"),
			Want: []Token{
				{Type: TokCmd, Value: []byte("w")},
				{Type: TokArg, Value: []byte(" filename.txt")}},
		},
		{
			Input: []byte("s/re/repl"),
			Want: []Token{
				{Type: TokCmd, Value: []byte("s")},
				{Type: TokArg, Value: []byte("re")},
				{Type: TokArg, Value: []byte("repl")},
			},
		},
		{
			Input: []byte("s/re/"),
			Want: []Token{
				{Type: TokCmd, Value: []byte("s")},
				{Type: TokArg, Value: []byte("re")},
				{Type: TokArg, Value: []byte{}},
			},
		},
		{
			Input: []byte("s/re"),
			Want:  nil,
		},
	}
	for i, c := range cases {
		fmt.Println(i)
		got, _ := Tokenize(c.Input)
		assert.Equal(t, c.Want, got)
	}
}
