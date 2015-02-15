package services

import (
	"fmt"
	"github.com/silenteh/gantryos/core/proto"
	"testing"
)

func TestInitTcpServer(t *testing.T) {
	readerChannel := make(chan *proto.Envelope)
	writerChannel := make(chan *proto.Envelope)

	// create the master
	m := newMaster("127.0.0.1", "6065", readerChannel, writerChannel)
	if m.master == nil {
		t.Error("Error creating the maste server")
	}

	// listen for TCP connections
	m.initTcpServer()

	if m.tcpServer == nil {
		t.Error("Error the maste cannot listen for TCP connections")
	}

	// stopping the master
	m.StopMaster()

	fmt.Println("Master server initialization: OK")

}

func TestNewMaster(t *testing.T) {
	masterReaderChannel := make(chan *proto.Envelope)
	//writerChannel := make(chan *proto.Envelope, 64)

	//slaveReaderChannel := make(chan *proto.Envelope, 64)
	slaveWriterChannel := make(chan *proto.Envelope)

	// create the master
	m := newMaster("127.0.0.1", "6070", masterReaderChannel, nil)
	if m.master == nil {
		t.Error("Error creating the maste server")
	}

	// master start the reader
	m.initTcpServer()

	if m.tcpServer == nil {
		t.Fatalf("Error the maste cannot listen for TCP connections")
	}
	// ==========================================================================
	// setup a client and send to the master the messages
	slave := newSlave("127.0.0.1", "6070", nil, slaveWriterChannel)
	// start the tcp connection with the master
	slave.initTcpClient()

	// start the writer loop
	slave.startSlaveWriter()

	// send an heartbeat
	slave.pingMaster()
	envelope := <-masterReaderChannel
	if envelope.GetHeartbeat() == nil {
		t.Fatal("Error getting the heartbeat message from the slave")
	}
	if envelope.GetHeartbeat().GetSlave().GetHostname() == "" {
		t.Fatal("Error getting the heartbeat message from the slave")
	}

	// send a register with the master
	slave.joinMaster()
	envelope = <-masterReaderChannel
	if envelope.GetRegisterSlave() == nil {
		t.Error("Error getting the register message from the slave")
	}
	if envelope.GetRegisterSlave().GetSlave().GetHostname() == "" {
		t.Error("Error getting the register message from the slave")
	}

	slave.reRegisterMaster()
	envelope = <-masterReaderChannel
	if envelope.GetReRegisterSlave() == nil {
		t.Error("Error getting the RE-register message from the slave")
	}
	if envelope.GetReRegisterSlave().GetSlave().GetHostname() == "" {
		t.Error("Error getting the RE-register message from the slave")
	}
	// --------------------------------------------------------------------------
	// Get the messages from the master readerChannel we should have them all and in the same order

	// ==========================================================================

	// stop the slave

	//close(slaveReaderChannel)
	slave.tcpClient.Disconnect()
	//close(slaveWriterChannel)
	m.tcpServer.Stop()
	// stopping the master
	//close(masterReaderChannel)
	//close(writerChannel)

	fmt.Println("Master Server messaging: OK")

}
