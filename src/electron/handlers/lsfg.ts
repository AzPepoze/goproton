import { ipcMain } from "electron";
import { saveLsfgProfile, findLsfgProfileForGame } from "../launcher/lsfg.js";
import fs from "fs";
import os from "os";
import path from "path";

export function registerLsfgHandlers() {
	ipcMain.handle("GetLsfgProfileForGame", (_, gamePath) => {
		const { profile } = findLsfgProfileForGame(gamePath);
		return profile;
	});

	ipcMain.handle("SaveLsfgProfile", (_, gamePath, opts) => {
		try {
			saveLsfgProfile(gamePath, opts);
			return null;
		} catch (e: any) {
			return e.toString();
		}
	});

	ipcMain.handle("RemoveProfile", (_, gamePath) => {
		try {
			const configPath = path.join(os.homedir(), ".config", "lsfg-vk", "conf.toml");
			if (fs.existsSync(configPath)) {
				console.log("RemoveProfile called for:", gamePath);
				// A proper implementation would parse and remove the specific profile
			}
			return null;
		} catch (e: any) {
			return e.toString();
		}
	});
}
