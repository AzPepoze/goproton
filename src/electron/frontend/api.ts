// Electron IPC Bridge

declare global {
	interface Window {
		ipc: {
			invoke: (channel: string, ...args: any[]) => Promise<any>;
			on: (channel: string, listener: (event: any, ...args: any[]) => void) => () => void;
			off: (channel: string, listener: (...args: any[]) => void) => void;
		};
	}
}

const call = (name: string, ...args: any[]): Promise<any> => {
	if (window.ipc) {
		return window.ipc.invoke(name, ...args);
	} else {
		console.warn(`IPC not available. Call to ${name} failed.`);
		return Promise.resolve(null);
	}
};

// --- App Methods ---

export const GetInitialLauncherPath = () => call("GetInitialLauncherPath");
export const GetInitialGamePath = () => call("GetInitialGamePath");
export const GetShouldEditLsfg = () => call("GetShouldEditLsfg");
export const CloseWindow = () => call("CloseWindow");
export const GetExeIcon = (path: string) => call("GetExeIcon", path);
export const ScanProtonVersions = () => call("ScanProtonVersions");

export const RunGame = (opts: any, showLogs: boolean) => call("RunGame", opts, showLogs);

export const RunPrefixTool = (prefixPath: string, toolName: string, protonPattern: string) =>
	call("RunPrefixTool", prefixPath, toolName, protonPattern);
export const PickFile = () => call("PickFile");
export const PickFolder = () => call("PickFolder");
export const GetConfig = (exePath: string) => call("GetConfig", exePath);
export const SavePrefixConfig = (prefixName: string, opts: any) => call("SavePrefixConfig", prefixName, opts);
export const LoadPrefixConfig = (prefixName: string) => call("LoadPrefixConfig", prefixName);
export const ListPrefixes = () => call("ListPrefixes");
export const CreatePrefix = (name: string) => call("CreatePrefix", name);
export const GetPrefixBaseDir = () => call("GetPrefixBaseDir");
export const GetUtilsStatus = () => call("GetUtilsStatus");
export const GetSystemToolsStatus = () => call("GetSystemToolsStatus");
export const InstallLsfg = () => call("InstallLsfg");
export const DetectLosslessDll = () => call("DetectLosslessDll");
export const GetListGpus = () => call("GetListGpus");
export const PickFileCustom = (title: string, filters: any) => call("PickFileCustom", title, filters);
export const UninstallLsfg = () => call("UninstallLsfg");
export const CleanupProcesses = () => call("CleanupProcesses");
export const GetTotalRam = () => call("GetTotalRam");
export const GetProtonVariants = () => call("GetProtonVariants");
export const GetProtonReleases = (variantID: string) => call("GetProtonReleases", variantID);
export const InstallProtonVersion = (url: string, version: string) => call("InstallProtonVersion", url, version);
export const GetLsfgProfileForGame = (gamePath: string) => call("GetLsfgProfileForGame", gamePath);

export const SaveLsfgProfile = (
	gamePath: string,
	multiplier: number,
	perfMode: boolean,
	dllPath: string,
	gpu: string,
	flowScale: string,
	pacing: string,
	allowFp16: boolean,
) => call("SaveLsfgProfile", gamePath, multiplier, perfMode, dllPath, gpu, flowScale, pacing, allowFp16);

export const RemoveProfile = (gamePath: string) => call("RemoveProfile", gamePath);
export const EditLsfgConfigForGame = (gamePath: string) => call("EditLsfgConfigForGame", gamePath);

// Runtime Mocks (Events)
export const EventsOn = (eventName: string, callback: (...args: any[]) => void) => {
	if (window.ipc) {
		window.ipc.on(eventName, (event, ...args) => callback(...args));
	}
};

export const EventsOff = (eventName: string) => {
	// Implementation depends on storing the listener ref, skipping for now
};

export const WindowHide = () => {
	call("CloseWindow");
};

export const BrowserOpenURL = (url: string) => {
	window.open(url, "_blank");
};
