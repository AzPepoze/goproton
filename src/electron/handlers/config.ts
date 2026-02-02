import { ipcMain } from "electron";
import fs from "fs";
import path from "path";
import { getGameConfigFile } from "../launcher/utils.js";

export function registerConfigHandlers() {
	ipcMain.handle("GetConfig", async (_, exePath) => {
		if (!exePath) return null;
		const configPath = getGameConfigFile(exePath);
		// Check legacy path too (next to exe)
		const legacyPath = exePath + ".goproton.json";
		const actualPath = fs.existsSync(configPath) ? configPath : fs.existsSync(legacyPath) ? legacyPath : null;

		if (actualPath) {
			try {
				const data = fs.readFileSync(actualPath, "utf8");
				return JSON.parse(data);
			} catch (e) {
				console.error(e);
				return null;
			}
		}
		return null;
	});

	ipcMain.handle("SaveConfig", async (_, opts) => {
		try {
			const exePath = opts.LauncherPath || opts.GamePath;
			if (!exePath) return "Missing Game/Launcher Path";
			const configPath = getGameConfigFile(exePath);
			fs.mkdirSync(path.dirname(configPath), { recursive: true });
			fs.writeFileSync(configPath, JSON.stringify(opts, null, 2));
			return null;
		} catch (e: any) {
			return e.toString();
		}
	});
}
