<script lang="ts">
	import playIcon from "../../../icons/play.svg";
	import settingsIcon from "../../../icons/settings.svg";
	import rocketIcon from "../../../icons/rocket.svg";

	export let game: any;
	export let icon: string = "";
	export let isRunning: boolean = false;
	export let isSelectionMode: boolean = false;
	export let isSelected: boolean = false;
	export let onLaunch: (game: any) => void = () => {};
	export let onConfigure: (game: any) => void = () => {};
	export let onSelect: (game: any) => void = () => {};

	function handleLaunch() {
		if (isSelectionMode) {
			onSelect(game);
			return;
		}
		onLaunch(game);
	}

	function handleConfigure() {
		onConfigure(game);
	}
</script>

<div
	class="list-card"
	class:running={isRunning}
	class:selection-mode={isSelectionMode}
	class:selected={isSelected}
	on:click={handleLaunch}
	role="button"
	tabindex="0"
	on:keydown={(e) => e.key === "Enter" && handleLaunch()}
>
	{#if isSelectionMode}
		<div class="selection-checkbox">
			<div class="checkbox" class:checked={isSelected}>
				{#if isSelected}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="16"
						height="16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="3"
						stroke-linecap="round"
						stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg
					>
				{/if}
			</div>
		</div>
	{/if}

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

		&.selected {
			border-color: var(--accent-primary);
			background: rgba(var(--accent-primary-rgb, 255, 255, 255), 0.05);

			.checkbox {
				background: var(--accent-primary) !important;
				border-color: var(--accent-primary) !important;
				color: #000;
			}
		}
	}

	.selection-checkbox {
		flex-shrink: 0;

		.checkbox {
			width: 24px;
			height: 24px;
			border: 2px solid rgba(255, 255, 255, 0.1);
			border-radius: 6px;
			background: rgba(0, 0, 0, 0.3);
			display: flex;
			align-items: center;
			justify-content: center;
			transition: all 0.2s;
			color: transparent;

			svg {
				stroke: currentColor;
			}
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
