package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"goproton/pkg/launcher"
	"goproton/pkg/lsfg_utils"
)

func (a *App) RunGame(opts launcher.LaunchOptions, showLogs bool) error {
	// DEBUG: Log what was received from frontend - write to file
	launcher.DebugLog("[APP.RunGame] Called with options:")
	launcher.DebugLog("[APP.RunGame]   LauncherPath: " + opts.LauncherPath)
	launcher.DebugLog("[APP.RunGame]   MainExecutablePath: " + opts.MainExecutablePath)
	launcher.DebugLog("[APP.RunGame]   HaveGameExe: " + fmt.Sprintf("%v", opts.HaveGameExe))
	launcher.DebugLog("[APP.RunGame]   EnableLsfgVk: " + fmt.Sprintf("%v", opts.EnableLsfgVk))
	launcher.DebugLog("[APP.RunGame] Full opts: " + fmt.Sprintf("%+v", opts))

	// Pre-flight check: Does the game exist?
	if _, err := os.Stat(opts.MainExecutablePath); os.IsNotExist(err) {
		return fmt.Errorf("game executable not found at: %s", opts.MainExecutablePath)
	}

	// CRITICAL: Normalize options before saving to ensure consistency
	// If HaveGameExe=false, MainExecutablePath MUST equal LauncherPath (launcher-only mode)
	// If HaveGameExe=true, MainExecutablePath must be different from LauncherPath (separate game exe)
	if !opts.HaveGameExe && opts.LauncherPath != "" {
		launcher.DebugLog("[APP.RunGame] NORMALIZE: HaveGameExe=false, enforcing MainExecutablePath=LauncherPath")
		opts.MainExecutablePath = opts.LauncherPath
	}

	// Auto-save config when launching
	_ = launcher.SaveGameConfig(opts)

	// If LSFG-VK enabled, ensure profile exists in config
	if opts.EnableLsfgVk {
		launcher.DebugLog("[APP.RunGame] LSFG-VK enabled, ensuring profile exists")
		configPath, err := lsfg_utils.GetConfigPath()
		if err == nil {
			// Check if profile exists
			_, _, err := lsfg_utils.FindProfileForGameAtPath(opts.MainExecutablePath, configPath)
			if err != nil {
				// Profile doesn't exist, create it
				launcher.DebugLog("[APP.RunGame] No profile found, creating one with current options")

				// If GPU is blank, use the first available GPU
				gpu := opts.LsfgGpu
				if gpu == "" {
					gpuList := launcher.GetListGpus()
					if len(gpuList) > 0 {
						gpu = gpuList[0]
						launcher.DebugLog("[APP.RunGame] GPU was blank, using first GPU: " + gpu)
					}
				}

				_ = lsfg_utils.SaveProfileToPath(opts.MainExecutablePath, configPath,
					parseMultiplier(opts.LsfgMultiplier),
					opts.LsfgPerfMode,
					opts.LsfgDllPath,
					gpu,
					opts.LsfgFlowScale,
					opts.LsfgPacing,
					opts.LsfgAllowFp16)
			}
		}
	}

	instanceName := "goproton-instance"

	// Try to find the binary relative to the current UI executable
	exeDir := "."
	if exe, err := os.Executable(); err == nil {
		exeDir = filepath.Dir(exe)
	}

	potentialPaths := []string{
		filepath.Join(exeDir, instanceName),                 // Same dir as UI
		filepath.Join(exeDir, "../../../bin", instanceName), // Dev mode: ui/build/bin -> root/bin
		"./bin/" + instanceName,                             // Current dir/bin
		"./" + instanceName,                                 // Current dir
		"/usr/bin/" + instanceName,                          // System installed (fallback)
	}

	foundPath := ""

	for _, p := range potentialPaths {
		if _, err := os.Stat(p); err == nil {
			foundPath = p
			break
		}
	}

	if foundPath == "" {
		return fmt.Errorf("instance manager (%s) not found. Checked: %v. Please run 'make build' first", instanceName, potentialPaths)
	}

	args := []string{
		"--game", opts.MainExecutablePath,
		"--launcher", opts.LauncherPath,
		"--prefix", opts.PrefixPath,
		"--proton-pattern", opts.ProtonPattern,
		"--proton-path", opts.ProtonPath,
	}
	if opts.EnableMangoHud {
		args = append(args, "--mango")
	}
	if opts.EnableGamemode {
		args = append(args, "--gamemode")
	}
	if opts.EnableLsfgVk {
		args = append(args, "--lsfg")
		args = append(args, "--lsfg-mult", opts.LsfgMultiplier)
		if opts.LsfgPerfMode {
			args = append(args, "--lsfg-perf")
		}
		if opts.LsfgDllPath != "" {
			args = append(args, "--lsfg-dll-path", opts.LsfgDllPath)
		}
	}
	if opts.EnableMemoryMin {
		args = append(args, "--memory-min")
		if opts.MemoryMinValue != "" {
			args = append(args, "--memory-min-value", opts.MemoryMinValue)
		}
	}
	if opts.EnableGamescope {
		args = append(args, "--gamescope")
		args = append(args, "--gs-w", opts.GamescopeW)
		args = append(args, "--gs-h", opts.GamescopeH)
		args = append(args, "--gs-r", opts.GamescopeR)
	}
	if !showLogs {
		args = append(args, "--logs=false")
	}
	cmd := exec.Command(foundPath, args...)

	// Start Detached
	if err := cmd.Start(); err != nil {
		fmt.Printf("!!! CRITICAL ERROR: Failed to start instance binary (%s): %v\n", foundPath, err)
		return fmt.Errorf("failed to start instance manager: %w", err)
	}
	go cmd.Process.Release()

	return nil
}

