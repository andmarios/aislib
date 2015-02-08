package ais

// Some fields are common across different type of messages. Thus here are functions
// to decode them.

// cbnCoordinates takes the start of the coordinates block and returns coordinates in
// decimal degrees
func cbnCoordinates(first int, data []byte) (float64, float64) {
	lon := float64((int32(bitsToInt(first, first+27, data)) << 4)) / 16
	lat := float64((int32(bitsToInt(first+28, first+54, data)) << 5)) / 32

	return CoordinatesMin2Deg(lon, lat)
}

// cbnSpeed takes the start of the speed block and returns speed in knots or 1023.
func cbnSpeed(first int, data []byte) float32 {
	speed := float32(bitsToInt(first, first + 9, data))
	if speed < 1022 {
		speed /= 10
	}
	return speed
}

// cbnBool decodes a bool value
func cbnBool(bit int, data []byte) bool {
	if bitsToInt(bit, bit, data) == 1 {
		return true
	}
	return false
}
