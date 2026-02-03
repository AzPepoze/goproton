<script lang="ts">
	import { onMount } from "svelte";
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
		GetInitialLauncherPath,
		DetectLosslessDll,
		GetExeIcon,
	} from "../api";
	import type { launcher } from "../models";
	import Dropdown from "../components/Dropdown.svelte";
	import ConfigForm from "../components/ConfigForm.svelte";
	import SlideButton from "../components/SlideButton.svelte";
	import Modal from "../components/Modal.svelte";
	import ExecutableSelector from "../components/ExecutableSelector.svelte";
	import { notifications } from "../notificationStore";
	import { runState } from "../stores/runState";
	import { get } from "svelte/store";
	import { WindowHide } from "../api";

	// Component State
	let mounted = false;

	// Game Selection
	// selectedGameExePath: tracks only when user EXPLICITLY selects a separate game exe
	// This is DIFFERENT from options.GamePath which is the actual path to execute
	let selectedGameExePath = "";
	let gameIcon = "";
	let launcherIcon = "";
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

	// Config
	let options: launcher.LaunchOptions = {
		GamePath: "",
		LauncherPath: "",
		UseGameExe: false,
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
		LsfgGpu: "",
		LsfgFlowScale: "0.8",
		LsfgPacing: "none",
		LsfgAllowFp16: false,
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

			// Auto-detect Lossless.dll if not already set
			if (!options.LsfgDllPath) {
				try {
					const dll = await DetectLosslessDll();
					if (dll) {
						options.LsfgDllPath = dll;
					}
				} catch (err) {
					console.error("Failed to detect Lossless.dll:", err);
				}
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
		options.LsfgGpu = config.LsfgGpu || "";
		options.LsfgFlowScale = config.LsfgFlowScale || "0.8";
		options.LsfgPacing = config.LsfgPacing || "none";
		options.LsfgAllowFp16 = config.LsfgAllowFp16 || false;
		// IMPORTANT: Only apply LauncherPath from config if we don't already have one selected
		if (!options.LauncherPath && config.LauncherPath) {
			options.LauncherPath = config.LauncherPath;
		}
		options.UseGameExe = config.UseGameExe === true; // Default to false (launcher-only)

		// CRITICAL: If UseGameExe is false, GamePath must equal LauncherPath (launcher-only mode)
		// If UseGameExe is true, GamePath should be a separate game exe
		if (!options.UseGameExe && options.LauncherPath) {
			options.GamePath = options.LauncherPath;
			console.log("[CONFIG] UseGameExe=false: enforcing GamePath=LauncherPath (launcher-only mode)");
			console.log("[CONFIG]   GamePath set to:", options.LauncherPath);
		} else if (options.UseGameExe && config.GamePath) {
			// Only load GamePath from config if UseGameExe is true (separate game exe)
			options.GamePath = config.GamePath;
			console.log("[CONFIG] UseGameExe=true: loaded separate game exe from config");
			console.log("[CONFIG]   GamePath set to:", config.GamePath);
		}

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
					selectedGameExePath = s.gamePath;
					options.GamePath = s.gamePath;
				}
				if (s.gameIcon) gameIcon = s.gameIcon;
				if (s.launcherIcon) launcherIcon = s.launcherIcon;
				if (s.prefixPath) prefixPath = s.prefixPath;
				if (s.selectedPrefixName) selectedPrefixName = s.selectedPrefixName;
				if (s.selectedProton) selectedProton = s.selectedProton;
				if (s.options) {
					options = { ...options, ...s.options };
				}
			}

			if (options.LauncherPath) {
				await loadConfigForGame(options.LauncherPath);
				if (!launcherIcon) {
					const icon = await GetExeIcon(options.LauncherPath);
					if (icon) launcherIcon = icon;
				}
			}

			const initialPath = await GetInitialLauncherPath();
			if (initialPath) {
				// Only set as game/launcher if not already set, or if explicitly passed from tray
				if (!options.LauncherPath && !options.GamePath) {
					// No prior state - set initial path as launcher
					options.LauncherPath = initialPath;
					const icon = await GetExeIcon(initialPath);
					if (icon) launcherIcon = icon;
				} else if (!options.GamePath || options.GamePath === options.LauncherPath) {
					// Prior state has launcher but no game - set initial path as game
					selectedGameExePath = initialPath;
					options.GamePath = initialPath;
					const icon = await GetExeIcon(initialPath);
					if (icon) gameIcon = icon;
					await loadConfigForGame(initialPath);
				}
			}

			// Load icons for any paths that don't have icons yet
			if (options.LauncherPath && !launcherIcon) {
				const icon = await GetExeIcon(options.LauncherPath);
				if (icon) launcherIcon = icon;
			}
			if (options.GamePath && !gameIcon) {
				const icon = await GetExeIcon(options.GamePath);
				if (icon) gameIcon = icon;
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

	onMount(() => {
		runState.set({
			gamePath: selectedGameExePath,
			gameIcon,
			launcherIcon,
			prefixPath,
			selectedPrefixName,
			selectedProton,
			options,
		});
	});

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
				console.log("[GAME] Selected game exe:", path);
				console.log("[GAME] Current LauncherPath before game selection:", options.LauncherPath);
				selectedGameExePath = path;
				// Use object spread to trigger Svelte reactivity
				options = { ...options, GamePath: path };
				console.log("[GAME] Set options.GamePath to:", options.GamePath);
				console.log("[GAME] LauncherPath after game selection:", options.LauncherPath);
				console.log("[GAME] Full options object:", JSON.stringify(options));
				// NOTE: Do NOT load config for game exe
				// Game exe is only for LSFG profile matching
				// Configuration is ALWAYS saved under launcher exe path only
			}
		} catch (err) {
			console.error("[GAME] Error loading game:", err);
		}
	}

	async function handleBrowseLauncher() {
		try {
			const path = await PickFile();
			if (path) {
				console.log("[LAUNCHER] Selected launcher exe:", path);
				options = { ...options, LauncherPath: path };
				console.log("[LAUNCHER] Set options.LauncherPath to:", options.LauncherPath);
				console.log("[LAUNCHER] Full options object after assignment:", JSON.stringify(options));

				// Only set GamePath if user has not explicitly selected a separate game exe
				if (!selectedGameExePath) {
					console.log(
						"[LAUNCHER] No separate game exe selected by user, initializing GamePath to launcher",
					);
					options = { ...options, GamePath: path };
					console.log("[LAUNCHER] Set GamePath to launcher path:", options.GamePath);
				} else {
					console.log(
						"[LAUNCHER] User already selected separate game exe, keeping GamePath:",
						selectedGameExePath,
					);
				}

				// Load config for the launcher
				// applyConfigToOptions will enforce UseGameExe if true
				console.log("[LAUNCHER] Loading config for launcher path...");
				await loadConfigForGame(path);
				console.log("[LAUNCHER] Config loaded, final GamePath:", options.GamePath);
			}
		} catch (err) {
			console.error("[LAUNCHER] Error selecting launcher:", err);
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
		if (!options.LauncherPath) {
			notifications.add("Please select a launcher executable.", "error");
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
		console.log("\n============ PROCEED TO LAUNCH ============");

		// DEBUG: Log state at execution time
		console.log("[EXECUTE] Step 1 - Initial state");
		console.log("[EXECUTE]   options.LauncherPath:", options.LauncherPath);
		console.log("[EXECUTE]   options.GamePath:", options.GamePath);
		console.log("[EXECUTE]   selectedGameExePath variable:", selectedGameExePath);
		console.log("[EXECUTE]   Full options:", JSON.stringify(options));

		const tool = protonVersions.find((p) => p.DisplayName === selectedProton);
		let cleanName = selectedProton;
		if (cleanName.startsWith("(Steam) ")) {
			cleanName = cleanName.substring(8);
		}

		// Config is ALWAYS saved to launcher path via SaveGameConfig backend logic
		// GamePath remains the actual game executable to run
		// LauncherPath is provided to SaveGameConfig for config storage
		options.PrefixPath = prefixPath;
		options.ProtonPattern = cleanName;
		options.ProtonPath = tool ? tool.Path : "";

		console.log("[EXECUTE] Step 2 - Final options object before RunGame:");
		console.log(JSON.stringify(options, null, 2));
		console.log("============ ABOUT TO CALL RunGame ============\n");

		try {
			console.log("[EXECUTE] Calling RunGame with LauncherPath:", options.LauncherPath);
			console.log("[EXECUTE] Calling RunGame with GamePath:", options.GamePath);
			console.log("[EXECUTE] Calling RunGame with full options:", JSON.stringify(options, null, 2));
			await RunGame(options, showLogsWindow);
			WindowHide();
		} catch (err) {
			console.error("[EXECUTE] Launch failed:", err);
			notifications.add(`Launch failed: ${err}`, "error");
		}
	}
</script>

<div class="run-container">
	<div class="header-row">
		<h1 class="page-title">Launch Configuration</h1>
	</div>

	<!-- Executable Selector Component -->
	<ExecutableSelector
		launcherPath={options.LauncherPath}
		gamePath={options.GamePath}
		bind:useGameExe={options.UseGameExe}
		bind:launcherIcon
		bind:gameIcon
		onBrowseLauncher={handleBrowseLauncher}
		onBrowseGame={handleBrowseGame}
	/>

	<!-- Main Form Container -->
	<div class="form-container">
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
		padding: 32px;
	}
	.form-container {
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 24px;
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
