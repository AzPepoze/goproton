package launcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	LsfgRepo = "PancakeTAS/lsfg-vk"
)

type UtilsStatus struct {
	IsLsfgInstalled bool   `json:"isLsfgInstalled"`
	LsfgVersion     string `json:"lsfgVersion"`
}

type SystemToolsStatus struct {
	HasGamescope bool `json:"hasGamescope"`
	HasMangoHud  bool `json:"hasMangoHud"`
	HasGameMode  bool `json:"hasGameMode"`
}

func GetSystemToolsStatus() SystemToolsStatus {
	return SystemToolsStatus{
		HasGamescope: isCommandAvailable("gamescope"),
		HasMangoHud:  isCommandAvailable("mangohud"),
		HasGameMode:  isCommandAvailable("gamemoderun"),
	}
}

func GetUtilsStatus() UtilsStatus {
	return UtilsStatus{
		IsLsfgInstalled: IsLsfgInstalled(),
		LsfgVersion:     "1.0.0",
	}
}

func IsLsfgInstalled() bool {
	vulkanDir := "/usr/share/vulkan/implicit_layer.d"
	_, err := os.Stat(filepath.Join(vulkanDir, "VkLayer_LSFGVK_frame_generation.json"))
	return err == nil
}

func InstallLsfg(onProgress func(int, string)) error {
	onProgress(0, "Fetching release info from GitHub...")
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases", LsfgRepo))
	if err != nil {
		return fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	var releases []struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return fmt.Errorf("failed to decode releases: %w", err)
	}

	var downloadURL, assetName string
	found := false

	// Search through releases (newest first)
	for _, release := range releases {
		if found {
			break
		}
		// Look for compatible assets in decreasing preference order
		for _, asset := range release.Assets {
			name := strings.ToLower(asset.Name)
			// Prefer x86_64 tar.zst, then linux tar.xz
			if (strings.Contains(name, "x86_64") && strings.HasSuffix(name, ".tar.zst")) ||
				(strings.Contains(name, "linux") && strings.HasSuffix(name, ".tar.xz")) {
				downloadURL = asset.BrowserDownloadURL
				assetName = asset.Name
				found = true
				break
			}
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("lsfg-vk suitable linux asset not found")
	}

	onProgress(5, fmt.Sprintf("Downloading %s...", assetName))
	ext := ".tar.xz"
	if strings.HasSuffix(assetName, ".tar.zst") {
		ext = ".tar.zst"
	}
	tmpFile := filepath.Join(os.TempDir(), "lsfg-vk-dl"+ext)

	err = downloadFileWithProgress(downloadURL, tmpFile, func(current, total int64) {
		percent := float64(current) / float64(total) * 80.0
		onProgress(5+int(percent), "Downloading...")
	})
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpFile)

	onProgress(85, "Extracting files...")
	extractTmp, err := os.MkdirTemp("", "lsfg-extract")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(extractTmp)

	extractCmd := []string{"-xf", tmpFile, "-C", extractTmp}
	if strings.HasSuffix(tmpFile, ".tar.zst") {
		extractCmd = []string{"--use-compress-program=unzstd", "-xf", tmpFile, "-C", extractTmp}
	}

	cmd := exec.Command("tar", extractCmd...)
	if _, err := cmd.CombinedOutput(); err != nil {
		cmd = exec.Command("tar", "-xf", tmpFile, "-C", extractTmp)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("extraction failed: %s", string(output))
		}
	}

	onProgress(88, "Installing to system directories (requires sudo)...")

	// Copy entire lsfg directory contents to /usr using pkexec (handles glob expansion)
	shellCmd := fmt.Sprintf("cp -r %s/* /usr", extractTmp)
	cpCmd := exec.Command("pkexec", "sh", "-c", shellCmd)
	if output, err := cpCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install lsfg: %v, output: %s", err, string(output))
	}

	onProgress(100, "Installation complete!")
	return nil
}

func UninstallLsfg() error {
	vulkanDir := "/usr/share/vulkan/implicit_layer.d"
	libDir := "/usr/lib"

	// Remove manifest
	manifest := filepath.Join(vulkanDir, "VkLayer_LSFGVK_frame_generation.json")
	exec.Command("pkexec", "rm", "-f", manifest).Run()

	// Remove library
	lib := filepath.Join(libDir, "liblsfg-vk-layer.so")
	exec.Command("pkexec", "rm", "-f", lib).Run()

	return nil
}

type progressWriter struct {
	total, current int64
	onProgress     func(current, total int64)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.current += int64(n)
	pw.onProgress(pw.current, pw.total)
	return n, nil
}

func downloadFileWithProgress(url string, dest string, onProgress func(current, total int64)) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "GoProton-App")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	pw := &progressWriter{total: resp.ContentLength, onProgress: onProgress}
	_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
	return err
}
