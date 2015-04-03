// +build integration

package ovsdb

import (
	//"fmt"
	"testing"
	//"time"
)

//var ovsdbTestServer string = "192.168.1.107"

func TestNewVSwitch(t *testing.T) {

	//defaultBridge := "default"
	testBridge := "br2"

	manager, err := NewOVSDBClient(ovsdbTestServer, "6633")
	if err != nil {
		t.Error(err)
	}

	defer manager.Close()

	vswitch, err := InitVSwitch(manager)
	if err != nil {
		t.Error(err)
	}

	if vswitch.Id == "" {
		t.Error("New vswitch cannot have an empty Id")
	}

	if vswitch.Name == "" {
		t.Error("New vswitch cannot have an empty Name")
	}

	vpc := vswitch.AddVPC(testBridge, "192.168.10.0/24", 1)

	if _, err := vpc.AddPort(testBridge, vswitch.Id, vpc.Name, "", 1, manager); err != nil {
		t.Error(err)
	}

	if p2, err := vpc.AddPort(testBridge, vswitch.Id, vpc.Name, "123456789", 1, manager); err != nil {
		t.Error(err)
	} else {
		//fmt.Println(p2)
		if p2.Interfaces["br2_2"].ContainerId == "" {
			t.Error("ContainerId should not be empty !")
		}

		if p2.Interfaces["br2_2"].ContainerId != "123456789" {
			t.Error("ContainerId should be 123456789 !")
		}
	}

	if len(vpc.Ports) < 2 {
		t.Error("VPC should have at least 2 ports !")
	}

	//fmt.Println(vswitch.toJson())

	if port := vswitch.VPCs[testBridge].Ports[testBridge]; port.Id != "" {

		if len(port.Interfaces) < 1 {
			t.Error("Port should have at least 1 interface !")
		}
	}

	// check if the load part works before deleting

	loadedVswitch, err := InitVSwitch(manager)
	if err != nil {
		t.Error(err)
	}

	//fmt.Println(loadedVswitch.toJson())

	if loadedVswitch.Id == "" {
		t.Error("Loaded vswitch does not have a valid ID")
	}

	if len(loadedVswitch.VPCs) < 2 {
		t.Error("There should be at least 2 VPCs")
	}

	vpc0 := loadedVswitch.VPCs["0"]

	if len(vpc0.Ports) < 1 {
		t.Error("Default VPC has at least 1 port")
	}

	vpc1 := loadedVswitch.VPCs["1"]

	if len(vpc1.Ports) < 2 {
		t.Error("Additional VPC with vlan 1 has at least 2 ports")
	}

	// Check if we can load the containerId and Iface
	if vpc1.Ports["br2_2"].Interfaces["br2_2"].ContainerId == "" {
		t.Error("ContainerId should not be empty !")
	}

	if vpc1.Ports["br2_2"].Interfaces["br2_2"].ContainerId == "" {
		t.Error("ContainerId should be 123456789 !")
	}

	// ============================================================
	// check if the port really still exixts
	if portUUID, err := getPortUUID(testBridge, manager); err != nil {
		t.Error(err)
	} else {
		if portUUID != "" {
			t.Error("The deleted port cannot have a valid UUID !")
		}
	}

	// ====================================================================
	// delete the VPC
	if err := vpc.DeleteVPC(&vswitch, manager); err != nil {
		t.Error(err)
	}

	if _, ok := vswitch.VPCs["default"].Ports[testBridge]; ok {
		t.Error("Port cannot be still in memory on the vswitch VPCs map!")
	}

	if err := vswitch.DeleteVSwitch(); err != nil {
		t.Error(err)
	}

	if bridgeUUID, err := getBridgeUUID("br2", manager); err != nil {
		t.Error(err)
	} else {
		if bridgeUUID != "" {
			t.Error("The deleted bridge cannot have a valid UUID !")
		}
	}
}
