package main

import (
	"flag"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/core/services"
	"github.com/silenteh/gantryos/core/state"
	"os"
	"os/signal"
	"time"
)

// writes the requests to the slaves
var masterWriterChannel = make(chan *proto.Envelope, 1024)

// writes the requests to the slaves
var masterReaderChannel = make(chan *proto.Envelope, 1024)

// ==========================================================

// writes the requests to the master
var slaveWriterChannel = make(chan *proto.Envelope, 1024)

// receives the requests from the master
var slaveReaderChannel = make(chan *proto.Envelope, 1024)

func init() {
	flag.Parse()
}

func main() {

	// =============================================================
	// Exit channel
	channelCtrlC := make(chan os.Signal, 1)
	signal.Notify(channelCtrlC, os.Interrupt, os.Kill)
	// =============================================================

	stateDb, err := state.InitSlaveDB("./gantryos.db")
	if err != nil {
		log.Fatal(err)
	}

	defer stateDb.Close()

	// start the master
	masterIp := "127.0.0.1"
	masterPort := "6060"
	services.StartMaster(masterIp, masterPort, masterReaderChannel, masterWriterChannel)
	log.Infoln("Master started at", masterIp, "on port", masterPort)

	// wait for binding
	time.Sleep(1 * time.Second)

	// start the slave
	services.StartSlave(masterIp, masterPort, slaveReaderChannel, slaveWriterChannel, stateDb)
	log.Infoln("Slave started")

	// first flush
	log.Flush()

	// wait for a signal to exit the app
	<-channelCtrlC

	// final flush
	log.Flush()

}
