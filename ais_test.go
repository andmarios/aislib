package ais

import (
	"fmt"
	"testing"
)

func TestAisPosition(t *testing.T) {
	cases := []struct {
		payload   string
		want AisPositionMessage
	} {
		{
			"38u<a<?PAA2>P:WfuAO9PW<P0PuQ",
			AisPositionMessage{Type: 3, Repeat: 0, MMSI: 601041200, Status: 15, Turn: -127, Speed: 8.1, Accuracy: false, Lon: 31.130165, Lat: -29.784113333333334},
		},
		{
			"13P:v?h009Ogbr4NkiITkU>L089D",
			AisPositionMessage{Type: 1, Repeat: 0, MMSI: 235060799, Status: 0, Turn: 0, Speed: 0.9, Accuracy: false, Lon: -3.56725, Lat: 53.84251666666667},
		},
	}
	for _, c := range cases {
		got := DecodeAisPosition(c.payload)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("DecodeAisPosition(payload string)")
		}
	}
}

func TestCoordinatesDeg2Human(t *testing.T) {
	cases := []struct {
		lon, lat float64
		want string
	} {
		{-3.56725, 53.84251666666667, "  3°34.0350'W  53°50.5510N"},
		{31.130165, -29.784113333333334,  " 31°07.8099'E  29°47.0468S"},
	}
	for _, c := range cases {
		got := CoordinatesDeg2Human(c.lon, c.lat)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("CoordinatesDeg2Human(lon, lat float64)")
		}
	}
}


