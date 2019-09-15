package ui

import (
	"image"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/diamondburned/discordgo"
	"github.com/disintegration/imaging"
	"github.com/mitchellh/go-wordwrap"
	"gitlab.com/diamondburned/real6cord/cache"
)

// ImageLines sets the image to be 3 lines big. This means the minimum padding
// should be 4 lines away. Here's how it works.
//
// First, we print the author name. That's one line. As we print the content,
// that's n more lines. We should now print the new lines for the next message.
// Where? To calculate that, we get max(ImageLines - 1, n + 1)
const ImageLines = 3

const TextLeftPadding = 2 // characters to pad from the avatar

type MessagePrinter struct {
	Current *discordgo.Channel
	Stdout  io.Writer

	avatars *cache.Avatar
	users   *cache.Users
	printMu sync.Mutex
	lastAu  string

	termSz   *Size
	cW, cH   int
	imgSz    int
	imgChars int

	dg *discordgo.Session
}

func NewMessagePrinter(dg *discordgo.Session) *MessagePrinter {
	sz, err := getSize()
	if err != nil {
		panic(err)
	}

	w, h := sz.CalculateCharSize()

	imgSz := ImageLines * h
	imgChars := imgSz / w

	println(imgChars)

	avatars := cache.NewAvatarStore()
	avatars.ImageOptions = append(avatars.ImageOptions, func(img image.Image) image.Image {
		return imaging.Fit(img, imgSz, imgSz, imaging.Linear)
	})

	m := &MessagePrinter{
		dg:       dg,
		Stdout:   NewStatefulPrinter(os.Stdout),
		avatars:  avatars,
		users:    cache.NewUserCache(dg),
		termSz:   sz,
		cW:       w,
		cH:       h,
		imgSz:    imgSz,
		imgChars: imgChars,
	}

	err = AddToResizeHandlers(func(sz *Size) {
		m.termSz = sz
		// TODO: redraw
	})

	if err != nil {
		panic(err)
	}

	return m
}

func (p *MessagePrinter) obtLock() func() {
	p.printMu.Lock()
	write(p.Stdout, "\033[?25l")
	return func() {
		write(p.Stdout, "\033[?25h")
		p.printMu.Unlock()
	}
}

func (p *MessagePrinter) GetHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	// This variable is used to calculate the padding when a new author needs to
	// be drawn.
	// Say for example, I have printed 1 line of author, 1 line of content
	// and the image is 3 lines high. I would need to print 2 new lines. If
	// I have 2 lines of content, I only need 1. If I have 3 or more, I
	// would still need 1.
	var paddingCounter int

	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		defer p.obtLock()()

		if p.lastAu != m.Author.ID {
			// Write the padding
			write(p.Stdout, strings.Repeat("\n", paddingCounter))

			b, err := p.avatars.DownloadAvatar(m.Author)
			if err != nil {
				panic(err)
			}

			// Add 4 lines,
			write(p.Stdout, strings.Repeat("\n", ImageLines+1))

			// Move the cursor back up that line
			MoveCursorUp(p.Stdout, ImageLines+1)

			// Print the SIXEL image, which should move the cursor to the 4th
			// line
			p.Stdout.Write(b)

			// Move the cursor back 3 lines up and to start of line.
			MoveCursorToLineStart(p.Stdout)
			MoveCursorUp(p.Stdout, ImageLines)

			// At this point, we should be on the image's left corner.

			// Get the Author information
			username, color := p.users.DiscordThis(m.Message)

			// Print the author's name
			MoveCursorRight(p.Stdout, p.imgChars+TextLeftPadding)
			write(p.Stdout, ColorString(GetRGBInt(color), username))

			// Move the cursor to the start of next line
			MoveCursorDown(p.Stdout, 1)
			MoveCursorToLineStart(p.Stdout)

			// Let's move the cursor to where we could draw the message.

			// Oh, this first.
			p.lastAu = m.Author.ID
		}

		// Wrap the content
		wrapped := wordwrap.WrapString(m.Content,
			uint(int(p.termSz.Col)-p.imgChars-TextLeftPadding)-1) // safety measure

		// Split the warpped content into lines
		lines := strings.Split(wrapped, "\n")

		// Print all those lines
		for _, line := range lines {
			// Move the cursor to the right and print this line
			MoveCursorRight(p.Stdout, p.imgChars+TextLeftPadding)
			write(p.Stdout, line+"\n") // new line should reset cursor to start of line
		}

		paddingCounter = max(1, ImageLines-len(lines))
	}
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
