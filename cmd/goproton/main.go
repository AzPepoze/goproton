package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Simple wrapper to run the Wails UI
	uiPath := "./goproton-ui"
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		uiPath = "./ui/build/bin/goproton-ui"
	}

	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		fmt.Println("Error: UI binary not found. Please run 'make build' first.")
		os.Exit(1)
	}

	cmd := exec.Command(uiPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
