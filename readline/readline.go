package readline

import (
	"fmt"
	"os"
	"slices"
)

const (
	KEY_NULL  = 0   /* NULL */
	CTRL_A    = 1   /* Ctrl+a */
	CTRL_B    = 2   /* Ctrl-b */
	CTRL_C    = 3   /* Ctrl-c */
	CTRL_D    = 4   /* Ctrl-d */
	CTRL_E    = 5   /* Ctrl-e */
	CTRL_F    = 6   /* Ctrl-f */
	CTRL_H    = 8   /* Ctrl-h */
	TAB       = 9   /* Tab */
	CTRL_K    = 11  /* Ctrl+k */
	CTRL_L    = 12  /* Ctrl+l */
	ENTER     = 10  /* Enter */
	CTRL_N    = 14  /* Ctrl-n */
	CTRL_P    = 16  /* Ctrl-p */
	CTRL_T    = 20  /* Ctrl-t */
	CTRL_U    = 21  /* Ctrl+u */
	CTRL_W    = 23  /* Ctrl+w */
	ESC       = 27  /* Escape */
	BACKSPACE = 127 /* Backspace */
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

func tab2space(tabLen int) string {
	b := []byte{}
	for range tabLen {
		b = append(b, ' ')
	}
	return string(b)
}

func (rl *rl) refreshLine() {
	cursorPos := rl.cursorPos
	cmd := []byte{}
	cmd = append(cmd, '\r')
	cmd = append(cmd, []byte(rl.prompt)...)
	for i, c := range rl.buf.b {
		if c != '\t' {
			cmd = append(cmd, byte(c))
			continue
		}
		cmd = append(cmd, []byte(tab2space(4))...)
		if i < rl.cursorPos {
			cursorPos += 3 // minus 1 than tab len, bcs rl.cursorPos is already ahead by 1
		}
	}
	cmd = append(cmd, []byte("\x1b[0K")...)
	cmd = fmt.Appendf(cmd, "\r\x1b[%dC", cursorPos-1+len(rl.prompt)-1)
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
	seq := make([]byte, 4)
	_, err := os.Stdin.Read(seq)
	if err != nil {
		return false
	}
	if seq[0] >= 0x20 && seq[0] < 0x7f {
		rl.insert(seq[0])
		return true
	}
	switch seq[0] {
	case ENTER:
		return false
	case BACKSPACE, 8:
		rl.backspace()
	case TAB:
		rl.insert('\t')
	case ESC:
		if seq[1] == '[' {
			if seq[2] >= '0' && seq[2] <= '9' && seq[3] == '~' {
				switch seq[2] {
				case '1': // HOME key
					rl.cursorPos = 0
				case '4': // END key
					rl.cursorPos = len(rl.buf.b)
				}
			}
			switch seq[2] {
			case 'C':
				rl.cursorRight()
			case 'D':
				rl.cursorLeft()
			case 'H':
				rl.cursorPos = 0
			case 'F':
				rl.cursorPos = len(rl.buf.b)
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
	defer rl.refreshLine()

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
