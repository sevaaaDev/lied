package parser

import (
	"lied/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		cases := []struct {
			input []lexer.Token
			want  CommandNode
		}{
			{
				input: []lexer.Token{
					{Type: lexer.TokDigits, Value: []byte("10")},
					{Type: lexer.TokCmd, Value: []byte("p")}},
				want: CommandNode{
					lineRange: &LineRange{
						start: LineNumber{rawVal: "10"},
						end:   LineNumber{rawVal: "10"},
					},
					cmd: "p",
				},
			},
			{
				input: []lexer.Token{
					{Type: lexer.TokCmd, Value: []byte("p")}},
				want: CommandNode{
					lineRange: nil,
					cmd:       "p",
				},
			},
			{
				input: []lexer.Token{
					{Type: lexer.TokDigits, Value: []byte("10")}},
				want: CommandNode{
					lineRange: &LineRange{
						start: LineNumber{rawVal: "10"},
						end:   LineNumber{rawVal: "10"},
					},
					cmd: "set",
				},
			},
			{
				input: []lexer.Token{
					{Type: lexer.TokDigits, Value: []byte("10")},
					{Type: lexer.TokComma},
					{Type: lexer.TokCmd, Value: []byte("p")}},
				want: CommandNode{
					lineRange: &LineRange{
						start: LineNumber{rawVal: "10"},
						end:   LineSymbol{rawVal: "$"},
					},
					cmd: "set",
				},
			},
			{
				input: []lexer.Token{
					{Type: lexer.TokComma},
					{Type: lexer.TokDigits, Value: []byte("10")},
					{Type: lexer.TokCmd, Value: []byte("p")}},
				want: CommandNode{
					lineRange: &LineRange{
						start: LineNumber{rawVal: "1"},
						end:   LineNumber{rawVal: "10"},
					},
					cmd: "set",
				},
			},
			{
				input: []lexer.Token{
					{Type: lexer.TokSymbol, Value: []byte(".")},
					{Type: lexer.TokComma},
					{Type: lexer.TokDigits, Value: []byte("10")},
					{Type: lexer.TokCmd, Value: []byte("p")}},
				want: CommandNode{
					lineRange: &LineRange{
						start: LineSymbol{rawVal: "."},
						end:   LineNumber{rawVal: "10"},
					},
					cmd: "set",
				},
			},
		}
		for _, c := range cases {
			got, _ := Parse(c.input)
			assert.Equal(t, c.want, got)
		}
	})
	t.Run("sad", func(t *testing.T) {
		cases := []struct {
			input []lexer.Token
			err   bool
		}{
			{
				input: []lexer.Token{
					{Type: lexer.TokDigits, Value: []byte("10")},
					{Type: lexer.TokCmd, Value: []byte("p")},
					{Type: lexer.TokCmd, Value: []byte("p")}},
				err: true,
			},
			{
				input: []lexer.Token{
					{Type: lexer.TokDigits, Value: []byte("10")},
					{Type: lexer.TokComma},
					{Type: lexer.TokComma}},
				err: true,
			},
		}
		for _, c := range cases {
			_, got := Parse(c.input)
			if c.err {
				assert.NotNil(t, got)
			}
		}
	})
}
