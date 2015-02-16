package services

import (
	"fmt"
	"github.com/silenteh/gantryos/core/proto"
	mock "github.com/silenteh/gantryos/utils/testing"
	"testing"
)

func TestInitTcpClient(t *testing.T) {

	// get a master tcp
	masterReaderChannel := make(chan *proto.Envelope)
	masterWriterChannel := make(chan *proto.Envelope)

	// create the master
	m := newMaster("127.0.0.1", "6066", masterReaderChannel, masterWriterChannel)
	if m.master == nil {
		t.Error("Error creating the maste server")
	}

	// listen for TCP connections
	m.initTcpServer()

	// start the master writer
	m.startMasterWriter()

	if m.tcpServer == nil {
		t.Error("Error the maste cannot listen for TCP connections")
	}
	// ====================================================================

	// creatwe the slave
	slaveWriterChannel := make(chan *proto.Envelope, 1)
	slaveReaderChannel := make(chan *proto.Envelope, 1)

	// instanciate the slave
	slave := newSlave("127.0.0.1", "6066", slaveReaderChannel, slaveWriterChannel)

	// init the slave TCP
	slave.initTcpClient()

	// init the slave reader
	slave.startSlaveReader()

	// init the slave writer
	slave.startSlaveWriter()

	// mock a task status
	taskStatus := mock.MakeTaskStatus()
	taskStatus.Id = "123456789"
	//fmt.Printf("%s\n", taskStatus)

	// queue the taskstatus for writing
	slave.taskStateChange(taskStatus)

	// try to get the message from the master reader channel
	eventEnvelope := <-masterReaderChannel

	if eventEnvelope == nil {
		t.Error("Slave cannot write to the master - nil envelope")
	}

	if eventEnvelope.GetTaskStatusMessage() == nil {
		t.Error("Slave cannot write to the master")
	}

	if eventEnvelope.GetTaskStatusMessage().GetTaskStatus().GetGantryTaskId() != "123456789" {
		t.Error("Slave cannot write to the master - wrong task id")
	}

	// Now the master asks to run a task to the slave

	taskInfo := mock.MakeGolangHelloTask()
	m.taskRequest(taskInfo)

	// make sure the data is in the channel !
	taskRequest := <-slaveReaderChannel

	fmt.Printf("%s", taskRequest)

	//===========================================================================================
	// stopping the master
	m.StopMaster()
	slave.StopSlave()

	fmt.Println("- Slave writes to master and viceversa: OK")

}
