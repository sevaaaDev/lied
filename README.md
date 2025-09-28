# Lied - Line Editor

A simple line-based text editor that can do basic things
- write text (of course)
- substitute regexp (s/re/repl/)
- delete line
- copy line
- [ ] move line (you can currently achive this by doing copy and delete)
- [ ] join line

## Quick Tutorial
The editor will automatically put you in append mode, where everything you enter will be appended to the current line.
To enter a command, prefix it with ":" (colon).
Each command may accept line range prefix, which can be written in this format `start,end`.

Here are the list of commands with its default range prefix:
- `.,.p` = print
- `.,.d` = delete
- `1,$w` FILE = save to FILE
- `q`    = quit
- `.,.s` = substitute
- `.,.t.`= transfer line to target

. (dot) means current line, $ (dollar) means last line. 

The `start` and `end` range defaults to 1 and $ respectively when omitted. 
Therefore, `3,p` means line 3 until last line, and `,3p` means line 1 until line 3.
A command containing only line range will set the current line to the end range (e.g. `3,4` set current line to line 4).
A range may accept a single number without comma, which causes the start and end range to be that same number (e.g. `3p` is equal to `3,3p`)

example:
- `:,p` print whole buffer
- `:2d` delete line 2
- `:s/moon/sun` change word 'moon' to 'sun' in current line
- `:w main.go` save to file main.go
- `:42` set current line to line 42
- `:1,5t10` copy line 1-5 and paste it to line 10

## Install
you can install from source
1. clone this repo
2. `git checkout v0.2.0`
3. `go install`