func (a *App) GetAllGames() ([]GameInfo, error) {
	configs, err := launcher.ListGameConfigs()
	if err != nil {
		return nil, err
	}

	games := make([]GameInfo, 0)
	for _, cfg := range configs {
		name := filepath.Base(cfg.MainExecutablePath)
		// Strip extension
		name = strings.TrimSuffix(name, filepath.Ext(name))

		// Clean path for comparison
		cleanedPath := filepath.Clean(cfg.MainExecutablePath)
		if abs, err := filepath.Abs(cleanedPath); err == nil {
			cleanedPath = abs
		}

		games = append(games, GameInfo{
			Name:   name,
			Path:   cleanedPath,
			Config: cfg,
		})
	}
	return games, nil
}

func (a *App) GetRunningSessions() ([]RunningSession, error) {
	// Use pgrep to find goproton-instance processes
	// Try full name first
	out, _ := exec.Command("pgrep", "goproton-instance").Output()
	if len(out) == 0 {
		out, _ = exec.Command("pgrep", "goproton-instan").Output() // Fallback to truncated
	}

	pids := strings.Split(strings.TrimSpace(string(out)), "\n")
	sessions := make([]RunningSession, 0)

	for _, pidStr := range pids {
		if pidStr == "" {
			continue
		}
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Read /proc/[pid]/cmdline for robust argument parsing
		cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
		content, err := os.ReadFile(cmdlinePath)
		if err != nil {
			continue
		}

		// cmdline uses null bytes to separate arguments
		args := strings.Split(string(content), "\x00")
		gamePath := ""
		for i, arg := range args {
			if arg == "--game" && i+1 < len(args) {
				gamePath = args[i+1]
				break
			}
		}

		if gamePath != "" {
			name := filepath.Base(gamePath)
			name = strings.TrimSuffix(name, filepath.Ext(name))

			// Clean path for comparison
			cleanedPath := filepath.Clean(gamePath)
			if abs, err := filepath.Abs(cleanedPath); err == nil {
				cleanedPath = abs
			}

			sessions = append(sessions, RunningSession{
				Pid:      pid,
				GamePath: cleanedPath,
				GameName: name,
			})
		}
	}
	return sessions, nil
}

func (a *App) KillSession(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(os.Interrupt) // Try SIGINT first for clean exit
}

func (a *App) RunPrefixTool(prefixPath, toolName, protonPattern string) error {
	opts := launcher.LaunchOptions{
		MainExecutablePath: toolName,
		PrefixPath:         prefixPath,
		ProtonPattern:      protonPattern,
	}
	cmdArgs, env := launcher.BuildCommand(opts)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = env
	return cmd.Start()
}

func (a *App) GetConfig(exePath string) (*launcher.LaunchOptions, error) {
	return launcher.LoadGameConfig(exePath)
}

func (a *App) SavePrefixConfig(prefixName string, opts launcher.LaunchOptions) error {
	return launcher.SavePrefixConfig(prefixName, opts)
}

func (a *App) LoadPrefixConfig(prefixName string) (*launcher.LaunchOptions, error) {
	return launcher.LoadPrefixConfig(prefixName)
}

// parseMultiplier converts a string multiplier to int
func parseMultiplier(mult string) int {
	var val int = 2
	if mult != "" {
		if _, err := fmt.Sscanf(mult, "%d", &val); err != nil {
			val = 2 // default
		}
	}
	return val
}
