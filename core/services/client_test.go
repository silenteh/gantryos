package services

import (
	"fmt"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/models"
	"strconv"
	"testing"
	"time"
)

var master = models.NewMaster("1", "127.0.0.1", "localhost-master", 6060)
var slave = models.NewSlave("1", "127.0.0.1", "localhost-slave", 6061, true, false)

var masterChannel chan *proto.Envelope

func TestConnect(t *testing.T) {

	masterChannel = make(chan *proto.Envelope, 1024)

	tcpServer := newGantryTCPServer(master.Ip, strconv.Itoa(master.Port), masterChannel, nil)
	tcpServer.StartTCP()
	fmt.Println("Server started")

	time.Sleep(2 * time.Second)

	tcpClient := newGantryTCPClient(master.Ip, strconv.Itoa(master.Port))

	err := tcpClient.Connect()
	if err != nil {
		tcpServer.Stop()
		t.Fatal(err)
	}

	fmt.Println("Client connected")

	e := slave.NewHeartBeat() //models.NewHeartBeat(slave)

	for i := 0; i < 1024; i++ {
		err = tcpClient.Write(e)
		if err != nil {
			fmt.Println(i)
			t.Error(err)
		}
		//time.Sleep(1 * time.Microsecond)
	}

	data := <-masterChannel

	fmt.Println("Client received data from the channel")

	switch {
	case data.Heartbeat != nil:
		fmt.Println("Server Heartbeat proto message received")
		fmt.Println(data.Heartbeat.GetSlave().GetHostname())
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
