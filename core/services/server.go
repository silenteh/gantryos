package services

import (
	"bufio"
	"fmt"
	protobuf "github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
	"io"
	"net"
)

type gantryTCPServer struct {
	LocalAddr       string
	LocalPort       string
	envelopeChannel chan *proto.Envelope // from this channel we can read all the data coming from the clients
	conn            *net.TCPListener
}

type gantryUDPServer struct {
	LocalAddr string
	LocalPort string
	conn      *net.UDPConn
}

// this is the module responsible for setting up a communication channel (TCP or UDP)
// where the data (protobuf, or JSON) can be exchanged

func newGantryTCPServer(ip, port string, dataChannel chan *proto.Envelope) *gantryTCPServer {
	return &gantryTCPServer{
		LocalAddr:       ip,
		LocalPort:       port,
		envelopeChannel: dataChannel,
	}
}

func newGantryUDPServer(ip, port string) *gantryUDPServer {
	return &gantryUDPServer{
		LocalAddr: ip,
		LocalPort: port,
	}
}

func (s *gantryTCPServer) StartTCP() {

	// the for loop blocks the current thread therefore starts a differen one
	go func(server *gantryTCPServer) {

		addr, err := net.ResolveTCPAddr("tcp", s.LocalAddr+":"+s.LocalPort)

		if err != nil {
			log.Fatalln(err)
		}

		ln, err := net.ListenTCP("tcp", addr)

		// assign the conn to stop the server
		s.conn = ln

		if err != nil {
			log.Fatalln(err)
		}

		log.Infoln("GantryOS Master is listening on", s.LocalAddr, ":", s.LocalPort)

		for {
			if conn, err := ln.AcceptTCP(); err == nil {
				go handleTCPConnection(conn, s.envelopeChannel)
			}
		}

	}(s)

}

func (s *gantryTCPServer) Stop() error {
	return s.conn.Close()
}

func (s *gantryUDPServer) Stop() error {
	return s.conn.Close()
}

// Handles incoming requests.
func handleTCPConnection(conn *net.TCPConn, dataChannel chan *proto.Envelope) {

	//fmt.Println("Server receiving data")
	var err error

	defer conn.Close()
	reader := bufio.NewReader(conn)
	//writer := bufio.NewWriter(&buffer)

	sizeByte, err := reader.ReadByte()
	totalSize := int(sizeByte)
	//fmt.Println("Total size:", totalSize)

	//fmt.Println("Server reading data - 1")

	buffer := make([]byte, totalSize)
	totalRead, err := io.ReadFull(reader, buffer)
	if err != nil || totalRead != totalSize {
		fmt.Println("Error reading data")
	}
	//_, err = buffer.Write(data)

	//data, err := conn(reader.)
	//fmt.Println("Server reading data - 2")
	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Errorln(err)
	// 	return
	// }

	envelope := new(proto.Envelope)

	err = protobuf.Unmarshal(buffer, envelope)
	if err != nil {
		fmt.Println("Failed to parse client sent data")
		log.Errorln("Failed to parse client sent data:", err)
		return
	}

	fmt.Println("Server: writing data to channel")
	// it's a buffered channel, therefore this method does not block (unless the channel is full)
	// TODO: notify when the channel is full
	dataChannel <- envelope

}

//==============================================================================================================

func startUDP(ip string, port string) {

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
