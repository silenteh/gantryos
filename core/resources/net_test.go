package resources

import (
	"fmt"
	"testing"
)

func TestNetStats(t *testing.T) {
	stats := netStats()

	if stats["eth0"].InterfaceName == "" {
		t.Fatal("Error reading the network stats information")
	}

	if stats["eth0"].RXBytes == 0 {
		t.Fatal("Error reading the network stats information")
	}

	fmt.Println("Network Stats: OK")

}
