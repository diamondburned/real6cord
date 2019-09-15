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
	d, err := discordgo.New(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	c, err := ui.NewCLI(d)
	if err != nil {
		log.Fatalln(err)
	}

	p := ui.NewMessagePrinter(d)
	handler := p.GetHandler()

	if err := d.Open(); err != nil {
		log.Fatalln(err)
	}

	defer d.Close()

	msgs, err := d.ChannelMessages(channelID, 16, "", "", "")
	if err != nil {
		log.Fatalln(err)
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID < msgs[j].ID
	})

	for _, m := range msgs {
		handler(d, &discordgo.MessageCreate{m})
	}

	out := c.Stdout()

	d.AddHandler(func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID != channelID {
			return
		}

		handler(d, m)
		out.Write([]byte{})
	})

	// Start the CLI
	c.Start(channelID)
}
