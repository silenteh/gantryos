package tasks

import (
	dockerclient "github.com/fsouza/go-dockerclient"
	log "github.com/golang/glog"
	gantrylog "github.com/silenteh/gantryos/core/logging"
	"github.com/silenteh/gantryos/core/proto"
	docker "github.com/silenteh/gantryos/core/tasks/docker"
	"github.com/silenteh/gantryos/models"

	"github.com/silenteh/gantryos/core/state"

	"bufio"
	//"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	container_image_exists_error   = 0
	container_pulling_image        = 1
	container_pulling_image_failed = 2
	container_starting             = 3
	container_starting_failed      = 4
	container_started              = 5
	channels_size                  = 1024
)

var (
	taskIndexBucket      = "task2container"
	containerIndexBucket = "container2task"
)

type dockerService struct {
	client           *dockerclient.Client
	dockerEvents     chan *dockerclient.APIEvents
	taskStatusEvents chan *models.TaskStatus
	taskLookup       taskIndex
}

// init a docker task and a docker client
func StartDockerService(stateDB state.StateDB) (TaskInterface, error) {
	var task TaskInterface
	var err error

	// init docker client
	endpoint := os.Getenv("DOCKER_HOST")
	log.Infoln("DOCKER_HOST env var set to:", endpoint)
	if endpoint == "" {
		endpoint = "unix:///var/run/docker.sock"
		log.Infoln("Falling back to local docker daemon at:", "/var/run/docker.sock")
	}

	certsPath := os.Getenv("DOCKER_CERT_PATH")

	var client *dockerclient.Client
	if certsPath != "" {

		cert := certsPath + "/cert.pem"
		key := certsPath + "/key.pem"
		ca := certsPath + "/ca.pem"
		client, err = dockerclient.NewTLSClient(endpoint, cert, key, ca)
	} else {
		client, err = dockerclient.NewClient(endpoint)
	}

	if err != nil || client == nil {
		return task, err
	}

	if err = client.Ping(); err != nil {
		return task, err
	}

	// channel to receive the events of the containers
	dockerEvents := make(chan *dockerclient.APIEvents, channels_size)
	taskEvents := make(chan *models.TaskStatus, channels_size)

	if err := client.AddEventListener(dockerEvents); err != nil {
		return task, err
	}

	t := dockerService{
		dockerEvents:     dockerEvents,
		taskStatusEvents: taskEvents,
		client:           client,
		taskLookup:       NewTaskIndex(taskIndexBucket, containerIndexBucket, stateDB),
	}

	// monitor docker events
	t.startMonitor()
	// assign the dockerService to the interface
	task = t

	return task, nil
}

func (t dockerService) StopService() {
	// remove the event listener
	t.client.RemoveEventListener(t.dockerEvents)

	// close the docker event channel
	close(t.dockerEvents)

	// close the task status channel
	close(t.taskStatusEvents)
}

// ==========================================================================================
// ==== DOCKER

func (t dockerService) Start(taskInfo *proto.TaskInfo) (string, error) {

	// image name
	image := taskInfo.GetContainer().GetImage()

	signalTaskStatus(t, taskInfo, container_starting, nil)

	// create the puller to download the image
	puller := docker.NewDockerPuller(t.client)

	// check if the image is already downloaded
	hasImage, err := puller.IsImagePresent(image)
	if err != nil {
		signalTaskStatus(t, taskInfo, container_pulling_image_failed, err)
		return "", err
	}

	// force pull if required by the task or if the image is not present
	if !hasImage || taskInfo.GetContainer().GetForcePullImage() {

		signalTaskStatus(t, taskInfo, container_pulling_image, nil)

		if err := puller.Pull(image); err != nil {
			signalTaskStatus(t, taskInfo, container_pulling_image_failed, err)
			return "", err
		}
	}

	// at this point we have the image so create the container with all the info we have in the task and then start it
	containerId, err := startDockerContainer(t, taskInfo)
	if err != nil || containerId == "" {
		signalTaskStatus(t, taskInfo, container_starting_failed, err)
		return "", err
	}

	// ORDER is important

	// update the taskIndex
	taskInfo.TaskId = &containerId
	t.taskLookup.AddTaskId(containerId, taskInfo.GetGantryTaskId())

	// send a signal
	signalTaskStatus(t, taskInfo, container_started, nil)

	// start to tail the logs
	go tailLogs(t, containerId)

	return containerId, nil

}

