package ui

import (
	"fmt"
	"io"
)

func MoveCursorUp(w io.Writer, line int) {
	fmt.Fprintf(w, "\033[%dA", line)
}

func MoveCursorDown(w io.Writer, line int) {
	fmt.Fprintf(w, "\033[%dB", line)
}

func MoveCursorLeft(w io.Writer, line int) {
	fmt.Fprintf(w, "\033[%dD", line)
}

func MoveCursorRight(w io.Writer, line int) {
	fmt.Fprintf(w, "\033[%dC", line)
}

func MoveCursorLine(w io.Writer, line int) {
	fmt.Fprintf(w, "\033[%dH", line)
}

func MoveCursorToLineStart(w io.Writer) {
	MoveCursorLeft(w, 9999)
}

func MoveCursorTo(w io.Writer, row, col int) {
	fmt.Fprintf(w, "\033[%d;%d", row, col)
}

func ClearLine(w io.Writer) {
	io.WriteString(w, "\033[2K")
}

func HideCursor(w io.Writer) {
	io.WriteString(w, "\033[?25l")
}

func ShowCursor(w io.Writer) {
	io.WriteString(w, "\033[?25h")
}

func write(w io.Writer, s string) {
	w.Write([]byte(s))
}
