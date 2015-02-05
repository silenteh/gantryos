package models

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
)

type Slave struct {
	Id         string
	Ip         string
	Port       int
	Hostname   string
	Resources  resources // resources available on the slave
	Checkpoint bool
}

func (s *Slave) toProtoBuf() *proto.SlaveInfo {

	slave := new(proto.SlaveInfo)
	slave.Id = &s.Id
	slave.Ip = &s.Ip
	slave.Port = protobuf.Uint32(uint32(s.Port))
	slave.Hostname = &s.Hostname
	slave.Checkpoint = protobuf.Bool(s.Checkpoint)
	slave.Resources = s.Resources.toProtoBuf()

	return slave
}
