package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// UtilsStatus returns LSFG installation and utility status
type UtilsStatus struct {
	IsLsfgInstalled bool   `json:"isLsfgInstalled"`
	LsfgVersion     string `json:"lsfgVersion"`
}

// GetShaderCacheSize calculates the size of the Steam shader cache
func GetShaderCacheSize() string {
	home, _ := os.UserHomeDir()
	paths := []string{
		filepath.Join(home, ".steam/root/steamapps/shadercache"),
		filepath.Join(home, ".local/share/Steam/steamapps/shadercache"),
	}

	var totalSize int64
	for _, p := range paths {
		filepath.Walk(p, func(_ string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				totalSize += info.Size()
			}
			return nil
		})
	}

	if totalSize == 0 {
		return "0 MB"
	}

	return fmt.Sprintf("%.1f GB", float64(totalSize)/(1024*1024*1024))
}

// ClearShaderCache removes all Steam shader cache files
func ClearShaderCache() error {
	home, _ := os.UserHomeDir()
	paths := []string{
		filepath.Join(home, ".steam/root/steamapps/shadercache"),
		filepath.Join(home, ".local/share/Steam/steamapps/shadercache"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			entries, _ := os.ReadDir(p)
			for _, entry := range entries {
				os.RemoveAll(filepath.Join(p, entry.Name()))
			}
		}
	}
	return nil
}
