// +build integration

package vswitch

import (
	"fmt"
	"testing"
	//"time"
)

var ovsdbTestServer string = "192.168.1.107"

//var ovsdbTestServer string = "192.168.59.104"

func TestGetBridgeUUID(t *testing.T) {

	fmt.Println("")

	testBridge := "br1"
	additionalPort := testBridge + "_1"

	client, err := newOVSDBClient(ovsdbTestServer, "6633")
	if err != nil {
		t.Error(err)
	}

	defer client.Close()

	// add a bridge for testing
	if _, err := client.AddBridge(testBridge, false); err != nil {
		t.Error(err)
	}

	bridgeUUID, err := getBridgeUUID(testBridge, client)
	if err != nil {
		t.Error(err)
	} else {
		addPort(additionalPort, bridgeUUID, "test_vpc", 5, VInterface{}, client)

		// this should throw an error !
		if _, _, err := addPort(additionalPort, bridgeUUID, "test_vpc", 5, VInterface{}, client); err == nil {
			t.Error("Error not thrown by the addPort ops !")
		}
	}

	if _, err := getPortUUID(testBridge, client); err != nil {
		t.Error(err)
	}

	if _, err := getInterfaceUUID(testBridge, client); err != nil {
		t.Error(err)
	}

	vswitch, err := getAllBridgePorts(bridgeUUID, client.GetRootUUID(), client)
	if err != nil {
		t.Error(err)
	}

	//fmt.Println(vswitch.toJson())

	ports := vswitch.VPCs[5].Ports

	if len(ports) == 0 {
		t.Error("There must be at least 1 ports !")
	}

	for _, port := range ports {
		ifaces, err := getAllPortInterfaces(port.Id, client)
		if len(ifaces) == 0 || err != nil {
			t.Error("Could not retrieve the interfaces from the port")
		}
	}

	if err := client.DeleteBridge(testBridge); err != nil {
		t.Error(err)
	}
}

func TestGetVPort(t *testing.T) {
	testBridge := "br4"
	additionalInterface := testBridge + "_additional"

	client, err := newOVSDBClient(ovsdbTestServer, "6633")
	if err != nil {
		t.Error(err)
	}

	defer client.Close()

	// add a bridge for testing
	if _, err := client.AddBridge(testBridge, false); err != nil {
		t.Error(err)
	}

	// this should throw the error !
	if _, err := client.AddBridge(testBridge, false); err == nil {
		t.Error("Error not thrown by the vswitch manager AddBridge method")
	}

	_, err = getBridgeUUID(testBridge, client)
	if err != nil {
		t.Error(err)
	}

	portUUID, err := getPortUUID(testBridge, client)
	if err != nil {
		t.Error(err)
	}

	if _, err := getInterfaceUUID(testBridge, client); err != nil {
		t.Error(err)
	}

	if _, err := addInterface(additionalInterface, portUUID, "", client); err != nil {
		t.Error(err)
	}

	// this should throw the error !
	if _, err := addInterface(additionalInterface, portUUID, "", client); err == nil {
		t.Error("Error not thrown by the addInterface")
	}

	port, vpcName, vlan, err := getVPort(portUUID, client)
	if err != nil {
		t.Error(err)
	}

	if vpcName != "default" {
		t.Error("getVPort should return the default VPC in this case")
	}

	if vlan != 0 {
		t.Error("getVPort should return the default VLAN(0) in this case")
	}

	if port.Id == "" {
		t.Error("Port ID cannot be empty")
	}

	if port.Name == "" {
		t.Error("Port Name cannot be empty")
	}

	if port.Name != testBridge {
		t.Error("Port Name cannot be different from " + testBridge)
	}

	if len(port.Interfaces) < 1 {
		t.Error("Port interface could not be loaded")
	}

	if port.Interfaces[testBridge].Name == "" {
		t.Error("Port interface is basically empty and has not value!")
	}

	if port.Interfaces[testBridge].Id == "" {
		t.Error("Port interface has no ID")
	}

	if port.Interfaces[additionalInterface].Name == "" {
		t.Error("Port additional interface is basically empty and has not value!")
	}

	if port.Interfaces[additionalInterface].Id == "" {
		t.Error("Port additional interface has no ID")
	}

	if err := client.DeleteBridge(testBridge); err != nil {
		t.Error(err)
	}

}
