package launcher

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type LaunchOptions struct {
	GamePath        string
	PrefixPath      string
	ProtonPattern   string
	ProtonPath      string
	EnableGamescope bool
	GamescopeW      string
	GamescopeH      string
	GamescopeR      string
	EnableMangoHud  bool
	EnableGamemode  bool
	EnableLsfgVk    bool
	LsfgMultiplier  string
	LsfgPerfMode    bool
	LsfgDllPath     string
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func GetBaseDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "GoProton")
}

func GetConfigDir() string {
	return filepath.Join(GetBaseDir(), "config", "executables")
}

func GetPrefixBaseDir() string {
	return filepath.Join(GetBaseDir(), "Prefixes")
}

func GetConfigPath(exePath string) string {
	h := sha1.New()
	h.Write([]byte(exePath))
	filename := hex.EncodeToString(h.Sum(nil)) + ".json"
	return filepath.Join(GetConfigDir(), filename)
}

func SaveGameConfig(opts LaunchOptions) error {
	path := GetConfigPath(opts.GamePath)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(opts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadGameConfig(exePath string) (*LaunchOptions, error) {
	path := GetConfigPath(exePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var opts LaunchOptions
	if err := json.Unmarshal(data, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func ListPrefixes() ([]string, error) {
	base := GetPrefixBaseDir()
	if err := os.MkdirAll(base, 0755); err != nil {
		return []string{}, nil
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		return []string{}, nil
	}
	prefixes := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			prefixes = append(prefixes, entry.Name())
		}
	}
	if len(prefixes) == 0 {
		defaultPath := filepath.Join(base, "Default")
		_ = os.MkdirAll(defaultPath, 0755)
		prefixes = append(prefixes, "Default")
	}
	return prefixes, nil
}

func CreatePrefix(name string) error {
	path := filepath.Join(GetPrefixBaseDir(), name)
	return os.MkdirAll(path, 0755)
}

func BuildCommand(opts LaunchOptions) ([]string, []string) {
	var cmdArgs []string
	env := os.Environ()
	home, _ := os.UserHomeDir()

	// Environment Variables (Based on original working code)
	env = append(env, fmt.Sprintf("WINEPREFIX=%s", ExpandPath(opts.PrefixPath)))
	env = append(env, fmt.Sprintf("UMU_PROTON_PATTERN=%s", opts.ProtonPattern))
	
	if opts.ProtonPath != "" {
		env = append(env, fmt.Sprintf("PROTONPATH=%s", ExpandPath(opts.ProtonPath)))
	}

			// Prepare LSFG Environment

			var lsfgEnv []string

			if opts.EnableLsfgVk {

				lsfgEnv = append(lsfgEnv, "ENABLE_LSFG=1", "DISABLE_LSFGVK=0", "LSFG_LEGACY=1")

				lsfgEnv = append(lsfgEnv, "VK_LOADER_LAYERS_ENABLE=VK_LAYER_LSFGVK_frame_generation")

				lsfgEnv = append(lsfgEnv, "VK_INSTANCE_LAYERS=VK_LAYER_LSFGVK_frame_generation")

				lsfgEnv = append(lsfgEnv, "LSFG_HUD=1", "LSFG_LOG_LEVEL=debug")

				

				lsfgDir := filepath.Join(home, "GoProton", "tools", "lsfg")

				lsfgEnv = append(lsfgEnv, "VK_ADD_LAYER_PATH="+lsfgDir)

				// Force mount the tools directory into the container

				lsfgEnv = append(lsfgEnv, "PRESSURE_VESSEL_FILESYSTEMS_RW="+lsfgDir)

				

				if opts.LsfgMultiplier != "" { lsfgEnv = append(lsfgEnv, "LSFG_MULTIPLIER="+opts.LsfgMultiplier) }

				if opts.LsfgPerfMode { lsfgEnv = append(lsfgEnv, "LSFG_PERFORMANCE_MODE=1") }

				

				dllPath := ExpandPath(opts.LsfgDllPath)

				if decoded, err := filepath.EvalSymlinks(dllPath); err == nil { dllPath = decoded }

				if dllPath != "" { lsfgEnv = append(lsfgEnv, "LSFG_DLL_PATH="+dllPath) }

			}

		
		// Logic for MangoHud and LSFG Environment placement
	if !opts.EnableGamescope {
		if opts.EnableMangoHud { env = append(env, "MANGOHUD=1") }
		if opts.EnableLsfgVk { env = append(env, lsfgEnv...) }
	}

	// Wrapper: Gamemode
	if opts.EnableGamemode && isCommandAvailable("gamemoderun") {
		cmdArgs = append(cmdArgs, "gamemoderun")
	}

	// Wrapper: Gamescope
	if opts.EnableGamescope && isCommandAvailable("gamescope") {
		cmdArgs = append(cmdArgs, "gamescope")
		if opts.GamescopeW != "" { cmdArgs = append(cmdArgs, "-W", opts.GamescopeW) }
		if opts.GamescopeH != "" { cmdArgs = append(cmdArgs, "-H", opts.GamescopeH) }
		if opts.GamescopeR != "" { cmdArgs = append(cmdArgs, "-r", opts.GamescopeR) }
		cmdArgs = append(cmdArgs, "--", "env")
		
		if opts.EnableMangoHud { cmdArgs = append(cmdArgs, "MANGOHUD=1") }
		if opts.EnableLsfgVk { cmdArgs = append(cmdArgs, lsfgEnv...) }
	}

	// Base Command: umu-run
	cmdArgs = append(cmdArgs, "umu-run", opts.GamePath)

	return cmdArgs, env
}

func RunGameWithLog(cmdArgs []string, env []string, onLog func(string), onExit func()) (*os.Process, error) {
	if len(cmdArgs) == 0 { return nil, fmt.Errorf("empty command") }
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil { return nil, err }
	readPipe := func(r io.Reader, prefix string) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() { onLog(fmt.Sprintf("%s%s", prefix, scanner.Text())) }
	}
	go readPipe(stdout, "")
	go readPipe(stderr, "")
	go func() {
		cmd.Wait()
		onLog("\n--- Process Exited ---")
		if onExit != nil { onExit() }
	}()
	return cmd.Process, nil
}

func StopProcessGroup(proc *os.Process) error {
	if proc == nil { return nil }
	pgid := -proc.Pid
	_ = syscall.Kill(pgid, syscall.SIGINT)
	time.Sleep(200 * time.Millisecond)
	return syscall.Kill(pgid, syscall.SIGTERM)
}

func FormatCommandForDisplay(cmdArgs []string, opts LaunchOptions) string {
	var sb strings.Builder
	sb.WriteString("WINEPREFIX=" + opts.PrefixPath + " ")
	sb.WriteString("UMU_PROTON_PATTERN=" + opts.ProtonPattern + " ")
	if opts.EnableMangoHud { sb.WriteString("MANGOHUD=1 ") }
	sb.WriteString(strings.Join(cmdArgs, " "))
	return sb.String()
}
