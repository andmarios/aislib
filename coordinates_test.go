// Copyright (c) 2015, Marios Andreopoulos.
//
// This file is part of aislib.
//
//  Aislib is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
//  Aislib is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
// along with aislib.  If not, see <http://www.gnu.org/licenses/>.

package aislib

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
