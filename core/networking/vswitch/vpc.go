package vswitch

// TODO: TOO LINUX SPECIFIC !!!

import (
	"encoding/json"
	"errors"
	//"fmt"
	log "github.com/Sirupsen/logrus"
	//"math/rand"
	"github.com/silenteh/gantryos/core/resources"
	"github.com/silenteh/gantryos/utils"
	"net"
	"os"
	"strconv"
	"sync"
)

// MODELS
var defaultNetwork string = "192.168.100.0/24"

type networkConf struct {
	DHCP         bool          // are the info set via DHCP ?
	IP           string        // ip address and netmask
	Net          string        // ip network
	IPNet        string        // ip  + net
	GatewayIP    string        // gateway
	GatewayNet   string        // gateway
	GatewayIPNet string        // gateway
	Dns          []string      // dns config
	Iface        net.Interface // this is the interface we attach the vswitch to
}

// local slave switch
type Vswitch struct {
	RootId  string      // ovsdb root UUID
	mutex   sync.Mutex  // mutex for adding and removing entries to the VPCs map
	VPCs    map[int]Vpc // VPC vlan
	manager *vswitchManager
}

// can contain multiple VPCs
// when the VPC
type Vpc struct {
	Id          string           // bridge UUID
	STPEnabled  bool             // stp
	Name        string           // description name must be unique
	VLan        int              // vlan ID
	mutex       sync.Mutex       // mutex for adding and removing entries to the Ports map
	Ports       map[string]VPort // the key here is the port name
	networkConf networkConf      // the network config for the bridge (here we are interested in the GW because it will be the net config for the bridge)
}

// each VPC has multiple ports
type VPort struct {
	Id           string       // uuid
	Name         string       // the port name
	Tag          int          // the vlan tag
	Interfaces   []VInterface // Name of the interface
	NetConf      networkConf  // network configuration
	ContainerPID string       // the container PID - necessary for linking the container namespace
	TaskId       string       // the task ID - to know to which task this port belongs to. If the process dies we can remove it
	IsGateway    bool         // indicates if this is a bridge port
}

// each port has an interface
type VInterface struct {
	Id   string // uuid
	Name string // the name of the interface: ethN
	Type string // the type of the interface - mostly will be internal
}

// ===================================================

func NewNetworkConf(dhcp bool, ipAddress, gw string, dns []string) (networkConf, error) {

	netConf := networkConf{}

	var gwIp string
	var gwNet string

	var ip string
	var ipnet string

	if gw != "" {
		gwIpP, gwNetP, err := net.ParseCIDR(gw)
		if err != nil {
			return netConf, err
		}
		gwIp = gwIpP.String()
		gwNet = gwNetP.String()
	}

	if ipAddress != "" {
		ipP, ipNetP, err := net.ParseCIDR(ipAddress)
		if err != nil {
			return netConf, err
		}
		ip = ipP.String()
		ipnet = ipNetP.String()
	}

	allDns := []string{}
	for _, d := range dns {
		allDns = append(allDns, net.ParseIP(d).String())
	}

	return networkConf{
		DHCP:         dhcp,
		IP:           ip,
		Net:          ipnet,
		IPNet:        ipAddress,
		GatewayIP:    gwIp,
		GatewayNet:   gwNet,
		GatewayIPNet: gw,
		Dns:          allDns,
	}, nil
}

func newVPC(vlan int) Vpc {
	vpc := Vpc{
		Name:  createVPCName(vlan),
		VLan:  vlan,
		Ports: make(map[string]VPort),
	}

	return vpc
}

func newVPort() VPort {
	port := VPort{
		Interfaces: []VInterface{},
	}
	return port
}

func (port VPort) hasInterfaces() bool {
	return port.totalInterfaces() > 0
}

func (port VPort) totalInterfaces() int {
	return len(port.Interfaces)
}

