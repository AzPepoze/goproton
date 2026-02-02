package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var debugLogFile *os.File

// InitDebugLog initializes the debug log file
func InitDebugLog() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	logPath := filepath.Join(homeDir, "GoProton", "debug.log")
	os.MkdirAll(filepath.Dir(logPath), 0755)

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		debugLogFile = f
		DebugLog("=== DEBUG LOG STARTED ===")
	}
}

// DebugLog writes a message to both console and debug log file
func DebugLog(msg string) {
	if debugLogFile == nil {
		InitDebugLog()
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] %s\n", timestamp, msg)

	// Write to console
	fmt.Print(line)

	// Write to file
	if debugLogFile != nil {
		debugLogFile.WriteString(line)
		debugLogFile.Sync()
	}
}
