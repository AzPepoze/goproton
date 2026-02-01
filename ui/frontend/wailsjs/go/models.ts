export namespace frontend {
	export class FileFilter {
		DisplayName: string;
		Pattern: string;

		static createFrom(source: any = {}) {
			return new FileFilter(source);
		}

		constructor(source: any = {}) {
			if ("string" === typeof source) source = JSON.parse(source);
			this.DisplayName = source["DisplayName"];
			this.Pattern = source["Pattern"];
		}
	}
}

export namespace launcher {
	export class GitHubAsset {
		name: string;
		browser_download_url: string;
		size: number;
		content_type: string;

		static createFrom(source: any = {}) {
			return new GitHubAsset(source);
		}

		constructor(source: any = {}) {
			if ("string" === typeof source) source = JSON.parse(source);
			this.name = source["name"];
			this.browser_download_url = source["browser_download_url"];
			this.size = source["size"];
			this.content_type = source["content_type"];
		}
	}
	export class GitHubRelease {
		tag_name: string;
		name: string;
		published_at: string;
		html_url: string;
		assets: [];

		static createFrom(source: any = {}) {
			return new GitHubRelease(source);
		}

		constructor(source: any = {}) {
			if ("string" === typeof source) source = JSON.parse(source);
			this.tag_name = source["tag_name"];
			this.name = source["name"];
			this.published_at = source["published_at"];
			this.html_url = source["html_url"];
			this.assets = this.convertValues(source["assets"], GitHubAsset);
		}

		convertValues(a: any, classs: any, asMap: boolean = false): any {
			if (!a) {
				return a;
			}
			if (a.slice && a.map) {
				return (a as any[]).map((elem) => this.convertValues(elem, classs));
			} else if ("object" === typeof a) {
				if (asMap) {
					for (const key of Object.keys(a)) {
						a[key] = new classs(a[key]);
					}
					return a;
				}
				return new classs(a);
			}
			return a;
		}
	}
	export class LaunchOptions {
		GamePath: string;
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
		EnableMemoryMin: boolean;
		MemoryMinValue: string;

		static createFrom(source: any = {}) {
			return new LaunchOptions(source);
		}

		constructor(source: any = {}) {
			if ("string" === typeof source) source = JSON.parse(source);
			this.GamePath = source["GamePath"];
			this.PrefixPath = source["PrefixPath"];
			this.ProtonPattern = source["ProtonPattern"];
			this.ProtonPath = source["ProtonPath"];
			this.CustomArgs = source["CustomArgs"];
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
			this.EnableMemoryMin = source["EnableMemoryMin"];
			this.MemoryMinValue = source["MemoryMinValue"];
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
			if ("string" === typeof source) source = JSON.parse(source);
			this.Name = source["Name"];
			this.Path = source["Path"];
			this.IsSteam = source["IsSteam"];
			this.DisplayName = source["DisplayName"];
		}
	}
	export class ProtonVariant {
		ID: string;
		Name: string;
		Description: string;
		RepoOwner: string;
		RepoName: string;

		static createFrom(source: any = {}) {
			return new ProtonVariant(source);
		}

		constructor(source: any = {}) {
			if ("string" === typeof source) source = JSON.parse(source);
			this.ID = source["ID"];
			this.Name = source["Name"];
			this.Description = source["Description"];
			this.RepoOwner = source["RepoOwner"];
			this.RepoName = source["RepoName"];
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
			if ("string" === typeof source) source = JSON.parse(source);
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
			if ("string" === typeof source) source = JSON.parse(source);
			this.isLsfgInstalled = source["isLsfgInstalled"];
			this.lsfgVersion = source["lsfgVersion"];
		}
	}
}
