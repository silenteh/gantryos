package ovsdb

import (
	"encoding/json"
	"errors"
	//"fmt"
	//"reflect"
	"regexp"
)

type UUID struct {
	GoUuid string `json:"uuid"`
}

func ParseOVSDBUUID(data interface{}) string {

	if array, ok := data.([]interface{}); ok {
		uuid := array[1].(string)
		return uuid //string(array[1])
	}

	return ""

}

func ParseOVSDBOpsResult(data interface{}) OperationResult {

	if res, ok := data.(OperationResult); ok {
		return res
	}

	return OperationResult{}

}

// <set> notation requires special marshaling
func (u UUID) MarshalJSON() ([]byte, error) {
	var uuidSlice []string
	err := u.validateUUID()
	if err == nil {
		uuidSlice = []string{"uuid", u.GoUuid}
	} else {
		uuidSlice = []string{"named-uuid", u.GoUuid}
	}

	return json.Marshal(uuidSlice)
}

func (u *UUID) UnmarshalJSON(b []byte) (err error) {
	var ovsUuid []string
	if err := json.Unmarshal(b, &ovsUuid); err == nil {
		u.GoUuid = ovsUuid[1]
	}
	return err
}

func (u UUID) validateUUID() error {
	if len(u.GoUuid) != 36 {
		return errors.New("uuid exceeds 36 characters")
	}

	var validUUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	if !validUUID.MatchString(u.GoUuid) {
		return errors.New("uuid does not match regexp")
	}

	return nil
}

func newNamedUUID(id string) []string {

	//uuidMap := make(map[string]string)
	//uuidMap["named-uuid"] = id
	//return []map[string]string{uuidMap}
	return []string{"named-uuid", id}

}
