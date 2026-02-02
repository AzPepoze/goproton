import https from "https";
import fs from "fs";
import path from "path";
import os from "os";
import { exec } from "child_process";
import { promisify } from "util";
import type { ProtonVariant } from "./types.js";

const execAsync = promisify(exec);

export const LSFG_REPO = "PancakeTAS/lsfg-vk";

export const KNOWN_PROTON_VARIANTS: ProtonVariant[] = [
	{
		ID: "ge-proton",
		Name: "GE-Proton (GloriousEggroll)",
		Description: "The most popular custom Proton build. Includes many game fixes and codec patches.",
		RepoOwner: "GloriousEggroll",
		RepoName: "proton-ge-custom",
	},
	{
		ID: "proton-cachyos",
		Name: "Proton-CachyOS",
		Description: "Optimized for performance with CachyOS patches and schedulers.",
		RepoOwner: "CachyOS",
		RepoName: "proton-cachyos",
	},
	{
		ID: "kron4ek",
		Name: "Proton-Kron4ek",
		Description: "Vanilla builds and TKG builds. Often smaller and faster updates.",
		RepoOwner: "Kron4ek",
		RepoName: "Proton-Builds",
	},
	{
		ID: "luxtorpeda",
		Name: "Luxtorpeda (Native Tools)",
		Description: "Runs Windows games using native Linux engines (e.g. GZDoom, ScummVM).",
		RepoOwner: "luxtorpeda-dev",
		RepoName: "luxtorpeda",
	},
];

async function fetchJson(url: string): Promise<any> {
	return new Promise((resolve, reject) => {
		https.get(url, { headers: { "User-Agent": "GoProton-App" } }, (res) => {
			let data = "";
			res.on("data", (chunk) => (data += chunk));
			res.on("end", () => {
				if (res.statusCode !== 200) reject(new Error(`Failed to fetch: ${res.statusCode}`));
				else resolve(JSON.parse(data));
			});
		}).on("error", reject);
	});
}

function downloadFile(url: string, dest: string, onProgress: (current: number, total: number) => void): Promise<void> {
	return new Promise((resolve, reject) => {
		https.get(url, { headers: { "User-Agent": "GoProton-App" } }, (res) => {
			if (res.statusCode === 302 || res.statusCode === 301) {
				// Handle redirect
				downloadFile(res.headers.location!, dest, onProgress).then(resolve).catch(reject);
				return;
			}
			if (res.statusCode !== 200) {
				reject(new Error(`Server returned ${res.statusCode}`));
				return;
			}

			const total = parseInt(res.headers["content-length"] || "0", 10);
			let current = 0;
			const file = fs.createWriteStream(dest);

			res.on("data", (chunk) => {
				current += chunk.length;
				file.write(chunk);
				onProgress(current, total);
			});

			res.on("end", () => {
				file.end();
				resolve();
			});

			res.on("error", (err) => {
				file.close();
				fs.unlinkSync(dest);
				reject(err);
			});
		}).on("error", reject);
	});
}

export async function installLsfg(onProgress: (percent: number, msg: string) => void): Promise<void> {
	onProgress(0, "Fetching release info from GitHub...");
	const releases = await fetchJson(`https://api.github.com/repos/${LSFG_REPO}/releases`);

	let downloadURL = "";
	let assetName = "";

	for (const release of releases) {
		for (const asset of release.assets) {
			const name = asset.name.toLowerCase();
			if (
				(name.includes("x86_64") && name.endsWith(".tar.zst")) ||
				(name.includes("linux") && name.endsWith(".tar.xz"))
			) {
				downloadURL = asset.browser_download_url;
				assetName = asset.name;
				break;
			}
		}
		if (downloadURL) break;
	}

	if (!downloadURL) throw new Error("lsfg-vk suitable linux asset not found");

	onProgress(5, `Downloading ${assetName}...`);
	const ext = assetName.endsWith(".tar.zst") ? ".tar.zst" : ".tar.xz";
	const tmpFile = path.join(os.tmpdir(), `lsfg-vk-dl${ext}`);

	await downloadFile(downloadURL, tmpFile, (current, total) => {
		const percent = Math.round((current / total) * 80);
		onProgress(5 + percent, "Downloading...");
	});

	onProgress(85, "Extracting files...");
	const extractTmp = fs.mkdtempSync(path.join(os.tmpdir(), "lsfg-extract-"));

	try {
		const tarArgs =
			ext === ".tar.zst"
				? ["--use-compress-program=unzstd", "-xf", tmpFile, "-C", extractTmp]
				: ["-xf", tmpFile, "-C", extractTmp];
		try {
			await execAsync(`tar ${tarArgs.join(" ")}`);
		} catch (err) {
			// Fallback to basic tar if unzstd fails or not needed
			await execAsync(`tar -xf ${tmpFile} -C ${extractTmp}`);
		}

		onProgress(88, "Installing to system directories (requires sudo)...");
		// Using pkexec to copy files
		await execAsync(`pkexec sh -c "cp -r ${extractTmp}/* /usr"`);

		onProgress(100, "Installation complete!");
	} finally {
		fs.rmSync(extractTmp, { recursive: true, force: true });
		if (fs.existsSync(tmpFile)) fs.unlinkSync(tmpFile);
	}
}

export async function installProton(url: string, onProgress: (percent: number, msg: string) => void): Promise<void> {
	const home = os.homedir();
	let targetBase = path.join(home, ".steam/root/compatibilitytools.d");
	if (!fs.existsSync(path.join(home, ".steam/root"))) {
		targetBase = path.join(home, ".local/share/Steam/compatibilitytools.d");
	}

	if (!fs.existsSync(targetBase)) {
		fs.mkdirSync(targetBase, { recursive: true });
	}

	onProgress(0, "Downloading...");
	const tmpFile = path.join(os.tmpdir(), `proton-dl-${Date.now()}.tar.xz`);

	try {
		await downloadFile(url, tmpFile, (current, total) => {
			const percent = Math.round((current / total) * 50);
			onProgress(percent, `Downloading... ${Math.round((current / total) * 100)}%`);
		});

		onProgress(50, "Extracting (using system tools)...");
		await execAsync(`tar -xf ${tmpFile} -C "${targetBase}"`);

		onProgress(100, "Installation Complete!");
	} finally {
		if (fs.existsSync(tmpFile)) fs.unlinkSync(tmpFile);
	}
}
