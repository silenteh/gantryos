package services

// this package sends data to the slaves

import (
	"github.com/silenteh/gantryos/models"
)

//====================================================================

// this method is used for registering with the master
func (ms *masterServer) taskRequest(task *models.Task) {
	e := ms.master.RunTask(task)
	ms.writerChannel <- e
}
