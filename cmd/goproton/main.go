package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Find UI binary
	uiPath := "./goproton-ui"
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		uiPath = "./ui/build/bin/goproton-ui"
	}
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		uiPath = "/usr/bin/goproton-ui"
	}

	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		fmt.Println("Error: UI binary not found. Please run 'make build' first.")
		os.Exit(1)
	}

	// If a path argument is provided, set it as env variable for UI to pre-select
	if len(os.Args) > 1 {
		gamePath := os.Args[1]

		// Check if file exists
		if _, err := os.Stat(gamePath); os.IsNotExist(err) {
			fmt.Printf("Error: File not found: %s\n", gamePath)
			os.Exit(1)
		}

		os.Setenv("GOPROTON_GAME_PATH", gamePath)
	}

	// Launch UI
	cmd := exec.Command(uiPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Run()
}
