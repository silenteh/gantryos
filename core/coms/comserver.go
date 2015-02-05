package coms

import (
	"bufio"
	"fmt"
	protobuf "github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
	"io/ioutil"
	"net"
)

type gantryTCPServer struct {
	LocalAddr string
	LocalPort string
	conn      *net.TCPListener
}

type gantryUDPServer struct {
	LocalAddr string
	LocalPort string
	conn      *net.UDPConn
}

// this is the module responsible for setting up a communication channel (TCP or UDP)
// where the data (protobuf, or JSON) can be exchanged

func NewGantryTCPServer(ip, port string) *gantryTCPServer {
	return &gantryTCPServer{
		LocalAddr: ip,
		LocalPort: port,
	}
}

func NewGantryUDPServer(ip, port string) *gantryUDPServer {
	return &gantryUDPServer{
		LocalAddr: ip,
		LocalPort: port,
	}
}

func (s *gantryTCPServer) StartTCP() {

	addr, err := net.ResolveTCPAddr("tcp", s.LocalAddr+":"+s.LocalPort)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Starting to listen on", s.LocalAddr, ":", s.LocalPort)
	ln, err := net.ListenTCP("tcp", addr)

	// assign the conn to stop the server
	s.conn = ln

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("OK: Listening on", s.LocalAddr, ":", s.LocalPort)

	for {
		if conn, err := ln.Accept(); err == nil {
			fmt.Println("Got a new connection request")
			go handleTCPConnection(conn)
		}
	}

}

func (s *gantryTCPServer) Stop() error {
	return s.conn.Close()
}

func (s *gantryUDPServer) Stop() error {
	return s.conn.Close()
}

func StartUDP(ip string, port string) {

	addr, err := net.ResolveUDPAddr("udp4", ip+":"+port)

	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatalln(err)
	}

	go handleUDPConnection(conn)

}

// Handles incoming requests.
func handleTCPConnection(conn net.Conn) {

	defer conn.Close()
	reader := bufio.NewReader(conn)

	data, err := ioutil.ReadAll(reader)

	fmt.Println("Reading connection data")

	if err != nil {
		log.Errorln(err)
		return
	}

	envelope := new(proto.Envelope)

	fmt.Println(protobuf.Unmarshal(data, envelope))

}

func handleUDPConnection(conn *net.UDPConn) {
	defer conn.Close()

	// set the buffer size
	if err := conn.SetReadBuffer(512); err != nil {
		log.Errorln(err)
	}

	buffer := make([]byte, 1500)

	if totalBytes, _, err := conn.ReadFromUDP(buffer); err != nil {
		log.Errorln(err)
	} else {

		envelope := new(proto.Envelope)
		protobuf.Unmarshal(buffer[:totalBytes], envelope)
	}

}
