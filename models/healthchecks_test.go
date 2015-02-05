package models

import (
	"fmt"
	"testing"
)

func TestNewHealthCheck(t *testing.T) {
	hc := NewHealthCheck(80, "/mypath", 10, 5, 6, 7, 3, []int{200, 301, 302})

	if hc.Port != 80 {
		t.Fatal("MakeHealthCheck port not assigned")
	}

	if hc.Path != "/mypath" {
		t.Fatal("MakeHealthCheck path not assigned")
	}

	if hc.IntervalSeconds != 10 {
		t.Fatal("MakeHealthCheck interval second not assigned")
	}

	if hc.TimeoutSeconds != 5 {
		t.Fatal("MakeHealthCheck timeout second not assigned")
	}

	if hc.DelaySeconds != 6 {
		t.Fatal("MakeHealthCheck delay second not assigned")
	}

	if hc.GracePeriodSeconds != 7 {
		t.Fatal("MakeHealthCheck grace period not assigned")
	}

	if hc.ConsecuritveFailures != 3 {
		t.Fatal("MakeHealthCheck consecutive failures not assigned")
	}

	if !arrayContains(hc.ExpectedStatuses, 200) {
		t.Fatal("MakeHealthCheck expected statuses not assigned")
	}

	if !arrayContains(hc.ExpectedStatuses, 301) {
		t.Fatal("MakeHealthCheck expected statuses not assigned")
	}

	if !arrayContains(hc.ExpectedStatuses, 302) {
		t.Fatal("MakeHealthCheck expected statuses not assigned")
	}

	hc = NewHealthCheck(0, "", 0, 0, 0, 0, 0, []int{200, 301, 302})

	if hc.Port != 80 {
		t.Fatal("MakeHealthCheck port not assigned")
	}

	if hc.Path != "/" {
		t.Fatal("MakeHealthCheck path not assigned")
	}

	if hc.IntervalSeconds != 10 {
		t.Fatal("MakeHealthCheck interval second not assigned")
	}

	if hc.TimeoutSeconds != 5 {
		t.Fatal("MakeHealthCheck timeout second not assigned")
	}

	if hc.DelaySeconds != 5 {
		t.Fatal("MakeHealthCheck delay second not assigned")
	}

	if hc.GracePeriodSeconds != 5 {
		t.Fatal("MakeHealthCheck grace period not assigned")
	}

	if hc.ConsecuritveFailures != 1 {
		t.Fatal("MakeHealthCheck consecutive failures not assigned")
	}

	if !arrayContains(hc.ExpectedStatuses, 200) {
		t.Fatal("MakeHealthCheck expected statuses not assigned")
	}

	if !arrayContains(hc.ExpectedStatuses, 301) {
		t.Fatal("MakeHealthCheck expected statuses not assigned")
	}

	if !arrayContains(hc.ExpectedStatuses, 302) {
		t.Fatal("MakeHealthCheck expected statuses not assigned")
	}

	fmt.Println("- NewHealthCheck: SUCCESS")
}

func arrayContains(array []int, match int) bool {
	for _, v := range array {
		if v == match {
			return true
		}
	}
	return false
}
