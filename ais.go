// Package ais provides functions and types to work with AIS (Automatic Identification System) sentences (radio messages) and messages in AIVDM/AIVDO format.
//
// An AIS sentence is one line, it is the actual radio message.
// An AIS message is the payload that is carried by one or more consecutive AIS sentences.
package ais

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// A Message stores the important properties of a AIS message, including only information useful for decoding: Type, Payload, Padding Bits
// A Message should come after processing one or more AIS radio sentences (checksum check, concatenate payloads spanning across sentences, etc).
type Message struct {
	Type    uint8
	Payload string
	Padding uint8
}

// FailedSentence includes an AIS sentence that failed to process (e.g wrong checksum) and the reason it failed.
type FailedSentence struct {
	Sentence string
	Issue    string
}

// A PositionMessage is a decoded AIS position message (messages of type 1, 2 or 3).
// Please have a look at <http://catb.org/gpsd/AIVDM.html> and at <http://www.navcen.uscg.gov/?pageName=AISMessagesA>
type PositionMessage struct {
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

// MessageType returns the type of an AIS message
func MessageType(payload string) uint8 {
	data := []byte(payload[:1])
	return decodeAisChar(data[0])
}

// Router accepts AIS radio sentences and process them. It checks their checksum,
// and AIS identifiers. If they are valid it tries to assemble the payload if it spans
// on multiple sentences. Upon success it returns the AIS Message at the out channel.
// Failed sentences go to the err channel.
// If the in channel is closed, then it sends a message with type 255 at the out channel.
// Your function can check for this message to know when it is safe to exit the program.
func Router(in chan string, out chan Message, failed chan FailedSentence) {
	count, ccount, padding := 0, 0, 0
	size, id := "0", "0"
	payload := ""
	var cache [5]string
	var err error
	aisIdentifiers := map[string]bool{
		"ABVD": true, "ADVD": true, "AIVD": true, "ANVD": true, "ARVD": true,
		"ASVD": true, "ATVD": true, "AXVD": true, "BSVD": true, "SAVD": true,
	}
	for sentence := range in {
		tokens := strings.Split(sentence, ",") // I think this takes the major portion of time for this function (after benchmarking)

		if !Nmea183ChecksumCheck(sentence) { // Checksum check
			failed <- FailedSentence{sentence, "Checksum failed"}
			continue
		}

		if !aisIdentifiers[tokens[0][1:5]] { // Check for valid AIS identifier
			failed <- FailedSentence{sentence, "Sentence isn't AIVDM/AIVDO"}
			continue
		}

		if tokens[1] == "1" { // One sentence message, process it immediately
			padding, _ = strconv.Atoi(tokens[6][:1])
			out <- Message{MessageType(tokens[5]), tokens[5], uint8(padding)}
		} else { // Message spans across sentences.
			ccount, err = strconv.Atoi(tokens[2])
			if err != nil {
				failed <- FailedSentence{sentence, "HERE " + tokens[2]}
				continue
			}
			if ccount != count+1 || // If there are sentences with wrong sequence number in cache send them as failed
				tokens[3] != id && count != 0 || // If there are sentences with different sequence id in cache , send old parts as failed
				tokens[1] != size && count != 0 { // If there messages with wrogn size in cache, send them as failed
				for i := 0; i <= count; i++ {
					failed <- FailedSentence{cache[i], "Incomplete/out of order span sentence"}
				}
				count = 0
				payload = ""
			}
			payload += tokens[5]
			cache[ccount-1] = sentence
			count++
			if ccount == 1 { // First message in sequence, get size and id
				size = tokens[1]
				id = tokens[3]
			} else if size == tokens[2] && count == ccount { // Last message in sequence, send it and clean up.
				padding, _ = strconv.Atoi(tokens[6][:1])
				out <- Message{MessageType(payload), payload, uint8(padding)}
				count = 0
				payload = ""
			}
		}
	}
	out <- Message{255, "", 0}
}

// DecodePositionMessage decodes [the payload of] an AIS position message (type 1/2/3)
func DecodePositionMessage(payload string) (PositionMessage, error) {
	data := []byte(payload)
	var m PositionMessage

	m.Type = decodeAisChar(data[0])

	if m.Type != 1 && m.Type != 2 && m.Type != 3 {
		return m, errors.New("Message isn't Position Report.")
	}

	m.Repeat = decodeAisChar(data[1]) >> 4

	m.MMSI = uint32(decodeAisChar(data[1])) << 28
	m.MMSI = m.MMSI >> 2
	m.MMSI += uint32(decodeAisChar(data[2]))<<20 | uint32(decodeAisChar(data[3]))<<14 | uint32(decodeAisChar(data[4]))<<8 | uint32(decodeAisChar(data[5]))<<2
	m.MMSI += uint32(decodeAisChar(data[6])) >> 4

	m.Status = (decodeAisChar(data[6]) << 4) >> 4

	m.Turn = float32(int8(decodeAisChar(data[7])<<2 | decodeAisChar(data[8])>>4))
	if m.Turn != 0 && m.Turn <= 126 && m.Turn >= -126 {
		sign := float32(1)
		if math.Signbit(float64(m.Turn)) {
			sign = -1
		}
		m.Turn = sign * (m.Turn / 4.733) * (m.Turn / 4.733)

	}

	m.Speed = float32(uint16(decodeAisChar(data[8]))<<12>>6 | uint16(decodeAisChar(data[9])))
	if m.Speed < 1022 {
		m.Speed = m.Speed / 10
	}

	m.Accuracy = false
	if decodeAisChar(data[10])>>5 == 1 {
		m.Accuracy = true
	}

	m.Lon = float64((int32(decodeAisChar(data[10]))<<27 | int32(decodeAisChar(data[11]))<<21 |
		int32(decodeAisChar(data[12]))<<15 | int32(decodeAisChar(data[13]))<<9 | int32(decodeAisChar(data[14]))>>1<<4)) / 16
	m.Lat = float64((int32(decodeAisChar(data[14]))<<31 | int32(decodeAisChar(data[15]))<<25 |
		int32(decodeAisChar(data[16]))<<19 | int32(decodeAisChar(data[17]))<<13 | int32(decodeAisChar(data[18]))<<7 | int32(decodeAisChar(data[19]))>>4<<5)) / 32
	m.Lon, m.Lat = CoordinatesMin2Deg(m.Lon, m.Lat)

	m.Course = float32(uint16(decodeAisChar(data[19]))<<12>>4|uint16(decodeAisChar(data[20]))<<2|uint16(decodeAisChar(data[21]))>>4) / 10

	m.Heading = uint16(decodeAisChar(data[21]))<<12>>7 | uint16(decodeAisChar(data[22]))>>1

	m.Second = decodeAisChar(data[22])<<7>>2 | decodeAisChar(data[23])>>1

	m.Maneuver = decodeAisChar(data[23])<<7>>6 | decodeAisChar(data[24])>>5

	m.RAIM = false
	if decodeAisChar(data[24])<<6>>7 == 1 {
		m.RAIM = true
	}

	return m, nil
}

// GetReferenceTime takes [the payload of] an AIS Base Station message (type 4) and returns the time data of it.
func GetReferenceTime(payload string) (time.Time, error) {
	data := []byte(payload)

	year := uint16(decodeAisChar(data[6]))<<12>>2 | uint16(decodeAisChar(data[7]))<<4 | uint16(decodeAisChar(data[8]))>>2

	if year == 0 {
		var t time.Time
		return t, errors.New("station doesn't report time")
	}

	month := decodeAisChar(data[8])<<6>>4 | decodeAisChar(data[9])>>4

	day := decodeAisChar(data[9])<<4>>3 | decodeAisChar(data[10])>>5

	hour := decodeAisChar(data[10]) << 3 >> 3

	minute := decodeAisChar(data[11])

	second := decodeAisChar(data[12])

	timeString := fmt.Sprintf("%d/%d/%d %d:%d:%d", year, month, day, hour, minute, second)
	t, _ := time.Parse("2006/1/2 15:4:5", timeString)

	return t, nil
}

// CoordinatesMin2Deg translates coordinates (lon, lat) in decimal minutes (×10^4) to decimal degrees.
// AIS data use decimal minutes but decimal degrees (DD) is a more universal format and easier to handle. Almost every third party asks for this format.
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

	degrees := float64(int(minLon / 600000))
	minutes := float64(minLon-600000*degrees) / 10000
	lon := degrees + minutes/60

	degrees = float64(int(minLat / 600000))
	minutes = float64(minLat-600000*degrees) / 10000
	lat := degrees + minutes/60

	return lonSign * lon, latSign * lat
}

