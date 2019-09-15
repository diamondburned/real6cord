package ui

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

var tty *os.File
var winch chan os.Signal

var ResizeHandlers []func(sz *Size)

type Size struct {
	Row    int
	Col    int
	Xpixel int
	Ypixel int
}

func initty() {
	f, err := os.Open("/dev/tty")
	if err != nil {
		panic(err)
	}

	tty = f
}

func AddToResizeHandlers(f func(*Size)) error {
	sz, err := getSize()
	if err != nil {
		return err
	}

	if winch == nil {
		winch = make(chan os.Signal)
		signal.Notify(winch, syscall.SIGWINCH)
		go pollWINCH(winch)
	}

	ResizeHandlers = append(ResizeHandlers, f)
	f(sz)

	return nil
}

func pollWINCH(winch chan os.Signal) {
	for range winch {
		sz, err := getSize()
		if err != nil {
			continue
		}

		for _, h := range ResizeHandlers {
			h(sz)
		}
	}
}

func getSize() (*Size, error) {
	initty()

	var dim [4]uint16
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL,
		uintptr(tty.Fd()), uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&dim)), 0, 0, 0); err != 0 {

		return nil, err
	}

	return &Size{
		Row:    int(dim[0]),
		Col:    int(dim[1]),
		Xpixel: int(dim[2]),
		Ypixel: int(dim[3]),
	}, nil
}

func (s *Size) CalculateCharSize() (int, int) {
	return int(s.Xpixel / s.Col), int(s.Ypixel / s.Row)
}
