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

import "strings"

// decodeAisChar takes a byte a returns the six bit field of AIS data.
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

	if len(payload)*6 < last+1 { // There are strange messages out there, this seems to be what decoders do
		return 0
	}
	for i := 0; i <= forTimes; i++ {
		// Instead of calling decodeAisChar we do the calculation manually here for speed.
		temp = uint32(payload[from+i]) - 48
		if temp > 40 {
			temp -= 8
		}

		// Depending on which byte in sequence we processing, do the appropriate shifts.
		if i == 0 { // For the first byte we (may) have to clean leftmost bits and shift to position
			remain = uint(first % 6)
			processed = 5 - remain
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

// bitsToString decodes text from an AIS payload. Text is packed in six bit ASCII
func bitsToString(first, last int, payload []byte) string {
	length := (last - first + 1) / 6 // How many characters we expect
	start := first / 6               // At which byte the first character starts
	var text [64]byte                // Not sure which the maximum text field size is, but this should be enough
	char := uint8(0)

	// Some times we get truncated text fields. Since text fields have constant size,
	// it is frequent that they aren't fully occupied. Transmitters use this to send shorter messages.
	// We should handle this gracefully, adjusting the length of the text we expect to read.
	if len(payload)*6 < last+1 {
		if len(payload)*6 < first+5 { // Haven't seen this case yet (text field missing) but better be prepared
			return ""
		}
		// Do not simplify this. It uses the uint type rounding method to get correct results
		length = (len(payload)*6 - first) / 6
	}

	remain := first % 6

	// In this if/else there is some code duplication but I think the speed enhancement is worth it.
	// The other way around would need 2*length branches. Now we have only 2.
	// decodeAisChar function should be safe to use here since we check the payload's length
	if remain < 6 {
		shiftLeftMost := uint8(remain + 2)
		shiftRightMost := uint8(6 - remain)
		for i := 0; i < length; i++ {
			char = decodeAisChar(payload[start+i])<<shiftLeftMost>>2 |
				decodeAisChar(payload[start+i+1])>>shiftRightMost
			if char < 32 {
				char += 64
			}
			text[i] = char
		}
	} else {
		for i := 0; i < length; i++ {
			char = decodeAisChar(payload[start+i])
			if char < 32 {
				char += 64
			}
			text[i] = char
		}
	}

	// We convert to string and trim the righmost spaces and @ according to the format specs.
	return strings.TrimRight(string(text[:length]), "@ ")
}
