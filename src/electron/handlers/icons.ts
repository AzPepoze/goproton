import { ipcMain } from "electron";
import fs from "fs";
import { exec } from "child_process";
import path from "path";
import os from "os";

export function registerIconHandlers() {
	ipcMain.handle("GetExeIcon", async (_, exePath) => {
		if (!exePath || !fs.existsSync(exePath)) {
			console.warn("Icon path does not exist:", exePath);
			return "";
		}
		try {
			const iconData = await extractPEIcon(exePath);
			return iconData || "";
		} catch (e) {
			console.error("Failed to extract PE icon:", e);
			return "";
		}
	});
}

function extractPEIcon(exePath: string): Promise<string | null> {
	return new Promise((resolve) => {
		const tempDir = path.join(os.tmpdir(), `goproton-icon-${Date.now()}`);

		// Create temp directory
		fs.mkdirSync(tempDir, { recursive: true });

		// Try wrestool first
		exec(`wrestool -x --output="${tempDir}" "${exePath}"`, (err) => {
			if (!err) {
				// Check for extracted ICO files
				fs.readdir(tempDir, (readErr, files) => {
					if (!readErr && files.length > 0) {
						const icoFile = files.find((f) => f.endsWith(".ico"));
						if (icoFile) {
							const icoPath = path.join(tempDir, icoFile);
							fs.readFile(icoPath, (readFileErr, data) => {
								fs.rm(tempDir, { recursive: true }, () => {});
								if (!readFileErr && data.length > 0) {
									console.log(`Successfully extracted icon from ${exePath} using wrestool`);
									resolve("data:image/x-icon;base64," + data.toString("base64"));
									return;
								}
								tryIcoExtract();
							});
							return;
						}
					}
					tryIcoExtract();
				});
			} else {
				console.log("wrestool failed, trying icoextract...");
				tryIcoExtract();
			}
		});

		function tryIcoExtract() {
			const icoPath = path.join(tempDir, "icon.ico");
			exec(`icoextract "${exePath}" "${icoPath}"`, (err) => {
				if (!err) {
					fs.readFile(icoPath, (readErr, data) => {
						fs.rm(tempDir, { recursive: true }, () => {});
						if (!readErr && data.length > 0) {
							console.log(`Successfully extracted icon from ${exePath} using icoextract`);
							resolve("data:image/x-icon;base64," + data.toString("base64"));
						} else {
							console.warn(`icoextract succeeded but failed to read file: ${readErr}`);
							resolve(null);
						}
					});
				} else {
					console.warn(`icoextract failed: ${err.message}`);
					fs.rm(tempDir, { recursive: true }, () => {});
					resolve(null);
				}
			});
		}
	});
}
