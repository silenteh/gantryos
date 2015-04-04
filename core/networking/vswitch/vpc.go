package vswitch

// TODO: TOO LINUX SPECIFIC !!!

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	//"math/rand"
	"github.com/silenteh/gantryos/utils"
	"net"
	"os"
	"strconv"
)

// MODELS
var defaultNetwork string = "192.168.100.0/24"

// type address struct {
// 	DHCP    bool          // are the info set via DHCP ?
// 	Address net.IPAddr    // ip address and netmask
// 	Gateway net.IPAddr    // gateway
// 	Dns     []net.IPAddr  // dns config
// 	Iface   net.Interface // this is the interface we attach the vswitch to
// }

// local slave switch
type Vswitch struct {
	RootId     string // ovsdb root UUID
	Id         string // bridge UUID
	Name       string
	STPEnabled bool
	VPCs       map[int]vpc // VPC vlan
	manager    *vswitchManager
}

// can contain multiple VPCs
type vpc struct {
	Name  string           // description name must be unique
	VLan  int              // vlan ID
	Ports map[string]VPort // the key here is the port name

	// // network part
	// Network *net.IPNet        // network of the VPC
	// UsedIPs map[string]string // map of IP -> interface name
	// Gateway net.IPAddr        // gateway of this VPC
	// Dns     []net.IPAddr      // dns config

}

// each VPC has multiple ports
type VPort struct {
	Id         string // uuid
	Name       string
	Interfaces map[string]VInterface // Name of the interface
	Address    net.IPAddr
}

// each port has an interface
type VInterface struct {
	Id             string // uuid
	Name           string
	Type           string
	ContainerId    string
	ContainerIface string
}

// ===================================================

func newVPC(name string, vlan int) vpc {
	vpc := vpc{
		Name:  name,
		VLan:  vlan,
		Ports: make(map[string]VPort),
		//UsedIPs: make(map[string]string),
	}

	return vpc
}

func newVPort() VPort {
	port := VPort{
		Interfaces: make(map[string]VInterface),
	}
	return port
}

func (vswitch *Vswitch) toJson() string {
	data, err := json.Marshal(vswitch)
	if err != nil {
		return ""
	}
	return string(data)
}

// ===================================================

func (vswitch *Vswitch) Close() error {
	if vswitch.manager != nil {
		return vswitch.manager.Close()
	}
	return nil
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
	vswitch.AddVPC(name, defaultNetwork, vlan)
	if _, err := vswitch.AddPort(bridgeUUID, name, INTERFACE_INTERNAL, "", vlan); err != nil {
		return err
	}

	return nil
}

func (vswitch *Vswitch) AddVPC(name, cidr string, vlan int) *vpc {

	vpc := vpc{
		Name:  name,
		VLan:  vlan,
		Ports: make(map[string]VPort),
	}

	vswitch.VPCs[vlan] = vpc
	return &vpc

}

func (vswitch *Vswitch) DeleteVPC(vpcName int) error {

	vpc, ok := vswitch.VPCs[vpcName]
	if !ok {
		return errors.New("Not such vpc: " + strconv.Itoa(vpcName))
	}

	// loop through ports
	for pk, p := range vpc.Ports {
		// loop thorugh interfaces
		for ik, i := range p.Interfaces {
			// if err := deleteInterface(p.Id, i.Id, manager); err != nil {
			// 	return err
			// }
			// if i.Type != "internal" {
			// 	utils.ExecCommand(false, "ip", "link", "delete", i.Name)
			// }
			log.Info("Deleting interface: ", i.Name)
			delete(p.Interfaces, ik)
		}
		// first try to delete physically the port if it succeeds then delete the in memory info
		if err := deletePort(vswitch.Id, p.Id, vswitch.manager); err != nil {
			return err
		}
		log.Info("Deleting port: ", p.Name)
		delete(vpc.Ports, pk)

	}

	delete(vswitch.VPCs, vpc.VLan)

	return nil

}

