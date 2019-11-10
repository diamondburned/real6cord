package main

import (
	"log"
	"os"
	"sort"

	"github.com/diamondburned/discordgo"
	"gitlab.com/diamondburned/real6cord/ui"
)

const channelID = "361916911682060289"

func main() {
	f, err := os.OpenFile("/tmp/real6cord.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err == nil {
		defer f.Close()
		log.SetOutput(f)
	}

	d, err := discordgo.New(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	c, err := ui.NewCLI(d)
	if err != nil {
		log.Fatalln(err)
	}

	stop := make(chan struct{})
	go func() {
		c.Start()
		stop <- struct{}{}
	}()
	defer c.Close()

	p := ui.NewMessagePrinter(d, c.Instance)
	p.Stdout = ui.NewPrinterMu(os.Stdout)
	handler := p.GetHandler()

	if err := d.Open(); err != nil {
		log.Fatalln(err)
	}

	defer d.Close()

	ch, err := d.State.Channel(channelID)
	if err == nil {
		c.SetChannel(ch)
	}

	msgs, err := d.ChannelMessages(channelID, 16, "", "", "")
	if err != nil {
		log.Fatalln(err)
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID < msgs[j].ID
	})

	for _, m := range msgs {
		handler(d, &discordgo.MessageCreate{
			Message: m,
		})
	}

	d.AddHandler(func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID != channelID {
			return
		}

		handler(d, m)
	})

	<-stop
}
