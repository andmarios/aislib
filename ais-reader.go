package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"encoding/hex"
)

// Please have a look at http://catb.org/gpsd/AIVDM.html
type aisMessageT123 struct {
	Type     uint8
	Repeat   uint8
	MMSI     uint32
	Status   uint8   // navigation status (enumerated type)
	Turn     float32 // rate of turn - ROT (sc - Special Calc I3)
	Speed    float32 // speed over ground - SOG (sc U3)
	Accuracy bool    // position accuracy
	Lon      float32 // (sc I4)
	Lat      float32 // (sc I4)
	Course   float32 //course over ground - COG (sc U1)
	Heading  uint16  // true heading - HDG
	Second   uint8   // timestamp
	Manuveur uint8  // maneuver indicator (enumerated)
	RAIM     bool   // RAIM flag
	Radio    uint32 // Radio status
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanLines)

	for in.Scan() {

		line := in.Text()
		tokens := strings.Split(line, ",")

		log.Println("Line length:", len(line), "Tokens:", len(tokens), "Payload:", tokens[5], "Checksum:", nmea183ChecksumCheck(line))
	}

}

func nmea183ChecksumCheck(sentence string) bool {
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
