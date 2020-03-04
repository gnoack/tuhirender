package tuhi_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gnoack/wacomrender/tuhi"
)

func eq(a, b tuhi.File) bool {
	return reflect.DeepEqual(a, b)
}

func TestRoundtripMarshal(t *testing.T) {
	f := tuhi.File{
		Version:    1,
		Devicename: "my sketchpad",
		Sessionid:  "123",
		Dimensions: []int{5000, 6000},
		Timestamp:  1583349599,
		Strokes: []tuhi.Stroke{
			{[]tuhi.Point{{[]int{10, 20}, 62000}}},
			{[]tuhi.Point{{[]int{15, 30}, 62000}}},
		},
	}

	m, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("Can't marshal: %v", err)
	}

	var res tuhi.File
	err = json.Unmarshal(m, &res)
	if err != nil {
		t.Fatalf("Can't unmarshal: %v", err)
	}

	if !eq(f, res) {
		t.Errorf("Mismatching roundtrip results, got %v, want %v", res, f)
	}
}
