package ui

import (
	"testing"

	"github.com/diamondburned/discordgo"
)

func TestMessagePrinter(t *testing.T) {
	return

	var msgs = []discordgo.Message{
		{
			Content: "Hello, world!",
			Author: &discordgo.User{
				ID:       "170132746042081280",
				Username: "diamondburned",
				Avatar:   "894665c92101436a3d523207e1683160",
			},
		},
		{
			Content: "Second message",
			Author: &discordgo.User{
				ID:       "170132746042081280",
				Username: "diamondburned",
				Avatar:   "894665c92101436a3d523207e1683160",
			},
		},
		{
			Content: "Third message, but has\nmultiline\nto test, thonk.",
			Author: &discordgo.User{
				ID:       "170132746042081280",
				Username: "diamondburned",
				Avatar:   "894665c92101436a3d523207e1683160",
			},
		},
		{
			Content: "Ok.",
			Author: &discordgo.User{
				ID:       "302203621293162497",
				Username: "xent",
				Avatar:   "14a6121178e205defbcb3a5f54f9dc8a",
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Image: &discordgo.MessageEmbedImage{
						URL:      "https://media.discordapp.net/attachments/361916911682060289/642874709851242526/unknown.png",
						ProxyURL: "https://media.discordapp.net/attachments/361916911682060289/642874709851242526/unknown.png",
						Width:    125,
						Height:   200,
					},
				},
			},
		},
	}

	p := NewMessagePrinter(nil, nil)
	handler := p.GetHandler()

	for _, m := range msgs {
		handler(nil, &discordgo.MessageCreate{&m})
	}
}
