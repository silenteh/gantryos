package tasks

import (
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/models"
)

type TaskInterface interface {
	Start(taskInfo *proto.TaskInfo) (string, error)    // used to start a container starting (it does all the operations, like pull and start)
	Stop(containerId string, removeVolumes bool) error // stops the container and removes the stopped container
	Status(containerId string) error                   // get the status of the container
	CleanContainers() error                            // cleans up the remained containers which are not running
	CleanImages(image string) error                    // cleans up the images which are not used and are not being pulled (check if there is a pull running)
	GetEventChannel() chan *models.TaskStatus          // get the events channel
	StopService()                                      // used to close the channels
}
