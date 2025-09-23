package context

import (
	"fmt"
)

type Context struct {
	Buffer      [][]byte
	CurrentLine int
	Commands    map[string]Command
}

type Command struct {
	Name string
	Desc string
	Run  func(Context, [2]int) error
}

func NewContext() Context {
	return Context{
		Buffer:      [][]byte{},
		CurrentLine: 0,
		Commands: map[string]Command{
			"p": {
				Name: "print",
				Desc: "print buffer",
				Run:  cmdPrint,
			},
		},
	}
}

func cmdPrint(ctx Context, lineRange [2]int) error {
	if len(lineRange) == 0 {
		lineRange = [2]int{ctx.CurrentLine,ctx.CurrentLine}
	}
	if lineRange[1] > len(ctx.Buffer) {
		return fmt.Errorf("Print: Invalid Range")
	}
	for i := lineRange[0]; i <= lineRange[1]; i++ {
		fmt.Println(string(ctx.Buffer[i-1]))
	}
	return nil
}
