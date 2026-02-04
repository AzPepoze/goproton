package core

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// GetListGpus detects available GPUs on the system using vulkaninfo
func GetListGpus() []string {
	gpus := []string{}
	detected := make(map[string]bool)

	// Use vulkaninfo to get GPU information
	if vulkanGpus := detectVulkanGpus(); len(vulkanGpus) > 0 {
		for _, gpu := range vulkanGpus {
			if !detected[gpu] {
				gpus = append(gpus, gpu)
				detected[gpu] = true
			}
		}
	}

	return gpus
}

// GetGpuUsage returns the current GPU load percentage
func GetGpuUsage() string {
	// NVIDIA check
	if out, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits").Output(); err == nil {
		return strings.TrimSpace(string(out)) + "%"
	}

	// AMD/Intel check via sysfs
	// Loop through potential card slots (0-5) to find an active one
	for i := 0; i <= 5; i++ {
		p := fmt.Sprintf("/sys/class/drm/card%d/device/gpu_busy_percent", i)
		if data, err := os.ReadFile(p); err == nil {
			val := strings.TrimSpace(string(data))
			if val != "" {
				// If we find a non-zero value, we probably found the active GPU
				if val != "0" {
					return val + "%"
				}
			}
		}
		// Alternative for older Intel
		pIntel := fmt.Sprintf("/sys/class/drm/card%d/device/i915_gpu_busy100", i)
		if data, err := os.ReadFile(pIntel); err == nil {
			val := strings.TrimSpace(string(data))
			if val != "" {
				if val != "0" {
					return val + "%"
				}
			}
		}
	}
	return "0%"
}

// detectVulkanGpus uses vulkaninfo to get GPU information
func detectVulkanGpus() []string {
	cmd := exec.Command("vulkaninfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	var gpus []string
	lines := strings.Split(string(output), "\n")

	// Pattern to match: GPU id = X (GPU_NAME)
	// Capture everything between "GPU id = <number> (" and the last ")"
	gpuPattern := regexp.MustCompile(`GPU\s+id\s*=\s*\d+\s*\((.+)\)`)

	for _, line := range lines {
		matches := gpuPattern.FindStringSubmatch(line)
		if len(matches) >= 2 {
			gpu := strings.TrimSpace(matches[1])
			if len(gpu) > 0 && !contains(gpus, gpu) {
				gpus = append(gpus, gpu)
			}
		}
	}
	return gpus
}

// contains checks if a slice contains a string
func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
