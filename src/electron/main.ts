import { app, BrowserWindow, shell } from "electron";
import path from "path";
import fs from "fs";
import { fileURLToPath } from "url";
import { spawn } from "child_process";
import { registerAllHandlers } from "./handlers/index.js";

// Support __dirname in ESM
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

app.commandLine.appendSwitch("--js-flags", "--max-old-space-size=512");
app.commandLine.appendSwitch("--no-zygote");
app.commandLine.appendSwitch("--no-sandbox");

// Set paths based on packaging state
const isDev = !app.isPackaged;
const preloadPath = path.join(__dirname, "preload.js");
const indexPath = isDev ? "http://localhost:5173" : path.join(__dirname, "../frontend/index.html");

let win: BrowserWindow | null;

const createWindow = () => {
	win = new BrowserWindow({
		width: 1024,
		height: 768,
		backgroundColor: "#18181b",
		webPreferences: {
			preload: preloadPath,
			contextIsolation: true,
			nodeIntegration: false,
			backgroundThrottling: true,
			spellcheck: false,
			offscreen: false,
			enableWebSQL: false,
		},
		title: "GoProton",
		icon: path.join(__dirname, isDev ? "../../src/public" : "../frontend", "favicon.ico"),
	});

	// Register all IPC handlers
	registerAllHandlers(win, __dirname);

	win.setAutoHideMenuBar(true);
	win.setMenuBarVisibility(false);

	// Open links in system browser instead of new Electron window
	win.webContents.setWindowOpenHandler(({ url }) => {
		if (url.startsWith("http://") || url.startsWith("https://")) {
			shell.openExternal(url);
		}
		return { action: "deny" };
	});

	if (isDev) {
		win.loadURL(indexPath);
	} else {
		win.loadFile(indexPath);
	}
};

app.on("window-all-closed", () => {
	app.quit();
});

app.whenReady().then(() => {
	// Parse Command Line Args
	const args = process.argv;
	let isDebugMode = false;
	const goprotonArgs: string[] = [];

	if (app.isPackaged) {
		if (args.length > 1) {
			for (let i = 1; i < args.length; i++) {
				if (args[i] === "--debug") {
					isDebugMode = true;
				} else {
					goprotonArgs.push(args[i]);
				}
			}
		}
	} else {
		if (args.length > 2) {
			for (let i = 2; i < args.length; i++) {
				if (args[i] === "--debug") {
					isDebugMode = true;
				} else {
					goprotonArgs.push(args[i]);
				}
			}
		}
	}

	if (process.env.RUN_FROM_GOPROTON === "true") {
		createWindow();
		return;
	}

	const goprotonPath = path.join(app.getAppPath(), "../goproton");

	if (fs.existsSync(goprotonPath)) {
		console.log("Starting goproton process:", goprotonPath);

		const goprotonProcess = spawn(goprotonPath, goprotonArgs, {
			stdio: "inherit",
			detached: false,
			env: {
				...process.env,
			},
		});

		goprotonProcess.on("close", (code) => {
			console.log("goproton process exited with code:", code);
			// Exit unless in debug mode
			if (!isDebugMode) {
				app.quit();
			} else {
				console.log("Debug mode: keeping app alive, press Ctrl+C to exit");
			}
		});

		goprotonProcess.on("error", (err) => {
			console.error("Failed to run goproton:", err);
			if (!isDebugMode) {
				app.quit();
			}
		});
	} else {
		console.error("goproton binary not found at:", goprotonPath);
		app.quit();
	}
});
