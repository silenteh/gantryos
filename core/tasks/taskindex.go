package tasks

// here we handle the mapping between gantry task IDs and containers task IDs

// this is the index that we need to lookup to map against gantry tasks ids and container ids
// this map needs to be persisted and should survive to restarts or crashes
var taskToContainerIndex = make(map[string]*string)
var containerToTaskIndex = make(map[string]*string)

// stores the container id according to the gantry task id
func addTaskId(containerId, gantryId string) {
	containerToTaskIndex[containerId] = &gantryId
	taskToContainerIndex[gantryId] = &containerId

}

func removeTaskId(containerId, gantryId string) {
	containerToTaskIndex[containerId] = nil
	taskToContainerIndex[gantryId] = nil
}

// gets the container id from the gantry task id
func getTaskId(containerId string) string {

	if data := containerToTaskIndex[containerId]; data == nil {
		return ""
	} else {
		return *data
	}
}

// gets the container id from the gantry task id
func getContainerId(gantryId string) string {
	//return *taskToContainerIndex[gantryId]

	if data := taskToContainerIndex[gantryId]; data == nil {
		return ""
	} else {
		return *data
	}
}
