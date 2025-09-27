package context

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
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
	Run  func(*Context, *[2]int, []string) error
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
			"d": {
				Name: "del",
				Desc: "print buffer",
				Run:  cmdDelete,
			},
			"s": {
				Name: "sub",
				Desc: "print buffer",
				Run:  cmdSubstitute,
			},
			"": {
				Name: "set",
				Desc: "print buffer",
				Run:  cmdSet,
			},
		},
	}
}

func cmdPrint(ctx *Context, lineRange *[2]int, _ []string) error {
	if lineRange == nil {
		lineRange = &[2]int{ctx.CurrentLine, ctx.CurrentLine}
	}
	if lineRange[0] <= 0 || lineRange[1] > len(ctx.Buffer) {
		return fmt.Errorf("invalid address")
	}
	printPseudoLine := false
	for i := lineRange[0]; i <= lineRange[1]; i++ {
		if printPseudoLine {
			fmt.Println("*  │")
			printPseudoLine = false
		}
		fmt.Printf("%3d│%s\n", i, string(ctx.Buffer[i-1]))
		if i == ctx.CurrentLine {
			printPseudoLine = !printPseudoLine
		}
	}
	return nil
}

func cmdQuit(_ *Context, _ *[2]int, _ []string) error {
	os.Exit(0)
	return nil
}

func cmdWrite(ctx *Context, _ *[2]int, args []string) error {
	filename := ctx.Filename
	if len(args) > 0 {
		filename = strings.TrimSpace(args[0])
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	bytesWritten := 0
	for _, v := range ctx.Buffer {
		n, _ := file.Write(v)
		bytesWritten += n
		n, _ = file.Write([]byte{10})
		bytesWritten += n
	}
	file.Sync()
	fmt.Println(bytesWritten)
	return nil
}

func cmdDelete(ctx *Context, lineRange *[2]int, _ []string) error {
	if lineRange == nil {
		lineRange = &[2]int{ctx.CurrentLine, ctx.CurrentLine}
	}
	if lineRange[0] <= 0 || lineRange[1] > len(ctx.Buffer) {
		return fmt.Errorf("invalid address")
	}
	ctx.Buffer = slices.Delete(ctx.Buffer, lineRange[0]-1, lineRange[1])
	if ctx.CurrentLine > len(ctx.Buffer) {
		ctx.CurrentLine = len(ctx.Buffer)
	}
	return nil
}

func cmdSet(ctx *Context, lineRange *[2]int, _ []string) error {
	if lineRange[1] > len(ctx.Buffer) {
		return fmt.Errorf("invalid address")
	}
	ctx.CurrentLine = lineRange[1]
	cmdPrint(ctx, lineRange, nil)
	return nil
}

func replacerFunc(repl string, n int) func([]byte) []byte {
	occurences := 0
	return func(b []byte) []byte {
		occurences++
		if n == 0 || occurences == n {
			return []byte(repl)
		}
		return b
	}
}
func cmdSubstitute(ctx *Context, lineRange *[2]int, args []string) error {
	if lineRange == nil {
		lineRange = &[2]int{ctx.CurrentLine, ctx.CurrentLine}
	}
	if lineRange[0] <= 0 || lineRange[1] > len(ctx.Buffer) {
		return fmt.Errorf("invalid address")
	}
	if len(args) == 0 {
		return fmt.Errorf("invalid arguments")
	}
	regex := args[0]
	repl := ""
	if len(args) >= 2 {
		repl = args[1]
	}
	suffix := 1
	// TODO: probably need refactor
	if len(args) >= 3 && len(args[2]) != 0 {
		if len(args[2]) > 1 {
			return fmt.Errorf("invalid arguments")
		}
		if unicode.IsDigit(rune(args[2][0])) {
			v, _ := strconv.Atoi(args[2])
			suffix = v
		} else if args[2] == "g" {
			suffix = 0
		} else {
			return fmt.Errorf("invalid arguments")
		}
	}
	re, err := regexp.CompilePOSIX(regex)
	if err != nil {
		return err
	}
	for i := lineRange[0]; i <= lineRange[1]; i++ {
		newLine := re.ReplaceAllFunc(ctx.Buffer[i-1], replacerFunc(repl, suffix))
		ctx.Buffer[i-1] = newLine
		cmdPrint(ctx, &[2]int{i, i}, nil)
	}
	return nil
}
