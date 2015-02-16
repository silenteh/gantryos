package models

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
	"time"
)

type TaskStatus struct {
	Id        string
	TaskId    string
	TaskState proto.TaskState
	Message   string
	Slave     *Slave
	Timestamp time.Time
	Healthy   bool // this is based on the health checks
}

// TODO: implement healthy and message
func NewTaskStatus(id, taskId, message string, taskState proto.TaskState, slave *Slave) *TaskStatus {

	t := TaskStatus{
		Id:        id,
		TaskId:    taskId,
		TaskState: taskState,
		Message:   message,
		Timestamp: time.Now().UTC(),
		Slave:     slave,
		Healthy:   true,
	}

	return &t

}

func NewTaskStatusNoSlave(id, taskId, message string, taskState proto.TaskState) *TaskStatus {

	t := TaskStatus{
		Id:        id,
		TaskId:    taskId,
		TaskState: taskState,
		Message:   message,
		Timestamp: time.Now().UTC(),
		Healthy:   true,
	}

	return &t

}

func (t *TaskStatus) ToProtoBuf() *proto.TaskStatus {

	//e := new(proto.Envelope)

	taskState := new(proto.TaskStatus)
	taskState.GantryTaskId = &t.Id
	taskState.TaskId = &t.TaskId
	taskState.Healthy = protobuf.Bool(t.Healthy)
	taskState.Timestamp = protobuf.Float64(float64(t.Timestamp.Unix()))
	taskState.Slave = t.Slave.ToProtoBuf()
	taskState.State = &t.TaskState
	//e.TaskStatus = taskState

	return taskState

}
