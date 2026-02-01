<script lang="ts">
	import { onMount } from "svelte";
	import { WindowHide } from "../../wailsjs/runtime/runtime";
	import {
		PickFile,
		PickFolder,
		ScanProtonVersions,
		RunGame,
		GetConfig,
		ListPrefixes,
		GetPrefixBaseDir,
		GetSystemToolsStatus,
		LoadPrefixConfig,
		GetInitialGamePath,
		GetExeIcon,
	} from "../../wailsjs/go/main/App";
	import type { launcher } from "../../wailsjs/go/models";
	import Dropdown from "../components/Dropdown.svelte";
	import ConfigForm from "../components/ConfigForm.svelte";
	import SlideButton from "../components/SlideButton.svelte";
	import Modal from "../components/Modal.svelte";
	import { notifications } from "../notificationStore";
	import { runState } from "../stores/runState";
	import { get } from "svelte/store";

	// Component State
	let mounted = false;

	// Game Selection
	let gamePath = "";
	let gameIcon = "";
	let prefixPath = "";
	let baseDir = "";
	let selectedPrefixName = "Default";

	// Prefix & Utilities
	let availablePrefixes: string[] = [];

	// Proton
	let protonVersions: launcher.ProtonTool[] = [];
	let protonOptions: string[] = [];
	let selectedProton = "";
	let isLoadingProton = true;

	// UI State
	let showLogsWindow = false;
	let showValidationModal = false;
	let missingToolsList: string[] = [];
	let systemStatus: launcher.SystemToolsStatus = { hasGamescope: false, hasMangoHud: false, hasGameMode: false };
	let iconLoadFailed = false;

	// Config
	let options: launcher.LaunchOptions = {
		GamePath: "",
		PrefixPath: "",
		ProtonPattern: "",
		ProtonPath: "",
		CustomArgs: "",
		EnableGamescope: false,
		GamescopeW: "1920",
		GamescopeH: "1080",
		GamescopeR: "60",
		EnableMangoHud: false,
		EnableGamemode: false,
		EnableLsfgVk: false,
		LsfgMultiplier: "2",
		LsfgPerfMode: false,
		LsfgDllPath: "",
		EnableMemoryMin: false,
		MemoryMinValue: "4G",
	};

	async function loadConfigForGame(path: string) {
		try {
			const config = await GetConfig(path);
			if (config) {
				prefixPath = config.PrefixPath;
				if (prefixPath.startsWith(baseDir)) {
					selectedPrefixName = prefixPath.replace(baseDir + "/", "");
				} else {
					selectedPrefixName = "Custom...";
				}
				applyConfigToOptions(config);
			} else {
				await loadConfigForPrefix(selectedPrefixName);
			}
		} catch (err) {}
	}

	async function loadConfigForPrefix(name: string) {
		if (name === "Custom...") return;
		try {
			const config = await LoadPrefixConfig(name);
			if (config) {
				const savedGamePath = options.GamePath;
				const savedPrefixPath = options.PrefixPath;
				applyConfigToOptions(config);
				options.GamePath = savedGamePath;
				if (savedPrefixPath) options.PrefixPath = savedPrefixPath;
			}
		} catch (err) {}
	}

	function applyConfigToOptions(config: launcher.LaunchOptions) {
		const match = protonVersions.find((p) => p.Path === config.ProtonPath);
		if (match) {
			selectedProton = match.DisplayName;
		} else if (config.ProtonPattern) {
			selectedProton = config.ProtonPattern;
		}

		options.CustomArgs = config.CustomArgs || "";
		options.EnableMangoHud = config.EnableMangoHud;
		options.EnableGamemode = config.EnableGamemode;
		options.EnableLsfgVk = config.EnableLsfgVk;
		options.LsfgMultiplier = config.LsfgMultiplier || "2";
		options.LsfgPerfMode = config.LsfgPerfMode;
		options.LsfgDllPath = config.LsfgDllPath || "";
		options.EnableGamescope = config.EnableGamescope;
		options.GamescopeW = config.GamescopeW || "1920";
		options.GamescopeH = config.GamescopeH || "1080";
		options.GamescopeR = config.GamescopeR || "60";
		options.EnableMemoryMin = config.EnableMemoryMin;
		options.MemoryMinValue = config.MemoryMinValue || "4G";
	}

	onMount(async () => {
		try {
			const s = get(runState);
			if (s) {
				if (s.gamePath) {
					gamePath = s.gamePath;
					options.GamePath = s.gamePath;
				}
				if (s.gameIcon) gameIcon = s.gameIcon;
				if (s.prefixPath) prefixPath = s.prefixPath;
				if (s.selectedPrefixName) selectedPrefixName = s.selectedPrefixName;
				if (s.selectedProton) selectedProton = s.selectedProton;
				if (s.options) options = { ...options, ...s.options };
			}

			const initialPath = await GetInitialGamePath();
			if (initialPath && !gamePath) {
				gamePath = initialPath;
				options.GamePath = initialPath;
				const icon = await GetExeIcon(initialPath);
				if (icon) {
					gameIcon = icon;
				} else {
					iconLoadFailed = true;
					notifications.info(
						"Could not extract icon. Please install 'icoutils' or 'icoextract' to display game icons.",
					);
				}
				await loadConfigForGame(initialPath);
			}

			const [tools, prefixes, base, sysStatus] = await Promise.all([
				ScanProtonVersions(),
				ListPrefixes(),
				GetPrefixBaseDir(),
				GetSystemToolsStatus(),
			]);
			if (tools) {
				protonVersions = tools;
				protonOptions = tools.map((t) => t.DisplayName);
				if (protonOptions.length > 0 && !selectedProton) {
					selectedProton = protonOptions[0];
				}
			}
			availablePrefixes = Array.isArray(prefixes) ? prefixes : ["Default"];
			baseDir = base;
			systemStatus = sysStatus;

			if (!prefixPath) {
				prefixPath = baseDir + "/Default";
				selectedPrefixName = "Default";
				await loadConfigForPrefix("Default");
			}
		} catch (err) {
			console.error("Failed to initialize:", err);
		} finally {
			isLoadingProton = false;
			mounted = true;
		}
	});

	$: if (mounted) {
		runState.set({ gamePath, gameIcon, prefixPath, selectedPrefixName, selectedProton, options });
	}

	async function handlePrefixChange(name: string) {
		if (name !== "Custom...") {
			prefixPath = baseDir + "/" + name;
			selectedPrefixName = name;
			await loadConfigForPrefix(name);
		}
	}

	async function handleBrowseGame() {
		try {
			const path = await PickFile();
			if (path) {
				gamePath = path;
				options.GamePath = path;
				iconLoadFailed = false;
				const icon = await GetExeIcon(path);
				if (icon) {
					gameIcon = icon;
				} else {
					gameIcon = "";
					iconLoadFailed = true;
					notifications.info(
						"Could not extract icon. Please install 'icoutils' or 'icoextract' to display game icons.",
					);
				}
				await loadConfigForGame(path);
			}
		} catch (err) {
			console.error("Error loading game:", err);
			iconLoadFailed = true;
		}
	}

	async function handleBrowsePrefix() {
		try {
			const path = await PickFolder();
			if (path) {
				prefixPath = path;
				selectedPrefixName = "Custom...";
			}
		} catch (err) {
			console.error(err);
		}
	}

	function handleProtonChange(value: string) {
		selectedProton = value;
	}

	async function handleLaunch() {
		if (!gamePath) {
			notifications.add("Please select a game executable.", "error");
			return;
		}

		if (options.EnableLsfgVk && !options.LsfgDllPath) {
			notifications.add("LSFG-VK requires Lossless.dll.", "error");
			return;
		}

		missingToolsList = [];
		if (options.EnableGamescope && !systemStatus.hasGamescope) missingToolsList.push("Gamescope");
		if (options.EnableMangoHud && !systemStatus.hasMangoHud) missingToolsList.push("MangoHud");
		if (options.EnableGamemode && !systemStatus.hasGameMode) missingToolsList.push("GameMode");
		if (missingToolsList.length > 0) {
			showValidationModal = true;
			return;
		}
		await proceedToLaunch();
	}

	async function proceedToLaunch() {
		showValidationModal = false;
		const tool = protonVersions.find((p) => p.DisplayName === selectedProton);
		let cleanName = selectedProton;
		if (cleanName.startsWith("(Steam) ")) {
			cleanName = cleanName.substring(8);
		}

		options.GamePath = gamePath;
		options.PrefixPath = prefixPath;
		options.ProtonPattern = cleanName;
		options.ProtonPath = tool ? tool.Path : "";

		try {
			await RunGame(options, showLogsWindow);
			WindowHide();
		} catch (err) {
			console.error("Launch failed:", err);
			notifications.add(`Launch failed: ${err}`, "error");
		}
	}
