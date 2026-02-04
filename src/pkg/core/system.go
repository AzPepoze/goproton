package core

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// SystemToolsStatus returns whether system tools are available
type SystemToolsStatus struct {
	HasGamescope bool `json:"hasGamescope"`
	HasMangoHud  bool `json:"hasMangoHud"`
	HasGameMode  bool `json:"hasGameMode"`
}

// SystemInfo represents detailed system information
type SystemInfo struct {
	OS     string `json:"os"`
	Kernel string `json:"kernel"`
	CPU    string `json:"cpu"`
	GPU    string `json:"gpu"`
	RAM    string `json:"ram"`
	Driver string `json:"driver"`
}

// SystemUsage represents real-time system usage
type SystemUsage struct {
	CPU string `json:"cpu"`
	RAM string `json:"ram"`
	GPU string `json:"gpu"`
}

var (
	lastTotal int64
	lastIdle  int64
)

// GetSystemUsage returns real-time CPU % and RAM usage
func GetSystemUsage() SystemUsage {
	usage := SystemUsage{CPU: "0%", RAM: "0%", GPU: "0%"}

	// 1. CPU Usage
	if f, err := os.Open("/proc/stat"); err == nil {
		scanner := bufio.NewScanner(f)
		if scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				var total int64
				for i := 1; i < len(fields); i++ {
					val, _ := strconv.ParseInt(fields[i], 10, 64)
					total += val
				}
				idle, _ := strconv.ParseInt(fields[4], 10, 64)

				if lastTotal > 0 {
					totalDelta := total - lastTotal
					idleDelta := idle - lastIdle
					if totalDelta > 0 {
						cpuPercent := 100 * (float64(totalDelta) - float64(idleDelta)) / float64(totalDelta)
						usage.CPU = fmt.Sprintf("%.1f%%", cpuPercent)
					}
				}
				lastTotal = total
				lastIdle = idle
			}
		}
		f.Close()
	}

	// 2. RAM Usage
	if f, err := os.Open("/proc/meminfo"); err == nil {
		scanner := bufio.NewScanner(f)
		var total, avail int64
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "MemTotal:") {
				fmt.Sscanf(line, "MemTotal: %d kB", &total)
			} else if strings.HasPrefix(line, "MemAvailable:") {
				fmt.Sscanf(line, "MemAvailable: %d kB", &avail)
			}
			if total > 0 && avail > 0 {
				break
			}
		}
		f.Close()

		if total > 0 {
			used := total - avail
			usage.RAM = fmt.Sprintf("%.1f GB / %d GB (%.0f%%)", float64(used)/1024/1024, total/1024/1024, 100*float64(used)/float64(total))
		}
	}

	// 3. GPU Usage
	usage.GPU = GetGpuUsage()

	return usage
}

// GetSystemToolsStatus returns availability of system tools
func GetSystemToolsStatus() SystemToolsStatus {
	return SystemToolsStatus{
		HasGamescope: isCommandAvailable("gamescope"),
		HasMangoHud:  isCommandAvailable("mangohud"),
		HasGameMode:  isCommandAvailable("gamemoderun"),
	}
}

// GetSystemInfo collects system hardware and software details
func GetSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:     "Unknown",
		Kernel: "Unknown",
		CPU:    "Unknown",
		GPU:    "Unknown",
		RAM:    "Unknown",
		Driver: "Unknown",
	}

	// 1. OS Info
	if f, err := os.Open("/etc/os-release"); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				info.OS = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				break
			}
		}
		f.Close()
	}

	// 2. Kernel Info
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		info.Kernel = strings.TrimSpace(string(out))
	}

	// 3. CPU Info
	if out, err := exec.Command("sh", "-c", "lscpu | grep 'Model name' | cut -d':' -f2 | xargs").Output(); err == nil {
		info.CPU = strings.TrimSpace(string(out))
	}

	// 4. GPU Info
	gpus := GetListGpus()
	if len(gpus) > 0 {
		info.GPU = gpus[0]
	}

	// 5. RAM Info
	if f, err := os.Open("/proc/meminfo"); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			var memTotalKb int64
			line := scanner.Text()
			if strings.HasPrefix(line, "MemTotal:") {
				fmt.Sscanf(line, "MemTotal: %d kB", &memTotalKb)
				info.RAM = fmt.Sprintf("%d GB", memTotalKb/1024/1024)
				break
			}
		}
		f.Close()
	}

	// 6. Driver Info
	if out, err := exec.Command("sh", "-c", "glxinfo -B | grep 'OpenGL version string' | cut -d':' -f2 | xargs").Output(); err == nil {
		info.Driver = strings.TrimSpace(string(out))
	} else if out, err := exec.Command("sh", "-c", "vulkaninfo --summary | grep -m 1 'driverVersion' | awk '{print $3}'").Output(); err == nil {
		info.Driver = strings.TrimSpace(string(out))
	}

	return info
}
