package models

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
)

type slave struct {
	Id         string
	Ip         string
	Port       int
	Hostname   string
	Resources  resources // resources available on the slave
	Checkpoint bool
}

func NewSlave(id, ip, hostname string, port int, checkpoint bool, res resources) *slave {

	slave := new(slave)
	slave.Id = id
	slave.Ip = ip
	slave.Hostname = hostname
	slave.Port = port
	slave.Checkpoint = checkpoint
	slave.Resources = res

	return slave

}

func (s *slave) ToProtoBuf() *proto.SlaveInfo {

	slave := new(proto.SlaveInfo)
	slave.Id = &s.Id
	slave.Ip = &s.Ip
	slave.Port = protobuf.Uint32(uint32(s.Port))
	slave.Hostname = &s.Hostname
	slave.Checkpoint = protobuf.Bool(s.Checkpoint)
	slave.Resources = s.Resources.ToProtoBuf()

	return slave
}
