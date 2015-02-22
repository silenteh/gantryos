package tasks

import (
	"fmt"
	"testing"
)

func TestAddTaskId(t *testing.T) {

	containerId := "12345"
	gantryId := "67890"

	addTaskId(containerId, gantryId)

	storedGantryId := getTaskId(containerId)

	if storedGantryId != gantryId {
		t.Fatal("Task index storage error - cannot retrieve the gantry id from the container id")
	}

	storedContainerId := getContainerId(gantryId)

	if storedContainerId != containerId {
		t.Fatal("Task index storage error - cannot retrieve the container id from the gantry id")
	}

	fmt.Println("- task index add and retrieve: OK")
}

func TestRemoveTaskId(t *testing.T) {
	containerId := "12345"
	gantryId := "67890"

	addTaskId(containerId, gantryId)

	removeTaskId(containerId, gantryId)

	storedGantryId := getTaskId(containerId)

	if storedGantryId != "" {
		t.Fatal("Task index storage remove error - retrived gantry id is not an empty string")
	}

	storedContainerId := getContainerId(gantryId)

	if storedContainerId != "" {
		t.Fatal("Task index storage remove error - retrived container id is not an empty string")
	}

	fmt.Println("- task index removal: OK")

}
