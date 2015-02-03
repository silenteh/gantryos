package coms

import (
	protobuf "github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"

	"bufio"
	"io/ioutil"
	"net"
)

type gantryTCPServer struct {
	RemoteAddr string
	RemotePort string
	conn       *net.TCPListener
}

type gantryUDPServer struct {
	RemoteAddr string
	RemotePort string
	conn       *net.UDPConn
}

// this is the module responsible for setting up a communication channel (TCP or UDP)
// where the data (protobuf, or JSON) can be exchanged

func NewGantryTCPServer(ip, port string) *gantryTCPServer {
	return &gantryTCPServer{
		RemoteAddr: ip,
		RemotePort: port,
	}
}

func NewGantryUDPServer(ip, port string) *gantryUDPServer {
	return &gantryUDPServer{
		RemoteAddr: ip,
		RemotePort: port,
	}
}

func (s *gantryTCPServer) StartTCP() {

	addr, err := net.ResolveTCPAddr("ipv4", s.RemoteAddr+":"+s.RemotePort)

	if err != nil {
		log.Fatalln(err)
	}

	ln, err := net.ListenTCP("ipv4", addr)

	// assign the conn to stop the server
	s.conn = ln

	if err != nil {
		log.Fatalln(err)
	}
	for {
		if conn, err := ln.Accept(); err == nil {
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

	if err != nil {
		log.Errorln(err)
		return
	}

	envelope := new(proto.Envepope)

	protobuf.Unmarshal(data, envelope)
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

		envelope := new(proto.Envepope)
		protobuf.Unmarshal(buffer[:totalBytes], envelope)
	}

}
