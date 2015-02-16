package services

import (
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/config"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/core/resources"
	"github.com/silenteh/gantryos/models"
	"time"
)

type slaveServer struct {
	slave         *models.Slave
	writerChannel chan *proto.Envelope
	readerChannel chan *proto.Envelope
	tcpClient     *gantryTCPClient
	masterIp      string
	masterPort    string
}

func newSlave(masterIp, masterPort string, readerChannel chan *proto.Envelope, writerChannel chan *proto.Envelope) slaveServer {

	slave := slaveServer{}

	// Port
	port := 6051
	if config.GantryOSConfig.Slave.Port != 0 {
		port = config.GantryOSConfig.Slave.Port
	}
	// ==============================================

	// IP
	ip := "127.0.0.1"
	if config.GantryOSConfig.Slave.IP != "" {
		ip = config.GantryOSConfig.Slave.IP
	}
	// ==============================================

	// Hostname
	hostname := resources.GetHostname()
	// ==============================================

	// Slave ID
	slaveId := config.GantryOSSlaveId

	slaveInfo := models.NewSlave(slaveId.Id, ip, hostname, port, config.GantryOSConfig.Slave.Checkpoint, slaveId.Registered)

	// assign properties
	slave.masterIp = masterIp
	slave.masterPort = masterPort
	slave.readerChannel = readerChannel
	slave.writerChannel = writerChannel
	slave.slave = slaveInfo

	return slave

}

func (s *slaveServer) initTcpClient() {

	slaveInstance := newGantryTCPClient(s.masterIp, s.masterPort)
	err := slaveInstance.Connect()
	if err != nil {
		log.Fatalln("Cannot connect to master", s.masterIp, "on port", s.masterPort, " => ", err)
	}

	s.tcpClient = slaveInstance

}

func (s *slaveServer) startSlaveWriter() {

	go writerSlaveLoop(s)
}

func (s *slaveServer) startSlaveReader() {
	go s.tcpClient.read(s.readerChannel)
}

func writerSlaveLoop(s *slaveServer) {

	for {
		envelope := <-s.writerChannel
		// means the channel got closed
		if envelope == nil {
			break
		}
		if err := s.tcpClient.Write(envelope); err != nil {
			log.Errorln(err)
			// re-queue
			s.writerChannel <- envelope
			// disconnect and ignore the error
			s.tcpClient.Disconnect()
			//reconnect
			s.tcpClient.Connect()

			continue
		}
	}
}

func (s *slaveServer) startHeartBeats() {
	go func(slave *slaveServer) {
		for {
			s.pingMaster()
			time.Sleep(15 * time.Second)
		}
	}(s)
}

// Public methods

func StartSlave(masterIp, masterPort string, readerChannel chan *proto.Envelope, writerChannel chan *proto.Envelope) slaveServer {

	// create a new slave client
	slave := newSlave(masterIp, masterPort, readerChannel, writerChannel)

	// init the TCP connection with the amster
	slave.initTcpClient()

	// start the loop for writing to the TCP connection
	slave.startSlaveWriter()

	// start the loop for reading data coming from the master
	slave.startSlaveReader()

	// start the loop for processing the data read from the master
	slave.startSlaveListener()

	// start to send the heartbeats to the master
	slave.startHeartBeats()

	return slave

}

func (s *slaveServer) StopSlave() {

	if s.tcpClient != nil {
		s.tcpClient.Disconnect()
	}
	log.Infoln("GantryOS slave stopped.")

}
