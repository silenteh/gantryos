package main

import (
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/core/services"
	"os"
	"os/signal"
	"time"
)

var channel = make(chan *proto.Envelope, 1024)

func main() {

	// =============================================================
	// Exit channel
	channelCtrlC := make(chan os.Signal, 1)
	signal.Notify(channelCtrlC, os.Interrupt, os.Kill)
	// =============================================================

	// start the master
	masterIp := "127.0.0.1"
	masterPort := "6060"
	services.StartMaster(masterIp, masterPort, channel)
	log.Infoln("Master started at", masterIp, "on port", masterPort)

	// wait for binding
	time.Sleep(1 * time.Second)

	// start the slave
	services.StartSlave(masterIp, masterPort, channel)
	log.Infoln("Slave started")

	// first flush
	log.Flush()

	// wait for a signal to exit the app
	<-channelCtrlC

	// final flush
	log.Flush()

}
