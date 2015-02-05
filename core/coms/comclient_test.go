package coms

import (
	"fmt"
	protobuf "github.com/gogo/protobuf/proto"
	//"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/models"
	"strconv"
	"testing"
	"time"
)

var master = models.NewMaster("1", "127.0.0.1", "localhost-master", 6060)
var slave = models.NewSlave("1", "127.0.0.1", "localhost-slave", 6061, true, nil)

func TestConnect(t *testing.T) {

	tcpServer := NewGantryTCPServer(master.Ip, strconv.Itoa(master.Port))
	go tcpServer.StartTCP()

	time.Sleep(3 * time.Second)

	tcpClient := NewGantryTCPClient(slave.Ip, strconv.Itoa(master.Port))

	conn, err := tcpClient.Connect()
	if err != nil {
		tcpServer.Stop()
		t.Fatal(err)
	}

	// taskInfo := proto.TaskInfo{}
	// taskInfo.TaskId = protobuf.String("1234")
	// taskInfo.TaskName = protobuf.String("TEST_TASK")

	e := models.NewEnvelope()
	e.MasterInfo = master.ToProtoBuf()

	data, err := protobuf.Marshal(e)
	if err != nil {
		tcpServer.Stop()
		fmt.Println(tcpClient.Disconnect())
		t.Fatal(err)
	}

	fmt.Println("Client connected")

	err = conn.WriteMessage(data)
	if err != nil {
		tcpServer.Stop()
		fmt.Println(tcpClient.Disconnect())
		t.Fatal(err)
	}
	fmt.Println("Message written")

	tcpServer.Stop()
	fmt.Println("Writing data to TCP server succeeded")
}

func startTcpTestServer() {

}
