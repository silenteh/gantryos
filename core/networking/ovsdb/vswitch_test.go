// +build integration

package ovsdb

import (
	"fmt"
	"testing"
	"time"
)

var ovsdbTestServer string = "192.168.1.107"

func TestNewOVSDBClient(t *testing.T) {

	client, err := NewOVSDBClient(ovsdbTestServer, "6633")
	if err != nil {
		t.Error(err)
		fmt.Println("Client failed to instantiate or connect. All following requests will fail")
	}

	defer client.Close()

	time.Sleep(5 * time.Second)

	if ping, err := client.Echo(); err != nil {
		t.Error(err)
	} else {
		if ping[0] != "ping" {
			t.Error("Failed to receive the pong")
		}
	}

	dbs, err := client.ListDBs()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(dbs)
	if len(dbs) == 0 {
		t.Fatal("Error getting the list of DBs")
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

	//time.Sleep(3 * time.Second)

	if err := client.AddBridge("br0"); err != nil {
		t.Error(err)
	}

	// // set VLAN
	//time.Sleep(3 * time.Second)

	if err := client.DeleteBridge("br0"); err != nil {
		t.Error(err)
	}

	//time.Sleep(5 * time.Second)

	fmt.Println("- OVSDB Client")

	time.Sleep(5 * time.Second)

}
