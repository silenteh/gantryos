package config

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/utils"
)

type gantryOSConfig struct {
	Master masterConfig `json:"master"`
	Slave  slaveConfig  `json:"slave"`
}

type masterConfig struct {
	IP   string
	Port int
}

type slaveConfig struct {
	IP         string
	Port       int
	Checkpoint bool
}

type slaveInfo struct {
	Id         string
	Registered bool // flag which says if it already registered
}

var GantryOSConfig gantryOSConfig = loadConfig()
var GantryOSSlaveId slaveInfo = loadSlaveId()

func loadConfig() gantryOSConfig {

	var localServerConfig gantryOSConfig
	var configFilePath string
	for _, path := range []string{"", "../", "../../"} {
		if utils.FileExists(path + "config.json") {
			configFilePath = path + "config.json"
			break
		}
	}
	if configFilePath == "" {
		panic("config.json not found!")
	}

	configFile := utils.ReadFile(configFilePath)
	if err := json.Unmarshal(configFile, &localServerConfig); err != nil {
		log.Fatalln("parsing config file: ", err.Error())
	}

	return localServerConfig
}

func loadSlaveId() slaveInfo {

	var configFilePath string
	for _, path := range []string{"", "../", "../../"} {
		if utils.FileExists(path + "slave.json") {
			configFilePath = path + "slave.json"
			break
		}
	}

	// create a new ID and write the file
	if configFilePath == "" {
		id := uuid.NewRandom().String()
		json := `{"Id":"` + id + `", "Registered" : false}`
		err := utils.WriteFile("slave.json", []byte(json), 0644)
		if err != nil {
			log.Fatalln(err)
		}
		return slaveInfo{Id: id, Registered: false}
	}

	configFile := utils.ReadFile(configFilePath)
	var slaveId slaveInfo
	if err := json.Unmarshal(configFile, &slaveId); err != nil {
		log.Fatalln("parsing slave id file: ", err.Error())
	}
	slaveId.Registered = true

	return slaveId
}
