package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"unicode"
)

func printBuf(buf [][]byte) {
	for i, l := range buf {
		fmt.Printf("%2d|%s\n", i+1, string(l))
	}
}

func parseCommand(input []byte) (int, error) {
	var buf []byte
	for _, v := range input {
		if !unicode.IsDigit(rune(v)) {
			return 0, fmt.Errorf("found non digit")
		}
		buf = append(buf, v)
	}
	i, _ := strconv.Atoi(string(buf))
	return i, nil
}

func readFile(filename string) [][]byte {
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
	}
	scanner := bufio.NewScanner(file)
	var buf [][]byte
	for scanner.Scan() {
		buf = append(buf, scanner.Bytes())
	}
	return buf
}

func writeFile(filename string, buf [][]byte) {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
	}
	for _, v := range buf {
		file.Write(v)
	file.Write([]byte{10})
	}
	file.Sync()
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("need file")
		os.Exit(1)
	}
	buf := readFile(os.Args[1])
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)
	index := 0
	for {
		var line []byte
		fmt.Printf("  |")
		for scanner.Scan() {
			if scanner.Bytes()[0] == '\n' {
				break
			}
			line = append(line, scanner.Bytes()...)
		}
		linenum, err := parseCommand(line[1:])
		if err != nil {
			if string(line) == ":q" {
				break
			}
			if string(line) == ":p" {
				printBuf(buf)
				continue
			}
		} else {
			index = linenum
			continue
		}
		buf = slices.Insert(buf, index, line)
		index++
	}
	printBuf(buf)
	writeFile(os.Args[1], buf)
}
