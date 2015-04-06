package vswitch

//import "fmt"
import "strconv"

func parseOVSString(data interface{}) string {
	if parsed, ok := data.(string); ok {
		return parsed
	}
	return ""
}

func parseVLANTag(data interface{}) int {
	if set, ok := data.([]interface{}); ok {
		ovsset, _ := newOvsSet(set)
		for _, singleInter := range ovsset.GoSet {
			if single, ok := singleInter.(string); ok {
				vlan, _ := strconv.ParseInt(single, 0, 64)
				return int(vlan)
			}
		}
	}
	if single, ok := data.(float64); ok {
		return int(single)
	}
	return 0
}

func parseOVSMap(data interface{}) map[string]string {
	gosMap := make(map[string]string)
	if array, ok := data.([]interface{}); ok {
		if len(array) > 0 {
			if ovsMap, ok := array[1].([]interface{}); ok {
				for _, item := range ovsMap {
					if stringArray, ok := item.([]interface{}); ok {
						if len(stringArray) > 1 {
							gosMap[parseOVSString(stringArray[0])] = parseOVSString(stringArray[1])
						}
					}
				}
			}
		}
	}
	return gosMap
}

func parseOVSDBUUID(data interface{}) string {

	if array, ok := data.([]interface{}); ok {
		uuid := array[1].(string)
		return uuid //string(array[1])
	}
	return ""
}

func parseOVSDBOpsResult(data interface{}) operationResult {

	if res, ok := data.(operationResult); ok {
		return res
	}
	return operationResult{}
}
