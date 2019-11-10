package ui

import (
	"fmt"
	"image"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/benpye/readline"
	"github.com/diamondburned/discordgo"
	"github.com/disintegration/imaging"
	"github.com/mitchellh/go-wordwrap"
	"gitlab.com/diamondburned/real6cord/cache"
)

// AvatarLines sets the avatar to be 3 lines big. This means the minimum padding
// should be 4 lines away. Here's how it works.
//
// First, we print the author name. That's one line. As we print the content,
// that's n more lines. We should now print the new lines for the next message.
// Where? To calculate that, we get max(ImageLines - 1, n + 1)
const AvatarLines = 3

// ImageLines sets the avatar to be at most 3 lines big. This means the minimum
const ImageLines = 8

const EmojiLines = 2

const TextLeftPadding = 2 // characters to pad from the avatar

const EmbedPill = Bold + "│" + Reset
const EmbedColor = 0x00B0F4

type MessagePrinter struct {
	Current *discordgo.Channel
	Stdout  io.Writer

	users   *cache.Users
	printMu sync.Mutex
	lastAu  string

	termSz *Size
	cW, cH int

	avatars  *cache.Avatar
	avaChars int

	images     *cache.ImageStore
	imgMaxSz   int
	thumbMaxSz int

	emojis   *cache.ImageStore
	emoChars int

	rl *readline.Instance

	dg *discordgo.Session
}

