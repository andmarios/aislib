package ais

import "testing"

func BenchmarkNmea183ChecksumCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Nmea183ChecksumCheck("!AIVDM,1,1,,B,38u<a<?PAA2>P:WfuAO9PW<P0PuQ,0*6F")
	}
}
