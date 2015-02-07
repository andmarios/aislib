package ais

import (
	"fmt"
	"testing"
)

// We don't decode radio field, so testing on radio field isn't really verified.
// Radio is a 19bit field which we decode as uint32. One would need to decode further
// these 19 bits but its uneccesary for now.
func TestDecodeClassAPositionReport(t *testing.T) {
	cases := []struct {
		payload string
		want    ClassAPositionReport
	}{
		{
			"38u<a<?PAA2>P:WfuAO9PW<P0PuQ",
			ClassAPositionReport{
				PositionReport: PositionReport{Type: 3, Repeat: 0, MMSI: 601041200, Speed: 8.1,
					Accuracy: false, Lon: 31.130165, Lat: -29.784113333333334, Course: 243.4,
					Heading: 230, Second: 16, RAIM: false, Radio: 135009},
				Status: 15, Turn: -127, Maneuver: 0},
		},
		{
			"13P:v?h009Ogbr4NkiITkU>L089D",
			ClassAPositionReport{
				PositionReport: PositionReport{Type: 1, Repeat: 0, MMSI: 235060799, Speed: 0.9,
					Accuracy: false, Lon: -3.56725, Lat: 53.84251666666667, Course: 123,
					Heading: 167, Second: 14, RAIM: false, Radio: 33364},
				Status: 0, Turn: 0, Maneuver: 0},
		},
		{
			"13n@oD0PB@0IRqvQj@W;EppH088t19uvPT",
			ClassAPositionReport{
				PositionReport: PositionReport{Type: 1, Repeat: 0, MMSI: 258226000, Speed: 14.4,
					Accuracy: false, Lon: 5.580478333333334, Lat: 59.0441, Course: 290.3,
					Heading: 284, Second: 12, RAIM: false, Radio: 33340},
				Status: 0, Turn: -127, Maneuver: 0},
		},
	}
	for _, c := range cases {
		got, _ := DecodeClassAPositionReport(c.payload)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("DecodeClassAPositionReport(payload string)")
		}
	}
}

func BenchmarkDecodeClassAPositionReport(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecodeClassAPositionReport("38u<a<?PAA2>P:WfuAO9PW<P0PuQ")
	}
}
