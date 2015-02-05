package masters

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/coms"
	"github.com/silenteh/gantryos/core/proto"
)

var envelopeChannel = make(chan *proto.Envelope, 1024)
var master *coms.GantryTCPServer

func Start(ip, port string) {

	Master = coms.NewGantryTCPServer(ip, port, envelopeChannel)
	Master.StartTCP()

	// start the listener to detect the client requests
	go listener()

}

func Stop() {

	if Master != nil {
		close(envelopeChannel)
		Master.Stop()
	}
	log.Infoln("GantryOS master stopped.")

}

// this listener received the proto.Envelope data
// the it detects which sub-field is available
// and therefore which request the client is doing
// this method blocks
func listener() {

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
