import { writable } from "svelte/store";
import type { launcher } from "../../wailsjs/go/models";

export interface RunState {
	gamePath: string;
	gameIcon: string;
	prefixPath: string;
	selectedPrefixName: string;
	selectedProton: string;
	options: launcher.LaunchOptions;
}

const defaultOptions: launcher.LaunchOptions = {
	GamePath: "",
	PrefixPath: "",
	ProtonPattern: "",
	ProtonPath: "",
	CustomArgs: "",
	EnableGamescope: false,
	GamescopeW: "1920",
	GamescopeH: "1080",
	GamescopeR: "60",
	EnableMangoHud: false,
	EnableGamemode: false,
	EnableLsfgVk: false,
	LsfgMultiplier: "2",
	LsfgPerfMode: false,
	LsfgDllPath: "",
	EnableMemoryMin: false,
	MemoryMinValue: "4G",
};

const initial: RunState = {
	gamePath: "",
	gameIcon: "",
	prefixPath: "",
	selectedPrefixName: "Default",
	selectedProton: "",
	options: defaultOptions,
};

export const runState = writable<RunState>(initial);
