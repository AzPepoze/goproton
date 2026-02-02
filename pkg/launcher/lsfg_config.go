package launcher

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type LsfgGlobalConfig struct {
	Version   int    `toml:"version"`
	AllowFP16 bool   `toml:"allow_fp16"`
	DLL       string `toml:"dll"`
}

type LsfgConfigProfile struct {
	Name            string      `toml:"name"`
	ActiveIn        interface{} `toml:"active_in"` // Can be string or []string
	Multiplier      int         `toml:"multiplier"`
	PerformanceMode bool        `toml:"performance_mode"`
	GPU             string      `toml:"gpu"`
	FlowScale       float32     `toml:"flow_scale"`
	Pacing          string      `toml:"pacing"`
}

type LsfgConfigFile struct {
	Version  int                 `toml:"version"`
	Global   LsfgGlobalConfig    `toml:"global"`
	Profiles []LsfgConfigProfile `toml:"profile"`
}

// GetLsfgProfilePath returns the path to store LSFG profile for a game exe
// Uses format: GoProton/config/lsfg/exename-hash.toml
func GetLsfgProfilePath(gamePath string) string {
	h := sha1.New()
	h.Write([]byte(gamePath))
	hash := hex.EncodeToString(h.Sum(nil))[:8] // First 8 chars of hash
	exeName := filepath.Base(gamePath)
	// Remove .exe extension for readability
	baseName := strings.TrimSuffix(exeName, ".exe")
	baseName = strings.TrimSuffix(baseName, ".EXE")

	filename := baseName + "-" + hash + ".toml"
	configDir := filepath.Join(GetBaseDir(), "config", "lsfg")
	return filepath.Join(configDir, filename)
}

// GetLsfgConfigPath returns the path to the lsfg-vk config file
// lsfg-vk reads from ~/.config/lsfg-vk/conf.toml by default
func GetLsfgConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "lsfg-vk", "conf.toml"), nil
}

// FindLsfgProfileForGame finds the profile that applies to the given game exe
// FindLsfgProfileForGameAtPath finds the profile for a game at a specific config path
func FindLsfgProfileForGameAtPath(gamePath, configPath string) (*LsfgConfigProfile, int, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to read LSFG config: %w", err)
	}

	var config LsfgConfigFile
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, -1, fmt.Errorf("failed to parse LSFG config: %w", err)
	}

	exeName := strings.ToLower(filepath.Base(gamePath))

	// Find the matching profile
	for i, profile := range config.Profiles {
		if matchesProfile(exeName, profile.ActiveIn) {
			return &profile, i, nil
		}
	}

	return nil, -1, fmt.Errorf("no matching LSFG profile found for %s", exeName)
}

// FindLsfgProfileForGame finds the profile for a game using LSFG_CONFIG env var
func FindLsfgProfileForGame(gamePath string) (*LsfgConfigProfile, int, error) {
	configPath, err := GetLsfgConfigPath()
	if err != nil {
		return nil, -1, err
	}
	return FindLsfgProfileForGameAtPath(gamePath, configPath)
}

// matchesProfile checks if the exe matches the profile's active_in list
func matchesProfile(exeName string, activeIn interface{}) bool {
	exeName = strings.ToLower(exeName)

	switch v := activeIn.(type) {
	case string:
		return strings.EqualFold(v, exeName)
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok {
				if strings.EqualFold(s, exeName) {
					return true
				}
			}
		}
	}
	return false
}

// EditLsfgProfileForGame opens the LSFG config file with the profile highlighted
func EditLsfgProfileForGame(gamePath string) error {
	configPath, err := GetLsfgConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("LSFG config file not found at %s", configPath)
	}

	// Open with default editor
	return openFileWithEditor(configPath)
}

// openFileWithEditor opens a file with the system default editor
func openFileWithEditor(filePath string) error {
	// Try xdg-open first (Linux)
	if cmd := exec.Command("xdg-open", filePath); cmd != nil {
		return cmd.Start()
	}
	return fmt.Errorf("failed to open editor")
}

