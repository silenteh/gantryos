package protoutils

import (
	//protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/models"
)

func makeProtoResource(resources []models.Resource) []*proto.Resource {

	total := len(resources)
	if total == 0 {
		return nil
	}

	protoResources := make([]*proto.Resource, total)

	for index, resource := range resources {
		protoRes := proto.Resource{}
		protoRes.Type = &resource.ValueType
		protoRes.ResourceType = &resource.Type
		switch resource.Type {
		case models.CPU_RESOURCE_TYPE:
			protoRes.Scalar = newScalarValue(resource.Value) //proto.Value_Scalar{protobuf.Float64(resource.Value)}
		case models.MEM_RESOURCE_TYPE:
			protoRes.Scalar = newScalarValue(resource.Value)
		case models.GPU_RESOURCE_TYPE:
			protoRes.Scalar = newScalarValue(resource.Value)
		case models.PORTS_RESOURCE_TYPE:
			protoRes.Ranges = newRangeValues(resource.Ranges) //proto.Value_Ranges{Range: NewRangeValue(resource.Ranges, end)}
		}
		protoResources[index] = &protoRes
	}

	return protoResources
}

func newScalarValue(value *float64) *proto.Value_Scalar {
	scalar := new(proto.Value_Scalar)
	scalar.Value = value
	return scalar
}

func newRangeValues(rangeValues []models.RangeValue) *proto.Value_Ranges {
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
