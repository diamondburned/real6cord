package vt

import "fmt"

func MoveCursorUp(line int) {
	fmt.Printf("\033[%dA", line)
}

func MoveCursorDown(line int) {
	fmt.Printf("\033[%dB", line)
}

func MoveCursorLeft(line int) {
	fmt.Printf("\033[%dD", line)
}

func MoveCursorRight(line int) {
	fmt.Printf("\033[%dC", line)
}

func MoveCursorLine(line int) {
	fmt.Printf("\033[%dH", line)
}

func MoveCursorToLineStart() {
	MoveCursorLeft(9999)
}

func MoveCursorTo(row, col int) {
	fmt.Printf("\033[%d;%d", row, col)
}
