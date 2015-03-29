package ovsdb

import (
	"errors"
)

// MODELS

// local slave switch
type vswitch struct {
	Id   string // ovsdb root ID
	Name string
	VPCs []vpc
}

// can contain multiple VPCs
type vpc struct {
	Name    string  // description name must be unique
	Network string  // network range
	VLan    int     // vlan ID
	Ports   []vPort // all ports have that ID
}

// each VPC has multiple ports
type vPort struct {
	Id         string // uuid
	Name       string
	Interfaces []vInterface
}

// each port has an interface
type vInterface struct {
	Id   string // uuid
	Name string
}

// func LoadVSwitch(bridgeName string, manager vswitchManager) (*vswitch, error) {

// 	condition := NewCondition("name", "==", bridgeName)

// 	selectBridgeOp := Operation{
// 		Op:    "select",
// 		Table: "Bridge",
// 		Where: []interface{}{condition},
// 	}

// 	operations := []Operation{insertBridgeOp, mutateOp}

// 	results, err := manager.Transact("Open_vSwitch", operations...)

// }

func NewVSwitch(rootUUID, bridgeName string, stpEnabled bool, manager vswitchManager) (*vswitch, error) {

	vswitch := vswitch{
		Id:   rootUUID,
		Name: bridgeName,
		VPCs: []vpc{},
	}

	insertBridgeUUID := bridgeName + "_gantryos"

	// bridge definition and properties
	bridge := make(map[string]interface{})
	bridge["name"] = bridgeName
	bridge["stp_enable"] = stpEnabled

	// assign the port to the bridge
	//bridge["ports"] = newNamedUUID(portUUID)

	// create the operation
	insertBridgeOp := Operation{
		Op:       "insert",
		Table:    "Bridge",
		Row:      bridge,
		UUIDName: insertBridgeUUID,
	}

	// Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{insertBridgeUUID}}
	mutateSet, _ := NewOvsSet(mutateUuid)
	mutation := NewMutation("bridges", "insert", mutateSet)
	condition := NewCondition("_uuid", "==", UUID{rootUUID})

	// simple mutate operation
	mutateOp := Operation{
		Op:        "mutate",
		Table:     "Open_vSwitch",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []Operation{insertBridgeOp, mutateOp, newCommitOp()}

	results, err := manager.Transact("Open_vSwitch", "NEW_VSWITCH", operations...)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("Error adding a new vswitch")
	}

	vswitch.Id = results[0].UUID.GoUuid

	return &vswitch, nil

}