func NewMessagePrinter(dg *discordgo.Session, rl *readline.Instance) *MessagePrinter {
	m := &MessagePrinter{
		dg:     dg,
		Stdout: os.Stdout,
		users:  cache.NewUserCache(dg),
		rl:     rl,
	}

	// Get the terminal size in lines/columns and in pixels
	sz, err := getSize()
	if err != nil {
		panic(err)
	}

	m.termSz = sz

	// Calculate the width and height of a character
	m.cW, m.cH = sz.CalculateCharSize()

	// Calculate the avatar size
	avaSz := AvatarLines * m.cH
	// Calculate the avatar column size
	m.avaChars = avaSz / m.cW

	m.avatars = cache.NewAvatarStore()
	m.avatars.AddImageOptions(func(img image.Image) image.Image {
		return imaging.Fit(img, avaSz, avaSz, imaging.Linear)
	})

	// Calculate the maximum image size
	m.imgMaxSz = ImageLines * m.cH
	m.thumbMaxSz = avaSz

	m.images = cache.NewImageStore()
	// m.images.AddImageOptions(func(img image.Image) image.Image {
	//	return imaging.Fit(img, m.imgMaxSz, m.imgMaxSz, imaging.Linear)
	// })

	// Calculate the maximum emoji size
	emojiSz := EmojiLines * m.cH
	m.emoChars = emojiSz / m.cW

	m.emojis = cache.NewImageStore()
	m.emojis.AddImageOptions(func(img image.Image) image.Image {
		return imaging.Fit(img, emojiSz, emojiSz, imaging.Linear)
	})

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
	HideCursor(p.Stdout)

	return func() {
		ShowCursor(p.Stdout)
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

	// This variable tracks the current line relative from beneath the author
	// name.
	var currentLine int

	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		defer p.obtLock()()

		// We need to move the cursor up one line to make up for the new line we
		// insert at the end of this callback.
		if paddingCounter > 0 { // If handler is called more than once
			MoveCursorUp(p.Stdout, 1)

			if currentLine == 1 {
				// Current line is 1, so the cursor was moved down twice. We
				// need to move it up twice.
				MoveCursorUp(p.Stdout, 1)
			}
		}

		// If the line we moved to is outside the avatar, we'll need to clean it
		// to get rid of the prompt.
		if currentLine >= AvatarLines-1 {
			ClearLine(p.Stdout)
		}

		if p.lastAu != m.Author.ID {
			// Reset currentLine
			currentLine = 0

			// Write the padding
			for i := 0; i < paddingCounter; i++ {
				io.WriteString(p.Stdout, "\n")
				ClearLine(p.Stdout)
			}

			b, err := p.avatars.DownloadAvatar(m.Author)
			if err != nil {
				panic(err)
			}

			// Add 4 lines,
			for i := 0; i < AvatarLines+1; i++ {
				io.WriteString(p.Stdout, "\n")
				ClearLine(p.Stdout)
			}

			// Move the cursor back up that many lines
			MoveCursorUp(p.Stdout, AvatarLines+1)
			MoveCursorToLineStart(p.Stdout)

			// Print the SIXEL image, which should move the cursor to the 4th
			// line
			p.Stdout.Write(b)

			// Move the cursor back 3 lines up and to start of line.
			MoveCursorUp(p.Stdout, AvatarLines)
			MoveCursorToLineStart(p.Stdout)

			// At this point, we should be on the image's left corner.

			// Get the Author information
			username, color := p.users.DiscordThis(m.Message)

			// Print the author's name
			MoveCursorRight(p.Stdout, p.avaChars+TextLeftPadding)
			write(p.Stdout, ColorString(GetRGBInt(color), username))

			// Move the cursor to the start of next line
			MoveCursorDown(p.Stdout, 1)
			MoveCursorToLineStart(p.Stdout)

			// Let's move the cursor to where we could draw the message.

			// Oh, this first.
			p.lastAu = m.Author.ID
		}

		// Reset cursor to line start
		MoveCursorToLineStart(p.Stdout)

		// Wrap the content and split the warpped content into lines
		lines := p.wrap(m.Content)

		// If the content is just a single emoji
		if emoji := getOnlyEmojiURL(m.Content); emoji != "" {
			i, err := p.emojis.Download(emoji)
			if err == nil {
				// Add the total line height of the images into lineCount and
				// currentLine
				for i := 0; i < EmojiLines; i++ {
					currentLine++
					write(p.Stdout, "\n")

					if currentLine > AvatarLines-1 {
						ClearLine(p.Stdout)
					}
				}

				MoveCursorUp(p.Stdout, EmojiLines)

				MoveCursorRight(p.Stdout, p.avaChars+TextLeftPadding)
				p.Stdout.Write(i.SIXEL)

				lines = lines[:0]
			}
		}

		// Total lines printed
		var lineCount = len(lines)

		// Print all those lines
		for _, line := range lines {
			// Move the cursor to the right and print this line
			p.writePad(line) // new line should reset cursor to start of line

			// If the text we're drawing isn't in a line where the avatar is,
			// we'll need to clear the line of the prompt.
			if currentLine > AvatarLines-1 {
				ClearLine(p.Stdout)
			}

			currentLine++
		}

		// Clear lines to write after rendering these embeds and attachments
		lines = lines[:0]

		for _, e := range m.Embeds {
			// Offset to move horizontally and vertically, used to draw a
			// thumbnail next to text
			var horizOffset int
			var vertOffset int

			if e.Type != "image" {
				color := GetRGBInt(int64(e.Color))
				lines := 0 // lines so far
				pill := ColorString(color, EmbedPill)

				if e.Author != nil && e.Author.Name != "" {
					p.writePad(pill + "[] " + ColorString(GetRGBInt(EmbedColor),
						truncateString(e.Author.Name, 80-1)))
					lines++
				}

				// Write the title, if there's one
				if e.Title != "" {
					p.writePad(pill + ColorString(
						GetRGBInt(EmbedColor), truncateString(e.Title, 80-1)))
					lines++
				}

				// Write the description, if there is one
				if e.Description != "" {
					wrapped := p.wrapCustom(e.Description, 80-1)
					lines += len(wrapped)

					for _, line := range wrapped {
						p.writePad(pill + line)
						horizOffset = max(len(line)+1, horizOffset)
					}
				}

				// Write fields, if there are any
				if len(e.Fields) > 0 {
					p.writePad(pill)
					lines++
				}

				for _, f := range e.Fields {
					p.writePad(pill + truncateString(f.Name, 60-1))
					p.writePad(pill + "  " + truncateString(f.Value, 60-1-2))

					lines += 2
				}

				if e.Thumbnail != nil {
					// draw extra embed pills
					for i := lines; i < AvatarLines; i++ {
						p.writePad(pill)
					}
				}

				horizOffset = min(horizOffset, 80)
				vertOffset = max(AvatarLines, lines)

				currentLine += vertOffset
				lineCount += vertOffset
			}

			// Number of lines to skip after drawing
			// var skipLines = 1

			// Draw the thumbnail
			if e.Thumbnail != nil {
				var maxSize = p.thumbMaxSz
				if e.Type == "image" {
					maxSize = p.imgMaxSz
				}

				url, imgLines, w, h := p.getEmbedImageURL(
					e.Thumbnail.ProxyURL,
					e.Thumbnail.Width, e.Thumbnail.Height,
					maxSize,
				)

				i, err := p.images.Download(url, func(img image.Image) image.Image {
					return imaging.Fit(img, w, h, imaging.Linear)
				})

				if err != nil {
					// This might be slow, but it should be pretty rare
					wrappedErr := p.wrap("Failed to download " +
						path.Base(e.URL) + ": " + err.Error())

					// Add the error into the lines
					lines = append(lines, wrappedErr...)
					lineCount += len(wrappedErr)

					continue
				}

				// At this point, the cursor is at the start of a new line. We pad
				// the cursor.
				MoveCursorUp(p.Stdout, vertOffset)
				MoveCursorRight(p.Stdout, p.avaChars+TextLeftPadding+horizOffset+2)

				// Then, we write the SIXEL image
				p.Stdout.Write(i.SIXEL)

				// Add the total line height of the images into lineCount and
				// currentLine
				lineCount += imgLines
				currentLine += imgLines

				/*
					// Bring the cursor back
					MoveCursorUp(p.Stdout, imgLines)
					MoveCursorLeft(p.Stdout, leftPad)

					// Set the number of lines to skip later
					if imgLines > skipLines {
						skipLines = imgLines
					}
				*/
			}
		}

		for _, a := range m.Attachments {
			if !(a.Width > 0 && a.Height > 0) {
				continue
			}

			url, imgLines, w, h := p.getEmbedImageURL(
				a.ProxyURL, a.Width, a.Height, p.imgMaxSz,
			)

			i, err := p.images.Download(url, func(img image.Image) image.Image {
				return imaging.Fit(img, w, h, imaging.Linear)
			})

			if err != nil {
				// This might be slow, but it should be pretty rare
				wrappedErr := p.wrap("Failed to download " +
					path.Base(a.URL) + ": " + err.Error())

				// Add the error into the lines
				lines = append(lines, wrappedErr...)
				lineCount += len(wrappedErr)

				continue
			}

			MoveCursorRight(p.Stdout, p.avaChars+TextLeftPadding)
			p.Stdout.Write(i.SIXEL)

			lineCount += imgLines
			currentLine += imgLines
		}

		// Calculate the padding for the next message
		paddingCounter = max(1, AvatarLines-lineCount)

		write(p.Stdout, "\n")
		if currentLine == 1 {
			// Print another line so the prompt goes to the bottom
			write(p.Stdout, "\n")
		}

		if p.rl != nil {
			p.rl.Refresh()
		}
	}
}

