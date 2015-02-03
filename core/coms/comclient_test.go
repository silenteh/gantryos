package coms

import (
	"fmt"
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
	"testing"
)

func TestConnect(t *testing.T) {

	tcpServer := NewGantryTCPServer("127.0.0.1", "6060")
	tcpClient := NewGantryTCPClient("127.0.0.1", "6060")

	conn, err := tcpClient.Connect()
	if err != nil {
		tcpServer.Stop()
		t.Fatal(err)
	}

	taskInfo := proto.TaskInfo{}
	taskInfo.TaskId = protobuf.String("1234")
	taskInfo.TaskName = protobuf.String("TEST_TASK")

	data, err := protobuf.Marshal(taskInfo)
	if err != nil {
		tcpServer.Stop()
		t.Fatal(err)
	}

	err = conn.WriteMessage(data)
	if err != nil {
		tcpServer.Stop()
		t.Fatal(err)
	}

	tcpServer.Stop()
	fmt.Println("Writing data to TCP server succeeded")
}

func startTcpTestServer() {

}
