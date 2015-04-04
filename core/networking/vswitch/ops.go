package vswitch

import (
	"errors"
	//"fmt"
	"strconv"
)

var INTERFACE_SYSTEM = ""
var INTERFACE_INTERNAL = "internal"
var INTERFACE_VLAN = "vlan"

func newCommitOp() commitOperation {
	commit := commitOperation{
		Durable: false,
		Op:      "commit",
	}
	return commit
}

// generates the bridge + port + interface operations
func addBridgeOps(bridgeName, rootUUID string, stpEnabled bool) transactOperations {

	var containerId string
	var containerIface string

	// interface first, because we need the UUID
	insertInterfaceOp, insertInterfaceUUID := newInterfaceOp(bridgeName, INTERFACE_INTERNAL)

	// then the port using the interface UUID
	insertPortOp, insertPortUUID := newPortOp(bridgeName, insertInterfaceUUID, "default", containerId, containerIface, 0)

	// finally the bridge with the port UUID
	insertBridgeOp, insertBridgeUUID := newBridgeOp(bridgeName, insertPortUUID, stpEnabled)

	// Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{insertBridgeUUID}}
	mutateSet, _ := newOvsSet(mutateUuid)
	mutation := newMutation("bridges", "insert", mutateSet)
	condition := newCondition("_uuid", "==", UUID{rootUUID})

	// simple mutate operation
	mutateOp := operation{
		Op:        "mutate",
		Table:     "Open_vSwitch",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []operation{insertInterfaceOp, insertPortOp, insertBridgeOp, mutateOp}

	return newTransactArgs("Open_vSwitch", true, operations...)

}

// Add a bridge + port + interface
func addFullBridge(bridgeName, rootUUID string, stpEnabled bool, manager *vswitchManager) (string, error) {
	//fmt.Println("ROOT UUID", manager.GetRootUUID())
	ops := addBridgeOps(bridgeName, manager.GetRootUUID(), false)

	results, err := manager.Transact("Open_vSwitch", "ADD_BRIDGE", ops)
	if err != nil {
		return "", err
	}
	err = checkForErrors(results)

	res := parseOVSDBOpsResult(results[0])

	return res.UUID.GoUuid, err
}

// Add a bridge only without port and interfaces
func addBridge(bridgeName, rootUUID string, stpEnabled bool, manager *vswitchManager) (string, error) {

	namedUuid := bridgeName + "_gantryos"

	// bridge definition and properties
	bridge := make(map[string]interface{})
	bridge["name"] = bridgeName
	bridge["stp_enable"] = stpEnabled

	// assign the port to the bridge
	//bridge["ports"] = newNamedUUID(portUUID)

	// create the operation
	insertOp := operation{
		Op:       "insert",
		Table:    "Bridge",
		Row:      bridge,
		UUIDName: namedUuid,
	}

	// Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{namedUuid}}
	mutateSet, _ := newOvsSet(mutateUuid)
	mutation := newMutation("bridges", "insert", mutateSet)
	condition := newCondition("_uuid", "==", UUID{rootUUID})

	// simple mutate operation
	mutateOp := operation{
		Op:        "mutate",
		Table:     "Open_vSwitch",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []operation{insertOp, mutateOp}

	ops := newTransactArgs("Open_vSwitch", true, operations...)

	results, err := manager.Transact("Open_vSwitch", "ADD_BRIDGE", ops)
	if err != nil {
		return "", err
	}

	err = checkForErrors(results)

	res := parseOVSDBOpsResult(results[0])

	return res.UUID.GoUuid, err

}

// Add a port only with a single interface attached to it (necessary by ovsdb)
// returns the Port UUID, Interface UUID, error
func addPort(portName, bridgeUUID, vpcName string, vlan int, vInt VInterface, manager *vswitchManager) (string, string, error) {

	// interface first, because we need the UUID
	insertInterfaceOp, insertInterfaceUUID := newInterfaceOp(portName, vInt.Type)

	// then the port using the interface UUID
	insertPortOp, insertPortUUID := newPortOp(portName, insertInterfaceUUID, vInt.ContainerId, vInt.ContainerIface, vpcName, vlan)

	// Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{insertPortUUID}}
	mutateSet, _ := newOvsSet(mutateUuid)
	mutation := newMutation("ports", "insert", mutateSet)
	condition := newCondition("_uuid", "==", UUID{bridgeUUID})

	// simple mutate operation
	mutateOp := operation{
		Op:        "mutate",
		Table:     "Bridge",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []operation{insertInterfaceOp, insertPortOp, mutateOp}

	ops := newTransactArgs("Open_vSwitch", true, operations...)

	results, err := manager.Transact("Open_vSwitch", "ADD_FIRST_PORT", ops)
	if err != nil {
		return "", "", err
	}

	//fmt.Println(results)

	err = checkForErrors(results)

	portRes := parseOVSDBOpsResult(results[1])
	intRes := parseOVSDBOpsResult(results[0])

	return portRes.UUID.GoUuid, intRes.UUID.GoUuid, err

}

// Add an additional interface to a port
func addInterface(interfaceName, portUUID, intType string, manager *vswitchManager) (string, error) {

	// interface first, because we need the UUID
	insertInterfaceOp, insertInterfaceUUID := newInterfaceOp(interfaceName, intType)

	mutateUuid := []UUID{UUID{insertInterfaceUUID}}
	mutateSet, _ := newOvsSet(mutateUuid)
	mutation := newMutation("interfaces", "insert", mutateSet)
	condition := newCondition("_uuid", "==", UUID{portUUID})

	// simple mutate operation
	mutateOp := operation{
		Op:        "mutate",
		Table:     "Port",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []operation{insertInterfaceOp, mutateOp}

	ops := newTransactArgs("Open_vSwitch", true, operations...)

	results, err := manager.Transact("Open_vSwitch", "ADD_INTERFACE", ops)
	if err != nil {
		return "", err
	}

	err = checkForErrors(results)
	res := parseOVSDBOpsResult(results[0])

	return res.UUID.GoUuid, err

}

func newBridgeOp(bridgeName, portUUID string, stpEnabled bool) (operation, string) {

	namedUuid := bridgeName + "_gantryos"

	// bridge definition and properties
	bridge := make(map[string]interface{})
	bridge["name"] = bridgeName
	bridge["stp_enable"] = stpEnabled

	// assign the port to the bridge
	bridge["ports"] = newNamedUUID(portUUID)

	// create the operation
	insertOp := operation{
		Op:       "insert",
		Table:    "Bridge",
		Row:      bridge,
		UUIDName: namedUuid,
	}

	return insertOp, namedUuid

}

func newPortOp(bridgeName, interfaceUUID, vpcName, containerId, containerIface string, vlan int) (operation, string) {
	namedUuid := bridgeName + "_port"

	// port definition
	port := make(map[string]interface{})
	port["name"] = bridgeName
	port["tag"] = vlan

	// ovs map
	gosMap := make(map[string]string)
	gosMap["gos-vpc-name"] = vpcName
	gosMap["gos-vpc-id"] = strconv.Itoa(vlan)
	if containerId != "" {
		gosMap["container_id"] = containerId
	}
	if containerIface != "" {
		gosMap["container_iface"] = containerIface
	}

	if ovsMap, err := newOvsMap(gosMap); err == nil {
		port["external_ids"] = ovsMap
	}

	// port["external_ids"] =  //make(map[string]string)

	// if gosMap, ok := port["external_ids"].(map[string]string); ok {
	// 	gosMap["gos-vpc-name"] = vpcName
	// 	gosMap["gos-vpc-id"] = strconv.Itoa(vlan)
	// }

	// assign the interface to the port
	port["interfaces"] = newNamedUUID(interfaceUUID)

	// create the operation
	insertOp := operation{
		Op:       "insert",
		Table:    "Port",
		Row:      port,
		UUIDName: namedUuid,
	}
	return insertOp, namedUuid
}

func newInterfaceOp(bridgeName, intType string) (operation, string) {
	namedUuid := bridgeName + "_interface"
	// simple insert operation
	brInterface := make(map[string]interface{})
	brInterface["name"] = bridgeName
	if intType != "" {
		brInterface["type"] = intType
	}

	insertOp := operation{
		Op:       "insert",
		Table:    "Interface",
		Row:      brInterface,
		UUIDName: namedUuid,
	}
	return insertOp, namedUuid
}

func deleteBridge(rootUUID, bridgeUUID string, manager *vswitchManager) error {

	deleteBridgeoperation := operation{
		Op:    "delete",
		Table: "Bridge",
		Where: []interface{}{newCondition("_uuid", "==", UUID{bridgeUUID})},
	}

	// // Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{bridgeUUID}}
	mutateSet, _ := newOvsSet(mutateUuid)
	mutation := newMutation("bridges", "delete", mutateSet)
	condition := newCondition("_uuid", "==", UUID{rootUUID})

	// // simple mutate operation
	mutateOp := operation{
		Op:        "mutate",
		Table:     "Open_vSwitch",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []operation{deleteBridgeoperation, mutateOp}

	ops := newTransactArgs("Open_vSwitch", true, operations...)

	results, err := manager.Transact("Open_vSwitch", "DELETE_BRIDGE", ops)

	if err != nil {
		return err
	}
	err = checkForErrors(results)

	return err
}

// Get the bridge UUID from the bridge name
func getBridgeUUID(bridgeName string, manager *vswitchManager) (string, error) {
	condition := newCondition("name", "==", bridgeName)

	selectBridgeOp := operation{
		Op:    "select",
		Table: "Bridge",
		Where: []interface{}{condition},
	}
	operations := []operation{selectBridgeOp}

	ops := newTransactArgs("Open_vSwitch", false, operations...)

	results, err := manager.Transact("Open_vSwitch", "GET_BRIDGE_UUID", ops)
	if err != nil {
		return "", err
	}

	err = checkForErrors(results)

	uuidBridge := results[0].UUID.GoUuid
	if len(results[0].Rows) > 0 {
		uuidBridge = parseOVSDBUUID(results[0].Rows[0]["_uuid"])
	}

	return uuidBridge, err

}

// Get the port UUID from the port name
func getPortUUID(portName string, manager *vswitchManager) (string, error) {
	condition := newCondition("name", "==", portName)

	selectPortOp := operation{
		Op:    "select",
		Table: "Port",
		Where: []interface{}{condition},
	}
	operations := []operation{selectPortOp}

	ops := newTransactArgs("Open_vSwitch", false, operations...)

	results, err := manager.Transact("Open_vSwitch", "GET_PORT_UUID", ops)
	if err != nil {
		return "", err
	}

	err = checkForErrors(results)

	uuidPort := results[0].UUID.GoUuid
	if len(results[0].Rows) > 0 {
		uuidPort = parseOVSDBUUID(results[0].Rows[0]["_uuid"])
	}

	return uuidPort, nil

}

// returns the UUIDs of the ports setup on the bridge
func getAllBridgePorts(bridgeUUID, rootUUID string, manager *vswitchManager) (Vswitch, error) {

	vswitch := Vswitch{
		RootId: rootUUID,
		Id:     bridgeUUID,
		VPCs:   make(map[int]vpc),
	}

	condition := newCondition("_uuid", "==", UUID{bridgeUUID})

	selectPortOp := operation{
		Op:    "select",
		Table: "Bridge",
		Where: []interface{}{condition},
	}
	operations := []operation{selectPortOp}

	ops := newTransactArgs("Open_vSwitch", false, operations...)

	data, err := manager.Transact("Open_vSwitch", "GET_ALL_PORT_UUID", ops)
	if err != nil {
		return vswitch, err
	}

	err = checkForErrors(data)

	if len(data[0].Rows) > 0 {
		vswitch.Name = parseOVSString(data[0].Rows[0]["name"])
		set, err := newOvsSet(data[0].Rows[0]["ports"])
		if err != nil {
			return vswitch, err
		}

		// get the VPorts UUIDs
		uuids := set.GetUUIDs()

		// range over them
		for _, portUUID := range uuids {
			// get the VPort info
			port, vpcName, vlan, portError := getVPort(portUUID, manager)

			if portError != nil {
				err = portError
			}
			// check if we have already a VPC associated with the vswitch
			// if not create it
			if _, ok := vswitch.VPCs[vlan]; !ok {
				vswitch.VPCs[vlan] = vpc{
					Name:  vpcName,
					VLan:  vlan,
					Ports: make(map[string]VPort),
				}

			}

			// if the VPort has Id and Name then add it
			if port.Id != "" && port.Name != "" {
				vswitch.VPCs[vlan].Ports[port.Name] = port
			}
		}
	}

	return vswitch, err

}

func getVPort(portUUID string, manager *vswitchManager) (VPort, string, int, error) {
	port := newVPort()
	var vpcName string
	var vlan int = -1
	var containerId string
	var containerIface string

	condition := newCondition("_uuid", "==", UUID{portUUID})

	selectPortOp := operation{
		Op:    "select",
		Table: "Port",
		Where: []interface{}{condition},
	}
	operations := []operation{selectPortOp}

	ops := newTransactArgs("Open_vSwitch", false, operations...)

	data, err := manager.Transact("Open_vSwitch", "GET_PORT_INFO", ops)
	if err != nil {
		return port, vpcName, vlan, err
	}

	//fmt.Println(data[0].Rows)

	err = checkForErrors(data)

	port.Id = portUUID
	if len(data[0].Rows) > 0 {
		port.Name = parseOVSString(data[0].Rows[0]["name"])

		// get VPC info
		if ovsMap, ok := data[0].Rows[0]["external_ids"].([]interface{}); ok {
			gosMap := parseOVSMap(ovsMap)
			//fmt.Println(gosMap)
			vpcName = gosMap["gos-vpc-name"]
			tempVlan, _ := strconv.ParseInt(gosMap["gos-vpc-id"], 0, 64)
			vlan = int(tempVlan)
			containerId = gosMap["container_id"]
			containerIface = gosMap["container_iface"]
		}

		set, err := newOvsSet(data[0].Rows[0]["interfaces"])
		if err != nil {
			return port, vpcName, vlan, err
		}
		for _, intUUID := range set.GetUUIDs() {
			vInt, intErr := getVInterface(intUUID, manager)

			if intErr != nil {
				err = intErr
			}

			vInt.ContainerId = containerId
			vInt.ContainerIface = containerIface

			if vInt.Name != "" {
				port.Interfaces[vInt.Name] = vInt
			}
		}
	}

	return port, vpcName, vlan, err
}

func getVInterface(interfaceUUID string, manager *vswitchManager) (VInterface, error) {
	vInt := VInterface{}

	condition := newCondition("_uuid", "==", UUID{interfaceUUID})

	selectPortOp := operation{
		Op:    "select",
		Table: "Interface",
		Where: []interface{}{condition},
	}
	operations := []operation{selectPortOp}

	ops := newTransactArgs("Open_vSwitch", false, operations...)

	data, err := manager.Transact("Open_vSwitch", "GET_INTERFACE_INFO_UUID", ops)
	if err != nil {
		return vInt, err
	}

	//fmt.Println(data[0].Rows)

	err = checkForErrors(data)

	vInt.Id = interfaceUUID
	if len(data[0].Rows) > 0 {
		vInt.Name = parseOVSString(data[0].Rows[0]["name"])
	}

	return vInt, err

}

// returns the UUIDs of the ports setup on the bridge
func getAllPortInterfaces(portUUID string, manager *vswitchManager) ([]string, error) {
	condition := newCondition("_uuid", "==", UUID{portUUID})

	selectPortOp := operation{
		Op:    "select",
		Table: "Port",
		Where: []interface{}{condition},
	}
	operations := []operation{selectPortOp}

	ops := newTransactArgs("Open_vSwitch", false, operations...)

	data, err := manager.Transact("Open_vSwitch", "GET_ALL_INTERFACES_UUID", ops)
	if err != nil {
		return nil, err
	}

	err = checkForErrors(data)

	if len(data[0].Rows) > 0 {
		set, intErr := newOvsSet(data[0].Rows[0]["interfaces"])
		if intErr != nil {
			err = intErr
		}
		return set.GetUUIDs(), err
	}

	return []string{}, err

}

func getInterfaceUUID(interfaceName string, manager *vswitchManager) (string, error) {
	condition := newCondition("name", "==", interfaceName)

	selectInterfaceOp := operation{
		Op:    "select",
		Table: "Interface",
		Where: []interface{}{condition},
	}
	operations := []operation{selectInterfaceOp}

	ops := newTransactArgs("Open_vSwitch", false, operations...)

	data, err := manager.Transact("Open_vSwitch", "GET_INTERFACE_UUID", ops)
	if err != nil {
		return "", err
	}

	err = checkForErrors(data)

	uuidInterface := data[0].UUID.GoUuid
	if len(data[0].Rows) > 0 {
		uuidInterface = parseOVSDBUUID(data[0].Rows[0]["_uuid"])
	}

	return uuidInterface, err

}

func deletePort(bridgeUUID, portUUID string, manager *vswitchManager) error {
	deletePortoperation := operation{
		Op:    "delete",
		Table: "Port",
		Where: []interface{}{newCondition("_uuid", "==", UUID{portUUID})},
	}

	// Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{portUUID}}
	mutateSet, _ := newOvsSet(mutateUuid)
	mutation := newMutation("ports", "delete", mutateSet)
	condition := newCondition("_uuid", "==", UUID{bridgeUUID})

	// simple mutate operation
	mutateOp := operation{
		Op:        "mutate",
		Table:     "Bridge",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []operation{mutateOp, deletePortoperation}
	ops := newTransactArgs("Open_vSwitch", true, operations...)
	data, err := manager.Transact("Open_vSwitch", "DELETE_PORT", ops)
	if err != nil {
		return err
	}
	return checkForErrors(data)

}

func deleteInterface(portUUID, interfaceUUID string, manager *vswitchManager) error {
	deleteInterfaceoperation := operation{
		Op:    "delete",
		Table: "Interface",
		Where: []interface{}{newCondition("_uuid", "==", UUID{interfaceUUID})},
	}

	// Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	mutateUuid := []UUID{UUID{interfaceUUID}}
	mutateSet, _ := newOvsSet(mutateUuid)
	mutation := newMutation("interfaces", "delete", mutateSet)
	condition := newCondition("_uuid", "==", UUID{portUUID})

	// simple mutate operation
	mutateOp := operation{
		Op:        "mutate",
		Table:     "Port",
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}

	operations := []operation{deleteInterfaceoperation, mutateOp}
	ops := newTransactArgs("Open_vSwitch", true, operations...)
	data, err := manager.Transact("Open_vSwitch", "DELETE_INTERFACE", ops)

	if err != nil {
		return err
	}
	return checkForErrors(data)

}

func selectBridgeOp(bridgeName string) []operation {

	condition := newCondition("name", "==", bridgeName)

	selectBridgeOp := operation{
		Op:    "select",
		Table: "Bridge",
		Where: []interface{}{condition},
	}
	selectPortOp := operation{
		Op:    "select",
		Table: "Port",
		Where: []interface{}{condition},
	}
	selectInterfaceOp := operation{
		Op:    "select",
		Table: "Interface",
		Where: []interface{}{condition},
	}
	return []operation{selectBridgeOp, selectPortOp, selectInterfaceOp}
}

// Check if a bridge exists
func bridgeExists(bridgeName string, manager *vswitchManager) (bool, string) {

	id, err := getBridgeUUID(bridgeName, manager)
	if err == nil && id != "" {
		return true, id
	}

	return false, ""

}

func checkForErrors(results []operationResult) error {
	var err error
	for _, result := range results {
		errorOp := parseOVSDBOpsResult(result)
		if errorOp.Error != "" || errorOp.Details != "" {
			return errors.New(errorOp.Error + ":" + errorOp.Details)
		}
	}
	return err
}
