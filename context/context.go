package context

import (
	"fmt"
	"os"
)

type Context struct {
	Buffer      [][]byte
	CurrentLine int
	Commands    map[string]Command
	Filename    string
}

type Command struct {
	Name string
	Desc string
	Run  func(*Context, *[2]int) error
}

func NewContext() *Context {
	return &Context{
		Buffer:      [][]byte{},
		CurrentLine: 0,
		Commands: map[string]Command{
			"p": {
				Name: "print",
				Desc: "print buffer",
				Run:  cmdPrint,
			},
			"q": {
				Name: "quit",
				Desc: "print buffer",
				Run:  cmdQuit,
			},
			"w": {
				Name: "write",
				Desc: "print buffer",
				Run:  cmdWrite,
			},
		},
	}
}

func cmdPrint(ctx *Context, lineRange *[2]int) error {
	if lineRange == nil {
		lineRange = &[2]int{ctx.CurrentLine, ctx.CurrentLine}
	}
	if lineRange[0] == 0 || lineRange[1] > len(ctx.Buffer) {
		return fmt.Errorf("invalid address")
	}
	for i := lineRange[0]; i <= lineRange[1]; i++ {
		fmt.Printf("%2d|%s\n", i, string(ctx.Buffer[i-1]))
	}
	return nil
}

func cmdQuit(_ *Context, _ *[2]int) error {
	os.Exit(0)
	return nil
}

func cmdWrite(ctx *Context, _ *[2]int) error {
	file, err := os.OpenFile(ctx.Filename, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file :%s", err)
	}
	defer file.Close()
	for _, v := range ctx.Buffer {
		file.Write(v)
		file.Write([]byte{10})
	}
	file.Sync()
	fmt.Println("wrote to", ctx.Filename)
	return nil
}
