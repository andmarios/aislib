// Copyright (c) 2015, Marios Andreopoulos.
//
// This file is part of aislib.
//
//  Aislib is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
//  Aislib is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
// along with aislib.  If not, see <http://www.gnu.org/licenses/>.

package aislib

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
