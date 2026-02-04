package lsfg_utils

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ManifestPath = "/usr/share/vulkan/implicit_layer.d/VkLayer_LSFGVK_frame_generation.json"
)

// IsInstalled checks if LSFG-VK is installed
func IsInstalled() bool {
	_, err := os.Stat(ManifestPath)
	return err == nil
}

// GetVersion reads the manifest and returns version info
func GetVersion() string {
	data, err := os.ReadFile(ManifestPath)
	if err != nil {
		return "unknown"
	}

	var manifest map[string]interface{}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "unknown"
	}

	// Try to find version info in the manifest
	if meta, ok := manifest["manifest_version"]; ok {
		return fmt.Sprintf("%v", meta)
	}

	return "installed"
}

// GetStatus returns the current installation status
func GetStatus() Status {
	return Status{
		IsInstalled: IsInstalled(),
		Version:     GetVersion(),
	}
}
