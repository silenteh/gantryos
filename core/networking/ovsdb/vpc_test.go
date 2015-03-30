// +build integration

package ovsdb

import (
	"fmt"
	"testing"
	"time"
)

var ovsdbTestServer string = "192.168.1.107"

func TestNewVSwitch(t *testing.T) {

	testBridge := "br2"

	manager, err := NewOVSDBClient(ovsdbTestServer, "6633")
	if err != nil {
		t.Error(err)
	}

	defer manager.Close()

	vswitch, err := NewVSwitch(manager.GetRootUUID(), testBridge, false, manager)
	if err != nil {
		t.Error(err)
	}

	if vswitch.Id == "" {
		t.Error("New vswitch cannot have an empty Id")
	}

	if vswitch.Name == "" {
		t.Error("New vswitch cannot have an empty Name")
	}

	time.Sleep(1 * time.Second)

	vswitch.AddVPC("default", "192.168.2.0/24", 1)

	if err := vswitch.VPCs["default"].AddPort("br2", vswitch.Id, 1, vswitch.Manager); err != nil {
		t.Error(err)
	}

	if port := vswitch.VPCs["default"].Ports["br2"]; port.Id != "" {
		port.AddInterface("br2_additional", manager)
	}

	time.Sleep(1 * time.Second)

	fmt.Println("DONE")

}
