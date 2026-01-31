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

	"go-proton/pkg/launcher"

	"github.com/getlantern/systray"
)

var (
	gamePath, prefixPath, protonPath, protonPattern string
	mango, gamemode, gamescope, lsfg, lsfgPerf, memoryMin     bool
	gsW, gsH, gsR, lsfgMult, lsfgDllPath, memoryMinValue            string
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

func sendNotification(title, message string) {
	_ = exec.Command("notify-send", "-a", "GoProton", title, message).Run()
}

func main() {
	flag.StringVar(&gamePath, "game", "", "Path to the game executable")
	flag.StringVar(&prefixPath, "prefix", "", "Path to the WINEPREFIX")
	flag.StringVar(&protonPath, "proton-path", "", "Full path to the Proton tool")
	flag.StringVar(&protonPattern, "proton-pattern", "", "Proton pattern for UMU")
	flag.BoolVar(&mango, "mango", false, "Enable MangoHud")
	flag.BoolVar(&gamemode, "gamemode", false, "Enable GameMode")
	flag.BoolVar(&gamescope, "gamescope", false, "Enable Gamescope")
	flag.BoolVar(&lsfg, "lsfg", false, "Enable LSFG-VK")
	flag.StringVar(&lsfgMult, "lsfg-mult", "2", "LSFG Multiplier")
	flag.BoolVar(&lsfgPerf, "lsfg-perf", false, "Enable LSFG Performance Mode")
	flag.StringVar(&lsfgDllPath, "lsfg-dll-path", "", "Path to Lossless.dll")
	flag.BoolVar(&memoryMin, "memory-min", false, "Enable Memory Protection (min RAM)")
	flag.StringVar(&memoryMinValue, "memory-min-value", "", "Memory Protection Value (e.g. 4G)")
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

	sendNotification("GoProton", "Launching "+exeName+"...")

	mStatus := systray.AddMenuItem("Running: "+exeName, "")
	mStatus.Disable()
	systray.AddSeparator()
	mKill := systray.AddMenuItem("End Process", "Stop this game")

	opts := launcher.LaunchOptions{
		GamePath: gamePath, PrefixPath: prefixPath, ProtonPattern: protonPattern, ProtonPath: protonPath,
		EnableMangoHud: mango, EnableGamemode: gamemode, EnableGamescope: gamescope,
		GamescopeW: gsW, GamescopeH: gsH, GamescopeR: gsR,
		EnableLsfgVk: lsfg, LsfgMultiplier: lsfgMult, LsfgPerfMode: lsfgPerf, LsfgDllPath: lsfgDllPath,
		EnableMemoryMin: memoryMin, MemoryMinValue: memoryMinValue,
	}
	cmdArgs, env := launcher.BuildCommand(opts)

	log.Printf("--- EXECUTION START ---")
	log.Printf("COMMAND: %s", strings.Join(cmdArgs, " "))
	log.Printf("ENABLED FEATURES:")
	if mango { log.Printf("  [+] MangoHud") }
	if gamemode { log.Printf("  [+] GameMode") }
	if gamescope { log.Printf("  [+] Gamescope (%sx%s@%s)", gsW, gsH, gsR) }
	if lsfg { log.Printf("  [+] LSFG-VK (x%s, PerfMode:%v)", lsfgMult, lsfgPerf) }
	if memoryMin { log.Printf("  [+] Memory Protection (Min: %s)", memoryMinValue) }

	log.Printf("CUSTOM ENVIRONMENT VARIABLES:")
	sysEnv := os.Environ()
	for _, e := range env {
		isSystem := false
		for _, s := range sysEnv {
			if e == s {
				isSystem = true
				break
			}
		}
		if !isSystem {
			log.Printf("  %s", e)
		}
	}
	log.Printf("-----------------------")

	// Force sync to disk so the terminal tail sees these headers immediately
	if logFileHandle != nil {
		_ = logFileHandle.Sync()
	}

	gameCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	gameCmd.Env = env
	gameCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	gameCmd.Stdout = logFileHandle
	gameCmd.Stderr = logFileHandle

	if err := gameCmd.Start(); err != nil {
		log.Printf("!!! ERROR: Failed to start game: %v\n", err)
		sendNotification("Launch Error", "Failed to start "+exeName+": "+err.Error())
		systray.Quit()
		return
	}

	var termCmd *exec.Cmd
	if showLogs {
		term := findTerminal()
		if term != "" {
			// Use cat/tail and filter out annoying spam logs
			filterExpr := "setpriority|vk_xwayland_wait_ready|vk_wsi_force_swapchain"
			tailCmd := fmt.Sprintf("sleep 0.2; cat %s | grep -E -v '%s'; tail --pid %d -f %s | grep -E -v '%s'; echo; echo '---------------------------------------'; echo 'Process finished. Press Enter to close...'; read", logPath, filterExpr, gameCmd.Process.Pid, logPath, filterExpr)
			var cmd *exec.Cmd
			if strings.Contains(term, "kitty") || strings.Contains(term, "alacritty") {
				cmd = exec.Command(term, "--", "bash", "-c", tailCmd)
			} else {
				cmd = exec.Command(term, "-e", "bash", "-c", tailCmd)
			}
			termCmd = cmd
			termCmd.Start()
		}
	}

	go func() {
		<-mKill.ClickedCh
		if gameCmd.Process != nil {
			launcher.StopProcessGroup(gameCmd.Process)
		}
	}()

	go func() {
		err := gameCmd.Wait()
		log.Printf("Game process exited with: %v\n", err)

		if err != nil {
			sendNotification("Process Exited", fmt.Sprintf("%s exited with error: %v", exeName, err))
		}

		time.Sleep(1 * time.Second)
		systray.Quit()
	}()
}

func onExit() {
	if logFileHandle != nil {
		logFileHandle.Close()
	}
}
