package resources

import (
	"fmt"
	"testing"
)

func TestGetHostname(t *testing.T) {
	hostname := GetHostname()
	if hostname == "unknown" {
		t.Fatal("Error detecting the hostname")
	}
	fmt.Println("TestGetHostname: OK")
}
