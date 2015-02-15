package tasks

import (
	dockerclient "github.com/fsouza/go-dockerclient"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/core/proto"
	//docker "github.com/silenteh/gantryos/core/tasks/docker"
	"fmt"
	"strconv"
	"strings"
)

func startMonitor(task dockerTask) {

	//task.client.AddEventListener(listener)

}

// returns the container ID and an error in case
func startDockerContainer(task dockerTask) (string, error) {

	config, hostConfig := newDockerConfig(task.info)

	createContainerOptions := dockerclient.CreateContainerOptions{
		Name:       task.info.GetTaskId(),
		Config:     &config,
		HostConfig: &hostConfig,
	}

	var err error

	container, err := task.client.CreateContainer(createContainerOptions)

	err = task.client.StartContainer(container.ID, &hostConfig)

	return container.ID, err

}

// stops and remove the container
func stopDockerContainer(task dockerTask, containerId string) error {
	var err error
	err = task.client.StopContainer(containerId, 30) // waits max 30 seconds
	if err != nil {

		killOptions := dockerclient.KillContainerOptions{
			ID:     containerId,
			Signal: dockerclient.SIGKILL,
		}
		// try to kill it
		err = task.client.KillContainer(killOptions)
	}

	if err != nil {
		// try to remove the container
		dockerRemoveOpts := dockerclient.RemoveContainerOptions{
			ID:            containerId,
			Force:         true,  // removes the container even though it's running
			RemoveVolumes: false, // this forces removal of mounted volumes - we need a flag on the task
		}

		err = task.client.RemoveContainer(dockerRemoveOpts)
	}

	return err
}

func statusDockerContainer(task dockerTask, containerId string) (proto.TaskState, error) {
	container, err := task.client.InspectContainer(containerId)
	if err != nil {
		return proto.TaskState_TASK_LOST, err
	}

	switch {
	case container.State.Pid > 0:
		return proto.TaskState_TASK_RUNNING, nil
	case container.State.OOMKilled:
		return proto.TaskState_TASK_FAILED, nil
	default:
		return proto.TaskState_TASK_FINISHED, nil
	}

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

	config := dockerclient.Config{
		Hostname:        taskInfo.GetContainer().GetHostname(),
		Domainname:      taskInfo.GetContainer().GetDomainName(),
		User:            taskInfo.GetContainer().GetUser().GetName(),
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

	fmt.Printf("%s\n", config)
	fmt.Printf("%s\n", hostConfig)

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
				//HostIP:   port.HostIP,
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
