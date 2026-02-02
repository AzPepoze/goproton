export namespace launcher {
    export class LaunchOptions {
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
        LsfgAllowFp16: boolean;
        EnableMemoryMin: boolean;
        MemoryMinValue: string;

        static createFrom(source: any = {}) {
            return new LaunchOptions(source);
        }

        constructor(source: any = {}) {
            if ('string' === typeof source) source = JSON.parse(source);
            Object.assign(this, source);
        }
    }

    export class ProtonTool {
        Name: string;
        Path: string;
        IsSteam: boolean;
        DisplayName: string;
        
        constructor(source: any = {}) { Object.assign(this, source); }
    }
    
    export class UtilsStatus {
        isLsfgInstalled: boolean;
        lsfgVersion: string;
        constructor(source: any = {}) { Object.assign(this, source); }
    }
    
    export class SystemToolsStatus {
        hasGamescope: boolean;
        hasMangoHud: boolean;
        hasGameMode: boolean;
        constructor(source: any = {}) { Object.assign(this, source); }
    }
    
    export class ProtonVariant {
        ID: string;
        Name: string;
        Description: string;
        RepoOwner: string;
        RepoName: string;
        constructor(source: any = {}) { Object.assign(this, source); }
    }
    
    export class GitHubRelease {
        tag_name: string;
        name: string;
        published_at: string;
        html_url: string;
        assets: any[];
        constructor(source: any = {}) { Object.assign(this, source); }
    }
    
    export class LsfgProfile {
        game_name: string;
        game_path: string;
        launcher_path: string;
        multiplier: string;
        performance_mode: boolean;
        dll_path: string;
        gpu: string;
        flow_scale: string;
        pacing: string;
        allow_fp16: boolean;
        constructor(source: any = {}) { Object.assign(this, source); }
    }
}