func (t dockerService) Stop(containerId string, removeVolumes bool) error {
	err := stopDockerContainer(t, removeVolumes, containerId)
	if err == nil {
		taskId := t.taskLookup.GetTaskId(containerId)
		t.taskLookup.RemoveTaskId(containerId, taskId)

		opts := dockerclient.RemoveContainerOptions{
			ID:            containerId,
			RemoveVolumes: removeVolumes,
			// A flag that indicates whether Docker should remove the container even if it is currently running.
			// At this point the container should have been stopped anyway
			Force: true,
		}

		if err := t.client.RemoveContainer(opts); err != nil {
			return err
		}
	}
	return err
}

func (t dockerService) Status(containerId string) (int, error) {

	return statusDockerContainer(t, containerId)
}

func (t dockerService) CleanContainers() error {
	return nil
}

func (t dockerService) CleanImages(image string) error {

	return nil
}

func (t dockerService) GetEventChannel() chan *models.TaskStatus {
	return t.taskStatusEvents
}

func (task dockerService) startMonitor() {

	go func(t dockerService) {
		for {
			// wait for a docker event
			dEvent := <-t.dockerEvents

			// this means the channel is closed
			if dEvent == nil {
				break
			}

			// get the task ID from the in memory index
			gantryTaskId := t.taskLookup.GetTaskId(dEvent.ID)

			// convert it to a model TaskStatus
			taskStatus := &models.TaskStatus{
				Id:        gantryTaskId,
				TaskId:    dEvent.ID,
				Message:   dEvent.Status,
				Timestamp: time.Now().UTC(),
			}

			// various docker events
			switch dEvent.Status {
			case "create", "exec_create":
				taskStatus.TaskState = proto.TaskState_TASK_STARTING
				break
			case "restart", "start", "exec_start", "unpause":
				taskStatus.TaskState = proto.TaskState_TASK_RUNNING
				break
			case "oom", "die":
				taskStatus.TaskState = proto.TaskState_TASK_FAILED
				break
			case "stop", "destroy", "kill":
				taskStatus.TaskState = proto.TaskState_TASK_FINISHED
				break
			case "paused":
				taskStatus.TaskState = proto.TaskState_TASK_PAUSED
				break
			case "extract":
				taskStatus.TaskState = proto.TaskState_TASK_CLONING_IMAGE
				break
			default:
				continue
			}

			t.taskStatusEvents <- taskStatus
		}
	}(task)

}

func signalTaskStatus(task dockerService, taskInfo *proto.TaskInfo, state int, err error) {
	taskId := taskInfo.GetTaskId()
	image := taskInfo.GetContainer().GetImage()
	taskStatus := models.NewTaskStatusNoSlave(taskInfo.GetGantryTaskId(), taskId, "", proto.TaskState_TASK_FAILED)

	switch state {
	case container_starting:
		taskStatus.Message = "Starting new container with image: " + image
		taskStatus.TaskState = proto.TaskState_TASK_STARTING
		break
	case container_image_exists_error:
		taskStatus.Message = "Assessing existance of the container image " + image + " failed"
		if err != nil {
			taskStatus.Message += " with error" + err.Error()
		}
		break
	case container_pulling_image:
		taskStatus.Message = "Pulling the image " + image
		taskStatus.TaskState = proto.TaskState_TASK_CLONING_IMAGE
		break
	case container_pulling_image_failed:
		taskStatus.Message = "Pulling the image " + image + "failed"
		if err != nil {
			taskStatus.Message += "with error " + err.Error()
		}
		break
	case container_starting_failed:
		taskStatus.Message = "Container from image " + image + " not started"
		if err != nil {
			taskStatus.Message += "with error " + err.Error()
		}
		break
	case container_started:
		taskStatus.Message = "Container id" + taskInfo.GetTaskId() + " from image " + image + " started successfully"
		taskStatus.TaskState = proto.TaskState_TASK_RUNNING
		break
	}

	task.taskStatusEvents <- taskStatus
}

