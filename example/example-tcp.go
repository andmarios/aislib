package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/marine-travel/marine-ais"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// Here are saved as JSON string the ships seen in the last 5-second period.
var serveJSON string

type shipData struct {
	Data  ais.ClassAPositionReport
	Human string
}

func main() {

	// Create an AIS router process to decode radio sentences
	send := make(chan string, 1024)
	receive := make(chan ais.Message, 1024)
	failed := make(chan ais.FailedSentence, 1024)
	go ais.Router(send, receive, failed)

	// Create a handler-process that reads messages from router, decodes and saves the payload
	seen := make(map[uint32]shipData)
	proceed := make(chan bool)
	go func() {
		var message ais.Message
		var problematic ais.FailedSentence
		for {
			select {
			case message = <-receive:
				if message.Type >= 1 && message.Type <= 3 {
					m, _ := ais.DecodeClassAPositionReport(message.Payload)
					seen[m.MMSI] = shipData{m, ais.PrintClassAPositionReport(m)}
				}
			case problematic = <-failed:
				log.Println(problematic)
			case _ = <-proceed: // Unbuffered channel used for synchronization (as mutex for [seen])
				<-proceed
			}
		}
	}()

	// Create a process that every five seconds refreshes [serveJSON] with new data
	go func() {
		var jsonBuf bytes.Buffer
		for _ = range time.Tick(5 * time.Second) {
			proceed <- true
			ships := seen
			seen = make(map[uint32]shipData)
			proceed <- true
			for _, s := range ships {
				j, _ := json.Marshal(s)
				jsonBuf.Write(j)
				jsonBuf.WriteString(",")
			}
			length := len(jsonBuf.String())
			if length > 10 {
				serveJSON = "[" + jsonBuf.String()[:length-1] + "]"
				jsonBuf.Reset()
				fmt.Println(len(serveJSON))
			} else {
				serveJSON = "[]"
			}
		}
	}()

	// Connect to a remote AIS server. Read AIS sentences and forward them to the AIS router.
	// If connection drops wait and reconnect.
	remote := "ais1.shipraiser.net:6492"
	go func() {
		sleep := 10 // How many seconds to sleep after a timeout
		sleepD := time.Duration(sleep) * time.Second
		for {
			serverAddr, err := net.ResolveTCPAddr("tcp", remote)
			if err != nil {
				log.Println(err, errors.New("(retrying in " + strconv.Itoa(sleep) + " seconds)"))
				time.Sleep(sleepD)
				continue
			}
			conn, err := net.DialTCP("tcp", nil, serverAddr)
			if err != nil {
				log.Println(err, errors.New("(retrying in " + strconv.Itoa(sleep) + " seconds)"))
				time.Sleep(sleepD)
				continue
			}
			defer conn.Close()

			connbuf := bufio.NewScanner(conn)
			connbuf.Split(bufio.ScanLines)
			for connbuf.Scan() {
				send <- connbuf.Text()
				conn.SetReadDeadline(time.Now().Add(15 * time.Second))
			}
			log.Println(remote + ": connection broken (retrying in " + strconv.Itoa(sleep) + " seconds)")
			time.Sleep(sleepD)
		}
	}()

	// Create a server to listen for files/data requests
	http.HandleFunc("/data", dataHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)

}

// Function to serve the ships JSON string
func dataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", serveJSON)
}






