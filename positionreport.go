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
	"math"
)

// A PositionReport is the generic structure of a Position Report, containing the common fields
// between Class A and B reports.
type PositionReport struct {
	Type     uint8
	Repeat   uint8
	MMSI     uint32
	Speed    float32 // speed over ground - SOG (sc U3)
	Accuracy bool    // position accuracy
	Lon      float64 // (sc I4)
	Lat      float64 // (sc I4)
	Course   float32 //course over ground - COG (sc U1)
	Heading  uint16  // true heading - HDG
	Second   uint8   // timestamp
	RAIM     bool    // RAIM flag
	Radio    uint32  // Radio status
}

// A ClassAPositionReport is a decoded AIS position message (messages of type 1, 2 or 3).
// Please have a look at http://catb.org/gpsd/AIVDM.html and at
// http://www.navcen.uscg.gov/?pageName=AISMessagesA
type ClassAPositionReport struct {
	PositionReport
	Status   uint8   // navigation status (enumerated type)
	Turn     float32 // rate of turn - ROT (sc - Special Calc I3)
	Maneuver uint8   // maneuver indicator (enumerated)
}

// A ClassBPositionReport is a decoded AIS position message (type 18).
type ClassBPositionReport struct {
	PositionReport
	CSUnit   bool
	Display  bool
	DSC      bool
	Band     bool
	Msg22    bool
	Assigned bool
}

// Navigation status codes
var NavigationStatusCodes = [...]string{
	"Under way using engine", "At anchor", "Not under command", "Restricted maneuverability",
	"Constrained by her draught", "Moored", "Aground", "Engaged in fishing", "Under way sailing",
	"status code reserved", "status code reserved", "status code reserved",
	"status code reserved", "status code reserved", "AIS-SART is active", "Not defined",
}

// DecodeClassAPositionReport decodes [the payload of] an AIS position message (type 1/2/3)
func DecodeClassAPositionReport(payload string) (ClassAPositionReport, error) {
	data := []byte(payload)
	var m ClassAPositionReport

	m.Type = decodeAisChar(data[0])
	if m.Type != 1 && m.Type != 2 && m.Type != 3 {
		return m, errors.New("Message isn't Class A Position Report (type 1, 2 or 3).")
	}

	// !!! This is the first decoding function written. Original decoding
	// routines are left here as comments, in order to help anyone diving
	// into binary field decoding.

	//m.Repeat = decodeAisChar(data[1]) >> 4
	m.Repeat = uint8(bitsToInt(6, 7, data))

	//m.MMSI = uint32(decodeAisChar(data[1]))<<28>>2 | uint32(decodeAisChar(data[2]))<<20 |
	//	uint32(decodeAisChar(data[3]))<<14 | uint32(decodeAisChar(data[4]))<<8 |
	//	uint32(decodeAisChar(data[5]))<<2 | uint32(decodeAisChar(data[6]))>>4
	m.MMSI = bitsToInt(8, 37, data)

	//m.Status = (decodeAisChar(data[6]) << 4) >> 4
	m.Status = uint8(bitsToInt(38, 41, data))

	//m.Turn = float32(int8(decodeAisChar(data[7])<<2 | decodeAisChar(data[8])>>4))
	m.Turn = float32(int8(bitsToInt(42, 49, data)))
	if m.Turn != 0 && m.Turn <= 126 && m.Turn >= -126 {
		sign := float32(1)
		if math.Signbit(float64(m.Turn)) {
			sign = -1
		}
		m.Turn = sign * (m.Turn / 4.733) * (m.Turn / 4.733)

	}

	//m.Speed = float32(uint16(decodeAisChar(data[8]))<<12>>6 | uint16(decodeAisChar(data[9])))
	m.Speed = cbnSpeed(50, data)

	//m.Accuracy = false
	//if decodeAisChar(data[10])>>5 == 1 {
	//	m.Accuracy = true
	//}
	m.Accuracy = cbnBool(60, data)

	// Old method 1
	//m.Lon = float64((int32(decodeAisChar(data[10]))<<27 | int32(decodeAisChar(data[11]))<<21 |
	//	int32(decodeAisChar(data[12]))<<15 | int32(decodeAisChar(data[13]))<<9 |
	//	int32(decodeAisChar(data[14]))>>1<<4)) / 16
	//m.Lat = float64((int32(decodeAisChar(data[14]))<<31 | int32(decodeAisChar(data[15]))<<25 |
	//	int32(decodeAisChar(data[16]))<<19 | int32(decodeAisChar(data[17]))<<13 |
	//	int32(decodeAisChar(data[18]))<<7 | int32(decodeAisChar(data[19]))>>4<<5)) / 32
	// Old method 2
	//m.Lon = float64((int32(bitsToInt(61, 88, data)) << 4)) / 16
	//m.Lat = float64((int32(bitsToInt(89, 115, data)) << 5)) / 32
	// Finish or both old methods
	//m.Lon, m.Lat = CoordinatesMin2Deg(m.Lon, m.Lat)
	m.Lon, m.Lat = cbnCoordinates(61, data)

	//m.Course = float32(uint16(decodeAisChar(data[19]))<<12>>4|uint16(decodeAisChar(data[20]))<<2|
	//	uint16(decodeAisChar(data[21]))>>4) / 10
	m.Course = float32(bitsToInt(116, 127, data)) / 10

	//m.Heading = uint16(decodeAisChar(data[21]))<<12>>7 | uint16(decodeAisChar(data[22]))>>1
	m.Heading = uint16(bitsToInt(128, 136, data))

	//m.Second = decodeAisChar(data[22])<<7>>2 | decodeAisChar(data[23])>>1
	m.Second = uint8(bitsToInt(137, 142, data))

	//m.Maneuver = decodeAisChar(data[23])<<7>>6 | decodeAisChar(data[24])>>5
	m.Maneuver = uint8(bitsToInt(143, 144, data))

	//m.RAIM = false
	//if decodeAisChar(data[24])<<6>>7 == 1 {
	//	m.RAIM = true
	//}
	m.RAIM = cbnBool(148, data)

	m.Radio = bitsToInt(149, 167, data)
	return m, nil
}

// DecodeClassBPositionReport decodes [the payload of] an AIS position message (type 18)
func DecodeClassBPositionReport(payload string) (ClassBPositionReport, error) {
	data := []byte(payload)
	var m ClassBPositionReport

	m.Type = decodeAisChar(data[0])
	if m.Type != 18 {
		return m, errors.New("Message isn't Class B Position Report (type 18).")
	}

	m.Repeat = uint8(bitsToInt(6, 7, data))

	m.MMSI = bitsToInt(8, 37, data)

	m.Speed = cbnSpeed(46, data)

	m.Accuracy = cbnBool(56, data)

	m.Lon, m.Lat = cbnCoordinates(57, data)

	m.Course = float32(bitsToInt(112, 123, data)) / 10

	m.Heading = uint16(bitsToInt(124, 132, data))

	m.Second = uint8(bitsToInt(133, 138, data))

	m.CSUnit = cbnBool(141, data)
	m.Display = cbnBool(142, data)
	m.DSC = cbnBool(143, data)
	m.Band = cbnBool(144, data)
	m.Msg22 = cbnBool(145, data)
	m.Assigned = cbnBool(146, data)

	m.RAIM = cbnBool(147, data)
	if decodeAisChar(data[24])<<6>>7 == 1 {
		m.RAIM = true
	}

	m.Radio = bitsToInt(148, 167, data)
	return m, nil
}
