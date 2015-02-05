package models

import (
	"fmt"
	"testing"
)

func TestNewEnvironmentVariable(t *testing.T) {
	p := NewEnvironmentVariable("host", "localhost")
	if p.Name != "host" {
		t.Fatal("MakeEnvironmentVariable key error")
	}

	if p.Value != "localhost" {
		t.Fatal("MakeEnvironmentVariable value error")
	}

	fmt.Println("- NewEnvironmentVariable: SUCCESS")
}
