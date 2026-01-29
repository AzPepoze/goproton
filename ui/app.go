package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"goproton-wails/launcher"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ScanProtonVersions scans for available Proton versions
func (a *App) ScanProtonVersions() ([]launcher.ProtonTool, error) {
	return launcher.GetProtonTools()
}

// RunGame launches the standalone instance manager (detached)
func (a *App) RunGame(opts launcher.LaunchOptions, showLogs bool) error {
	instanceBin := "goproton-instance"
	
	// Robust path finding
	potentialPaths := []string{
		"./" + instanceBin,          // Current dir (Production)
		"../" + instanceBin,         // Parent dir (Dev Mode)
		"/usr/bin/" + instanceBin,    // System path
	}

	foundPath := ""
	for _, p := range potentialPaths {
		if _, err := os.Stat(p); err == nil {
			foundPath = p
			break
		}
	}

	if foundPath == "" {
		return fmt.Errorf("instance manager (%s) not found. Please run 'make build' first", instanceBin)
	}

	// Prepare flags
	args := []string{
		"--game", opts.GamePath,
		"--prefix", opts.PrefixPath,
		"--proton-pattern", opts.ProtonPattern,
		"--proton-path", opts.ProtonPath,
	}

	if opts.EnableMangoHud { args = append(args, "--mango") }
	if opts.EnableGamemode { args = append(args, "--gamemode") }
	if opts.EnableGamescope {
		args = append(args, "--gamescope")
		args = append(args, "--gs-w", opts.GamescopeW)
		args = append(args, "--gs-h", opts.GamescopeH)
		args = append(args, "--gs-r", opts.GamescopeR)
	}
	if !showLogs { args = append(args, "--logs=false") }

	cmd := exec.Command(foundPath, args...)
	
	// Start Detached
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	// Release the process so it keeps running independently
	go cmd.Process.Release()
	
	return nil
}

// RunPrefixTool runs tools inside the prefix (can still be handled by umu-run directly or instance)
func (a *App) RunPrefixTool(prefixPath, toolName, protonPattern string) error {
	opts := launcher.LaunchOptions{
		GamePath:       toolName,
		PrefixPath:     prefixPath,
		ProtonPattern:  protonPattern,
	}
	// We run prefix tools through the instance manager too if we want a tray, 
	// but usually prefix tools don't need persistent management. 
	// Let's run it normally for now.
	cmdArgs, env := launcher.BuildCommand(opts)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = env
	return cmd.Start()
}

// PickFile opens a file dialog
func (a *App) PickFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Game Executable",
		Filters: []runtime.FileFilter{
			{DisplayName: "Executables (*.exe)", Pattern: "*.exe"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
}

// PickFolder opens a directory dialog
func (a *App) PickFolder() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Prefix Directory",
	})
}