package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"go-proton/pkg/launcher"

	"github.com/pelletier/go-toml/v2"
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

func (a *App) ScanProtonVersions() ([]launcher.ProtonTool, error) {
	return launcher.GetProtonTools()
}

func (a *App) RunGame(opts launcher.LaunchOptions, showLogs bool) error {
	// DEBUG: Log what was received from frontend - write to file
	launcher.DebugLog("[APP.RunGame] Called with options:")
	launcher.DebugLog("[APP.RunGame]   LauncherPath: " + opts.LauncherPath)
	launcher.DebugLog("[APP.RunGame]   GamePath: " + opts.GamePath)
	launcher.DebugLog("[APP.RunGame]   UseGameExe: " + fmt.Sprintf("%v", opts.UseGameExe))
	launcher.DebugLog("[APP.RunGame]   EnableLsfgVk: " + fmt.Sprintf("%v", opts.EnableLsfgVk))
	launcher.DebugLog("[APP.RunGame] Full opts: " + fmt.Sprintf("%+v", opts))

	// Pre-flight check: Does the game exist?
	if _, err := os.Stat(opts.GamePath); os.IsNotExist(err) {
		return fmt.Errorf("game executable not found at: %s", opts.GamePath)
	}

	// CRITICAL: Normalize options before saving to ensure consistency
	// If UseGameExe=false, GamePath MUST equal LauncherPath (launcher-only mode)
	// If UseGameExe=true, GamePath must be different from LauncherPath (separate game exe)
	if !opts.UseGameExe && opts.LauncherPath != "" {
		launcher.DebugLog("[APP.RunGame] NORMALIZE: UseGameExe=false, enforcing GamePath=LauncherPath")
		opts.GamePath = opts.LauncherPath
	}

	// Auto-save config when launching
	_ = launcher.SaveGameConfig(opts)

	// If LSFG-VK enabled, ensure profile exists in config
	if opts.EnableLsfgVk {
		launcher.DebugLog("[APP.RunGame] LSFG-VK enabled, ensuring profile exists")
		configPath, err := launcher.GetLsfgConfigPath()
		if err == nil {
			// Check if profile exists
			_, _, err := launcher.FindLsfgProfileForGameAtPath(opts.GamePath, configPath)
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

				_ = launcher.SaveLsfgProfileToPath(opts.GamePath, configPath,
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
		"--game", opts.GamePath,
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

	// Close window after starting game
	go func() {
		runtime.Quit(a.ctx)
	}()

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

func (a *App) SavePrefixConfig(prefixName string, opts launcher.LaunchOptions) error {
	return launcher.SavePrefixConfig(prefixName, opts)
}

func (a *App) LoadPrefixConfig(prefixName string) (*launcher.LaunchOptions, error) {
	return launcher.LoadPrefixConfig(prefixName)
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
	return launcher.InstallLsfg(func(percent int, msg string) {
		runtime.EventsEmit(a.ctx, "lsfg-install-progress", map[string]interface{}{
			"percent": percent,
			"message": msg,
		})
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

// GetListGpus returns a list of available GPUs on the system
func (a *App) GetListGpus() []string {
	return launcher.GetListGpus()
}

// PickFileCustom opens a file dialog with custom filters

func (a *App) PickFileCustom(title string, filters []runtime.FileFilter) (string, error) {

	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{

		Title: title,

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

func (a *App) GetProtonVariants() []launcher.ProtonVariant {
	return launcher.GetKnownVariants()
}

func (a *App) GetProtonReleases(variantID string) ([]launcher.GitHubRelease, error) {
	return launcher.FetchReleases(variantID)
}

func (a *App) InstallProtonVersion(url, version string) error {
	return launcher.InstallProton(url, version, func(percent int, msg string) {
		runtime.EventsEmit(a.ctx, "install-proton-progress", map[string]interface{}{
			"percent": percent,
			"message": msg,
		})
	})
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

// GetLsfgProfileForGame retrieves the LSFG profile for a specific game
// Returns nil if no profile exists (expected for unconfigured games)
func (a *App) GetLsfgProfileForGame(gamePath string) (*LsfgProfileData, error) {
	profile, _, err := launcher.FindLsfgProfileForGame(gamePath)
	if err != nil {
		// No profile found - return nil (not an error)
		// This is normal for games that haven't been configured yet
		return nil, nil
	}

	// Get global config for DLL and AllowFp16
	var dllPath string
	var allowFp16 bool
	configPath, err := launcher.GetLsfgConfigPath()
	if err == nil {
		data, err := os.ReadFile(configPath)
		if err == nil {
			var config launcher.LsfgConfigFile
			if err := toml.Unmarshal(data, &config); err == nil {
				dllPath = config.Global.DLL
				allowFp16 = config.Global.AllowFP16
			}
		}
	}

	return &LsfgProfileData{
		Name:            profile.Name,
		Multiplier:      profile.Multiplier,
		PerformanceMode: profile.PerformanceMode,
		GPU:             profile.GPU,
		FlowScale:       profile.FlowScale,
		Pacing:          profile.Pacing,
		DllPath:         dllPath,
		AllowFp16:       allowFp16,
	}, nil
}

// SaveLsfgProfile saves an LSFG profile for a specific game
func (a *App) SaveLsfgProfile(gamePath string, multiplier int, perfMode bool, dllPath, gpu, flowScale, pacing string, allowFp16 bool) error {
	// If GPU is blank, use the first available GPU
	if gpu == "" {
		gpuList := launcher.GetListGpus()
		if len(gpuList) > 0 {
			gpu = gpuList[0]
			launcher.DebugLog("SaveLsfgProfile() GPU was blank, using first GPU: " + gpu)
		}
	}

	// Get the config file path (defaults to ~/.config/lsfg-vk/conf.toml)
	configPath, err := launcher.GetLsfgConfigPath()
	if err != nil {
		return err
	}

	// Save to the config file
	return launcher.SaveLsfgProfileToPath(gamePath, configPath, multiplier, perfMode, dllPath, gpu, flowScale, pacing, allowFp16)
}

// RemoveProfile removes a profile from the lsfg-vk config
func (a *App) RemoveProfile(gamePath string) error {
	return launcher.RemoveProfileFromConfig(gamePath)
}

// EditLsfgConfigForGame returns the game path (used to trigger fullscreen editor in frontend)
func (a *App) EditLsfgConfigForGame(gamePath string) error {
	// Verify the game has a profile
	_, _, err := launcher.FindLsfgProfileForGame(gamePath)
	return err
}
