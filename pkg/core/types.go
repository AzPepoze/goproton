package core

type LaunchOptions struct {
	MainExecutablePath string `json:"MainExecutablePath"`
	LauncherPath       string `json:"LauncherPath"`
	HaveGameExe        bool   `json:"HaveGameExe"`
	PrefixPath         string `json:"PrefixPath"`
	ProtonPattern      string `json:"ProtonPattern"`
	ProtonPath         string `json:"ProtonPath"`
	CustomArgs         string `json:"CustomArgs"`
	EnableGamescope    bool   `json:"EnableGamescope"`
	GamescopeW         string `json:"GamescopeW"`
	GamescopeH         string `json:"GamescopeH"`
	GamescopeR         string `json:"GamescopeR"`
	EnableMangoHud     bool   `json:"EnableMangoHud"`
	EnableGamemode     bool   `json:"EnableGamemode"`
	EnableLsfgVk       bool   `json:"EnableLsfgVk"`
	LsfgMultiplier     string `json:"LsfgMultiplier"`
	LsfgPerfMode       bool   `json:"LsfgPerfMode"`
	LsfgDllPath        string `json:"LsfgDllPath"`
	LsfgGpu            string `json:"LsfgGpu"`
	LsfgFlowScale      string `json:"LsfgFlowScale"`
	LsfgPacing         string `json:"LsfgPacing"`
	LsfgAllowFp16      bool   `json:"LsfgAllowFp16"`
	EnableMemoryMin    bool   `json:"EnableMemoryMin"`
	MemoryMinValue     string `json:"MemoryMinValue"`
}

type LsfgProfile struct {
	GameName           string `toml:"game_name"`
	MainExecutablePath string `toml:"main_executable_path"`
	LauncherPath       string `toml:"launcher_path"`
	Multiplier         string `toml:"multiplier"`
	PerformanceMode    bool   `toml:"performance_mode"`
	DllPath            string `toml:"dll_path"`
	Gpu                string `toml:"gpu"`
	FlowScale          string `toml:"flow_scale"`
	Pacing             string `toml:"pacing"`
	AllowFp16          bool   `toml:"allow_fp16"`
}
