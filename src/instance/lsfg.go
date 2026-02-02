package main

import (
	"fmt"
	"os"
	"path/filepath"
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

