export namespace frontend {
	
	export class FileFilter {
	    DisplayName: string;
	    Pattern: string;
	
	    static createFrom(source: any = {}) {
	        return new FileFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DisplayName = source["DisplayName"];
	        this.Pattern = source["Pattern"];
	    }
	}

}

export namespace launcher {
	
	export class LaunchOptions {
	    GamePath: string;
	    PrefixPath: string;
	    ProtonPattern: string;
	    ProtonPath: string;
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
	
	    static createFrom(source: any = {}) {
	        return new LaunchOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.GamePath = source["GamePath"];
	        this.PrefixPath = source["PrefixPath"];
	        this.ProtonPattern = source["ProtonPattern"];
	        this.ProtonPath = source["ProtonPath"];
	        this.EnableGamescope = source["EnableGamescope"];
	        this.GamescopeW = source["GamescopeW"];
	        this.GamescopeH = source["GamescopeH"];
	        this.GamescopeR = source["GamescopeR"];
	        this.EnableMangoHud = source["EnableMangoHud"];
	        this.EnableGamemode = source["EnableGamemode"];
	        this.EnableLsfgVk = source["EnableLsfgVk"];
	        this.LsfgMultiplier = source["LsfgMultiplier"];
	        this.LsfgPerfMode = source["LsfgPerfMode"];
	        this.LsfgDllPath = source["LsfgDllPath"];
	    }
	}
	export class ProtonTool {
	    Name: string;
	    Path: string;
	    IsSteam: boolean;
	    DisplayName: string;
	
	    static createFrom(source: any = {}) {
	        return new ProtonTool(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Path = source["Path"];
	        this.IsSteam = source["IsSteam"];
	        this.DisplayName = source["DisplayName"];
	    }
	}
	export class SystemToolsStatus {
	    hasGamescope: boolean;
	    hasMangoHud: boolean;
	    hasGameMode: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SystemToolsStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hasGamescope = source["hasGamescope"];
	        this.hasMangoHud = source["hasMangoHud"];
	        this.hasGameMode = source["hasGameMode"];
	    }
	}
	export class UtilsStatus {
	    isLsfgInstalled: boolean;
	    lsfgVersion: string;
	
	    static createFrom(source: any = {}) {
	        return new UtilsStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isLsfgInstalled = source["isLsfgInstalled"];
	        this.lsfgVersion = source["lsfgVersion"];
	    }
	}

}

