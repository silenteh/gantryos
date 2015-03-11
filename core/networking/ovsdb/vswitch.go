package networking

// http://json-rpc.org/wiki/specification

import (
	//"fmt"
	//"github.com/socketplane/libovsdb"

	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"reflect"
	//"io"

	//"net/rpc"
	//"net/rpc/jsonrpc"
)

//==========================================
// RPC Codec
// ----------------------------------------------------------------------------
// Request and Response
// ----------------------------------------------------------------------------
// An Error is a wrapper for a JSON interface value. It can be used by either
// a service's handler func to write more complex JSON data to an error field
// of a server's response, or by a client to read it.
type Error struct {
	Data interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v", e.Data)
}

type ovsdbOperation string

type monitor struct {
	string
	int
	ovsdbOperation
}

type vswitchManager struct {
	host     string
	port     string
	client   *rpcjJsonClient
	handlers []NotificationHandler
	schema   map[string]DatabaseSchema
	cache    map[string]map[string]Row
}

type NotificationHandler interface {
	// RFC 7047 section 4.1.6 Update Notification
	Update(context interface{}, tableUpdates TableUpdates)

	// RFC 7047 section 4.1.9 Locked Notification
	Locked([]interface{})

	// RFC 7047 section 4.1.10 Stolen Notification
	Stolen([]interface{})

	// RFC 7047 section 4.1.11 Echo Notification
	Echo([]interface{})
}

func NewOVSDBClient(host, port string) (*vswitchManager, error) {

	manager := vswitchManager{}
	// c, err := jsonrpc.Dial("tcp4", host+":"+port)
	// if err != nil {
	// 	return &manager, err
	// }

	c := NewRPCJsonClient(host, port)
	if err := c.Connect(); err != nil {
		return &manager, err
	}

	manager.host = host
	manager.port = port
	manager.client = c
	manager.schema = make(map[string]DatabaseSchema)
	manager.cache = make(map[string]map[string]Row)
	manager.handlers = []NotificationHandler{}

	// monitor and register the changes in the manager cache
	notifier := Notifier{manager: manager}
	manager.Register(notifier)

	// get the schema
	if _, err := manager.GetSchema("Open_vSwitch"); err != nil {
		return &manager, err
	}

	// start to monitor and populate the cache
	initial, err := manager.MonitorAll("Open_vSwitch", "")
	if err != nil {
		return &manager, err
	}
	manager.populateCache(*initial)

	return &manager, nil

}

func (manager vswitchManager) Close() error {
	return manager.client.Close()
}

func (manager vswitchManager) Echo() ([]string, error) {
	var resp []string
	return resp, manager.client.Call("echo", []string{"ping"}, &resp)
}

func (manager vswitchManager) ListDBs() []string {
	var resp []string
	manager.client.Call("list_dbs", []string{""}, &resp)
	return resp
}

func (manager vswitchManager) GetSchema(db string) (*DatabaseSchema, error) {

	var reply DatabaseSchema
	err := manager.client.Call("get_schema", []string{db}, &reply)
	if err != nil {
		return nil, err
	} else {
		manager.schema[db] = reply
	}
	return &reply, err
}

func (manager vswitchManager) AddBridge(bridgeName string) {

	// SELECT EXAMPLE
	// condition := NewCondition("_uuid", "==", UUID{manager.GetRootUUID()})
	// insertBridge := selectBaseOp(condition)

	// INSERT INTERFACE
	operations := addBridgeOps(bridgeName, manager.GetRootUUID(), false)

	reply, err := manager.Transact("Open_vSwitch", operations...)
	if err != nil {
		fmt.Println(err)
	}

	if len(reply) < len(operations) {
		fmt.Println("Number of Replies should be atleast equal to number of Operations")
	}
	//ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			fmt.Println("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			//ok = false
		} else if o.Error != "" {
			fmt.Println("Transaction Failed due to an error :", o.Error)
			//ok = false
		}
	}

	//fmt.Println(reply[0])

	//manager.AddPort(bridgeName)
	if len(reply) > 0 {
		fmt.Println("Bridge Addition Successful : ", reply[0].UUID.GoUuid)
	}
}

