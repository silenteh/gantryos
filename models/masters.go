package models

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
)

type master struct {
	Id       string
	Ip       string
	Port     int
	Hostname string
}

func NewMaster(id, ip, hostname string, port int) *master {
	m := new(master)
	m.Id = id
	m.Ip = ip
	m.Hostname = hostname
	m.Port = port
	return m
}

func (m *master) ToProtoBuf() *proto.MasterInfo {

	master := new(proto.MasterInfo)
	master.Id = &m.Id
	master.Ip = &m.Ip
	master.Port = protobuf.Uint32(uint32(m.Port))
	master.Hostname = &m.Hostname

	return master

}
