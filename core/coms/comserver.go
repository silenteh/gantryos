package coms

import (
	log "github.com/golang/glog"

	"net"
	"os"
)

// this is the module responsible for setting up a communication channel (TCP or UDP)
// where the data (protobuf, or JSON) can be exchanged

func StartTCP(port string) {

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Errorln("Failed to listen: %s", err)
	}
	for {
		if conn, err := ln.Accept(); err == nil {
			go handleTCPConnection(conn)
		}
	}

}

func StartUDP(port int) {

}

// Handles incoming requests.
func handleTCPRequest(conn net.Conn) {

	defer conn.Close()
	reader := bufio.NewReader(conn)

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
