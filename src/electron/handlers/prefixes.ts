import { ipcMain } from "electron";
import fs from "fs";
import path from "path";
import { getPrefixBaseDir, getPrefixConfigPath } from "../launcher/utils.js";

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
}
