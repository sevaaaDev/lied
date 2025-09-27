# Lied - Line Editor

A simple line-based text editor that can do basic things
- write text (of course)
- substitute regexp (s/re/repl/)
- delete line
- [ ] copy/cut line
- [ ] join line

## Quick Tutorial
The editor will automatically put you in append mode, where everything you enter will be appended to the current line.
To enter a command, prefix it with ":" (colon).
Each command may accept line range prefix, similar to ed/ex.
Here are the list of commands with its default range prefix:
- 1,$p = print
- .,.d = delete
- 1,$w FILE = save to FILE
- q = quit
- .,.s = substitute

. (dot) means current line, $ (dollar) means last line. 
The start and end range defaults to 1 and $ respectively. 
Therefore, `3,p` means line 3 until last line, and `,3p` means line 1 until line 3.
A command containing only line range will set the current line to the end range (e.g. `3,4` set current line to line 4).
A range may accept a single number without comma, which causes the start and end range to be that same number (e.g. `3p` is equal to `3,3p`)

example:
- `:,p` print whole buffer
- `:2d` delete line 2
- `:s/moon/sun` change word 'moon' to 'sun' in current line
- `:w main.go` save to file main.go
- `:32` set current line to line 32

## Goal
apart from basic editing experience, i also want to add syntax highlighting. idk if its worth it
