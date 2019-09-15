package ui

import (
	"io"
	"os"

	"github.com/benpye/readline"
	"github.com/diamondburned/discordgo"
)

type CLIContext struct {
	*readline.Instance
	dg *discordgo.Session
}

var grl *readline.Instance

func NewCLI(dg *discordgo.Session) (*CLIContext, error) {
	initty()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            ">> ",
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		UniqueEditLine:    true,
		Stdin:             os.Stdin,
		Stdout:            NewPrinterMu(os.Stdout),
	})

	if err != nil {
		return nil, err
	}

	grl = rl

	return &CLIContext{
		Instance: rl,
		dg:       dg,
	}, nil
}

// TODO: singleton for channelID
func (c *CLIContext) Start(chID string) {
	defer c.Close()

ReadLoop:
	for {
		line, err := c.Readline()
		switch err {
		case readline.ErrInterrupt:
			if len(line) == 0 {
				break ReadLoop
			}

			continue ReadLoop
		case io.EOF:
			break ReadLoop
		}

		//go c.dg.ChannelMessageSend(chID, line)
		//break
		print(line)
	}
}
