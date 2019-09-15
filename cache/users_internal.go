package cache

import (
	"log"
	"sort"

	"github.com/diamondburned/discordgo"
	"github.com/rivo/tview"
)

// UserStoreArray is an array
type UserStoreArray []User

func safeAuthor(u *discordgo.User) (string, string) {
	if u != nil {
		return u.Username, u.ID
	}

	return "invalid user", ""
}

func (s *Users) getUserData(u *discordgo.User, chID string) (name string, color int64) {
	color = 16711422
	name, id := safeAuthor(u)

	channel, err := s.dg.State.Channel(chID)
	if err != nil {
		if channel, err = s.dg.Channel(chID); err != nil {
			return
		}
	}

	if channel.GuildID == "" {
		return
	}

	member, err := s.dg.State.Member(channel.GuildID, id)
	if err != nil {
		if member, err = s.dg.GuildMember(channel.GuildID, id); err != nil {
			return
		}
	}

	if member.Nick != "" {
		name = tview.Escape(member.Nick)
	}

	color = s.getUserColor(channel.GuildID, member.Roles)

	return
}

func (s *Users) getUserColor(guildID string, rls []string) int64 {
	g, err := s.dg.State.Guild(guildID)
	if err != nil {
		if g, err = s.dg.Guild(guildID); err != nil {
			log.Println(err)
			return 16711422
		}
	}

	roles := g.Roles
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})

	for _, role := range roles {
		for _, roleID := range rls {
			if role.ID == roleID && role.Color != 0 {
				return int64(role.Color)
			}
		}
	}

	return 16711422
}
