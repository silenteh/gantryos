package services

import (
	"fmt"
	log "github.com/golang/glog"
	//"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/core/tasks"
	"github.com/silenteh/gantryos/models"
)

func (slave *slaveServer) startSlaveListener() error {

	// init the docker service
	dockerService, err := tasks.StartDockerService()
	if err != nil {
		return err
	}

	// start to forward the task events to the master
	go startDockerEvents(slave, dockerService.GetEventChannel())

	// start its

	go func(s *slaveServer, t tasks.TaskInterface) {

		// when the function exists then it stops also the events channel and therefore also
		defer t.StopService()

		for {

			// this blocks until a message is available in the queue
			// message processing is sequential
			// all messages are idempotent
			data := <-s.readerChannel

			// means the channel was closed therefore exit gthe for loop
			if data == nil {
				break
			}

			switch {
			case data.GetRunTask() != nil:

				fmt.Println("Got Run TASK")
				// get the task information
				taskInfo := data.GetRunTask().GetTask()

				// try to start the task
				_, err := t.Start(taskInfo)

				if err != nil {
					log.Errorln(err)
				}

				// all events are sent back from the docker events listener
				// TODO: test the taskId

				continue
			case data.Heartbeat != nil:
				log.Infoln(data.Heartbeat.GetSlave().GetHostname())
				fmt.Println("Got heartbeat")
			default:
				log.Errorln("Got an unknown request from the GatryOS slave")
				fmt.Println("Got an unknown request from the GatryOS slave")
				continue
			}

		}
	}(slave, dockerService)

	return nil
}

func startDockerEvents(slave *slaveServer, eventsChannel chan *models.TaskStatus) {

	for {
		event := <-eventsChannel
		// means the channel was closed
		if event == nil {
			break
		}
		// write the message to the channel
		slave.taskStateChange(event)

	}

}
