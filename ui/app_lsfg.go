package main

import (
	"os"
	"path/filepath"

	"goproton/pkg/launcher"
	"goproton/pkg/lsfg_utils"

	"github.com/pelletier/go-toml/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetUtilsStatus() launcher.UtilsStatus {
	return launcher.UtilsStatus{
		IsLsfgInstalled: lsfg_utils.IsInstalled(),
		LsfgVersion:     lsfg_utils.GetVersion(),
	}
}

func (a *App) InstallLsfg() error {
	return lsfg_utils.Install(func(percent int, msg string) {
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

// UninstallLsfg
func (a *App) UninstallLsfg() error {
	return lsfg_utils.Uninstall(launcher.DebugLog)
}

// GetLsfgProfileForGame retrieves the LSFG profile for a specific game
// Returns nil if no profile exists (expected for unconfigured games)
func (a *App) GetLsfgProfileForGame(mainExePath string) (*LsfgProfileData, error) {
	profile, _, err := lsfg_utils.FindProfileForGame(mainExePath)
	if err != nil {
		// No profile found - return nil (not an error)
		// This is normal for games that haven't been configured yet
		return nil, nil
	}

	// Get global config for DLL and AllowFp16
	var dllPath string
	var allowFp16 bool
	configPath, err := lsfg_utils.GetConfigPath()
	if err == nil {
		data, err := os.ReadFile(configPath)
		if err == nil {
			var config lsfg_utils.ConfigFile
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
func (a *App) SaveLsfgProfile(mainExePath string, multiplier int, perfMode bool, dllPath, gpu, flowScale, pacing string, allowFp16 bool) error {
	// If GPU is blank, use the first available GPU
	if gpu == "" {
		gpuList := launcher.GetListGpus()
		if len(gpuList) > 0 {
			gpu = gpuList[0]
			launcher.DebugLog("SaveLsfgProfile() GPU was blank, using first GPU: " + gpu)
		}
	}

	// Get the config file path (defaults to ~/.config/lsfg-vk/conf.toml)
	configPath, err := lsfg_utils.GetConfigPath()
	if err != nil {
		return err
	}

	// Save to the config file
	return lsfg_utils.SaveProfileToPath(mainExePath, configPath, multiplier, perfMode, dllPath, gpu, flowScale, pacing, allowFp16)
}

// RemoveProfile removes a profile from the lsfg-vk config
func (a *App) RemoveProfile(mainExePath string) error {
	return lsfg_utils.RemoveProfileFromConfig(mainExePath)
}

// EditLsfgConfigForGame returns the game path (used to trigger fullscreen editor in frontend)
func (a *App) EditLsfgConfigForGame(mainExePath string) error {
	// Verify the game has a profile
	_, _, err := lsfg_utils.FindProfileForGame(mainExePath)
	return err
}
