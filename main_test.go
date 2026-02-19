package main

import (
	"os/exec"
	"testing"
)

func TestMainPackage(t *testing.T) {
	t.Run("builds successfully", func(t *testing.T) {
		cmd := exec.Command("go", "build", "-o", "/dev/null", ".")
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to build main package: %v", err)
		}
	})
}

func TestMainFunction(t *testing.T) {
	t.Run("main function exists", func(t *testing.T) {
	})
}
