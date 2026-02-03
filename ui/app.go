package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"goproton/pkg/launcher"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type LsfgProfileData struct {
	Name            string  `json:"name"`
	Multiplier      int     `json:"multiplier"`
	PerformanceMode bool    `json:"performanceMode"`
	GPU             string  `json:"gpu"`
	FlowScale       float32 `json:"flowScale"`
	Pacing          string  `json:"pacing"`
	DllPath         string  `json:"dllPath"`
	AllowFp16       bool    `json:"allowFp16"`
}

type GameInfo struct {
	Name     string                 `json:"name"`
	Path     string                 `json:"path"`
	Icon     string                 `json:"icon"`
	Config   launcher.LaunchOptions `json:"config"`
	IsRecent bool                   `json:"isRecent"`
}

type RunningSession struct {
	Pid      int    `json:"pid"`
	GamePath string `json:"gamePath"`
	GameName string `json:"gameName"`
}

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetInitialLauncherPath() string {
	path := os.Getenv("GOPROTON_LAUNCHER_PATH")
	return path
}

func (a *App) GetInitialGamePath() string {
	path := os.Getenv("GOPROTON_GAME_PATH")
	return path
}

func (a *App) GetShouldEditLsfg() bool {
	shouldEdit := os.Getenv("GOPROTON_EDIT_LSFG")
	return shouldEdit == "1"
}

func (a *App) CloseWindow() {
	runtime.Quit(a.ctx)
	// Kill the process to ensure it doesn't stay running
	os.Exit(0)
}

func (a *App) GetExeIcon(exePath string) string {
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return ""
	}

	tempDir, err := os.MkdirTemp("", "goproton-icon-*")
	if err != nil {
		return ""
	}
	defer os.RemoveAll(tempDir)

	// Try using wrestool first
	cmd := exec.Command("wrestool", "-x", "--output="+tempDir, exePath)
	if err := cmd.Run(); err == nil {
		matches, _ := filepath.Glob(filepath.Join(tempDir, "*.ico"))
		if len(matches) > 0 && tryReadIcon(matches[0]) != "" {
			return tryReadIcon(matches[0])
		}
	}

	// Fallback: try icoextract
	cmd = exec.Command("icoextract", exePath, filepath.Join(tempDir, "icon.ico"))
	if err := cmd.Run(); err == nil {
		icoPath := filepath.Join(tempDir, "icon.ico")
		if icon := tryReadIcon(icoPath); icon != "" {
			return icon
		}

		// Check for alternative names
		matches, _ := filepath.Glob(filepath.Join(tempDir, "*.ico"))
		if len(matches) > 0 {
			return tryReadIcon(matches[0])
		}
	}

	return ""
}

func tryReadIcon(icoPath string) string {
	data, err := os.ReadFile(icoPath)
	if err != nil || len(data) == 0 {
		return ""
	}
	return "data:image/x-icon;base64," + base64.StdEncoding.EncodeToString(data)
}

func (a *App) GetSystemToolsStatus() launcher.SystemToolsStatus {
	return launcher.GetSystemToolsStatus()
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

// GetTotalRam returns the total system RAM in GB
func (a *App) GetTotalRam() (int, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var memTotalKb int
	// Read line by line
	var line string
	for {
		_, err := fmt.Fscanf(file, "%s %d kB\n", &line, &memTotalKb)
		if err != nil || line == "MemTotal:" {
			break
		}
	}

	if memTotalKb == 0 {
		return 0, fmt.Errorf("failed to parse MemTotal")
	}

	// Convert to GB
	return memTotalKb / 1024 / 1024, nil
}
