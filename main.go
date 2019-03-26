package main

import (
	"bufio"
	"flag"
	"image"
	"net/http"
	"os"

	_ "image/png"

	"gitlab.com/diamondburned/real6cord/roundImage"
	"github.com/diamondburned/discordgo"
	"github.com/mattn/go-sixel"
)

func main() {
	var (
		token = flag.String("t", "", "Token")
		ch    = flag.Int("ch", 0, "Channel ID")

		avatars = make(map[int64]*image.Image)
	)

	flag.Parse()

	if *token == "" || *ch == 0 {
		panic("Invalid flags, check -h")
	}

	d, err := discordgo.New(*token)
	if err != nil {
		panic(err)
	}

	if err := d.Open(); err != nil {
		panic(err)
	}

	var ready = make(chan struct{})

	d.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		ready <- struct{}{}
	})

	<-ready

	msgs, err := d.ChannelMessages(int64(*ch), 8, 0, 0, 0)
	if err != nil {
		panic(err)
	}

	for _, m := range msgs {
		if _, ok := avatars[m.Author.ID]; ok {
			continue
		}

		r, err := http.Get(m.Author.AvatarURL("32"))
		if err != nil {
			println(err.Error())
			continue
		}

		println("Downloaded", m.Author.Username)

		img, _, err := image.Decode(r.Body)
		if err != nil {
			r.Body.Close()
			println(err.Error())
			continue
		}

		r.Body.Close()

		c := crop.Round(img)

		avatars[m.Author.ID] = &c
	}

	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()

	for _, i := range avatars {
		enc := sixel.NewEncoder(buf)
		enc.Dither = false

		if err := enc.Encode(*i); err != nil {
			panic(err)
		}
	}
}
