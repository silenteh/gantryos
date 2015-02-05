package models

import "github.com/silenteh/gantryos/core/proto"

type Task struct {
	Id        string
	Name      string
	Version   string
	Slave     Slave
	Resources resources  // resources required by the task
	Command   *command   // this is in case we want to execute a single exec
	Container *container // this is the definition of the container properties
	Discovery *discovery // information about the task/service for discovery mechanism
	Labels    labels
}

func (t *Task) toProtoBufMessage() *proto.TaskInfo {

	taskInfo := new(proto.TaskInfo)
	taskInfo.TaskId = &t.Id
	taskInfo.TaskName = &t.Name
	taskInfo.TaskVersion = &t.Version
	taskInfo.Slave = t.Slave.toProtoBuf()
	taskInfo.Resources = t.Resources.toProtoBuf()
	taskInfo.Command = t.Command.toProtoBuf()

	// container

	taskInfo.Discovery = t.Discovery.toProfoBuf()
	taskInfo.Labels = t.Labels.toProtoBuf()

	return taskInfo
}
