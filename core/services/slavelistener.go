package services

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
)

func slaveListener(envelopeChannel chan *proto.Envelope) {
	for {

		// this blocks until a message is available in the queue
		// message processing is sequential
		// all messages are idempotent
		data := <-envelopeChannel

		switch {
		case data.MasterInfo != nil:
		case data.SlaveInfo != nil:
		case data.Heartbeat != nil:
			fmt.Println("HEARTBEAT RECEIVED !")
			log.Infoln(data.Heartbeat.GetSlave().Hostname)
		default:
			log.Errorln("Got an unknown request from the GatryOS slave")
			continue
		}

	}
}
