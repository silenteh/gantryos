package ovsdb

import (
	"fmt"
)

func newCommitOp() CommitOperation {
	commit := CommitOperation{
		Durable: false,
		Op:      "commit",
	}
	return commit
}

func addBridgeOps(bridgeName, rootUUID string, stpEnabled bool) TransactOperations {

	// interface first, because we need the UUID
	insertInterfaceOp, insertInterfaceUUID := newInterfaceOp(bridgeName)

	// then the port using the interface UUID
	insertPortOp, insertPortUUID := newPortOp(bridgeName, insertInterfaceUUID)

	// finally the bridge with the port UUID
	insertBridgeOp, insertBridgeUUID := newBridgeOp(bridgeName, insertPortUUID, stpEnabled)

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

	operations := []Operation{insertInterfaceOp, insertPortOp, insertBridgeOp, mutateOp}

	return NewTransactArgs("Open_vSwitch", true, operations...)

}

func addBridge(bridgeName, rootUUID string, stpEnabled bool, manager *vswitchManager) (string, error) {

	namedUuid := bridgeName + "_gantryos"

	// bridge definition and properties
	bridge := make(map[string]interface{})
	bridge["name"] = bridgeName
	bridge["stp_enable"] = stpEnabled

	// assign the port to the bridge
	//bridge["ports"] = newNamedUUID(portUUID)

	// create the operation
	insertOp := Operation{
		Op:       "insert",
		Table:    "Bridge",
		Row:      bridge,
		UUIDName: namedUuid,
	}

	// Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{namedUuid}}
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

	operations := []Operation{insertOp, mutateOp}

	ops := NewTransactArgs("Open_vSwitch", true, operations...)

	results, err := manager.Transact("Open_vSwitch", "ADD_BRIDGE", ops)
	if err != nil {
		return "", err
	}

	res := ParseOVSDBOpsResult(results[0])

	return res.UUID.GoUuid, err

}

