package ais

import (
	"strconv"
	"strings"
)

// A Message stores the important properties of a AIS message, including only information useful
// for decoding: Type, Payload, Padding Bits
// A Message should come after processing one or more AIS radio sentences (checksum check,
// concatenate payloads spanning across sentences, etc).
type Message struct {
	Type    uint8
	Payload string
	Padding uint8
}

// FailedSentence includes an AIS sentence that failed to process (e.g wrong checksum) and the reason
// it failed.
type FailedSentence struct {
	Sentence string
	Issue    string
}

// Router accepts AIS radio sentences and process them. It checks their checksum,
// and AIS identifiers. If they are valid it tries to assemble the payload if it spans
// on multiple sentences. Upon success it returns the AIS Message at the out channel.
// Failed sentences go to the err channel.
// If the in channel is closed, then it sends a message with type 255 at the out channel.
// Your function can check for this message to know when it is safe to exit the program.
func Router(in chan string, out chan Message, failed chan FailedSentence) {
	count, ccount, padding := 0, 0, 0
	size, id := "0", "0"
	payload := ""
	var cache [5]string
	var err error
	aisIdentifiers := map[string]bool{
		"ABVD": true, "ADVD": true, "AIVD": true, "ANVD": true, "ARVD": true,
		"ASVD": true, "ATVD": true, "AXVD": true, "BSVD": true, "SAVD": true,
	}
	for sentence := range in {
		if len(sentence) == 0 { // Do not process empty lines
			failed <- FailedSentence{sentence, "Empty line"}
			continue
		}
		tokens := strings.Split(sentence, ",") // I think this takes the major portion of time for this function (after benchmarking)

		if !Nmea183ChecksumCheck(sentence) { // Checksum check
			failed <- FailedSentence{sentence, "Checksum failed"}
			continue
		}

		if !aisIdentifiers[tokens[0][1:5]] { // Check for valid AIS identifier
			failed <- FailedSentence{sentence, "Sentence isn't AIVDM/AIVDO"}
			continue
		}

		if tokens[1] == "1" { // One sentence message, process it immediately
			padding, _ = strconv.Atoi(tokens[6][:1])
			out <- Message{MessageType(tokens[5]), tokens[5], uint8(padding)}
			if count > 1 { // Invalidate cache
				for i := 0; i < count; i++ {
					failed <- FailedSentence{cache[i], "Incomplete/out of order span sentence"}
				}
				count = 0
				payload = ""
			}
		} else { // Message spans across sentences.
			ccount, err = strconv.Atoi(tokens[2])
			if err != nil {
				failed <- FailedSentence{sentence, "HERE " + tokens[2]}
				continue
			}
			if ccount != count+1 || // If there are sentences with wrong seq.number in cache send them as failed
				tokens[3] != id && count != 0 || // If there are sentences with different sequence id in cache , send old parts as failed
				tokens[1] != size && count != 0 { // If there messages with wrong size in cache, send them as failed
				for i := 0; i < count; i++ {
					failed <- FailedSentence{cache[i], "Incomplete/out of order span sentence"}
				}
				if ccount != 1 { // The current one is invalid too
					failed <- FailedSentence{sentence, "Incomplete/out of order span sentence"}
					count = 0
					payload = ""
					continue
				}
				count = 0
				payload = ""
			}
			payload += tokens[5]
			cache[ccount-1] = sentence
			count++
			if ccount == 1 { // First message in sequence, get size and id
				size = tokens[1]
				id = tokens[3]
			} else if size == tokens[2] && count == ccount { // Last message in sequence, send it and clean up.
				padding, _ = strconv.Atoi(tokens[6][:1])
				out <- Message{MessageType(payload), payload, uint8(padding)}
				count = 0
				payload = ""
			}
		}
	}
	out <- Message{255, "", 0}
}
