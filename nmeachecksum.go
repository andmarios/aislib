package ais

import "encoding/hex"

// Nmea183ChecksumCheck performs a checksum check for NMEA183 sentences.
// AIS messages are NMEA183 encoded.
func Nmea183ChecksumCheck(sentence string) bool {
	length := len(sentence)
	if length < 5 { // Sentence isn't long enough to have a csum, avoid bounds out of range
		return false
	}

	// Read the checksum from the AIS sentence
	csum, err := hex.DecodeString(sentence[length-2:])
	if err != nil {
		return false
	}

	// The checksum is calculated from the whole sentence except
	// the first and last three characters
	bline := []byte(sentence[1 : length-3])
	ccsum := bline[0]
	// The checksum is calculated by XOR'ing all the characters
	for i := 1; i < len(bline); i++ {
		ccsum ^= bline[i]
	}

	if csum[0] == ccsum {
		return true
	}
	return false
}
