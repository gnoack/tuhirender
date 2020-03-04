// Package tuhi defines the Tuhi data format.
//
// The tuhi.File type can be marshalled with encoding/json.
package tuhi

import (
	"image"
	"math"
)

type Point struct {
	Position image.Point
	Pressure float64
}

type Stroke struct {
	Points []Point
}

type File struct {
	Version    int
	Devicename string
	Sessionid  string
	Dimensions []int
	Timestamp  int
	Strokes    []Stroke
}

func (f File) Bounds() image.Rectangle {
	return image.Rect(0, 0, f.Dimensions[0], f.Dimensions[1])
}

func (f File) DrawingBounds() image.Rectangle {
	// Find min and max x and y values in all points.
	minx := math.Inf(1)
	miny := math.Inf(1)
	maxx := math.Inf(-1)
	maxy := math.Inf(-1)

	for _, s := range f.Strokes {
		for _, p := range s.Points {
			minx = math.Min(minx, float64(p.Position.X))
			miny = math.Min(miny, float64(p.Position.Y))
			maxx = math.Max(maxx, float64(p.Position.X))
			maxy = math.Max(maxy, float64(p.Position.Y))
		}
	}

	return image.Rect(int(minx), int(miny), int(maxx), int(maxy))
}
