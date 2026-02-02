import { ipcMain, app, BrowserWindow } from "electron";
import fs from "fs";
import path from "path";
import { spawn } from "child_process";

export function registerLauncherHandlers(win: BrowserWindow | null, __dirname: string) {
	ipcMain.handle("RunGame", async (_, optsJson, showLogs) => {
		try {
			const opts = typeof optsJson === "string" ? JSON.parse(optsJson) : optsJson;

			// Find instance binary
			let instanceBin = "";
			if (app.isPackaged) {
				instanceBin = path.join(process.resourcesPath, "goproton-instance");
			} else {
				instanceBin = path.resolve(__dirname, "../../bin/goproton-instance");
			}

			if (!fs.existsSync(instanceBin)) {
				console.error(`Instance binary not found at: ${instanceBin}`);
				instanceBin = "goproton-instance";
			}

			const args = [
				`--game=${opts.GamePath}`,
				`--launcher=${opts.LauncherPath || ""}`,
				`--prefix=${opts.PrefixPath}`,
				`--proton-pattern=${opts.ProtonPattern || ""}`,
				`--proton-path=${opts.ProtonPath || ""}`,
			];

			if (opts.EnableMangoHud) args.push("--mango");
			if (opts.EnableGamemode) args.push("--gamemode");
			if (opts.EnableGamescope) {
				args.push("--gamescope");
				if (opts.GamescopeW) args.push(`--gs-w=${opts.GamescopeW}`);
				if (opts.GamescopeH) args.push(`--gs-h=${opts.GamescopeH}`);
				if (opts.GamescopeR) args.push(`--gs-r=${opts.GamescopeR}`);
			}
			if (opts.EnableLsfgVk) {
				args.push("--lsfg");
				args.push(`--lsfg-mult=${opts.LsfgMultiplier || "2"}`);
				if (opts.LsfgPerfMode) args.push("--lsfg-perf");
				if (opts.LsfgDllPath) args.push(`--lsfg-dll-path=${opts.LsfgDllPath}`);
			}

			if (opts.EnableMemoryMin) {
				args.push("--memory-min");
				if (opts.MemoryMinValue) args.push(`--memory-min-value=${opts.MemoryMinValue}`);
			}

			// Pass logs flag explicitly
			if (showLogs === true) {
				args.push("--logs=true");
			} else if (showLogs === false) {
				args.push("--logs=false");
			}

			console.log("Spawning instance manager:", instanceBin, args.join(" "));

			const child = spawn(instanceBin, args, {
				detached: true,
				stdio: "ignore",
			});
			child.unref();

			app.quit();
			return null;
		} catch (e: any) {
			console.error("RunGame Error:", e);
			return e.toString();
		}
	});
}
