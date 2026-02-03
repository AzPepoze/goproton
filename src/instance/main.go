package main

import (
	"log"
	"os"

	"github.com/getlantern/systray"
)

func main() {
	initFlags()

	if gamePath == "" {
		os.Exit(1)
	}

	logPath := getLogPath()
	var err error
	logFileHandle, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFileHandle)
		// Trim log file to last 500 lines to keep queue behavior
		_ = trimLogFile(logPath, 500)
	}

	systray.Run(func() { onReady(logPath) }, onExit)
}

// onExit is called when systray is closed
func onExit() {
	if logFileHandle != nil {
		logFileHandle.Close()
	}
}
