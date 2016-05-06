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
	"strconv"
)

// PrintBaseStationReport returns a formatted string of a BaseStationReport. Mainly to help
// developers with understanding base position reports.
func (m BaseStationReport) String() string {
	accuracy := "High accuracy (<10m)"
	if m.Accuracy == false {
		accuracy = "Low accuracy (>10m)"
	}

	raim := "not in use"
	if m.RAIM == true {
		raim = "in use"
	}

	message :=
		fmt.Sprintf("=== Base Station Report ===\n") +
			fmt.Sprintf(" Repeat       : %d\n", m.Repeat) +
			fmt.Sprintf(" MMSI         : %09d [%s]\n", m.MMSI, DecodeMMSI(m.MMSI)) +
			fmt.Sprintf(" Time         : %s\n", m.Time.String()) +
			fmt.Sprintf(" Accuracy     : %s\n", accuracy) +
			fmt.Sprintf(" Coordinates  : %s\n", CoordinatesDeg2Human(m.Lon, m.Lat)) +
			fmt.Sprintf(" EPFD         : %s\n", EpfdFixTypes[m.EPFD]) +
			fmt.Sprintf(" RAIM         : %s\n", raim)

	return message
}

// PrintClassAPositionReport returns a formatted string with the detailed data of a AIS position message.
// Its main use is to act as a guide for any developer wishing to correctly parse an AIS position message,
// since some parts of a message are enumareted, and other parts although they mainly are numeric values,
// for certain values they can have a non-numeric meaning.
func (m ClassAPositionReport) String() string {
	turn := ""
	switch {
	case m.Turn == 0:
		turn = "not turning"
	case m.Turn == 127:
		turn = "right at more than 5deg/30s"
	case m.Turn == -127:
		turn = "left at more than 5deg/30s"
	case m.Turn == -128:
		turn = "no turn information"
	case m.Turn > 0 && m.Turn < 127:
		turn = "right at " + strconv.FormatFloat(float64(m.Turn), 'f', 3, 32)
	case m.Turn < 0 && m.Turn > -127:
		turn = "left at " + strconv.FormatFloat(float64(-m.Turn), 'f', 3, 32)
	}

	speed := ""
	switch {
	case m.Speed <= 102:
		speed = strconv.FormatFloat(float64(m.Speed), 'f', 1, 32) + " knots"
	case m.Speed == 1022:
		speed = ">102.2 knots"
	case m.Speed == 1023:
		speed = "information not available"
	}

	accuracy := "High accuracy (<10m)"
	if m.Accuracy == false {
		accuracy = "Low accuracy (>10m)"
	}

	course := ""
	switch {
	case m.Course < 360:
		course = fmt.Sprintf("%.1f°", m.Course)
	case m.Course == 360:
		course = "not available"
	case m.Course > 360:
		course = "please report this to developer"
	}

	heading := ""
	switch {
	case m.Heading <= 359:
		heading = fmt.Sprintf("%d°", m.Heading)
	case m.Heading == 511:
		heading = "not available"
	case m.Heading != 511 && m.Heading >= 360:
		heading = "please report this to developer"
	}

	maneuver := ""
	switch {
	case m.Maneuver == 0:
		maneuver = "not available"
	case m.Maneuver == 1:
		maneuver = "no special maneuver"
	case m.Maneuver == 2:
		maneuver = "special maneuver"
	}

	raim := "not in use"
	if m.RAIM == true {
		raim = "in use"
	}

	message :=
		fmt.Sprintf("=== Class A Position Report (%d) ===\n", m.Type) +
			fmt.Sprintf(" Repeat       : %d\n", m.Repeat) +
			fmt.Sprintf(" MMSI         : %09d [%s]\n", m.MMSI, DecodeMMSI(m.MMSI)) +
			fmt.Sprintf(" Nav.Status   : %s\n", NavigationStatusCodes[m.Status]) +
			fmt.Sprintf(" Turn (ROT)   : %s\n", turn) +
			fmt.Sprintf(" Speed (SOG)  : %s\n", speed) +
			fmt.Sprintf(" Accuracy     : %s\n", accuracy) +
			fmt.Sprintf(" Coordinates  : %s\n", CoordinatesDeg2Human(m.Lon, m.Lat)) +
			fmt.Sprintf(" Course (COG) : %s\n", course) +
			fmt.Sprintf(" Heading (HDG): %s\n", heading) +
			fmt.Sprintf(" Manuever ind.: %s\n", maneuver) +
			fmt.Sprintf(" RAIM         : %s\n", raim)

	return message
}