// CoordinatesDeg2Human takes coordinates (lon, lat) in decimal degrees (DD),
// formats them as degrees minutes and returns them as string.
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
	minutes := 60 * (degLon - degrees)

	if degrees > 180 {
		coordinates = "longitude not available, "
	} else if lonSign > 0 {
		coordinates = fmt.Sprintf("%3.0f°%07.4f'%s", degrees, minutes, "E")
	} else {
		coordinates = fmt.Sprintf("%3.0f°%07.4f'%s", degrees, minutes, "W")
	}

	degrees = math.Floor(degLat)
	minutes = 60 * (degLat - degrees)

	if degrees > 90 {
		coordinates += "latitude not available"
	} else if latSign > 0 {
		coordinates += fmt.Sprintf(" %3.0f°%07.4f%s", degrees, minutes, "N")
	} else {
		coordinates += fmt.Sprintf(" %3.0f°%07.4f%s", degrees, minutes, "S")
	}

	return coordinates
}

// PrintPositionData returns a formatted string with the detailed data of a AIS position message.
// Its main use is to act as a guide for any developer wishing to correctly parse an AIS position message,
// since some parts of a message are enumareted, and other parts although they mainly are numeric values,
// for certain values they can have a non-numeric meaning.
func PrintPositionData(m PositionMessage) string {

	status := []string{"Under way using engine", "At anchor", "Not under command", "Restricted maneuverability", "Constrained by her draught",
		"Moored", "Aground", "Engaged in fishing", "Under way sailing", "status code reserved", "status code reserved", "status code reserved",
		"status code reserved", "status code reserved", "AIS-SART is active", "Not defined"}

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
		fmt.Sprintf("=== Message Type %d ===\n", m.Type) +
			fmt.Sprintf(" Repeat       : %d\n", m.Repeat) +
			fmt.Sprintf(" MMSI         : %d\n", m.MMSI) +
			fmt.Sprintf(" Nav.Status   : %s\n", status[m.Status]) +
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

// Nmea183ChecksumCheck performs a checksum check for NMEA183 sentences.
// AIS messages are NMEA183 encoded.
func Nmea183ChecksumCheck(sentence string) bool {
	length := len(sentence)

	csum, err := hex.DecodeString(sentence[length-2:])

	if err != nil {
		return false
	}

	bline := []byte(sentence[1 : length-3])
	ccsum := bline[0]

	for i := 1; i < len(bline); i++ {
		ccsum ^= bline[i]
	}

	if csum[0] == ccsum {
		return true
	}
	return false
}
