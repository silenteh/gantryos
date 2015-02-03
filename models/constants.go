package models

import "github.com/silenteh/gantryos/core/proto"

const (
	CPU_RESOURCE_TYPE   = proto.ResourceType_CPU
	MEM_RESOURCE_TYPE   = proto.ResourceType_MEMORY
	PORTS_RESOURCE_TYPE = proto.ResourceType_PORTS
	DISK_RESOURCE_TYPE  = proto.ResourceType_DISK
	NET_RESOURCE_TYPE   = proto.ResourceType_BANDWIDTH
	GPU_RESOURCE_TYPE   = proto.ResourceType_GPU
)
