package services

import (
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
)

var masterInstance *gantryTCPServer
var mReaderChannel chan *proto.Envelope
var mWriterChannel chan *proto.Envelope

func StartMaster(ip, port string, readerChannel chan *proto.Envelope, writerChannel chan *proto.Envelope) {

	// this is done only so that we can easily close the channel before exiting.
	// TODO: improve it, because it's ugly
	// Solutions: do we really need to care about closing the channel once the process exists ?
	// also we can listen to specific messages which then close the channel
	mReaderChannel = readerChannel

	masterInstance = newGantryTCPServer(ip, port, mReaderChannel, writerChannel)
	masterInstance.StartTCP()

	// start the listener to detect the client calls
	go masterListener(mReaderChannel)

}

func StopMaster() {

	if masterInstance != nil {
		close(mWriterChannel)
		close(mReaderChannel)
		masterInstance.Stop()
	}
	log.Infoln("GantryOS master stopped.")

}
