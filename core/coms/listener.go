package coms

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
)

// this listener received the proto.Envelope data
// the it detects which sub-field is available
// and therefore which request the client is doing
// this method blocks
func initListener(envelopeChannel chan *proto.Envelope) {
	go listener(envelopeChannel)
}

func listener(envelopeChannel chan *proto.Envelope) {
	for {
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
