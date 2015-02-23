// +build integration

package tasks

import (
	"flag"
	"fmt"
	"github.com/silenteh/gantryos/models"
	mock "github.com/silenteh/gantryos/utils/testing"
	"os"
	"testing"
	"time"
)

func init() {
	os.Setenv("DOCKER_HOST", "tcp://192.168.59.103:2376")
	os.Setenv("DOCKER_CERT_PATH", "/Users/silenteh/.boot2docker/certs/boot2docker-vm")
	flag.Parse()
}

func TestStartDockerService(t *testing.T) {
	fmt.Println("Running integration tests...")

	// start the docker service
	service, err := StartDockerService()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("# Docker client started successfully")

	events := service.GetEventChannel()

	go func(c chan *models.TaskStatus) {
		for {
			data := <-c
			if data == nil {
				fmt.Println("Events channel closed.")
				break
			}

			fmt.Println(data.Message)
		}
	}(events)

	// mock a task
	taskInfo := mock.MakeGolangHelloTaskWithVolume()

	// fo the integration tests disable the force pull to speed them up
	taskInfo.Container.ForcePull = false

	containerId, err := service.Start(taskInfo.ToProtoBuf())
	if err != nil {
		t.Error(err)
	}

	if containerId == "" {
		t.Error("Container id is an empty string")
	}

	fmt.Println("# Container started successfully")

	if err = service.Status(containerId); err != nil {
		t.Error(err)
	}

	fmt.Println("# Container status success")

	time.Sleep(5 * time.Second)

	removeVolumes := true
	if err = service.Stop(containerId, removeVolumes); err != nil {
		t.Error(err)
	}

	fmt.Println("# Container stopped successfully")

	service.StopService()
	//close(events)

	fmt.Println("Integration tests: OK")
}
