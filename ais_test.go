package ais

import (
	"fmt"
	"testing"
	"time"
)

func TestDecodePositionMessage(t *testing.T) {
	cases := []struct {
		payload string
		want    PositionMessage
	}{
		{
			"38u<a<?PAA2>P:WfuAO9PW<P0PuQ",
			PositionMessage{Type: 3, Repeat: 0, MMSI: 601041200, Status: 15, Turn: -127, Speed: 8.1, Accuracy: false, Lon: 31.130165, Lat: -29.784113333333334, Course: 243.4, Heading: 230, Second: 16, Maneuver: 0, RAIM: false},
		},
		{
			"13P:v?h009Ogbr4NkiITkU>L089D",
			PositionMessage{Type: 1, Repeat: 0, MMSI: 235060799, Status: 0, Turn: 0, Speed: 0.9, Accuracy: false, Lon: -3.56725, Lat: 53.84251666666667, Course: 123, Heading: 167, Second: 14, Maneuver: 0, RAIM: false},
		},
		{
			"13n@oD0PB@0IRqvQj@W;EppH088t19uvPT",
			PositionMessage{Type: 1, Repeat: 0, MMSI: 258226000, Status: 0, Turn: -127, Speed: 14.4, Accuracy: false, Lon: 5.580478333333334, Lat: 59.0441, Course: 290.3, Heading: 284, Second: 12, Maneuver: 0, RAIM: false},
		},
	}
	for _, c := range cases {
		got, _ := DecodePositionMessage(c.payload)
		if got != c.want {
			fmt.Println("Got : ", got)
			fmt.Println("Want: ", c.want)
			t.Errorf("DecodePositionMessage(payload string)")
		}
	}
}

func BenchmarkDecodePositionMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecodePositionMessage("38u<a<?PAA2>P:WfuAO9PW<P0PuQ")
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
			[]string{"!AIVDM,2,1,5,A,533iFNT00003W;3G;384iT<T400000000000001?88?73v0ik0RC1H11H30H,0*44", "!AIVDM,2,2,5,A,51CU0E2CkP0,2*0C"},
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
		if i % 2 == 0 {
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
