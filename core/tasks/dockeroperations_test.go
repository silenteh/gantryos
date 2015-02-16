package tasks

import (
	"fmt"
	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/silenteh/gantryos/models"
	mock "github.com/silenteh/gantryos/utils/testing"
	"testing"
)

func TestMakeDockerVolumesBinds(t *testing.T) {

	task := mock.MakeGolangHelloTask()

	vol, err := models.NewContainerVolume("/tmp", "/var/tmp", false, models.CONTAINER_VOLUME_RW)
	if err != nil {
		t.Fatal("Error creating the volume bindings")
	}
	vols := models.NewContainerVolumes(vol)

	task.Container.Volumes = vols

	bindings := makeDockerVolumesBinds(task.ToProtoBuf())

	if len(bindings) == 0 {
		t.Fatal("Error creating the volume bindings")
	}

	if bindings[0] != "/var/tmp:/tmp" {
		t.Fatal("Error creating the volume bindings")
	}

	fmt.Println("- makeDockerVolumesBinds: OK")

}

func TestMakeDockerPortsAndBindings(t *testing.T) {
	task := mock.MakeGolangHelloTask()
	mapping, binding := makeDockerPortsAndBindings(task.ToProtoBuf())
	if len(mapping) == 0 {
		t.Fatal("Error generating port mappings")
	}

	var port dockerclient.Port = "8080/tcp"

	if mapping[port] != struct{}{} {
		t.Fatal("Error generating port mappings")
	}

	if len(binding) == 0 {
		t.Fatal("Error generating port mappings")
	}

}
