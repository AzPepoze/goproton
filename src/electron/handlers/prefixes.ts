import { ipcMain } from "electron";
import fs from "fs";
import path from "path";
import { spawn } from "child_process";
import { getPrefixBaseDir, getPrefixConfigPath, expandPath, scanProtonVersions } from "../launcher/utils.js";

export function registerPrefixHandlers() {
	ipcMain.handle("GetPrefixBaseDir", () => getPrefixBaseDir());

	ipcMain.handle("SavePrefixConfig", async (_, prefixName, opts) => {
		try {
			const configPath = getPrefixConfigPath(prefixName);
			fs.mkdirSync(path.dirname(configPath), { recursive: true });
			fs.writeFileSync(configPath, JSON.stringify(opts, null, 2));
			return null;
		} catch (e: any) {
			return e.toString();
		}
	});

	ipcMain.handle("GetPrefixConfig", async (_, prefixName) => {
		const configPath = getPrefixConfigPath(prefixName);
		if (fs.existsSync(configPath)) {
			try {
				return JSON.parse(fs.readFileSync(configPath, "utf8"));
			} catch {
				return null;
			}
		}
		return null;
	});

	ipcMain.handle("LoadPrefixConfig", async (_, prefixName) => {
		const configPath = getPrefixConfigPath(prefixName);
		if (fs.existsSync(configPath)) {
			try {
				return JSON.parse(fs.readFileSync(configPath, "utf8"));
			} catch {
				return null;
			}
		}
		return null;
	});

	ipcMain.handle("ListPrefixes", () => {
		const base = getPrefixBaseDir();
		if (!fs.existsSync(base)) fs.mkdirSync(base, { recursive: true });
		try {
			const entries = fs.readdirSync(base, { withFileTypes: true });
			const prefixes = entries.filter((e) => e.isDirectory()).map((e) => e.name);
			if (prefixes.length === 0) {
				fs.mkdirSync(path.join(base, "Default"), { recursive: true });
				return ["Default"];
			}
			return prefixes;
		} catch {
			return ["Default"];
		}
	});

	ipcMain.handle("CreatePrefix", (_, name) => {
		try {
			const prefixPath = path.join(getPrefixBaseDir(), name);
			fs.mkdirSync(prefixPath, { recursive: true });
			return null;
		} catch (e: any) {
			return e.toString();
		}
	});

	ipcMain.handle("RunPrefixTool", async (_, prefixPath, toolName, protonPattern) => {
		try {
			// Find the Proton installation matching the pattern
			const protonTools = scanProtonVersions();
			let selectedProton = null;

			// If a specific pattern is provided, find matching proton
			if (protonPattern) {
				selectedProton = protonTools.find(
					(tool) => protonPattern.includes(tool.Name) || tool.Name.includes(protonPattern),
				);
			}

			// Fallback to first available Proton
			if (!selectedProton && protonTools.length > 0) {
				selectedProton = protonTools[0];
			}

			if (!selectedProton) {
				return "No Proton installation found";
			}

			const protonPath = path.join(selectedProton.Path, "proton");
			if (!fs.existsSync(protonPath)) {
				return `Proton binary not found at ${protonPath}`;
			}

			// Set up environment
			const env = { ...process.env };
			env["WINEPREFIX"] = expandPath(prefixPath);

			// Determine the tool to run
			let toolCmd: string;
			switch (toolName.toLowerCase()) {
				case "winecfg":
					toolCmd = "winecfg";
					break;
				case "regedit":
					toolCmd = "regedit";
					break;
				case "cmd":
					toolCmd = "cmd";
					break;
				case "winetricks":
					toolCmd = "winetricks";
					break;
				case "taskmgr":
					toolCmd = "taskmgr";
					break;
				case "explorer":
					toolCmd = "explorer";
					break;
				default:
					return `Unknown tool: ${toolName}`;
			}

			// Spawn the tool through Proton
			const child = spawn(protonPath, ["run", toolCmd], {
				env,
				detached: true,
				stdio: "ignore",
			});
			child.unref();

			return null;
		} catch (e: any) {
			console.error("RunPrefixTool Error:", e);
			return e.toString();
		}
	});
}
