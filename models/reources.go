package models

import (
	"github.com/silenteh/gantryos/core/proto"
)

type resources []*resource

type resource struct {
	Type      proto.ResourceType // type of resource
	ValueType proto.Value_Type   // type of resource
	Value     *float64           // cpu - mem
	Ranges    []RangeValue       // ports - bandwidth
}

type RangeValue struct {
	Begin uint64
	End   uint64
}

// ======================================================================================

func (r *resource) toProtoBuf() *proto.Resource {

	resource := new(proto.Resource)
	resource.ResourceType = &r.Type
	resource.Type = &r.ValueType
	resource.Scalar = newScalar(r.Value)
	resource.Ranges = newRanges(r.Ranges)
	return resource
}

func (rs resources) toProtoBuf() []*proto.Resource {

	protoResources := make([]*proto.Resource, len(rs))
	for index, res := range rs {
		protoResources[index] = res.toProtoBuf()
	}

	return protoResources

}

// ======================================================================================

func NewCPUResource(cpuValue float64) *resource {
	return &resource{
		Type:      CPU_RESOURCE_TYPE,
		ValueType: proto.Value_SCALAR,
		Value:     &cpuValue,
		Ranges:    nil,
	}
}

func NewMEMResource(memValue float64) *resource {
	return &resource{
		Type:      MEM_RESOURCE_TYPE,
		ValueType: proto.Value_SCALAR,
		Value:     &memValue,
		Ranges:    nil,
	}
}

func NewPortsResource(begin, end uint64) *resource {
	ports := []RangeValue{
		RangeValue{Begin: begin, End: end},
	}
	return &resource{
		Type:      PORTS_RESOURCE_TYPE,
		ValueType: proto.Value_RANGES,
		Value:     nil,
		Ranges:    ports,
	}
}

func newScalar(value *float64) *proto.Value_Scalar {
	return &proto.Value_Scalar{
		Value: value,
	}
}

func newRanges(rangeValues []RangeValue) *proto.Value_Ranges {
	totalRanges := len(rangeValues)

	protoValueRanges := new(proto.Value_Ranges)

	ranges := make([]*proto.Value_Range, totalRanges)

	for index, rng := range rangeValues {
		singleRange := &proto.Value_Range{
			Begin: &rng.Begin,
			End:   &rng.End,
		}
		ranges[index] = singleRange
	}

	protoValueRanges.Range = ranges

	return protoValueRanges
}
