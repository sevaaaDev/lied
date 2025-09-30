package main

import (
	"bufio"
	"fmt"
	"lied/context"
	"slices"

	// "lied/lexer"
	// "lied/parser"
	"os"
	// "slices"
	// "strings"

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

type Cursor struct {
	Pos int
}
type Buffer struct {
	value []byte
	Cursor
}

// TODO: make method for delete, insert and stuff for buffer and cursor

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
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("failed setting raw mode")
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	text := make([]byte, 0, 50)
	text = append(text, []byte("hello world")...)
	cursor := Cursor{
		Pos: len(text),
	}
	bigB := make([]byte, 0, 10)
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer logFile.Close()
	for len(bigB) < 10 {
		fmt.Print("\033[0E")
		fmt.Print("\033[2K") // clear current line
		fmt.Print(string(text))
		fmt.Print("\033[0E") // move to start of line
		if cursor.Pos > 0 {
			fmt.Printf("\033[%dC", cursor.Pos)
		}
		b := make([]byte, 3)
		_, err = os.Stdin.Read(b)
		if err != nil {
			fmt.Println("failed reading")
			os.Exit(2)
		}
		if b[0] == 27 && b[1] == '[' {
			switch b[2] {
			case 'C':
				if cursor.Pos < len(text) {
					fmt.Fprintln(logFile, "moving: ", len(text))
					cursor.Pos++
				}
			case 'D':
				if cursor.Pos > 0 {
					cursor.Pos--
				}
			}
			continue
		}
		if b[0] == 3 {
			break
		}
		if b[0] <= 0x1f || b[0] == 0x7f {
			switch b[0] {
			case 0x0d:
				return
			case 0x7f:
				text = slices.Delete(text, cursor.Pos-1, cursor.Pos)
				cursor.Pos--
			}
			continue
		}
		text = slices.Insert(text, cursor.Pos, b[0]) // make sure insert normal char, rn 0x0a is inserted too
		fmt.Fprintln(logFile, "adding: ", len(text))
		cursor.Pos++
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
