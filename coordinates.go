package ais

import (
	"fmt"
	"math"
)

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
