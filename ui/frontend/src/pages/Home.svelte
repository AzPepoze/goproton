<script lang="ts">
	import {
		CleanupProcesses,
		GetAllGames,
		GetRunningSessions,
		KillSession,
		RunGame,
		GetSystemInfo,
		GetSystemUsage,
	} from "../../wailsjs/go/backend/App";
	import { onMount, onDestroy } from "svelte";
	import { notifications } from "../notificationStore";
	import { navigationCommand } from "../stores/navigationStore";
	import { runState } from "../stores/runState";
	import { loadExeIcon } from "../lib/iconService";
	import GameCard from "../components/GameCard.svelte";
	import trashIcon from "../icons/trash.svg";
	import type { core } from "../../wailsjs/go/models";

	let games = [];
	let sessions = [];
	let sessionInterval;
	let usageInterval;
	let gameIcons = {};
	let sysInfo: core.SystemInfo = {
		os: "Loading...",
		kernel: "",
		cpu: "",
		gpu: "",
		ram: "",
		driver: "",
	};
	let sysUsage: core.SystemUsage = {
		cpu: "0%",
		ram: "0%",
		gpu: "0%",
	};

	async function refreshData() {
		try {
			const fetchedGames = await GetAllGames();
			games = fetchedGames || [];
			const fetchedSessions = await GetRunningSessions();
			sessions = fetchedSessions || [];

			GetSystemInfo().then((info) => {
				sysInfo = info;
			});

			// Fetch icons for games
			for (const game of games) {
				const path = game.path || game.config.LauncherPath;
				if (path && !gameIcons[path]) {
					loadExeIcon(path).then((icon) => {
						if (icon) {
							gameIcons = { ...gameIcons, [path]: icon };
						}
					});
				}
			}
		} catch (err) {
			console.error("Failed to refresh home data:", err);
		}
	}

	onMount(() => {
		refreshData();
		sessionInterval = setInterval(async () => {
			try {
				const fetchedSessions = await GetRunningSessions();
				sessions = fetchedSessions || [];
			} catch (err) {
				console.error("Failed to fetch sessions in interval:", err);
			}
		}, 3000);

		usageInterval = setInterval(async () => {
			try {
				sysUsage = await GetSystemUsage();
			} catch (err) {
				console.error("Failed to fetch usage:", err);
			}
		}, 2000);
	});

	onDestroy(() => {
		if (sessionInterval) clearInterval(sessionInterval);
		if (usageInterval) clearInterval(usageInterval);
	});

	async function handleCleanup() {
		try {
			await CleanupProcesses();
			notifications.add("System cleaned successfully!", "success");
			refreshData();
		} catch (err) {
			notifications.add(`Cleanup failed: ${err}`, "error");
		}
	}

	async function handleQuickLaunch(game) {
		try {
			notifications.add(`Launching ${game.name}...`, "info");
			await RunGame(game.config, false); // No logs for quick launch
			refreshData();
		} catch (err) {
			notifications.add(`Launch failed: ${err}`, "error");
		}
	}

	function handleConfigure(game) {
		runState.set({
			options: game.config,
		});
		navigationCommand.set({ page: "run" });
	}

	function isGameRunning(game) {
		const path = game.path || game.config.LauncherPath;
		return sessions.some((s) => s.gamePath === path);
	}

	async function handleKillSession(pid, name) {
		try {
			await KillSession(pid);
			notifications.add(`Terminated session: ${name}`, "success");
			refreshData();
		} catch (err) {
			notifications.add(`Failed to kill session: ${err}`, "error");
		}
	}
</script>

