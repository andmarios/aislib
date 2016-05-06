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
	"fmt"
	"time"
)

// A BaseStationReport is a decoded AIS base station report (message type 4)
type BaseStationReport struct {
	Repeat   uint8
	MMSI     uint32
	Time     time.Time
	Accuracy bool
	Lon      float64
	Lat      float64
	EPFD     uint8 // Enum type
	RAIM     bool
	Radio    uint32
}

// EPFD Fix Codes
var EpfdFixTypes = [...]string{
	"Undefined", "GPS", "GLONASS", "Combined GPS/GLONASS", "Loran-C",
	"Chayka", "Integrated Navigation System", "Surveyed", "Galileo",
	"not defined", "not defined", "not defined", "not defined",
	"not defined", "not defined", "not defined",
}

// DecodeBaseStationReport decodes the payload of a Type 4 AIS message
func DecodeBaseStationReport(payload string) (BaseStationReport, error) {
	data := []byte(payload)
	var m BaseStationReport

	mType := decodeAisChar(data[0])
	if mType != 4 {
		return m, errors.New("Message isn't Base Station Report (type 4).")
	}

	//m.Repeat = decodeAisChar(data[1]) >> 4
	m.Repeat = uint8(bitsToInt(6, 7, data))

	//m.MMSI = uint32(decodeAisChar(data[1]))<<28>>2 | uint32(decodeAisChar(data[2]))<<20 |
	//	uint32(decodeAisChar(data[3]))<<14 | uint32(decodeAisChar(data[4]))<<8 |
	//	uint32(decodeAisChar(data[5]))<<2 | uint32(decodeAisChar(data[6]))>>4
	m.MMSI = bitsToInt(8, 37, data)

	m.Time, _ = GetReferenceTime(payload) // Some base stations do not report time, for this case we do not consider it as error

	m.Accuracy = cbnBool(78, data)

	m.Lon, m.Lat = cbnCoordinates(79, data)

	m.EPFD = uint8(bitsToInt(134, 137, data))

	m.RAIM = cbnBool(148, data)

	m.Radio = bitsToInt(149, 167, data)
	return m, nil
}

// GetReferenceTime takes [the payload of] an AIS Base Station message (type 4)
// and returns the time data of it. It is a separate function from DecodeBaseStationReport
// because it can be useful to set a timeframe for our received AIS messages.
func GetReferenceTime(payload string) (time.Time, error) {
	data := []byte(payload)

	//year := uint16(decodeAisChar(data[6]))<<12>>2 | uint16(decodeAisChar(data[7]))<<4 |
	//	uint16(decodeAisChar(data[8]))>>2
	year := bitsToInt(38, 51, data)
	if year == 0 {
		var t time.Time
		return t, errors.New("station doesn't report time")
	}

	//month := decodeAisChar(data[8])<<6>>4 | decodeAisChar(data[9])>>4
	//day := decodeAisChar(data[9])<<4>>3 | decodeAisChar(data[10])>>5
	//hour := decodeAisChar(data[10]) << 3 >> 3
	//minute := decodeAisChar(data[11])
	//second := decodeAisChar(data[12])
	month := bitsToInt(52, 55, data)
	day := bitsToInt(56, 60, data)
	hour := bitsToInt(61, 65, data)
	minute := bitsToInt(66, 71, data)
	second := bitsToInt(72, 77, data)

	timeString := fmt.Sprintf("%d/%d/%d %d:%d:%d", year, month, day, hour, minute, second)
	t, _ := time.Parse("2006/1/2 15:4:5", timeString)

	return t, nil
}
