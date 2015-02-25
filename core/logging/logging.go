package logging

import (
	//"bufio"
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/silenteh/gantryos/config"
	"log/syslog"
	"os"
)

type LogInterface interface {
	Info(msg string)        // used to start a container starting (it does all the operations, like pull and start)
	Error(msg string)       // stops the container and removes the stopped container
	ToFileWriter() *os.File // writes to a file - it returns the handle so that it can be closed
}

func newSlaveContainerLogger() *logrus.Logger {
	var log = logrus.New()

	// With the default log.Formatter = new(logrus.TextFormatter) when a TTY is not attached,
	// the output is compatible with the logfmt format:
	log.Formatter = new(logrus.TextFormatter)

	switch {
	// Syslog
	case config.GantryOSConfig.Slave.SyslogServer.Hostname != "":
		proto := config.GantryOSConfig.Slave.SyslogServer.Protocol
		address := config.GantryOSConfig.Slave.SyslogServer.Hostname + ":" + config.GantryOSConfig.Slave.SyslogServer.Port

		// LOG_INFO
		syslogInfoHook, err := logrus_syslog.NewSyslogHook(proto, address, syslog.LOG_INFO, "")
		if err != nil {
			log.Error("Unable to connect to " + address + " syslog daemon")
		} else {
			log.Hooks.Add(syslogInfoHook)
		}

		// LOG_INFO
		syslogErrorHook, err := logrus_syslog.NewSyslogHook(proto, address, syslog.LOG_ERR, "")
		if err != nil {
			log.Error("Unable to connect to " + address + " syslog daemon")
		} else {
			log.Hooks.Add(syslogErrorHook)
		}
	// Logstash
	case config.GantryOSConfig.Slave.LogstashServer.Hostname != "":
		log.Formatter = new(logrus.JSONFormatter)

	}
	return log
}

func NewMasterLogger() *logrus.Logger {

	//config.GantryOSConfig.Slave.
	return nil

}
