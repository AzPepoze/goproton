package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// isCommandAvailable checks if a command is available in the system PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// buildLsfgEnv constructs environment variables for LSFG-VK if enabled
func buildLsfgEnv(opts LaunchOptions) []string {
	if !opts.EnableLsfgVk {
		return nil
	}

	DebugLog("buildLsfgEnv() called")
	// lsfg-vk reads from ~/.config/lsfg-vk/conf.toml by default
	return []string{}
}

// buildBaseEnv constructs base environment variables for game execution
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

// BuildCommand constructs the complete command and environment for game execution
// It returns the command arguments and environment variables based on the launch options
func BuildCommand(opts LaunchOptions) ([]string, []string) {
	DebugLog("BuildCommand() called")
	DebugLog("  LauncherPath: " + opts.LauncherPath)
	DebugLog("  GamePath: " + opts.GamePath)
	DebugLog("  UseGameExe: " + fmt.Sprintf("%v", opts.UseGameExe))
	DebugLog("  EnableLsfgVk: " + fmt.Sprintf("%v", opts.EnableLsfgVk))

	var cmdArgs []string
	env := append(os.Environ(), buildBaseEnv(opts)...)
	lsfgEnv := buildLsfgEnv(opts)

	// Add options if not using gamescope
	if !opts.EnableGamescope {
		if opts.EnableMangoHud {
			env = append(env, "MANGOHUD=1")
		}
		env = append(env, lsfgEnv...)
	}

	// Add gamemode if available
	if opts.EnableGamemode && isCommandAvailable("gamemoderun") {
		cmdArgs = append(cmdArgs, "gamemoderun")
	}

	// Add gamescope if enabled
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

	cmdArgs = append(cmdArgs, "umu-run")

	// Resolve executable path
	exePath := opts.LauncherPath
	if exePath == "" {
		DebugLog("LauncherPath is empty, falling back to GamePath")
		exePath = opts.GamePath
	}
	DebugLog("Final exePath: " + exePath)
	cmdArgs = append(cmdArgs, exePath)

	// Add custom arguments
	if opts.CustomArgs != "" {
		args := strings.Fields(opts.CustomArgs)
		cmdArgs = append(cmdArgs, args...)
	}

	// Add memory protection wrapper if enabled
	if opts.EnableMemoryMin && opts.MemoryMinValue != "" && isCommandAvailable("systemd-run") {
		wrappedArgs := []string{"systemd-run", "--user", "--scope", fmt.Sprintf("-pMemoryMin=%s", opts.MemoryMinValue), "--"}
		cmdArgs = append(wrappedArgs, cmdArgs...)
	}

	return cmdArgs, env
}

// StopProcessGroup sends SIGINT and SIGTERM to a process group to shut it down gracefully
func StopProcessGroup(proc *os.Process) error {
	if proc == nil {
		return nil
	}
	pgid := -proc.Pid
	_ = syscall.Kill(pgid, syscall.SIGINT)
	time.Sleep(200 * time.Millisecond)
	return syscall.Kill(pgid, syscall.SIGTERM)
}


