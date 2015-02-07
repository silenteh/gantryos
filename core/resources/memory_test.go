package resources

import (
	"fmt"
	"testing"
)

func TestTotalRam(t *testing.T) {
	ram := totalRam()

	if ram <= 1 {
		t.Fatal("Error calculating the total ram")
	}

	fmt.Println("totalRam: OK")
}
