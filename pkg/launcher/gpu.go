package launcher

import (
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