func (p *MessagePrinter) getEmbedImageURL(url string, w, h, max int,
) (newURL string, lines, newW, newH int) {

	var (
		resizeW int
		resizeH int
	)

	if h > w {
		resizeH = max
		resizeW = max * w / h
	} else {
		resizeW = max
		resizeH = max * h / w
	}

	lines = resizeH / p.cH

	return strings.Split(url, "?")[0] + fmt.Sprintf(
		"?width=%d&height=%d",
		resizeW, resizeH,
	), lines, resizeW, resizeH
}

func (p *MessagePrinter) wrap(text string) []string {
	return strings.Split(wordwrap.WrapString(text,
		uint(int(p.termSz.Col)-p.avaChars-TextLeftPadding)-1), "\n") // safety measure
}

func (p *MessagePrinter) wrapCustom(text string, width int) []string {
	return strings.Split(wordwrap.WrapString(text, uint(width)), "\n")
}

func (p *MessagePrinter) writePad(text string) {
	MoveCursorRight(p.Stdout, p.avaChars+TextLeftPadding)
	write(p.Stdout, text+"\n")
}

var emojiRegex = regexp.MustCompile(`<(a?):(.+?):(\d+)>`)

func getOnlyEmojiURL(from string) string {
	emojiIDs := emojiRegex.FindAllStringSubmatch(from, -1)
	if len(emojiIDs) != 1 {
		return ""
	}

	if emojiIDs[0][1] != "" {
		// is animated, return
		return ""
	}

	if from != emojiIDs[0][0] {
		// message has more than just the emoji
		return ""
	}

	return `https://cdn.discordapp.com/emojis/` + emojiIDs[0][3] + ".png"
}

func truncateString(str string, width int) string {
	if len(str) > width {
		return str[:width-1] + "…"
	}

	return str
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