</script>

<div class="run-container">
	<div class="header-row">
		<h1 class="page-title">Launch Configuration</h1>
	</div>

	<div class="form-container">
		<div class="form-group game-exe-group">
			<label for="gameExe">Game Executable</label>
			<div class="game-exe-wrapper">
				<div class="game-icon-display">
					{#if gameIcon && !iconLoadFailed}
						<img
							src={gameIcon}
							alt="Game Icon"
							class="game-icon"
							on:load={() => {
								console.log("Icon loaded successfully");
								iconLoadFailed = false;
							}}
							on:error={(e) => {
								console.error("Icon failed to load:", e);
								iconLoadFailed = true;
							}}
						/>
					{:else}
						<div class="game-icon-placeholder">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="32"
								height="32"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
							>
								<rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect>
								<line x1="8" y1="21" x2="16" y2="21"></line>
								<line x1="12" y1="17" x2="12" y2="21"></line>
							</svg>
						</div>
					{/if}
				</div>
				<div class="input-group game-exe-input-group">
					<input
						id="gameExe"
						type="text"
						bind:value={gamePath}
						placeholder="Select .exe file..."
						class="input"
					/>
					<button on:click={handleBrowseGame} class="btn">Browse</button>
				</div>
			</div>
		</div>

		<div class="form-group">
			<label for="winePrefix">WINEPREFIX</label>
			<div class="input-group">
				<div class="dropdown-wrapper">
					<Dropdown
						options={[...availablePrefixes, "Custom..."]}
						bind:value={selectedPrefixName}
						onChange={handlePrefixChange}
					/>
				</div>
				<button on:click={handleBrowsePrefix} class="btn">Browse</button>
			</div>
			{#if selectedPrefixName === "Custom..." || !prefixPath.startsWith(baseDir)}
				<div class="path-display">{prefixPath}</div>
			{/if}
		</div>

		<div class="form-group">
			<label for="protonVersion">Proton Version</label>
			<div id="protonVersion">
				<Dropdown
					options={protonOptions}
					bind:value={selectedProton}
					placeholder={isLoadingProton ? "Scanning..." : "Select Version"}
					disabled={isLoadingProton}
					onChange={handleProtonChange}
				/>
			</div>
		</div>

		<ConfigForm bind:options />

		<div class="form-group">
			<SlideButton bind:checked={showLogsWindow} label="Show Logs" subtitle="Open logs in terminal" />
		</div>

		<Modal show={showValidationModal} title="Missing Dependencies" onClose={() => (showValidationModal = false)}>
			<div class="warning-modal-content">
				<div class="warning-icon">⚠️</div>
				<p>The following requested features are not installed on your system:</p>
				<div class="missing-list">
					{#each missingToolsList as tool}
						<span class="tool-tag">{tool}</span>
					{/each}
				</div>
				<p class="question">Do you want to launch the game without these features?</p>
				<div class="modal-actions">
					<button class="btn secondary" on:click={() => (showValidationModal = false)}>Cancel</button>
					<button class="btn primary" on:click={proceedToLaunch}>Launch Anyway</button>
				</div>
			</div>
		</Modal>

		<div class="action-area">
			<button class="btn primary launch-btn" on:click={handleLaunch}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2.5"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="run-icon"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg
				>
				<span>LAUNCH GAME</span>
			</button>
		</div>
	</div>
</div>

<style lang="scss">
	.run-container {
		display: flex;
		flex-direction: column;
		height: 100%;
		padding: 32px;
		overflow: hidden;
	}
	.form-container {
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 24px;
		overflow-y: auto;
		padding-right: 8px;
	}
	.game-exe-group {
		.game-exe-wrapper {
			display: flex;
			gap: 16px;
			align-items: flex-end;
		}
		.game-exe-input-group {
			flex: 1;
			flex-direction: row;
		}
	}
	.game-icon-display {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 80px;
		height: 80px;
		border-radius: 12px;
		background: rgba(0, 0, 0, 0.3);
		border: 1px solid var(--glass-border);
		flex-shrink: 0;
		.game-icon {
			width: 100%;
			height: 100%;
			border-radius: 12px;
			object-fit: contain;
			padding: 4px;
		}
		.game-icon-placeholder {
			width: 100%;
			height: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
			color: var(--text-muted);
			svg {
				width: 40px;
				height: 40px;
				opacity: 0.5;
			}
		}
	}
	.form-group label {
		display: block;
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--text-muted);
		margin-bottom: 8px;
	}
	.input-group {
		display: flex;
		gap: 12px;
		width: 100%;
		.input {
			flex: 1;
		}
		.dropdown-wrapper {
			flex: 1;
		}
	}
	.path-display {
		margin-top: 8px;
		font-size: 0.75rem;
		color: var(--text-dim);
		word-break: break-all;
		padding: 8px;
		background: rgba(0, 0, 0, 0.2);
		border-radius: 6px;
	}
	.warning-modal-content {
		text-align: center;
		.warning-icon {
			font-size: 3rem;
			margin-bottom: 16px;
		}
		p {
			color: var(--text-main);
			line-height: 1.5;
		}
		.missing-list {
			margin: 16px 0;
			display: flex;
			flex-wrap: wrap;
			justify-content: center;
			gap: 12px;
			.tool-tag {
				background: rgba(239, 68, 68, 0.1);
				color: #ef4444;
				padding: 6px 16px;
				border-radius: 20px;
				font-size: 0.9rem;
				border: 1px solid rgba(239, 68, 68, 0.2);
				font-weight: bold;
			}
		}
		.question {
			margin-top: 24px;
			font-weight: 600;
			color: var(--accent-secondary);
		}
	}
	.modal-actions {
		display: flex;
		gap: 12px;
		margin-top: 32px;
		button {
			flex: 1;
			padding: 12px;
			font-weight: 600;
		}
	}
	.action-area {
		margin-top: 16px;
		margin-bottom: 32px;
	}
	.launch-btn {
		width: 100%;
		padding: 18px;
		font-size: 1.25rem;
		font-weight: 800;
		border-radius: 14px;
		text-transform: uppercase;
		letter-spacing: 2px;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 12px;
		background: #ffffff;
		color: #000000 !important;
		border: none;
		box-shadow: 0 8px 24px rgba(255, 255, 255, 0.15);
		cursor: pointer;
		transition: all 0.2s;

		&:hover {
			background: #f4f4f5;
			transform: translateY(-2px);
			box-shadow: 0 12px 32px rgba(255, 255, 255, 0.25);
		}

		&:active {
			transform: translateY(0);
		}

		.run-icon {
			width: 24px;
			height: 24px;
			fill: #000000;
			stroke: #000000;
		}
	}
	.input {
		background: rgba(0, 0, 0, 0.25);
		border: 1px solid var(--glass-border);
		color: var(--text-main);
		padding: 12px 16px;
		border-radius: 10px;
		outline: none;
		transition: all 0.2s;
		&:focus {
			border-color: var(--text-muted);
			background: rgba(0, 0, 0, 0.4);
		}
	}
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 10px 20px;
		border-radius: 10px;
		font-weight: 600;
		font-size: 0.9rem;
		cursor: pointer;
		transition: all 0.2s ease;
		border: 1px solid var(--glass-border);
		background: rgba(255, 255, 255, 0.05);
		color: var(--text-main);
		&:hover {
			background: rgba(255, 255, 255, 0.1);
			border-color: var(--glass-border-bright);
		}
		&.primary {
			background: var(--accent-primary);
			border: none;
			color: #000;
			&:hover {
				background: var(--accent-secondary);
			}
		}
	}
	.page-title {
		font-size: 2rem;
		font-weight: bold;
		color: var(--text-main);
		margin: 0 0 24px 0;
	}
</style>
