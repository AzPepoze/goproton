package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"go-proton/internal/launcher"

	"github.com/getlantern/systray"
)

var (
	gamePath, prefixPath, protonPath, protonPattern string
	mango, gamemode, gamescope                      bool
	gsW, gsH, gsR                                   string
	showLogs                                        bool
	logFileHandle                                   *os.File
)

func getLogPath() string {
	u, err := user.Current()
	if err != nil { return "/tmp/goproton-instance.log" }
	logDir := filepath.Join(u.HomeDir, "GoProton/logs")
	os.MkdirAll(logDir, 0755)
	cleanupLogs(logDir, 10)
	timestamp := time.Now().Format("20060102-150405")
	exeName := filepath.Base(gamePath)
	return filepath.Join(logDir, fmt.Sprintf("%s-%s.log", exeName, timestamp))
}

func cleanupLogs(dir string, keep int) {
	entries, err := os.ReadDir(dir)
	if err != nil { return }
	var files []os.FileInfo
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".log") {
			info, err := entry.Info()
			if err == nil { files = append(files, info) }
		}
	}
	if len(files) <= keep { return }
	sort.Slice(files, func(i, j int) bool { return files[i].ModTime().Before(files[j].ModTime()) })
	toDelete := len(files) - keep
	for i := 0; i < toDelete; i++ { os.Remove(filepath.Join(dir, files[i].Name())) }
}

func main() {
	flag.StringVar(&gamePath, "game", "", "Path to the game executable")
	flag.StringVar(&prefixPath, "prefix", "", "Path to the WINEPREFIX")
	flag.StringVar(&protonPath, "proton-path", "", "Full path to the Proton tool")
	flag.StringVar(&protonPattern, "proton-pattern", "", "Proton pattern for UMU")
	flag.BoolVar(&mango, "mango", false, "Enable MangoHud")
	flag.BoolVar(&gamemode, "gamemode", false, "Enable GameMode")
	flag.BoolVar(&gamescope, "gamescope", false, "Enable Gamescope")
	flag.StringVar(&gsW, "gs-w", "1920", "Width")
	flag.StringVar(&gsH, "gs-h", "1080", "Height")
	flag.StringVar(&gsR, "gs-r", "60", "Refresh Rate")
	flag.BoolVar(&showLogs, "logs", true, "Show terminal logs")
	flag.Parse()

	if gamePath == "" { os.Exit(1) }

	logPath := getLogPath()
	var err error
	logFileHandle, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFileHandle)
	}

	log.Printf("Starting instance for: %s\n", gamePath)
	systray.Run(func() { onReady(logPath) }, onExit)
}

func findTerminal() string {
	if t := os.Getenv("TERMINAL"); t != "" {
		parts := strings.Fields(t)
		if len(parts) > 0 {
			if p, err := exec.LookPath(parts[0]); err == nil { return p }
		}
	}
	terms := []string{"kitty", "alacritty", "gnome-terminal", "konsole", "xfce4-terminal", "xterm"}
	for _, t := range terms {
		if p, err := exec.LookPath(t); err == nil { return p }
	}
	return ""
}

func onReady(logPath string) {
	exeName := filepath.Base(gamePath)
	systray.SetTitle("GoProton: " + exeName)
	systray.SetTooltip("Running: " + exeName)

	mStatus := systray.AddMenuItem("Running: "+exeName, "")
	mStatus.Disable()
	systray.AddSeparator()
	mKill := systray.AddMenuItem("End Process", "Stop this game")

	opts := launcher.LaunchOptions{
		GamePath: gamePath, PrefixPath: prefixPath, ProtonPattern: protonPattern, ProtonPath: protonPath,
		EnableMangoHud: mango, EnableGamemode: gamemode, EnableGamescope: gamescope,
		GamescopeW: gsW, GamescopeH: gsH, GamescopeR: gsR,
	}
	cmdArgs, env := launcher.BuildCommand(opts)

	// 1. Start the actual Game Process (Directly managed by Go)
	gameCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	gameCmd.Env = env
	gameCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	
	// Pipe game output to our log file
	gameCmd.Stdout = logFileHandle
	gameCmd.Stderr = logFileHandle

	if err := gameCmd.Start(); err != nil {
		log.Printf("Failed to start game: %v\n", err)
		systray.Quit()
		return
	}

	// 2. If showLogs is requested, open a terminal that just follows the log file
	var termCmd *exec.Cmd
	if showLogs {
		term := findTerminal()
		if term != "" {
			// terminal -e tail -f logPath
			// We add a small delay to ensure log file exists and has content
			tailCmd := fmt.Sprintf("sleep 0.5; tail -f %s", logPath)
			termCmd = exec.Command(term, "-e", "bash", "-c", tailCmd)
			termCmd.Start()
		}
	}

	// Handle Graceful Stop from Tray
	go func() {
		<-mKill.ClickedCh
		log.Println("Graceful stop requested from tray")
		if gameCmd.Process != nil {
			launcher.StopProcessGroup(gameCmd.Process)
		}
	}()

	// Wait for Game Process Exit
	go func() {
		err := gameCmd.Wait()
		log.Printf("Game process exited with: %v\n", err)
		
		// Cleanup terminal if it was opened
		if termCmd != nil && termCmd.Process != nil {
			_ = termCmd.Process.Kill()
		}
		
systray.Quit()
	}()
}

func onExit() {
	if logFileHandle != nil {
		logFileHandle.Close()
	}
}
