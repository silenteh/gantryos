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

type portMapping struct {
	HostPort      int
	ContainerPort int
	Protocol      string
}

func (p *port) toProtoBuf() *proto.Port {

	port := new(proto.Port)
	port.Name = &p.Name
	port.Number = protobuf.Uint32(uint32(p.Number))
	port.Protocol = &p.Protocol
	return port
}

func (p ports) toProtoBuf() *proto.Ports {

	ports := new(proto.Ports)

	portsProto := make([]*proto.Port, len(p))
	for index, res := range p {
		portsProto[index] = res.toProtoBuf()
	}

	ports.Ports = portsProto

	return ports
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
