package vswitch

import "encoding/json"

type row struct {
	Fields map[string]interface{}
}

func (r *row) UnmarshalJSON(b []byte) (err error) {
	r.Fields = make(map[string]interface{})
	var raw map[string]interface{}
	err = json.Unmarshal(b, &raw)
	for key, val := range raw {
		val, err = ovsSliceToGoNotation(val)
		if err != nil {
			return err
		}
		r.Fields[key] = val
	}
	return err
}
