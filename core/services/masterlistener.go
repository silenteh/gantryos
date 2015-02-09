package services

import (
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
)

// this listener received the proto.Envelope data
// the it detects which sub-field is available
// and therefore which request the client is doing
// this method blocks

func masterListener(envelopeChannel chan *proto.Envelope) {
	for {
		data := <-envelopeChannel

		switch {
		case data.MasterInfo != nil:
		case data.RegisterSlave != nil:
			slave := data.RegisterSlave.GetSlave()
			log.Infoln(slave.Hostname, "registered as a slave with id", slave.Id)
		case data.ReRegisterSlave != nil:
			slave := data.ReRegisterSlave.GetSlave()
			log.Infoln(slave.Hostname, "registered as a slave with id", slave.Id)
		case data.SlaveInfo != nil:
		case data.Heartbeat != nil:
			slave := data.Heartbeat.GetSlave()
			log.Infoln("Slave id", slave.Id, "with hostname", slave.Hostname, "sent an heartbeat.")
		default:
			log.Errorln("Got an unknown request from the GatryOS slave")
			continue
		}

	}
}

// ********************** SLAVE ***********************************************************
