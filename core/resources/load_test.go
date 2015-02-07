package resources

import (
	"fmt"
	"testing"
)

func TestLoadAverage(t *testing.T) {
	load := loadAverage()
	if load.Minute <= 0 {
		t.Fatalf("Error detecting the operating system LOAD on %s", detectOS())
	}

	fmt.Println("loadAverage: OK")
}
