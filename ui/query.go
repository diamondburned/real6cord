package ui

import (
	"fmt"
	"time"
)

// https://github.com/gizak/termui/pull/233/files#diff-61ca5d3d7b39f5b633e6774d6a31aea5
func queryTerm(qs string) (ret []rune) {
	ch := make(chan struct{})
	runes := make(chan rune)

	grl.Terminal.Steal = func(r rune) {
		runes <- r
	}

	defer func() {
		grl.Terminal.Steal = nil
	}()

	go func() {
		// query terminal
		fmt.Printf(qs)

		for r := range runes {
			// handle key event
			switch r {
			case 'c', 't', 'R':
				ret = append(ret, r)
				goto afterLoop
			default:
				ret = append(ret, r)
			}
		}
	afterLoop:
		ch <- struct{}{}
	}()

	var timer = time.NewTimer(5000 * time.Microsecond)
	defer timer.Stop()

	select {
	case <-ch:
		close(runes)
		close(ch)
	case <-timer.C:
	}

	return
}
