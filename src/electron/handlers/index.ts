import { BrowserWindow } from "electron";
import { registerEnvHandlers } from "./env.js";
import { registerDialogHandlers } from "./dialogs.js";
import { registerIconHandlers } from "./icons.js";
import { registerSystemHandlers } from "./system.js";
import { registerConfigHandlers } from "./config.js";
import { registerPrefixHandlers } from "./prefixes.js";
import { registerLauncherHandlers } from "./launcher.js";
import { registerProtonHandlers } from "./proton.js";
import { registerUtilsHandlers } from "./utils.js";
import { registerLsfgHandlers } from "./lsfg.js";

export function registerAllHandlers(win: BrowserWindow | null, __dirname: string) {
	// Register all handler groups
	registerEnvHandlers();
	registerDialogHandlers(win);
	registerIconHandlers();
	registerSystemHandlers();
	registerConfigHandlers();
	registerPrefixHandlers();
	registerLauncherHandlers(win, __dirname);
	registerProtonHandlers();
	registerUtilsHandlers(win);
	registerLsfgHandlers();
}
