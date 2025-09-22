package parser

import (
	"lied/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input []lexer.Token
		want  CommandNode
	}{
		{
			input: []lexer.Token{
				{Type: lexer.TokDigits, Value: []byte("10")},
				{Type: lexer.TokCmd, Value: []byte("p")}},
			want: CommandNode{
				lineRange: LineRange{
					start: LineNumber{rawVal: "10"},
					end:   LineNumber{rawVal: "10"},
				},
				cmd: "p",
			},
		},
	}
	for _, c := range cases {
		got, _ := Parse(c.input)
		assert.Equal(t, c.want, got)
	}
}
