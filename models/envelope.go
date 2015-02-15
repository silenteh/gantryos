package models

import "github.com/silenteh/gantryos/core/proto"

func newEnvelope() *proto.Envelope {
	e := new(proto.Envelope)
	return e
}

func newSlaveEnvelope(s *Slave) *proto.Envelope {
	e := new(proto.Envelope)
	e.SenderId = &s.Id
	return e
}

func newMasterEnvelope(m *Master) *proto.Envelope {
	e := new(proto.Envelope)
	e.SenderId = &m.Id
	return e
}
