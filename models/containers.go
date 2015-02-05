package models

import "github.com/silenteh/gantryos/core/proto"

const (
	NONE_NETWORK    = proto.ContainerInfo_NONE
	HOST_NETWORK    = proto.ContainerInfo_HOST
	BRIDGED_NETWORK = proto.ContainerInfo_BRIDGE
	VLAN_NETWORK    = proto.ContainerInfo_VLAN
	GRE_NETWORK     = proto.ContainerInfo_GRE
)

type container struct {
	Name                 string
	Image                string
	ForcePull            bool
	Network              proto.ContainerInfo_Network // BRIDGED_NETWORK || HOST_NETWORK || ...
	Hostname             string
	Resources            resources
	Volumes              containerVolumes
	PortsMapping         portsMapping
	HealthCheck          healthchecks
	EnvironmentVariables environmentVariables
	CMDs                 parameters // this is to override the execution of the container
}

func (c *container) ToProtoBuf() *proto.ContainerInfo {

	cinfo := new(proto.ContainerInfo)
	cinfo.Image = &c.Image
	cinfo.ForcePullImage = &c.ForcePull
	cinfo.Network = &c.Network
	cinfo.Hostname = &c.Hostname
	cinfo.Volumes = c.Volumes.ToProtoBuf()
	cinfo.PortMappings = c.PortsMapping.ToProtoBuf()
	cinfo.Environments = c.EnvironmentVariables.ToProtoBuf()
	cinfo.Parameters = c.CMDs.ToProtoBuf()
	return cinfo
}
