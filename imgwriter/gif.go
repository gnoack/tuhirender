package imgwriter

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
)

type gifWriter struct {
	g gif.GIF
	i image.Image
}

func toPaletted(img image.Image) *image.Paletted {
	pimg := image.NewPaletted(img.Bounds(), palette.Plan9)
	draw.FloydSteinberg.Draw(pimg, img.Bounds(), img, image.Point{})
	return pimg
}

func (gw *gifWriter) Step() {
	gw.g.Image = append(gw.g.Image, toPaletted(gw.i))
	gw.g.Delay = append(gw.g.Delay, 10)
}

func (gw *gifWriter) Write(w io.Writer) error {
	return gif.EncodeAll(w, &gw.g)
}

func NewGif(img image.Image) *gifWriter {
	return &gifWriter{i: img}
}
