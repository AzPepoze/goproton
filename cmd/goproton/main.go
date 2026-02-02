package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Find UI binary - prefer local bin first
	uiPath := "./bin/goproton-ui"
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

		if _, err := os.Stat(gamePath); os.IsNotExist(err) {
			fmt.Printf("Error: File not found: %s\n", gamePath)
			os.Exit(1)
		}

		absPath, err := filepath.Abs(gamePath)
		if err != nil {
			fmt.Printf("Error: Failed to resolve absolute path: %s\n", err)
			os.Exit(1)
		}

		os.Setenv("GOPROTON_LAUNCHER_PATH", absPath)
		fmt.Printf("Pre-selecting launcher path: %s\n", absPath)
	}

	// Launch UI
	cmd := exec.Command(uiPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	newEnv := []string{}
	for _, env := range os.Environ() {
		newEnv = append(newEnv, env)
	}

	cmd.Run()
}
