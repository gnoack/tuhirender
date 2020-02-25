package imgwriter

import (
	"image"
	"image/png"
	"io"
)

type pngWriter struct {
	i image.Image
}

func (*pngWriter) Step() {}

func (pw *pngWriter) Write(w io.Writer) error {
	return png.Encode(w, pw.i)
}

func NewPng(img image.Image) *pngWriter {
	return &pngWriter{i: img}
}
