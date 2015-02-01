package ais

import (
	"encoding/hex"
	"math"
	"fmt"
	"strconv"
	"errors"
)

// Please have a look at <http://catb.org/gpsd/AIVDM.html> and at <http://www.navcen.uscg.gov/?pageName=AISMessagesA>
// This is the struct of AIS messages of types 1/2/3.
type AisPositionMessage struct {
	Type     uint8
	Repeat   uint8
	MMSI     uint32
	Status   uint8   // navigation status (enumerated type)
	Turn     float32 // rate of turn - ROT (sc - Special Calc I3)
	Speed    float32 // speed over ground - SOG (sc U3)
	Accuracy bool    // position accuracy
	Lon      float64 // (sc I4)
	Lat      float64 // (sc I4)
	Course   float32 //course over ground - COG (sc U1)
	Heading  uint16  // true heading - HDG
	Second   uint8   // timestamp
	Maneuver uint8   // maneuver indicator (enumerated)
	RAIM     bool    // RAIM flag
	Radio    uint32  // Radio status
}

func decodeAisChar(character byte) byte {
	character -= 48
	if character > 40 {
		character -= 8
	}
	return character
}

func AisMessageType(payload string) uint8 {
	data := []byte(payload)
	return decodeAisChar(data[0])
}

func DecodeAisPosition(payload string) (AisPositionMessage, error) {
	data := []byte(payload)

	var m AisPositionMessage

	m.Type = decodeAisChar(data[0])

	if m.Type != 1 && m.Type != 2 && m.Type != 3 {
		return m, errors.New("Message isn't Position Report.")
	}

	m.Repeat = decodeAisChar(data[1]) >> 4

	m.MMSI = uint32(decodeAisChar(data[1])) << 28
	m.MMSI = m.MMSI >> 2
	m.MMSI += uint32(decodeAisChar(data[2])) << 20 | uint32(decodeAisChar(data[3])) << 14 | uint32(decodeAisChar(data[4])) << 8 | uint32(decodeAisChar(data[5])) << 2
	m.MMSI += uint32(decodeAisChar(data[6])) >> 4

	m.Status = ( decodeAisChar(data[6]) << 4 ) >> 4

	m.Turn = float32(int8(decodeAisChar(data[7]) << 2 | decodeAisChar(data[8]) >> 4))
	if m.Turn != 0 && m.Turn <= 126 && m.Turn >= -126 {
		sign := float32(1)
		if math.Signbit(float64(m.Turn)) {
			sign = -1
		}
		m.Turn = sign * (m.Turn/4.733) * (m.Turn/4.733)

	}

	m.Speed = float32(uint16(decodeAisChar(data[8])) << 12 >> 6 | uint16(decodeAisChar(data[9])))
	if m.Speed < 1022 {
		m.Speed = m.Speed / 10
	}

	m.Accuracy = false
	if decodeAisChar(data[10]) >> 5 == 1 {
		m.Accuracy = true
	}

	m.Lon = float64((int32(decodeAisChar(data[10])) << 27 | int32(decodeAisChar(data[11])) << 21 |
		int32(decodeAisChar(data[12])) << 15 | int32(decodeAisChar(data[13])) << 9 | int32(decodeAisChar(data[14])) >> 1 << 4)) / 16
	m.Lat = float64((int32(decodeAisChar(data[14])) << 31 | int32(decodeAisChar(data[15])) << 25 |
		int32(decodeAisChar(data[16])) << 19 | int32(decodeAisChar(data[17])) << 13 | int32(decodeAisChar(data[18])) << 7 | int32(decodeAisChar(data[19])) >> 4 << 5 )) / 32
	m.Lon, m.Lat = CoordinatesMin2Deg(m.Lon, m.Lat)

	m.Course = float32(uint16(decodeAisChar(data[19])) << 12 >> 4 | uint16(decodeAisChar(data[20])) << 2 | uint16(decodeAisChar(data[21])) >> 4) / 10

	m.Heading = uint16(decodeAisChar(data[21])) << 12 >> 7 | uint16(decodeAisChar(data[22])) >> 1

	m.Second = decodeAisChar(data[22]) << 7 >> 2 | decodeAisChar(data[23]) >> 1

	m.Maneuver = decodeAisChar(data[23]) << 7 >> 6 | decodeAisChar(data[24]) >> 5

	m.RAIM = false
	if decodeAisChar(data[24]) << 6 >> 7 == 1 {
		m.RAIM = true
	}

	return m, nil
}

