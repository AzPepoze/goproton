<script lang="ts">
	import { onMount, onDestroy } from "svelte";
	import {
		GetSystemInfo,
		GetSystemUsage,
		CleanupProcesses,
		GetShaderCacheSize,
		ClearShaderCache,
	} from "../../wailsjs/go/backend/App";
	import type { core } from "../../wailsjs/go/models";
	import trashIcon from "../icons/trash.svg";
	import rocketIcon from "../icons/rocket.svg";
	import StatusUtilityButton from "./StatusUtilityButton.svelte";

	let isExpanded = false;
	let isCleaning = false;
	let isClearingCache = false;
	let showCleanupSuccess = false;
	let showCacheSuccess = false;
	let sysInfo: core.SystemInfo = { os: "", kernel: "", cpu: "", gpu: "", ram: "", driver: "" };
	let sysUsage: core.SystemUsage = { cpu: "0%", ram: "0%", gpu: "0%" };
	let shaderCacheSize = "...";
	let usageInterval;

	async function fetchData() {
		try {
			const [info, usage, cache] = await Promise.all([
				GetSystemInfo(),
				GetSystemUsage(),
				GetShaderCacheSize(),
			]);
			sysInfo = info;
			sysUsage = usage;
			shaderCacheSize = cache;
		} catch (err) {
			console.error("Failed to fetch status drawer data:", err);
		}
	}

	onMount(() => {
		fetchData();
		usageInterval = setInterval(async () => {
			try {
				sysUsage = await GetSystemUsage();
			} catch (e) {}
		}, 3000);
	});

	onDestroy(() => {
		if (usageInterval) clearInterval(usageInterval);
	});

	async function handleCleanup() {
		if (isCleaning) return;
		isCleaning = true;
		showCleanupSuccess = false;
		try {
			await CleanupProcesses();
			await fetchData();
			// Faster pop
			setTimeout(() => {
				showCleanupSuccess = true;
				// Longer visibility
				setTimeout(() => {
					showCleanupSuccess = false;
				}, 2000);
			}, 100);
		} catch (err) {
			console.error(`Cleanup failed: ${err}`);
		} finally {
			setTimeout(() => {
				isCleaning = false;
			}, 1500);
		}
	}

	async function handleClearCache() {
		if (isClearingCache) return;
		isClearingCache = true;
		showCacheSuccess = false;
		try {
			await ClearShaderCache();
			const newCache = await GetShaderCacheSize();
			shaderCacheSize = newCache;
			setTimeout(() => {
				showCacheSuccess = true;
				setTimeout(() => {
					showCacheSuccess = false;
				}, 2000);
			}, 100);
		} catch (err) {
			console.error(`Failed to clear cache: ${err}`);
		} finally {
			setTimeout(() => {
				isClearingCache = false;
			}, 1500);
		}
	}
</script>

<div class="status-drawer-wrapper" class:expanded={isExpanded}>
	<button class="toggle-btn" on:click={() => (isExpanded = !isExpanded)}>
		<span class="trigger-text">{isExpanded ? "CLOSE DRAWER" : "SYSTEM STATUS & UTILITIES"}</span>
	</button>

	<div class="drawer-content">
		<div class="status-grid">
			<!-- OS -->
			<div class="status-box">
				<div class="box-header">
					<div class="icon-label">
						<svg
							class="mini-icon"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2.5"
							stroke-linecap="round"
							stroke-linejoin="round"
							><rect width="20" height="14" x="2" y="3" rx="2" /><line
								x1="8"
								x2="16"
								y1="21"
								y2="21"
							/><line x1="12" x2="12" y1="17" y2="21" /></svg
						>
						<span class="label">SYSTEM</span>
					</div>
				</div>
				<div class="system-info">
					<span class="os-text" title={sysInfo.os}>{sysInfo.os}</span>
					<span class="kernel-text">Kernel: {sysInfo.kernel}</span>
				</div>
			</div>

			<!-- CPU -->
			<div class="status-box">
				<div class="box-header">
					<div class="icon-label">
						<svg
							class="mini-icon"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2.5"
							stroke-linecap="round"
							stroke-linejoin="round"
							><rect width="16" height="16" x="4" y="4" rx="2" /><rect
								width="6"
								height="6"
								x="9"
								y="9"
							/><path d="M15 2v2" /><path d="M15 20v2" /><path d="M2 15h2" /><path
								d="M2 9h2"
							/><path d="M20 15h2" /><path d="M20 9h2" /><path d="M9 2v2" /><path
								d="M9 20v2"
							/></svg
						>
						<span class="label">CPU</span>
					</div>
					<span class="usage">{sysUsage.cpu}</span>
				</div>
				<div class="progress-bg">
					<div class="progress-fill" style="width: {sysUsage.cpu}"></div>
				</div>
				<span class="info-text" title={sysInfo.cpu}>{sysInfo.cpu}</span>
			</div>

			<!-- RAM -->
			<div class="status-box">
				<div class="box-header">
					<div class="icon-label">
						<svg
							class="mini-icon"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2.5"
							stroke-linecap="round"
							stroke-linejoin="round"
							><path d="M6 19v-3" /><path d="M10 19v-3" /><path d="M14 19v-3" /><path
								d="M18 19v-3"
							/><path d="M8 11V9" /><path d="M16 11V9" /><rect
								width="18"
								height="12"
								x="3"
								y="7"
								rx="2"
							/><path d="M3 13h18" /></svg
						>
						<span class="label">RAM</span>
					</div>
					<span class="usage"
						>{sysUsage.ram.includes("(")
							? sysUsage.ram.split("(").pop().replace(")", "")
							: "0%"}</span
					>
				</div>
				<div class="progress-bg">
					<div
						class="progress-fill"
						style="width: {sysUsage.ram.includes('(')
							? sysUsage.ram.split('(').pop().replace(')', '')
							: '0%'}"
					></div>
				</div>
				<span class="info-text">{sysUsage.ram.split(" / ")[0]} used</span>
			</div>

			<!-- GPU -->
			<div class="status-box">
				<div class="box-header">
					<div class="icon-label">
						<svg
							class="mini-icon"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2.5"
							stroke-linecap="round"
							stroke-linejoin="round"
							><path d="M18 8V6a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h2" /><rect
								width="12"
								height="8"
								x="10"
								y="10"
								rx="2"
							/><path d="M14 10V8" /><path d="M18 10V8" /><path d="M14 20v-2" /><path
								d="M18 20v-2"
							/><path d="M22 14h-2" /><path d="M22 18h-2" /></svg
						>
						<span class="label">GPU</span>
					</div>
					<span class="usage">{sysUsage.gpu}</span>
				</div>
				<div class="progress-bg">
					<div
						class="progress-fill"
						style="width: {sysUsage.gpu}; background: var(--accent-secondary, #b197fc)"
					></div>
				</div>
				<span class="info-text" title="{sysInfo.gpu} ({sysInfo.driver})">{sysInfo.gpu}</span>
			</div>
		</div>

		<div class="divider"></div>

		<div class="utilities-row">
			<StatusUtilityButton
				icon={trashIcon}
				title="Cleanup System"
				subtitle="Terminate all running processes"
				isPulsing={isCleaning}
				showSuccess={showCleanupSuccess}
				btnClass="cleanup"
				on:click={handleCleanup}
			/>

			<StatusUtilityButton
				icon={rocketIcon}
				title="Clear Shader Cache"
				subtitle={shaderCacheSize}
				isPulsing={isClearingCache}
				showSuccess={showCacheSuccess}
				btnClass="cache"
				on:click={handleClearCache}
			/>
		</div>
	</div>
