package models

import (
	"fmt"
	"testing"
)

func TestNewContainerVolume(t *testing.T) {
	v, err := NewContainerVolume("/data", "/data/container0", false, CONTAINER_VOLUME_RO)
	if err != nil {
		t.Fatal(err)
	}

	if v.ContainerPath != "/data" {
		t.Fatal("NewContainerVolume: Container path not set correctly")
	}

	if v.HostPath != "/data/container0" {
		t.Fatal("NewContainerVolume: Host path not set correctly")
	}

	if v.Persistent {
		t.Fatal("NewContainerVolume: TEMP ERROR - we do not support persistence yet")
	}

	if v.Permission != CONTAINER_VOLUME_RO {
		t.Fatal("NewContainerVolume: Wrong volume permissions")
	}

	fmt.Println("- NewContainerVolume: SUCCESS")

}
