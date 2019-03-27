package main

import (
	"bytes"
	"flag"
	"image"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"

	_ "image/png"

	"github.com/diamondburned/discordgo"
	"github.com/disintegration/imaging"
	"github.com/mattn/go-sixel"
	crop "gitlab.com/diamondburned/real6cord/roundImage"
	"gitlab.com/diamondburned/real6cord/tui"
)

var (
	avatars   = make(map[int64][]byte)
	avatarsMu sync.Mutex

	messageBuffer = make(chan *discordgo.Message, 3)
	lastAuthor    *discordgo.User

	d *discordgo.Session
)

func main() {
	var (
		token = flag.String("t", "", "Token")
		ch    = flag.Int("ch", 0, "Channel ID")
	)

	flag.Parse()

	if *token == "" || *ch == 0 {
		panic("Invalid flags, check -h")
	}

	err := tui.InitializeUI()
	if err != nil {
		panic(err)
	}

	go messageRenderer()

	d, err = discordgo.New(*token)
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

	msgs, err := d.ChannelMessages(int64(*ch), 16, 0, 0, 0)
	if err != nil {
		panic(err)
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID < msgs[j].ID
	})

	for _, m := range msgs {
		messageBuffer <- m
	}

	d.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID != int64(*ch) {
			return
		}

		messageBuffer <- m.Message
	})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func getAvatar(u *discordgo.User, refresh bool) ([]byte, error) {
	if a, ok := avatars[u.ID]; ok && !refresh {
		return a, nil
	}

	r, err := http.Get(u.AvatarURL("32"))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	img, _, err := image.Decode(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body.Close()

	if img.Bounds().Dx() != 32 {
		img = imaging.Resize(img, 32, 32, imaging.Linear)
	}

	c := crop.Round(img)

	var b bytes.Buffer

	enc := sixel.NewEncoder(&b)
	enc.Dither = false

	if err := enc.Encode(c); err != nil {
		return nil, err
	}

	bytes := b.Bytes()

	avatarsMu.Lock()
	defer avatarsMu.Unlock()

	avatars[u.ID] = bytes
	return bytes, nil
}
