package context

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
	return nil
}
