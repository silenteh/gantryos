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
	//total := 0

	for {
		envlope := <-envelopeChannel

		// dereference
		data := *envlope

		switch {
		case data.MasterInfo != nil:
			continue
		case data.RegisterSlave != nil:
			slave := data.RegisterSlave.GetSlave()
			log.Infoln(slave.GetHostname(), "registered as a slave with id", slave.GetId())
			continue
		case data.ReRegisterSlave != nil:
			slave := data.ReRegisterSlave.GetSlave()
			log.Infoln(slave.GetHostname(), "re-registered as a slave with id", slave.GetId())
			continue
		case data.SlaveInfo != nil:
			continue
		case data.Heartbeat != nil:
			slave := data.Heartbeat.GetSlave()
			log.Infoln("Slave id", slave.GetId(), "with hostname", slave.GetHostname(), "sent an heartbeat.")
			continue
		default:
			log.Errorln("Got an unknown request from the GatryOS slave")
			continue
		}
	}
}

// ********************** SLAVE ***********************************************************
