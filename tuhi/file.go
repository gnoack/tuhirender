// Package tuhi defines the Tuhi data format.
//
// The tuhi.File type can be serialized and deserialized using the JSON module.
package tuhi

import (
	"image"
	"math"
)

// TODO: Implement UnmarshalJSON so fields can be named X and Y.
type Point struct {
	Position []int
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
			minx = math.Min(minx, float64(p.Position[0]))
			miny = math.Min(miny, float64(p.Position[1]))
			maxx = math.Max(maxx, float64(p.Position[0]))
			maxy = math.Max(maxy, float64(p.Position[1]))
		}
	}

	return image.Rect(int(minx), int(miny), int(maxx), int(maxy))
}
