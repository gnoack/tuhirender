package main

import (
	"encoding/json"
	"flag"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/gnoack/path"
	"github.com/gnoack/wacomrender/imgwriter"
	"github.com/gnoack/wacomrender/tuhi"
)

var (
	outfile         = flag.String("o", "out.png", "output file")
	width           = flag.Float64("width", 800.0, "image width to scale to")
	scaleToFit      = flag.Bool("fit", false, "scale image to fit?")
	border          = flag.Float64("border", 10.0, "border pixels in target image")
	format          = flag.String("format", "png", "output format")
	cycleColors     = flag.Bool("debug", false, "cycle colors for debugging")
	fixedPressure   = flag.Float64("fixed-pressure", -1.0, "to set a fixed pen pressure")
	scalePressure   = flag.Float64("scale-pressure", 1.0, "to set a scaling factor for pen pressure")
	simplify        = flag.Bool("simplify", false, "simplify strokes")
	simplifyEpsilon = flag.Float64("simplify.epsilon", 50, "epsilon for simplification algorithm")
)

var colors = []color.Color{
	color.RGBA{255, 0, 0, 255},
	color.RGBA{0, 255, 0, 255},
	color.RGBA{0, 0, 255, 255},
}

func ggPt(pt tuhi.Point) (float64, float64) {
	return float64(pt.Position.X), float64(pt.Position.Y)
}

// return ctx, scale
func makeCtxWithScaling(f tuhi.File) (*gg.Context, float64) {
	var r image.Rectangle
	if *scaleToFit {
		r = f.DrawingBounds()
	} else {
		r = f.Bounds()
	}
	scale := *width / float64(r.Dx())
	r = r.Inset(-int(*border / scale))

	dc := gg.NewContext(
		int(float64(r.Dx())*scale),
		int(float64(r.Dy())*scale),
	)
	dc.Scale(scale, scale)
	dc.Translate(-float64(r.Min.X), -float64(r.Min.Y))
	return dc, scale
}

// Set pressure sensitivity attributes;
// Line width depends on pressure and overall image scale.
func setPressure(dc *gg.Context, pressure float64, scale float64) {
	// TODO: Opacity scales linearly for pressures below the threshold,
	// and is fully opaque above.  (Not enabled because line drawing is shit)

	pressure = *scalePressure * pressure
	if *fixedPressure > 0.0 {
		pressure = *fixedPressure * 65000
	}
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

func simplifyStroke(s tuhi.Stroke) tuhi.Stroke {
	epsilon := *simplifyEpsilon * *simplifyEpsilon
	f := func(i int) (x, y int) {
		pos := s.Points[i].Position
		return pos.X, pos.Y
	}
	indices := path.Simplify(path.OfIntPoints(f, len(s.Points)), epsilon)
	var out []tuhi.Point
	for _, i := range indices {
		out = append(out, s.Points[i])
	}
	return tuhi.Stroke{Points: out}
}

func simplifyStrokes(strokes []tuhi.Stroke) []tuhi.Stroke {
	var out []tuhi.Stroke
	for _, s := range strokes {
		out = append(out, simplifyStroke(s))
	}
	return out
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

	if *simplify {
		f.Strokes = simplifyStrokes(f.Strokes)
	}

	for _, stroke := range f.Strokes {
		prev := stroke.Points[0]
		for idx, pt := range stroke.Points[1:] {
			if *cycleColors {
				// A new color for each path segment.
				dc.SetColor(colors[idx%len(colors)])
			}
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
