package models

import (
	"github.com/silenteh/gantryos/core/proto"
)

type Resource struct {
	Type      proto.ResourceType // type of resource
	ValueType proto.Value_Type   // type of resource
	Value     *float64           // cpu - mem
	Range     []float64          // ports - bandwidth
}

func NewCPUResource(cpuValue float64) *Resource {
	return &Resource{
		Type:      CPU_RESOURCE_TYPE,
		ValueType: proto.Value_SCALAR,
		Value:     &cpuValue,
		Range:     nil,
	}
}

func NewMEMResource(memValue float64) *Resource {
	return &Resource{
		Type:      MEM_RESOURCE_TYPE,
		ValueType: proto.Value_SCALAR,
		Value:     &memValue,
		Range:     nil,
	}
}

func NewPortsResource(begin, end float64) *Resource {
	ports := []float64{begin, end}
	return &Resource{
		Type:      PORTS_RESOURCE_TYPE,
		ValueType: proto.Value_RANGES,
		Value:     nil,
		Range:     ports,
	}
}
