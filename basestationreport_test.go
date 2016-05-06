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
	"time"
)

func TestDecodeBaseStationReport(t *testing.T) {
	caseTime1, _ := time.Parse("2006/1/2 15:4:5", "2015/2/4 0:33:51")
	caseTime2, _ := time.Parse("2006/1/2 15:4:5", "2015/2/4 0:33:50")
	cases := []struct {
		payload string
		want    BaseStationReport
	}{
		{
			"402R3KiutR0Qk156V4QQTOA00<0;",
			BaseStationReport{
				Repeat: 0, MMSI: 2655087, Time: caseTime1, Accuracy: false, Lon: 15.09579,
				Lat: 58.588368333333335, EPFD: 1, RAIM: false, Radio: 49163,
			},
		},
		{
			"4025boiutR0Qj0qgK<OodKW00@N1",
			BaseStationReport{
				Repeat: 0, MMSI: 2190047, Time: caseTime2, Accuracy: false, Lon: 12.613716666666667,
				Lat: 55.69725, EPFD: 7, RAIM: false, Radio: 67457,
			},
		},
	}
	for _, c := range cases {
		got, _ := DecodeBaseStationReport(c.payload)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("DecodeClassAPositionReport(payload string)")
		}
	}
}

func BenchmarkDecodeBaseStationReport(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecodeBaseStationReport("402R3KiutR0Qk156V4QQTOA00<0;")
	}
}

func TestGetReferenceTime(t *testing.T) {
	cases := []struct {
		payload, reference string
	}{
		{"4025;PAuho;N>0NJbfMRhNA00D3l", "2012/3/14 11:30:14"},
		{"403tDGiuho;P5<tSF0l4Q@000l67", "2012/3/14 11:32:5"},
	}
	for _, c := range cases {
		got, _ := GetReferenceTime(c.payload)
		want, _ := time.Parse("2006/1/2 15:4:5", c.reference)
		if got != want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", want)
			t.Errorf("GetReferenceTime(payload string)")
		}
	}
}

func BenchmarkGetReferenceTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetReferenceTime("4025;PAuho;N>0NJbfMRhNA00D3l")
	}
}