<div class="home-container">
	<div class="hero-section">
		<h1 class="welcome-text">WELCOME TO <span class="highlight">GOPROTON!</span></h1>
	</div>

	<h2 class="section-title">System Status</h2>
	<div class="system-info-grid">
		<!-- CPU Card -->
		<div class="info-card glass">
			<div class="card-header-small">
				<svg class="card-icon-small" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="16" height="16" x="4" y="4" rx="2"/><rect width="6" height="6" x="9" y="9"/><path d="M15 2v2"/><path d="M15 20v2"/><path d="M2 15h2"/><path d="M2 9h2"/><path d="M20 15h2"/><path d="M20 9h2"/><path d="M9 2v2"/><path d="M9 20v2"/></svg>
				<span class="label">Processor</span>
			</div>
			<div class="card-content">
				<span class="value" title={sysInfo.cpu}>{sysInfo.cpu}</span>
				<div class="usage-mini">
					<span class="u-val">{sysUsage.cpu} Usage</span>
					<div class="progress-bar-mini">
						<div class="fill" style="width: {sysUsage.cpu}"></div>
					</div>
				</div>
			</div>
		</div>

		<!-- RAM Card -->
		<div class="info-card glass">
			<div class="card-header-small">
				<svg class="card-icon-small" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 19v-3"/><path d="M10 19v-3"/><path d="M14 19v-3"/><path d="M18 19v-3"/><path d="M8 11V9"/><path d="M16 11V9"/><rect width="18" height="12" x="3" y="7" rx="2"/><path d="M3 13h18"/></svg>
				<span class="label">Memory</span>
			</div>
			<div class="card-content">
				<span class="value">{sysInfo.ram} Total</span>
				<div class="usage-mini">
					<span class="u-val">{sysUsage.ram.split(" / ")[0]} Used</span>
					<div class="progress-bar-mini">
						<div
							class="fill"
							style="width: {sysUsage.ram.includes('(')
								? sysUsage.ram.split('(').pop().replace(')', '')
								: '0%'}"
						></div>
					</div>
				</div>
			</div>
		</div>

		<!-- GPU Card -->
		<div class="info-card glass">
			<div class="card-header-small">
				<svg class="card-icon-small" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8V6a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h2"/><rect width="12" height="8" x="10" y="10" rx="2"/><path d="M14 10V8"/><path d="M18 10V8"/><path d="M14 20v-2"/><path d="M18 20v-2"/><path d="M22 14h-2"/><path d="M22 18h-2"/></svg>
				<span class="label">Graphics</span>
			</div>
			<div class="card-content">
				<span class="value" title={sysInfo.gpu}>{sysInfo.gpu}</span>
				<div class="usage-mini">
					<span class="u-val">{sysUsage.gpu} Load</span>
					<div class="progress-bar-mini">
						<div class="fill" style="width: {sysUsage.gpu}; background: var(--accent-secondary, #b197fc)"></div>
					</div>
				</div>
				<span class="sub-value" style="margin-top: -4px" title={sysInfo.driver}>{sysInfo.driver}</span>
			</div>
		</div>

		<!-- OS Card -->
		<div class="info-card glass">
			<div class="card-header-small">
				<svg class="card-icon-small" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="14" x="2" y="3" rx="2"/><line x1="8" x2="16" y1="21" y2="21"/><line x1="12" x2="12" y1="17" y2="21"/></svg>
				<span class="label">System</span>
			</div>
			<div class="card-content single">
				<span class="value" title={sysInfo.os}>{sysInfo.os}</span>
				<span class="sub-value">Kernel: {sysInfo.kernel}</span>
			</div>
		</div>
	</div>

	{#if sessions.length > 0}
		<div class="sessions-section">
			<h2 class="section-title">Running Sessions</h2>
			<div class="sessions-grid">
				{#each sessions as session}
					<div class="session-card">
						<div class="session-info">
							<div class="session-title">{session.gameName}</div>
							<div class="session-pid">PID: {session.pid}</div>
						</div>
						<button
							class="kill-btn"
							on:click={() => handleKillSession(session.pid, session.gameName)}
						>
							Terminate
						</button>
					</div>
				{/each}
			</div>
		</div>
	{/if}

	<div class="quick-launch-section">
		<h2 class="section-title">Quick Launch</h2>
		{#if games.length === 0}
			<div class="empty-state">
				<p>
					No games configured yet. Go to <span
						class="link"
						on:click={() => navigationCommand.set({ page: "run" })}>Run</span
					> to add one.
				</p>
			</div>
		{:else}
			<div class="games-grid">
				{#each games as game}
					<GameCard
						{game}
						icon={gameIcons[game.path || game.config.LauncherPath]}
						isRunning={isGameRunning(game)}
						on:launch={() => handleQuickLaunch(game)}
						on:configure={() => handleConfigure(game)}
					/>
				{/each}
			</div>
		{/if}
	</div>

	<div class="utils-section">
		<h2 class="section-title">Utilities</h2>
		<div class="grid">
			<button class="card hoverable cleanup-card" on:click={handleCleanup}>
				<div class="card-header">
					<span class="card-icon text-warning">
						<img src={trashIcon} alt="cleanup" class="svg-icon" />
					</span>
					<h3>Cleanup System</h3>
				</div>
				<p>Forcefully terminate all running game instances and background services.</p>
			</button>
		</div>
	</div>
</div>

<style lang="scss">
	.home-container {
		display: flex;
		flex-direction: column;
		height: 100%;
		padding: 40px;
		overflow-y: auto;
		background-color: transparent;
		gap: 40px;
	}

	.hero-section {
		text-align: center;
		margin-bottom: 20px;

		.welcome-text {
			font-size: 3rem;
			font-weight: 900;
			color: #fff;
			margin: 0;
			letter-spacing: -1px;
			line-height: 1.1;

			.highlight {
				color: transparent;
				-webkit-text-stroke: 1px rgba(255, 255, 255, 0.8);
				background: linear-gradient(180deg, #fff 0%, rgba(255, 255, 255, 0.2) 100%);
				-webkit-background-clip: text;
				background-clip: text;
			}
		}
	}

	.system-info-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
		gap: 20px;
		margin-bottom: 20px;
	}

	.info-card {
		padding: 20px;
		border-radius: 20px;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid var(--glass-border);
		display: flex;
		flex-direction: column;
		gap: 12px;
		transition: all 0.3s cubic-bezier(0.23, 1, 0.32, 1);
		min-width: 0;

		&:hover {
			background: rgba(255, 255, 255, 0.06);
			border-color: var(--glass-border-bright);
			transform: translateY(-2px);
		}

		.card-header-small {
			display: flex;
			align-items: center;
			gap: 10px;

			.card-icon-small {
				width: 18px;
				height: 18px;
				color: var(--text-dim);
				opacity: 0.8;
			}

			.label {
				font-size: 0.7rem;
				font-weight: 800;
				color: var(--text-dim);
				text-transform: uppercase;
				letter-spacing: 1.5px;
			}
		}

		.card-content {
			display: flex;
			flex-direction: column;
			gap: 10px;

			&.single {
				gap: 4px;
			}
		}

		.value {
			font-size: 1rem;
			font-weight: 700;
			color: var(--text-main);
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}

		.sub-value {
			font-size: 0.8rem;
			color: var(--text-muted);
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}
	}

	.usage-mini {
		display: flex;
		flex-direction: column;
		gap: 6px;

		.u-val {
			font-size: 0.85rem;
			font-weight: 700;
			color: var(--accent-primary);
		}
	}

	.progress-bar-mini {
		height: 6px;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 10px;
		overflow: hidden;

		.fill {
			height: 100%;
			background: var(--accent-primary);
			transition: width 0.3s ease;
		}
	}

	.section-title {
		font-size: 1.2rem;
		font-weight: 800;
		color: rgba(255, 255, 255, 0.4);
		text-transform: uppercase;
		letter-spacing: 2px;
		margin-bottom: 20px;
	}

	.sessions-section {
		display: flex;
		flex-direction: column;
		gap: 20px;
		background: linear-gradient(135deg, rgba(239, 68, 68, 0.1) 0%, rgba(239, 68, 68, 0.02) 100%);
		padding: 24px;
		border-radius: 24px;
		border: 1px solid rgba(239, 68, 68, 0.2);
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
		animation: slide-down 0.5s cubic-bezier(0.23, 1, 0.32, 1);

		.section-title {
			margin-bottom: 0;
			color: var(--danger);
		}
	}

	.sessions-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 16px;
	}

	.session-card {
		background: rgba(0, 0, 0, 0.3);
		border: 1px solid rgba(255, 255, 255, 0.05);
		border-radius: 16px;
		padding: 14px 20px;
		display: flex;
		justify-content: space-between;
		align-items: center;
		transition: all 0.3s;

		&:hover {
			border-color: rgba(239, 68, 68, 0.4);
			background: rgba(0, 0, 0, 0.5);
			transform: translateX(4px);
		}

		.session-info {
			display: flex;
			flex-direction: column;
			gap: 2px;
		}

		.session-title {
			font-weight: 800;
			color: #fff;
			font-size: 1rem;
			letter-spacing: -0.3px;
		}

		.session-pid {
			font-size: 0.7rem;
			color: rgba(255, 255, 255, 0.4);
			font-family: monospace;
			font-weight: 600;
		}

		.kill-btn {
			background: var(--danger);
			color: #fff;
			padding: 8px 16px;
			border: none;
			border-radius: 10px;
			font-size: 0.75rem;
			font-weight: 800;
			cursor: pointer;
			transition: all 0.2s;
			box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);

			&:hover {
				filter: brightness(1.2);
				transform: translateY(-2px);
				box-shadow: 0 6px 16px rgba(239, 68, 68, 0.4);
			}

			&:active {
				transform: translateY(0);
			}
		}
	}

	@keyframes slide-down {
		from {
			opacity: 0;
			transform: translateY(-20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.games-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
		gap: 32px;
	}

	.empty-state {
		background: rgba(255, 255, 255, 0.02);
		border: 1px dashed var(--glass-border);
		border-radius: 12px;
		padding: 32px;
		text-align: center;
		color: var(--text-muted);

		.link {
			color: var(--accent-color);
			text-decoration: underline;
			cursor: pointer;
		}
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
		gap: 32px;
	}

	.card {
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid var(--glass-border);
		border-radius: 20px;
		padding: 24px;
		transition: all 0.3s;
		text-align: left;

		.card-header {
			display: flex;
			align-items: center;
			gap: 16px;
			margin-bottom: 12px;
		}

		.card-icon {
			width: 24px;
			height: 24px;
			display: flex;
			align-items: center;

			.svg-icon {
				width: 100%;
				height: 100%;
				filter: invert(72%) sepia(85%) saturate(1008%) hue-rotate(359deg) brightness(101%) contrast(93%);
			}
		}

		h3 {
			font-size: 1.1rem;
			font-weight: 700;
			margin: 0;
			color: #fff;
		}

		p {
			font-size: 0.9rem;
			color: var(--text-muted);
			margin: 0;
			line-height: 1.4;
		}

		&.hoverable {
			cursor: pointer;

			&:hover {
				background: rgba(255, 255, 255, 0.05);
				border-color: var(--glass-border-bright);
				transform: translateY(-4px);
			}
		}

		&.cleanup-card:hover {
			border-color: rgba(245, 158, 11, 0.4);
			background: rgba(245, 158, 11, 0.05);
		}
	}

	.text-warning {
		color: #f59e0b;
	}
</style>
