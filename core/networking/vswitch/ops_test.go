// +build integration

package vswitch

import (
	_ "fmt"
	"strconv"
	"testing"
	//"time"
)

// var ovsdbHost string = "192.168.59.104"
var ovsdbHost string = "192.168.1.113"
var ovsdbPort string = "6633"

// this function is supposed to add a bridge + port + interface
func TestAddFullBridge(t *testing.T) {

	bridgeName := "br0_test"
	manager, err := newOVSDBClient(ovsdbHost, ovsdbPort)
	if err != nil {
		t.Fatal(err)
	}

	defer manager.Close()

	// add the bridge
	bridgeUUID, err := addFullBridge(bridgeName, true, manager)
	if err != nil {
		t.Error(err)
	}
	defer deleteBridge(bridgeUUID, manager)
	//fmt.Println(bridgeUUID)

	// check the UUID
	if bridgeUUID == "" {
		t.Error("Bridge UUID cannot be empty")
	}

	// try to look for the bridge
	storedUUID, err := getBridgeUUID(bridgeName, manager)
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(storedUUID)

	// compare the 2 UUIDs
	if bridgeUUID != storedUUID {
		t.Error("The 2 bridge UUIDs are different !")
	}

	// VPC CHECK
	// check if the we can create a VPC from this bridge
	vpc, err := getVPC(bridgeUUID, manager)
	if err != nil {
		t.Error(err)
	}

	if len(vpc.Ports) != 1 {
		t.Error("Wrong number of ports on the loaded VPC: expected 1 got " + strconv.Itoa(len(vpc.Ports)))
	}
	// =====================================================================

	// load the port
	vport, vlan, err := getVPort(vpc.Ports[bridgeName].Id, manager)
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(vport.Id)

	// check that the Port tag and the VPC vlan are the same
	if vlan != vpc.VLan {
		t.Error("The port effective VLAN and the VPC VLAN ids are different !")
	}

	// check that the UUIDS correspond
	if vport.Id != vpc.Ports[vport.Name].Id {
		t.Error("The loaded port has a different UUID than the VPC port")
	}

	if vport.Name == "" {
		t.Error("The loaded port has an empty name")
	}

	if vport.Tag != vlan && vport.Tag != vpc.VLan {
		t.Error("The port VLAN TAG and the VPC VLAN values are different !")
	}

	if !vport.hasInterfaces() {
		t.Error("The port should have exactly 1 interface, via hasInterfaces method")
	}

	if vport.totalInterfaces() != 1 {
		t.Error("The port should have exactly 1 interface, via totalInterfaces method")
	}

	// =====================================================================

	// INTERFACE CHECK
	// check that the vport has 1 interface
	if len(vport.Interfaces) != 1 {
		t.Error("The port should have exactly 1 interface")
	}

	// load the interface via ops
	vint, err := getVInterface(vport.Interfaces[0].Id, manager)
	if err != nil {
		t.Error(err)
	}

	// check that the UUIDS correspond
	if vint.Id != vpc.Ports[vport.Name].Interfaces[0].Id {
		t.Error("The loaded interface has a different UUID than the VPC port interface")
	}

	if vint.Name == "" {
		t.Error("The loaded interface has an empty name")
	}

	if vint.Type != vpc.Ports[vport.Name].Interfaces[0].Type {
		t.Error("The interface type mismatch between the VPC port one and the loaded one")
	}

	// =====================================================================

}

// here we test the adding of a port
func TestAddPort(t *testing.T) {

	bridgeName := "br1_test"
	additionalPortName := "eth0"
	vlan := 0
	containerPID := ""
	taskID := ""

	manager, err := newOVSDBClient(ovsdbHost, ovsdbPort)
	if err != nil {
		t.Fatal(err)
	}

	defer manager.Close()

	// add the bridge
	bridgeUUID, err := addFullBridge(bridgeName, true, manager)
	if err != nil {
		t.Error(err)
	}
	defer deleteBridge(bridgeUUID, manager)

	// try to add a port with the same name of the default one which has the bridge name in it
	// THIS MUST FAIL !
	portUUID, intUUID, err := addPort(bridgeName, bridgeName, bridgeUUID, vlan, INTERFACE_INTERNAL, containerPID, taskID, manager)
	if err == nil {
		t.Error("We should not be able to add another port with the same name !!")
	}

	// add a port
	portUUID, intUUID, err = addPort(additionalPortName, additionalPortName, bridgeUUID, vlan, INTERFACE_INTERNAL, containerPID, taskID, manager)
	if err != nil {
		t.Error(err)
	}

	// PORT
	// check that the port and an interface was successfully added
	port, portVlan, err := getVPort(portUUID, manager)
	if err != nil {
		t.Error(err)
	}

	if port.Id != portUUID {
		t.Error("Port UUIDs do not match")
	}

	if portVlan != vlan {
		t.Error("Port vlan info do not match")
	}

	if port.Tag != vlan {
		t.Error("Port vlan info do not match")
	}

	if port.Name != additionalPortName {
		t.Error("Port name info do not match with bridge name")
	}
	// ================================================================

	// INTERFACE
	if len(port.Interfaces) != 1 {
		t.Error("Expected exactly one interface on the port got: " + strconv.Itoa(len(port.Interfaces)))
	}

	vint := port.Interfaces[0]
	if vint.Id != intUUID {
		t.Error("Interface UUIDs do not match")
	}

	if vint.Name != additionalPortName {
		t.Error("Interface name does not match")
	}

}

func TestGetVswitch(t *testing.T) {

	vlan := 1

	// this should create the default switch with VLAN 0
	vswitch, err := InitVSwitch(ovsdbHost, ovsdbPort)

	if err != nil {
		t.Fatal(err)
	}

	defer vswitch.Close()

	// add a new VPC
	vpc, err := vswitch.AddVPC(vlan)
	//fmt.Println(vpc)
	if err != nil {
		t.Error(err)
	}

	if len(vswitch.VPCs) != 2 {
		t.Error("There should be exactly 2 VPCs")
	}

	if len(vpc.Ports) != 1 {
		t.Error("There should be exactly 1 Port instead we got: " + strconv.Itoa(len(vpc.Ports)))
	}

	// Load the VPC
	storedVpc, err := getVPC(vpc.Id, vswitch.manager)
	//fmt.Println(storedVpc)
	if err != nil {
		t.Error(err)
	}

	if storedVpc.Id != vpc.Id {
		t.Error("VPC IDs do not match")
	}

	// --------------------------------------
	var port VPort
	for _, v := range vpc.Ports {
		port = v
	}

	if port.Id == "" {
		t.Error("Port Id cannot be empty")
	}
	// --------------------------------------
	var storedPort VPort
	for _, v := range storedVpc.Ports {
		storedPort = v
	}

	if storedPort.Id == "" {
		t.Error("Stored Port Id cannot be empty")
	}
	// --------------------------------------

	if storedPort.Id != port.Id {
		t.Error("The port IDs do not match")
	}

	if storedPort.Name != port.Name {
		t.Error("The port names do not match")
	}

	// CLEANUP
	// delete the newly created VPC
	if err := vswitch.DeleteVPC(vlan); err != nil {
		t.Error(err)
	}

	// delete the DEFAULT created VPC
	if err := vswitch.DeleteVPC(0); err != nil {
		t.Error(err)
	}

	if len(vswitch.VPCs) != 0 {
		t.Error("There should be exactly 0 VPCs")
	}

}

func test() {

}
