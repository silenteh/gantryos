package networking

func addBridgeOps(bridgeName, rootUUID string, stpEnabled bool) []Operation {

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

	return []Operation{insertInterfaceOp, insertPortOp, insertBridgeOp, mutateOp}

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

// func newAutoAttachOp(bridgeName string) (Operation, string) {
// 	namedUuid := bridgeName + "_autoattach"
// 	emptyRow := make(map[string]interface{})
// 	insertOp := Operation{
// 		Op:       "insert",
// 		Table:    "AutoAttach",
// 		Row:      emptyRow,
// 		UUIDName: namedUuid,
// 	}
// 	return insertOp, namedUuid
// }

func selectBaseOp(condition interface{}) Operation {
	// simple insert operation
	//brInterface := make(map[string]interface{})
	//brInterface["name"] = name
	//brInterface["type"] = "internal"
	insertOp := Operation{
		Op:    "select",
		Table: "Open_vSwitch",
		Where: []interface{}{condition},
		//Row:      brInterface,
		//UUIDName: name + "_interface",
	}
	return insertOp
}
