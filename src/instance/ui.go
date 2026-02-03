package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/getlantern/systray"
)

// onReady is called when systray is ready
func onReady(logPath string) {
	exeName := filepath.Base(gamePath)
	exeNameClean := strings.TrimSuffix(exeName, filepath.Ext(exeName))
	launcherName := filepath.Base(launcherPath)

	systray.SetTitle("GoProton: " + exeNameClean)
	systray.SetTooltip("Running: " + exeNameClean)

	sendNotification("GoProton", "Launching "+exeNameClean+" ("+launcherName+")...")

	// Setup UI
	mStatus := systray.AddMenuItem("Running: "+exeNameClean, "")
	mStatus.Disable()
	systray.AddSeparator()

	if lsfg {
		setupLsfgMenu()
	}

	mKill := systray.AddMenuItem("End Process", "Stop this game")

	// Start game
	opts := buildLaunchOptions()
	cmdArgs, env := BuildCommand(opts)

	logGameStartup(cmdArgs)

	gameCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	gameCmd.Env = env
	gameCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	gameCmd.Stdout = logFileHandle
	gameCmd.Stderr = logFileHandle

	if err := gameCmd.Start(); err != nil {
		log.Printf("!!! ERROR: Failed to start game: %v\n", err)
		sendNotification("Launch Error", "Failed to start "+exeNameClean+" ("+launcherName+"): "+err.Error())
		systray.Quit()
		return
	}

	// Show logs in terminal if enabled
	if showLogs {
		startLogTerminal(logPath, gameCmd.Process.Pid)
	}

	// Periodically trim log file to keep it manageable (queue: last 500 lines)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			_ = trimLogFile(logPath, 500)
		}
	}()

	// Setup kill handler
	go func() {
		<-mKill.ClickedCh
		if gameCmd.Process != nil {
			StopProcessGroup(gameCmd.Process)
		}
	}()

	// Wait for game to exit
	go func() {
		err := gameCmd.Wait()
		log.Printf("Game process exited with: %v\n", err)

		if err != nil {
			sendNotification("Process Exited", fmt.Sprintf("%s exited with error: %v", exeNameClean, err))
		}

		time.Sleep(1 * time.Second)
		systray.Quit()
	}()
}

// setupLsfgMenu configures the LSFG-VK configuration menu item
func setupLsfgMenu() {
	mLsfgEdit := systray.AddMenuItem("Edit LSFG-VK Config", "Open LSFG-VK configuration")
	log.Printf("LSFG menu item created, waiting for clicks...")

	go func() {
		for {
			log.Printf("LSFG menu handler: waiting for click...")
			<-mLsfgEdit.ClickedCh
			log.Printf("LSFG menu handler: click received!")

			profile, idx, err := FindLsfgProfileForGame(gamePath)
			if err != nil {
				log.Printf("LSFG menu handler: error finding profile: %v", err)
				sendNotification("LSFG-VK Config", fmt.Sprintf("Could not find profile for this game: %v", err))
				continue
			}
			log.Printf("Found LSFG profile for game: %s (index: %d)", profile.Name, idx)

			launchLsfgUI()
			log.Printf("LSFG menu handler: UI launch completed, ready for next click")
		}
	}()
}

// launchLsfgUI launches the goproton-ui in LSFG edit mode
func launchLsfgUI() {
	uiBinary := "goproton-ui"

	// Try to find UI binary in the same directory as this executable (packed scenario)
	if exePath, err := os.Executable(); err == nil {
		dir := filepath.Dir(exePath)

		// Check for goproton-ui in the same directory
		localBinary := filepath.Join(dir, "goproton-ui")
		if _, err := os.Stat(localBinary); err == nil {
			uiBinary = localBinary
			log.Printf("Found local goproton-ui binary: %s", localBinary)
		} else {
			// If not found locally, try parent directory (for electron resources structure)
			parentDir := filepath.Dir(dir)
			parentBinary := filepath.Join(parentDir, "goproton-ui")
			if _, err := os.Stat(parentBinary); err == nil {
				uiBinary = parentBinary
				log.Printf("Found goproton-ui in parent directory: %s", parentBinary)
			}
		}
	}

	uiCmd := exec.Command(uiBinary)
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOPROTON_GAME_PATH=%s", gamePath))
	env = append(env, "GOPROTON_EDIT_LSFG=1")
	uiCmd.Env = env
	uiCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	log.Printf("Launching goproton-ui with GOPROTON_GAME_PATH=%s and GOPROTON_EDIT_LSFG=1", gamePath)

	if err := uiCmd.Start(); err != nil {
		sendNotification("LSFG-VK Config", fmt.Sprintf("Failed to launch UI: %v", err))
		log.Printf("Error launching UI: %v", err)
	} else {
		log.Printf("UI launched successfully (PID: %d)", uiCmd.Process.Pid)
		go uiCmd.Process.Release()
	}
}
