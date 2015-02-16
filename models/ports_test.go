package models

import (
	"fmt"
	"testing"
)

func TestNewPortsMapping(t *testing.T) {

	portMappings := NewPortMapping(8080, 8080, "tcp")
	portsMapping := NewPortsMapping(portMappings)

	if len(portsMapping) == 0 {
		t.Fatal("Error generating ports mapping")
	}

	if len(portsMapping) != 1 {
		t.Fatal("Error generating ports mapping")
	}

	fmt.Println("- NewPortsMapping: SUCCESS")
}
