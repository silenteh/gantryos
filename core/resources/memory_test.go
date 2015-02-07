package resources

import (
	"fmt"
	"testing"
)

func TestTotalRam(t *testing.T) {
	ram := totalRam()

	if ram <= 1 {
		t.Fatalf("Error calculating the total ram on %s", detectOS())
	}

	fmt.Println("totalRam: OK")
}
