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

// A StaticDataReport is a decoded AIS static data report (message type 24)
type StaticDataReport struct {
	Repeat uint8
	MMSI   uint32
	PartNo uint8
	//PartA
	VesselName string
	//PartB
	ShipType      uint8
	VendorID      string
	UnitModelCode uint8
	SerialNumber  uint32
	CallSign      string
	// optional with MothershipMMSI
	ToBow          uint16 // Dimension to bow
	ToStern        uint16 // Dimension to stern
	ToPort         uint8  // Dimension to port
	ToStarboard    uint8  // Dimension to starboard
	MothershipMMSI uint32
}

// DecodeStaticDataReport decodes the payload of a Type 24 AIS message
func DecodeStaticDataReport(payload string) (StaticDataReport, error) {
	data := []byte(payload)
	var m StaticDataReport

	mType := decodeAisChar(data[0])
	if mType != 24 {
		return m, errors.New("Message isn't Static Data Station Report (type 24).")
	}

	m.Repeat = uint8(bitsToInt(6, 7, data))

	// m.MMSI = uint32(decodeAisChar(data[1]))<<28>>2 | uint32(decodeAisChar(data[2]))<<20 |
	//	uint32(decodeAisChar(data[3]))<<14 | uint32(decodeAisChar(data[4]))<<8 |
	//	uint32(decodeAisChar(data[5]))<<2 | uint32(decodeAisChar(data[6]))>>4
	m.MMSI = bitsToInt(8, 37, data)

	m.PartNo = uint8(bitsToInt(38, 39, data))
	if m.PartNo == 0 {
		m.VesselName = bitsToString(112, 231, data)
	} else {
		m.ShipType = uint8(bitsToInt(40, 47, data))
		m.VendorID = bitsToString(48, 65, data)

		m.UnitModelCode = uint8(bitsToInt(66, 69, data))
		m.SerialNumber = uint32(bitsToInt(70, 89, data))
		m.CallSign = bitsToString(90, 131, data)

		// its an auxiliary craft
		if m.MMSI >= 980000000 {
			m.MothershipMMSI = bitsToInt(132, 161, data)
		} else {
			m.ToBow = uint16(bitsToInt(132, 140, data))
			m.ToStern = uint16(bitsToInt(141, 149, data))
			m.ToPort = uint8(bitsToInt(150, 155, data))
			m.ToStarboard = uint8(bitsToInt(156, 161, data))
		}
	}

	return m, nil
}
