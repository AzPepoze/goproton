<script lang="ts">
	import { fly } from "svelte/transition";
	import { KillSession } from "../../bindings/goproton-wails/backend/app";
	import { notifications } from "../notificationStore";
	import { createEventDispatcher } from "svelte";

	export let sessions = [];
	let isExpanded = false;

	const dispatch = createEventDispatcher();

	async function handleKillSession(pid, name) {
		try {
			await KillSession(pid);
			notifications.add(`Terminated session: ${name}`, "success");
			dispatch("refresh");
		} catch (err) {
			notifications.add(`Failed to kill session: ${err}`, "error");
		}
	}
</script>

<div class="sessions-drawer-wrapper" class:expanded={isExpanded}>
	<button class="toggle-btn" on:click={() => (isExpanded = !isExpanded)}>
		<div class="indicator" class:active={sessions.length > 0}></div>
		<span class="trigger-text">SESSIONS</span>
		{#if sessions.length > 0}
			<span class="count">{sessions.length}</span>
		{/if}
	</button>

	<div class="drawer-content">
		<h2 class="section-title">Running Sessions</h2>

		{#if sessions.length === 0}
			<div class="empty-sessions">
				<p>No active sessions</p>
			</div>
		{:else}
			<div class="sessions-list">
				{#each sessions as session}
					<div class="session-card" in:fly={{ x: 50, duration: 400 }}>
						<div class="session-info">
							<div class="session-title">{session.gameName}</div>
							<div class="session-pid">PID: {session.pid}</div>
						</div>
						<button
							class="kill-btn"
							on:click={() => handleKillSession(session.pid, session.gameName)}
						>
							Kill
						</button>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style lang="scss">
	.sessions-drawer-wrapper {
		position: fixed;
		top: 20px;
		right: 20px;
		bottom: 100px; /* Leave space for bottom drawer trigger */
		width: 300px;
		background: rgba(18, 18, 22, 0.98);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 20px;
		transform: translateX(calc(100% - 48px));
		transition: all 0.5s cubic-bezier(0.23, 1, 0.32, 1);
		z-index: 101;
		display: flex;
		box-shadow: -10px 0 40px rgba(0, 0, 0, 0.4);
		overflow: hidden;

		&.expanded {
			transform: translateX(0);
		}
	}

	.toggle-btn {
		width: 48px;
		height: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		background: rgba(60, 60, 65, 0.2);
		border: none;
		border-right: 1px solid rgba(255, 255, 255, 0.05);
		transition: all 0.2s ease;
		gap: 20px;
		padding: 20px 0;

		.trigger-text {
			writing-mode: vertical-rl;
			font-size: 0.7rem;
			font-weight: 900;
			color: rgba(255, 255, 255, 0.4);
			letter-spacing: 2px;
			text-transform: uppercase;
		}

		.indicator {
			width: 8px;
			height: 8px;
			border-radius: 50%;
			background: rgba(255, 255, 255, 0.1);

			&.active {
				background: #ef4444;
				box-shadow: 0 0 10px rgba(239, 68, 68, 0.5);
				animation: pulse 2s infinite;
			}
		}

		.count {
			background: #ef4444;
			color: white;
			font-size: 0.65rem;
			font-weight: 900;
			padding: 4px;
			border-radius: 6px;
			min-width: 15px;
			text-align: center;
		}

		&:hover {
			background: rgba(80, 80, 85, 0.4);
			.trigger-text {
				color: #fff;
			}
		}
	}

	.drawer-content {
		flex: 1;
		padding: 24px;
		display: flex;
		flex-direction: column;
		gap: 20px;
		overflow-y: auto;
		min-width: 252px;
	}

	.section-title {
		font-size: 0.9rem;
		font-weight: 800;
		color: #ef4444;
		text-transform: uppercase;
		letter-spacing: 1.5px;
		margin: 0;
	}

	.sessions-list {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.session-card {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.05);
		border-radius: 12px;
		padding: 12px;
		display: flex;
		justify-content: space-between;
		align-items: center;
		transition: all 0.3s;

		&:hover {
			background: rgba(255, 255, 255, 0.05);
			border-color: rgba(239, 68, 68, 0.3);
		}

		.session-info {
			display: flex;
			flex-direction: column;
			gap: 2px;
			min-width: 0;
		}

		.session-title {
			font-weight: 700;
			color: #fff;
			font-size: 0.85rem;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}

		.session-pid {
			font-size: 0.65rem;
			color: rgba(255, 255, 255, 0.4);
			font-family: monospace;
		}

		.kill-btn {
			background: rgba(239, 68, 68, 0.2);
			color: #ef4444;
			padding: 4px 10px;
			border: 1px solid rgba(239, 68, 68, 0.3);
			border-radius: 8px;
			font-size: 0.7rem;
			font-weight: 800;
			cursor: pointer;
			transition: all 0.2s;

			&:hover {
				background: #ef4444;
				color: white;
			}
		}
	}

	.empty-sessions {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		color: rgba(255, 255, 255, 0.2);
		font-size: 0.8rem;
		font-style: italic;
	}

	@keyframes pulse {
		0% {
			opacity: 1;
		}
		50% {
			opacity: 0.4;
		}
		100% {
			opacity: 1;
		}
	}
</style>
