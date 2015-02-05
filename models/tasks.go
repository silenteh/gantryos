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

func (t *Task) ToProtoBuf() *proto.TaskInfo {

	taskInfo := new(proto.TaskInfo)
	taskInfo.TaskId = &t.Id
	taskInfo.TaskName = &t.Name
	taskInfo.TaskVersion = &t.Version
	taskInfo.Slave = t.Slave.ToProtoBuf()
	taskInfo.Resources = t.Resources.ToProtoBuf()
	taskInfo.Command = t.Command.ToProtoBuf()

	// container
	taskInfo.Container = t.Container.ToProtoBuf()

	taskInfo.Discovery = t.Discovery.ToProtoBuf()
	taskInfo.Labels = t.Labels.ToProtoBuf()

	return taskInfo
}
