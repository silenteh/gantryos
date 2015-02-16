package testing

import (
	"fmt"
	"testing"
)

func TestMakeGolangHelloTaskInfo(t *testing.T) {

	task := MakeGolangHelloTask().ToProtoBuf()
	if task == nil {
		t.Fatal("Error creating the golang task for data mocking")
	}

	if task.GetContainer() == nil {
		t.Fatal("Error creating the golang task for data mocking - nil container")
	}

	if task.GetTaskName() == "" {
		t.Fatal("Error creating the golang task for data mocking - task empty name")
	}

	if task.GetContainer().GetImage() == "" {
		t.Fatal("Error creating the golang task for data mocking - task empty image")
	}

	if task.GetSlave() == nil {
		t.Fatal("Error creating the golang task for data mocking - task empty slave")
	}

	fmt.Println("- MakeGolangHelloTaskInfo: SUCCESS")

}

func TestMakeSlave(t *testing.T) {

	slave := MakeSlave(true)

	if slave == nil {
		t.Fatal("Error mocking the slave")
	}

	if slave.Id == "" {
		t.Fatal("Error mocking the slave - empty ID")
	}

	if slave.Hostname == "" {
		t.Fatal("Error mocking the slave - empty Hostname")
	}

	fmt.Println("- MakeSlave: SUCCESS")

}
