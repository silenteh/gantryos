package resources

import (
	"fmt"
	"testing"
)

func TestLoadAverage(t *testing.T) {
	load := loadAverage()
	if load.Minute == 0 {
		t.Fatal("Error detecting the operating system LOAD")
	}

	fmt.Println("loadAverage: OK")
}
