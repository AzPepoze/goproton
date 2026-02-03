import type { launcher } from "../../wailsjs/go/models";

export interface ConfigLoadResult {
	success: boolean;
	data: Partial<launcher.LaunchOptions> | null;
	error: string;
}

export interface ConfigSaveResult {
	success: boolean;
	error: string;
}

/**
 * Handles config loading with error recovery
 */
export async function loadConfigSafely<T>(
	loaderFn: () => Promise<T | null>,
	fallback: T,
): Promise<{ data: T; error: string }> {
	try {
		const data = await loaderFn();
		return {
			data: data ?? fallback,
			error: "",
		};
	} catch (err) {
		console.error("Config load error:", err);
		return {
			data: fallback,
			error: String(err),
		};
	}
}

/**
 * Handles config saving with error recovery
 */
export async function saveConfigSafely<T>(
	data: T,
	saverFn: (data: T) => Promise<void>,
): Promise<{ success: boolean; error: string }> {
	try {
		await saverFn(data);
		return { success: true, error: "" };
	} catch (err) {
		console.error("Config save error:", err);
		return {
			success: false,
			error: String(err),
		};
	}
}

/**
 * Validates if required paths are set in LaunchOptions
 */
export function validateLaunchOptions(options: launcher.LaunchOptions): { valid: boolean; errors: string[] } {
	const errors: string[] = [];

	if (!options.MainExecutablePath && !options.LauncherPath) {
		errors.push("Game path or launcher path is required");
	}

	if (options.EnableLsfgVk) {
		if (!options.LsfgDllPath) {
			errors.push("LSFG-VK enabled but DLL path not set");
		}
		if (!options.LsfgGpu) {
			errors.push("LSFG-VK enabled but GPU not selected");
		}
	}

	return {
		valid: errors.length === 0,
		errors,
	};
}
