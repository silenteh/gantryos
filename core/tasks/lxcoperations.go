package tasks

import (
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/models"
)

type lxcTask struct {
	info             *proto.TaskInfo
	taskStatusEvents chan *models.TaskStatus
}

// ==========================================================================================
// ==== LXC

func (t lxcTask) Start(taskInfo *proto.TaskInfo) (string, error) {
	return "", nil
}

func (t lxcTask) Stop(taskId string) error {
	return nil
}

func (t lxcTask) Status(taskId string) (proto.TaskState, error) {

	return -1, nil
}

func (t lxcTask) CleanContainers() error {
	return nil
}

func (t lxcTask) CleanImages() error {
	return nil
}

func (t lxcTask) Monitor() {

}

func (t lxcTask) StopService() {
	close(t.taskStatusEvents)
}
