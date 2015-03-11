package networking

import (
	"fmt"
	"testing"
)

func TestNewOVSDBClient(t *testing.T) {

	client, err := NewOVSDBClient("192.168.1.117", "6633")
	if err != nil {
		//t.Error(err)
		return
	}

	defer client.Close()

	if ping, err := client.Echo(); err != nil {
		t.Error(err)
	} else {
		if ping[0] != "ping" {
			t.Error("Failed to receive the pong")
		}
	}

	dbs := client.ListDBs()
	fmt.Println(dbs)
	if len(dbs) == 0 {
		t.Error("Error getting the list of DBs")
	}

	if dbs[0] != "Open_vSwitch" {
		t.Error("Error Open_vSwitch not found")
	}

	_, err = client.GetSchema(dbs[0])
	if err != nil {
		t.Error(err)
	}

	root := client.GetRootUUID()
	if root == "" {
		t.Error("Got an empty root UUID !")
	}

	client.AddBridge("br0")

	client.Close()

	fmt.Println("- OVSDB Client")

}
