package services

import (
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/config"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/core/resources"
	"github.com/silenteh/gantryos/models"
	"strconv"
)

type masterServer struct {
	master        *models.Master
	writerChannel chan *proto.Envelope
	readerChannel chan *proto.Envelope
	tcpServer     *gantryTCPServer
}

func newMaster(masterIp, masterPort string, readerChannel chan *proto.Envelope, writerChannel chan *proto.Envelope) masterServer {

	// Port
	port := 6050
	if config.GantryOSConfig.Master.Port != 0 {
		port = config.GantryOSConfig.Master.Port
	}
	if masterPort != "" {
		if newPort, err := strconv.Atoi(masterPort); err == nil {
			port = newPort
		}
	}
	// ==============================================

	// IP
	ip := "127.0.0.1"
	if config.GantryOSConfig.Master.IP != "" {
		ip = config.GantryOSConfig.Master.IP
	}
	if masterIp != "" {
		ip = masterIp
	}
	// ==============================================

	// Hostname
	hostname := resources.GetHostname()
	// ==============================================

	// Slave ID
	masterId := config.GantryOSMasterId

	master := models.NewMaster(masterId.Id, ip, hostname, port)

	ms := masterServer{}
	ms.master = master
	ms.writerChannel = writerChannel
	ms.readerChannel = readerChannel

	return ms

}

func (m *masterServer) initTcpServer() {
	masterInstance := newGantryTCPServer(m.master.Ip, strconv.Itoa(m.master.Port), m.readerChannel, m.writerChannel)
	m.tcpServer = masterInstance
	m.tcpServer.StartTCP()
}

func StartMaster(masterIp, masterPort string, readerChannel chan *proto.Envelope, writerChannel chan *proto.Envelope) masterServer {

	// cerate a new master instance
	ms := newMaster(masterIp, masterPort, readerChannel, writerChannel)

	// init the TCP server
	ms.initTcpServer()

	// start the listener to detect the client calls
	ms.startListener()

	return ms

}

func (m *masterServer) StopMaster() {

	if m.tcpServer != nil {
		m.tcpServer.Stop()
	}
	log.Infoln("GantryOS master stopped.")

}
