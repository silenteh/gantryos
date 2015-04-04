package vswitch

import "encoding/json"

type transactOperations []interface{}

// Operation represents an operation according to RFC7047 section 5.2
type operation struct {
	Op        string                   `json:"op"`
	Table     string                   `json:"table,omitempty"`
	Row       map[string]interface{}   `json:"row,omitempty"`
	Rows      []map[string]interface{} `json:"rows,omitempty"`
	Columns   []string                 `json:"columns,omitempty"`
	Mutations []interface{}            `json:"mutations,omitempty"`
	Timeout   int                      `json:"timeout,omitempty"`
	Where     []interface{}            `json:"where,omitempty"`
	Until     string                   `json:"until,omitempty"`
	UUIDName  string                   `json:"uuid-name,omitempty"`
	Durable   bool                     `json:"durable,omitempty"`
}

type commitOperation struct {
	Op      string `json:"op"`
	Durable bool   `json:"durable"`
}

// MonitorRequest represents a monitor request according to RFC7047
/*
 * We cannot use MonitorRequests by inlining the MonitorRequest Map structure till GoLang issue #6213 makes it.
 * The only option is to go with raw map[string]interface{} option :-( that sucks !
 * Refer to client.go : MonitorAll() function for more details
 */

type monitorRequests struct {
	Requests map[string]monitorRequest `json:"requests,overflow"`
}

// MonitorRequest represents a monitor request according to RFC7047
type monitorRequest struct {
	Columns []string      `json:"columns,omitempty"`
	Select  monitorSelect `json:"select,omitempty"`
}

// MonitorSelect represents a monitor select according to RFC7047
type monitorSelect struct {
	Initial bool `json:"initial,omitempty"`
	Insert  bool `json:"insert,omitempty"`
	Delete  bool `json:"delete,omitempty"`
	Modify  bool `json:"modify,omitempty"`
}

/*
 * We cannot use TableUpdates directly by json encoding by inlining the TableUpdate Map
 * structure till GoLang issue #6213 makes it.
 *
 * The only option is to go with raw map[string]map[string]interface{} option :-( that sucks !
 * Refer to client.go : MonitorAll() function for more details
 */
type tableUpdates struct {
	Updates map[string]tableUpdate `json:"updates,overflow"`
}

type tableUpdate struct {
	Rows map[string]rowUpdate `json:"rows,overflow"`
}

type rowUpdate struct {
	Uuid UUID `json:"-,omitempty"`
	New  row  `json:"new,omitempty"`
	Old  row  `json:"old,omitempty"`
}

// OvsdbError is an OVS Error Condition
type ovsdbError struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// NewCondition creates a new condition as specified in RFC7047
func newCondition(column string, function string, value interface{}) []interface{} {
	return []interface{}{column, function, value}
}

// NewMutation creates a new mutation as specified in RFC7047
func newMutation(column string, mutator string, value interface{}) []interface{} {
	return []interface{}{column, mutator, value}
}

type transactResponse struct {
	Result []operationResult `json:"result"`
	Error  string            `json:"error"`
}

type operationResult struct {
	Count   int                      `json:"count,omitempty"`
	Error   string                   `json:"error,omitempty"`
	Details string                   `json:"details,omitempty"`
	UUID    UUID                     `json:"uuid,omitempty"`
	Rows    []map[string]interface{} `json:"rows,omitempty"`
}

func ovsSliceToGoNotation(val interface{}) (interface{}, error) {
	switch val.(type) {
	case []interface{}:
		sl := val.([]interface{})
		bsliced, err := json.Marshal(sl)
		if err != nil {
			return nil, err
		}

		switch sl[0] {
		case "uuid":
			var uuid UUID
			err = json.Unmarshal(bsliced, &uuid)
			return uuid, err
		case "set":
			var oSet ovsSet
			err = json.Unmarshal(bsliced, &oSet)
			return oSet, err
		case "map":
			var oMap ovsMap
			err = json.Unmarshal(bsliced, &oMap)
			return oMap, err
		}
		return val, nil
	}
	return val, nil
}

// TODO : add Condition, Function, Mutation and Mutator notations