// SaveLsfgProfileToPath saves a profile to the LSFG-VK config file at a specific path
func SaveLsfgProfileToPath(gamePath, configPath string, multiplier int, perfMode bool, dllPath, gpu, flowScale, pacing string, allowFp16 bool) error {
	DebugLog("SaveLsfgProfileToPath() called for game: " + gamePath)
	DebugLog("SaveLsfgProfileToPath() saving to: " + configPath)

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Read existing config or create new one
	var config LsfgConfigFile
	if data, err := os.ReadFile(configPath); err == nil {
		if err := toml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse existing LSFG config: %w", err)
		}
	} else {
		// Create new config with version 2
		config = LsfgConfigFile{
			Version: 2,
			Global: LsfgGlobalConfig{
				Version:   2,
				AllowFP16: allowFp16,
				DLL:       dllPath,
			},
			Profiles: []LsfgConfigProfile{},
		}
	}

	// Update global settings
	config.Version = 2
	config.Global.Version = 2
	config.Global.DLL = dllPath
	config.Global.AllowFP16 = allowFp16

	// Extract just the exe name for active_in
	exeName := filepath.Base(gamePath)

	// Find if profile already exists for this exe
	found := false
	for i, profile := range config.Profiles {
		if matchesProfile(strings.ToLower(exeName), profile.ActiveIn) {
			// Update existing profile
			config.Profiles[i].Multiplier = multiplier
			config.Profiles[i].PerformanceMode = perfMode
			config.Profiles[i].GPU = gpu
			config.Profiles[i].FlowScale = parseFlowScale(flowScale)
			config.Profiles[i].Pacing = "none" // Always use 'none' pacing
			found = true
			break
		}
	}

	if !found {
		// Create new profile
		newProfile := LsfgConfigProfile{
			Name:            strings.TrimSuffix(exeName, filepath.Ext(exeName)),
			ActiveIn:        exeName, // Store just the exe name
			Multiplier:      multiplier,
			PerformanceMode: perfMode,
			GPU:             gpu,
			FlowScale:       parseFlowScale(flowScale),
			Pacing:          "none", // Always use 'none' pacing
		}
		config.Profiles = append(config.Profiles, newProfile)
	}

	// Write back to file
	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal LSFG config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write LSFG config: %w", err)
	}

	return nil
}

// SaveLsfgProfileToGlobal saves a profile to the lsfg-vk config file
// Edits the config file to add or update profile for the game
func SaveLsfgProfileToGlobal(gamePath string, multiplier int, perfMode bool, dllPath, gpu, flowScale, pacing string, allowFp16 bool) error {
	configPath, err := GetLsfgConfigPath()
	if err != nil {
		DebugLog("SaveLsfgProfileToGlobal() error getting config path: " + err.Error())
		return err
	}
	return SaveLsfgProfileToPath(gamePath, configPath, multiplier, perfMode, dllPath, gpu, flowScale, pacing, allowFp16)
}

// RemoveProfileFromConfig removes a profile for a game from the config file
func RemoveProfileFromConfig(gamePath string) error {
	configPath, err := GetLsfgConfigPath()
	if err != nil {
		return err
	}

	DebugLog("RemoveProfileFromConfig() called for game: " + gamePath)
	DebugLog("RemoveProfileFromConfig() config path: " + configPath)

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read LSFG config: %w", err)
	}

	var config LsfgConfigFile
	if err := toml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse LSFG config: %w", err)
	}

	// Extract just the exe name for matching
	exeName := filepath.Base(gamePath)

	// Find and remove the profile for this exe
	found := false
	for i, profile := range config.Profiles {
		if matchesProfile(strings.ToLower(exeName), profile.ActiveIn) {
			DebugLog("RemoveProfileFromConfig() found profile for " + exeName + ", removing it")
			// Remove profile by slicing
			config.Profiles = append(config.Profiles[:i], config.Profiles[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("no profile found for %s", exeName)
	}

	// Write back to file
	data, err = toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal LSFG config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write LSFG config: %w", err)
	}

	return nil
}

// parseFlowScale parses a flowScale string to float32
func parseFlowScale(flowScale string) float32 {
	var val float32 = 1.0
	if flowScale != "" {
		// Try to parse as float
		var f float64
		if _, err := fmt.Sscanf(flowScale, "%f", &f); err == nil {
			val = float32(f)
		}
	}
	return val
}
