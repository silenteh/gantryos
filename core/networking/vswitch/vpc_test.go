package vswitch

import (
	//mock "github.com/silenteh/gantryos/utils/testing"
	"testing"
)

func TestNewNetworkConf(t *testing.T) {
	//3309
	dhcp := false
	gwipNet := "192.168.10.1/24"
	ipNet := "192.168.10.10/24"

	netconfBridge, err := NewNetworkConf(dhcp, ipNet, gwipNet, []string{"8.8.8.8"})
	if err != nil {
		t.Error(err)
	}

	if netconfBridge.GatewayIP == "" {
		t.Error("Gateway IP should not be empty")
	}

	if netconfBridge.GatewayIPNet != gwipNet {
		t.Error("Gateway IP and NET do not match")
	}

	if netconfBridge.IP == "" {
		t.Error("Gateway IP should not be empty")
	}

	if netconfBridge.IPNet != ipNet {
		t.Error("Gateway IP and NET do not match")
	}

	// add a task

	// vswitch, err := InitVSwitch(ovsdbHost, ovsdbPort)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// defer vswitch.Close()

	//task := mock.MakeGolangHelloTask()

}
