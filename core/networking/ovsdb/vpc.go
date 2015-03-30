package ovsdb

import (
//"errors"
//"fmt"
)

// MODELS

// local slave switch
type vswitch struct {
	RootId string // ovsdb root UUID
	Id     string // bridge UUID
	Name   string
	VPCs   map[string]vpc
}

// can contain multiple VPCs
type vpc struct {
	Name    string           // description name must be unique
	Network string           // network range
	VLan    int              // vlan ID
	Ports   map[string]vPort // all ports have that ID
}

// each VPC has multiple ports
type vPort struct {
	Id         string // uuid
	Name       string
	Interfaces map[string]vInterface
}

// each port has an interface
type vInterface struct {
	Id   string // uuid
	Name string
}

func (vswitch *vswitch) AddVPC(name, network string, vlan int) {
	vpc := vpc{
		Name:    name,
		Network: network,
		VLan:    vlan,
		Ports:   make(map[string]vPort),
	}

	vswitch.VPCs[name] = vpc
}

func (vpc vpc) AddPort(portName, bridgeUUID string, vlan int, manager *vswitchManager) error {
	port := vPort{
		Name:       portName,
		Interfaces: make(map[string]vInterface),
	}

	id, err := addPort(portName, bridgeUUID, vlan, manager)
	if err != nil {
		return err
	}

	port.Id = id
	vpc.Ports[portName] = port

	return nil
}

func (port vPort) AddInterface(interfaceName string, manager *vswitchManager) error {

	interfaceUUID, err := addInterface(interfaceName, port.Id, manager)
	if err != nil {
		return err
	}

	vint := vInterface{
		Id:   interfaceUUID,
		Name: interfaceName,
	}

	port.Interfaces[interfaceName] = vint

	return nil
}

// func loadVSwitch(bridgeName string, manager vswitchManager) (*vswitch, error) {

// 	condition := NewCondition("name", "==", bridgeName)

// 	selectBridgeOp := Operation{
// 		Op:    "select",
// 		Table: "Bridge",
// 		Where: []interface{}{condition},
// 	}

// 	operations := []Operation{insertBridgeOp, mutateOp}

// 	results, err := manager.Transact("Open_vSwitch", operations...)

// }

func NewVSwitch(rootUUID, bridgeName string, stpEnabled bool, manager *vswitchManager) (*vswitch, error) {

	exists, id := bridgeExists(bridgeName, manager)

	vswitch := vswitch{
		RootId: rootUUID,
		Name:   bridgeName,
		VPCs:   make(map[string]vpc),
	}

	if exists {
		vswitch.Id = id
		return &vswitch, nil
	}

	bridgeUUID, err := addBridge(bridgeName, rootUUID, stpEnabled, manager)

	if err != nil {
		return &vswitch, err
	}

	vswitch.Id = bridgeUUID

	return &vswitch, nil

}