</div>

<style lang="scss">
	.status-drawer-wrapper {
		position: fixed;
		bottom: 20px;
		background: rgba(18, 18, 22, 1);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 20px;
		transform: translateY(calc(100% - 55px));
		transition: all 0.5s cubic-bezier(0.23, 1, 0.32, 1);
		z-index: 100;
		padding: 0 20px 20px 20px;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.4);
		overflow: hidden;
		margin-right: 20px;
		width: -webkit-fill-available;

		&.expanded {
			transform: translateY(0);
		}
	}

	.toggle-btn {
		width: 100%;
		height: 48px;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		background: rgba(60, 60, 65, 0.4);
		border: 1px solid rgba(255, 255, 255, 0.05);
		border-radius: 12px;
		margin: 10px 0;
		transition: all 0.2s ease;

		.trigger-text {
			font-size: 0.8rem;
			font-weight: 700;
			color: rgba(255, 255, 255, 0.6);
			letter-spacing: 1px;
			text-transform: uppercase;
		}

		&:hover {
			background: rgba(80, 80, 85, 0.6);
			border-color: rgba(255, 255, 255, 0.1);

			.trigger-text {
				color: #fff;
			}
		}

		&:active {
			background: rgba(45, 45, 50, 0.6);
			transform: scale(0.995);
		}
	}

	.drawer-content {
		padding-top: 10px;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.status-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 16px;
	}

	.status-box {
		background: rgba(255, 255, 255, 0.03);
		padding: 12px;
		border-radius: 12px;
		border: 1px solid rgba(255, 255, 255, 0.05);
		display: flex;
		flex-direction: column;
		min-width: 0;

		.box-header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 8px;

			.icon-label {
				display: flex;
				align-items: center;
				gap: 6px;
				color: var(--text-dim);

				.mini-icon {
					width: 14px;
					height: 14px;
					opacity: 0.8;
				}

				.label {
					font-size: 0.7rem;
					font-weight: 900;
				}
			}

			.usage {
				font-size: 0.9rem;
				font-weight: 700;
				color: var(--accent-primary);
			}
		}

		.info-text {
			display: block;
			font-size: 0.7rem;
			color: var(--text-muted);
			margin-top: 8px;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}

		.system-info {
			display: flex;
			flex-direction: column;
			gap: 2px;
			overflow: hidden;

			.os-text {
				font-size: 0.8rem;
				font-weight: 700;
				color: var(--text-main);
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
			}

			.kernel-text {
				font-size: 0.65rem;
				color: var(--text-muted);
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
			}
		}
	}

	.progress-bg {
		height: 4px;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 2px;
		overflow: hidden;

		.progress-fill {
			height: 100%;
			background: var(--accent-primary);
			transition: width 0.3s ease;
		}
	}

	.divider {
		height: 1px;
		background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
	}

	.utilities-row {
		display: flex;
		gap: 16px;
		align-items: center;
	}
</style>
