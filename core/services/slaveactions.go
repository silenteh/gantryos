package services

import (
	"github.com/silenteh/gantryos/models"
)

//====================================================================

// this method is used for registering with the master
func (s *slaveServer) joinMaster() {
	e := s.slave.RegisterSlaveMessage()
	s.writerChannel <- e
}

// this method is used to re-register with the master
func (s *slaveServer) reRegisterMaster() {
	e := s.slave.ReRegisterSlaveMessage()
	s.writerChannel <- e
}

// this method is used to disconect from the master
func disconnectMaster() {

}

// this method is used to tell the master we are still alive
// IT BLOCKS !
func (s *slaveServer) pingMaster() {
	e := s.slave.NewHeartBeatMessage()
	s.writerChannel <- e
}

// this method is used to offer resources to the master
func (s *slaveServer) resourceOffer() {

}

// this method is used to tell the master a task has changed its state
func (s *slaveServer) taskStateChange(taskStatus *models.TaskStatus) {

	e := s.slave.NewTaskStatusMessage(taskStatus)
	s.writerChannel <- e

}

// this method is used to answer the master about the inquiry for a specific task
func (s *slaveServer) reconciliateTask(taskId string) {

}

// this method is used to answer the master about the inquiry for a set of tasks
func (s *slaveServer) reconciliateTasks(tasks []string) {

}
