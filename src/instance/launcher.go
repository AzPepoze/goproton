package main

import (
	"log"
	"strings"
)

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
