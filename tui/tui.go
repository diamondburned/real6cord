package tui

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

// Screen contains screen info
type Screen struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

var (
	tty *os.File

	// Global screen, always up to date
	screen *Screen

	errorUninitialized = errors.New("UI not initialized")
)

// InitializeUI starts stuff
func InitializeUI() (err error) {
	tty, err = os.Open("/dev/tty")
	if err != nil {
		return err
	}

	s, err := getSize()
	if err != nil {
		return err
	}

	screen = s

	sigwinch := make(chan os.Signal)
	signal.Notify(sigwinch, syscall.SIGWINCH)

	go func() {
		for {
			<-sigwinch

			s, err := getSize()
			if err != nil {
				panic(err)
			}

			screen = s
		}
	}()

	return nil
}

func GetScreen() (s Screen, err error) {
	if screen == nil {
		err = errorUninitialized
		return
	}

	return *screen, nil
}

func getSize() (s *Screen, err error) {
	s = new(Screen)

	_, _, eint := syscall.Syscall(syscall.SYS_IOCTL,
		tty.Fd(),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(s)),
	)

	if eint != 0 {
		err = eint
	}

	return
}
