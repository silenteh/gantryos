package services

import (
	"fmt"
	"github.com/silenteh/gantryos/core/proto"
	"testing"
)

func TestJoinMaster(t *testing.T) {

	writerChannel := make(chan *proto.Envelope, 1)

	// instanciate the slave
	slave := newSlave("127.0.0.1", "6050", nil, writerChannel)

	// queue a register slave message for writing
	slave.joinMaster()

	//get the message and check if the format is correct
	envelope := <-slave.writerChannel

	if envelope.RegisterSlave == nil {
		t.Error("Error sending the slave registration message")
	}

	if envelope.RegisterSlave.GetSlave().GetHostname() == "" {
		t.Error("Error sending the slave registration message - SlaveInfo not present")
	}

	close(writerChannel)

	fmt.Println("JoinMaster: OK")

}

func TestReRegisterMaster(t *testing.T) {

	writerChannel := make(chan *proto.Envelope, 1)

	// instanciate the slave
	slave := newSlave("127.0.0.1", "6050", nil, writerChannel)

	// queue a register slave message for writing
	slave.reRegisterMaster()

	//get the message and check if the format is correct
	envelope := <-slave.writerChannel

	if envelope.ReRegisterSlave == nil {
		t.Error("Error sending the slave RE-registration message")
	}

	if envelope.ReRegisterSlave.GetSlave().GetHostname() == "" {
		t.Error("Error sending the slave Re-registration message - SlaveInfo not present")
	}

	close(writerChannel)

	fmt.Println("ReRegisterMaster: OK")

}

func TestPingMaster(t *testing.T) {

	writerChannel := make(chan *proto.Envelope, 1)

	// instanciate the slave
	slave := newSlave("127.0.0.1", "6050", nil, writerChannel)

	// queue a register slave message for writing
	slave.pingMaster()

	//get the message and check if the format is correct
	envelope := <-slave.writerChannel

	if envelope.Heartbeat == nil {
		t.Error("Error sending the slave heartbeat message")
	}

	if envelope.Heartbeat.GetSlave().GetHostname() == "" {
		t.Error("Error sending the slave heartbeat message - SlaveInfo not present")
	}

	close(writerChannel)

	fmt.Println("Slave heartbeat: OK")

}
