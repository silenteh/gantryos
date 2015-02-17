package services

import (
	"fmt"
	"github.com/silenteh/gantryos/core/proto"
	mock "github.com/silenteh/gantryos/utils/testing"
	"testing"
)

func TestInitTcpClient(t *testing.T) {

	// get a master tcp
	masterReaderChannel := make(chan *proto.Envelope, 2)
	masterWriterChannel := make(chan *proto.Envelope, 2)

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
	slaveWriterChannel := make(chan *proto.Envelope, 2)
	slaveReaderChannel := make(chan *proto.Envelope, 2)

	// instanciate the slave
	slave := newSlave("127.0.0.1", "6066", slaveReaderChannel, slaveWriterChannel)
	slave.slave.Id = "test_slave_id_123456789"

	// init the slave TCP
	slave.initTcpClient()

	// init the slave reader
	slave.startSlaveReader()

	// init the slave writer
	slave.startSlaveWriter()

	// this is needed to register the slave with the master
	slave.joinMaster()

	// read the register message first and discard it
	eventEnvelope := <-masterReaderChannel

	// mock a task status
	taskStatus := mock.MakeTaskStatus()
	taskStatus.Id = "123456789"
	//fmt.Printf("%s\n", taskStatus)

	// queue the taskstatus for writing
	slave.taskStateChange(taskStatus)

	// try to get the message from the master reader channel

	// read the task state change message
	eventEnvelope = <-masterReaderChannel

	//fmt.Printf("\n\n%s\n\n", eventEnvelope)

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

	// // make sure the data is in the channel !
	taskRequestEnvelope := <-slaveReaderChannel

	if taskRequestEnvelope.GetSenderId() != m.master.Id {
		t.Error("Error matching the master and the sender id")
	}

	taskRequest := taskRequestEnvelope.GetRunTask()

	if taskRequest == nil {
		t.Error("Master cannot send data to the slave")
	}

	//fmt.Printf("\n\n%s\n\n", taskRequest)

	//===========================================================================================
	// stopping the master
	m.StopMaster()
	slave.StopSlave()

	fmt.Println("- Slave writes to master and viceversa: OK")

}
