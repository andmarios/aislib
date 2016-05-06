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

import "testing"

func BenchmarkNmea183ChecksumCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Nmea183ChecksumCheck("!AIVDM,1,1,,B,38u<a<?PAA2>P:WfuAO9PW<P0PuQ,0*6F")
	}
}