func (vswitch *Vswitch) DeleteVSwitch() error {

	for _, vpc := range vswitch.VPCs {
		if err := vswitch.DeleteVPC(vpc.VLan); err != nil {
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

func (vswitch *Vswitch) AddPort(bridgeUUID, vpcName, interfaceType, containerId string, vlan int) (VPort, error) {

	port := newVPort()
	name, err := vswitch.buildPortName(vlan)
	if err != nil {
		return port, err
	}

	vInt := VInterface{
		Name: name,
		Type: interfaceType,
	}

	// add the containerID and the container iface
	iFaceName, err := vswitch.buildIfaceName(vlan)
	if err != nil {
		return port, err
	}
	if containerId != "" {
		vInt.ContainerId = containerId
		vInt.ContainerIface = iFaceName
	}

	portUUID, intUUID, err := addPort(name, bridgeUUID, vpcName, vlan, vInt, vswitch.manager)
	if err != nil {
		return port, err
	}

	vInt.Id = intUUID

	port.Id = portUUID
	port.Name = name
	port.Interfaces[name] = vInt

	vpc, ok := vswitch.VPCs[vlan]
	if !ok {
		return port, errors.New("Not such vpc: " + strconv.Itoa(vlan))
	}

	vpc.Ports[name] = port

	return port, nil
}

func (vswitch *Vswitch) AddContainerPort(portName, bridgeUUID, vpcName, containerId string, vlan int) (VPort, error) {

	port := newVPort()
	name, err := vswitch.buildPortName(vlan)
	if err != nil {
		return port, err
	}

	vInt := VInterface{
		Name: name,
	}

	log.Infoln("ip", "link", "add", name, "type", "veth", "peer", "name", name+"_c")
	log.Infoln(utils.ExecCommand(false, "ip", "link", "add", name, "type", "veth", "peer", "name", name+"_c"))

	// add the containerID and the container iface
	// add the containerID and the container iface
	iFaceName, err := vswitch.buildIfaceName(vlan)
	if err != nil {
		return port, err
	}
	if containerId != "" {
		vInt.ContainerId = containerId
		vInt.ContainerIface = iFaceName
	}

	portUUID, intUUID, err := addPort(name, bridgeUUID, vpcName, vlan, vInt, vswitch.manager)
	if err != nil {
		return port, err
	}

	vInt.Id = intUUID

	port.Id = portUUID
	port.Name = name
	port.Interfaces[name] = vInt

	vpc, ok := vswitch.VPCs[vlan]
	if !ok {
		return port, errors.New("Not such vpc: " + strconv.Itoa(vlan))
	}

	vpc.Ports[name] = port

	return port, nil
}

func (port *VPort) Up(containerPIDid string) error {
	// /proc/3887/ns/net /var/run/netns/3887
	if err := os.Symlink("/proc/"+containerPIDid+"/ns/net", "/var/run/netns/"+containerPIDid); err != nil {
		return err
	}

	// ip link add "${PORTNAME}_l" type veth peer name "${PORTNAME}_c"
	// create an interface

	// brings up the interface
	log.Infoln("ip", "link", "set", port.Name, "up")
	log.Infoln(utils.ExecCommand(false, "ip", "link", "set", port.Name, "up"))

	// add an ip address
	log.Infoln("ip", "addr", "add", "192.168.2.1", "dev", port.Name)
	log.Infoln(utils.ExecCommand(false, "ip", "addr", "add", "192.168.2.1", "dev", port.Name))

	// add the ip route
	log.Infoln("route", "add", "-net", "192.168.2.0/24", "dev", port.Name)
	log.Infoln(utils.ExecCommand(false, "route", "add", "-net", "192.168.2.0/24", "dev", port.Name))

	// # Move "${PORTNAME}_c" inside the container and changes its name.
	// ip link set "${PORTNAME}_c" netns "$PID"
	log.Infoln("ip", "link", "set", port.Name+"_c", "netns", containerPIDid)
	log.Infoln(utils.ExecCommand(false, "ip", "link", "set", port.Name+"_c", "netns", containerPIDid))

	// ip netns exec "$PID" ip link set dev "${PORTNAME}_c" name "$INTERFACE"
	log.Infoln("ip", "netns", "exec", containerPIDid, "ip", "link", "set", "dev", port.Name+"_c", "name", port.Interfaces[port.Name].Name+"_c")
	log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "ip", "link", "set", "dev", port.Name+"_c", "name", port.Interfaces[port.Name].Name+"_c"))

	// ip netns exec "$PID" ip link set "$INTERFACE" up
	log.Infoln("ip", "netns", "exec", containerPIDid, "ip", "link", "set", port.Interfaces[port.Name].Name+"_c", "up")
	log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "ip", "link", "set", port.Interfaces[port.Name].Name+"_c", "up"))

	// if [ -n "$ADDRESS" ]; then
	//    ip netns exec "$PID" ip addr add "$ADDRESS" dev "$INTERFACE"
	// fi
	log.Infoln("ip", "netns", "exec", containerPIDid, "ip", "addr", "add", "192.168.2.111", "dev", port.Interfaces[port.Name].Name+"_c")
	log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "ip", "addr", "add", "192.168.2.111", "dev", port.Interfaces[port.Name].Name+"_c"))

	// if [ -n "$GATEWAY" ]; then
	//    ip netns exec "$PID" ip route add default via "$GATEWAY"
	// fi
	//log.Infoln("ip", "netns", "exec", containerPIDid, "ip", "route", "add", "default", "via", "192.168.2.1")
	//log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "ip", "route", "add", "default", "via", "192.168.2.1"))

	return nil
}

