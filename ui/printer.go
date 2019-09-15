package ui

import (
	"io"
	"sync"
)

var printMu sync.Mutex

type PrinterMu struct {
	w io.Writer
}

func NewPrinterMu(w io.Writer) *PrinterMu {
	return &PrinterMu{w}
}

func (p *PrinterMu) Write(b []byte) (int, error) {
	printMu.Lock()
	defer printMu.Unlock()

	return p.w.Write(b)
}
