package launcher

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// LaunchOptions defines the configuration for launching a game
type LaunchOptions struct {
	GamePath       string
	PrefixPath     string
	ProtonPattern  string
	ProtonPath     string // Full path to the specific proton folder
	EnableGamescope bool
	GamescopeW     string
	GamescopeH     string
	GamescopeR     string
	EnableMangoHud bool
	EnableGamemode bool
}

// BuildCommand generates the command slice and environment variables
func BuildCommand(opts LaunchOptions) ([]string, []string) {
	var cmdArgs []string
	
	// Environment Variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("WINEPREFIX=%s", ExpandPath(opts.PrefixPath)))
	env = append(env, fmt.Sprintf("UMU_PROTON_PATTERN=%s", opts.ProtonPattern))
	
	// Critical: Tell umu-run where to look for this custom proton
	if opts.ProtonPath != "" {
		// PROTONPATH should point to the specific proton folder
		// e.g. /home/user/.steam/root/compatibilitytools.d/dwproton-10.0-14
		env = append(env, fmt.Sprintf("PROTONPATH=%s", opts.ProtonPath))
	}
	
	if opts.EnableMangoHud {
		env = append(env, "MANGOHUD=1")
	}

	// Wrapper: Gamemode
	if opts.EnableGamemode {
		cmdArgs = append(cmdArgs, "gamemoderun")
	}

	// Wrapper: Gamescope
	if opts.EnableGamescope {
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
		cmdArgs = append(cmdArgs, "--")
	}

	// Base Command: umu-run
	cmdArgs = append(cmdArgs, "umu-run", opts.GamePath)

	return cmdArgs, env
}

// RunGameWithLog executes the command and streams stdout/stderr to the onLog callback
// Returns the process object (for killing) or an error
func RunGameWithLog(cmdArgs []string, env []string, onLog func(string), onExit func()) (*os.Process, error) {
	if len(cmdArgs) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = env
	
	// Create a new Process Group so we can kill the whole tree later
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// Stream Output in Goroutines
	readPipe := func(r io.Reader, prefix string) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			onLog(fmt.Sprintf("%s%s", prefix, scanner.Text()))
		}
	}

	go readPipe(stdout, "")
	go readPipe(stderr, "")

	// Wait in background
	go func() {
		cmd.Wait()
		onLog("\n--- Process Exited ---")
		if onExit != nil {
			onExit()
		}
	}()

	return cmd.Process, nil
}

// StopProcessGroup requests the process and all its children to stop gracefully.
// It sends SIGINT followed by SIGTERM to allow proper cleanup and saving.
func StopProcessGroup(proc *os.Process) error {
	if proc == nil {
		return nil
	}
	// Target the Process Group
	pgid := -proc.Pid

	// 1. Send SIGINT (Equivalent to Ctrl+C)
	_ = syscall.Kill(pgid, syscall.SIGINT)
	
	// Give it a short moment to react to SIGINT
	time.Sleep(200 * time.Millisecond)

	// 2. Send SIGTERM (Proper shutdown request)
	return syscall.Kill(pgid, syscall.SIGTERM)
}

// FormatCommandForDisplay converts command slice and env to a readable string
func FormatCommandForDisplay(cmdArgs []string, opts LaunchOptions) string {
	var sb strings.Builder
	
	sb.WriteString("WINEPREFIX=" + opts.PrefixPath + " ")
	sb.WriteString("UMU_PROTON_PATTERN=" + opts.ProtonPattern + " ")
	if opts.EnableMangoHud {
		sb.WriteString("MANGOHUD=1 ")
	}
	
	sb.WriteString(strings.Join(cmdArgs, " "))
	
	return sb.String()
}