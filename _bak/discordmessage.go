package main

import (
	"strings"

	"github.com/diamondburned/discordgo"
	"github.com/eidolon/wordwrap"
	"gitlab.com/diamondburned/real6cord/tui"
	"gitlab.com/diamondburned/real6cord/vt"
)

const contentPadding = 7

type message struct {
	*discordgo.Message

	DrawAuthor *drawAuthor
}

type drawAuthor struct {
	Role int
	Name string
}

func messageRenderer() {
	for m := range messageBuffer {
		var dA *drawAuthor

		if lastAuthor == nil || (lastAuthor != nil && lastAuthor.ID != m.Author.ID) {
			lastAuthor = m.Author

			username, color := us.DiscordThis(m)
			dA = &drawAuthor{
				Role: color,
				Name: username,
			}
		}

		drawMessage(m, dA)
	}
}

var authorCounter int

func drawMessage(m *discordgo.Message, d *drawAuthor) {
	s, _ := tui.GetScreen()

	switch {
	case d != nil:
		a, err := getAvatar(m.Author, false)
		if err != nil {
			panic(err)
		}

		vt.MoveCursorLine(int(s.Row))

		if authorCounter == 1 {
			authorCounter = 0
		}

		print(strings.Repeat("\n", 4+authorCounter))
		vt.MoveCursorUp(4)

		print(string(a))

		vt.MoveCursorLine(int(s.Row))
		vt.MoveCursorUp(4)
		vt.MoveCursorRight(contentPadding)

		print(vt.ColorString(
			vt.GetRGBInt(int64(d.Role)),
			d.Name,
		))

		authorCounter = 0
	default:
		authorCounter++
	}

	vt.MoveCursorDown(1)
	vt.MoveCursorToLineStart()

	wrapper := wordwrap.Wrapper(int(s.Col)-contentPadding-1, true)
	for _, l := range strings.Split(wrapper(m.Content), "\n") {
		vt.MoveCursorRight(contentPadding)
		print(l)
	}
}
