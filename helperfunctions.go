package ais

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
