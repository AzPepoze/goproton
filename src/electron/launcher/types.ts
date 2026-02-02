export interface LaunchOptions {
    GamePath: string;
    LauncherPath: string;
    UseGameExe: boolean;
    PrefixPath: string;
    ProtonPattern: string;
    ProtonPath: string;
    CustomArgs: string;
    EnableGamescope: boolean;
    GamescopeW: string;
    GamescopeH: string;
    GamescopeR: string;
    EnableMangoHud: boolean;
    EnableGamemode: boolean;
    EnableLsfgVk: boolean;
    LsfgMultiplier: string;
    LsfgPerfMode: boolean;
    LsfgDllPath: string;
    LsfgGpu: string;
    LsfgFlowScale: string;
    LsfgPacing: string;
    AllowFp16: boolean;
    EnableMemoryMin: boolean;
    MemoryMinValue: string;
}

export interface LsfgGlobalConfig {
    version: number;
    allow_fp16: boolean;
    dll: string;
}

export interface LsfgConfigProfile {
    name: string;
    active_in: string | string[];
    multiplier: number;
    performance_mode: boolean;
    gpu: string;
    flow_scale: number;
    pacing: string;
}

export interface LsfgConfigFile {
    version: number;
    global: LsfgGlobalConfig;
    profile: LsfgConfigProfile[];
}

export interface ProtonVariant {
    ID: string;
    Name: string;
    Description: string;
    RepoOwner: string;
    RepoName: string;
}

export interface ProtonTool {
    Name: string;
    Path: string;
    IsSteam: boolean;
    DisplayName: string;
}

export interface UtilsStatus {
    isLsfgInstalled: boolean;
    lsfgVersion: string;
}

export interface SystemToolsStatus {
    hasGamescope: boolean;
    hasMangoHud: boolean;
    hasGameMode: boolean;
}
