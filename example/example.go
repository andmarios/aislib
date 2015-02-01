package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"projects.30ohm.com/mrsaccess/ais"
	"strings"
	//	"time"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanLines)

	for in.Scan() {

		line := in.Text()

		if ais.Nmea183ChecksumCheck(line) {

			tokens := strings.Split(line, ",")
			if tokens[0] == "!AIVDM" && // Sentence is ais data
				tokens[1] == "1" && // Payload doesn't span across two sentences (ok for messages 1/2/3)
				tokens[6][:1] == "0" { // Message doesn't need weird padding (ok for messages 1/2/3)

				messageType := ais.AisMessageType(tokens[5])
				if messageType >= 1 && messageType <= 3 {

					message, err := ais.DecodeAisPosition(tokens[5])
					if err != nil {
						log.Println(err)
					} else {
						fmt.Println(ais.PrintAisPositionData(message))
					}
				} else if messageType == 4 {
					t, err := ais.GetReferenceTime(tokens[5])
					if err != nil {
						log.Println(err)
					} else {
						fmt.Println("=== Reference Time ===")
						fmt.Println(t)
						fmt.Println()
					}
				}
			} else {
				log.Println("There was an error with message:", line)
			}

		} else {
			log.Println("Checksum failed:", line)
		}
	}

}
