package tasks

import (
	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/silenteh/gantryos/core/proto"
	docker "github.com/silenteh/gantryos/core/tasks/docker"
	"os"
)

type dockerTask struct {
	client *dockerclient.Client
	info   *proto.TaskInfo
	//events chan<- *dockerclient.APIEvents
}

type lxcTask struct {
	info *proto.TaskInfo
}

type taskInterface interface {
	Start() (string, proto.TaskState, error) // used to start a container starting (it does all the operations, like pull and start)
	Stop() error                             // stops the container and removes the stopped container
	Status() (proto.TaskState, error)        // get the status of the container
	CleanContainers() error                  // cleans up the remained containers which are not running
	CleanImages() error                      // cleans up the images which are not used and are not being pulled (check if there is a pull running)
	Monitor()                                // monitors the container daemon for containers events
}

func MakeTask(taskInfo *proto.TaskInfo) taskInterface {
	var task taskInterface
	if taskInfo.GetContainer().GetType() == proto.ContainerInfo_DOCKER {

		// init docker client
		endpoint := os.Getenv("DOCKER_HOST")
		if endpoint == "" {
			endpoint = "unix:///var/run/docker.sock"
		}
		client, _ := dockerclient.NewClient(endpoint)

		t := dockerTask{}
		t.info = taskInfo
		t.client = client
		task = t
		return task
	}

	if taskInfo.GetContainer().GetType() == proto.ContainerInfo_DOCKER {

		t := lxcTask{}
		t.info = taskInfo
		task = t
		return task
	}

	return task
}

// ==========================================================================================
// ==== DOCKER

func (t dockerTask) Start() (string, proto.TaskState, error) {

	// image name
	image := t.info.GetContainer().GetImage()

	// create the puller to download the image
	puller := docker.NewDockerPuller(t.client)

	// check if the image is already downloaded
	hasImage, err := puller.IsImagePresent(image)
	if err != nil {
		return "", proto.TaskState_TASK_FAILED, err
	}

	// force pull if required by the task or if the image is not present
	if !hasImage || t.info.GetContainer().GetForcePullImage() {
		if err := puller.Pull(image); err != nil {
			return "", proto.TaskState_TASK_FAILED, err
		}
	}

	// at this point we have the image so create the container with all the info we have in the task and then start it
	containerId, err := startDockerContainer(t)
	if err != nil {
		return "", proto.TaskState_TASK_FAILED, err
	}

	return containerId, proto.TaskState_TASK_RUNNING, nil

}

func (t dockerTask) Stop() error {
	return stopDockerContainer(t, t.info.GetTaskId())
}

func (t dockerTask) Status() (proto.TaskState, error) {

	return statusDockerContainer(t, t.info.GetTaskId())
}

func (t dockerTask) CleanContainers() error {
	return nil
}

func (t dockerTask) CleanImages() error {
	return nil
}

func (t dockerTask) Monitor() {

}

// ==========================================================================================
// ==== LXC

func (t lxcTask) Start() (string, proto.TaskState, error) {
	return "", proto.TaskState_TASK_FAILED, nil
}

func (t lxcTask) Stop() error {
	return nil
}

func (t lxcTask) Status() (proto.TaskState, error) {

	return -1, nil
}

func (t lxcTask) CleanContainers() error {
	return nil
}

func (t lxcTask) CleanImages() error {
	return nil
}

func (t lxcTask) Monitor() {

}
