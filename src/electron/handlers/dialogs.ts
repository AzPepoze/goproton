import { ipcMain, dialog, BrowserWindow } from "electron";

export function registerDialogHandlers(win: BrowserWindow | null) {
	ipcMain.handle("PickFile", async () => {
		if (!win) return "";
		const result = await dialog.showOpenDialog(win, {
			title: "Select Game Executable",
			filters: [
				{ name: "Executables", extensions: ["exe"] },
				{ name: "All Files", extensions: ["*"] },
			],
			properties: ["openFile"],
		});
		return result.canceled ? "" : result.filePaths[0];
	});

	ipcMain.handle("PickFolder", async () => {
		if (!win) return "";
		const result = await dialog.showOpenDialog(win, {
			title: "Select Prefix Directory",
			properties: ["openDirectory"],
		});
		return result.canceled ? "" : result.filePaths[0];
	});

	ipcMain.handle("PickFileCustom", async (_, title) => {
		if (!win) return "";
		const result = await dialog.showOpenDialog(win, {
			title: title || "Select File",
			properties: ["openFile"],
		});
		return result.canceled ? "" : result.filePaths[0];
	});
}
