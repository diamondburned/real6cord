package cache

import (
	"sync"

	"github.com/diamondburned/discordgo"
)

// User is used for one user
type User struct {
	ID      string
	Discrim string
	Name    string
	Nick    string
	Color   int64
}

// Users stores multiple users
type Users struct {
	sync.RWMutex
	Guilds map[string]UserStoreArray

	dg *discordgo.Session
}

func NewUserCache(dg *discordgo.Session) *Users {
	return &Users{
		Guilds: map[string]UserStoreArray{},
		dg:     dg,
	}
}

// Populated returns a bool on whether or not the array
// alraedy is populated
func (s *Users) Populated(guildID string) bool {
	if s == nil {
		return false
	}

	if guildID == "" {
		return true
	}

	return len(s.Guilds[guildID]) > 0
}

// InStore checks if a user is in the store
func (s *Users) InStore(guildID, id string) bool {
	if s == nil {
		return false
	}

	if _, u := s.GetUser(guildID, id); u != nil {
		return true
	}

	return false
}

// DiscordThis interfaces with DiscordGo
func (s *Users) DiscordThis(m *discordgo.Message) (n string, c int64) {
	n = "invalid user"
	c = 16777215

	if m.Author == nil || s == nil {
		return
	}

	if m.GuildID == "" {
		channel, err := s.dg.State.Channel(m.ChannelID)
		if err != nil {
			channel, err = s.dg.Channel(m.ChannelID)
		}
		if err != nil {
			return
		}

		if channel.GuildID == "" {
			return m.Author.Username, c
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

	nick, color := s.getUserData(m.Author, m.ChannelID)
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
func (s *Users) GetUser(guildID, id string) (int, *User) {
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
func (s *Users) RemoveUser(guildID, id string) {
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
func (s *Users) UpdateUser(guildID, id, name, nick, discrim string, color int64) {
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
