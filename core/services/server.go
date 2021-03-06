package services

import (
	"bufio"
	"fmt"
	protobuf "github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/utils"
	"io"
	"net"
	"sync"
)

// tracks down the connections that the slaves made to the server so we can execute a specific task on a specific slave
var slaveConnections = make(map[string]*net.TCPConn)
var counterMutex sync.Mutex
var totalSlaveConnections = 1

type gantryTCPServer struct {
	stopMutex     sync.Mutex
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
		s.stopMutex.Lock()
		s.conn = ln
		s.stopMutex.Unlock()

		if err != nil {
			log.Fatalln(err)
		}

		log.Infoln("GantryOS Master is listening on", s.LocalAddr, ":", s.LocalPort)

		for {
			if conn, err := ln.AcceptTCP(); err == nil {
				// we accept a maximum of 640000 concurrent connections
				// each slave creates 1 connection, therefore it should be enough for handling up to 64k slaves !
				counterMutex.Lock()
				if totalSlaveConnections < 640000 {
					log.Infoln("Amount of slave connections:", totalSlaveConnections)
					go handleTCPConnection(conn, s.readerChannel)
				} else {
					log.Errorln("Too many connections from slaves. Stopped accepting new connections.")
				}
				counterMutex.Unlock()
			}
		}

	}(s)

}

func (s *gantryTCPServer) Stop() error {
	s.stopMutex.Lock()
	defer s.stopMutex.Unlock()
	return s.conn.Close()
}

func (s *gantryUDPServer) Stop() error {
	return s.conn.Close()
}

// Handles incoming requests.
func handleTCPConnection(conn *net.TCPConn, dataChannel chan *proto.Envelope) {

	// new slave connection was successfully created
	counterMutex.Lock()
	totalSlaveConnections++
	counterMutex.Unlock()

	var slaveId string

	// close the connection in case this exists
	defer conn.Close()

	for {

		// new buffered reader
		reader := bufio.NewReader(conn)

		// get the lenght
		lenght := make([]byte, 4)
		_, err := reader.Read(lenght)
		totalSize := utils.BytesToInt(lenght)

		buffer := make([]byte, totalSize)
		_, err = io.ReadFull(reader, buffer) //io.ReadFull(reader, buffer)
		if err != nil && err != io.EOF {
			log.Errorln(err)
			break
		}

		envelope := new(proto.Envelope)

		err = protobuf.Unmarshal(buffer, envelope)
		if err != nil {
			log.Errorln("Failed to parse client sent data:", err)
			continue
		}

		// detect the slaveid
		if envelope.GetRegisterSlave() != nil {
			slaveId = envelope.GetRegisterSlave().GetSlave().GetId()
			fmt.Println("Adding slave to the index:", slaveId)
			slaveConnections[slaveId] = conn
		}

		// detect the slaveid
		if envelope.GetReRegisterSlave() != nil {
			slaveId = envelope.GetReRegisterSlave().GetSlave().GetId()
			fmt.Println("Adding slave to the index:", slaveId)
			slaveConnections[slaveId] = conn
		}

		// it's a buffered channel, therefore this method does not block (unless the channel is full)
		// TODO: notify when the channel is full
		dataChannel <- envelope
	}

	// remove the connection
	if slaveId != "" {
		fmt.Println("Removing slave from index:", slaveId)
		slaveConnections[slaveId] = nil
	}

	// at this point the connection will be closed therefore decrease the counter
	counterMutex.Lock()
	totalSlaveConnections--
	counterMutex.Unlock()

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
