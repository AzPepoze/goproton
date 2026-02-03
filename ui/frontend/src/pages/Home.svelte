<script lang="ts">
	import { GetAllGames, GetRunningSessions, KillSession, RunGame } from "../../wailsjs/go/backend/App";
	import { onMount, onDestroy } from "svelte";
	import { fly } from "svelte/transition";
	import { notifications } from "../notificationStore";
	import { navigationCommand } from "../stores/navigationStore";
	import { runState } from "../stores/runState";
	import { loadExeIcon } from "../lib/iconService";
	import GameCard from "../components/GameCard.svelte";
	import StatusDrawer from "../components/StatusDrawer.svelte";
	import Modal from "../components/Modal.svelte";

	let games = [];
	let sessions = [];
	let sessionInterval;
	let gameIcons = {};
	let showHelpModal = false;
	let currentView: "grid" | "list-grid" = "grid";
	let searchQuery = "";

	$: filteredGames = games.filter((game) => game.name.toLowerCase().includes(searchQuery.toLowerCase()));

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
		runState.update((s) => ({
			...s,
			options: game.config,
		}));
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
	{#if sessions.length > 0}
		<div class="sessions-section">
			<h2 class="section-title">Running Sessions</h2>
			<div class="sessions-grid">
				{#each sessions as session}
					<div class="session-card" in:fly={{ y: -20, duration: 400 }}>
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
		<div class="section-header">
			<h2 class="section-title">Quick Launch</h2>
			<button class="help-btn" on:click={() => (showHelpModal = true)} title="How it works">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2.5"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="12" cy="12" r="10" />
					<path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
					<line x1="12" y1="17" x2="12.01" y2="17" />
				</svg>
			</button>

			<div class="search-container">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="16"
					height="16"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="search-icon"
				>
					<circle cx="11" cy="11" r="8"></circle>
					<line x1="21" y1="21" x2="16.65" y2="16.65"></line>
				</svg>
				<input type="text" placeholder="Search games..." bind:value={searchQuery} class="search-input" />
				{#if searchQuery}
					<button class="clear-search" on:click={() => (searchQuery = "")}>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="14"
							height="14"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<line x1="18" y1="6" x2="6" y2="18"></line>
							<line x1="6" y1="6" x2="18" y2="18"></line>
						</svg>
					</button>
				{/if}
			</div>

			<div class="view-switcher">
				<button
					class="view-btn"
					class:active={currentView === "grid"}
					on:click={() => (currentView = "grid")}
					title="Grid View"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="18"
						height="18"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"
						></rect><rect x="14" y="14" width="7" height="7"></rect><rect
							x="3"
							y="14"
							width="7"
							height="7"
						></rect></svg
					>
				</button>
				<button
					class="view-btn"
					class:active={currentView === "list-grid"}
					on:click={() => (currentView = "list-grid")}
					title="List View"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="18"
						height="18"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						><line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"
						></line><line x1="8" y1="18" x2="21" y2="18"></line><line x1="3" y1="6" x2="3.01" y2="6"
						></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line
							x1="3"
							y1="18"
							x2="3.01"
							y2="18"
						></line></svg
					>
				</button>
			</div>
		</div>

		{#if games.length === 0}
			<div class="empty-state">
				<p>
					No games configured yet. Go to <button
						class="link-btn"
						on:click={() => navigationCommand.set({ page: "run" })}>Run</button
					> to add one.
				</p>
			</div>
		{:else}
			<div
				class="games-container"
				class:grid-view={currentView === "grid"}
				class:list-view={currentView === "list-grid"}
			>
				{#if filteredGames.length === 0 && games.length > 0}
					<div class="no-results">
						<p>No games matching "{searchQuery}"</p>
						<button class="link-btn" on:click={() => (searchQuery = "")}>Clear search</button>
					</div>
				{:else}
					<div class="games-grid">
						{#each filteredGames as game}
							<GameCard
								{game}
								icon={gameIcons[game.path || game.config.LauncherPath]}
								isRunning={isGameRunning(game)}
								view={currentView}
								on:launch={() => handleQuickLaunch(game)}
								on:configure={() => handleConfigure(game)}
							/>
						{/each}
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>

<Modal show={showHelpModal} title="How it works" onClose={() => (showHelpModal = false)}>
	<div class="help-content">
		<section>
			<h3>Adding Games</h3>
			<p>
				Go to the <strong>Run</strong> page, select your game executable (and launcher if applicable),
				configure your settings, and click <strong>LAUNCH GAME</strong>.
			</p>
			<p>
				After the first run, the game will automatically appear here in <strong>Quick Launch</strong>.
			</p>
		</section>

		<section>
			<h3>Quick Launch</h3>
			<p>Click on any game card in this section to start it immediately with its saved configuration.</p>
		</section>

		<section>
			<h3>Managing Sessions</h3>
			<p>Active game sessions are displayed at the top. You can terminate them if they become unresponsive.</p>
		</section>

		<section>
			<h3>CLI Usage</h3>
			<p>You can also launch games directly from your terminal or add them to your desktop entries:</p>
			<code class="help-code">goproton /path/to/game.exe</code>
		</section>
	</div>
</Modal>

<StatusDrawer />

<style lang="scss">
	.home-container {
		display: flex;
		flex-direction: column;
		height: 100%;
		width: 100%;
		padding: 0;
		background-color: transparent;
		gap: 20px;
		box-sizing: border-box;
		min-height: 0;
		overflow-x: hidden;
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
		flex-shrink: 0;
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
			margin-bottom: 10px;
			color: #ef4444;
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
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
			max-width: 200px;
		}

		.session-pid {
			font-size: 0.7rem;
			color: rgba(255, 255, 255, 0.4);
			font-family: monospace;
			font-weight: 600;
		}

		.kill-btn {
			background: #ef4444;
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

	.quick-launch-section {
		display: flex;
		flex-direction: column;
		flex: 1;
		min-height: 0;

		.section-header {
			display: flex;
			align-items: center;
			gap: 12px;
			margin-bottom: 20px;
			flex-wrap: wrap;

			.section-title {
				margin: 0;
				line-height: 1;
				white-space: nowrap;
			}
		}

		.view-switcher {
			display: flex;
			background: rgba(255, 255, 255, 0.05);
			padding: 4px;
			border-radius: 12px;
			gap: 4px;
			border: 1px solid rgba(255, 255, 255, 0.05);

			.view-btn {
				background: none;
				border: none;
				color: rgba(255, 255, 255, 0.4);
				padding: 6px;
				cursor: pointer;
				border-radius: 8px;
				display: flex;
				align-items: center;
				justify-content: center;
				aspect-ratio: 1 / 1;
				transition: all 0.2s;

				&:hover {
					color: #fff;
					background: rgba(255, 255, 255, 0.05);
				}

				&.active {
					color: #fff;
					background: rgba(255, 255, 255, 0.1);
					box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
				}
			}
		}

		.help-btn {
			background: none;
			border: none;
			color: rgba(255, 255, 255, 0.2);
			cursor: pointer;
			padding: 4px;
			display: flex;
			align-items: center;
			justify-content: center;
			transition: all 0.2s;
			border-radius: 50%;

			&:hover {
				color: var(--accent-primary, #fff);
				background: rgba(255, 255, 255, 0.05);
				transform: scale(1.1);
			}

			svg {
				width: 20px;
				height: 20px;
			}
		}

		.search-container {
			display: flex;
			align-items: center;
			background: rgba(255, 255, 255, 0.05);
			border: 1px solid rgba(255, 255, 255, 0.05);
			border-radius: 12px;
			padding: 4px 10px;
			gap: 8px;
			flex: 1;
			transition: all 0.3s;

			&:focus-within {
				background: rgba(255, 255, 255, 0.1);
				border-color: rgba(255, 255, 255, 0.2);
				box-shadow: 0 0 15px rgba(0, 0, 0, 0.2);
			}

			.search-icon {
				color: rgba(255, 255, 255, 0.3);
			}

			.search-input {
				background: none;
				border: none;
				color: #fff;
				font-size: 0.9rem;
				width: 100%;
				outline: none;

				&::placeholder {
					color: rgba(255, 255, 255, 0.2);
				}
			}

			.clear-search {
				background: none;
				border: none;
				color: rgba(255, 255, 255, 0.3);
				cursor: pointer;
				padding: 2px;
				display: flex;
				align-items: center;
				justify-content: center;
				border-radius: 4px;

				&:hover {
					color: #fff;
					background: rgba(255, 255, 255, 0.1);
				}
			}
		}
	}

	.help-content {
		display: flex;
		flex-direction: column;
		gap: 24px;
		color: var(--text-main, #eee);

		section {
			h3 {
				margin: 0 0 8px 0;
				font-size: 1.1rem;
				color: var(--accent-primary, #fff);
			}

			p {
				margin: 0;
				line-height: 1.6;
				font-size: 0.95rem;
				color: var(--text-dim, #aaa);

				strong {
					color: var(--text-main, #eee);
				}
			}

			.help-code {
				display: block;
				background: rgba(0, 0, 0, 0.3);
				padding: 12px;
				border-radius: 8px;
				font-family: monospace;
				font-size: 0.85rem;
				color: var(--accent-primary, #fff);
				margin-top: 10px;
				border: 1px solid rgba(255, 255, 255, 0.05);
			}

			& + section {
				padding-top: 16px;
				border-top: 1px solid rgba(255, 255, 255, 0.05);
			}
		}
	}

	.games-container {
		flex: 1;
		min-height: 0;
		display: flex;
		flex-direction: column;

		&.list-view {
			.games-grid {
				grid-template-columns: 1fr;
				gap: 16px;
			}
		}
	}

	.no-results {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 200px;
		color: rgba(255, 255, 255, 0.4);
		gap: 10px;

		p {
			margin: 0;
			font-size: 1rem;
		}

		.link-btn {
			background: none;
			border: none;
			color: var(--accent-color, #60a5fa);
			text-decoration: underline;
			cursor: pointer;
			padding: 0;
			font: inherit;
			font-weight: 600;

			&:hover {
				filter: brightness(1.2);
			}
		}
	}

	.games-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
		grid-auto-rows: min-content;
		gap: 32px;
		overflow-y: auto;
		overflow-x: hidden;
		padding: 10px;
		flex: 1;
		min-height: 0;
	}

	.empty-state {
		background: rgba(255, 255, 255, 0.02);
		border: 1px dashed var(--glass-border);
		border-radius: 12px;
		padding: 32px;
		text-align: center;
		color: var(--text-muted);

		.link-btn {
			background: none;
			border: none;
			color: var(--accent-color, #60a5fa);
			text-decoration: underline;
			cursor: pointer;
			padding: 0;
			font: inherit;
			font-weight: 600;

			&:hover {
				filter: brightness(1.2);
			}
		}
	}
</style>
