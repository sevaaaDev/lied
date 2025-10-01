package readline

import (
	"fmt"
	"os"
	"slices"
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
	cmd = fmt.Appendf(cmd, "\r\x1b[%dC", rl.cursorPos-1+len(rl.prompt)-1)
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
func (rl *rl) cursorRight() {
	if rl.cursorPos < len(rl.buf.b) {
		rl.cursorPos++
	}
}

func (rl *rl) SetPrompt(p string) {
	rl.prompt = p
}

func (rl *rl) readChar() bool {
	rl.refreshLine()
	seq := make([]byte, 3)
	n, err := os.Stdin.Read(seq)
	if err != nil {
		return false
	}
	if n == 1 {
		if seq[0] > 0x1f && seq[0] != 0x7f {
			rl.insert(seq[0])
			return true
		}
		switch seq[0] {
		case 0x0a:
			rl.refreshLine()
			return false
		case 0x7f:
			fallthrough
		case 0x08:
			rl.backspace()
		case 0x09:
			rl.insert('\t')
		}
	} else {
		// handle esc seq
		if seq[0] == 27 && seq[1] == '[' { // ESC SEQ
			switch seq[2] {
			case 'C':
				rl.cursorRight()
			case 'D':
				rl.cursorLeft()
			}
		}
	}
	return true
}

func (rl *rl) Readline() ([]byte, error) {
	rl.buf.b = make([]byte, 0)
	rl.cursorPos = len(rl.buf.b)

	oldstate, err := enableRawMode()
	if err != nil {
		return nil, err
	}
	defer print("\n")
	defer disableRawMode(oldstate)

	for rl.readChar() {
	}

	res := make([]byte, len(rl.buf.b))
	copy(res, rl.buf.b) // dunno if this is necessary
	return res, nil
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

func New() *rl {
	return &rl{}
}
