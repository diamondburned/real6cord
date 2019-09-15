package cache

import (
	"bytes"
	"image"
	"io"
	"net/http"
	"sync"

	"github.com/diamondburned/discordgo"
	"github.com/mattn/go-sixel"
	"gitlab.com/diamondburned/real6cord/imageutil"
)

func SIXELNoDithering(enc sixel.Encoder) {
	enc.Dither = false
}

type Avatar struct {
	store   map[string][]byte
	storeMu sync.Mutex

	EncodeOptions []func(enc *sixel.Encoder)
	ImageOptions  []func(img image.Image) image.Image
}

func NewAvatarStore() *Avatar {
	return &Avatar{
		store: map[string][]byte{},
		ImageOptions: []func(img image.Image) image.Image{
			imageutil.Round,
		},
	}
}

func (i *Avatar) download(id string, r io.Reader) ([]byte, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	for _, o := range i.ImageOptions {
		img = o(img)
	}

	var b bytes.Buffer

	enc := sixel.NewEncoder(&b)
	for _, o := range i.EncodeOptions {
		o(enc)
	}

	if err := enc.Encode(img); err != nil {
		return nil, err
	}

	bytes := b.Bytes()

	i.storeMu.Lock()
	defer i.storeMu.Unlock()

	i.store[id] = bytes
	return bytes, nil
}

func (i *Avatar) get(id string) []byte {
	i.storeMu.Lock()
	defer i.storeMu.Unlock()

	if b, ok := i.store[id]; ok {
		return b
	}

	return nil
}

func (i *Avatar) DownloadAvatar(u *discordgo.User) ([]byte, error) {
	if b := i.get(u.ID); b != nil {
		return b, nil
	}

	r, err := http.Get(u.AvatarURL("64"))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return i.download(u.ID, r.Body)
}