// returns the container ID and an error in case
func startDockerContainer(task dockerService, taskInfo *proto.TaskInfo) (string, error) {

	config, hostConfig := newDockerConfig(taskInfo)

	containerName := "gantryos-" + taskInfo.GetGantryTaskId()

	createContainerOptions := dockerclient.CreateContainerOptions{
		Name:       containerName,
		Config:     &config,
		HostConfig: &hostConfig,
	}

	var err error

	container, err := task.client.CreateContainer(createContainerOptions)

	err = task.client.StartContainer(container.ID, &hostConfig)

	return container.ID, err

}

// stops and remove the container
func stopDockerContainer(task dockerService, removeVolumesOnStop bool, containerId string) error {
	var err error
	err = task.client.StopContainer(containerId, 30) // waits max 30 seconds
	if err != nil {

		killOptions := dockerclient.KillContainerOptions{
			ID:     containerId,
			Signal: dockerclient.SIGTERM,
		}
		// try to kill it
		err = task.client.KillContainer(killOptions)
	}

	if err != nil {
		// try to remove the container
		dockerRemoveOpts := dockerclient.RemoveContainerOptions{
			ID:            containerId,
			Force:         true,                // removes the container even though it's running
			RemoveVolumes: removeVolumesOnStop, // this forces removal of mounted volumes
		}

		err = task.client.RemoveContainer(dockerRemoveOpts)
	}

	return err
}

// returns the PID of the process running the container and/or the error
func statusDockerContainer(task dockerService, containerId string) (int, error) {

	PID := -1

	// get the gantry task Id from the container id
	gantryId := task.taskLookup.GetTaskId(containerId)
	// create a task status
	taskStatus := models.NewTaskStatusNoSlave(containerId, "", "", proto.TaskState_TASK_LOST)

	if gantryId == "" {
		err := errors.New("No gantry task id found linked to the container id:" + containerId + " - task LOST ?")
		taskStatus.Message = err.Error()
		task.taskStatusEvents <- taskStatus
	}

	if containerId == "" {
		err := errors.New("No container found linked to the task id:" + gantryId + " - task LOST ?")
		taskStatus.Message = err.Error()
		task.taskStatusEvents <- taskStatus
		return PID, err
	}

	// it measn we have found the container ID
	taskStatus.TaskId = containerId

	// inspect the container
	container, err := task.client.InspectContainer(containerId)

	if err != nil {
		taskStatus.Message = "Could not inspect the container id:" + containerId
		task.taskStatusEvents <- taskStatus
		return PID, err
	}

	PID = container.State.Pid
	taskStatus.PID = PID

	switch {
	case container.State.Pid > 0:
		taskStatus.Message = "Container id" + containerId + " RUNNING"
		taskStatus.TaskState = proto.TaskState_TASK_RUNNING
		break
	case container.State.OOMKilled:
		taskStatus.Message = "Container id" + containerId + " killed because ran OUT of MEMORY"
		taskStatus.TaskState = proto.TaskState_TASK_FAILED
		break
	}

	task.taskStatusEvents <- taskStatus
	return PID, nil

}

