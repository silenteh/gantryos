package coms

import (
	log "github.com/golang/glog"

	"net"
	"os"
)

// this is the module responsible for setting up a communication channel (TCP or UDP)
// where the data (protobuf, or JSON) can be exchanged

func StartTCP(port int) {

}

func StartUDP(port int) {

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 512)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}
