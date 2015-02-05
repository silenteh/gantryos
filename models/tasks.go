package models

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/silenteh/gantryos/core/proto"
)

type task struct {
	Id        string
	Name      string
	Version   string
	Slave     slave
	Resources resources  // resources required by the task
	Command   *command   // this is in case we want to execute a single exec
	Container *container // this is the definition of the container properties
	Discovery *discovery // information about the task/service for discovery mechanism
	Labels    labels
}

func NewTask(name, version string, s slave, res resources, cmd command, cont container, disc discovery, lbls labels) *task {
	t := new(task)
	t.Id = uuid.NewRandom().String()
	t.Name = name
	t.Version = version
	t.Slave = s
	t.Resources = res
	t.Command = &cmd
	t.Container = &cont
	t.Discovery = &disc
	t.Labels = lbls

	return t
}

func (t *task) ToProtoBuf() *proto.TaskInfo {

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
