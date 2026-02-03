package launcher

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func loadConfig(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Detect format based on file extension
	if strings.HasSuffix(path, ".toml") {
		return toml.Unmarshal(data, v)
	}
	return json.Unmarshal(data, v)
}

func saveConfig(path string, v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	var data []byte
	var err error

	// Detect format based on file extension
	if strings.HasSuffix(path, ".toml") {
		data, err = toml.Marshal(v)
	} else {
		data, err = json.MarshalIndent(v, "", "  ")
	}

	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SavePrefixConfig(prefixName string, opts LaunchOptions) error {
	path := GetPrefixConfigPath(prefixName)
	return saveConfig(path, opts)
}

func LoadPrefixConfig(prefixName string) (*LaunchOptions, error) {
	path := GetPrefixConfigPath(prefixName)
	var opts LaunchOptions

	if err := loadConfig(path, &opts); err != nil {
		return &LaunchOptions{
			LsfgMultiplier: "2",
			MemoryMinValue: "4G",
			GamescopeW:     "1920",
			GamescopeH:     "1080",
			GamescopeR:     "60",
		}, nil
	}
	return &opts, nil
}

func SaveGameConfig(opts LaunchOptions) error {
	// Use LauncherPath for config storage if available, otherwise fall back to MainExecutablePath
	configExePath := opts.MainExecutablePath
	if opts.LauncherPath != "" {
		configExePath = opts.LauncherPath
	}
	path := GetGameConfigFile(configExePath)
	return saveConfig(path, opts)
}

func LoadGameConfig(exePath string) (*LaunchOptions, error) {
	path := GetGameConfigFile(exePath)
	var opts LaunchOptions
	if err := loadConfig(path, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

// ListGameConfigs returns a list of all saved game configurations
func ListGameConfigs() ([]LaunchOptions, error) {
	configDir := GetConfigDir()
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return make([]LaunchOptions, 0), nil
	}

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	configs := make([]LaunchOptions, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		configPath := filepath.Join(configDir, entry.Name(), "config.json")
		var opts LaunchOptions
		if err := loadConfig(configPath, &opts); err == nil {
			configs = append(configs, opts)
		}
	}

	return configs, nil
}

// SaveLsfgProfile saves an LSFG profile for a game
func SaveLsfgProfile(gameName string, profile LsfgProfile) error {
	// Use the main executable path to compute the directory with hash
	configPath := GetConfigPath(profile.MainExecutablePath)

	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	// Save as lsfg_vk.toml in the executables directory structure
	profilePath := filepath.Join(configPath, "lsfg_vk.toml")
	return saveConfig(profilePath, profile)
}

// LoadLsfgProfile loads an LSFG profile for a game
func LoadLsfgProfile(gamePath string) (*LsfgProfile, error) {
	// Use the game path to compute the directory with hash
	configPath := GetConfigPath(gamePath)

	profilePath := filepath.Join(configPath, "lsfg_vk.toml")
	var profile LsfgProfile
	if err := loadConfig(profilePath, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}
