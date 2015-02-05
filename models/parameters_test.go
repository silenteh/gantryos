package models

import (
	"fmt"
	"testing"
)

func TestNewParameter(t *testing.T) {
	p := NewParameter("host", "127.0.0.1")
	if p.Key != "host" {
		t.Fatal("NewParameter key error")
	}

	if p.Value != "127.0.0.1" {
		t.Fatal("NewParameter value error")
	}

	fmt.Println("- NewParameter: SUCCESS")
}
