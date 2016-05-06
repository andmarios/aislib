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
				PositionReport: PositionReport{
					Type: 3, Repeat: 0, MMSI: 601041200, Speed: 8.1,
					Accuracy: false, Lon: 31.130165, Lat: -29.784113333333334, Course: 243.4,
					Heading: 230, Second: 16, RAIM: false, Radio: 135009},
				Status: 15, Turn: -127, Maneuver: 0},
		},
		{
			"13P:v?h009Ogbr4NkiITkU>L089D",
			ClassAPositionReport{
				PositionReport: PositionReport{
					Type: 1, Repeat: 0, MMSI: 235060799, Speed: 0.9,
					Accuracy: false, Lon: -3.56725, Lat: 53.84251666666667, Course: 123,
					Heading: 167, Second: 14, RAIM: false, Radio: 33364},
				Status: 0, Turn: 0, Maneuver: 0},
		},
		{
			"13n@oD0PB@0IRqvQj@W;EppH088t19uvPT",
			ClassAPositionReport{
				PositionReport: PositionReport{
					Type: 1, Repeat: 0, MMSI: 258226000, Speed: 14.4,
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

func TestDecodeClassBPositionReport(t *testing.T) {
	cases := []struct {
		payload string
		want    ClassBPositionReport
	}{
		{
			"B3ujWF0000DdVU8O:1H03wi5oP06",
			ClassBPositionReport{
				PositionReport: PositionReport{
					Type: 18, Repeat: 0, MMSI: 266119000, Speed: 0,
					Accuracy: false, Lon: 18.085243333333334, Lat: 59.32718333333333, Course: 0,
					Heading: 511, Second: 34, RAIM: true, Radio: 917510},
				CSUnit: true, Display: false, DSC: true, Band: true, Msg22: true, Assigned: false},
		},
		{
			"B3uIwBP008=QHv8Cerc;wwjUWP06",
			ClassBPositionReport{
				PositionReport: PositionReport{
					Type: 18, Repeat: 0, MMSI: 265715530, Speed: 0,
					Accuracy: true, Lon: 11.81546, Lat: 58.07772333333333, Course: 326.3,
					Heading: 511, Second: 37, RAIM: true, Radio: 917510},
				CSUnit: true, Display: false, DSC: true, Band: true, Msg22: false, Assigned: false},
		},
	}
	for _, c := range cases {
		got, _ := DecodeClassBPositionReport(c.payload)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("DecodeClassBPositionReport(payload string)")
		}
	}
}

func BenchmarkDecodeClassBPositionReport(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecodeClassBPositionReport("B3ujWF0000DdVU8O:1H03wi5oP06")
	}
}
