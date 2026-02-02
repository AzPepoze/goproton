import { ipcMain } from "electron";
import os from "os";
import path from "path";
import fs from "fs";
import { getListGpus } from "../launcher/utils.js";

export function registerSystemHandlers() {
	ipcMain.handle("GetTotalRam", () => {
		const totalMem = os.totalmem();
		return Math.round(totalMem / (1024 * 1024 * 1024));
	});

	ipcMain.handle("GetListGpus", () => getListGpus());

	ipcMain.handle("DetectLosslessDll", () => {
		const home = os.homedir();
		const paths = [
			path.join(home, ".steam/root/steamapps/common/Lossless Scaling/Lossless.dll"),
			path.join(home, ".local/share/Steam/steamapps/common/Lossless Scaling/Lossless.dll"),
		];
		for (const p of paths) {
			if (fs.existsSync(p)) return p;
		}
		return "";
	});
}
