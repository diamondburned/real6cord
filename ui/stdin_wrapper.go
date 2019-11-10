package ui

import (
	"bufio"
	"io"
	"os"
	"sync"
	"syscall"
)

type stdinwrapper struct {
	mu     sync.Mutex
	frozen bool
	reader *bufio.Reader
	cont   chan struct{}
	done   chan struct{}
	paused bool

	in  io.Reader
	out io.Writer
}

var stdin *stdinwrapper

func init() {
	in, err := os.Open("/dev/tty")
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile("/dev/tty", syscall.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}

	stdin = &stdinwrapper{
		reader: bufio.NewReader(in),
		cont:   make(chan struct{}),
		in:     in,
		out:    out,
	}
}

func (s *stdinwrapper) Freeze() func() {
	s.frozen = true
	s.out.Write([]byte{byte(0)})
	<-s.done
	s.paused = true

	return func() {
		s.frozen = false
		s.paused = false
		s.cont <- struct{}{}
	}
}

func (s *stdinwrapper) Read(b []byte) (int, error) {
	n, err := s.in.Read(b)
	if err != nil {
		return n, err
	}

	if !s.paused && s.frozen {
		s.done <- struct{}{}
		<-s.cont
	}

	return n, nil
}
