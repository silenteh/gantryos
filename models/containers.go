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
	DockerNetwork        proto.ContainerInfo_Network // BRIDGED_NETWORK || HOST_NETWORK || ...
	CPUs                 float64
	Memory               float64
	Hostname             string
	Volumes              containerVolumes
	PortMapping          []*portMapping
	HealthCheck          healthchecks
	EnvironmentVariables environmentVariables
	CMDs                 parameters // this is to override the execution of the container
}
