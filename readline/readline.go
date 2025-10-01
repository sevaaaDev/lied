package readline

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"slices"
	"strings"
)

type buffer struct {
	b []byte
}

func (b *buffer) insert(pos int, v ...byte) {
	b.b = slices.Insert(b.b, pos, v...)
}
func (b *buffer) delete(pos int) {
	b.b = slices.Delete(b.b, pos, pos+1)
}
func (b *buffer) delBefore(pos int) {
	b.b = slices.Delete(b.b, pos-1, pos)
}

type rl struct {
	buf       buffer
	cursorPos int
	prompt    string
}

func (rl *rl) refreshLine() {
	cmd := []byte{}
	cmd = append(cmd, '\r')
	cmd = append(cmd, []byte(rl.prompt)...)
	cmd = append(cmd, rl.buf.b...)
	cmd = append(cmd, []byte("\x1b[0K")...)
	cmd = fmt.Appendf(cmd, "\r\x1b[%dC", rl.cursorPos+len(rl.prompt))
	fmt.Print(string(cmd))
}
func (rl *rl) insert(v ...byte) {
	rl.buf.insert(rl.cursorPos, v...)
	rl.cursorPos++
}
func (rl *rl) backspace() {
	if rl.cursorPos > 0 {
		rl.buf.delBefore(rl.cursorPos)
		rl.cursorPos--
	}
}
func (rl *rl) cursorLeft() {
	if rl.cursorPos > 0 {
		rl.cursorPos--
	}
}
func (rl rl) cursorRight() {
	if rl.cursorPos < len(rl.buf.b) {
		rl.cursorPos++
	}
}

func (rl *rl) Readline() error {
	// makeraw
	for {
		rl.refreshLine()
		seq := make([]byte, 3)
		n, err := os.Stdin.Read(seq)
		if err != nil {
			return err
		}
		if n == 1 {
			if seq[0] == 0x08 {
				rl.backspace()
				continue
			}
			rl.insert(seq[0])
			// handle 1 char
		} else {
			// handle esc seq
		}

	}
}

type foo struct {
	value     []byte
	CursorPos int
}

func (rl *foo) CursorLeft() {
	if rl.CursorPos > 0 {
		rl.CursorPos--
	}
}
func (rl *foo) CursorRight() {
	if rl.CursorPos < len(rl.value) {
		rl.CursorPos++
	}
}

func (rl *foo) Insert(v ...byte) {
	rl.value = slices.Insert(rl.value, rl.CursorPos, v...)
	rl.CursorRight()
}

func (rl *foo) Delete(i int, j int) {
	rl.value = slices.Delete(rl.value, i, j)
}

func (rl *foo) Backspace() {
	if rl.CursorPos > 0 {
		rl.Delete(rl.CursorPos-1, rl.CursorPos)
		rl.CursorLeft()
	}
}

func (rl *foo) Value() string {
	return strings.ReplaceAll(string(rl.value), "\t", "        ")
}
func (rl *foo) Set(buf []byte) {
	rl.value = buf
	rl.CursorPos = len(buf)
}

func (rl *foo) Len() int {
	return len(rl.value)
}

func (rl *foo) print(prompt string) {
	print("\033[2K") // clear current line
	print("\033[0G")
	print(prompt)
	line := rl.Value() // TODO: Handle tab char
	print(line)
	print("\033[", len(line)+1+4, "G")
}

var ctrlC = fmt.Errorf("Pressed ^C")

func (rl *foo) Readline(prompt string) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer fmt.Print("\n")
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	for {
		rl.print(prompt)
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
				rl.CursorRight()
			case 'D':
				rl.CursorLeft()
			}
			continue
		}
		if chars[0] <= 0x1f || chars[0] == 0x7f { // CTRL
			switch chars[0] {
			case 0x0d: // CR
				fmt.Printf("\033[%dG", rl.Len()+1)
				return nil
			case 0x7f:
				rl.Backspace()
			case 0x09:
				rl.Insert('\t')
			}
			continue
		}
		rl.Insert(chars[0])
	}

}

// Readline instances
// struct readline {
// 	buf buffer
//	cursorpos uint
//	prompt string
// 	readline()
// }

// struct buffer {
// 	realbuffer []slice
//	insert(pos, v)
//	delete(pos)
//	backspace(pos)
// }

func New(value []byte) *foo {
	return &foo{
		value:     value,
		CursorPos: len(value),
	}
}
