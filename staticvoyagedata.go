package ais

import (
	"errors"
	"fmt"
	//	"math"
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

	/*

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
		m.Speed = float32(bitsToInt(50, 59, data))
		if m.Speed < 1022 {
			m.Speed = m.Speed / 10
		}

		m.Accuracy = false
		if decodeAisChar(data[10])>>5 == 1 {
			m.Accuracy = true
		}

		//m.Lon = float64((int32(decodeAisChar(data[10]))<<27 | int32(decodeAisChar(data[11]))<<21 |
		//	int32(decodeAisChar(data[12]))<<15 | int32(decodeAisChar(data[13]))<<9 |
		//	int32(decodeAisChar(data[14]))>>1<<4)) / 16
		m.Lon = float64((int32(bitsToInt(61, 88, data)) << 4)) / 16
		//m.Lat = float64((int32(decodeAisChar(data[14]))<<31 | int32(decodeAisChar(data[15]))<<25 |
		//	int32(decodeAisChar(data[16]))<<19 | int32(decodeAisChar(data[17]))<<13 |
		//	int32(decodeAisChar(data[18]))<<7 | int32(decodeAisChar(data[19]))>>4<<5)) / 32
		m.Lat = float64((int32(bitsToInt(89, 115, data)) << 5)) / 32
		m.Lon, m.Lat = CoordinatesMin2Deg(m.Lon, m.Lat)

		//m.Course = float32(uint16(decodeAisChar(data[19]))<<12>>4|uint16(decodeAisChar(data[20]))<<2|
		//	uint16(decodeAisChar(data[21]))>>4) / 10
		m.Course = float32(bitsToInt(116, 127, data)) / 10

		//m.Heading = uint16(decodeAisChar(data[21]))<<12>>7 | uint16(decodeAisChar(data[22]))>>1
		m.Heading = uint16(bitsToInt(128, 136, data))

		//m.Second = decodeAisChar(data[22])<<7>>2 | decodeAisChar(data[23])>>1
		m.Second = uint8(bitsToInt(137, 142, data))

		m.Maneuver = decodeAisChar(data[23])<<7>>6 | decodeAisChar(data[24])>>5

		m.RAIM = false
		if decodeAisChar(data[24])<<6>>7 == 1 {
			m.RAIM = true
		}

		m.Radio = bitsToInt(149, 167, data)
	*/
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

	message :=
		fmt.Sprintf("=== Static and Voyage Related Data ===\n") +
			fmt.Sprintf(" Repeat       : %d\n", m.Repeat) +
			fmt.Sprintf(" MMSI         : %09d [%s]\n", m.MMSI, DecodeMMSI(m.MMSI)) +
			fmt.Sprintf(" AIS Version  : %d\n", m.AisVersion) +
			fmt.Sprintf(" IMO number   : %s\n", imo) +
			fmt.Sprintf(" Call Sign    : %s\n", m.Callsign) +
			fmt.Sprintf(" Vessel Name  : %s\n", m.VesselName) +
			fmt.Sprintf(" Ship Type    : %s\n", ShipType[int(m.ShipType)])

	return message
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
