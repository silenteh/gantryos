package resources

import (
	"fmt"
	"testing"
)

func TestGetTotalCPUsCount(t *testing.T) {
	totalCPUS := GetTotalCPUsCount()
	if totalCPUS <= 0 {
		t.Fatal("Error detecting the amount of CPUs")
	}
	fmt.Println("GetTotalCPUsCount: OK")
}
