package models

import "github.com/silenteh/gantryos/core/proto"

const (
	DISCOVERY_CONSUL = 1
	DISCOVERY_ETCD   = 2
)

type discovery struct {
	Type        int
	Ports       ports
	Name        string
	Environment string
	Location    string
	Version     string
	Labels      labels
}

func NewDiscovery(discoveryType int, prts ports, name, env, location, version string, lbls labels) *discovery {

	return &discovery{
		Type:        discoveryType,
		Ports:       prts,
		Name:        name,
		Environment: env,
		Location:    location,
		Version:     version,
		Labels:      lbls,
	}

}

// TODO: implement the Discovery TYPE
// we need to create an ENUM in PROTOBUF
func (d discovery) ToProtoBuf() *proto.DiscoveryInfo {
	dProto := new(proto.DiscoveryInfo)
	dProto.Name = &d.Name
	dProto.Environment = &d.Environment
	dProto.Location = &d.Location
	dProto.Version = &d.Version
	dProto.Ports = d.Ports.ToProtoBuf()
	dProto.Labels = d.Labels.ToProtoBuf()
	return dProto
}