func addPort(portName, bridgeUUID string, vlan int, manager *vswitchManager) (string, error) {

	// interface first, because we need the UUID
	insertInterfaceOp, insertInterfaceUUID := newInterfaceOp(portName)

	// then the port using the interface UUID
	insertPortOp, insertPortUUID := newPortOp(portName, insertInterfaceUUID)

	// Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{insertPortUUID}}
	mutateSet, _ := NewOvsSet(mutateUuid)
	mutation := NewMutation("ports", "insert", mutateSet)
	condition := NewCondition("_uuid", "==", UUID{bridgeUUID})

	// simple mutate operation
	mutateOp := Operation{
		Op:        "mutate",
		Table:     "Bridge",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []Operation{insertInterfaceOp, insertPortOp, mutateOp}

	ops := NewTransactArgs("Open_vSwitch", true, operations...)

	results, err := manager.Transact("Open_vSwitch", "ADD_FIRST_PORT", ops)
	if err != nil {
		return "", err
	}

	//fmt.Println(results)

	res := ParseOVSDBOpsResult(results[1])

	return res.UUID.GoUuid, err

}

func addInterface(interfaceName, portUUID string, manager *vswitchManager) (string, error) {

	// interface first, because we need the UUID
	insertInterfaceOp, insertInterfaceUUID := newInterfaceOp(interfaceName)

	mutateUuid := []UUID{UUID{insertInterfaceUUID}}
	mutateSet, _ := NewOvsSet(mutateUuid)
	mutation := NewMutation("interfaces", "insert", mutateSet)
	condition := NewCondition("_uuid", "==", UUID{portUUID})

	// simple mutate operation
	mutateOp := Operation{
		Op:        "mutate",
		Table:     "Port",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []Operation{insertInterfaceOp, mutateOp}

	ops := NewTransactArgs("Open_vSwitch", true, operations...)

	results, err := manager.Transact("Open_vSwitch", "ADD_INTERFACE", ops)
	if err != nil {
		return "", err
	}

	fmt.Println(results)

	res := ParseOVSDBOpsResult(results[0])

	return res.UUID.GoUuid, err

}

func newBridgeOp(bridgeName, portUUID string, stpEnabled bool) (Operation, string) {

	namedUuid := bridgeName + "_gantryos"

	// bridge definition and properties
	bridge := make(map[string]interface{})
	bridge["name"] = bridgeName
	bridge["stp_enable"] = stpEnabled

	// assign the port to the bridge
	bridge["ports"] = newNamedUUID(portUUID)

	// create the operation
	insertOp := Operation{
		Op:       "insert",
		Table:    "Bridge",
		Row:      bridge,
		UUIDName: namedUuid,
	}

	return insertOp, namedUuid

}

func newPortOp(bridgeName, interfaceUUID string) (Operation, string) {
	namedUuid := bridgeName + "_port"

	// port definition
	port := make(map[string]interface{})
	port["name"] = bridgeName

	// assign the interface to the port
	port["interfaces"] = newNamedUUID(interfaceUUID)

	// create the operation
	insertOp := Operation{
		Op:       "insert",
		Table:    "Port",
		Row:      port,
		UUIDName: namedUuid,
	}
	return insertOp, namedUuid
}

func newInterfaceOp(bridgeName string) (Operation, string) {
	namedUuid := bridgeName + "_interface"
	// simple insert operation
	brInterface := make(map[string]interface{})
	brInterface["name"] = bridgeName
	brInterface["type"] = "internal"
	insertOp := Operation{
		Op:       "insert",
		Table:    "Interface",
		Row:      brInterface,
		UUIDName: namedUuid,
	}
	return insertOp, namedUuid
}

func deleteBridge(rootUUID, bridgeUUID string, manager *vswitchManager) error {

	deleteBridgeOperation := Operation{
		Op:    "delete",
		Table: "Bridge",
		Where: []interface{}{NewCondition("_uuid", "==", UUID{bridgeUUID})},
	}

	// // Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{bridgeUUID}}
	mutateSet, _ := NewOvsSet(mutateUuid)
	mutation := NewMutation("bridges", "delete", mutateSet)
	condition := NewCondition("_uuid", "==", UUID{rootUUID})

	// // simple mutate operation
	mutateOp := Operation{
		Op:        "mutate",
		Table:     "Open_vSwitch",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []Operation{deleteBridgeOperation, mutateOp}

	ops := NewTransactArgs("Open_vSwitch", true, operations...)

	_, err := manager.Transact("Open_vSwitch", "DELETE_BRIDGE", ops)
	return err
}

func getBridgeUUID(bridgeName string, manager *vswitchManager) (string, error) {
	condition := NewCondition("name", "==", bridgeName)

	selectBridgeOp := Operation{
		Op:    "select",
		Table: "Bridge",
		Where: []interface{}{condition},
	}
	operations := []Operation{selectBridgeOp}

	ops := NewTransactArgs("Open_vSwitch", false, operations...)

	data, err := manager.Transact("Open_vSwitch", "GET_BRIDGE_UUID", ops)
	if err != nil {
		return "", err
	}

	uuidBridge := data[0].UUID.GoUuid
	if len(data[0].Rows) > 0 {
		uuidBridge = ParseOVSDBUUID(data[0].Rows[0]["_uuid"])
	}

	return uuidBridge, nil

}

func getPortUUID(portName string, manager *vswitchManager) (string, error) {
	condition := NewCondition("name", "==", portName)

	selectPortOp := Operation{
		Op:    "select",
		Table: "Port",
		Where: []interface{}{condition},
	}
	operations := []Operation{selectPortOp}

	ops := NewTransactArgs("Open_vSwitch", false, operations...)

	data, err := manager.Transact("Open_vSwitch", "GET_PORT_UUID", ops)
	if err != nil {
		return "", err
	}

	uuidPort := data[0].UUID.GoUuid
	if len(data[0].Rows) > 0 {
		uuidPort = ParseOVSDBUUID(data[0].Rows[0]["_uuid"])
	}

	return uuidPort, nil

}

func getAllPorts(portName string, manager *vswitchManager) (string, error) {
	condition := NewCondition("name", "==", portName)

	selectPortOp := Operation{
		Op:    "select",
		Table: "Port",
		Where: []interface{}{condition},
	}
	operations := []Operation{selectPortOp}

	ops := NewTransactArgs("Open_vSwitch", false, operations...)

	data, err := manager.Transact("Open_vSwitch", "GET_PORT_UUID", ops)
	if err != nil {
		return "", err
	}

	uuidPort := data[0].UUID.GoUuid
	if len(data[0].Rows) > 0 {
		uuidPort = ParseOVSDBUUID(data[0].Rows[0]["_uuid"])
	}

	return uuidPort, nil

}

func getInterfaceUUID(interfaceName string, manager *vswitchManager) (string, error) {
	condition := NewCondition("name", "==", interfaceName)

	selectInterfaceOp := Operation{
		Op:    "select",
		Table: "Interface",
		Where: []interface{}{condition},
	}
	operations := []Operation{selectInterfaceOp}

	ops := NewTransactArgs("Open_vSwitch", false, operations...)

	data, err := manager.Transact("Open_vSwitch", "GET_INTERFACE_UUID", ops)
	if err != nil {
		return "", err
	}

	uuidInterface := data[0].UUID.GoUuid
	if len(data[0].Rows) > 0 {
		uuidInterface = ParseOVSDBUUID(data[0].Rows[0]["_uuid"])
	}

	return uuidInterface, nil

}

func deletePort(portUUID string, manager *vswitchManager) error {
	deletePortOperation := Operation{
		Op:    "delete",
		Table: "Port",
		Where: []interface{}{NewCondition("_uuid", "==", UUID{portUUID})},
	}
	operations := []Operation{deletePortOperation}
	ops := NewTransactArgs("Open_vSwitch", true, operations...)
	_, err := manager.Transact("Open_vSwitch", "DELETE_PORT", ops)
	return err

}

func deleteInterface(interfaceUUID string, manager *vswitchManager) error {
	deleteInterfaceOperation := Operation{
		Op:    "delete",
		Table: "Interface",
		Where: []interface{}{NewCondition("_uuid", "==", UUID{interfaceUUID})},
	}
	operations := []Operation{deleteInterfaceOperation}
	ops := NewTransactArgs("Open_vSwitch", true, operations...)
	_, err := manager.Transact("Open_vSwitch", "DELETE_PORT", ops)
	return err

}

func selectBridgeOp(bridgeName string) []Operation {
	// simple insert operation
	//brInterface := make(map[string]interface{})
	//brInterface["name"] = name

	condition := NewCondition("name", "==", bridgeName)

	selectBridgeOp := Operation{
		Op:    "select",
		Table: "Bridge",
		Where: []interface{}{condition},
	}
	selectPortOp := Operation{
		Op:    "select",
		Table: "Port",
		Where: []interface{}{condition},
	}
	selectInterfaceOp := Operation{
		Op:    "select",
		Table: "Interface",
		Where: []interface{}{condition},
	}
	return []Operation{selectBridgeOp, selectPortOp, selectInterfaceOp}
}

func bridgeExists(bridgeName string, manager *vswitchManager) (bool, string) {

	id, err := getBridgeUUID(bridgeName, manager)
	if err == nil && id != "" {
		return true, id
	}

	return false, ""

}
