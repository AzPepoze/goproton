package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// LSFG Configuration Types

// LsfgGlobalConfig represents global LSFG-VK frame generation settings
type LsfgGlobalConfig struct {
	Version   int    `toml:"version"`
	AllowFP16 bool   `toml:"allow_fp16"`
	DLL       string `toml:"dll"`
}

// LsfgConfigProfile represents a per-game LSFG-VK frame generation profile
type LsfgConfigProfile struct {
	Name            string      `toml:"name"`
	ActiveIn        interface{} `toml:"active_in"` // Can be string or []string
	Multiplier      int         `toml:"multiplier"`
	PerformanceMode bool        `toml:"performance_mode"`
	GPU             string      `toml:"gpu"`
	FlowScale       float32     `toml:"flow_scale"`
	Pacing          string      `toml:"pacing"`
}

// LsfgConfigFile represents the complete LSFG-VK configuration file
type LsfgConfigFile struct {
	Version  int                 `toml:"version"`
	Global   LsfgGlobalConfig    `toml:"global"`
	Profiles []LsfgConfigProfile `toml:"profile"`
}

// LSFG Configuration Functions

// GetLsfgConfigPath returns the path to the lsfg-vk config file.
// By default lsfg-vk reads from ~/.config/lsfg-vk/conf.toml
func GetLsfgConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "lsfg-vk", "conf.toml"), nil
}

// FindLsfgProfileForGame finds the LSFG profile that applies to the given game path
func FindLsfgProfileForGame(gamePath string) (*LsfgConfigProfile, int, error) {
	configPath, err := GetLsfgConfigPath()
	if err != nil {
		return nil, -1, err
	}
	return findLsfgProfileForGameAtPath(gamePath, configPath)
}

// findLsfgProfileForGameAtPath finds the profile for a game at a specific config path
func findLsfgProfileForGameAtPath(gamePath, configPath string) (*LsfgConfigProfile, int, error) {
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

// matchesProfile checks if the executable name matches a profile's active_in configuration
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

// EnsureLsfgProfileExists creates a profile for the game if one doesn't already exist
// Uses provided options for initial configuration, or defaults if not provided
func EnsureLsfgProfileExists(gamePath string, opts LaunchOptions) error {
	configPath, err := GetLsfgConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get LSFG config path: %w", err)
	}

	// Try to find existing profile
	_, _, err = findLsfgProfileForGameAtPath(gamePath, configPath)
	if err == nil {
		// Profile already exists
		return nil
	}

	// Profile doesn't exist, create it
	DebugLog(fmt.Sprintf("Creating new LSFG profile for %s", gamePath))

	exeName := filepath.Base(gamePath)

	// Read existing config or create new one
	var config LsfgConfigFile
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config doesn't exist, create new
		config = LsfgConfigFile{
			Version:  2,
			Global:   LsfgGlobalConfig{Version: 2, AllowFP16: true, DLL: opts.LsfgDllPath},
			Profiles: []LsfgConfigProfile{},
		}
	} else {
		// Config exists, try to parse it
		if err := toml.Unmarshal(data, &config); err != nil {
			// Parsing failed, start fresh
			config = LsfgConfigFile{
				Version:  2,
				Global:   LsfgGlobalConfig{Version: 2, AllowFP16: true, DLL: opts.LsfgDllPath},
				Profiles: []LsfgConfigProfile{},
			}
		} else {
			// Successfully parsed - preserve existing global config, only update DLL if provided
			if opts.LsfgDllPath != "" {
				config.Global.DLL = opts.LsfgDllPath
			}
		}
	}

	// Parse multiplier from options
	multiplier := 2
	if opts.LsfgMultiplier != "" {
		if m, err := fmt.Sscanf(opts.LsfgMultiplier, "%d", &multiplier); m == 1 && err == nil {
			// Parsing succeeded
		}
	}

	// Create profile name from executable name (without extension)
	profileName := strings.TrimSuffix(exeName, filepath.Ext(exeName))

	// Parse flow scale from options
	flowScale := 1.0
	if opts.LsfgFlowScale != "" {
		if f, err := strconv.ParseFloat(opts.LsfgFlowScale, 32); err == nil {
			flowScale = f
		}
	}

	// Create new profile with values from options
	newProfile := LsfgConfigProfile{
		Name:            profileName,
		ActiveIn:        exeName,
		Multiplier:      multiplier,
		PerformanceMode: opts.LsfgPerfMode,
		GPU:             opts.LsfgGpu,
		FlowScale:       float32(flowScale),
		Pacing:          opts.LsfgPacing,
	}

	config.Profiles = append(config.Profiles, newProfile)

	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create LSFG config directory: %w", err)
	}

	// Write config
	data, err = toml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal LSFG config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write LSFG config: %w", err)
	}

	DebugLog(fmt.Sprintf("Successfully created LSFG profile for %s", exeName))
	return nil
}