func (manager vswitchManager) AddPort(bridgeName string) {
	//namedUuid := "gantryos"
	// bridge row to insert
	//bridge := make(map[string]interface{})
	//bridge["name"] = bridgeName
	// 	// simple insert operation
	port := make(map[string]interface{})
	//interfaces := make(map[string]interface{})
	//interfaces["named-uuid"] = "new_interface"
	port["name"] = bridgeName
	//port["interfaces"] = []map[string]interface{}{}
	insertPort := Operation{
		Op:       "insert",
		Table:    "Port",
		Row:      port,
		UUIDName: bridgeName + "_port",
	}

	// Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
	// mutateUuid := []UUID{UUID{bridgeName + "_port"}}
	// mutateSet, _ := NewOvsSet(mutateUuid)
	// mutation := NewMutation("ports", "insert", mutateSet)
	// condition := NewCondition("_uuid", "==", UUID{"gantryos"})

	// // simple mutate operation
	// mutateOp := Operation{
	// 	Op:        "mutate",
	// 	Table:     "Bridge",
	// 	Mutations: []interface{}{mutation},
	// 	Where:     []interface{}{condition},
	// }

	// // 	// simple insert operation
	// brInterface := make(map[string]interface{})
	// brInterface["name"] = bridgeName
	// brInterface["type"] = "internal"
	// insertInterface := Operation{
	// 	Op:       "insert",
	// 	Table:    "Interface",
	// 	Row:      brInterface,
	// 	UUIDName: "new_interface",
	// }

	operations := []Operation{insertPort}
	reply, _ := manager.Transact("Open_vSwitch", operations...)

	if len(reply) < len(operations) {
		fmt.Println("Number of Replies should be atleast equal to number of Operations")
	}
	//ok := true
	for i, o := range reply {
		if o.Error != "" && i < len(operations) {
			fmt.Println("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", operations[i])
			//ok = false
		} else if o.Error != "" {
			fmt.Println("Transaction Failed due to an error :", o.Error)
			//ok = false
		}
	}
	// if ok {
	// 	fmt.Println("Bridge Addition Successful : ", reply[0].UUID.GoUuid)
	// }
}

// func (manager vswitchManager) AddBridge(name string) interface{} {

// 	bridges := make(map[string]interface{})

// 	firstSwitch := make(map[string]interface{})
// 	firstSwitch["named-uuid"] = "test_switch"
// 	bridges["bridges"] = []map[string]interface{}{firstSwitch}
// 	// simple insert operation
// 	insertSwitch := Operation{
// 		Op:       "insert",
// 		Table:    "Open_vSwitch",
// 		Row:      bridges,
// 		UUIDName: "test_switch",
// 	}

// 	// simple insert operation
// 	brInterface := make(map[string]interface{})
// 	brInterface["name"] = name
// 	brInterface["type"] = "internal"
// 	insertInterface := Operation{
// 		Op:       "insert",
// 		Table:    "Interface",
// 		Row:      brInterface,
// 		UUIDName: "new_interface",
// 	}

// 	// simple insert operation
// 	port := make(map[string]interface{})
// 	interfaces := make(map[string]interface{})
// 	interfaces["named-uuid"] = "new_interface"
// 	port["name"] = name
// 	port["interfaces"] = []map[string]interface{}{interfaces}
// 	insertPort := Operation{
// 		Op:       "insert",
// 		Table:    "Port",
// 		Row:      port,
// 		UUIDName: "new_port",
// 	}

// 	// simple insert operation
// 	bridge := make(map[string]interface{})
// 	bridgesMap := make(map[string]interface{})
// 	bridgesMap["named-uuid"] = "new_port"
// 	bridge["name"] = name
// 	bridge["ports"] = bridgesMap
// 	insertBridge := Operation{
// 		Op:       "insert",
// 		Table:    "Bridge",
// 		Row:      bridge,
// 		UUIDName: "new_bridge",
// 	}

// 	// namedUuid := "new_switch"
// 	// // bridge row to insert
// 	// bridge := make(map[string]interface{})
// 	// bridge["name"] = name

// 	// // simple insert operation
// 	// insertOp := Operation{
// 	// 	Op:       "insert",
// 	// 	Table:    "Bridge",
// 	// 	Row:      bridge,
// 	// 	UUIDName: namedUuid,
// 	// }

// 	// // Inserting a Bridge row in Bridge table requires mutating the open_vswitch table.
// 	// mutateUuid := []UUID{UUID{namedUuid}}
// 	// mutateSet, _ := NewOvsSet(mutateUuid)
// 	// mutation := NewMutation("bridges", "insert", mutateSet)
// 	//condition := NewCondition("_uuid", "==", UUID{"00000000-0000-0000-0000-000000000000"})

// 	// // simple mutate operation
// 	// mutateOp := Operation{
// 	// 	Op:        "mutate",
// 	// 	Table:     "Open_vSwitch",
// 	// 	Mutations: []interface{}{mutation},
// 	// 	Where:     []interface{}{},
// 	// }

// 	operations := []Operation{insertSwitch, insertInterface, insertPort, insertBridge}
// 	args := NewTransactArgs("Open_vSwitch", operations...)
// 	//reply, _ := //ovs.Transact("Open_vSwitch", operations...)

// 	var resp interface{}

// 	manager.client.Call("transact", args, &resp)
// 	return resp

