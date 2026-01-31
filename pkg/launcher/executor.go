package launcher

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func buildLsfgEnv(opts LaunchOptions) []string {
	if !opts.EnableLsfgVk {
		return nil
	}

	home, _ := os.UserHomeDir()
	lsfgDir := filepath.Join(home, "GoProton", "tools", "lsfg")
	lsfgLayers := "VK_LAYER_LSFGVK_frame_generation:VK_LAYER_LS_frame_generation"

	env := []string{
		"ENABLE_LSFG=1",
		"DISABLE_LSFGVK=0",
		"LSFG_LEGACY=1",
		"VK_LOADER_LAYERS_ENABLE=" + lsfgLayers,
		"VK_INSTANCE_LAYERS=" + lsfgLayers,
		"LSFG_HUD=1",
		"LSFG_LOG_LEVEL=debug",
		"VK_ADD_LAYER_PATH=" + lsfgDir,
		"PRESSURE_VESSEL_FILESYSTEMS_RW=" + lsfgDir,
	}

	if opts.LsfgMultiplier != "" {
		env = append(env, "LSFG_MULTIPLIER="+opts.LsfgMultiplier)
	}
	if opts.LsfgPerfMode {
		env = append(env, "LSFG_PERFORMANCE_MODE=1")
	}

	dllPath := ExpandPath(opts.LsfgDllPath)
	if decoded, err := filepath.EvalSymlinks(dllPath); err == nil {
		dllPath = decoded
	}
	if dllPath != "" {
		env = append(env, "LSFG_DLL_PATH="+dllPath)
	}

	return env
}

func buildBaseEnv(opts LaunchOptions) []string {
	baseEnv := []string{
		fmt.Sprintf("WINEPREFIX=%s", ExpandPath(opts.PrefixPath)),
		fmt.Sprintf("UMU_PROTON_PATTERN=%s", opts.ProtonPattern),
	}

	if opts.ProtonPath != "" {
		baseEnv = append(baseEnv, fmt.Sprintf("PROTONPATH=%s", ExpandPath(opts.ProtonPath)))
	}

	return baseEnv
}

func BuildCommand(opts LaunchOptions) ([]string, []string) {
	var cmdArgs []string
	env := append(os.Environ(), buildBaseEnv(opts)...)

	lsfgEnv := buildLsfgEnv(opts)

	if !opts.EnableGamescope {
		if opts.EnableMangoHud {
			env = append(env, "MANGOHUD=1")
		}
		env = append(env, lsfgEnv...)
	}

	if opts.EnableGamemode && isCommandAvailable("gamemoderun") {
		cmdArgs = append(cmdArgs, "gamemoderun")
	}

	if opts.EnableGamescope && isCommandAvailable("gamescope") {
		cmdArgs = append(cmdArgs, "gamescope")
		if opts.GamescopeW != "" {
			cmdArgs = append(cmdArgs, "-W", opts.GamescopeW)
		}
		if opts.GamescopeH != "" {
			cmdArgs = append(cmdArgs, "-H", opts.GamescopeH)
		}
		if opts.GamescopeR != "" {
			cmdArgs = append(cmdArgs, "-r", opts.GamescopeR)
		}
		cmdArgs = append(cmdArgs, "--", "env")

		if opts.EnableMangoHud {
			env = append(env, "MANGOHUD=1")
		}
		env = append(env, lsfgEnv...)
	}

	cmdArgs = append(cmdArgs, "umu-run", opts.GamePath)

	if opts.CustomArgs != "" {
		args := strings.Fields(opts.CustomArgs)
		cmdArgs = append(cmdArgs, args...)
	}

	if opts.EnableMemoryMin && opts.MemoryMinValue != "" && isCommandAvailable("systemd-run") {
		wrappedArgs := []string{"systemd-run", "--user", "--scope", fmt.Sprintf("-pMemoryMin=%s", opts.MemoryMinValue), "--"}
		cmdArgs = append(wrappedArgs, cmdArgs...)
	}

	return cmdArgs, env
}

func RunGameWithLog(cmdArgs []string, env []string, onLog func(string), onExit func()) (*os.Process, error) {
	if len(cmdArgs) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	readPipe := func(r io.Reader, prefix string) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			onLog(fmt.Sprintf("%s%s", prefix, scanner.Text()))
		}
	}

	go readPipe(stdout, "")
	go readPipe(stderr, "")

	go func() {
		_ = cmd.Wait()
		onLog("\n--- Process Exited ---")
		if onExit != nil {
			onExit()
		}
	}()

	return cmd.Process, nil
}

func StopProcessGroup(proc *os.Process) error {
	if proc == nil {
		return nil
	}
	pgid := -proc.Pid
	_ = syscall.Kill(pgid, syscall.SIGINT)
	time.Sleep(200 * time.Millisecond)
	return syscall.Kill(pgid, syscall.SIGTERM)
}

func FormatCommandForDisplay(cmdArgs []string, opts LaunchOptions) string {
	var sb strings.Builder
	if opts.EnableMemoryMin && opts.MemoryMinValue != "" {
		sb.WriteString(fmt.Sprintf("[MemMin:%s] ", opts.MemoryMinValue))
	}
	sb.WriteString("WINEPREFIX=" + opts.PrefixPath + " ")
	sb.WriteString("UMU_PROTON_PATTERN=" + opts.ProtonPattern + " ")
	if opts.EnableMangoHud {
		sb.WriteString("MANGOHUD=1 ")
	}
	sb.WriteString(strings.Join(cmdArgs, " "))
	return sb.String()
}
