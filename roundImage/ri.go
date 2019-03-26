package crop

import (
	"image"
	"image/color"
	"image/draw"
)

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(
		c.p.X-c.r,
		c.p.Y-c.r,
		c.p.X+c.r,
		c.p.Y+c.r,
	)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

// Round round-crops an image
func Round(src image.Image) image.Image {
	r := src.Bounds().Dx() / 2

	var dst = image.NewRGBA(image.Rect(
		0, 0,
		r*2, r*2,
	))

	draw.DrawMask(
		dst,
		src.Bounds(),
		src,
		image.ZP,
		&circle{
			p: image.Point{X: r, Y: r},
			r: r,
		},
		image.ZP,
		draw.Src,
	)

	return image.Image(dst)
}
