package lsfg_utils

const (
	Repo = "PancakeTAS/lsfg-vk"
)

type Status struct {
	IsInstalled bool   `json:"isInstalled"`
	Version     string `json:"version"`
}

type GlobalConfig struct {
	Version   int    `toml:"version"`
	AllowFP16 bool   `toml:"allow_fp16"`
	DLL       string `toml:"dll"`
}

type ConfigProfile struct {
	Name            string      `toml:"name"`
	ActiveIn        interface{} `toml:"active_in"` // Can be string or []string
	Multiplier      int         `toml:"multiplier"`
	PerformanceMode bool        `toml:"performance_mode"`
	GPU             string      `toml:"gpu"`
	FlowScale       float32     `toml:"flow_scale"`
	Pacing          string      `toml:"pacing"`
}

type ConfigFile struct {
	Version  int             `toml:"version"`
	Global   GlobalConfig    `toml:"global"`
	Profiles []ConfigProfile `toml:"profile"`
}
