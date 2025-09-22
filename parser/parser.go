package parser

import (
	"fmt"
	"lied/context"
	"lied/lexer"
	"strconv"
)

type Line interface {
	Eval(context.Context) int
}

type LineNumber struct {
	rawVal string
}

func (l LineNumber) Eval(_ context.Context) int {
	val, _ := strconv.Atoi(l.rawVal)
	return val
}

type LineSymbol struct {
	rawVal string
}

func (l LineSymbol) Eval(ctx context.Context) int {
	if l.rawVal == "$" {
		return len(ctx.Buffer)
	}
	if l.rawVal == "." {
		return ctx.CurrentLine
	}
	return 1
}

type LineRange struct {
	start Line
	end   Line
}

type CommandNode struct {
	lineRange LineRange
	cmd       string
}

func (c CommandNode) Eval(ctx context.Context) error {
	return nil
}

func Parse(tokens []lexer.Token) (CommandNode, error) {
	i := 0
	lineRange, errLr := parseLineRange(tokens, &i)
	cmdType, errCmd := parseCmdType(tokens, &i)
	cmdNode := CommandNode{}
	if errLr != nil && errCmd != nil {
		return cmdNode, fmt.Errorf("Parser: nothing parsed, empty command")
	}
	if errLr == nil {
		cmdNode.lineRange = lineRange
	}
	if errCmd == nil {
		cmdNode.cmd = cmdType
	}
	return cmdNode, nil
}

func parseLineRange(tokens []lexer.Token, i *int) (LineRange, error) {
	start, errStart := parseLine(tokens, i)
	if tokens[*i].Type != lexer.TokComma {
		if errStart == nil {
			return LineRange{start: start, end: start}, nil
		}
		return LineRange{}, fmt.Errorf("Parser: no line range found")
	}
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

func parseLine(tokens []lexer.Token, i *int) (Line, error) {
	if tokens[*i].Type == lexer.TokDigits {
		l := LineNumber{rawVal: string(tokens[*i].Value)}
		*i++
		return l, nil
	}
	if tokens[*i].Type == lexer.TokSymbol {
		l := LineSymbol{rawVal: string(tokens[*i].Value)}
		*i++
		return l, nil
	}
	*i++
	return nil, fmt.Errorf("Parser: no address found")
}

func parseCmdType(tokens []lexer.Token, i *int) (string, error) {
	if tokens[*i].Type == lexer.TokCmd {
		c := string(tokens[*i].Value)
		*i++
		return c, nil
	}
	*i++
	return "", fmt.Errorf("Parser: no cmd found")
}
