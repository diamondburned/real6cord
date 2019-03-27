package main

import (
	"log"
	"sort"
	"sync"

	"github.com/diamondburned/discordgo"
	"github.com/rivo/tview"
)

// User is used for one user
type User struct {
	ID      int64
	Discrim string
	Name    string
	Nick    string
	Color   int
}

// UserStore stores multiple users
type UserStore struct {
	sync.RWMutex
	Guilds map[int64]UserStoreArray
}

// UserStoreArray is an array
type UserStoreArray []User

var us = &UserStore{
	Guilds: map[int64]UserStoreArray{},
}

// Populated returns a bool on whether or not the array
// alraedy is populated
func (s *UserStore) Populated(guildID int64) bool {
	if s == nil {
		return false
	}

	if guildID == 0 {
		return true
	}

	return len(s.Guilds[guildID]) > 0
}

// InStore checks if a user is in the store
func (s *UserStore) InStore(guildID, id int64) bool {
	if s == nil {
		return false
	}

	if _, u := s.GetUser(guildID, id); u != nil {
		return true
	}

	return false
}

// DiscordThis interfaces with DiscordGo
func (s *UserStore) DiscordThis(m *discordgo.Message) (n string, c int) {
	n = "invalid user"
	c = 16777215

	if m.Author == nil || s == nil {
		return
	}

	if m.GuildID == 0 {
		channel, err := d.State.Channel(m.ChannelID)
		if err != nil {
			return
		}

		m.GuildID = channel.GuildID
	}

	_, user := s.GetUser(m.GuildID, m.Author.ID)
	if user != nil {
		n = user.Name
		c = user.Color

		if user.Nick != "" {
			n = user.Nick
		}

		return
	}

	nick, color := getUserData(m.Author, m.ChannelID)
	s.UpdateUser(
		m.GuildID,
		m.Author.ID,
		m.Author.Username,
		nick,
		m.Author.Discriminator,
		color,
	)

	n = m.Author.Username
	c = color

	if nick != "" {
		n = nick
	}

	return
}

// GetUser returns the index and user for that ID
func (s *UserStore) GetUser(guildID, id int64) (int, *User) {
	s.RLock()
	defer s.RUnlock()

	if v, ok := s.Guilds[guildID]; ok {
		for i, u := range v {
			if u.ID == id {
				return i, &u
			}
		}
	}

	return 0, nil
}

// RemoveUser removes the user from the store
func (s *UserStore) RemoveUser(guildID, id int64) {
	var index int

	s.Lock()
	defer s.Unlock()

	if v, ok := s.Guilds[guildID]; ok {
		for i, u := range v {
			if u.ID == id {
				index = i
				goto Remove
			}
		}
	}

	return

Remove:
	var st = s.Guilds[guildID]

	st[len(st)-1], st[index] = st[index], st[len(st)-1]
	s.Guilds[guildID] = st[:len(st)-1]
}

// UpdateUser updates an user
func (s *UserStore) UpdateUser(guildID, id int64, name, nick, discrim string, color int) {
	if s == nil {
		return
	}

	if i, u := s.GetUser(guildID, id); u != nil {
		if name != "" {
			u.Name = name
		}

		if nick != "" {
			u.Nick = nick
		}

		if discrim != "" {
			u.Discrim = discrim
		}

		if color > 0 {
			u.Color = color
		}

		s.Lock()
		defer s.Unlock()

		s.Guilds[guildID][i] = *u
	} else {
		s.Lock()
		defer s.Unlock()

		s.Guilds[guildID] = append(s.Guilds[guildID], User{
			ID:      id,
			Discrim: discrim,
			Name:    name,
			Nick:    nick,
			Color:   color,
		})
	}
}

func safeAuthor(u *discordgo.User) (string, int64) {
	if u != nil {
		return u.Username, u.ID
	}

	return "invalid user", 0
}

func getUserData(u *discordgo.User, chID int64) (name string, color int) {
	color = 16711422
	name, id := safeAuthor(u)

	if d == nil {
		return
	}

	channel, err := d.State.Channel(chID)
	if err != nil {
		if channel, err = d.Channel(chID); err != nil {
			log.Println(err)
			return
		}
	}

	if channel.GuildID == 0 {
		return
	}

	member, err := d.State.Member(channel.GuildID, id)
	if err != nil {
		if member, err = d.GuildMember(channel.GuildID, id); err != nil {
			log.Println(err)
			return
		}
	}

	if member.Nick != "" {
		name = tview.Escape(member.Nick)
	}

	color = getUserColor(channel.GuildID, member.Roles)

	return
}

func getUserColor(guildID int64, rls discordgo.IDSlice) int {
	g, err := d.State.Guild(guildID)
	if err != nil {
		if g, err = d.Guild(guildID); err != nil {
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
				return role.Color
			}
		}
	}

	return 16711422
}
