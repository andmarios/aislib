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

func TestDecodeStaticVoyageData(t *testing.T) {
	caseTime1, _ := time.Parse("1/2 15:4", "3/11 21:15")
	caseTime2, _ := time.Parse("1/2 15:4", "0/0 0:0")
	cases := []struct {
		payload string
		want    StaticVoyageData
	}{
		{
			"53uJur01rN?U<9@T001@tI@F000000000000000l0pA444mm?:1km1@SlQp000000000000",
			StaticVoyageData{
				Repeat: 0, MMSI: 265731560, AisVersion: 0, IMO: 8026361, Callsign: "SBTI",
				VesselName: "TOFTE", ShipType: 52, ToBow: 7, ToStern: 17, ToPort: 4, ToStarboard: 4,
				EPFD: 1, ETA: caseTime1, Draught: 40, Destination: "GOTEBORG", DTE: false,
			},
		},
		{
			"53m`0o400000hKGCON18E<=DF0:1",
			StaticVoyageData{
				Repeat: 0, MMSI: 257556700, AisVersion: 1, IMO: 0, Callsign: "LF5477",
				VesselName: "RESCUE B", ShipType: 0, ToBow: 0, ToStern: 0, ToPort: 0, ToStarboard: 0,
				EPFD: 0, ETA: caseTime2, Draught: 0, Destination: "", DTE: false,
			},
		},
	}
	for _, c := range cases {
		got, _ := DecodeStaticVoyageData(c.payload)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("DecodeStaticVoyageData(payload string)")
		}
	}
}

func BenchmarkDecodeStaticVoyageData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecodeStaticVoyageData("53uJur01rN?U<9@T001@tI@F000000000000000l0pA444mm?:1km1@SlQp000000000000")
	}
}
