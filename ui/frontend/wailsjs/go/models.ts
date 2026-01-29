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

}

