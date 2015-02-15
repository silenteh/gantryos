package services

import (
	//"fmt"
	log "github.com/golang/glog"
	//"github.com/silenteh/gantryos/core/proto"
	"github.com/silenteh/gantryos/core/tasks"
	"github.com/silenteh/gantryos/models"
)

func (slave *slaveServer) startSlaveListener() {
	go func(s *slaveServer) {
		for {

			// this blocks until a message is available in the queue
			// message processing is sequential
			// all messages are idempotent
			data := <-s.readerChannel

			switch {
			case data.GetRunTask() != nil:
				task := tasks.MakeTask(data.GetRunTask().GetTask())

				containerId, taskState, err := task.Start()

				taskStatus := models.NewTaskStatus(data.GetRunTask().GetTask().GetTaskId(), containerId, err.Error(), &taskState, s.slave)

				s.writerChannel <- taskStatus.ToProtoBuf()

				continue
			case data.Heartbeat != nil:
				log.Infoln(data.Heartbeat.GetSlave().GetHostname())
			default:
				log.Errorln("Got an unknown request from the GatryOS slave")
				continue
			}

		}
	}(slave)
}
