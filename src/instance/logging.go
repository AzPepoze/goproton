package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// getLogPath returns the path for the current log file
func getLogPath() string {
	u, err := user.Current()
	if err != nil {
		return "/tmp/goproton-instance.log"
	}
	logDir := filepath.Join(u.HomeDir, "GoProton/logs")
	os.MkdirAll(logDir, 0755)
	cleanupLogs(logDir, 10)
	timestamp := time.Now().Format("20060102-150405")
	exeName := filepath.Base(gamePath)
	return filepath.Join(logDir, fmt.Sprintf("%s-%s.log", exeName, timestamp))
}

// trimLogFile keeps only the last maxLines lines of a log file
func trimLogFile(filePath string, maxLines int) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) <= maxLines {
		return nil
	}

	// Keep only the last maxLines
	trimmed := lines[len(lines)-maxLines:]
	trimmedData := strings.Join(trimmed, "\n")

	return os.WriteFile(filePath, []byte(trimmedData), 0666)
}

// cleanupLogs removes old log files, keeping only the most recent 'keep' files
func cleanupLogs(dir string, keep int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var files []os.FileInfo
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".log") {
			info, err := entry.Info()
			if err == nil {
				files = append(files, info)
			}
		}
	}
	if len(files) <= keep {
		return
	}
	sort.Slice(files, func(i, j int) bool { return files[i].ModTime().Before(files[j].ModTime()) })
	toDelete := len(files) - keep
	for i := 0; i < toDelete; i++ {
		os.Remove(filepath.Join(dir, files[i].Name()))
	}
}
