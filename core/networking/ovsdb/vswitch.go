package ovsdb

// http://json-rpc.org/wiki/specification

import (
	//"fmt"
	//"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"reflect"
	"time"
	//"time"
	//"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	//"strconv"
)

//==========================================

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

type vswitchManager struct {
	host   string
	port   string
	client *rpcjJsonClient
	//handlers []NotificationHandler
	schema   map[string]DatabaseSchema
	cache    map[string]map[string]Row
	vpcCache vswitch
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

var transactCounter = 0

func NewOVSDBClient(host, port string) (*vswitchManager, error) {

	manager := vswitchManager{}

	c := NewRPCJsonClient(host, port)
	if err := c.Connect(); err != nil {
		return &manager, err
	}

	manager.host = host
	manager.port = port
	manager.client = &c
	manager.schema = make(map[string]DatabaseSchema)
	manager.cache = make(map[string]map[string]Row)

	// monitor and register the changes in the manager cache
	notifier := Notifier{manager: manager}
	manager.Register(notifier)

	//get the schema
	if _, err := manager.GetSchema("Open_vSwitch"); err != nil {
		return &manager, err
	}

	//start to monitor and populate the cache
	initial, err := manager.MonitorAll("Open_vSwitch", "")
	if err != nil {
		return &manager, err
	}

	//fmt.Println(initial)

	manager.populateCache(*initial)

	return &manager, nil

}

func (manager vswitchManager) Close() error {
	return manager.client.Close()
}

func (manager vswitchManager) Echo() ([]string, error) {
	var response []string
	dataChan, err := manager.client.Call("echo", []string{"ping"}, true)
	if err != nil {
		return response, err
	}

	// blocks
	select {
	case data := <-dataChan:
		err = json.Unmarshal(data, &response)
	case <-time.After(5 * time.Second):
		err = errors.New("Echo request timed out")
	}
	return response, err
}

func (manager vswitchManager) ListDBs() ([]string, error) {
	var response []string
	dataChan, err := manager.client.Call("list_dbs", []string{}, true)
	if err != nil {
		return response, err
	}

	// blocks
	select {
	case data := <-dataChan:
		err = json.Unmarshal(data, &response)
	case <-time.After(5 * time.Second):
		err = errors.New("List DBs request timed out")
	}

	return response, err
}

func (manager vswitchManager) GetSchema(db string) (*DatabaseSchema, error) {

	var dbSchema DatabaseSchema
	dataChan, err := manager.client.Call("get_schema", []string{db}, true)
	if err != nil {
		//fmt.Println("ERROR GETTING DB SCHEMA !", err)
		return nil, err
	}

	// blocks
	select {
	case data := <-dataChan:
		err = json.Unmarshal(data, &dbSchema)
		if err == nil {
			manager.schema[db] = dbSchema
		}
	case <-time.After(5 * time.Second):
		err = errors.New("Get schema request timed out")
	}

	return &dbSchema, err
}

func (manager vswitchManager) AddBridge(bridgeName string) error {

	// SELECT EXAMPLE
	// condition := NewCondition("_uuid", "==", UUID{manager.GetRootUUID()})
	// insertBridge := selectBaseOp(condition)

	// INSERT INTERFACE
	fmt.Println("ROOT UUID", manager.GetRootUUID())
	ops := addBridgeOps(bridgeName, manager.GetRootUUID(), false)

	_, err := manager.Transact("Open_vSwitch", "ADD_BRIDGE", ops)

	return err
}

func (manager *vswitchManager) DeleteBridge(bridgeName string) error {

	// get the uuid of the br0
	uuidBridge, err := getBridgeUUID(bridgeName, manager)
	if err != nil {
		return err
	}

	// PORT
	uuidPort, err := getPortUUID(bridgeName, manager)
	if err != nil {
		return err
	}

	// INTERFACE
	uuidInterface, err := getPortUUID(bridgeName, manager)
	if err != nil {
		return err
	}

	if err := deleteInterface(uuidInterface, manager); err != nil {
		return err
	}

	if err := deletePort(uuidPort, manager); err != nil {
		return err
	}

	if err := deleteBridge(manager.GetRootUUID(), uuidBridge, manager); err != nil {
		return err
	}

	return nil

}

func (manager vswitchManager) Monitor() interface{} {

	return nil
}

func (manager vswitchManager) Register(handler NotificationHandler) {
	manager.client.AddNotificationHandler(handler)
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

// // RFC 7047 : monitor
func (manager vswitchManager) monitor(database string, jsonContext interface{}, requests map[string]MonitorRequest) (*TableUpdates, error) {
	var reply TableUpdates

	args := NewMonitorArgs(database, jsonContext, requests)

	// This totally sucks. Refer to golang JSON issue #6213
	var response map[string]map[string]RowUpdate
	dataChan, err := manager.client.Call("monitor", args, true)

	data := <-dataChan
	err = json.Unmarshal(data, &response)

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

				switch table {
				case "Bridge":
				case "Port":
				case "Interface":
				}

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
	//fmt.Println("Got update from monitor")
	log.Infoln("Update received")
	n.manager.populateCache(tableUpdates)
	// for k, v := range tableUpdates.Updates {
	// 	log.Infoln(k, v)
	// 	fmt.Println(k, v)
	// }
}
func (n Notifier) Locked([]interface{}) {
	//fmt.Println("Got locked from monitor")
	log.Infoln("Locked")
}
func (n Notifier) Stolen([]interface{}) {
	//fmt.Println("Got stolen from monitor")
	log.Infoln("Stolen")
}
func (n Notifier) Echo([]interface{}) {
	//fmt.Println("Got echo from monitor")
	log.Infoln("Echo")

}

// func (n Notifier) Disconnected() {
// 	n.manager.client.Close()
// }

func (manager *vswitchManager) Transact(database, description string, operations TransactOperations) ([]OperationResult, error) {
	var err error
	var response []OperationResult
	// db, ok := manager.schema[database]
	// if !ok {
	// 	return nil, errors.New("invalid Database Schema")
	// }

	// if ok := db.validateOperations(operation...); !ok {
	// 	return nil, errors.New("Validation failed for the operation")
	// }

	//args := NewTransactArgs(database, operation...)

	// // Increment the counter
	//transactCounter++
	// // LOCK id
	//id := strconv.Itoa(transactCounter)

	// // LOCK response
	//var locked map[string]bool

	// // Lock call
	//manager.client.Call("lock", id, &locked)

	// fmt.Println("LOCKED: =====> ", locked)
	// for !locked["locked"] {
	// 	fmt.Println("Loop", "LOCKING.......")
	// 	manager.client.Call("lock", id, &locked)
	// 	time.Sleep(500 * time.Millisecond)
	// }

	//manager.lock(id)
	//defer manager.unlock(id)

	dataChan, err := manager.client.Call("transact", operations, true)
	if err != nil {
		fmt.Println("TRANSACT:", err)
		return response, err
	}

	// blocks
	select {
	case data := <-dataChan:
		err = json.Unmarshal(data, &response)
	case <-time.After(30 * time.Second):
		err = errors.New(description + ": Transact request timed out")
	}

	return response, err

	// UNLOCK Response
	//var unlocked map[string]bool

	// UNLOCK ID
	//id = NewLockArgs("unlock_" + strconv.Itoa(transactCounter))
	//manager.client.Call("unlock", id, &unlocked)

	// switch reply.(type) {
	// case string:
	// 	fmt.Println("STRING:", reply.(string))
	// 	return nil, errors.New("Got back a string")
	// case []interface{}:
	// 	fmt.Println("INTERFACE:", reply)
	// 	var result []OperationResult
	// 	data, _ := json.Marshal(reply)
	// 	json.Unmarshal(data, &result)
	// 	return result, nil
	// default:
	// 	fmt.Println(reply)
	// }

	// if err != nil {
	// 	//manager.client.Call("unlock", id, &unlocked)
	// 	//fmt.Println("UNLOCKED: =====> ", unlocked)
	// 	fmt.Println("FAILED", err)
	// 	return nil, err
	// }
	//manager.client.Call("unlock", id, &unlocked)
	//fmt.Println("UNLOCKED: =====> ", unlocked)
	//return nil, nil
}

func (manager *vswitchManager) lock(trasnactionId string) {
	var locked map[string]bool
	id := NewLockArgs("lock_" + trasnactionId)

	dataChan, err := manager.client.Call("lock", id, true)
	if err != nil {
		log.Errorln("LOCK:", err)
		return
	}

	// blocks
	select {
	case data := <-dataChan:
		err = json.Unmarshal(data, &locked)
	case <-time.After(10 * time.Second):
		err = errors.New("LOCK" + ": Transact request timed out")
	}
}

func (manager *vswitchManager) unlock(trasnactionId string) {
	var locked map[string]bool
	id := NewLockArgs("unlock_" + trasnactionId)

	dataChan, err := manager.client.Call("unlock", id, true)
	if err != nil {
		log.Errorln("UNLOCK:", err)
		return
	}

	// blocks
	select {
	case data := <-dataChan:
		err = json.Unmarshal(data, &locked)
	case <-time.After(10 * time.Second):
		err = errors.New("UNLOCK" + ": Transact request timed out")
	}
}
