package ais

import (
	"errors"
	"fmt"
	"math"
	"strconv"
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

	m.Maneuver = decodeAisChar(data[23])<<7>>6 | decodeAisChar(data[24])>>5

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

// PrintPositionData returns a formatted string with the detailed data of a AIS position message.
// Its main use is to act as a guide for any developer wishing to correctly parse an AIS position message,
// since some parts of a message are enumareted, and other parts although they mainly are numeric values,
// for certain values they can have a non-numeric meaning.
func PrintPositionData(m ClassAPositionReport) string {
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
