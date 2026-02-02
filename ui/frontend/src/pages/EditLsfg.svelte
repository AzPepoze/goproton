<script lang="ts">
	import Modal from "../components/Modal.svelte";
	import LsfgConfigForm from "../components/LsfgConfigForm.svelte";
	import {
		GetListGpus,
		GetLsfgProfileForGame,
		GetInitialGamePath,
		SaveLsfgProfile,
		DetectLosslessDll,
		CloseWindow,
		PickFileCustom,
	} from "../../wailsjs/go/main/App";
	import type { launcher } from "../../wailsjs/go/models";
	import { onMount } from "svelte";

	export let gamePath = "";

	let options: launcher.LaunchOptions = {
		GamePath: "",
		LauncherPath: "",
		UseGameExe: true,
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
		EnableLsfgVk: true,
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

	let loading = true;
	let error = "";
	let gpuList: string[] = [];
	let saving = false;
	let saveSuccess = false;
	let saveError = "";

	onMount(async () => {
		try {
			// Use passed gamePath or get from environment
			let currentGamePath = gamePath;
			if (!currentGamePath) {
				currentGamePath = await GetInitialGamePath();
			}

			if (!currentGamePath) {
				error = "No game path provided";
				loading = false;
				return;
			}

			options.GamePath = currentGamePath;

			// Load profile data
			const data = await GetLsfgProfileForGame(currentGamePath);
			if (data) {
				options.LsfgMultiplier = String(data.multiplier || 2);
				options.LsfgPerfMode = data.performanceMode || false;
				options.LsfgGpu = data.gpu || "";
				options.LsfgFlowScale = String(data.flowScale || 0.8);
				options.LsfgPacing = data.pacing || "none";
				options.LsfgDllPath = data.dllPath || "";
				options.LsfgAllowFp16 = data.allowFp16 || false;
			}

			// Auto-detect DLL if not set
			if (!options.LsfgDllPath) {
				try {
					const detected = await DetectLosslessDll();
					if (detected) {
						options.LsfgDllPath = detected;
						console.log("Auto-detected DLL:", detected);
					}
				} catch (e) {
					console.error("Failed to detect DLL:", e);
				}
			}

			// Load available GPUs
			try {
				const gpus = await GetListGpus();
				if (gpus && gpus.length > 0) {
					gpuList = gpus;
				}
			} catch (e) {
				console.error("Failed to detect GPUs:", e);
			}

			loading = false;
		} catch (err) {
			error = String(err);
			loading = false;
		}
	});

	async function handleApply() {
		saving = true;
		saveSuccess = false;
		saveError = "";

		try {
			await SaveLsfgProfile(
				options.GamePath,
				parseInt(options.LsfgMultiplier) || 2,
				options.LsfgPerfMode,
				options.LsfgDllPath,
				options.LsfgGpu,
				options.LsfgFlowScale,
				options.LsfgPacing,
				options.LsfgAllowFp16,
			);
			saveSuccess = true;
			saveError = "";

			// Auto-dismiss success message after 2 seconds
			setTimeout(() => {
				saveSuccess = false;
			}, 2000);
		} catch (err) {
			saveError = String(err);
			saveSuccess = false;
		} finally {
			saving = false;
		}
	}

	async function handleBrowseDll() {
		try {
			const path = await PickFileCustom("Select Lossless.dll", [
				{ DisplayName: "Lossless.dll", Pattern: "Lossless.dll" },
			]);
			if (path) options.LsfgDllPath = path;
		} catch (err) {
			console.error(err);
		}
	}

	function handleClose() {
		CloseWindow();
	}
</script>

{#if loading}
	<Modal show={true} title="LSFG-VK Configuration" onClose={handleClose} fullscreen={true}>
		<div class="loading-container">
			<p>Loading LSFG configuration...</p>
		</div>
	</Modal>
{:else if error}
	<Modal show={true} title="LSFG-VK Configuration" onClose={handleClose} fullscreen={true}>
		<div class="error-container">
			<p>Error loading profile: {error}</p>
		</div>
	</Modal>
{:else}
	<Modal show={true} title="LSFG-VK Configuration" onClose={handleClose} fullscreen={true}>
		<div class="modal-content">
			<div class="profile-info">
				<p class="game-path">{gamePath}</p>
			</div>

			<LsfgConfigForm {options} {gpuList} onDllBrowse={handleBrowseDll} />

			{#if saveSuccess}
				<div class="status-message success">✓ Configuration saved successfully</div>
			{:else if saveError}
				<div class="status-message error">✗ Error: {saveError}</div>
			{/if}

			<div class="actions">
				<button class="btn secondary" on:click={handleClose}>Close</button>
				<button class="btn primary" on:click={handleApply} disabled={saving}>
					{saving ? "Saving..." : "Apply"}
				</button>
			</div>
		</div>
	</Modal>
{/if}

<style lang="scss">
	.loading-container,
	.error-container {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 400px;
		text-align: center;

		p {
			color: var(--text-muted);
			font-size: 1.1rem;
		}
	}

	.error-container p {
		color: #ef4444;
	}

	.modal-content {
		max-width: 800px;
		margin: 0 auto;
	}

	.profile-info {
		margin-bottom: 16px;

		.game-path {
			margin: 0;
			color: var(--text-dim);
			font-size: 0.9rem;
			word-break: break-all;
		}
	}

	.actions {
		display: flex;
		justify-content: flex-end;
		gap: 12px;
		margin-top: 32px;
		padding-top: 24px;
		border-top: 1px solid var(--glass-border);
	}

	.status-message {
		padding: 12px 16px;
		border-radius: 8px;
		margin: 16px 0;
		font-weight: 500;
		text-align: center;

		&.success {
			background: rgba(34, 197, 94, 0.1);
			color: #22c55e;
		}

		&.error {
			background: rgba(239, 68, 68, 0.1);
			color: #ef4444;
		}
	}

	.btn {
		padding: 10px 20px;
		border: none;
		border-radius: 8px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
		font-size: 0.9rem;

		&:disabled {
			opacity: 0.5;
			cursor: not-allowed;
		}

		&.primary {
			background: #3b82f6;
			color: white;

			&:hover:not(:disabled) {
				background: #2563eb;
			}
		}

		&.secondary {
			background: transparent;
			border: 1px solid var(--glass-border);
			color: var(--text-main);

			&:hover:not(:disabled) {
				background: rgba(255, 255, 255, 0.05);
				border-color: var(--text-main);
			}
		}
	}
</style>
