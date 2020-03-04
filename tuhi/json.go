package tuhi

import (
	"encoding/json"
	"errors"
	"image"
)

type jsonPoint struct {
	Position []int
	Pressure float64
}

func (p Point) MarshalJSON() ([]byte, error) {
	jp := jsonPoint{
		Position: []int{p.Position.X, p.Position.Y},
		Pressure: p.Pressure,
	}
	return json.Marshal(jp)

}

func (p *Point) UnmarshalJSON(b []byte) error {
	var jp jsonPoint
	if err := json.Unmarshal(b, &jp); err != nil {
		return err
	}

	if len(jp.Position) != 2 {
		return errors.New("Wrong number of ints in point")
	}

	p.Position = image.Pt(jp.Position[0], jp.Position[1])
	p.Pressure = jp.Pressure
	return nil
}
