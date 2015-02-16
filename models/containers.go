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
	DomainName           string
	WorkingDir           string
	MacAddress           string // Sets the container's Ethernet device's MAC address
	VolumesFrom          string // Mount all volumes from the given container(s)
	Resources            resources
	Volumes              containerVolumes
	PortsMapping         portsMapping
	HealthCheck          healthchecks
	EnvironmentVariables environmentVariables
	CMDs                 arguments // this is to override the execution of the container
	Entrypoint           arguments
	SecurityOptions      arguments // these are labels for AppArmor or SELinux
	OnBuild              arguments
}

func NewContainer(name, image, hostname, domainName, workingDir string, forcePull bool,
	network proto.ContainerInfo_Network, volumes containerVolumes, portsMapping portsMapping,
	envVars environmentVariables) *container {

	container := &container{
		Name:                 name,
		Image:                image,
		ForcePull:            forcePull,
		Network:              network,
		Hostname:             hostname,
		DomainName:           domainName,
		WorkingDir:           workingDir,
		Volumes:              volumes,
		PortsMapping:         portsMapping,
		EnvironmentVariables: envVars,
	}

	return container

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
	cinfo.Cmd = c.CMDs
	cinfo.EntryPoint = c.Entrypoint
	cinfo.WorkingDir = &c.WorkingDir
	cinfo.MacAddress = &c.MacAddress
	cinfo.SecurityOptions = c.SecurityOptions
	cinfo.OnBuild = c.OnBuild
	cinfo.VolumesFrom = &c.VolumesFrom
	return cinfo
}
