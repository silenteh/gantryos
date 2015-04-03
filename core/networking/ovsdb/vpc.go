package ovsdb

import (
	"encoding/json"
	//"errors"
	log "github.com/Sirupsen/logrus"
	//"math/rand"
	"net"
	"strconv"

	//"github.com/silenteh/gantryos/utils"
)

// MODELS
var defaultNetwork string = "192.168.100.0/24"

type address struct {
	DHCP    bool          // are the info set via DHCP ?
	Address net.IPAddr    // ip address and netmask
	Gateway net.IPAddr    // gateway
	Dns     []net.IPAddr  // dns config
	Iface   net.Interface // this is the interface we attach the vswitch to
}

// local slave switch
type Vswitch struct {
	RootId     string // ovsdb root UUID
	Id         string // bridge UUID
	Name       string
	STPEnabled bool
	VPCs       map[string]vpc // VPC vlan
	manager    *vswitchManager
}

// can contain multiple VPCs
type vpc struct {
	Name  string           // description name must be unique
	VLan  int              // vlan ID
	Ports map[string]vPort // the key here is the port name

	// // network part
	// Network *net.IPNet        // network of the VPC
	// UsedIPs map[string]string // map of IP -> interface name
	// Gateway net.IPAddr        // gateway of this VPC
	// Dns     []net.IPAddr      // dns config

}

func newVPC(name string, vlan int) vpc {
	vpc := vpc{
		Name:  name,
		VLan:  vlan,
		Ports: make(map[string]vPort),
		//UsedIPs: make(map[string]string),
	}

	return vpc
}

// each VPC has multiple ports
type vPort struct {
	Id         string // uuid
	Name       string
	Interfaces map[string]vInterface // Name of the interface
	Address    net.IPAddr
}

func newVPort() vPort {
	port := vPort{
		Interfaces: make(map[string]vInterface),
	}
	return port
}

// each port has an interface
type vInterface struct {
	Id             string // uuid
	Name           string
	ContainerId    string
	ContainerIface string
}

func (vswitch *Vswitch) toJson() string {
	data, err := json.Marshal(vswitch)
	if err != nil {
		return ""
	}
	return string(data)
}

// called only if the default switch is not present
func (vswitch *Vswitch) initDefaultVswitch(name string) error {
	vlan := 0

	// add the bridge
	bridgeUUID, err := addBridge(name, vswitch.RootId, vswitch.STPEnabled, vswitch.manager)
	if err != nil {
		return err
	}
	vswitch.Id = bridgeUUID

	// ad a vpc
	vpc := vswitch.AddVPC(name, defaultNetwork, vlan)
	if _, err := vpc.AddPort(name, bridgeUUID, name, "", vlan, vswitch.manager); err != nil {
		return err
	}

	return nil
}

func (vswitch *Vswitch) AddVPC(name, cidr string, vlan int) *vpc {

	// _, network, err := net.ParseCIDR(cidr)
	// if err != nil {
	// 	_, network, _ = net.ParseCIDR(defaultNetwork)
	// }

	vpc := vpc{
		//Network: network,
		Name:  name,
		VLan:  vlan,
		Ports: make(map[string]vPort),
	}

	vswitch.VPCs[strconv.Itoa(vlan)] = vpc
	return &vpc

}

func (vpc *vpc) DeleteVPC(vswitch *Vswitch, manager *vswitchManager) error {

	// loop through ports
	for pk, p := range vpc.Ports {
		// loop thorugh interfaces
		for ik, i := range p.Interfaces {
			// if err := deleteInterface(p.Id, i.Id, manager); err != nil {
			// 	return err
			// }
			log.Info("Deleting interface: ", i.Name)
			delete(p.Interfaces, ik)
		}
		// first try to delete physically the port if it succeeds then delete the in memory info
		if err := deletePort(vswitch.Id, p.Id, manager); err != nil {
			return err
		}
		log.Info("Deleting port: ", p.Name)
		delete(vpc.Ports, pk)

	}

	delete(vswitch.VPCs, strconv.Itoa(vpc.VLan))

	return nil

}

func (vswitch *Vswitch) DeleteVSwitch() error {

	for _, vpc := range vswitch.VPCs {
		if err := vpc.DeleteVPC(vswitch, vswitch.manager); err != nil {
			return err
		}
	}

	if err := deleteBridge(vswitch.RootId, vswitch.Id, vswitch.manager); err != nil {
		return err
	}

	log.Info("Deleting switch: ", vswitch.Name)
	vswitch.Id = ""
	vswitch.Name = ""

	return nil

}

func (vpc *vpc) AddPort(portName, bridgeUUID, vpcName, containerId string, vlan int, manager *vswitchManager) (vPort, error) {

	name := vpc.buildPortName(portName)
	port := newVPort()

	vInt := vInterface{
		Name: name,
	}

	// add the containerID and the container iface
	if containerId != "" {
		vInt.ContainerId = containerId
		vInt.ContainerIface = vpc.buildIfaceName(portName)
	}

	portUUID, intUUID, err := addPort(name, bridgeUUID, vpcName, vlan, vInt, manager)
	if err != nil {
		return port, err
	}

	vInt.Id = intUUID

	port.Id = portUUID
	port.Name = name
	port.Interfaces[name] = vInt

	vpc.Ports[name] = port

	return port, nil
}

// func (port *vPort) Up(containerId string) error {
// 	executor := utils.New()
// 	cmd := executor.Command("ip netns add", containerId)

// }

// func (port *vPort) Down(containerId string) error {
// 	executor := utils.New()
// 	cmd := executor.Command("ip netns add", containerId)

// }

func (vpc *vpc) buildPortName(name string) string {
	return name + "_" + strconv.Itoa(len(vpc.Ports)+1)
}

func (vpc *vpc) buildIfaceName(name string) string {
	return "veth_" + vpc.buildPortName(name)
}

func InitVSwitch(manager *vswitchManager) (Vswitch, error) {

	defaultSwitch := "gos0"

	exists, bridgeUUID := bridgeExists(defaultSwitch, manager)
	if exists {
		return loadVSwitch(bridgeUUID, manager)
	}

	return newVSwitch(defaultSwitch, false, manager)

}

func loadVSwitch(bridgeUUID string, manager *vswitchManager) (Vswitch, error) {

	return getAllBridgePorts(bridgeUUID, manager.GetRootUUID(), manager)

}

func newVSwitch(bridgeName string, stpEnabled bool, manager *vswitchManager) (Vswitch, error) {

	vswitch := Vswitch{
		RootId:  manager.GetRootUUID(),
		Name:    bridgeName,
		VPCs:    make(map[string]vpc),
		manager: manager,
	}

	err := vswitch.initDefaultVswitch(bridgeName)

	return vswitch, err

}
