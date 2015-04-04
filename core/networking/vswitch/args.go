package vswitch

// NewGetSchemaArgs creates a new set of arguments for a get_schemas RPC
func newGetSchemaArgs(schema string) []interface{} {
	return []interface{}{schema}
}

// NewTransactArgs creates a new set of arguments for a transact RPC
func newTransactArgs(database string, addCommit bool, operations ...operation) transactOperations {
	var dbSlice = make([]interface{}, 1)
	dbSlice[0] = database

	var opsSlice []interface{} = make([]interface{}, len(operations))
	for i, d := range operations {
		opsSlice[i] = d
	}

	if addCommit {
		opsSlice = append(opsSlice, newCommitOp())
	}

	ops := append(dbSlice, opsSlice...)

	var transactOps transactOperations
	transactOps = ops
	return transactOps
}

// NewCancelArgs creates a new set of arguments for a cancel RPC
func newCancelArgs(id interface{}) []interface{} {
	return []interface{}{id}
}

// NewMonitorArgs creates a new set of arguments for a monitor RPC
func newMonitorArgs(database string, value interface{}, requests map[string]monitorRequest) []interface{} {
	return []interface{}{database, value, requests}
}

// NewMonitorCancelArgs creates a new set of arguments for a monitor_cancel RPC
func newMonitorCancelArgs(value interface{}) []interface{} {
	return []interface{}{value}
}

// NewLockArgs creates a new set of arguments for a lock, steal or unlock RPC
func newLockArgs(id interface{}) []interface{} {
	return []interface{}{id}
}
