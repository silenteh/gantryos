package services

import (
	log "github.com/golang/glog"
)

// this listener received the proto.Envelope data
// the it detects which sub-field is available
// and therefore which request the client is doing
// this method blocks

func (master *masterServer) startListener() {
	//total := 0

	go func(m *masterServer) {

		for {
			data := <-m.readerChannel
			//envlope := <-m.readerChannel
			// dereference
			//data := *envlope

			switch {
			case data.MasterInfo != nil:
				continue
			case data.RegisterSlave != nil:
				slave := data.RegisterSlave.GetSlave()
				log.Infoln(slave.GetHostname(), "registered as a slave with id", slave.GetId())
				continue
			case data.ReRegisterSlave != nil:
				slave := data.ReRegisterSlave.GetSlave()
				log.Infoln(slave.GetHostname(), "re-registered as a slave with id", slave.GetId())
				continue
			case data.SlaveInfo != nil:
				continue
			case data.Heartbeat != nil:
				slave := data.Heartbeat.GetSlave()
				log.Infoln("Slave id", slave.GetId(), "with hostname", slave.GetHostname(), "sent an heartbeat.")
				continue
			default:
				log.Errorln("Got an unknown request from the GatryOS slave")
				continue
			}
		}

	}(master)
}

// ********************** SLAVE ***********************************************************
