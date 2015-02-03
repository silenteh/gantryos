package protoutils

import (
	protobuf "github.com/gogo/protobuf/proto"
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
		protoRes.Type = resource.ValueType
		protoRes.ResourceType = resource.Type
		switch resource.Type {
		case models.CPU_RESOURCE_TYPE:
			protoRes.Scalar = protobuf.Float64(resource.Value)
		case models.MEM_RESOURCE_TYPE:
			protoRes.Scalar = protobuf.Float64(resource.Value)
		case models.GPU_RESOURCE_TYPE:
			protoRes.Scalar = protobuf.Float64(resource.Value)
		case models.PORTS_RESOURCE_TYPE:
			protoRes.Ranges = resource.Value
		}
	}

}
