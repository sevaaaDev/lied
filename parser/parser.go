package parser

import (
	"fmt"
	"lied/context"
	"lied/lexer"
	"strconv"
)

type Line interface {
	Eval(*context.Context) (int, error)
}

type LineNumber struct {
	rawVal string
}

func (l LineNumber) Eval(_ *context.Context) (int, error) {
	val, _ := strconv.Atoi(l.rawVal)
	return val, nil
}

type LineSymbol struct {
	rawVal string
}

func (l LineSymbol) Eval(ctx *context.Context) (int, error) {
	if l.rawVal == "$" {
		return len(ctx.Buffer), nil
	}
	if l.rawVal == "." {
		return ctx.CurrentLine, nil
	}
	return 0, nil
}

type LineRange struct {
	start Line
	end   Line
}

func (lr *LineRange) Eval(ctx *context.Context) (*[2]int, error) {
	var res [2]int
	var err error
	res[0], err = lr.start.Eval(ctx)
	if err != nil {
		return nil, err
	}
	res[1], err = lr.end.Eval(ctx)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type CommandNode struct {
	lineRange *LineRange
	cmd       string
}

func (c CommandNode) Eval(ctx *context.Context) error {
	var err error
	var lr *[2]int = nil
	if c.lineRange != nil {
		lr, err = c.lineRange.Eval(ctx)
		switch true {
		case err != nil:
			fallthrough
		case lr[0] > lr[1]:
			return fmt.Errorf("invalid line range")
		}
	}
	cmd, ok := ctx.Commands[c.cmd]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return cmd.Run(ctx, lr)
}

/* ============================================================
	Parsing part
   ============================================================ */

func peek(tokens []lexer.Token, want lexer.TokenType, currentIndex int) bool {
	if currentIndex == len(tokens) {
		return false
	}
	if tokens[currentIndex].Type != want {
		return false
	}
	return true
}

func Parse(tokens []lexer.Token) (CommandNode, error) {
	i := 0
	lineRange, errLr := parseLineRange(tokens, &i)
	cmdType, errCmd := parseCmdType(tokens, &i)
	cmdNode := CommandNode{cmd: "set"}
	if i <= len(tokens)-1 {
		return cmdNode, fmt.Errorf("invalid syntax")
	}
	if errLr != nil && errCmd != nil {
		return cmdNode, fmt.Errorf("nothing is parsed")
	}
	if errLr == nil {
		cmdNode.lineRange = &lineRange
	}
	if errCmd == nil {
		cmdNode.cmd = cmdType
	}
	return cmdNode, nil
}

func parseLineRange(tokens []lexer.Token, i *int) (LineRange, error) {
	start, errStart := parseLine(tokens, i)
	if peek(tokens, lexer.TokComma, *i) {
		*i++
		end, errEnd := parseLine(tokens, i)
		lr := LineRange{start: LineNumber{rawVal: "1"}, end: LineSymbol{rawVal: "$"}}
		if errStart == nil {
			lr.start = start
		}
		if errEnd == nil {
			lr.end = end
		}
		return lr, nil
	}
	if errStart == nil {
		return LineRange{start: start, end: start}, nil
	}
	return LineRange{}, fmt.Errorf("no line range found")
}

func parseLine(tokens []lexer.Token, i *int) (Line, error) {
	if peek(tokens, lexer.TokDigits, *i) {
		l := LineNumber{rawVal: string(tokens[*i].Value)}
		*i++
		return l, nil
	}
	if peek(tokens, lexer.TokSymbol, *i) {
		l := LineSymbol{rawVal: string(tokens[*i].Value)}
		*i++
		return l, nil
	}
	return nil, fmt.Errorf("no address found")
}

func parseCmdType(tokens []lexer.Token, i *int) (string, error) {
	if peek(tokens, lexer.TokCmd, *i) {
		c := string(tokens[*i].Value)
		*i++
		return c, nil
	}
	return "", fmt.Errorf("no cmd found")
}
