<script lang="ts">
	import { onMount } from "svelte";
	import { ScanProtonVersions } from "../../wailsjs/go/backend/App";
	import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
	import type { core } from "../../wailsjs/go/models";
	import { notifications } from "../notificationStore";

	import steamIcon from "../icons/steam.png";
	import externalIcon from "../icons/protron_forked.png";

	let protonVersions: core.ProtonTool[] = [];
	let isLoading = true;

	onMount(async () => {
		loadInstalledVersions();
	});

	async function loadInstalledVersions() {
		isLoading = true;
		try {
			protonVersions = await ScanProtonVersions();
		} catch (err) {
			console.error(err);
			notifications.error("Failed to scan versions");
		} finally {
			isLoading = false;
		}
	}

	function openExternal(url: string) {
		BrowserOpenURL(url);
	}
</script>

<div class="versions-container">
	<div class="header">
		<div class="title-group">
			<h1 class="page-title">Proton Versions</h1>
			<div class="count-badge">{protonVersions.length} Installed</div>
		</div>

		<div class="actions">
			<button class="btn secondary" on:click={() => openExternal("https://protondb.com")}> ProtonDB </button>
			<button class="btn secondary" on:click={() => openExternal("https://github.com/Vulong000/ProtonPlus")}>
				ProtonPlus
			</button>
			<button class="btn secondary" on:click={() => openExternal("https://github.com/DavidoTek/ProtonUp-Qt")}>
				ProtonUp
			</button>
		</div>
	</div>

	<div class="versions-list">
		{#if isLoading}
			<div class="loading">Scanning...</div>
		{:else}
			{#each protonVersions as tool}
				<div class="version-card">
					<div class="icon">
						<img src={tool.IsSteam ? steamIcon : externalIcon} alt="tool" class="tool-icon" />
					</div>
					<div class="info">
						<div class="name">{tool.Name}</div>
						<div class="path" title={tool.Path}>{tool.Path}</div>
					</div>
					<div class="type-badge" class:steam={tool.IsSteam}>
						{tool.IsSteam ? "Steam" : "External"}
					</div>
				</div>
			{/each}
		{/if}
	</div>
</div>

<style lang="scss">
	.versions-container {
		padding: 32px;
		height: 100%;
		display: flex;
		flex-direction: column;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 32px;
	}

	.title-group {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.actions {
		display: flex;
		gap: 12px;
	}

	.page-title {
		font-size: 2rem;
		font-weight: 800;
		color: var(--text-main);
		margin: 0;
	}

	.count-badge {
		background: rgba(255, 255, 255, 0.1);
		padding: 4px 12px;
		border-radius: 20px;
		font-size: 0.8rem;
		color: var(--text-muted);
		border: 1px solid var(--glass-border);
	}

	.versions-list {
		display: flex;
		flex-direction: column;
		gap: 16px;
		overflow-y: auto;
		padding-right: 8px;
		flex: 1;
	}

	.version-card {
		display: flex;
		align-items: center;
		gap: 20px;
		padding: 20px;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid var(--glass-border);
		border-radius: 16px;
		transition: all 0.2s;

		&:hover {
			background: rgba(255, 255, 255, 0.06);
			border-color: var(--glass-border-bright);
			transform: translateX(4px);
		}

		.icon {
			font-size: 1.5rem;
			background: rgba(0, 0, 0, 0.2);
			width: 48px;
			height: 48px;
			display: flex;
			align-items: center;
			justify-content: center;
			border-radius: 12px;

			.tool-icon {
				width: 24px;
				height: 24px;
				filter: brightness(0) invert(1);
				opacity: 0.8;
			}
		}

		.info {
			flex: 1;
			overflow: hidden;
		}

		.name {
			font-size: 1.1rem;
			font-weight: 700;
			color: var(--text-main);
			margin-bottom: 4px;
		}

		.path {
			font-size: 0.8rem;
			color: var(--text-dim);
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}

		.type-badge {
			font-size: 0.75rem;
			padding: 6px 12px;
			border-radius: 8px;
			background: rgba(255, 255, 255, 0.05);
			color: var(--text-muted);
			font-weight: 600;
			text-transform: uppercase;
			letter-spacing: 0.5px;

			&.steam {
				background: rgba(14, 165, 233, 0.15);
				color: var(--accent-primary);
				border: 1px solid rgba(14, 165, 233, 0.2);
			}
		}
	}

	.loading {
		text-align: center;
		color: var(--text-dim);
		margin-top: 48px;
	}
</style>
