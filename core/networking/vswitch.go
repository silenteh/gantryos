package networking

import (
	"errors"
	"github.com/silenteh/gantryos/networking/ovsdb"
	"github.com/silenteh/gantryos/utils"
)

func NewVswitch(bridgeName string) error {

	if result := utils.ExecCommand(false, "ovs-vsctl", "add-br", bridgeName); result != "" {
		return errors.New(result)
	}
	return nil

}

func AddPort(bridgeName, portName string) error {
	if result := utils.ExecCommand(false, "ovs-vsctl", "add-port", bridgeName, portName); result != "" {
		return errors.New(result)
	}
	return nil
}

func SetVLAN(bridgeName string) error {
	return nil
}

func SetDockerIP(bridgeName, dockerIp string) error {
	return nil
}

func SetDockerVLAN(bridgeName, dockerIp string, VLANId int) error {
	return nil
}
