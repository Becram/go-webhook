package main

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	err := os.Setenv("TEST", "test")
	if err != nil {
		t.Error("cannot set env")
	}
	s := env("TEST")

	if s != "test" {
		t.Error("failed")
	}
}

func TestEnvOrDefault(t *testing.T) {
	d := "defaultvalue"

	s := envOrDefault("UNDEFINED", d)

	if s != d {
		t.Error("failed")
	}

}
