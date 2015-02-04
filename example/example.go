package main

import (
	"bufio"
	"fmt"
	"github.com/marine-travel/marine-ais"
	"log"
	"os"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanLines)

	send := make(chan string, 1024)
	receive := make(chan ais.Message, 1024)
	failed := make(chan ais.FailedSentence, 1024)

	done := make(chan bool)

	go ais.Router(send, receive, failed)

	go func() {
		var message ais.Message
		var problematic ais.FailedSentence
		for {
			select {
			case message = <-receive:
				if message.Type >= 1 && message.Type <= 3 {
					positionMessage, _ := ais.DecodePositionMessage(message.Payload)
					fmt.Println(ais.PrintPositionData(positionMessage))
				} else if message.Type == 4 {
					t, err := ais.GetReferenceTime(message.Payload)
					if err != nil {
						log.Println(err)
					} else {
						fmt.Println("=== Reference Time ===")
						fmt.Println(t)
						fmt.Println()
					}
				} else if message.Type == 255 {
					done <- true
				} else {
					fmt.Printf("=== Message Type %2d ===\n", message.Type)
					fmt.Printf(" Unsupported type \n\n")
				}
			case problematic = <-failed:
				log.Println(problematic)
			}
		}
	}()

	for in.Scan() {
		send <- in.Text()
	}
	close(send)
	<-done
}
