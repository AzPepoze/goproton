package launcher

type LaunchOptions struct {
	GamePath        string
	PrefixPath      string
	ProtonPattern   string
	ProtonPath      string
	CustomArgs      string
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
	EnableMemoryMin bool
	MemoryMinValue  string
}
