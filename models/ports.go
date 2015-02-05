package models

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
)

type ports []*port

type port struct {
	Number   int
	Name     string
	Protocol string // * means could be TCP, UDP, both TCP_UDP or * or empty
}

type portsMapping []*portMapping

type portMapping struct {
	HostPort      int
	ContainerPort int
	Protocol      string
}

func (p *port) ToProtoBuf() *proto.Port {

	port := new(proto.Port)
	port.Name = &p.Name
	port.Number = protobuf.Uint32(uint32(p.Number))
	port.Protocol = &p.Protocol
	return port
}

func (p ports) ToProtoBuf() *proto.Ports {

	ports := new(proto.Ports)

	portsProto := make([]*proto.Port, len(p))
	for index, res := range p {
		portsProto[index] = res.ToProtoBuf()
	}

	ports.Ports = portsProto

	return ports
}

func (pm *portMapping) ToProtoBuf() *proto.ContainerInfo_PortMapping {
	portMapping := new(proto.ContainerInfo_PortMapping)
	portMapping.ContainerPort = protobuf.Uint32(uint32(pm.ContainerPort))
	portMapping.HostPort = protobuf.Uint32(uint32(pm.HostPort))
	portMapping.Protocol = &pm.Protocol
	return portMapping
}

func (pms portsMapping) ToProtoBuf() []*proto.ContainerInfo_PortMapping {

	portsProto := make([]*proto.ContainerInfo_PortMapping, len(pms))
	for index, res := range pms {
		portsProto[index] = res.ToProtoBuf()
	}

	return portsProto
}

func NewPort(number int, name, proto string) *port {
	if proto == "" || proto == "*" {
		proto = "TCP_UDP" // means both
	}

	return &port{
		Number:   number,
		Name:     name,
		Protocol: proto,
	}
}

func NewPortMapping(hostPort, containerPort int, proto string) *portMapping {
	if proto == "" || proto == "*" {
		proto = "TCP_UDP" // means both
	}

	return &portMapping{
		HostPort:      hostPort,
		ContainerPort: containerPort,
		Protocol:      proto,
	}
}
