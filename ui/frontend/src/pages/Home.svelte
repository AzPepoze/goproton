<script lang="ts">
	import {
		CleanupProcesses,
		GetAllGames,
		GetRunningSessions,
		KillSession,
		RunGame,
	} from "../../wailsjs/go/main/App";
	import { onMount, onDestroy } from "svelte";
	import { notifications } from "../notificationStore";
	import { navigationCommand } from "../stores/navigationStore";
	import { runState } from "../stores/runState";
	import { loadExeIcon } from "../lib/iconService";
	import GameCard from "../components/GameCard.svelte";
	import trashIcon from "../icons/trash.svg";

	let games = [];
	let sessions = [];
	let sessionInterval;
	let gameIcons = {};

	async function refreshData() {
		try {
			const fetchedGames = await GetAllGames();
			games = fetchedGames || [];
			const fetchedSessions = await GetRunningSessions();
			sessions = fetchedSessions || [];

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
	});

	onDestroy(() => {
		if (sessionInterval) clearInterval(sessionInterval);
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
		grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
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
