package main

import (
	"bytes"
	"image"
	"net/http"
	"testing"

	"github.com/mattn/go-sixel"
)

func BenchmarkSixel(b *testing.B) {
	r, err := http.Get("https://cdn.discordapp.com/avatars/170132746042081280/5a8a9f925f78304ac3e9d313f836cf26.png?size=32")
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(r.Body)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		var b bytes.Buffer

		enc := sixel.NewEncoder(&b)
		enc.Dither = false

		if err := enc.Encode(img); err != nil {
			panic(err)
		}

		print(b.String())
	}
}
