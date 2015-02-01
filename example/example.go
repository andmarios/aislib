package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"projects.30ohm.com/mrsaccess/ais"
	"fmt"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanLines)

	for in.Scan() {

		line := in.Text()

		if ais.Nmea183ChecksumCheck(line) {

			tokens := strings.Split(line, ",")
			if tokens[0] == "!AIVDM" && // Sentence is ais data
				tokens[1] == "1" &&     // Payload doesn't span across two sentences (ok for messages 1/2/3)
				tokens[6][:1] == "0" {  // Message doesn't need weird padding (ok for messages 1/2/3)

				//log.Println("Line length:", len(line), "Tokens:", len(tokens), "Payload:", tokens[5], "Checksum:", nmea183ChecksumCheck(line))

				message, err := ais.DecodeAisPosition(tokens[5])
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println(ais.PrintAisPositionData(message))
				}
			} else {
				log.Println("There was an error with message:", line)
			}

		} else {
			log.Println("Checksum failed:", line)
		}
	}

}
