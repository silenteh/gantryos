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

	fmt.Println("- makeDockerPortsAndBindings: OK")
}

func TestNewDockerConfig(t *testing.T) {

	task := mock.MakeGolangHelloTaskWithVolume()
	dc, hc := newDockerConfig(task.ToProtoBuf())

	// check the dockerConfig properties
	if len(dc.Env) == 0 || dc.Env[0] != "GANTRY=os" {
		t.Fatal("MakeGolangHelloTaskWithVolume has 1 environmnet variable: GANTRY=os")
	}

	// TODO HOST CONFIG additional properties
	if len(hc.Binds) == 0 || hc.Binds[0] != "/var/tmp:/tmp" {
		t.Fatal("Host Config volume bindings are wrong")
	}

	var port dockerclient.Port = "8080/tcp"

	portBinding := dockerclient.PortBinding{"0.0.0.0", "8080"}

	portBindingWrong := dockerclient.PortBinding{"0.0.0.0", "9090"}

	if len(hc.PortBindings) == 0 || hc.PortBindings[port][0] != portBinding {
		t.Fatal("Host Config port mapping bindings are wrong")
	}

	if len(hc.PortBindings) == 0 || hc.PortBindings[port][0] == portBindingWrong {
		t.Fatal("Host Config port mapping bindings are wrong")
	}

	fmt.Println("- newDockerConfig: OK")
}
