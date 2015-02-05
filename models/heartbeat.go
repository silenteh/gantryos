package models

import "github.com/silenteh/gantryos/core/proto"

func NewHeartBeat(s *slave) *proto.HeartbeatMessage {
	hb := new(proto.HeartbeatMessage)
	hb.Slave = s.ToProtoBuf()
	return hb
}
