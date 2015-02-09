package services

import (
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
)

var slaveEnvelopeChannel chan *proto.Envelope
var slaveInstance *gantryTCPClient

func StartSlave(masterIp, masterPort string, channel chan *proto.Envelope) {

	slaveEnvelopeChannel = channel

	slaveInstance = newGantryTCPClient(masterIp, masterPort)
	err := slaveInstance.Connect()
	if err != nil {
		log.Fatalln("Cannot connect to master", masterIp, "on port", masterPort, " => ", err)
	}
	// init the slave
	// it needs to send its information to the master
	if slaveInfo.NewRegistration {
		reRegisterMaster()
	} else {
		joinMaster()
	}
	// start the writers to pull from the queue and write to the master
	go writer(channel)

	// start the listener to detect the master requests
	go slaveListener(channel)

	// the method blocks therefore start a go routine
	go pingMaster()

}

func slaveSendMessage(envelope *proto.Envelope) {
	slaveEnvelopeChannel <- envelope
}

func writer(channel chan *proto.Envelope) {
	for {
		// this blocks until there is data in the channel
		envelope := <-slaveEnvelopeChannel
		if err := slaveInstance.Write(envelope); err != nil {
			// re-queue
			slaveEnvelopeChannel <- envelope
			// disconnect and ignore the error
			slaveInstance.Disconnect()
			//reconnect
			slaveInstance.Connect()
		}

	}
}

func StopSlave() {

	if slaveInstance != nil {
		close(slaveEnvelopeChannel)
		slaveInstance.Disconnect()
	}
	log.Infoln("GantryOS slave stopped.")

}
