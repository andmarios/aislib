package ais

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// StaticVoyageData is a type 5 AIS message (static and voyage related data)
// ETA is not reliable
type StaticVoyageData struct {
	Repeat      uint8
	MMSI        uint32
	AisVersion  uint8
	IMO         uint32 // IMO Ship ID number
	Callsign    string
	VesselName  string
	ShipType    uint8
	ToBow       uint16    // Dimension to bow
	ToStern     uint16    // Dimension to stern
	ToPort      uint8     // Dimension to port
	ToStarboard uint8     // Dimension to starboard
	EPFD        uint8     // Position Fix Type (enumeration declared at basestationreport.go)
	ETA         time.Time // Not reliable
	Draught     uint8     // Meters/10
	Destination string
	DTE         bool
}

// DecodeStaticVoyageData decodes [the payload of] an AIS Static and Voyage Related Data message (type 5)
func DecodeStaticVoyageData(payload string) (StaticVoyageData, error) {
	data := []byte(payload)
	var m StaticVoyageData

	mType := decodeAisChar(data[0])
	if mType != 5 {
		return m, errors.New("Message isn't Static and Voyage Related Data (type 5).")
	}
	m.Repeat = uint8(bitsToInt(6, 7, data))

	m.MMSI = uint32(bitsToInt(8, 37, data))

	m.AisVersion = uint8(bitsToInt(38, 39, data))

	m.IMO = uint32(bitsToInt(40, 69, data))

	m.Callsign = bitsToString(70, 111, data)

	m.VesselName = bitsToString(112, 231, data)

	m.ShipType = uint8(bitsToInt(232, 239, data))

	m.ToBow = uint16(bitsToInt(240, 248, data))
	m.ToStern = uint16(bitsToInt(249, 257, data))
	m.ToPort = uint8(bitsToInt(258, 263, data))
	m.ToStarboard = uint8(bitsToInt(264, 269, data))

	m.EPFD = uint8(bitsToInt(270, 273, data))

	// ETA does not include year, so we omit it too (it is set as 0000)
	// cyear := time.Now().Year()
	month := uint8(bitsToInt(274, 277, data))
	day := uint8(bitsToInt(278, 282, data))
	hour := uint8(bitsToInt(283, 287, data))
	minute := uint8(bitsToInt(288, 293, data))
	timeString := fmt.Sprintf("%d/%d %d:%d", month, day, hour, minute)
	m.ETA, _ = time.Parse("1/2 15:4", timeString)

	m.Draught = uint8(bitsToInt(294, 301, data))

	m.Destination = bitsToString(302, 421, data)

	m.DTE = cbnBool(422, data)

	return m, nil
}

// PrintPositionData returns a formatted string with the detailed data of a AIS position message.
// Its main use is to act as a guide for any developer wishing to correctly parse an AIS position message,
// since some parts of a message are enumareted, and other parts although they mainly are numeric values,
// for certain values they can have a non-numeric meaning.
func PrintStaticVoyageData(m StaticVoyageData) string {

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

// Ship types codes.
var ShipType = map[int]string{
	0:  "Not available",
	1:  "Reserved for future use",
	2:  "Reserved for future use",
	3:  "Reserved for future use",
	4:  "Reserved for future use",
	5:  "Reserved for future use",
	6:  "Reserved for future use",
	7:  "Reserved for future use",
	8:  "Reserved for future use",
	9:  "Reserved for future use",
	10: "Reserved for future use",
	11: "Reserved for future use",
	12: "Reserved for future use",
	13: "Reserved for future use",
	14: "Reserved for future use",
	15: "Reserved for future use",
	16: "Reserved for future use",
	17: "Reserved for future use",
	18: "Reserved for future use",
	19: "Reserved for future use",
	20: "Wing in ground (WIG)",
	21: "Wing in ground (WIG), Hazardous category A",
	22: "Wing in ground (WIG), Hazardous category B",
	23: "Wing in ground (WIG), Hazardous category C",
	24: "Wing in ground (WIG), Hazardous category D",
	25: "Wing in ground (WIG), Reserved for future use",
	26: "Wing in ground (WIG), Reserved for future use",
	27: "Wing in ground (WIG), Reserved for future use",
	28: "Wing in ground (WIG), Reserved for future use",
	29: "Wing in ground (WIG), Reserved for future use",
	30: "Fishing",
	31: "Towing",
	32: "Towing: length exceeds 200m or breadth exceeds 25m",
	33: "Dredging or underwater ops",
	34: "Diving ops",
	35: "Military ops",
	36: "Sailing",
	37: "Pleasure Craft",
	38: "Reserved",
	39: "Reserved",
	40: "High speed craft (HSC)",
	41: "High speed craft (HSC), Hazardous category A",
	42: "High speed craft (HSC), Hazardous category B",
	43: "High speed craft (HSC), Hazardous category C",
	44: "High speed craft (HSC), Hazardous category D",
	45: "High speed craft (HSC), Reserved for future use",
	46: "High speed craft (HSC), Reserved for future use",
	47: "High speed craft (HSC), Reserved for future use",
	48: "High speed craft (HSC), Reserved for future use",
	49: "High speed craft (HSC), No additional information",
	50: "Pilot Vessel",
	51: "Search and Rescue vessel",
	52: "Tug",
	53: "Port Tender",
	54: "Anti-pollution equipment",
	55: "Law Enforcement",
	56: "Spare - Local Vessel",
	57: "Spare - Local Vessel",
	58: "Medical Transport",
	59: "Noncombatant ship according to RR Resolution No. 18",
	60: "Passenger",
	61: "Passenger, Hazardous category A",
	62: "Passenger, Hazardous category B",
	63: "Passenger, Hazardous category C",
	64: "Passenger, Hazardous category D",
	65: "Passenger, Reserved for future use",
	66: "Passenger, Reserved for future use",
	67: "Passenger, Reserved for future use",
	68: "Passenger, Reserved for future use",
	69: "Passenger, No additional information",
	70: "Cargo",
	71: "Cargo, Hazardous category A",
	72: "Cargo, Hazardous category B",
	73: "Cargo, Hazardous category C",
	74: "Cargo, Hazardous category D",
	75: "Cargo, Reserved for future use",
	76: "Cargo, Reserved for future use",
	77: "Cargo, Reserved for future use",
	78: "Cargo, Reserved for future use",
	79: "Cargo, No additional information",
	80: "Tanker",
	81: "Tanker, Hazardous category A",
	82: "Tanker, Hazardous category B",
	83: "Tanker, Hazardous category C",
	84: "Tanker, Hazardous category D",
	85: "Tanker, Reserved for future use",
	86: "Tanker, Reserved for future use",
	87: "Tanker, Reserved for future use",
	88: "Tanker, Reserved for future use",
	89: "Tanker, No additional information",
	90: "Other Type",
	91: "Other Type, Hazardous category A",
	92: "Other Type, Hazardous category B",
	93: "Other Type, Hazardous category C",
	94: "Other Type, Hazardous category D",
	95: "Other Type, Reserved for future use",
	96: "Other Type, Reserved for future use",
	97: "Other Type, Reserved for future use",
	98: "Other Type, Reserved for future use",
	99: "Other Type, no additional information",
}
