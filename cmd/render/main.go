package main

import (
	"encoding/json"
	"flag"
	"image"
	"io/ioutil"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/gnoack/wacomrender/imgwriter"
	"github.com/gnoack/wacomrender/tuhi"
)

func ggPt(pt tuhi.Point) (float64, float64) {
	return float64(pt.Position[0]), float64(pt.Position[1])
}

var (
	outfile    = flag.String("o", "out.png", "output file")
	width      = flag.Float64("width", 800.0, "image width to scale to")
	scaleToFit = flag.Bool("fit", false, "scale image to fit?")
	format     = flag.String("format", "png", "output format")
)

// return ctx, scale
func makeCtxWithScaling(f tuhi.File) (*gg.Context, float64) {
	if *scaleToFit {
		r := f.Bounds()

		scale := *width / float64(r.Dx())
		dc := gg.NewContext(int(float64(r.Dx())*scale), int(float64(r.Dy())*scale))
		dc.Scale(scale, scale)
		dc.Translate(-float64(r.Min.X), -float64(r.Min.Y))

		return dc, scale
	} else {
		w, h := f.Size()
		scale := *width / float64(w)
		dc := gg.NewContext(int(float64(w)*scale), int(float64(h)*scale))
		dc.Scale(scale, scale)
		return dc, scale
	}
}

// Set pressure sensitivity attributes;
// Line width depends on pressure and overall image scale.
func setPressure(dc *gg.Context, pressure float64, scale float64) {
	// TODO: Opacity scales linearly for pressures below the threshold,
	// and is fully opaque above.  (Not enabled because line drawing is shit)

	// v := 1.0
	// const threshold = 18000
	// if pressure < threshold {
	// 	v = pressure / threshold
	// }
	//dc.SetRGBA(0, 0, 0, v)
	dc.SetLineWidth(pressure * scale * 0.007)
}

func newWriter(img image.Image) imgwriter.W {
	switch *format {
	case "gif":
		return imgwriter.NewGif(img)
	case "png":
		return imgwriter.NewPng(img)
	default:
		log.Fatalf("invalid image format %v", *format)
		return nil
	}
}

func main() {
	flag.Parse()

	var f tuhi.File
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Can't read input: %v", err)
	}
	if err = json.Unmarshal(buf, &f); err != nil {
		log.Fatalf("Can't unmarshal: %v", err)
	}

	dc, scale := makeCtxWithScaling(f)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)

	imgw := newWriter(dc.Image())

	for _, stroke := range f.Strokes {
		prev := stroke.Points[0]
		for _, pt := range stroke.Points[1:] {
			setPressure(dc, pt.Pressure, scale)
			dc.MoveTo(ggPt(prev))
			dc.LineTo(ggPt(pt))
			dc.Stroke()
			dc.ClosePath()

			prev = pt
		}
		imgw.Step()
	}

	w, err := os.Create(*outfile)
	if err != nil {
		log.Fatalf("Could not open output file: %v", err)
	}

	imgw.Write(w)
	if err != nil {
		log.Fatalf("Could not write output file: %v", err)
	}
}
