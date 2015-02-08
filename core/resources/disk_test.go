package resources

import (
	"fmt"
	"testing"
)

func TestLayout(t *testing.T) {
	disk := layout()

	if disk["/"].Device == "" {
		t.Fatal("Error getting the disk layout")
	}

	if disk["/"].Size == "0" {
		t.Fatal("Error getting the disk layout")
	}

	fmt.Println("Disk Layout: OK")
}
