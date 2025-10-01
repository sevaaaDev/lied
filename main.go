package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"lied/context"
	"lied/lexer"
	"lied/parser"

	"golang.org/x/term"
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

// TODO: make method for delete, insert and stuff for buffer and cursor

func main() {
	retcode := 0
	defer func() {
		os.Exit(retcode)
	}()
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
	readBuffer := NewReadBuffer(nil)
	for {
		prompt := "*a │"
		_ = readBuffer.Readline(prompt)
		line := []byte(readBuffer.Value())
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
	// for {
	// 	prompt := "*a │"
	// 	line := readline(scanner, prompt)
	// 	if len(line) == 0 || line[0] != ':' {
	// 		ctx.Buffer = slices.Insert(ctx.Buffer, ctx.CurrentLine, line)
	// 		ctx.CurrentLine++
	// 		continue
	// 	}
	// 	tokens, err := lexer.Tokenize(line[1:])
	// 	if err != nil {
	// 		fmt.Println("Lexer:", err)
	// 		continue
	// 	}
	// 	node, err := parser.Parse(tokens)
	// 	if err != nil {
	// 		fmt.Println("Parser:", err)
	// 		continue
	// 	}
	// 	err = node.Eval(ctx)
	// 	if err != nil {
	// 		fmt.Println("Exec:", err)
	// 		continue
	// 	}
	// }
}
