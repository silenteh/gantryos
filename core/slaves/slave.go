package slaves

import (
	//"fmt"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/coms"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/models"
)

var envelopeChannel = make(chan *proto.Envelope, 1024)
var slave *coms.GantryTCPClient

func Start(masterIp, masterPort string) {

	slave = coms.NewGantryTCPClient(masterIp, masterPort)

	// init the slave
	// it needs to send its information to the master

	// start the writers to pull from the queue and write to the master
	go writer()

}

func writer() {
	for {
		// this blocks until there is data in the channel
		envelope := <-envelopeChannel
		if err := slave.Write(envelope); err != nil {
			// re-queue
			envelopeChannel <- envelope
			// disconnect and ignore the error
			slave.Disconnect()
			//reconnect
			slave.Connect()
		}

	}
}

func Stop() {

	if slave != nil {
		close(envelopeChannel)
		slave.Disconnect()
	}
	log.Infoln("GantryOS slave stopped.")

}

//====================================================================

func initSlave() {

	slave := models.NewSlave(id, ip, hostname, port, checkpoint, res)
}
