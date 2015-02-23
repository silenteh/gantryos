package logging

import (
	log "github.com/golang/glog"
)

type gantrylog struct{}

func NewGantryLog() LogInterface {
	var li LogInterface
	gantry := gantrylog{}
	li = gantry
	return li
}

func (l gantrylog) Info(tag, data string) {
	//log.Infoln(...)
	log.Infoln(data)
}

func (l gantrylog) Error(tag, data string) {
	log.Errorln(data)
}
