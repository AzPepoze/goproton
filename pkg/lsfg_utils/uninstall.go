package lsfg_utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// filesToRemove lists all files installed by LSFG-VK v2.0.0
var filesToRemove = []string{
	"/usr/share/vulkan/implicit_layer.d/VkLayer_LSFGVK_frame_generation.json",
	"/usr/lib/liblsfg-vk-layer.so",
	"/usr/share/icons/hicolor/256x256/apps/gay.pancake.lsfg-vk-ui.png",
	"/usr/share/applications/gay.pancake.lsfg-vk-ui.desktop",
	"/usr/bin/lsfg-vk-cli",
	"/usr/bin/lsfg-vk-ui",
}

// Uninstall removes all LSFG-VK files from the system
func Uninstall(onLog func(string)) error {
	onLog("[Uninstall] Starting LSFG-VK removal...")

	files := strings.Join(filesToRemove, " ")
	cmd := exec.Command("pkexec", "sh", "-c", fmt.Sprintf("rm -f %s 2>&1", files))
	output, err := cmd.CombinedOutput()

	if err != nil {
		onLog(fmt.Sprintf("[Uninstall] Failed to remove files: %v", err))
		if len(output) > 0 {
			onLog(fmt.Sprintf("[Uninstall] Details: %s", string(output)))
		}
		return fmt.Errorf("uninstall failed: %w", err)
	}

	onLog("[Uninstall] Removal complete.")
	return nil
}