// 	// args := `
// 	//        {
// 	//            "row": {
// 	//                "bridges": [
// 	//                    "named-uuid",
// 	//                    "new_bridge"
// 	//                ]
// 	//            },
// 	//            "table": "Open_vSwitch",
// 	//            "uuid-name": "new_switch",
// 	//            "op": "insert"
// 	//        },
// 	//        {
// 	//            "row": {
// 	//                "name": "br1",
// 	//                "type": "internal"
// 	//            },
// 	//            "table": "Interface",
// 	//            "uuid-name": "new_interface",
// 	//            "op": "insert"
// 	//        },
// 	//        {
// 	//            "row": {
// 	//                "name": "br1",
// 	//                "interfaces": [
// 	//                    "named-uuid",
// 	//                    "new_interface"
// 	//                ]
// 	//            },
// 	//            "table": "Port",
// 	//            "uuid-name": "new_port",
// 	//            "op": "insert"
// 	//        },
// 	//        {
// 	//            "row": {
// 	//                "name": "br1",
// 	//                "ports": [
// 	//                    "named-uuid",
// 	//                    "new_port"
// 	//                ]
// 	//            },
// 	//            "table": "Bridge",
// 	//            "uuid-name": "new_bridge",
// 	//            "op": "insert"
// 	//        }
// 	//    `

// 	// var resp interface{}

// 	// manager.client.Call("transact", []string{"Open_vSwitch", args}, &resp)
// 	// // return resp

// }

func (manager vswitchManager) Monitor() interface{} {

	return nil
}

func (manager vswitchManager) Register(handler NotificationHandler) {
	manager.handlers = append(manager.handlers, handler)
}

// Convenience method to monitor every table/column
func (manager vswitchManager) MonitorAll(database string, jsonContext interface{}) (*TableUpdates, error) {
	schema, ok := manager.schema[database]
	if !ok {
		return nil, errors.New("invalid Database Schema")
	}

	requests := make(map[string]MonitorRequest)
	for table, tableSchema := range schema.Tables {
		var columns []string
		for column, _ := range tableSchema.Columns {
			columns = append(columns, column)
		}
		requests[table] = MonitorRequest{
			Columns: columns,
			Select: MonitorSelect{
				Initial: true,
				Insert:  true,
				Delete:  true,
				Modify:  true,
			}}
	}
	return manager.monitor(database, jsonContext, requests)
}

// RFC 7047 : monitor
func (manager vswitchManager) monitor(database string, jsonContext interface{}, requests map[string]MonitorRequest) (*TableUpdates, error) {
	var reply TableUpdates

	args := NewMonitorArgs(database, jsonContext, requests)

	// This totally sucks. Refer to golang JSON issue #6213
	var response map[string]map[string]RowUpdate
	err := manager.client.Call("monitor", args, &response)
	reply = getTableUpdatesFromRawUnmarshal(response)
	if err != nil {
		return nil, err
	}
	return &reply, err
}

func getTableUpdatesFromRawUnmarshal(raw map[string]map[string]RowUpdate) TableUpdates {
	var tableUpdates TableUpdates
	tableUpdates.Updates = make(map[string]TableUpdate)
	for table, update := range raw {
		tableUpdate := TableUpdate{update}
		tableUpdates.Updates[table] = tableUpdate
	}
	return tableUpdates
}

func (manager vswitchManager) GetRootUUID() string {
	for uuid, _ := range manager.cache["Open_vSwitch"] {
		return uuid
	}
	return ""
}

func (manager vswitchManager) populateCache(updates TableUpdates) {
	for table, tableUpdate := range updates.Updates {
		if _, ok := manager.cache[table]; !ok {
			manager.cache[table] = make(map[string]Row)

		}
		for uuid, row := range tableUpdate.Rows {
			empty := Row{}
			if !reflect.DeepEqual(row.New, empty) {
				manager.cache[table][uuid] = row.New
			} else {
				delete(manager.cache[table], uuid)
			}
		}
	}
}

type Notifier struct {
	manager vswitchManager
}

func (n Notifier) Update(context interface{}, tableUpdates TableUpdates) {
	fmt.Println("Got update from monitor")
	n.manager.populateCache(tableUpdates)
	for k, v := range tableUpdates.Updates {
		log.Infoln(k, v)
		fmt.Println(k, v)
	}
}
func (n Notifier) Locked([]interface{}) {
	fmt.Println("Got locked from monitor")
	log.Infoln("Locked")
}
func (n Notifier) Stolen([]interface{}) {
	fmt.Println("Got stolen from monitor")
	log.Infoln("Stolen")
}
func (n Notifier) Echo([]interface{}) {
	fmt.Println("Got echo from monitor")
	log.Infoln("Echo")
}

// func (n Notifier) Disconnected() {
// 	n.manager.client.Close()
// }

func (manager vswitchManager) Transact(database string, operation ...Operation) ([]OperationResult, error) {
	var reply []OperationResult
	db, ok := manager.schema[database]
	if !ok {
		return nil, errors.New("invalid Database Schema")
	}

	if ok := db.validateOperations(operation...); !ok {
		return nil, errors.New("Validation failed for the operation")
	}

	args := NewTransactArgs(database, operation...)
	err := manager.client.Call("transact", args, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
