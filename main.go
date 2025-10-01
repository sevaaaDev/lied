package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/sevaaadev/lied/context"
	"github.com/sevaaadev/lied/lexer"
	"github.com/sevaaadev/lied/parser"
	"github.com/sevaaadev/lied/readline"
)

func readFile(filename string) ([][]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var buf [][]byte
	for scanner.Scan() {
		buf = append(buf, scanner.Bytes())
	}
	return buf, nil
}

func main() {
	ctx := context.NewContext()
	if len(os.Args) > 1 {
		buf, err := readFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ctx.Buffer = buf
		ctx.Filename = os.Args[1]
	}
	ctx.CurrentLine = len(ctx.Buffer)

	rl := readline.New()
	prompt := "*a │"
	rl.SetPrompt(prompt)
	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(line) == 0 || line[0] != ':' {
			if ctx.Mode == context.M_APPEND {
				ctx.Buffer = slices.Insert(ctx.Buffer, ctx.CurrentLine, line)
				ctx.CurrentLine++
			}
			continue
		}
		tokens, err := lexer.Tokenize(line[1:])
		if err != nil {
			fmt.Println("Lexer:", err)
			continue
		}
		node, err := parser.Parse(tokens)
		if err != nil {
			fmt.Println("Parser:", err)
			continue
		}
		err = node.Eval(ctx)
		if err != nil {
			fmt.Println("Exec:", err)
			continue
		}
	}
}
