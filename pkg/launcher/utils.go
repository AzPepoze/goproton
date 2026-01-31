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
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	lsfgDir := filepath.Join(home, "GoProton", "tools", "lsfg")
	entries, err := os.ReadDir(lsfgDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			return true
		}
	}
	return false
}

func InstallLsfgWithLog(onProgress func(int, string)) error {
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

	patterns := []struct {
		contains string
		suffix   string
	}{
		{"x86_64", ".tar.zst"},
		{"linux", ".tar.xz"},
	}

	for _, pattern := range patterns {
		if found {
			break
		}
		for _, release := range releases {
			for _, asset := range release.Assets {
				name := strings.ToLower(asset.Name)
				if strings.Contains(name, pattern.contains) && strings.HasSuffix(name, pattern.suffix) {
					downloadURL = asset.BrowserDownloadURL
					assetName = asset.Name
					found = true
					break
				}
			}
			if found {
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

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}
	lsfgDir := filepath.Join(home, "GoProton", "tools", "lsfg")

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

	onProgress(90, "Installing files to ~/GoProton/tools/lsfg...")
	if err := os.MkdirAll(lsfgDir, 0755); err != nil {
		return fmt.Errorf("failed to create tools dir: %w", err)
	}

	err = filepath.Walk(extractTmp, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		name := info.Name()
		var dstPath string

		if strings.HasSuffix(name, ".so") || strings.HasSuffix(name, ".json") {
			dstPath = filepath.Join(lsfgDir, name)
		} else if strings.HasPrefix(name, "lsfg-vk-") && !strings.Contains(name, ".tar") {
			dstPath = filepath.Join(lsfgDir, name)
		} else {
			return nil
		}

		return moveFile(srcPath, dstPath)
	})
	if err != nil {
		return fmt.Errorf("failed to install files: %w", err)
	}

	onProgress(95, "Fixing manifest...")
	if err := fixLsfgManifest(lsfgDir); err != nil {
		return fmt.Errorf("failed to fix manifest: %w", err)
	}

	onProgress(100, "Installation complete!")
	return nil
}

func fixLsfgManifest(lsfgDir string) error {
	entries, err := os.ReadDir(lsfgDir)
	if err != nil {
		return err
	}

	var originalJson string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			originalJson = filepath.Join(lsfgDir, entry.Name())
			break
		}
	}

	if originalJson == "" {
		return fmt.Errorf("no JSON manifest found")
	}

	jsonBytes, err := os.ReadFile(originalJson)
	if err != nil {
		return err
	}

	var installedSoPath string
	if _, err := os.Stat(filepath.Join(lsfgDir, "liblsfg-vk.so")); err == nil {
		installedSoPath = filepath.Join(lsfgDir, "liblsfg-vk.so")
	} else if _, err := os.Stat(filepath.Join(lsfgDir, "liblsfg-vk-layer.so")); err == nil {
		installedSoPath = filepath.Join(lsfgDir, "liblsfg-vk-layer.so")
	} else {
		return fmt.Errorf("no .so library found")
	}

	content := string(jsonBytes)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, "\"library_path\"") {
			lines[i] = fmt.Sprintf("        \"library_path\": \"%s\",", installedSoPath)
			break
		}
	}

	return os.WriteFile(originalJson, []byte(strings.Join(lines, "\n")), 0644)
}

func moveFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return os.Chmod(dst, 0755)
}

func UninstallLsfg() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	lsfgDir := filepath.Join(home, "GoProton", "tools", "lsfg")
	return os.RemoveAll(lsfgDir)
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
