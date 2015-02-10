package services

import (
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
)

var sReaderChannel chan *proto.Envelope
var sWriterChannel chan *proto.Envelope

var slaveInstance *gantryTCPClient

func StartSlave(masterIp, masterPort string, readerChannel chan *proto.Envelope, writerChannel chan *proto.Envelope) {

	sReaderChannel = readerChannel
	sWriterChannel = writerChannel

	slaveInstance = newGantryTCPClient(masterIp, masterPort)
	err := slaveInstance.Connect()
	if err != nil {
		log.Fatalln("Cannot connect to master", masterIp, "on port", masterPort, " => ", err)
	}

	// init Slave: the method is in the slaveactions.go file
	// should be improved
	initSlave()

	// start the writers to pull from the queue and write to the master
	go writer(sWriterChannel)

	// start the listener to detect the master requests
	go slaveListener(readerChannel)

	// the method blocks therefore start a go routine
	go pingMaster()
}

func slaveSendMessage(envelope *proto.Envelope) {
	sWriterChannel <- envelope
}

func writer(channel chan *proto.Envelope) {
	for {
		// this blocks until there is data in the channel
		envelope := <-sWriterChannel
		if err := slaveInstance.Write(envelope); err != nil {
			log.Errorln(err)
			// re-queue
			sWriterChannel <- envelope
			// disconnect and ignore the error
			slaveInstance.Disconnect()
			//reconnect
			slaveInstance.Connect()

			continue
		}
	}
}

func StopSlave() {

	if slaveInstance != nil {
		close(sWriterChannel)
		close(sReaderChannel)
		slaveInstance.Disconnect()
	}
	log.Infoln("GantryOS slave stopped.")

}
