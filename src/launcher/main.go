package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	var uiBinary string

	// Try to find UI binary in the same directory as this executable (packed scenario)
	if exePath, err := os.Executable(); err == nil {
		dir := filepath.Dir(exePath)

		// Check for goproton-ui in the same directory
		localBinary := filepath.Join(dir, "goproton-ui")
		if _, err := os.Stat(localBinary); err == nil {
			uiBinary = localBinary
		}
	}

	// If not found locally, try parent directory (for electron resources structure)
	if uiBinary == "" {
		if exePath, err := os.Executable(); err == nil {
			dir := filepath.Dir(exePath)
			// Try parent directory (one level up)
			parentDir := filepath.Dir(dir)

			// Look for goproton-ui in parent
			uiBinary = filepath.Join(parentDir, "goproton-ui")
			if _, err := os.Stat(uiBinary); err != nil {
				// Not found, fallback to PATH
				uiBinary = "goproton-ui"
			}
		}
	}

	// Create the command to run goproton-ui with any passed arguments
	cmd := exec.Command(uiBinary, os.Args[1:]...)

	// Set up environment variables
	env := os.Environ()

	// If a launcher/game path was provided as first argument, pass it via env var
	if len(os.Args) > 1 {
		env = append(env, fmt.Sprintf("GOPROTON_LAUNCHER_PATH=%s", os.Args[1]))
	}

	// Check if running on Wayland and force XWayland
	if _, wayland := os.LookupEnv("WAYLAND_DISPLAY"); wayland {
		var filteredEnv []string
		for _, e := range env {
			// Remove Wayland and conflicting variables
			if !strings.HasPrefix(e, "GDK_BACKEND=") &&
				!strings.HasPrefix(e, "OZONE_PLATFORM=") &&
				!strings.HasPrefix(e, "ELECTRON_OZONE_PLATFORM_HINT=") &&
				!strings.HasPrefix(e, "XDG_SESSION_TYPE=") {
				filteredEnv = append(filteredEnv, e)
			}
		}
		env = filteredEnv

		// Set X11/XWayland-specific variables
		env = append(env, "OZONE_PLATFORM=x11")
		env = append(env, "GDK_BACKEND=x11")
		env = append(env, "ELECTRON_OZONE_PLATFORM_HINT=x11")
		env = append(env, "QT_QPA_PLATFORM=xcb")
		env = append(env, "XDG_SESSION_TYPE=x11")
	}

	env = append(env, "RUN_FROM_GOPROTON=true")

	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Run the UI
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to launch goproton-ui: %v\n", err)
		os.Exit(1)
	}
}
