package services

import (
	"bufio"
	protobuf "github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
	"io"
	"net"
)

var totalSlaveConnections = 1

type gantryTCPServer struct {
	LocalAddr     string
	LocalPort     string
	readerChannel chan *proto.Envelope // from this channel we can read all the data coming from the clients
	writerChannel chan *proto.Envelope // from this channel we can write all the data to the clients
	conn          *net.TCPListener
}

type gantryUDPServer struct {
	LocalAddr string
	LocalPort string
	conn      *net.UDPConn
}

// this is the module responsible for setting up a communication channel (TCP or UDP)
// where the data (protobuf, or JSON) can be exchanged

func newGantryTCPServer(ip, port string, readerChannel chan *proto.Envelope, writerChannel chan *proto.Envelope) *gantryTCPServer {
	return &gantryTCPServer{
		LocalAddr:     ip,
		LocalPort:     port,
		readerChannel: readerChannel,
		writerChannel: writerChannel,
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
				// we accept a maximum of 640000 concurrent connections
				// each slave creates 1 connection, therefore it should be enough for handling up to 64k slaves !
				if totalSlaveConnections < 640000 {
					log.Infoln("Amount of slave connections:", totalSlaveConnections)
					go handleTCPConnection(conn, s.readerChannel)
				} else {
					log.Errorln("Too many connections from slaves. Stopped accepting new connections.")
				}
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

	// new slave connection was successfully created
	totalSlaveConnections++

	// close the connection in case this exists
	defer conn.Close()

	for {

		reader := bufio.NewReader(conn)

		sizeByte, err := reader.ReadByte()
		totalSize := int(sizeByte)

		buffer := make([]byte, totalSize)
		totalRead, err := io.ReadFull(reader, buffer)
		if err != nil || totalRead != totalSize {
			log.Errorln("Error reading data from socket:", err)
			continue
		}

		envelope := new(proto.Envelope)

		err = protobuf.Unmarshal(buffer, envelope)
		if err != nil {
			log.Errorln("Failed to parse client sent data:", err)
			continue
		}
		// it's a buffered channel, therefore this method does not block (unless the channel is full)
		// TODO: notify when the channel is full
		dataChannel <- envelope
	}

	// at this point the connection will be closed therefore decrease the counter
	totalSlaveConnections--

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