func newDockerConfig(taskInfo *proto.TaskInfo) (dockerclient.Config, dockerclient.HostConfig) {

	// default options which do not make sense in gantryos

	// needed for the logs
	attachStdout := true
	attachStderr := true

	attachStdin := false
	tty := false
	openStdin := false
	stdinOnce := false
	networkDisabled := false

	if taskInfo.GetContainer().GetNetwork() == proto.ContainerInfo_NONE {
		networkDisabled = true
	}

	// dockerEnvVars
	envs := taskInfo.GetContainer().GetEnvironments()
	envVars := envs.GetVariables()
	dockerEnvVars := make([]string, len(envVars))
	for i, v := range envVars {
		dockerEnvVars[i] = v.GetName() + "=" + v.GetValue()
	}

	// port mapping
	portMapping, portBinding := makeDockerPortsAndBindings(taskInfo) //_ := taskInfo.GetContainer().GetPortMappings()

	// Volumes
	dockerVols := make(map[string]struct{})
	for _, v := range taskInfo.GetContainer().GetVolumes() {
		dockerVols[v.GetContainerPath()] = struct{}{}
	}

	cpu := int64(1024)               // 1024 shares
	mem := int64(512 * 1024 * 1024)  // 512 Mb
	swap := int64(512 * 1024 * 1024) // 512 Mb

	for _, res := range taskInfo.GetResources() {
		if res.GetResourceType() == proto.ResourceType_MEMORY {
			mem = int64(res.GetScalar().GetValue())
		}

		if res.GetResourceType() == proto.ResourceType_CPU {
			cpu = int64(res.GetScalar().GetValue())
		}

		if res.GetResourceType() == proto.ResourceType_SWAP {
			swap = int64(res.GetScalar().GetValue())
		}
	}

	config := dockerclient.Config{
		Hostname:        taskInfo.GetContainer().GetHostname(),
		Domainname:      taskInfo.GetContainer().GetDomainName(),
		User:            taskInfo.GetContainer().GetUser().GetName(),
		Memory:          mem,
		CPUShares:       cpu,
		MemorySwap:      swap,
		AttachStdin:     attachStdin,
		AttachStdout:    attachStdout,
		AttachStderr:    attachStderr,
		ExposedPorts:    portMapping,
		Tty:             tty,
		OpenStdin:       openStdin,
		StdinOnce:       stdinOnce,
		Env:             dockerEnvVars,
		Cmd:             taskInfo.GetContainer().GetCmd(),
		Image:           taskInfo.GetContainer().GetImage(),
		WorkingDir:      taskInfo.GetContainer().GetWorkingDir(),
		Entrypoint:      taskInfo.GetContainer().GetEntryPoint(),
		SecurityOpts:    taskInfo.GetContainer().GetSecurityOptions(),
		VolumesFrom:     taskInfo.GetContainer().GetVolumesFrom(),
		Volumes:         dockerVols,
		OnBuild:         taskInfo.GetContainer().GetOnBuild(),
		NetworkDisabled: networkDisabled,
	}

	hostConfig := dockerclient.HostConfig{
		PortBindings: portBinding,
		Binds:        makeDockerVolumesBinds(taskInfo),
		Privileged:   taskInfo.GetContainer().GetPrivileged(),
	}

	//fmt.Printf("%s\n", config)
	//fmt.Printf("%s\n", hostConfig)

	return config, hostConfig

	/*
		Hostname        string              `json:"Hostname,omitempty" yaml:"Hostname,omitempty"`
		Domainname      string              `json:"Domainname,omitempty" yaml:"Domainname,omitempty"`
		User            string              `json:"User,omitempty" yaml:"User,omitempty"`
		Memory          int64               `json:"Memory,omitempty" yaml:"Memory,omitempty"`
		MemorySwap      int64               `json:"MemorySwap,omitempty" yaml:"MemorySwap,omitempty"`
		CPUShares       int64               `json:"CpuShares,omitempty" yaml:"CpuShares,omitempty"`
		CPUSet          string              `json:"Cpuset,omitempty" yaml:"Cpuset,omitempty"`
		AttachStdin     bool                `json:"AttachStdin,omitempty" yaml:"AttachStdin,omitempty"`
		AttachStdout    bool                `json:"AttachStdout,omitempty" yaml:"AttachStdout,omitempty"`
		AttachStderr    bool                `json:"AttachStderr,omitempty" yaml:"AttachStderr,omitempty"`
		PortSpecs       []string            `json:"PortSpecs,omitempty" yaml:"PortSpecs,omitempty"`
		ExposedPorts    map[Port]struct{}   `json:"ExposedPorts,omitempty" yaml:"ExposedPorts,omitempty"`
		Tty             bool                `json:"Tty,omitempty" yaml:"Tty,omitempty"`
		OpenStdin       bool                `json:"OpenStdin,omitempty" yaml:"OpenStdin,omitempty"`
		StdinOnce       bool                `json:"StdinOnce,omitempty" yaml:"StdinOnce,omitempty"`
		Env             []string            `json:"Env,omitempty" yaml:"Env,omitempty"`
		Cmd             []string            `json:"Cmd,omitempty" yaml:"Cmd,omitempty"`
		DNS             []string            `json:"Dns,omitempty" yaml:"Dns,omitempty"` // For Docker API v1.9 and below only
		Image           string              `json:"Image,omitempty" yaml:"Image,omitempty"`
		Volumes         map[string]struct{} `json:"Volumes,omitempty" yaml:"Volumes,omitempty"`
		VolumesFrom     string              `json:"VolumesFrom,omitempty" yaml:"VolumesFrom,omitempty"`
		WorkingDir      string              `json:"WorkingDir,omitempty" yaml:"WorkingDir,omitempty"`
		Entrypoint      []string            `json:"Entrypoint,omitempty" yaml:"Entrypoint,omitempty"`
		NetworkDisabled bool                `json:"NetworkDisabled,omitempty" yaml:"NetworkDisabled,omitempty"`
		SecurityOpts    []string            `json:"SecurityOpts,omitempty" yaml:"SecurityOpts,omitempty"`
		OnBuild         []string            `json:"OnBuild,omitempty" yaml:"OnBuild,omitempty"`
	*/

	/*
		Binds           []string
		ContainerIDFile string
		LxcConf         []utils.KeyValuePair
		Privileged      bool
		PortBindings    nat.PortMap
		Links           []string
		PublishAllPorts bool
		Dns             []string
		DnsSearch       []string
		ExtraHosts      []string
		VolumesFrom     []string
		Devices         []DeviceMapping
		NetworkMode     NetworkMode
		IpcMode         IpcMode
		PidMode         PidMode
		CapAdd          []string
		CapDrop         []string
		RestartPolicy   RestartPolicy
		SecurityOpt     []string
		ReadonlyRootfs  bool
	*/

}

