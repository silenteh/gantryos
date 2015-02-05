package coms

import (
	"fmt"
	//protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/models"
	"strconv"
	"testing"
	"time"
)

var master = models.NewMaster("1", "127.0.0.1", "localhost-master", 6060)
var slave = models.NewSlave("1", "127.0.0.1", "localhost-slave", 6061, true, nil)
var dataChannel = make(chan *proto.Envelope, 10)

func TestConnect(t *testing.T) {

	//gantryMasterService := masters.Start(master.Ip, strconv.Itoa(master.Port))

	tcpServer := NewGantryTCPServer(master.Ip, strconv.Itoa(master.Port), dataChannel)
	tcpServer.StartTCP()

	time.Sleep(3 * time.Second)

	tcpClient := NewGantryTCPClient(master.Ip, strconv.Itoa(master.Port))

	err := tcpClient.Connect()
	if err != nil {
		tcpServer.Stop()
		t.Fatal(err)
	}

	fmt.Println("Client connected")

	heartbeat := models.NewHeartBeat(slave)

	e := models.NewEnvelope()
	e.Heartbeat = heartbeat

	for i := 0; i < 100; i++ {
		err = tcpClient.Write(e)
		if err != nil {
			tcpServer.Stop()
			fmt.Println(tcpClient.Disconnect())
			t.Fatal(err)
		}
		fmt.Println("Write:", i)
	}

	fmt.Println("Message written")

	data := <-dataChannel

	fmt.Println("Test: received data from the channel")

	switch {
	case data.Heartbeat != nil:
		fmt.Println("HEARTBEAT RECEIVED !")
		break
	default:
		tcpClient.Disconnect()
		tcpServer.Stop()
		t.Fatal("Unknwon message received")
		break
	}

	tcpClient.Disconnect()
	tcpServer.Stop()
	fmt.Println("Writing data to TCP server succeeded")
}

func startTcpTestServer() {

}
