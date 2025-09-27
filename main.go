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

func readFile(filename string) [][]byte {
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
	}
	scanner := bufio.NewScanner(file)
	var buf [][]byte
	for scanner.Scan() {
		buf = append(buf, scanner.Bytes())
	}
	return buf
}
func writeFile(filename string, buf [][]byte) {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
	}
	for _, v := range buf {
		file.Write(v)
		file.Write([]byte{10})
	}
	file.Sync()
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
	if len(os.Args) == 1 {
		fmt.Println("need file")
		os.Exit(1)
	}
	ctx := context.NewContext()
	buf := readFile(os.Args[1])
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)
	ctx.Buffer = buf
	ctx.CurrentLine = len(ctx.Buffer)
	ctx.Filename = os.Args[1]
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
