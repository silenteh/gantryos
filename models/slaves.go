package models

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
	sysresources "github.com/silenteh/gantryos/core/resources"
	"strconv"
)

type Slave struct {
	Id              string
	Ip              string
	Port            int
	Hostname        string
	Resources       resources // resources available on the slave
	Checkpoint      bool
	NewRegistration bool
}

func NewSlave(id, ip, hostname string, port int, checkpoint, newRegistration bool) *Slave {

	slave := new(Slave)
	slave.Id = id
	slave.Ip = ip
	slave.Hostname = hostname
	slave.Port = port
	slave.Checkpoint = checkpoint
	slave.NewRegistration = newRegistration

	// get the resources
	cpu := NewCPUResource(sysresources.GetTotalCPUsCount())
	mem := NewMEMResource(sysresources.GetTotalRam())
	ports := NewPortsResource(58000, 59000)

	res := make([]*resource, 3)
	res[0] = cpu
	res[1] = mem
	res[2] = ports

	slave.Resources = res
	// ===================================

	return slave
}

func (s *Slave) ToProtoBuf() *proto.SlaveInfo {

	slave := new(proto.SlaveInfo)
	slave.Id = &s.Id
	slave.Ip = &s.Ip
	slave.Port = protobuf.Uint32(uint32(s.Port))
	slave.Hostname = &s.Hostname
	slave.Checkpoint = protobuf.Bool(s.Checkpoint)
	slave.Resources = s.Resources.ToProtoBuf()

	return slave
}

func (s *Slave) RegisterSlaveMessage() *proto.Envelope {

	e := newSlaveEnvelope(s)
	m := new(proto.RegisterSlaveMessage)
	m.Slave = s.ToProtoBuf()
	e.RegisterSlave = m
	return e
}

func (s *Slave) ReRegisterSlaveMessage() *proto.Envelope {
	e := newSlaveEnvelope(s)
	m := new(proto.ReregisterSlaveMessage)
	m.Slave = s.ToProtoBuf()
	e.ReRegisterSlave = m
	return e
}

func (s *Slave) GetPortString() string {
	return strconv.Itoa(s.Port)
}

func (s *Slave) NewHeartBeatMessage() *proto.Envelope {

	e := newSlaveEnvelope(s)
	hb := new(proto.HeartbeatMessage)
	hb.Slave = s.ToProtoBuf()
	e.Heartbeat = hb
	return e
}

func (s *Slave) NewTaskStatusMessage(taskStatus *TaskStatus) *proto.Envelope {

	e := newSlaveEnvelope(s)
	taskStatus.Slave = s
	su := new(proto.TaskStatusMessage)
	su.TaskStatus = taskStatus.ToProtoBuf()
	e.TaskStatusMessage = su
	return e
}
