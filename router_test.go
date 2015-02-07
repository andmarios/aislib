package ais

import (
	"fmt"
	"testing"
)

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
