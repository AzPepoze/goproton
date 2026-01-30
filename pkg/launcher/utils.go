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
	home, _ := os.UserHomeDir()
	// Check in our dedicated tool directory
	jsonPath := filepath.Join(home, "GoProton", "tools", "lsfg", "VkLayer_LSFGVK_frame_generation.json")
	_, err := os.Stat(jsonPath)
	return err == nil
}

func InstallLsfgWithLog(onProgress func(string)) error {
	onProgress("Fetching release info from GitHub...")
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases", LsfgRepo))
	if err != nil { return err }
	defer resp.Body.Close()

	var releases []struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil { return err }

	var downloadURL, assetName string
	found := false
	for _, release := range releases {
		for _, asset := range release.Assets {
			name := strings.ToLower(asset.Name)
			// Target: lsfg-vk-x.x.x-linux.tar.xz
			if strings.Contains(name, "linux") && strings.HasSuffix(name, ".tar.xz") {
				downloadURL = asset.BrowserDownloadURL
				assetName = asset.Name
				found = true
				break
			}
		}
		if found { break }
	}

	if downloadURL == "" { return fmt.Errorf("lsfg-vk linux asset not found") }

	onProgress(fmt.Sprintf("Downloading %s...", assetName))
	tmpFile := filepath.Join(os.TempDir(), "lsfg-vk-dl.tar.xz")
	err = downloadFileWithProgress(downloadURL, tmpFile, func(current, total int64) {
		percent := float64(current) / float64(total) * 100
		onProgress(fmt.Sprintf("Downloading (%.1f%%)", percent))
	})
	if err != nil { return err }
	defer os.Remove(tmpFile)

	home, _ := os.UserHomeDir()
	lsfgDir := filepath.Join(home, "GoProton", "tools", "lsfg")
	onProgress("Extracting files...")

	extractTmp, _ := os.MkdirTemp("", "lsfg-extract")
	defer os.RemoveAll(extractTmp)

	cmd := exec.Command("tar", "-xf", tmpFile, "-C", extractTmp)
	if output, err := cmd.CombinedOutput(); err != nil { return fmt.Errorf("extraction failed: %s", string(output)) }

	onProgress("Installing files to ~/GoProton/tools/lsfg...")
	_ = os.MkdirAll(lsfgDir, 0755)

	// In the new version, files might be in subfolders. We'll find them and flatten or move them.
	err = filepath.Walk(extractTmp, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() { return nil }
		
		name := info.Name()
		var dstPath string
		
		if strings.HasSuffix(name, ".so") {
			dstPath = filepath.Join(lsfgDir, name)
		} else if strings.HasSuffix(name, ".json") {
			dstPath = filepath.Join(lsfgDir, name)
		} else if strings.HasPrefix(name, "lsfg-vk-") && !strings.HasSuffix(name, ".tar.xz") {
			dstPath = filepath.Join(lsfgDir, name)
		} else {
			return nil // Skip other files like desktop or icons for now to keep it simple
		}
		
		return moveFile(srcPath, dstPath)
	})
	if err != nil { return fmt.Errorf("failed to copy files: %w", err) }

	installedJsonPath := filepath.Join(lsfgDir, "VkLayer_LSFGVK_frame_generation.json")
	installedSoPath := filepath.Join(lsfgDir, "liblsfg-vk-layer.so")

	if _, err := os.Stat(installedJsonPath); err == nil {
		onProgress("Fixing library path in JSON to absolute path...")
		jsonBytes, err := os.ReadFile(installedJsonPath)
		if err == nil {
			content := string(jsonBytes)
			// Replace with absolute path in our tool directory
			newContent := strings.ReplaceAll(content, "\"library_path\": \"liblsfg-vk-layer.so\"", fmt.Sprintf("\"library_path\": \"%s\"", installedSoPath))
			// Just in case it's different
			if newContent == content {
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.Contains(line, "\"library_path\"") {
						lines[i] = fmt.Sprintf("        \"library_path\": \"%s\",", installedSoPath)
						break
					}
				}
				newContent = strings.Join(lines, "\n")
			}
			_ = os.WriteFile(installedJsonPath, []byte(newContent), 0644)
		}
	}

	onProgress("Installation complete!")
	return nil
}

func moveFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil { return err }
	defer in.Close()
	_ = os.MkdirAll(filepath.Dir(dst), 0755)
	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil { return err }
	defer out.Close()
	_, err = io.Copy(out, in)
	_ = os.Chmod(dst, 0755)
	return err
}

func UninstallLsfg() error {
	home, _ := os.UserHomeDir()
	lsfgDir := filepath.Join(home, "GoProton", "tools", "lsfg")
	return os.RemoveAll(lsfgDir)
}

type progressWriter struct {
	total, current int64
	onProgress func(current, total int64)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.current += int64(n)
	pw.onProgress(pw.current, pw.total)
	return n, nil
}

func downloadFileWithProgress(url string, dest string, onProgress func(current, total int64)) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "GoProton-App")
	resp, err := client.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { return fmt.Errorf("server returned %s", resp.Status) }
	out, err := os.Create(dest)
	if err != nil { return err }
	defer out.Close()
	pw := &progressWriter{total: resp.ContentLength, onProgress: onProgress}
	_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
	return err
}