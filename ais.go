// Package ais provides functions and types to work with AIS (Automatic Identification System)
// sentences (radio messages) and messages in AIVDM/AIVDO format.
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

// A Message stores the important properties of a AIS message, including only information useful
// for decoding: Type, Payload, Padding Bits
// A Message should come after processing one or more AIS radio sentences (checksum check,
// concatenate payloads spanning across sentences, etc).
type Message struct {
	Type    uint8
	Payload string
	Padding uint8
}

// FailedSentence includes an AIS sentence that failed to process (e.g wrong checksum) and the reason
// it failed.
type FailedSentence struct {
	Sentence string
	Issue    string
}

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

// Navigation status codes
var NavigationStatusCodes = [...]string{
	"Under way using engine", "At anchor", "Not under command", "Restricted maneuverability",
	"Constrained by her draught", "Moored", "Aground", "Engaged in fishing", "Under way sailing",
	"status code reserved", "status code reserved", "status code reserved",
	"status code reserved", "status code reserved", "AIS-SART is active", "Not defined",
}

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
			if count > 1 { // Invalidate cache
				for i := 0; i <= count; i++ {
					failed <- FailedSentence{cache[i], "Incomplete/out of order span sentence"}
				}
				count = 0
				payload = ""
			}
		} else { // Message spans across sentences.
			ccount, err = strconv.Atoi(tokens[2])
			if err != nil {
				failed <- FailedSentence{sentence, "HERE " + tokens[2]}
				continue
			}
			if ccount != count+1 || // If there are sentences with wrong seq.number in cache send them as failed
				tokens[3] != id && count != 0 || // If there are sentences with different sequence id in cache , send old parts as failed
				tokens[1] != size && count != 0 { // If there messages with wrong size in cache, send them as failed
				for i := 0; i < count; i++ {
					failed <- FailedSentence{cache[i], "Incomplete/out of order span sentence"}
				}
				if ccount != 1 { // The current one is invalid too
					failed <- FailedSentence{sentence, "Incomplete/out of order span sentence"}
					count = 0
					payload = ""
					continue
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

// DecodeClassAPositionReport decodes [the payload of] an AIS position message (type 1/2/3)
func DecodeClassAPositionReport(payload string) (ClassAPositionReport, error) {
	data := []byte(payload)
	var m ClassAPositionReport

	m.Type = decodeAisChar(data[0])
	if m.Type != 1 && m.Type != 2 && m.Type != 3 {
		return m, errors.New("Message isn't Position Report.")
	}

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
	return m, nil
}

// DecodeBaseStationReport decodes the payload of a Type 4 AIS message
func DecodeBaseStationReport(payload string) (BaseStationReport, error) {
	data := []byte(payload)
	var m BaseStationReport

	mType := decodeAisChar(data[0])
	if mType != 4 {
		return m, errors.New("Message isn't Base Station Report.")
	}

	//m.Repeat = decodeAisChar(data[1]) >> 4
	m.Repeat = uint8(bitsToInt(6, 7, data))

	//m.MMSI = uint32(decodeAisChar(data[1]))<<28>>2 | uint32(decodeAisChar(data[2]))<<20 |
	//	uint32(decodeAisChar(data[3]))<<14 | uint32(decodeAisChar(data[4]))<<8 |
	//	uint32(decodeAisChar(data[5]))<<2 | uint32(decodeAisChar(data[6]))>>4
	m.MMSI = bitsToInt(8, 37, data)

	m.Time, _ = GetReferenceTime(payload) // Some base stations do not report time, for this case we do not consider it as error

	m.Accuracy = false
	if decodeAisChar(data[13])>>5 == 1 {
		m.Accuracy = true
	}

	m.Lon = float64((int32(bitsToInt(79, 106, data)) << 4)) / 16
	m.Lat = float64((int32(bitsToInt(107, 133, data)) << 5)) / 32
	m.Lon, m.Lat = CoordinatesMin2Deg(m.Lon, m.Lat)

	m.EPFD = uint8(bitsToInt(134, 137, data))

	m.RAIM = false
	if bitsToInt(148, 148, data) == 1 {
		m.RAIM = true
	}

	m.Radio = bitsToInt(149, 167, data)
	return m, nil
}

// GetReferenceTime takes [the payload of] an AIS Base Station message (type 4)
// and returns the time data of it. It is a separate function from DecodeBaseStationReport
// because it can be useful to set a timeframe for our received AIS messages.
func GetReferenceTime(payload string) (time.Time, error) {
	data := []byte(payload)

	year := uint16(decodeAisChar(data[6]))<<12>>2 | uint16(decodeAisChar(data[7]))<<4 |
		uint16(decodeAisChar(data[8]))>>2
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
// AIS data use decimal minutes but decimal degrees (DD) is a more universal format and easier to
// handle. Almost every third party asks for this format.
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
			fmt.Sprintf(" MMSI         : %s\n", PrintMMSI(m.MMSI)) +
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

// PrintBaseStationReport returns a formatted string of a BaseStationReport. Mainly to help
// developers with understanding base position reports.
func PrintBaseStationReport(m BaseStationReport) string {
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
			fmt.Sprintf(" MMSI         : %s\n", PrintMMSI(m.MMSI)) +
			fmt.Sprintf(" Time         : %s\n", m.Time.String()) +
			fmt.Sprintf(" Accuracy     : %s\n", accuracy) +
			fmt.Sprintf(" Coordinates  : %s\n", CoordinatesDeg2Human(m.Lon, m.Lat)) +
			fmt.Sprintf(" EPFD         : %s\n", EpfdFixTypes[m.EPFD]) +
			fmt.Sprintf(" RAIM         : %s\n", raim)

	return message
}

// PrintMMSI returns a string with the type of the owner of the MMSI and its country
// Some MMSIs aren't valid. There is some more information in some MMSIs (the satellite
// equipment of the ship). We may add them in the future.
// Have a look at http://en.wikipedia.org/wiki/Maritime_Mobile_Service_Identity
func PrintMMSI(m uint32) string {
	mid := fmt.Sprintf("%09d", m)
	data := ""

	switch mid[0:1] {
	case "0":
		if mid[1:2] == "0" {
			m = m / 10000
			data = "Coastal Station, " + Mid[int(m)]
		} else {
			m = m / 100000
			data = "Group of ships,  " + Mid[int(m)]
		}
	case "1":
		m = m / 1000 - 111000
		data = "SAR —Search and Rescue Aircraft, " + Mid[int(m)]
	case "2", "3", "4", "5", "6", "7":
		m = m / 1000000
		data = "Ship, " + Mid[int(m)]
	case "8":
		m = m / 100000 - 8000
		data = "Diver's radio, " + Mid[int(m)]
	case "9":
		if mid[1:2] == "9" {
			m = m / 10000 - 99000
			data = "Aids to navigation, " + Mid[int(m)]
		} else if mid[1:2] == "8" {
			m = m / 10000 - 98000
			data = "Auxiliary craft associated with parent ship, " + Mid[int(m)]
		} else if mid[1:3] == "970" {
			m = m / 1000 - 970000
			data = "AIS SART —Search and Rescue Transmitter, " + Mid[int(m)]
		} else if mid[1:3] == "972" {
			data = "MOB —Man Overboard Device"
		} else if mid[1:3] == "974" {
			data = "EPIRB —Emergency Position Indicating Radio Beacon"
		}
	}
	return data + " [" + mid + "]"
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

// bitsToInt extracts certain bits from a payload.
// Payload consists of six bit packets, each one armored in one byte.
// The function seems simple enough but took me some hours to figure out.
// It is necessary since this is the most frequent part of the message decoding process
// and one can only write so many binary operations. We sacrifice performance a bit to
// significantly ease development. :-/
func bitsToInt(first, last int, payload []byte) uint32 {
	size := uint(last - first) // Bit fields start at 0
	processed, remain := uint(0), uint(0)
	result, temp := uint32(0), uint32(0)

	from := first / 6
	forTimes := last/6 - from

	for i := 0; i <= forTimes; i++ {
		// Instead of calling decodeAisChar we do the calculation manually here for speed.
		temp = uint32(payload[from+i]) - 48
		if temp > 40 {
			temp -= 8
		}

		// Depending on which byte in sequence we processing, do the appropriate shifts.
		if i == 0 { // For the first byte we (may) have to clean leftmost bits and shift to position
			remain = uint(first%6) + 1
			processed = 6 - remain
			temp = temp << (31 - processed) >> (31 - size)
		} else if i < forTimes { // For middle bytes we only shift to position
			processed = processed + 6
			temp = temp << (size - processed)
		} else { // For last byte we (may) clear rightmost bits
			remain = uint(last%6) + 1
			temp = temp >> (6 - remain)
		}
		result = result | temp
	}
	return result
}
