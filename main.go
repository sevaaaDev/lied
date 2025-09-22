package main

import (
	"bufio"
	"fmt"
	"lied/context"
	"lied/lexer"
	"lied/parser"
	"os"
	"slices"
	"strconv"
	"unicode"
)

func printBuf(buf [][]byte) {
	for i, l := range buf {
		fmt.Printf("%2d|%s\n", i+1, string(l))
	}
}

func parseCommand(input []byte) (int, error) {
	var buf []byte
	for _, v := range input {
		if !unicode.IsDigit(rune(v)) {
			return 0, fmt.Errorf("found non digit")
		}
		buf = append(buf, v)
	}
	i, _ := strconv.Atoi(string(buf))
	return i, nil
}

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
	ctx.CurrentLine = len(buf)
	for {
		var line []byte
		fmt.Printf("%2d|", ctx.CurrentLine+1)
		for scanner.Scan() {
			if scanner.Bytes()[0] == '\n' {
				break
			}
			line = append(line, scanner.Bytes()...)
		}
		if line[0] != ':' {
			ctx.Buffer = slices.Insert(ctx.Buffer, ctx.CurrentLine, line)
			ctx.CurrentLine++
			continue
		}
		tokens := lexer.Tokenize(line[1:])
		node, _ := parser.Parse(tokens)
		err := node.Eval(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
