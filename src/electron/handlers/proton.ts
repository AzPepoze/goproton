import { ipcMain } from "electron";
import { scanProtonVersions } from "../launcher/utils.js";

export function registerProtonHandlers() {
	ipcMain.handle("ScanProtonVersions", () => scanProtonVersions());
}
