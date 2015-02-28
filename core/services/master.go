package services

import (
	"errors"
	"fmt"
	protobuf "github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/config"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/core/resources"
	"github.com/silenteh/gantryos/core/state"
	"github.com/silenteh/gantryos/models"
	"github.com/silenteh/gantryos/utils"
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

func StartMaster(masterIp, masterPort string, readerChannel chan *proto.Envelope, writerChannel chan *proto.Envelope, stateDB state.StateDB) masterServer {

	// cerate a new master instance
	ms := newMaster(masterIp, masterPort, readerChannel, writerChannel)

	// init the TCP server
	ms.initTcpServer()

	// start the listener to detect the client calls
	ms.startMasterListener()

	// start writer
	ms.startMasterWriter()

	return ms

}

func (m *masterServer) startMasterWriter() {

	go writerMasterLoop(m)
}

func writerMasterLoop(m *masterServer) {

	for {
		envelope := <-m.writerChannel
		// means the channel got closed
		if envelope == nil {
			break
		}

		// switch between the possible messages the master can send to the slaves
		switch {
		case envelope.GetRunTask() != nil:
			// execute the task to the specific slave
			if err := write(envelope); err != nil {
				log.Errorln(err)
				fmt.Println(err)
			}

			continue
		case envelope.GetKillTask() != nil:
			// execute the task to the specific slave
			write(envelope)
			continue
		//case envelope.GetTaskStatus() != nil:
		//	// ask for the status of a specific task
		//	write(envelope)
		//	continue
		case envelope.GetReconcileTasks() != nil:
			// ask TO ALL SLAVES for the reconciliation of the tasks
			if envelope.GetReconcileTasks().GetSlaveId() == "" {
				writeAll(envelope)
			} else {
				// send it to a specific slave
				envelope.DestinationId = envelope.GetReconcileTasks().SlaveId
				write(envelope)
			}

			continue
		default:
			log.Errorln("Unkown envelope queued on master")
		}
	}
}

func write(envelope *proto.Envelope) error {

	slaveId := envelope.GetDestinationId()
	fmt.Println("Envelope slave Id", slaveId)

	if slaveId == "" {
		return errors.New("Unknown destination id, master skipping handling of data write")
	}

	conn := slaveConnections[slaveId]

	if conn != nil {
		data, err := protobuf.Marshal(envelope)
		if err != nil {
			return err
		}

		dataSize := len(data)

		data = append(utils.IntToBytes(dataSize), data...)

		_, err = conn.Write(data)
		return err
	} else {
		return errors.New("Unknown slave connection.")
	}

}

func writeAll(envelope *proto.Envelope) {

	for _, conn := range slaveConnections {
		if conn != nil {
			data, err := protobuf.Marshal(envelope)
			if err != nil {
				continue
			}

			dataSize := len(data)

			data = append(utils.IntToBytes(dataSize), data...)

			conn.Write(data)
		}
	}
}

func (m *masterServer) StopMaster() {

	if m.tcpServer != nil {
		m.tcpServer.Stop()
	}
	log.Infoln("GantryOS master stopped.")

}
