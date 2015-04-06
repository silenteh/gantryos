package tasks

import (
	"fmt"
	"github.com/silenteh/gantryos/core/state"
	"github.com/silenteh/gantryos/utils"
	"testing"
)

func TestAddTaskId(t *testing.T) {
	dbName := "test_index.db"

	utils.RemoveDir(dbName)

	stateDb, err := state.InitSlaveDB(dbName)
	if err != nil {
		t.Fatal(err)
	}

	index := NewTaskIndex("task_2_cont", "cont_2_task", stateDb)

	containerId := "12345"
	gantryId := "67890"

	index.AddTaskId(containerId, gantryId)

	storedGantryId := index.GetTaskId(containerId)

	if storedGantryId != gantryId {
		t.Error("Task index storage error - cannot retrieve the gantry id from the container id")
	}

	storedContainerId := index.GetContainerId(gantryId) //getContainerId(gantryId)

	if storedContainerId != containerId {
		t.Error("Task index storage error - cannot retrieve the container id from the gantry id")
	}

	stateDb.Close()
	utils.RemoveDir(dbName)

	fmt.Println("- task index add and retrieve: OK")
}

func TestRemoveTaskId(t *testing.T) {

	dbName := "test_index_removal.db"
	utils.RemoveDir(dbName)

	stateDb, err := state.InitSlaveDB("test_index_removal.db")
	if err != nil {
		t.Fatal(err)
	}

	index := NewTaskIndex("task_2_cont", "cont_2_task", stateDb)

	containerId := "12345"
	gantryId := "67890"

	index.AddTaskId(containerId, gantryId)

	index.RemoveTaskId(containerId, gantryId)

	storedGantryId := index.GetTaskId(containerId)

	if storedGantryId != "" {
		t.Error("Task index storage remove error - retrived gantry id is not an empty string")
	}

	storedContainerId := index.GetContainerId(gantryId)

	if storedContainerId != "" {
		t.Error("Task index storage remove error - retrived container id is not an empty string")
	}

	stateDb.Close()
	utils.RemoveDir("test_index_removal.db")

	fmt.Println("- task index removal: OK")

}
