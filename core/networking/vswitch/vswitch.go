package vswitch

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
	//handlers []notificationHandler
	schema map[string]databaseSchema
	cache  map[string]map[string]row
}

type notificationHandler interface {
	// RFC 7047 section 4.1.6 Update Notification
	Update(context interface{}, tableUpdates tableUpdates)

	// RFC 7047 section 4.1.9 Locked Notification
	Locked([]interface{})

	// RFC 7047 section 4.1.10 Stolen Notification
	Stolen([]interface{})

	// RFC 7047 section 4.1.11 Echo Notification
	Echo([]interface{})
}

var transactCounter = 0

func newOVSDBClient(host, port string) (*vswitchManager, error) {

	manager := vswitchManager{}

	c := newRPCJsonClient(host, port)
	if err := c.Connect(); err != nil {
		return &manager, err
	}

	manager.host = host
	manager.port = port
	manager.client = &c
	manager.schema = make(map[string]databaseSchema)
	manager.cache = make(map[string]map[string]row)

	// monitor and register the changes in the manager cache
	notifier := notifier{manager: manager}
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
		break
	case <-time.After(5 * time.Second):
		err = errors.New("List DBs request timed out")
		break
	}

	return response, err
}

func (manager vswitchManager) GetSchema(db string) (*databaseSchema, error) {

	var dbSchema databaseSchema
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
		break
	case <-time.After(5 * time.Second):
		err = errors.New("Get schema request timed out")
	}

	return &dbSchema, err
}

func (manager *vswitchManager) AddBridge(bridgeName string, stpEnabled bool) (string, error) {

	return addFullBridge(bridgeName, manager.GetRootUUID(), stpEnabled, manager)
}

func (manager *vswitchManager) DeleteBridge(bridgeName string) error {

	// get the uuid of the br0
	uuidBridge, err := getBridgeUUID(bridgeName, manager)
	if err != nil {
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

func (manager vswitchManager) Register(handler notificationHandler) {
	manager.client.addnotificationHandler(handler)
}

// Convenience method to monitor every table/column
func (manager vswitchManager) MonitorAll(database string, jsonContext interface{}) (*tableUpdates, error) {
	schema, ok := manager.schema[database]
	if !ok {
		return nil, errors.New("invalid Database Schema")
	}

	requests := make(map[string]monitorRequest)
	for table, tableSchema := range schema.Tables {
		var columns []string
		for column, _ := range tableSchema.Columns {
			columns = append(columns, column)
		}
		requests[table] = monitorRequest{
			Columns: columns,
			Select: monitorSelect{
				Initial: true,
				Insert:  true,
				Delete:  true,
				Modify:  true,
			}}
	}
	return manager.monitor(database, jsonContext, requests)
}

// // RFC 7047 : monitor
func (manager vswitchManager) monitor(database string, jsonContext interface{}, requests map[string]monitorRequest) (*tableUpdates, error) {
	var reply tableUpdates

	args := newMonitorArgs(database, jsonContext, requests)

	// This totally sucks. Refer to golang JSON issue #6213
	var response map[string]map[string]rowUpdate
	dataChan, err := manager.client.Call("monitor", args, true)

	data := <-dataChan
	err = json.Unmarshal(data, &response)

	reply = gettableUpdatesFromRawUnmarshal(response)
	if err != nil {
		return nil, err
	}
	return &reply, err
}

func gettableUpdatesFromRawUnmarshal(raw map[string]map[string]rowUpdate) tableUpdates {
	var tableUpdates tableUpdates
	tableUpdates.Updates = make(map[string]tableUpdate)
	for table, update := range raw {
		tableUpdate := tableUpdate{update}
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

func (manager vswitchManager) populateCache(updates tableUpdates) {
	for table, tableUpdate := range updates.Updates {
		if _, ok := manager.cache[table]; !ok {
			manager.cache[table] = make(map[string]row)

		}
		for uuid, trow := range tableUpdate.Rows {
			empty := row{}
			if !reflect.DeepEqual(trow.New, empty) {
				manager.cache[table][uuid] = trow.New

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

type notifier struct {
	manager vswitchManager
}

func (n notifier) Update(context interface{}, tableUpdates tableUpdates) {
	//fmt.Println("Got update from monitor")
	log.Infoln("Update received")
	n.manager.populateCache(tableUpdates)
	// for k, v := range tableUpdates.Updates {
	// 	log.Infoln(k, v)
	// 	fmt.Println(k, v)
	// }
}
func (n notifier) Locked([]interface{}) {
	//fmt.Println("Got locked from monitor")
	log.Infoln("Locked")
}
func (n notifier) Stolen([]interface{}) {
	//fmt.Println("Got stolen from monitor")
	log.Infoln("Stolen")
}
func (n notifier) Echo([]interface{}) {
	//fmt.Println("Got echo from monitor")
	log.Infoln("Handler got Echo")

}

// func (n notifier) Disconnected() {
// 	n.manager.client.Close()
// }

func (manager *vswitchManager) Transact(database, description string, operations transactOperations) ([]operationResult, error) {
	var err error
	var response []operationResult

	dataChan, err := manager.client.Call("transact", operations, true)

	if err != nil {
		fmt.Println("TRANSACT:", err)
		return response, err
	}

	// blocks
	select {
	case data := <-dataChan:
		err = json.Unmarshal(data, &response)
		break
	case <-time.After(5 * time.Second):
		err = errors.New(description + ": Transact request timed out")
	}

	return response, err
}

func (manager *vswitchManager) lock(trasnactionId string) {
	var locked map[string]bool
	id := newLockArgs("lock_" + trasnactionId)

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
	id := newLockArgs("unlock_" + trasnactionId)

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
