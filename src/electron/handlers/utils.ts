import { ipcMain, BrowserWindow } from "electron";
import { exec } from "child_process";
import { promisify } from "util";
import { getSystemToolsStatus, getUtilsStatus } from "../launcher/utils.js";
import { installLsfg, installProton, KNOWN_PROTON_VARIANTS } from "../launcher/downloader.js";

const execAsync = promisify(exec);

export function registerUtilsHandlers(win: BrowserWindow | null) {
	ipcMain.handle("CleanupProcesses", async () => {
		try {
			await Promise.all(
				["umu-run", "reaper", "gamescopereaper"].map((p) => execAsync(`pkill -f ${p}`).catch(() => {})),
			);
			return null;
		} catch (e: any) {
			return e.toString();
		}
	});

	ipcMain.handle("InstallLsfg", async () => {
		try {
			await installLsfg((percent, message) => {
				win?.webContents.send("lsfg-install-progress", { percent, message });
			});
			return null;
		} catch (e: any) {
			throw new Error(e.toString());
		}
	});

	ipcMain.handle("InstallProton", async (_, url) => {
		try {
			await installProton(url, (percent, message) => {
				win?.webContents.send("lsfg-install-progress", { percent, message });
			});
			return null;
		} catch (e: any) {
			throw new Error(e.toString());
		}
	});

	ipcMain.handle("GetUtilsStatus", () => getUtilsStatus());
	ipcMain.handle("GetSystemToolsStatus", () => getSystemToolsStatus());
	ipcMain.handle("GetProtonVariants", () => KNOWN_PROTON_VARIANTS);

	ipcMain.handle("UninstallLsfg", async () => {
		try {
			// Remove all LSFG-VK files that were installed
			const removalCommands = [
				"rm -f /usr/lib/liblsfg-vk-layer.so",
				"rm -f /usr/share/vulkan/implicit_layer.d/VkLayer_LSFGVK_frame_generation.json",
				"rm -f /usr/share/icons/hicolor/256x256/apps/gay.pancake.lsfg-vk-ui.png",
				"rm -f /usr/share/applications/gay.pancake.lsfg-vk-ui.desktop",
				"rm -f /usr/bin/lsfg-vk-cli",
				"rm -f /usr/bin/lsfg-vk-ui",
			].join(" && ");

			await execAsync(`pkexec sh -c "${removalCommands}"`);
			return null;
		} catch (e: any) {
			return e.toString();
		}
	});
}
