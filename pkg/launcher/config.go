package launcher

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func loadConfig(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func saveConfig(path string, v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
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
	path := GetConfigPath(opts.GamePath)
	return saveConfig(path, opts)
}

func LoadGameConfig(exePath string) (*LaunchOptions, error) {
	path := GetConfigPath(exePath)
	var opts LaunchOptions
	if err := loadConfig(path, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}
