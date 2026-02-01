<script lang="ts">
	import SlideButton from "./SlideButton.svelte";
	import Dropdown from "./Dropdown.svelte";
	import Modal from "./Modal.svelte";
	import RangeSlider from "./RangeSlider.svelte";
	import { PickFileCustom, GetTotalRam, DetectLosslessDll } from "../../wailsjs/go/main/App";
	import type { launcher } from "../../wailsjs/go/models";
	import { onMount } from "svelte";

	export let options: launcher.LaunchOptions;

	// Modals
	let showLsfgModal = false;
	let showGamescopeModal = false;
	let showMemoryModal = false;

	let memorySliderValue = 4;
	let systemRamTotal = 16;

	onMount(async () => {
		try {
			const ram = await GetTotalRam();
			if (ram > 0) systemRamTotal = ram;

			if (options.MemoryMinValue) {
				const numericMatch = options.MemoryMinValue.match(/([\d.]+)/);
				if (numericMatch) {
					const val = parseFloat(numericMatch[1]);
					if (options.MemoryMinValue.endsWith("M")) {
						memorySliderValue = val / 1024;
					} else {
						memorySliderValue = val;
					}
				}
			}

			// Auto-detect DLL on mount
			if (!options.LsfgDllPath) {
				const dll = await DetectLosslessDll();
				if (dll) options.LsfgDllPath = dll;
			}
		} catch (e) {
			console.error(e);
		}
	});

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
</script>

<div class="config-form">
	<div class="form-group">
		<label for="customArgs">Custom Arguments</label>
		<input
			id="customArgs"
			type="text"
			class="input"
			bind:value={options.CustomArgs}
			placeholder="e.g. -windowed -novid"
		/>
	</div>

	<div class="toggles-grid">
		<SlideButton bind:checked={options.EnableMangoHud} label="MangoHud" subtitle="Performance overlay" />
		<SlideButton bind:checked={options.EnableGamemode} label="GameMode" subtitle="Optimize priorities" />
		<SlideButton
			bind:checked={options.EnableLsfgVk}
			label="LSFG-VK"
			subtitle="Lossless Scaling Frame Generation"
			hasConfig={true}
			onConfig={() => (showLsfgModal = true)}
		/>
		<SlideButton
			bind:checked={options.EnableGamescope}
			label="Gamescope"
			subtitle="Micro-compositor"
			hasConfig={true}
			onConfig={() => (showGamescopeModal = true)}
		/>
		<SlideButton
			bind:checked={options.EnableMemoryMin}
			label="Memory Protect"
			subtitle="Prevent Swap (Min RAM)"
			hasConfig={true}
			onConfig={() => (showMemoryModal = true)}
		/>
	</div>

	<!-- LSFG Settings Modal -->
	<Modal show={showLsfgModal} title="LSFG-VK Configuration" onClose={() => (showLsfgModal = false)}>
		<div class="modal-form">
			<div class="form-group">
				<label for="lsfgDll">
					Lossless.dll Path (from Steam)
					{#if options.LsfgDllPath}
						<span class="status-badge success">Detected</span>
					{:else}
						<span class="status-badge error">Not Found</span>
					{/if}
				</label>
				<div class="input-group">
					<input
						id="lsfgDll"
						type="text"
						class="input sm"
						bind:value={options.LsfgDllPath}
						placeholder="Path to Lossless.dll..."
					/>
					<button class="btn sm" on:click={handleBrowseDll}>Browse</button>
				</div>
			</div>
			<div class="form-group">
				<label for="fpsMultiplier">FPS Multiplier</label>
				<div id="fpsMultiplier">
					<Dropdown
						options={["2", "3", "4"]}
						bind:value={options.LsfgMultiplier}
						onChange={(val) => (options.LsfgMultiplier = val)}
					/>
				</div>
			</div>
			<div class="form-group">
				<SlideButton
					bind:checked={options.LsfgPerfMode}
					label="Performance Mode"
					subtitle="Faster frame generation, slight quality loss"
				/>
			</div>
		</div>
	</Modal>

	<!-- Gamescope Settings Modal -->
	<Modal show={showGamescopeModal} title="Gamescope Configuration" onClose={() => (showGamescopeModal = false)}>
		<div class="modal-form">
			<div class="form-group">
				<label for="gamescopeWidth">Width (px)</label>
				<input
					id="gamescopeWidth"
					type="text"
					class="input"
					bind:value={options.GamescopeW}
					placeholder="e.g. 1920"
				/>
			</div>
			<div class="form-group">
				<label for="gamescopeHeight">Height (px)</label>
				<input
					id="gamescopeHeight"
					type="text"
					class="input"
					bind:value={options.GamescopeH}
					placeholder="e.g. 1080"
				/>
			</div>
			<div class="form-group">
				<label for="gamescopeRefresh">Refresh Rate (Hz)</label>
				<input
					id="gamescopeRefresh"
					type="text"
					class="input"
					bind:value={options.GamescopeR}
					placeholder="e.g. 60"
				/>
			</div>
			<p class="note">Note: Mouse visibility fix enabled automatically.</p>
		</div>
	</Modal>

	<!-- Memory Settings Modal -->
	<Modal show={showMemoryModal} title="Memory Protection" onClose={() => (showMemoryModal = false)}>
		<div class="modal-form">
			<div class="form-group">
				<label for="memorySlider">Minimum RAM Allocation</label>
				<RangeSlider
					value={memorySliderValue}
					max={systemRamTotal}
					snapValues={[2, 4, 6, 8, 12, 16, 24, 32, 48, 64]}
					onChange={(changedValue) => {
						memorySliderValue = changedValue;
						options.MemoryMinValue = changedValue + "G";
					}}
				/>
			</div>
			<p class="note">
				Guarantees {options.MemoryMinValue} of physical RAM for the game process, preventing swapping.
				<br />Values in Red Zone (>75%) might cause system instability.
			</p>
		</div>
	</Modal>
</div>

<style lang="scss">
	.toggles-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
		gap: 16px;
		margin-top: 8px;
	}
	.modal-form {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}
	.form-group label {
		display: block;
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--text-muted);
		margin-bottom: 8px;
	}

	.form-group .input {
		width: 100%;
		display: block;
	}
	.input-group {
		display: flex;
		gap: 12px;
		width: 100%;
		.input {
			flex: 1;
		}
	}
	.input.sm {
		padding: 6px 10px;
		font-size: 0.875rem;
	}
	.btn.sm {
		padding: 6px 12px;
		font-size: 0.8rem;
	}

	.status-badge {
		font-size: 0.65rem;
		padding: 2px 6px;
		border-radius: 4px;
		margin-left: 8px;
		text-transform: uppercase;
		font-weight: bold;
		&.success {
			background: rgba(16, 185, 129, 0.2);
			color: #10b981;
		}
		&.error {
			background: rgba(239, 68, 68, 0.2);
			color: #ef4444;
		}
	}
	.note {
		font-size: 0.75rem;
		color: var(--text-muted);
		font-style: italic;
		margin-top: 8px;
	}
</style>
