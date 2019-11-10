package ui

/*
type StatefulPrinter struct {
	lastCurLine int
	lastCurCol  int

	w io.Writer
}

func NewStatefulPrinter(w io.Writer) *StatefulPrinter {
	return &StatefulPrinter{
		w: w,
	}
}

func (p *StatefulPrinter) Write(b []byte) (int, error) {
	printMu.Lock()
	defer printMu.Unlock()

		oldLine, oldCol := getCursorPos()
		defer MoveCursorTo(p.w, oldLine, oldCol)

		if p.lastCurLine > 0 && p.lastCurCol > 0 {
			MoveCursorTo(p.w, p.lastCurLine, p.lastCurCol)
		}

	return p.w.Write(b)
}

func (p *StatefulPrinter) MoveCursorUp(line int) {
	fmt.Fprintf(p, "\033[%dA", line)
}

func (p *StatefulPrinter) MoveCursorDown(line int) {
	fmt.Fprintf(p, "\033[%dB", line)
}

func (p *StatefulPrinter) MoveCursorLeft(line int) {
	fmt.Fprintf(p, "\033[%dD", line)
}

func (p *StatefulPrinter) MoveCursorRight(line int) {
	fmt.Fprintf(p, "\033[%dC", line)
}

func (p *StatefulPrinter) MoveCursorLine(line int) {
	fmt.Fprintf(p, "\033[%dH", line)
}

func (p *StatefulPrinter) MoveCursorToLineStart() {
	MoveCursorLeft(p, 9999)
}

func (p *StatefulPrinter) MoveCursorTo(row, col int) {
	fmt.Fprintf(p, "\033[%d;%d", row, col)
}
*/
