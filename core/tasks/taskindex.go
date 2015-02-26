package tasks

import (
	log "github.com/Sirupsen/logrus"
	"github.com/silenteh/gantryos/core/state"
)

type taskIndex struct {
	task2ContainerBucket string
	container2TaskBucket string
	taskToContainerIndex map[string]*string
	containerToTaskIndex map[string]*string
	state                state.StateDB
}

// this is used on  each slave
// here we handle the mapping between gantry task IDs and containers task IDs

// this is the index that we need to lookup to map against gantry tasks ids and container ids
// this map needs to be persisted and should survive to restarts or crashes
//var taskToContainerIndex = make(map[string]*string)
//var containerToTaskIndex = make(map[string]*string)

func NewTaskIndex(task2ContainerBucket, container2TaskBucket string, stateDb state.StateDB) taskIndex {
	ti := taskIndex{}
	ti.container2TaskBucket = container2TaskBucket
	ti.task2ContainerBucket = task2ContainerBucket
	ti.state = stateDb
	ti.taskToContainerIndex = stateDb.GetAllKeyValues(task2ContainerBucket) //make(map[string]*string)
	ti.containerToTaskIndex = stateDb.GetAllKeyValues(container2TaskBucket)

	// pre-load all the data in memory

	return ti
}

// stores the container id according to the gantry task id
func (ti taskIndex) AddTaskId(containerId, gantryId string) {
	// store in memory
	ti.containerToTaskIndex[containerId] = &gantryId
	ti.taskToContainerIndex[gantryId] = &containerId

	// store in DB
	ti.state.Set(ti.container2TaskBucket, containerId, gantryId)
	ti.state.Set(ti.task2ContainerBucket, gantryId, containerId)
	log.Infoln("Added container id and task id to storage")
}

func (ti taskIndex) RemoveTaskId(containerId, gantryId string) {
	ti.containerToTaskIndex[containerId] = nil
	ti.taskToContainerIndex[gantryId] = nil

	// remove from db
	ti.state.Delete(ti.container2TaskBucket, containerId)
	ti.state.Delete(ti.task2ContainerBucket, gantryId)

	log.Infoln("Removed container id and task id to storage")
}

// gets the container id from the gantry task id
func (ti taskIndex) GetTaskId(containerId string) string {
	log.Infoln("Getting task id from container id")
	if data := ti.containerToTaskIndex[containerId]; data == nil {
		if data, err := ti.state.Get(ti.container2TaskBucket, containerId); err != nil {
			return ""
		} else {
			return data
		}

	} else {
		return *data
	}
}

// gets the container id from the gantry task id
func (ti taskIndex) GetContainerId(gantryId string) string {
	log.Infoln("Getting container id from task id")
	if data := ti.taskToContainerIndex[gantryId]; data == nil {
		if data, err := ti.state.Get(ti.task2ContainerBucket, gantryId); err != nil {
			return ""
		} else {
			return data
		}
	} else {
		return *data
	}
}
