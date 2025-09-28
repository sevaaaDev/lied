package main

import (
	"bufio"
	"fmt"
	"lied/context"
	"lied/lexer"
	"lied/parser"
	"os"
	"slices"
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

func readline(scanner *bufio.Scanner, prompt string) []byte {
	var line []byte
	fmt.Print(prompt)
	for scanner.Scan() {
		if scanner.Bytes()[0] == '\n' {
			break
		}
		line = append(line, scanner.Bytes()...)
	}
	return line
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
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)
	for {
		prompt := "*a │"
		line := readline(scanner, prompt)
		if len(line) == 0 || line[0] != ':' {
			ctx.Buffer = slices.Insert(ctx.Buffer, ctx.CurrentLine, line)
			ctx.CurrentLine++
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
