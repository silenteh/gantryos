package services

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
)

func slaveListener(masterToSlaveChannel chan *proto.Envelope) {
	for {

		// this blocks until a message is available in the queue
		// message processing is sequential
		// all messages are idempotent
		data := <-masterToSlaveChannel
		fmt.Println("Slave reader channel got data !")

		switch {
		case data.MasterInfo != nil:
		case data.SlaveInfo != nil:
		case data.Heartbeat != nil:
			log.Infoln(data.Heartbeat.GetSlave().GetHostname())
		default:
			log.Errorln("Got an unknown request from the GatryOS slave")
			continue
		}

	}
}
