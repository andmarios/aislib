package ais

import (
	"fmt"
	"testing"
)

func TestCoordinatesDeg2Human(t *testing.T) {
	cases := []struct {
		lon, lat float64
		want     string
	}{
		{-3.56725, 53.84251666666667, "  3°34.0350'W  53°50.5510N"},
		{31.130165, -29.784113333333334, " 31°07.8099'E  29°47.0468S"},
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

func ExampleCoordinatesDeg2Human() {
	fmt.Println(CoordinatesDeg2Human(-3.56725, 53.84251666666667))
	// Output:   3°34.0350'W  53°50.5510N
}
