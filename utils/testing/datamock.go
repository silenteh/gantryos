package testing

import (
	"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/core/resources"
	"github.com/silenteh/gantryos/models"
)

func MakeGolangHelloTask() *models.Task {
	slave := MakeSlave(true)

	name := "golang_test"
	image := "google/golang-hello"
	emptyVolumes := models.NewEmptyContainerVolumesSet()

	portMappings := models.NewPortMapping(8080, 8080, "tcp")
	portsMapping := models.NewPortsMapping(portMappings)

	env := models.NewEnvironmentVariable("GANTRY", "os")
	envs := models.NewEnvironmentVariables(env)

	container := models.NewContainer(name, image,
		"golangtest", "", "", true, proto.ContainerInfo_BRIDGE, emptyVolumes, portsMapping, envs)

	cpu := models.NewCPUResource(float64(2048))
	mem := models.NewCPUResource(float64(1024))

	allResources := models.MakeResources(cpu, mem)

	removeVolumesOnStop := true

	task := models.NewTask("golang_test_task", "1.0", removeVolumesOnStop, slave, allResources, nil, container, nil, nil)

	return task

}

func MakeGolangHelloTaskWithVolume() *models.Task {
	persistent := false
	vol, err := models.NewContainerVolume("/tmp", "/var/tmp", persistent, models.CONTAINER_VOLUME_RW)
	if err != nil {
		return nil
	}
	vols := models.NewContainerVolumes(vol)

	task := MakeGolangHelloTask()
	task.Container.Volumes = vols
	return task

}

func MakeSlave(registered bool) *models.Slave {
	// Port
	port := 7051
	// ==============================================

	// IP
	ip := "127.0.0.1"
	// ==============================================

	// Hostname
	hostname := resources.GetHostname()
	// ==============================================

	// Slave ID
	slaveId := "test_slave_id_123456789"

	checkpoint := true

	slaveInfo := models.NewSlave(slaveId, ip, hostname, port, checkpoint, registered)

	return slaveInfo
}

func MakeTaskStatus() *models.TaskStatus {

	return models.NewTaskStatusNoSlave("12345", "unwxiwnxiwnciencinec", "This is a test message from the TaskStatusChange", proto.TaskState_TASK_FAILED)

}
