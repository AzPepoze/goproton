import { writable } from "svelte/store";
import type { launcher } from "../models";

export interface RunState {
	gamePath: string;
	gameIcon: string;
	launcherIcon: string;
	prefixPath: string;
	selectedPrefixName: string;
	selectedProton: string;
	options: launcher.LaunchOptions;
}

const defaultOptions: launcher.LaunchOptions = {
	GamePath: "",
	LauncherPath: "",
	UseGameExe: false,
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
	LsfgGpu: "",
	LsfgFlowScale: "",
	LsfgPacing: "",
	LsfgAllowFp16: false,
	EnableMemoryMin: false,
	MemoryMinValue: "4G",
};

const initial: RunState = {
	gamePath: "",
	gameIcon: "",
	launcherIcon: "",
	prefixPath: "",
	selectedPrefixName: "Default",
	selectedProton: "",
	options: defaultOptions,
};

export const runState = writable<RunState>(initial);
