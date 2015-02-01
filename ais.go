package ais

import (
	"bufio"
	"log"
	"os"
	"strings"
	"encoding/hex"
	"math"
	"fmt"
	"strconv"
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
	Manuveur uint8   // maneuver indicator (enumerated)
	RAIM     bool    // RAIM flag
	Radio    uint32  // Radio status
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanLines)

	for in.Scan() {

		line := in.Text()

		if Nmea183ChecksumCheck(line) {

			tokens := strings.Split(line, ",")
			if tokens[0] == "!AIVDM" && // Sentence is AIS data
				tokens[1] == "1" &&     // Payload doesn't span across two sentences (ok for messages 1/2/3)
				tokens[6][:1] == "0" {  // Message doesn't need weird padding (ok for messages 1/2/3)

				//log.Println("Line length:", len(line), "Tokens:", len(tokens), "Payload:", tokens[5], "Checksum:", nmea183ChecksumCheck(line))
				PrintAisData(DecodeAisPosition(tokens[5]))
			} else {
				log.Println("There was an error with message:", line)
			}

		} else {
			log.Println("Checksum failed:", line)
		}
	}

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


func DecodeAisPosition(payload string) AisPositionMessage {
	data := []byte(payload)

	var m AisPositionMessage
	m.Type = decodeAisChar(data[0])

	if m.Type != 1 && m.Type != 2 && m.Type != 3 {
		log.Println("Message isn't Position Report.")
		return m
	}

	m.Repeat = decodeAisChar(data[1]) >> 4

	m.MMSI = uint32(decodeAisChar(data[1])) << 28
	m.MMSI = m.MMSI >> 2
	m.MMSI += uint32(decodeAisChar(data[2])) << 20 + uint32(decodeAisChar(data[3])) << 14 + uint32(decodeAisChar(data[4])) << 8 + uint32(decodeAisChar(data[5])) << 2
	m.MMSI += uint32(decodeAisChar(data[6])) >> 4

	m.Status = ( decodeAisChar(data[6]) << 4 ) >> 4

	m.Turn = float32(int8(decodeAisChar(data[7]) << 2 + decodeAisChar(data[8]) >> 4))
	if m.Turn != 0 && m.Turn <= 126 && m.Turn >= -126 {
		sign := float32(1)
		if math.Signbit(float64(m.Turn)) {
			sign = -1
		}
		m.Turn = sign * (m.Turn/4.733) * (m.Turn/4.733)

	}

	m.Speed = float32(uint16(decodeAisChar(data[8])) << 12 >> 6 + uint16(decodeAisChar(data[9])))
	if m.Speed < 1022 {
		m.Speed = m.Speed / 10
	}

	accuracy := decodeAisChar(data[10]) >> 5
	if accuracy == 0 {
		m.Accuracy = false
	} else {
		m.Accuracy = true
	}

	m.Lon = float64((int32(decodeAisChar(data[10])) << 27 + int32(decodeAisChar(data[11])) << 21 +
		int32(decodeAisChar(data[12])) << 15 + int32(decodeAisChar(data[13])) << 9 + int32(decodeAisChar(data[14])) >> 1 << 4)) / 16
	m.Lat = float64((int32(decodeAisChar(data[14])) << 31 + int32(decodeAisChar(data[15])) << 25 +
		int32(decodeAisChar(data[16])) << 19 + int32(decodeAisChar(data[17])) << 13 + int32(decodeAisChar(data[18])) << 7 + int32(decodeAisChar(data[19])) >> 4 << 5 )) / 32
	m.Lon, m.Lat = CoordinatesMin2Deg(m.Lon, m.Lat)
	//log.Println(CoordinatesDeg2Human(m.Lon, m.Lat))

	return m

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

func PrintAisData(m AisPositionMessage) {

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

	fmt.Printf("=== Message Type %d ===\n", m.Type)
	fmt.Printf(" Repeat     : %d\n", m.Repeat)
	fmt.Printf(" MMSI       : %d\n", m.MMSI)
	fmt.Printf(" Nav.Status : %s\n", status[m.Status])
	fmt.Printf(" Turn       : %s\n", turn)
	fmt.Printf(" Speed      : %s\n", speed)
	fmt.Printf(" Accuracy   : %s\n", accuracy)
	fmt.Printf(" Coordinates: %s\n", CoordinatesDeg2Human(m.Lon, m.Lat))
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







