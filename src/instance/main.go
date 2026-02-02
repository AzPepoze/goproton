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

	"github.com/getlantern/systray"
)

var (
	// Game and launcher paths
	gamePath      string
	launcherPath  string
	prefixPath    string
	protonPath    string
	protonPattern string

	// Feature flags
	mango     bool
	gamemode  bool
	gamescope bool
	lsfg      bool
	lsfgPerf  bool
	memoryMin bool
	showLogs  bool

	// Gamescope configuration
	gsW string
	gsH string
	gsR string

	// LSFG configuration
	lsfgMult    string
	lsfgDllPath string

	// Memory configuration
	memoryMinValue string

	// Logging
	logFileHandle *os.File
)

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

// trimLogFile keeps only the last 500 lines of a log file (queue behavior)
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

func sendNotification(title, message string) {
	_ = exec.Command("notify-send", "-a", "GoProton", title, message).Run()
}

func main() {
	flag.StringVar(&gamePath, "game", "", "Path to the game executable")
	flag.StringVar(&launcherPath, "launcher", "", "Path to the launcher executable")
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

func findTerminal() string {
	if t := os.Getenv("TERMINAL"); t != "" {
		parts := strings.Fields(t)
		if len(parts) > 0 {
			if p, err := exec.LookPath(parts[0]); err == nil {
				return p
			}
		}
	}
	terms := []string{"kitty", "alacritty", "gnome-terminal", "konsole", "xfce4-terminal", "xterm"}
	for _, t := range terms {
		if p, err := exec.LookPath(t); err == nil {
			return p
		}
	}
	return ""
}

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

	// Try to find local UI binary first
	if exePath, err := os.Executable(); err == nil {
		localBinary := filepath.Join(filepath.Dir(exePath), "goproton-ui")
		if _, err := os.Stat(localBinary); err == nil {
			uiBinary = localBinary
			log.Printf("Found local goproton-ui binary: %s", localBinary)
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

// buildLaunchOptions creates the launch options from command line flags
func buildLaunchOptions() LaunchOptions {
	return LaunchOptions{
		GamePath:        gamePath,
		LauncherPath:    launcherPath,
		PrefixPath:      prefixPath,
		ProtonPattern:   protonPattern,
		ProtonPath:      protonPath,
		EnableMangoHud:  mango,
		EnableGamemode:  gamemode,
		EnableGamescope: gamescope,
		GamescopeW:      gsW,
		GamescopeH:      gsH,
		GamescopeR:      gsR,
		EnableLsfgVk:    lsfg,
		LsfgMultiplier:  lsfgMult,
		LsfgPerfMode:    lsfgPerf,
		LsfgDllPath:     lsfgDllPath,
		EnableMemoryMin: memoryMin,
		MemoryMinValue:  memoryMinValue,
	}
}

// logGameStartup logs the command and enabled features
func logGameStartup(cmdArgs []string) {
	log.Printf("--- EXECUTION START ---")
	log.Printf("COMMAND: %s", strings.Join(cmdArgs, " "))
	log.Printf("ENABLED FEATURES:")

	if mango {
		log.Printf("  [+] MangoHud")
	}
	if gamemode {
		log.Printf("  [+] GameMode")
	}
	if gamescope {
		log.Printf("  [+] Gamescope (%sx%s@%s)", gsW, gsH, gsR)
	}
	if lsfg {
		log.Printf("  [+] LSFG-VK (x%s, PerfMode:%v)", lsfgMult, lsfgPerf)
	}
	if memoryMin {
		log.Printf("  [+] Memory Protection (Min: %s)", memoryMinValue)
	}

	log.Printf("-----------------------")

	// Sync to disk so terminal sees output immediately
	if logFileHandle != nil {
		_ = logFileHandle.Sync()
	}
}

// startLogTerminal opens a terminal window to display game logs
func startLogTerminal(logPath string, gamePID int) {
	term := findTerminal()
	if term == "" {
		return
	}

	filterExpr := "setpriority|vk_xwayland_wait_ready|vk_wsi_force_swapchain"
	tailCmd := fmt.Sprintf(
		"sleep 0.2; cat %s | grep -E -v '%s'; tail --pid %d -f %s | grep -E -v '%s'; echo; echo '---------------------------------------'; echo 'Process finished. Press Enter to close...'; read",
		logPath, filterExpr, gamePID, logPath, filterExpr,
	)

	var cmd *exec.Cmd
	if strings.Contains(term, "kitty") || strings.Contains(term, "alacritty") {
		cmd = exec.Command(term, "--", "bash", "-c", tailCmd)
	} else {
		cmd = exec.Command(term, "-e", "bash", "-c", tailCmd)
	}

	cmd.Start()
}

func onExit() {
	if logFileHandle != nil {
		logFileHandle.Close()
	}
}
