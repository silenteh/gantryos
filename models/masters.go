package models

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
	"strconv"
)

type Master struct {
	Id       string
	Ip       string
	Port     int
	Hostname string
}

func NewMaster(id, ip, hostname string, port int) *Master {
	m := new(Master)
	m.Id = id
	m.Ip = ip
	m.Hostname = hostname
	m.Port = port
	return m
}

func (m *Master) ToProtoBuf() *proto.MasterInfo {

	master := new(proto.MasterInfo)
	master.Id = &m.Id
	master.Ip = &m.Ip
	master.Port = protobuf.Uint32(uint32(m.Port))
	master.Hostname = &m.Hostname

	return master

}

func (m *Master) GetPortString() string {
	return strconv.Itoa(m.Port)
}
