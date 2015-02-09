package services

import (
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
)

var masterInstance *gantryTCPServer
var masterEnvelopeChannel chan *proto.Envelope

func StartMaster(ip, port string, channel chan *proto.Envelope) {

	// this is done only so that we can easily close the channel before exiting.
	// TODO: improve it, because it's ugly
	// Solutions: do we really need to care about closing the channel once the process exists ?
	// also we can listen to specific messages which then close the channel
	masterEnvelopeChannel = channel

	masterInstance = newGantryTCPServer(ip, port, channel)
	masterInstance.StartTCP()

	// start the listener to detect the client calls
	go masterListener(channel)

}

func StopMaster() {

	if masterInstance != nil {
		close(masterEnvelopeChannel)
		masterInstance.Stop()
	}
	log.Infoln("GantryOS master stopped.")

}
