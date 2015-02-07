package resources

import (
	"fmt"
	"testing"
)

func TestDetectOS(t *testing.T) {
	os := detectOS()

	if os == UNKNOWN {
		t.Fatal("Error detecting the operating system")
	}

	fmt.Println("detectOS: OK")
}
