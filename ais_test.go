package ais

import (
	"fmt"
	"testing"
	"time"
)

// We don't decode radio field, so testing on radio field isn't really verified.
// Radio is a 19bit field which we decode as uint32. One would need to decode further
// these 19 bits but its uneccesary for now.
func TestDecodeClassAPositionReport(t *testing.T) {
	cases := []struct {
		payload string
		want    ClassAPositionReport
	}{
		{
			"38u<a<?PAA2>P:WfuAO9PW<P0PuQ",
			ClassAPositionReport{
				PositionReport: PositionReport{Type: 3, Repeat: 0, MMSI: 601041200, Speed: 8.1,
					Accuracy: false, Lon: 31.130165, Lat: -29.784113333333334, Course: 243.4,
					Heading: 230, Second: 16, RAIM: false, Radio: 135009},
				Status: 15, Turn: -127, Maneuver: 0},
		},
		{
			"13P:v?h009Ogbr4NkiITkU>L089D",
			ClassAPositionReport{
				PositionReport: PositionReport{Type: 1, Repeat: 0, MMSI: 235060799, Speed: 0.9,
					Accuracy: false, Lon: -3.56725, Lat: 53.84251666666667, Course: 123,
					Heading: 167, Second: 14, RAIM: false, Radio: 33364},
				Status: 0, Turn: 0, Maneuver: 0},
		},
		{
			"13n@oD0PB@0IRqvQj@W;EppH088t19uvPT",
			ClassAPositionReport{
				PositionReport: PositionReport{Type: 1, Repeat: 0, MMSI: 258226000, Speed: 14.4,
					Accuracy: false, Lon: 5.580478333333334, Lat: 59.0441, Course: 290.3,
					Heading: 284, Second: 12, RAIM: false, Radio: 33340},
				Status: 0, Turn: -127, Maneuver: 0},
		},
	}
	for _, c := range cases {
		got, _ := DecodeClassAPositionReport(c.payload)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("DecodeClassAPositionReport(payload string)")
		}
	}
}

func BenchmarkDecodeClassAPositionReport(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecodeClassAPositionReport("38u<a<?PAA2>P:WfuAO9PW<P0PuQ")
	}
}

func TestDecodeBaseStationReport(t *testing.T) {
	caseTime1, _ := time.Parse("2006/1/2 15:4:5", "2015/2/4 0:33:51")
	caseTime2, _ := time.Parse("2006/1/2 15:4:5", "2015/2/4 0:33:50")
	cases := []struct {
		payload string
		want    BaseStationReport
	}{
		{
			"402R3KiutR0Qk156V4QQTOA00<0;",
			BaseStationReport{Repeat: 0, MMSI: 2655087, Time: caseTime1, Accuracy: false, Lon: 15.09579, Lat: 58.588368333333335, EPFD: 1, RAIM: false, Radio: 49163},
		},
		{
			"4025boiutR0Qj0qgK<OodKW00@N1",
			BaseStationReport{Repeat: 0, MMSI: 2190047, Time: caseTime2, Accuracy: false, Lon: 12.613716666666667, Lat: 55.69725, EPFD: 7, RAIM: false, Radio: 67457},
		},
	}
	for _, c := range cases {
		got, _ := DecodeBaseStationReport(c.payload)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("DecodeClassAPositionReport(payload string)")
		}
	}
}

func BenchmarkDecodeBaseStationReport(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecodeBaseStationReport("402R3KiutR0Qk156V4QQTOA00<0;")
	}
}

func TestCoordinatesDeg2Human(t *testing.T) {
	cases := []struct {
		lon, lat float64
		want     string
	}{
		{-3.56725, 53.84251666666667, "  3°34.0350'W  53°50.5510N"},
		{31.130165, -29.784113333333334, " 31°07.8099'E  29°47.0468S"},
	}
	for _, c := range cases {
		got := CoordinatesDeg2Human(c.lon, c.lat)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("CoordinatesDeg2Human(lon, lat float64)")
		}
	}
}

func ExampleCoordinatesDeg2Human() {
	fmt.Println(CoordinatesDeg2Human(-3.56725, 53.84251666666667))
	// Output:   3°34.0350'W  53°50.5510N
}

func TestGetReferenceTime(t *testing.T) {
	cases := []struct {
		payload, reference string
	}{
		{"4025;PAuho;N>0NJbfMRhNA00D3l", "2012/3/14 11:30:14"},
		{"403tDGiuho;P5<tSF0l4Q@000l67", "2012/3/14 11:32:5"},
	}
	for _, c := range cases {
		got, _ := GetReferenceTime(c.payload)
		want, _ := time.Parse("2006/1/2 15:4:5", c.reference)
		if got != want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", want)
			t.Errorf("GetReferenceTime(payload string)")
		}
	}
}

func BenchmarkGetReferenceTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetReferenceTime("4025;PAuho;N>0NJbfMRhNA00D3l")
	}
}

func TestRouter(t *testing.T) {
	cases := []struct {
		message  Message
		sentence []string
	}{
		{
			Message{3, "38u<a<?PAA2>P:WfuAO9PW<P0PuQ", 0},
			[]string{"!AIVDM,1,1,,B,38u<a<?PAA2>P:WfuAO9PW<P0PuQ,0*6F"},
		},
		{
			Message{5, "533iFNT00003W;3G;384iT<T400000000000001?88?73v0ik0RC1H11H30H51CU0E2CkP0", 2},
			[]string{"!AIVDM,2,1,5,A,533iFNT00003W;3G;384iT<T400000000000001?88?73v0ik0RC1H11H30H,0*44",
				"!AIVDM,2,2,5,A,51CU0E2CkP0,2*0C"},
		},
		{
			Message{8, "85Mwom1KfI?GR<NgcvM1Hg<P2FaGjRN<S22j;WN:IDle3f5Qsq6=620c;<gvsa8P?;j>Nl0oKaCLIdeFlr<Gh@Jc95:i>c0", 2},
			[]string{"!AIVDM,3,1,7,A,85Mwom1KfI?GR<NgcvM1Hg<P2FaGjRN<S22j;WN:IDl,0*3E",
				"!AIVDM,3,2,7,A,e3f5Qsq6=620c;<gvsa8P?;j>Nl0oKaCLIdeFlr<Gh@,0*3D",
				"!AIVDM,3,3,7,A,Jc95:i>c0,2*08"},
		},
	}

	send := make(chan string)
	receive := make(chan Message, 1024)
	failed := make(chan FailedSentence, 1024)

	go Router(send, receive, failed)

	for _, c := range cases {
		for _, m := range c.sentence {
			send <- m
		}
		got := <-receive
		if got != c.message {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.message)
			t.Errorf("Router(in chan string, out chan Message, failed chan FailedSentence)")
		}
	}
}

func BenchmarkRouter(b *testing.B) {
	send := make(chan string)
	receive := make(chan Message, 1024)
	failed := make(chan FailedSentence, 1024)

	go Router(send, receive, failed)

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			send <- "!AIVDM,1,1,,B,38u<a<?PAA2>P:WfuAO9PW<P0PuQ,0*6F"
			<-receive
		} else {
			send <- "!AIVDM,2,1,5,A,533iFNT00003W;3G;384iT<T400000000000001?88?73v0ik0RC1H11H30H,0*44"
			send <- "!AIVDM,2,2,5,A,51CU0E2CkP0,2*0C"
			<-receive
		}
	}
}

func BenchmarkMessageType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MessageType("38u<a<?PAA2>P:WfuAO9PW<P0PuQ")
	}
}

func BenchmarkNmea183ChecksumCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Nmea183ChecksumCheck("!AIVDM,1,1,,B,38u<a<?PAA2>P:WfuAO9PW<P0PuQ,0*6F")
	}
}

// I didn't found a MMSI decoder, so this test is verified only by me
func TestDecodeMMSI(t *testing.T) {
	cases := []struct {
		MMSI uint32
		reference string
	}{
		{227006760, "Ship, France"},
		{2573425, "Coastal Station, Norway"},
		{25634906, "Group of ships, Malta"},
		{842517724, "Diver's radio, Iraq (Republic of)"},
		{992351000, "Aids to navigation, United Kingdom of Great Britain and Northern Ireland"},
		{1000010000, "Invalid MMSI"},
		{972345000, "MOB —Man Overboard Device"},
		{970241023, "AIS SART —Search and Rescue Transmitter, Greece"},
		{971356034, "Invalid MMSI"},
	}
	for _, c := range cases {
		got := DecodeMMSI(c.MMSI)
		if got != c.reference {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.reference)
			t.Errorf("DecodeMMSI(m uint32)")
		}
	}
}

func BenchmarkDecodeMMSI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		switch {
		case  i%116 < 108:
			DecodeMMSI(316013198)
		case  i%116 < 115:
			DecodeMMSI(002241076)
		default:
			DecodeMMSI(992351000)
		}
	}
}

