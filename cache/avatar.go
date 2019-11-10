package cache

import (
	"net/http"

	"github.com/diamondburned/discordgo"
	"gitlab.com/diamondburned/real6cord/imageutil"
)

type Avatar struct {
	*ImageStore
}

func NewAvatarStore() *Avatar {
	store := NewImageStore()
	store.ImageOptions = []ImageOption{
		imageutil.Round,
	}

	return &Avatar{
		ImageStore: store,
	}
}

func (i *Avatar) DownloadAvatar(u *discordgo.User) ([]byte, error) {
	if b := i.get(u.ID); b != nil {
		return b.SIXEL, nil
	}

	r, err := http.Get(u.AvatarURL("64"))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	b, err := i.download(u.ID, r.Body, nil)
	if err != nil {
		return nil, err
	}

	return b.SIXEL, nil
}
