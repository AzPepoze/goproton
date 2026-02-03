package main

import "goproton/pkg/launcher"

// buildLaunchOptions creates the launch options from command line flags
func buildLaunchOptions() launcher.LaunchOptions {
	return launcher.LaunchOptions{
		MainExecutablePath: gamePath,
		LauncherPath:       launcherPath,
		PrefixPath:         prefixPath,
		ProtonPattern:      protonPattern,
		ProtonPath:         protonPath,
		EnableMangoHud:     mango,
		EnableGamemode:     gamemode,
		EnableGamescope:    gamescope,
		GamescopeW:         gsW,
		GamescopeH:         gsH,
		GamescopeR:         gsR,
		EnableLsfgVk:       lsfg,
		LsfgMultiplier:     lsfgMult,
		LsfgPerfMode:       lsfgPerf,
		LsfgDllPath:        lsfgDllPath,
		EnableMemoryMin:    memoryMin,
		MemoryMinValue:     memoryMinValue,
	}
}