func makeDockerPortsAndBindings(taskInfo *proto.TaskInfo) (map[dockerclient.Port]struct{}, map[dockerclient.Port][]dockerclient.PortBinding) {
	exposedPorts := map[dockerclient.Port]struct{}{}
	portBindings := map[dockerclient.Port][]dockerclient.PortBinding{}
	for _, port := range taskInfo.GetContainer().GetPortMappings() {
		exteriorPort := port.GetHostPort() //port.HostPort
		if exteriorPort == 0 {
			// No need to do port binding when HostPort is not specified
			continue
		}
		interiorPort := port.GetContainerPort() //port.ContainerPort
		// Some of this port stuff is under-documented voodoo.
		// See http://stackoverflow.com/questions/20428302/binding-a-port-to-a-host-interface-using-the-rest-api
		var protocol string
		switch strings.ToUpper(string(port.GetProtocol())) {
		case "UDP":
			protocol = "/udp"
		case "TCP":
			protocol = "/tcp"
		default:
			log.Warningf("Unknown protocol %q: defaulting to TCP", port.Protocol)
			protocol = "/tcp"
		}
		dockerPort := dockerclient.Port(strconv.Itoa(int(interiorPort)) + protocol)
		exposedPorts[dockerPort] = struct{}{}
		portBindings[dockerPort] = []dockerclient.PortBinding{
			{
				HostPort: strconv.Itoa(int(exteriorPort)),
				HostIP:   "0.0.0.0",
			},
		}
	}
	return exposedPorts, portBindings
}

func makeDockerVolumesBinds(taskInfo *proto.TaskInfo) []string {
	binds := []string{}
	for _, mount := range taskInfo.GetContainer().GetVolumes() {
		b := fmt.Sprintf("%s:%s", mount.GetHostPath(), mount.GetContainerPath())
		if mount.GetMode() == proto.Volume_RO {
			b += ":ro"
		}
		binds = append(binds, b)
	}
	return binds
}

func tailLogs(task dockerService, containerId string) {

	r, w := io.Pipe()

	// TODO: CONTAINER NAME !!!!
	loggerInstance := gantrylog.NewGantryContainerLog(containerId, containerId)
	f := loggerInstance.ToFileWriter()

	go func(writer *io.PipeWriter, s dockerService, clog gantrylog.LogInterface) {

		// when the goroutine exists close the PipeReader
		defer writer.Close()

		var ww io.Writer
		ww = writer

		opts := dockerclient.LogsOptions{
			Container:    containerId,
			OutputStream: ww,
			ErrorStream:  ww,
			Follow:       true,
			Stdout:       true,
			Stderr:       true,
		}

		if err := task.client.Logs(opts); err != nil {
			//clog.Error(err.Error())
			log.Errorln(err)
		}

	}(w, task, loggerInstance)

	go func(reader *io.PipeReader, containerId string, clog gantrylog.LogInterface, file *os.File) {
		// when the goroutine exists close the PipeReader
		defer reader.Close()

		// close the file handle
		defer file.Close()

		var rr io.Reader
		rr = reader

		scanner := bufio.NewScanner(rr)
		for scanner.Scan() {
			clog.Info(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Errorln(err)
		}
	}(r, containerId, loggerInstance, f)

}
