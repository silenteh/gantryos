package services

import (
	"github.com/silenteh/gantryos/config"
	"github.com/silenteh/gantryos/core/resources"
	"github.com/silenteh/gantryos/models"
	"time"
)

//====================================================================

var slaveInfo *models.Slave

func initSlave() {
	// Port
	port := 6061
	if config.GantryOSConfig.Slave.Port != 0 {
		port = config.GantryOSConfig.Slave.Port
	}
	// ==============================================

	// IP
	ip := "127.0.0.1"
	if config.GantryOSConfig.Slave.IP != "" {
		ip = config.GantryOSConfig.Slave.IP
	}
	// ==============================================

	// Hostname
	hostname := resources.GetHostname()
	// ==============================================

	// Slave ID
	slaveId := config.GantryOSSlaveId

	slaveInfo = models.NewSlave(slaveId.Id, ip, hostname, port, config.GantryOSConfig.Slave.Checkpoint, slaveId.Registered)

	// init the slave
	// it needs to send its information to the master
	if slaveId.Registered {
		reRegisterMaster()
	} else {
		joinMaster()
	}
}

// this method is used for registering with the master
func joinMaster() {
	m := slaveInfo.RegisterSlaveMessage()
	e := models.NewEnvelope()
	e.RegisterSlave = m
	slaveSendMessage(e)
}

// this method is used to re-register with the master
func reRegisterMaster() {
	m := slaveInfo.ReRegisterSlaveMessage()
	e := models.NewEnvelope()
	e.ReRegisterSlave = m
	slaveSendMessage(e)
}

// this method is used to disconect from the master
func disconnectMaster() {

}

// this method is used to tell the master we are still alive
// IT BLOCKS !
func pingMaster() {
	m := models.NewHeartBeat(slaveInfo)
	e := models.NewEnvelope()
	e.Heartbeat = m
	slaveSendMessage(e)
}

func startHeartBeats() {
	for {
		pingMaster()
		time.Sleep(5 * time.Second)
	}
}

// this method is used to offer resources to the master
func resourceOffer() {

}

// this method is used to tell the master a task has changed its state
func taskStateChange() {

}

// this method is used to answer the master about the inquiry for a specific task
func reconciliateTask(taskId string) {

}

// this method is used to answer the master about the inquiry for a set of tasks
func reconciliateTasks(tasks []string) {

}
