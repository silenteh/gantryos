package protoutils

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/models"
)

// TODO: check input parameters
func NewSlaveRegistrationRequest(id, ip, hostname string, port uint32, checkpoint bool, resources []models.Resource) *proto.Envelope {

	// create the base objects
	envelope := new(proto.Envelope)
	regRequest := new(proto.RegisterSlaveMessage)
	slaveInfo := new(proto.SlaveInfo)

	// assign properties to the slaveInfo
	slaveInfo.Id = protobuf.String(id)
	slaveInfo.Ip = protobuf.String(ip)
	slaveInfo.Port = protobuf.Uint32(port)
	slaveInfo.Hostname = protobuf.String(hostname)
	slaveInfo.Checkpoint = protobuf.Bool(checkpoint)
	slaveInfo.Resources = makeProtoResource(resources)

	// assign the slave info to the request
	regRequest.Slave = slaveInfo

	// wrap it to the envelope
	envelope.RegisterSlave = regRequest

	return envelope

}

func NewSlaveRe_RegistrationRequest(id, ip, hostname string, port uint32, checkpoint bool, resources []models.Resource, tasks []models.Task) *proto.Envelope {

	// create the base objects
	envelope := new(proto.Envelope)
	regRequest := new(proto.ReregisterSlaveMessage)
	slaveInfo := new(proto.SlaveInfo)

	// assign properties to the slaveInfo
	slaveInfo.Id = protobuf.String(id)
	slaveInfo.Ip = protobuf.String(ip)
	slaveInfo.Port = protobuf.Uint32(port)
	slaveInfo.Hostname = protobuf.String(hostname)
	slaveInfo.Checkpoint = protobuf.Bool(checkpoint)
	slaveInfo.Resources = makeProtoResource(resources)

	// assign the slave info to the request
	regRequest.Slave = slaveInfo
	//regRequest.Tasks

	// wrap it to the envelope
	envelope.ReRegisterSlave = regRequest

	return envelope

}
