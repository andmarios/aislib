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










