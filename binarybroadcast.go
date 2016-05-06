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
	"errors"
)

// BinaryBroadcast is a Type 8 message
type BinaryBroadcast struct {
	Repeat uint8
	MMSI   uint32
	DAC    uint16
	FID    uint8
	Data   string
}

// DecodeBinaryBroadcast decodes [the payload of] an AIS Binary Broadcast message (Type 8) but not its binary payload
func DecodeBinaryBroadcast(payload string) (BinaryBroadcast, error) {
	data := []byte(payload)
	var m BinaryBroadcast

	mType := decodeAisChar(data[0])
	if mType != 8 {
		return m, errors.New("Message isn't Binary Broadcast (type 8).")
	}

	m.Repeat = uint8(bitsToInt(6, 7, data))

	m.MMSI = uint32(bitsToInt(8, 37, data))

	m.DAC = uint16(bitsToInt(40, 49, data))
	m.FID = uint8(bitsToInt(50, 55, data))

	m.Data = payload // Data start at bit 56, but this way we simplify our code

	return m, nil
}

// Some Binary Broadcast types. The list isn't complete but I haven't searched for a better source
var BinaryBroadcastType = map[int]map[int]string{
	1: {
		11: "Meteorological/Hydrogological Data",
		13: "Fairway closed",
		15: "Extended ship and voyage",
		17: "VTS-Generated/Synthetic targets",
		19: "Marine traffic signals",
		21: "Weather observation from ship",
		22: "Area notice (broadcast)",
		24: "Extended ship and voyage",
		26: "Environmental",
		27: "Route info broadcast",
		29: "Text description broadcast",
		31: "Meteorological and Hydrological",
	},
	200: {
		10: "Ship static and voyage related data",
		23: "EMMA warning report",
		24: "Water levels",
		40: "Signal status",
	},
	316: {
		1:  "Weather Station or Wind or Water Level",
		2:  "Lockage Order or Estimated Lock Times",
		32: "Seaway Version Message",
	},
	366: {
		1:  "Weather Station or Wind or Water Level or PAWS Hydro / Current or PAWS Hydro / Salinity Temp or PAWS Vessel Procession Order",
		2:  "Lockage Order or Estimated Lock Times",
		32: "Seaway Version Message",
	},
}
