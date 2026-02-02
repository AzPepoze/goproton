import fs from "fs";
import path from "path";
import crypto from "crypto";
import os from "os";
import { getBaseDir } from "./utils.js";
import type { LsfgConfigFile, LsfgConfigProfile } from "./types.js";

export function getLsfgProfilePath(gamePath: string): string {
	const hash = crypto.createHash("sha1").update(gamePath).digest("hex").slice(0, 8);
	const exeName = path.basename(gamePath);
	const baseName = exeName.replace(/\.exe$/i, "");
	const filename = `${baseName}-${hash}.toml`;
	const configDir = path.join(getBaseDir(), "config", "lsfg");
	return path.join(configDir, filename);
}

export function getLsfgConfigPath(): string {
	return path.join(os.homedir(), ".config", "lsfg-vk", "conf.toml");
}

// Simple TOML Parser for the specific LSFG structure
function parseLsfgToml(data: string): LsfgConfigFile {
	const lines = data.split("\n");
	const config: any = { version: 2, global: {}, profile: [] };
	let currentSection = "";
	let currentProfile: any = null;

	for (let line of lines) {
		line = line.trim();
		if (!line || line.startsWith("#")) continue;

		if (line.startsWith("[") && line.endsWith("]")) {
			currentSection = line.slice(1, -1);
			if (currentSection === "profile") {
				if (currentProfile) config.profile.push(currentProfile);
				currentProfile = {};
			}
			continue;
		}

		const [key, ...valParts] = line.split("=");
		if (!key) continue;
		const val = valParts.join("=").trim();
		const k = key.trim();

		let parsedVal: any = val;
		if (val === "true") parsedVal = true;
		else if (val === "false") parsedVal = false;
		else if (val.startsWith('"') && val.endsWith('"')) parsedVal = val.slice(1, -1);
		else if (!isNaN(Number(val))) parsedVal = Number(val);
		else if (val.startsWith("[") && val.endsWith("]")) {
			// Simple array parser for active_in
			parsedVal = val
				.slice(1, -1)
				.split(",")
				.map((s) => s.trim().replace(/^"|"$/g, ""))
				.filter((s) => s);
		}

		if (currentSection === "global") {
			config.global[k] = parsedVal;
		} else if (currentSection === "profile") {
			currentProfile[k] = parsedVal;
		} else {
			config[k] = parsedVal;
		}
	}
	if (currentProfile) config.profile.push(currentProfile);

	return config as LsfgConfigFile;
}

function stringifyLsfgToml(config: LsfgConfigFile): string {
	let output = `version = ${config.version}\n\n`;
	output += `[global]\n`;
	output += `version = ${config.global.version}\n`;
	output += `allow_fp16 = ${config.global.allow_fp16}\n`;
	output += `dll = "${config.global.dll}"\n\n`;

	for (const p of config.profile) {
		output += `[[profile]]\n`;
		output += `name = "${p.name}"\n`;
		if (Array.isArray(p.active_in)) {
			output += `active_in = [${p.active_in.map((s) => `"${s}"`).join(", ")}]\n`;
		} else {
			output += `active_in = "${p.active_in}"\n`;
		}
		output += `multiplier = ${p.multiplier}\n`;
		output += `performance_mode = ${p.performance_mode}\n`;
		output += `gpu = "${p.gpu}"\n`;
		output += `flow_scale = ${p.flow_scale.toFixed(1)}\n`;
		output += `pacing = "${p.pacing}"\n\n`;
	}
	return output;
}

export function findLsfgProfileForGame(gamePath: string): { profile: LsfgConfigProfile | null; index: number } {
	const configPath = getLsfgConfigPath();
	if (!fs.existsSync(configPath)) return { profile: null, index: -1 };

	try {
		const data = fs.readFileSync(configPath, "utf8");
		const config = parseLsfgToml(data);
		const exeName = path.basename(gamePath).toLowerCase();

		for (let i = 0; i < config.profile.length; i++) {
			const p = config.profile[i];
			const activeIn = Array.isArray(p.active_in) ? p.active_in : [p.active_in];
			if (activeIn.some((s) => s.toLowerCase() === exeName)) {
				return { profile: p, index: i };
			}
		}
	} catch (e) {
		console.error("Failed to parse LSFG config:", e);
	}
	return { profile: null, index: -1 };
}

export function saveLsfgProfile(
	gamePath: string,
	opts: Partial<LsfgConfigProfile & { allowFp16: boolean; dllPath: string }>,
): void {
	const configPath = getLsfgConfigPath();
	const configDir = path.dirname(configPath);
	if (!fs.existsSync(configDir)) fs.mkdirSync(configDir, { recursive: true });

	let config: LsfgConfigFile = {
		version: 2,
		global: { version: 2, allow_fp16: opts.allowFp16 ?? true, dll: opts.dllPath ?? "" },
		profile: [],
	};

	if (fs.existsSync(configPath)) {
		try {
			config = parseLsfgToml(fs.readFileSync(configPath, "utf8"));
		} catch (e) {
			console.error("Failed to read existing LSFG config, creating new one");
		}
	}

	// Update global
	if (opts.allowFp16 !== undefined) config.global.allow_fp16 = opts.allowFp16;
	if (opts.dllPath !== undefined) config.global.dll = opts.dllPath;

	const exeName = path.basename(gamePath);
	const { index } = findLsfgProfileForGame(gamePath);

	const newProfile: LsfgConfigProfile = {
		name: exeName.replace(/\.[^/.]+$/, ""),
		active_in: exeName,
		multiplier: opts.multiplier ?? 2,
		performance_mode: opts.performance_mode ?? false,
		gpu: opts.gpu ?? "auto",
		flow_scale: opts.flow_scale ?? 1.0,
		pacing: "none",
	};

	if (index !== -1) {
		config.profile[index] = { ...config.profile[index], ...newProfile };
	} else {
		config.profile.push(newProfile);
	}

	fs.writeFileSync(configPath, stringifyLsfgToml(config));
}
