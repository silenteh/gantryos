package services

// this package sends data to the slaves

import (
	"github.com/silenteh/gantryos/models"
)

//====================================================================

// this method is used for registering with the master
func (ms *masterServer) taskRequest(t models.Task) {
	e := ms.master.RunTask(t)
	ms.writerChannel <- e
}
