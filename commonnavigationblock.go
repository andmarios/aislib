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
	speed := float32(bitsToInt(first, first+9, data))
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
