<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import playIcon from "../../icons/play.svg";
	import settingsIcon from "../../icons/settings.svg";
	import rocketIcon from "../../icons/rocket.svg";

	export let game: any;
	export let icon: string = "";
	export let isRunning: boolean = false;

	const dispatch = createEventDispatcher();

	function handleLaunch() {
		dispatch("launch", game);
	}

	function handleConfigure() {
		dispatch("configure", game);
	}
</script>

<div class="list-card" class:running={isRunning} on:click={handleLaunch}>
	<div class="icon-section">
		{#if icon}
			<img src={icon} alt={game.name} class="game-icon" />
		{:else}
			<div class="fallback-wrapper">
				<img src={rocketIcon} alt="rocket" class="game-icon-fallback" />
			</div>
		{/if}

		{#if isRunning}
			<div class="running-indicator-small">
				<span class="pulse"></span>
			</div>
		{/if}
	</div>

	<div class="content-section">
		<div class="info">
			<span class="game-name">{game.name}</span>
			<span class="game-path">{game.path || game.config.LauncherPath}</span>
		</div>

		<div class="actions">
			<button class="action-btn play" title="Play Now">
				<img src={playIcon} alt="play" />
			</button>
			<button class="action-btn config" title="Configure" on:click|stopPropagation={handleConfigure}>
				<img src={settingsIcon} alt="configure" />
			</button>
		</div>
	</div>
</div>

<style lang="scss">
	.list-card {
		display: flex;
		align-items: center;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.05);
		border-radius: 16px;
		padding: 12px 20px;
		gap: 20px;
		cursor: pointer;
		transition: all 0.3s cubic-bezier(0.23, 1, 0.32, 1);
		max-width: 100%;

		&:hover {
			background: rgba(255, 255, 255, 0.07);
			border-color: rgba(255, 255, 255, 0.2);
			transform: translateX(8px);

			.game-icon {
				transform: scale(1.1);
			}

			.play img {
				transform: scale(1.2);
			}
		}

		&.running {
			border-color: var(--success, #22c55e);
			background: linear-gradient(90deg, rgba(34, 197, 94, 0.1) 0%, rgba(34, 197, 94, 0.02) 100%);
		}
	}

	.icon-section {
		height: 64px;
		aspect-ratio: 1/1;
		border-radius: 12px;
		overflow: hidden;
		position: relative;
		flex-shrink: 0;
		background: #000;
		border: 1px solid rgba(255, 255, 255, 0.1);

		.game-icon {
			width: 100%;
			height: 100%;
			object-fit: cover;
			transition: transform 0.4s;
		}

		.fallback-wrapper {
			width: 100%;
			height: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
			opacity: 1;
		}

		.game-icon-fallback {
			width: 32px;
			height: 32px;
			filter: brightness(0) invert(1);
		}

		.running-indicator-small {
			position: absolute;
			top: 4px;
			right: 4px;
			width: 8px;
			height: 8px;
			background: var(--success, #22c55e);
			border-radius: 50%;
			box-shadow: 0 0 8px var(--success, #22c55e);

			.pulse {
				position: absolute;
				inset: 0;
				background: inherit;
				border-radius: inherit;
				animation: ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite;
			}
		}
	}

	.content-section {
		flex: 1;
		display: flex;
		justify-content: space-between;
		align-items: center;
		min-width: 0;
	}

	.info {
		display: flex;
		flex-direction: column;
		gap: 2px;
		min-width: 0;

		.game-name {
			font-weight: 700;
			color: #fff;
			font-size: 1.1rem;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}

		.game-path {
			font-size: 0.8rem;
			color: rgba(255, 255, 255, 0.4);
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
			max-width: 100%;
		}
	}

	.actions {
		display: flex;
		gap: 8px;
	}

	.action-btn {
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		width: 40px;
		height: 40px;
		border-radius: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: all 0.2s;

		img {
			width: 20px;
			height: 20px;
			filter: brightness(0) invert(1);
			transition: transform 0.2s;
		}

		&:hover {
			background: rgba(255, 255, 255, 0.15);
			border-color: rgba(255, 255, 255, 0.3);
		}

		&.play:hover {
			background: var(--success, #22c55e);
			border-color: transparent;
		}
	}

	@keyframes ping {
		75%,
		100% {
			transform: scale(2.5);
			opacity: 0;
		}
	}
</style>
