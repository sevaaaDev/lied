package main

import (
	"bufio"
	"errors"
	"fmt"
	"lied/context"
	"slices"
	"strings"

	"lied/lexer"
	"lied/parser"
	"os"

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

type Cursor struct {
	Pos int
}

type Buffer struct {
	value     []byte
	CursorPos int
}

func (b *Buffer) CursorLeft() {
	if b.CursorPos > 0 {
		b.CursorPos--
	}
}
func (b *Buffer) CursorRight() {
	if b.CursorPos < len(b.value) {
		b.CursorPos++
	}
}

func (b *Buffer) Insert(v ...byte) {
	b.value = slices.Insert(b.value, b.CursorPos, v...)
	b.CursorRight()
}

func (b *Buffer) Delete(i int, j int) {
	b.value = slices.Delete(b.value, i, j)
}

func (b *Buffer) Backspace() {
	if b.CursorPos > 0 {
		b.Delete(b.CursorPos-1, b.CursorPos)
		b.CursorLeft()
	}
}

func (b *Buffer) Value() string {
	return strings.ReplaceAll(string(b.value), "\t", "        ")
}
func (b *Buffer) Set(buf []byte) {
	b.value = buf
	b.CursorPos = len(buf)
}

func (b *Buffer) Len() int {
	return len(b.value)
}

func (b *Buffer) print(prompt string) {
	print("\033[2K") // clear current line
	print("\033[0G")
	print(prompt)
	line := b.Value() // TODO: Handle tab char
	print(line)
	print("\033[", len(line)+1+4, "G")
}

var ctrlC = fmt.Errorf("Pressed ^C")

func (b *Buffer) Readline(prompt string) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer fmt.Print("\n")
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	for {
		b.print(prompt)
		chars := make([]byte, 3)
		_, err := os.Stdin.Read(chars)
		if err != nil {
			return err
		}
		if chars[0] == 3 { // ^C
			return ctrlC
		}
		if chars[0] == 27 && chars[1] == '[' { // ESC SEQ
			switch chars[2] {
			case 'C':
				b.CursorRight()
			case 'D':
				b.CursorLeft()
			}
			continue
		}
		if chars[0] <= 0x1f || chars[0] == 0x7f { // CTRL
			switch chars[0] {
			case 0x0d: // CR
				fmt.Printf("\033[%dG", b.Len()+1)
				return nil
			case 0x7f:
				b.Backspace()
			case 0x09:
				b.Insert('\t')
			}
			continue
		}
		b.Insert(chars[0])
	}

}

func NewReadBuffer(value []byte) *Buffer {
	return &Buffer{
		value:     value,
		CursorPos: len(value),
	}
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
	readBuffer := NewReadBuffer([]byte("hello world"))
	for {
		readBuffer.Set([]byte{})
		if ctx.Mode == context.M_CHANGE {
			readBuffer.Set(ctx.Buffer[ctx.CurrentLine-1])
		}
		err := readBuffer.Readline("*a │")
		if err != nil {
			switch errors.Is(err, ctrlC) {
			case true:
				retcode = 130
				return
			case false:
				fmt.Println(err)
			}
		}
		line := []byte(readBuffer.Value())
		if ctx.Mode == context.M_CHANGE {
			ctx.Buffer[ctx.CurrentLine-1] = line
			continue
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
