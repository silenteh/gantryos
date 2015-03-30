// +build integration

package ovsdb

import (
	//"fmt"
	"testing"
	"time"
)

func TestGetBridgeUUID(t *testing.T) {

	testBridge := "br1"

	client, err := NewOVSDBClient(ovsdbTestServer, "6633")
	if err != nil {
		t.Error(err)
	}

	defer client.Close()

	time.Sleep(5 * time.Second)

	// ad a bridge for testing
	if err := client.AddBridge(testBridge); err != nil {
		t.Error(err)
	}

	bridgeUUID, err := getBridgeUUID(testBridge, client)
	if err != nil {
		t.Error(err)
	} else {
		addPort("add_port", bridgeUUID, 5, client)
	}

	if _, err := getPortUUID(testBridge, client); err != nil {
		t.Error(err)
	}

	if _, err := getInterfaceUUID(testBridge, client); err != nil {
		t.Error(err)
	}

	ports, err := getAllBridgePorts(bridgeUUID, client)
	if err != nil {
		t.Error(err)
	}

	if len(ports) == 0 {
		t.Error("There must be at least 2 ports !")
	}

	for _, portUUID := range ports {
		ifaces, err := getAllPortInterfaces(portUUID, client)
		if len(ifaces) == 0 || err != nil {
			t.Error("Could not retrieve the interfaces from the port")
		}
	}

	//time.Sleep(5 * time.Second)

	if err := client.DeleteBridge(testBridge); err != nil {
		t.Error(err)
	}

	//time.Sleep(5 * time.Second)
}
