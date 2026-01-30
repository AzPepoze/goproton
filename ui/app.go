package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"go-proton/pkg/launcher"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ScanProtonVersions() ([]launcher.ProtonTool, error) {
	return launcher.GetProtonTools()
}

func (a *App) RunGame(opts launcher.LaunchOptions, showLogs bool) error {
	// Pre-flight check: Does the game exist?
	if _, err := os.Stat(opts.GamePath); os.IsNotExist(err) {
		return fmt.Errorf("game executable not found at: %s", opts.GamePath)
	}

	// Auto-save config when launching
	_ = launcher.SaveGameConfig(opts)

	instanceBin := "goproton-instance"
	potentialPaths := []string{
		"./" + instanceBin,
		"../" + instanceBin,
		"/usr/bin/" + instanceBin,
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
	args := []string{
		"--game", opts.GamePath,
		"--prefix", opts.PrefixPath,
		"--proton-pattern", opts.ProtonPattern,
		"--proton-path", opts.ProtonPath,
	}
	if opts.EnableMangoHud { args = append(args, "--mango") }
	if opts.EnableGamemode { args = append(args, "--gamemode") }
	if opts.EnableLsfgVk {
		args = append(args, "--lsfg")
		args = append(args, "--lsfg-mult", opts.LsfgMultiplier)
		if opts.LsfgPerfMode { args = append(args, "--lsfg-perf") }
		if opts.LsfgDllPath != "" { args = append(args, "--lsfg-dll-path", opts.LsfgDllPath) }
	}
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
		fmt.Printf("!!! CRITICAL ERROR: Failed to start instance binary (%s): %v\n", foundPath, err)
		return fmt.Errorf("failed to start instance manager: %w", err)
	}
	go cmd.Process.Release()
	return nil
}

func (a *App) RunPrefixTool(prefixPath, toolName, protonPattern string) error {
	opts := launcher.LaunchOptions{
		GamePath:      toolName,
		PrefixPath:    prefixPath,
		ProtonPattern: protonPattern,
	}
	cmdArgs, env := launcher.BuildCommand(opts)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = env
	return cmd.Start()
}

func (a *App) PickFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Game Executable",
		Filters: []runtime.FileFilter{
			{DisplayName: "Executables (*.exe)", Pattern: "*.exe"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
}

func (a *App) PickFolder() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Prefix Directory",
	})
}

func (a *App) GetConfig(exePath string) (*launcher.LaunchOptions, error) {
	return launcher.LoadGameConfig(exePath)
}

func (a *App) ListPrefixes() ([]string, error) {
	return launcher.ListPrefixes()
}

func (a *App) CreatePrefix(name string) error {
	return launcher.CreatePrefix(name)
}

func (a *App) GetPrefixBaseDir() string {
	return launcher.GetPrefixBaseDir()
}

func (a *App) GetUtilsStatus() launcher.UtilsStatus {
	return launcher.GetUtilsStatus()
}

func (a *App) GetSystemToolsStatus() launcher.SystemToolsStatus {
	return launcher.GetSystemToolsStatus()
}

func (a *App) InstallLsfg() error {
	return launcher.InstallLsfgWithLog(func(msg string) {
		runtime.EventsEmit(a.ctx, "lsfg-install-progress", msg)
	})
}

// DetectLosslessDll tries to find Lossless.dll in common Steam paths

func (a *App) DetectLosslessDll() string {

	home, _ := os.UserHomeDir()

	// Common Steam paths

	paths := []string{

		filepath.Join(home, ".steam/root/steamapps/common/Lossless Scaling/Lossless.dll"),

		filepath.Join(home, ".local/share/Steam/steamapps/common/Lossless Scaling/Lossless.dll"),

	}

	for _, p := range paths {

		if _, err := os.Stat(p); err == nil {

			return p

		}

	}

	return ""

}



// PickFileCustom opens a file dialog with custom filters

func (a *App) PickFileCustom(title string, filters []runtime.FileFilter) (string, error) {

	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{

		Title:   title,

		Filters: filters,

	})

}



// UninstallLsfg

func (a *App) UninstallLsfg() error {

	return launcher.UninstallLsfg()

}

func (a *App) CleanupProcesses() error {
	commands := []string{
		"umu-run",
		"pressure-vessel",
		"gamescopereaper",
		"steam-runtime-launcher-service",
		"srt-bwrap",
		"reaper",
	}
	for _, cmd := range commands {
		// Use standard pkill (SIGTERM) to allow clean exit
		_ = exec.Command("pkill", "-f", cmd).Run()
	}
	return nil
}
