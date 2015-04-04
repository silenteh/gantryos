package services

import (
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/config"
	"github.com/silenteh/gantryos/core/networking/vswitch"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/core/resources"
	"github.com/silenteh/gantryos/core/state"
	"github.com/silenteh/gantryos/models"
	"time"
)

type slaveServer struct {
	slave         *models.Slave        // the slave model
	writerChannel chan *proto.Envelope // the writer channel to send info the to master
	readerChannel chan *proto.Envelope // the reader channels for communications from the master
	tcpClient     *gantryTCPClient     // the tcp client connection to the master
	masterIp      string               // the master IP
	masterPort    string               // the master port to send the messages to
	vswitch       vswitch.Vswitch      // the vswitch which contains the VPCs and other info
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
	slaveConfig := config.GantryOSSlaveId

	slaveInfo := models.NewSlave(slaveConfig.Id, ip, hostname, port, config.GantryOSConfig.Slave.Checkpoint, slaveConfig.Registered)

	// assign properties
	slave.masterIp = masterIp
	slave.masterPort = masterPort
	slave.readerChannel = readerChannel
	slave.writerChannel = writerChannel
	slave.slave = slaveInfo

	return slave

}

func (s *slaveServer) initVswitch() {
	// create the vswitch manager
	vswitchHost := "127.0.0.1"
	vswitchPort := "6633"
	if config.GantryOSConfig.Slave.VSwitchServer.Hostname != "" {
		vswitchHost = config.GantryOSConfig.Slave.VSwitchServer.Hostname
	}

	if config.GantryOSConfig.Slave.VSwitchServer.Port != "" {
		vswitchPort = config.GantryOSConfig.Slave.VSwitchServer.Port
	}

	vswitch, err := vswitch.InitVSwitch(vswitchHost, vswitchPort)
	if err != nil {
		log.Errorln(err)
	} else {
		s.vswitch = vswitch
	}

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

func StartSlave(masterIp, masterPort string, readerChannel chan *proto.Envelope, writerChannel chan *proto.Envelope, stateDB state.StateDB) slaveServer {

	// create a new slave client
	slave := newSlave(masterIp, masterPort, readerChannel, writerChannel)

	// init the VSwitch connection and default switch
	slave.initVswitch()

	// init the TCP connection with the amster
	slave.initTcpClient()

	// start the loop for writing to the TCP connection
	slave.startSlaveWriter()

	// start the loop for reading data coming from the master
	slave.startSlaveReader()

	// start the loop for processing the data read from the master
	slave.startSlaveListener(stateDB)

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