func (port *VPort) Down(containerPIDid string) error {

	log.Infoln("ip", "link", "delete", port.Name)
	log.Infoln(utils.ExecCommand(false, "ip", "link", "delete", port.Name))
	//log.Infoln(utils.ExecCommand(false, "ip", "link", "set", port.Name, "down"))

	if err := os.Remove("/var/run/netns/" + containerPIDid); err != nil {
		return err
	}
	return nil
}

func (vswitch *Vswitch) buildPortName(vlan int) (string, error) {
	vpc, ok := vswitch.VPCs[vlan]

	if !ok {
		return "", errors.New("Not such vpc: " + strconv.Itoa(vlan))
	}
	return "gos" + "_" + strconv.Itoa(len(vpc.Ports)+1), nil
}

func (vswitch *Vswitch) buildIfaceName(vlan int) (string, error) {
	ifaceName, err := vswitch.buildPortName(vlan)
	return "veth_" + ifaceName, err
}

// ======================================================================

func InitVSwitch(host, port string) (Vswitch, error) {

	// init the connection
	manager, err := newOVSDBClient(host, port)
	if err != nil {
		return Vswitch{}, err
	}

	defaultSwitch := "gos0"

	exists, bridgeUUID := bridgeExists(defaultSwitch, manager)
	if exists {
		vswitch, err := loadVSwitch(bridgeUUID, manager)
		vswitch.manager = manager
		return vswitch, err
	}

	vswitch, err := newVSwitch(defaultSwitch, false, manager)
	vswitch.manager = manager

	return vswitch, err

}

func loadVSwitch(bridgeUUID string, manager *vswitchManager) (Vswitch, error) {

	return getAllBridgePorts(bridgeUUID, manager.GetRootUUID(), manager)

}

func newVSwitch(bridgeName string, stpEnabled bool, manager *vswitchManager) (Vswitch, error) {

	vswitch := Vswitch{
		RootId:  manager.GetRootUUID(),
		Name:    bridgeName,
		VPCs:    make(map[int]vpc),
		manager: manager,
	}

	err := vswitch.initDefaultVswitch(bridgeName)

	return vswitch, err

}