func CoordinatesMin2Deg(minLon, minLat float64) (float64, float64) {
	lonSign := 1.0
	latSign := 1.0

	if math.Signbit(minLon) {
		minLon = -minLon
		lonSign = -1
	}
	if math.Signbit(minLat) {
		minLat = -minLat
		latSign = -1
	}

	degrees := float64(int(minLon/600000))
	minutes := float64(minLon - 600000*degrees)/10000
	lon := degrees + minutes/60

	degrees = float64(int(minLat/600000))
	minutes = float64(minLat - 600000*degrees)/10000
	lat := degrees + minutes/60

	return lonSign*lon, latSign*lat
}

func CoordinatesDeg2Human(degLon, degLat float64) string {
	lonSign := 1.0
	latSign := 1.0
	coordinates := ""

	if math.Signbit(degLon) {
		degLon = -degLon
		lonSign = -1
	}
	if math.Signbit(degLat) {
		degLat = -degLat
		latSign = -1
	}

	degrees := math.Floor(degLon)
	minutes := 60*(degLon - degrees)


	if degrees > 180 {
		coordinates = "longitude not available, "
	} else if lonSign > 0 {
		coordinates = fmt.Sprintf("%3.0f°%07.4f'%s", degrees, minutes, "E")
	} else {
		coordinates = fmt.Sprintf("%3.0f°%07.4f'%s", degrees, minutes, "W")
	}

	degrees = math.Floor(degLat)
	minutes = 60*(degLat - degrees)

	if degrees > 90 {
		coordinates += "latitude not available"
	}else if latSign > 0 {
		coordinates += fmt.Sprintf(" %3.0f°%07.4f%s", degrees, minutes, "N")
	} else {
		coordinates += fmt.Sprintf(" %3.0f°%07.4f%s", degrees, minutes, "S")
	}

	return coordinates
}

func PrintAisPositionData(m AisPositionMessage) string {

	status := []string{"Under way using engine", "At anchor", "Not under command", "Restricted maneuverability", "Constrained by her draught",
		"Moored", "Aground", "Engaged in fishing", "Under way sailing", "status code reserved", "status code reserved", "status code reserved",
		"status code reserved", "status code reserved", "AIS-SART is active", "Not defined"}

	turn := ""
	switch {
	case m.Turn == 0:
		turn = "not turning";
	case m.Turn == 127:
		turn = "right at more than 5deg/30s"
	case m.Turn == -127:
		turn = "left at more than 5deg/30s"
	case m.Turn == -128:
		turn = "no turn information"
	case m.Turn > 0 && m.Turn < 127:
		turn = "right at " + strconv.FormatFloat(float64(m.Turn), 'f', 3, 32)
	case m.Turn < 0 && m.Turn > -127:
		turn = "left at " + strconv.FormatFloat(float64(m.Turn), 'f', 3, 32)
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

	accuracy := ""
	switch {
	case m.Accuracy == true:
		accuracy = "High accuracy (<10m)"
	case m.Accuracy == false:
		accuracy = "Low accuracy (>10m)"
	}

	course := ""
	switch {
	case m.Course < 360:
		course = fmt.Sprintf("%.1f°", m.Course)
	case m.Course == 360:
		course = fmt.Sprintf("not available")
	case m.Course > 360:
		course = fmt.Sprintf("please report this to developer")
	}

	message :=
		fmt.Sprintf("=== Message Type %d ===\n", m.Type) +
		fmt.Sprintf(" Repeat      : %d\n", m.Repeat) +
		fmt.Sprintf(" MMSI        : %d\n", m.MMSI) +
		fmt.Sprintf(" Nav.Status  : %s\n", status[m.Status]) +
		fmt.Sprintf(" Turn (ROT)  : %s\n", turn) +
		fmt.Sprintf(" Speed (SOG) : %s\n", speed) +
		fmt.Sprintf(" Accuracy    : %s\n", accuracy) +
		fmt.Sprintf(" Coordinates : %s\n", CoordinatesDeg2Human(m.Lon, m.Lat)) +
		fmt.Sprintf(" Course (COG): %s\n", course)

	return message
}

func Nmea183ChecksumCheck(sentence string) bool {
	length := len(sentence)

	var csum []byte
	csum, err := hex.DecodeString(sentence[length-2:])

	if err != nil {
		return false
	}

	bline := []byte(sentence[1:length-3])
	ccsum := bline[0]
	for i := 1; i < len(bline); i++ {
		ccsum ^= bline[i]
	}

	if csum[0] == ccsum {
		return true
	}
	return false
}







