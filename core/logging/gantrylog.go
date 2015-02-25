package logging

import (
	//"bufio"
	"github.com/Sirupsen/logrus"
	"github.com/silenteh/gantryos/config"
	"os"
)

type gantrylog struct {
	log           *logrus.Logger
	containerId   string
	containerName string
	logDir        string
	logName       string
}

func NewGantryContainerLog(containerId, containerName string) LogInterface {
	var li LogInterface
	logDir := config.GantryOSConfig.Slave.ContainersLogDir
	gantry := gantrylog{
		log:           newSlaveContainerLogger(),
		containerId:   containerId,
		containerName: containerName,
		logDir:        logDir,
		logName:       logDir + "/" + containerName + ".log",
	}
	li = gantry
	return li
}

func (gl gantrylog) Info(msg string) {
	gl.log.WithFields(logrus.Fields{
		"cname": gl.containerName,
		"cid":   gl.containerId,
	}).Infoln(msg)
}

func (gl gantrylog) Error(msg string) {
	gl.log.WithFields(logrus.Fields{
		"cname": gl.containerName,
		"cid":   gl.containerId,
	}).Errorln(msg)
}

func (gl gantrylog) ToFileWriter() *os.File {

	fo, err := os.Create(gl.logName)
	if err != nil {
		gl.log.Errorln(err)
		return nil
	}

	gl.log.Out = fo

	return fo
	// // // open output file
	// fo, err := os.Create(gl.logName)
	// if err != nil {
	// 	gl.log.Errorln(err)
	// 	return nil, nil
	// }

	// // make a write buffer
	// w := bufio.NewWriter(fo)
	// w.WriteString("Logging for container " + gl.containerId + " initialized.\n")
	// w.Flush()

	// return fo, w
}