// PrintClassBPositionReport returns a formatted string with the detailed data of a AIS position message.
// Its main use is to act as a guide for any developer wishing to correctly parse an AIS position message,
// since some parts of a message are enumareted, and other parts although they mainly are numeric values,
// for certain values they can have a non-numeric meaning.
func (m ClassBPositionReport) String() string {
	speed := ""
	switch {
	case m.Speed <= 102:
		speed = strconv.FormatFloat(float64(m.Speed), 'f', 1, 32) + " knots"
	case m.Speed == 1022:
		speed = ">102.2 knots"
	case m.Speed == 1023:
		speed = "information not available"
	}

	accuracy := "High accuracy (<10m)"
	if m.Accuracy == false {
		accuracy = "Low accuracy (>10m)"
	}

	course := ""
	switch {
	case m.Course < 360:
		course = fmt.Sprintf("%.1f°", m.Course)
	case m.Course == 360:
		course = "not available"
	case m.Course > 360:
		course = "please report this to developer"
	}

	heading := ""
	switch {
	case m.Heading <= 359:
		heading = fmt.Sprintf("%d°", m.Heading)
	case m.Heading == 511:
		heading = "not available"
	case m.Heading != 511 && m.Heading >= 360:
		heading = "please report this to developer"
	}

	message :=
		fmt.Sprintf("=== Class B Position Report ===\n") +
			fmt.Sprintf(" Repeat       : %d\n", m.Repeat) +
			fmt.Sprintf(" MMSI         : %09d [%s]\n", m.MMSI, DecodeMMSI(m.MMSI)) +
			fmt.Sprintf(" Speed (SOG)  : %s\n", speed) +
			fmt.Sprintf(" Accuracy     : %s\n", accuracy) +
			fmt.Sprintf(" Coordinates  : %s\n", CoordinatesDeg2Human(m.Lon, m.Lat)) +
			fmt.Sprintf(" Course (COG) : %s\n", course) +
			fmt.Sprintf(" Heading (HDG): %s\n", heading) +
			fmt.Sprintf(" CS Unit      : %t\n", m.CSUnit) +
			fmt.Sprintf(" Display      : %t\n", m.Display) +
			fmt.Sprintf(" DSC          : %t\n", m.DSC) +
			fmt.Sprintf(" Band         : %t\n", m.Band) +
			fmt.Sprintf(" Message 22   : %t\n", m.Msg22) +
			fmt.Sprintf(" Assigned     : %t\n", m.Assigned) +
			fmt.Sprintf(" RAIM         : %t\n", m.RAIM)

	return message
}

// PrintStaticVoyageData returns a formatted string with the detailed data of a AIS Static and Voyage
// Related Data (message type 5). Its main use is to act as a guide for any developer wishing to
// correctly parse an AIS type 5 message since some parts are enumareted, and other parts although
// they mainly are numeric values, for certain values they can have a non-numeric meaning.
func (m StaticVoyageData) String() string {

	imo := ""
	if m.IMO == 0 {
		imo = "Inland Vessel"
	} else {
		imo = strconv.Itoa(int(m.IMO))
	}

	draught := ""
	if m.Draught == 0 {
		draught = "Not available"
	} else {
		draught = strconv.Itoa(10*int(m.IMO)) + " meters"
	}

	message :=
		fmt.Sprintf("=== Static and Voyage Related Data ===\n") +
			fmt.Sprintf(" Repeat       : %d\n", m.Repeat) +
			fmt.Sprintf(" MMSI         : %09d [%s]\n", m.MMSI, DecodeMMSI(m.MMSI)) +
			fmt.Sprintf(" AIS Version  : %d\n", m.AisVersion) +
			fmt.Sprintf(" IMO number   : %s\n", imo) +
			fmt.Sprintf(" Call Sign    : %s\n", m.Callsign) +
			fmt.Sprintf(" Vessel Name  : %s\n", m.VesselName) +
			fmt.Sprintf(" Ship Type    : %s\n", ShipType[int(m.ShipType)]) +
			fmt.Sprintf(" Dim to Bow   : %s\n", type5size2String(0, 511, int(m.ToBow))) +
			fmt.Sprintf(" Dim to Stern : %s\n", type5size2String(0, 511, int(m.ToStern))) +
			fmt.Sprintf(" Dim to Port  : %s\n", type5size2String(0, 511, int(m.ToPort))) +
			fmt.Sprintf(" Dim to StrBrd: %s\n", type5size2String(0, 511, int(m.ToStarboard))) +
			fmt.Sprintf(" EPFD         : %s\n", EpfdFixTypes[m.EPFD]) +
			fmt.Sprintf(" ETA          : %s\n", m.ETA.String()) +
			fmt.Sprintf(" Draught      : %s\n", draught) +
			fmt.Sprintf(" Destination  : %s\n", m.Destination)

	return message
}

// A small function to translate the size fields
func type5size2String(min, max, size int) string {
	s := ""
	switch size {
	case min:
		s = "Not available"
	case max:
		s = ">" + strconv.Itoa(max) + " meters"
	default:
		s = strconv.Itoa(size) + " meters"
	}
	return s
}

// PrintBinaryBroadcast returns a string with some data for a Binary Broadcast message
func (m BinaryBroadcast) String() string {

	message :=
		fmt.Sprintf("=== Binary Broadcast ===\n") +
			fmt.Sprintf(" Repeat       : %d\n", m.Repeat) +
			fmt.Sprintf(" MMSI         : %09d [%s]\n", m.MMSI, DecodeMMSI(m.MMSI)) +
			fmt.Sprintf(" DAC-FID      : %d-%d (%s)\n", m.DAC, m.FID, BinaryBroadcastType[int(m.DAC)][int(m.FID)])

	return message
}
