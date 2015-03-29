// +build integration

package ovsdb

import (
	"fmt"
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

	if uuid, err := getBridgeUUID(testBridge, client); err != nil {
		t.Error(err)
	} else {
		fmt.Println(uuid)
	}

	if uuid, err := getPortUUID(testBridge, client); err != nil {
		t.Error(err)
	} else {
		fmt.Println(uuid)
	}

	if uuid, err := getInterfaceUUID(testBridge, client); err != nil {
		t.Error(err)
	} else {
		fmt.Println(uuid)
	}

	//time.Sleep(5 * time.Second)

	if err := client.DeleteBridge(testBridge); err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
}
