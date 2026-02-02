import { ipcMain, app } from "electron";

export function registerEnvHandlers() {
	ipcMain.handle("GetInitialLauncherPath", () => process.env.GOPROTON_LAUNCHER_PATH || "");
	ipcMain.handle("GetInitialGamePath", () => process.env.GOPROTON_GAME_PATH || "");
	ipcMain.handle("GetShouldEditLsfg", () => process.env.GOPROTON_EDIT_LSFG === "1");
	ipcMain.handle("CloseWindow", () => app.quit());
}
