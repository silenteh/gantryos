package models

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/silenteh/gantryos/core/proto"
)

type Task struct {
	Id        string
	TaskId    string
	Name      string
	Version   string
	Slave     *Slave
	Resources *resources // resources required by the task
	Command   *command   // this is in case we want to execute a single exec
	Container *container // this is the definition of the container properties
	Discovery *discovery // information about the task/service for discovery mechanism
	Labels    *labels
}

func NewTask(name, version string, slave *Slave, res *resources, cmd *command, containerObject *container, serviceDiscovery *discovery, labelsObject *labels) *Task {
	t := new(Task)
	t.Id = uuid.NewRandom().String()
	t.Name = name
	t.Version = version
	t.Slave = slave
	t.Resources = res
	t.Command = cmd
	t.Container = containerObject
	t.Discovery = serviceDiscovery
	t.Labels = labelsObject

	return t
}

func (t *Task) ToProtoBuf() *proto.TaskInfo {

	taskInfo := new(proto.TaskInfo)
	taskInfo.GantryTaskId = &t.Id
	taskInfo.TaskId = &t.TaskId
	taskInfo.TaskName = &t.Name
	taskInfo.TaskVersion = &t.Version
	taskInfo.Slave = t.Slave.ToProtoBuf()

	taskInfo.Resources = nil
	if t.Resources != nil {
		taskInfo.Resources = t.Resources.ToProtoBuf()
	}

	taskInfo.Command = nil
	if t.Command != nil {
		taskInfo.Command = t.Command.ToProtoBuf()
	}

	taskInfo.Discovery = nil
	if t.Discovery != nil {
		taskInfo.Discovery = t.Discovery.ToProtoBuf()
	}

	taskInfo.Labels = nil
	if t.Labels != nil {
		taskInfo.Labels = t.Labels.ToProtoBuf()
	}

	// container
	taskInfo.Container = t.Container.ToProtoBuf()

	return taskInfo
}
