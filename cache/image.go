package cache

import (
	"bytes"
	"image"
	"io"
	"net/http"
	"sync"

	"github.com/mattn/go-sixel"
)

func SIXELNoDithering(enc sixel.Encoder) {
	enc.Dither = false
}

type ImageOption func(img image.Image) image.Image

type ImageStore struct {
	store   map[string]*Image
	storeMu sync.Mutex

	EncodeOptions []func(enc *sixel.Encoder)
	ImageOptions  []ImageOption
}

type Image struct {
	Original image.Image
	SIXEL    []byte
}

func NewImageStore() *ImageStore {
	return &ImageStore{
		store: map[string]*Image{},
	}
}

func (i *ImageStore) AddImageOptions(f ...ImageOption) {
	i.ImageOptions = append(i.ImageOptions, f...)
}

func (i *ImageStore) download(id string, r io.Reader, imgOpts []ImageOption) (*Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	for _, o := range append(i.ImageOptions, imgOpts...) {
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

	image := &Image{
		Original: img,
		SIXEL:    bytes,
	}

	i.store[id] = image
	return image, nil
}

func (i *ImageStore) get(id string) *Image {
	i.storeMu.Lock()
	defer i.storeMu.Unlock()

	if b, ok := i.store[id]; ok {
		return b
	}

	return nil
}

func (i *ImageStore) Download(url string, imageOpts ...ImageOption) (*Image, error) {
	if b := i.get(url); b != nil {
		return b, nil
	}

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return i.download(url, r.Body, imageOpts)
}
