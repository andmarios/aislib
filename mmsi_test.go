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

// I didn't found a MMSI decoder, so this test is verified only by me
func TestDecodeMMSI(t *testing.T) {
	cases := []struct {
		MMSI      uint32
		reference string
	}{
		{227006760, "Ship, France"},
		{2573425, "Coastal Station, Norway"},
		{25634906, "Group of ships, Malta"},
		{842517724, "Diver's radio, Iraq (Republic of)"},
		{992351000, "Aids to navigation, United Kingdom of Great Britain and Northern Ireland"},
		{1000010000, "Invalid MMSI"},
		{972345000, "MOB —Man Overboard Device"},
		{970241023, "AIS SART —Search and Rescue Transmitter, Greece"},
		{971356034, "Invalid MMSI"},
	}
	for _, c := range cases {
		got := DecodeMMSI(c.MMSI)
		if got != c.reference {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.reference)
			t.Errorf("DecodeMMSI(m uint32)")
		}
	}
}

func BenchmarkDecodeMMSI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		switch {
		case i%116 < 108:
			DecodeMMSI(316013198)
		case i%116 < 115:
			DecodeMMSI(002241076)
		default:
			DecodeMMSI(992351000)
		}
	}
}
