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

func TestNewEnvironmentVariables(t *testing.T) {
	env := NewEnvironmentVariable("GANTRY", "os")
	envs := NewEnvironmentVariables(env)

	if len(envs) == 0 {
		t.Fatal("Error generating environment variables object")
	}

	if len(envs) != 1 {
		t.Fatal("Error generating environment variables object")
	}

	fmt.Println("- NewEnvironmentVariables: SUCCESS")
}
