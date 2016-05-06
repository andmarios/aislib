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

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	ais "github.com/andmarios/aislib"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanLines)

	send := make(chan string, 1024*8)
	receive := make(chan ais.Message, 1024*8)
	failed := make(chan ais.FailedSentence, 1024*8)

	done := make(chan bool)

	go ais.Router(send, receive, failed)

	go func() {
		var message ais.Message
		var problematic ais.FailedSentence
		for {
			select {
			case message = <-receive:
				switch message.Type {
				case 1, 2, 3:
					t, _ := ais.DecodeClassAPositionReport(message.Payload)
					fmt.Println(t)
				case 4:
					t, _ := ais.DecodeBaseStationReport(message.Payload)
					fmt.Println(t)
				case 5:
					t, _ := ais.DecodeStaticVoyageData(message.Payload)
					fmt.Println(t)
				case 8:
					t, _ := ais.DecodeBinaryBroadcast(message.Payload)
					fmt.Println(t)
				case 18:
					t, _ := ais.DecodeClassBPositionReport(message.Payload)
					fmt.Println(t)
				case 255:
					done <- true
				default:
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