func (vswitch *Vswitch) toJson() string {
	data, err := json.Marshal(vswitch)
	if err != nil {
		log.Errorln(err)
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
func (vswitch *Vswitch) initDefaultVswitch() error {
	vlan := 0
	containerPID := ""
	taskID := ""

	// ad a vpc
	if _, err := vswitch.AddVPC(vlan); err != nil {
		return err
	}
	if _, err := vswitch.AddPort(INTERFACE_INTERNAL, containerPID, taskID, vlan); err != nil {
		return err
	}

	return nil
}

func (vswitch *Vswitch) AddVPC(vlan int) (Vpc, error) {

	containerPID := ""
	taskID := ""
	intType := INTERFACE_INTERNAL

	name := createVPCName(vlan)

	// if the vpc exists already the return it
	if vpc, ok := vswitch.VPCs[vlan]; ok {
		return vpc, nil
	}

	// otherwise create it new
	vpc := Vpc{
		Name:       name,
		VLan:       vlan,
		Ports:      make(map[string]VPort),
		STPEnabled: true,
	}

	// add the bridge associated with the VPC
	bridgeUUID, err := addBridge(vpc.Name, vpc.VLan, vpc.STPEnabled, vswitch.manager)
	if err != nil {
		return vpc, err
	}
	vpc.Id = bridgeUUID

	portUUID, intUUID, err := addPort(name, name, bridgeUUID, vlan, intType, containerPID, taskID, vswitch.manager)
	if err != nil {
		return vpc, err
	}

	port := newVPort()
	port.Id = portUUID
	port.Name = name
	port.Tag = vlan
	port.IsGateway = true

	vint := VInterface{
		Id:   intUUID,
		Name: name,
		Type: intType,
	}

	port.Interfaces = append(port.Interfaces, vint)
	vpc.Ports[name] = port

	vswitch.mutex.Lock()
	vswitch.VPCs[vlan] = vpc
	vswitch.mutex.Unlock()

	return vpc, nil

}

func setIP() {

}

func createVPCName(vlan int) string {
	if vlan <= 0 {
		return "gbr0"
	}
	return "gbr" + strconv.Itoa(vlan)
}

func (vswitch *Vswitch) DeleteVPC(vlan int) error {

	vpc, ok := vswitch.VPCs[vlan]
	if !ok {
		return errors.New("Not such vpc: " + strconv.Itoa(vlan))
	}

	vpc.mutex.Lock()
	defer vpc.mutex.Unlock()

	// loop through ports

	for pk, p := range vpc.Ports {
		// loop thorugh interfaces
		if resources.DetectPlatform().Type == resources.LINUX && p.Interfaces[0].Type == INTERFACE_SYSTEM {
			log.Infoln("ip", "link", "delete", p.Name)
			log.Infoln(utils.ExecCommand(false, "ip", "link", "delete", p.Name))
		}

		log.Info("Deleting interfaces")
		p.Interfaces = []VInterface{}

		// first try to delete physically the port if it succeeds then delete the in memory info
		if err := deletePort(vpc.Id, p.Id, vswitch.manager); err != nil {
			return err
		}
		log.Info("Deleting port: ", p.Name)
		delete(vpc.Ports, pk)

	}

	vswitch.mutex.Lock()
	defer vswitch.mutex.Unlock()

	// delete the bridge
	if err := deleteBridge(vpc.Id, vswitch.manager); err != nil {
		return err
	}

	// remove the VPC from the map
	delete(vswitch.VPCs, vpc.VLan)

	return nil

}

func (port VPort) Delete() error {

	netConf := port.NetConf

	// this means it's a bridge.
	if port.IsGateway {
		// iptables -t nat -A POSTROUTING -s 192.168.2.0/24 -j MASQUERADE
		log.Infoln("iptables", "-t", "nat", "-D", "POSTROUTING", "-s", netConf.Net, "-j", "MASQUERADE")
		log.Infoln(utils.ExecCommand(false, "iptables", "-t", "nat", "-D", "POSTROUTING", "-s", netConf.Net, "-j", "MASQUERADE"))

		// brind the interface down
		// bring up the ineterface
		log.Infoln("ip", "link", "set", port.Interfaces[0].Name, "down")
		log.Infoln(utils.ExecCommand(false, "ip", "link", "set", port.Interfaces[0].Name, "down"))

		return nil
	}

	// bring the interface up
	// ip netns exec "$PID" ip link set "$INTERFACE" up
	log.Infoln("ip", "netns", "exec", port.ContainerPID, "ip", "link", "set", port.Interfaces[0].Name, "down")
	log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", port.ContainerPID, "ip", "link", "set", port.Interfaces[0].Name, "down"))

	return nil

}

// func (vswitch *Vswitch) DeleteVSwitch() error {

// 	vswitch.mutex.Lock()
// 	defer vswitch.mutex.Unlock()
// 	for _, vpc := range vswitch.VPCs {
// 		if err := vswitch.DeleteVPC(vpc.VLan); err != nil {
// 			return err
// 		}
// 	}

// 	if err := deleteBridge(vswitch.RootId, vswitch.Id, vswitch.manager); err != nil {
// 		return err
// 	}

// 	log.Info("Deleting switch: ", vswitch.Name)
// 	vswitch.Id = ""
// 	vswitch.Name = ""

// 	return nil

// }

func (vswitch *Vswitch) AddPort(interfaceType, containerPID, taskID string, vlan int) (VPort, error) {

	port := newVPort()
	vswitch.mutex.Lock()
	defer vswitch.mutex.Unlock()
	vpc, ok := vswitch.VPCs[vlan]
	if !ok {
		return port, errors.New("Not such vpc: " + strconv.Itoa(vlan))
	}

	vpc.mutex.Lock()
	defer vpc.mutex.Unlock()

	portName := taskID
	if taskID == "" {
		portName = vpc.buildPortName()
	}

	vInt := VInterface{
		Name: portName,
		Type: interfaceType,
	}

	// add the containerID and the container iface
	ifaceName := port.buildIfaceName()

	if containerPID != "" && taskID != "" {
		port.ContainerPID = containerPID
		port.TaskId = taskID
	}

	if resources.DetectPlatform().Type == resources.LINUX && interfaceType == INTERFACE_SYSTEM {
		log.Infoln("### Linux OS detected, creating physical interface")
		log.Infoln("ip", "link", "add", portName, "type", "veth", "peer", "name", containerPortName(portName))
		log.Infoln(utils.ExecCommand(false, "ip", "link", "add", portName, "type", "veth", "peer", "name", containerPortName(portName)))
	}

	portUUID, intUUID, err := addPort(portName, ifaceName, vpc.Id, vlan, vInt.Type, port.ContainerPID, port.TaskId, vswitch.manager)
	if err != nil {
		log.Errorln(err)
		return port, err
	}

	vInt.Id = intUUID

	port.Id = portUUID
	port.Name = portName
	port.Interfaces = []VInterface{vInt}

	vpc.Ports[portName] = port

	return port, nil
}

// TODO: this is linux specific !
// look into libcontainer netlink for syscalls
func (port VPort) Up() error {
	// if the IS is NOT Linux then skip it
	if resources.DetectPlatform().Type != resources.LINUX {
		return nil
	}

	netConf := port.NetConf

	// check if the containerPID is set, if so enable the namespace
	if port.ContainerPID != "" {
		utils.CreateDir("/var/run/netns")
		if err := os.Symlink("/proc/"+port.ContainerPID+"/ns/net", "/var/run/netns/"+port.ContainerPID); err != nil {
			return err
		}

		// enable the ip inside the container
		// # Move "${PORTNAME}_c" inside the container and changes its name.
		// ip link set "${PORTNAME}_c" netns "$PID"
		intName := port.Interfaces[0].Name
		ip := netConf.IP
		gw := netConf.GatewayIPNet
		log.Infoln("ip", "link", "set", intName, "netns", port.ContainerPID)
		log.Infoln(utils.ExecCommand(false, "ip", "link", "set", intName, "netns", port.ContainerPID))

		// set its ip ad bring it up
		if ip != "" {

			// set the IP
			log.Infoln("ip", "netns", "exec", port.ContainerPID, "ip", "addr", "add", ip, "brd + dev", intName)
			log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", port.ContainerPID, "ip", "addr", "add", ip, "brd + dev", intName))

			// bring the interface up
			// ip netns exec "$PID" ip link set "$INTERFACE" up
			log.Infoln("ip", "netns", "exec", port.ContainerPID, "ip", "link", "set", intName, "up")
			log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", port.ContainerPID, "ip", "link", "set", intName, "up"))

			// set the gateway if present
			if gw != "" {
				// add the gateway
				log.Infoln("ip", "netns", "exec", port.ContainerPID, "ip", "route", "add", "default", "gw", gw)
				log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", port.ContainerPID, "ip", "route", "add", "default", "gw", gw))
			}
		}

		return nil

	}

	// add the IP of the GW
	log.Infoln("ip", "addr", "add", netConf.GatewayIPNet, "brd", "+", "dev", port.Interfaces[0].Name)
	log.Infoln(utils.ExecCommand(false, "ip", "addr", "add", netConf.GatewayIPNet, "brd", "+", "dev", port.Interfaces[0].Name))

	// bring up the ineterface
	log.Infoln("ip", "link", "set", port.Interfaces[0].Name, "up")
	log.Infoln(utils.ExecCommand(false, "ip", "link", "set", port.Interfaces[0].Name, "up"))

	// masquerade
	// iptables -t nat -A POSTROUTING -o port.Name -j MASQUERADE
	log.Infoln("iptables", "-t", "nat", "-A", "POSTROUTING", "-s", netConf.Net, "-j", "MASQUERADE")
	log.Infoln(utils.ExecCommand(false, "iptables", "-t", "nat", "-A", "POSTROUTING", "-s", netConf.Net, "-j", "MASQUERADE"))

	return nil

	// // ip link add "${PORTNAME}_l" type veth peer name "${PORTNAME}_c"
	// // create an interface

	// // ip netns exec "$PID" ip link set dev "${PORTNAME}_c" name "$INTERFACE"
	// if port.hasInterfaces() {
	// 	vint := port.Interfaces[0]
	// 	cPort := containerPortName(port.Name)
	// 	cInt := containerPortName(vint.Name)

	// 	// brings up the port
	// 	log.Infoln("ip", "link", "set", port.Name, "up")
	// 	log.Infoln(utils.ExecCommand(false, "ip", "link", "set", port.Name, "up"))

	// 	// add the ip route
	// 	log.Infoln("route", "add", "-net", netConf.Net, "dev", port.Name)
	// 	log.Infoln(utils.ExecCommand(false, "route", "add", "-net", netConf.Net, "dev", port.Name))

	// 	log.Infoln("ip", "netns", "exec", containerPIDid, "ip", "link", "set", "dev", cPort, "name", cInt)
	// 	log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "ip", "link", "set", "dev", cPort, "name", cInt))

	// 	// ip netns exec "$PID" ip link set "$INTERFACE" up
	// 	log.Infoln("ip", "netns", "exec", containerPIDid, "ip", "link", "set", cInt, "up")
	// 	log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "ip", "link", "set", cInt, "up"))

	// 	// if [ -n "$ADDRESS" ]; then
	// 	//    ip netns exec "$PID" ip addr add "$ADDRESS" dev "$INTERFACE"
	// 	// fi
	// 	log.Infoln("ip", "netns", "exec", containerPIDid, "ip", "addr", "add", "192.168.2.111", "dev", cInt)
	// 	log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "ip", "addr", "add", "192.168.2.111", "dev", cInt))

	// 	// normal network
	// 	// add the ip route
	// 	log.Infoln("route", "add", "-net", netConf.Net, "dev", cInt)
	// 	log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "route", "add", "-net", netConf.Net, "dev", cInt))

	// 	// gateway
	// 	if netConf.GatewayIP != "" {

	// 		// add an ip address
	// 		log.Infoln("ip", "addr", "add", netConf.GatewayIPNet, "brd", "+", "dev", port.Name)
	// 		log.Infoln(utils.ExecCommand(false, "ip", "addr", "add", netConf.GatewayIPNet, "brd", "+", "dev", port.Name))

	// 		// add the gateway
	// 		log.Infoln("ip", "netns", "exec", containerPIDid, "ip", "route", "add", "default", "via", netConf.GatewayIP)
	// 		log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "ip", "route", "add", "default", "via", netConf.GatewayIP))

	// 		// masquerade
	// 		// iptables -t nat -A POSTROUTING -o port.Name -j MASQUERADE
	// 		log.Infoln("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", "!"+port.Name, "-s", netConf.Net, "-j", "MASQUERADE")
	// 		log.Infoln(utils.ExecCommand(false, "iptables", "-t", "nat", "-A", "POSTROUTING", "-o", "!"+port.Name, "-s", netConf.Net, "-j", "MASQUERADE"))
	// 	}

	// }

	// if [ -n "$GATEWAY" ]; then
	//    ip netns exec "$PID" ip route add default via "$GATEWAY"
	// fi
	//log.Infoln("ip", "netns", "exec", containerPIDid, "ip", "route", "add", "default", "via", "192.168.2.1")
	//log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", containerPIDid, "ip", "route", "add", "default", "via", "192.168.2.1"))

}

func (port *VPort) Down() error {

	// If the OS is NOT Linux then skip it
	log.Infoln("*************", resources.OsModel())
	if resources.DetectPlatform().Type != resources.LINUX {
		return nil
	}

	if port.ContainerPID != "" {

		// bring the interface down
		// ip netns exec "$PID" ip link set "$INTERFACE" down
		log.Infoln("ip", "netns", "exec", port.ContainerPID, "ip", "link", "set", port.Interfaces[0].Name, "down")
		log.Infoln(utils.ExecCommand(false, "ip", "netns", "exec", port.ContainerPID, "ip", "link", "set", port.Interfaces[0].Name, "down"))

		if err := os.Remove("/var/run/netns/" + port.ContainerPID); err != nil {
			return err
		}
	} else {
		// iptables -t nat -A POSTROUTING -o port.Name -j MASQUERADE
		log.Infoln("iptables", "-t", "nat", "-D", "POSTROUTING", "-s", port.NetConf.Net, "-j", "MASQUERADE")
		log.Infoln(utils.ExecCommand(false, "iptables", "-t", "nat", "-D", "POSTROUTING", "-s", port.NetConf.Net, "-j", "MASQUERADE"))

		//log.Infoln("ip", "link", "delete", port.Name)
		//log.Infoln(utils.ExecCommand(false, "ip", "link", "delete", port.Name))
		log.Infoln(utils.ExecCommand(false, "ip", "link", "set", port.Interfaces[0].Name, "down"))
	}
	return nil
}

func (vpc *Vpc) buildPortName() string {
	total := len(vpc.Ports)
	name := "geth" + "_" + strconv.Itoa(total)
	if vpc.VLan >= 0 {
		name = name + "_vlan_" + strconv.Itoa(vpc.VLan)
	}
	return name
}

func containerPortName(portName string) string {
	//total := len(vpc.Ports) + 1
	//return "gos" + "_" + strconv.Itoa(total) + "_vlan_" + strconv.Itoa(vpc.VLan), nil
	return portName + "_c"
}

func (port *VPort) buildIfaceName() string {
	total := len(port.Interfaces)
	return "eth" + strconv.Itoa(total)
}

// ======================================================================

func InitVSwitch(host, port string) (Vswitch, error) {

	defaultSwitch := "br0"

	// init the connection
	manager, err := newOVSDBClient(host, port)
	if err != nil {
		return Vswitch{}, err
	}

	exists, bridgeUUID := bridgeExists(defaultSwitch, manager)
	if exists {
		vswitch, err := loadVSwitch(bridgeUUID, manager)
		vswitch.manager = manager
		return vswitch, err
	}

	vswitch := newVSwitch(manager)
	err = vswitch.initDefaultVswitch()

	return vswitch, err

}

func loadVSwitch(bridgeUUID string, manager *vswitchManager) (Vswitch, error) {

	return getVswitch(manager)

}

func newVSwitch(manager *vswitchManager) Vswitch {

	vswitch := Vswitch{
		RootId:  manager.GetRootUUID(),
		VPCs:    make(map[int]Vpc),
		manager: manager,
	}

	return vswitch

}
