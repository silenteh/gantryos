package resources

import (
	"fmt"
	"testing"
)

func TestGetRedhatishVersion(t *testing.T) {
	version := getRedhatishVersion("Red Hat Enterprise Linux Server release 6.3 (Santiago)")
	expectedVersion := "6.3"
	if version != expectedVersion {
		t.Fatal("GetRedhatishVersion - Failed to get version of: Red Hat Enterprise Linux Server release 6.3 (Santiago)")
	}

	version = getRedhatishVersion("Oracle Linux Server release 5.7")
	expectedVersion = "5.7"
	if version != expectedVersion {
		t.Fatal("GetRedhatishVersion - Failed to get version of: Oracle Linux Server release 5.7")
	}

	version = getRedhatishVersion("Fedora release 22 (Rawhide)")
	expectedVersion = "22 (rawhide)"
	if version != expectedVersion {
		t.Fatal("GetRedhatishVersion - Failed to get version of: Fedora release 22 (Rawhide)")
	}

	fmt.Println("GetRedhatishVersion: OK")
}

func TestGetRedhatishPlatform(t *testing.T) {
	version := getRedhatishPlatform("Parallels Cloud Server 6.0.5 (1784)")
	expectedVersion := "parallels"
	if version != expectedVersion {
		t.Fatal("GetRedhatishPlatform - Failed to get platform of: Parallels Cloud Server 6.0.5 (1784)")
	}

	version = getRedhatishPlatform("Red Hat Enterprise Linux Server release 6.3 (Santiago)")
	expectedVersion = "redhat"
	if version != expectedVersion {
		t.Fatal("GetRedhatishVersion - Failed to get version of: Red Hat Enterprise Linux Server release 6.3 (Santiago)")
	}

	fmt.Println("GetRedhatishPlatform: OK")
}
