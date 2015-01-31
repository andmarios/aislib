package main

import (
	"bufio"
//	"io"
	"log"
	"os"
	//	"unicode/utf8"
	"strings"
	"encoding/hex"
)

type aisMessageT123 struct {
	Type   int
	Repeat int
	MMSI   string
	status int //navigation status
	Turn int //rate of turn - ROT
	Speed int //speed over ground - SOG
	Accuracy boolean //position accuracy
	Lon int
	Lat int
	Course float //course over ground - COG
	Heading int // true heading - HDG
	Second int // timestamp
	Manuveur int // maneuver indicator
	Raim bool // RAIM flag
	Radio int // Radio status
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanLines)

	for in.Scan() {

		line := in.Text()
		//llength := len(line)
		tokens := strings.Split(line, ",")

		log.Println("Line length:", len(line), "Tokens:", len(tokens), "Payload:", tokens[5], "Checksum Passed:", nmea183ChecksumCheck(line))
	}

}

func nmea183ChecksumCheck(sentence string) bool {
	length := len(sentence)

	var csum []byte
	csum, err := hex.DecodeString(sentence[length-2:length])

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
