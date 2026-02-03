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

<div class="game-card" class:running={isRunning}>
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="game-icon-container" on:click={handleLaunch}>
		{#if isRunning}
			<div class="running-indicator">
				<span class="pulse"></span>
				<span class="indicator-text">RUNNING</span>
			</div>
		{/if}

		<div class="rainbow-glow"></div>

		<div class="icon-wrapper">
			{#if icon}
				<img src={icon} alt={game.name} class="game-icon" />
			{:else}
				<img src={rocketIcon} alt="rocket" class="game-icon-fallback" />
			{/if}
		</div>

		<div class="play-overlay">
			<img src={playIcon} alt="play" class="launch-icon-large" />
		</div>
	</div>

	<div class="game-footer">
		<span class="game-name" title={game.name}>{game.name}</span>
		<button class="config-btn" title="Configure" on:click|stopPropagation={handleConfigure}>
			<img src={settingsIcon} alt="configure" />
		</button>
	</div>
</div>

<style lang="scss">
	.game-card {
		display: flex;
		flex-direction: column;
		gap: 16px;
		transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
		position: relative;
		max-width: 180px;
		width: 100%;
		margin: 0 auto;

		&:hover {
			transform: scale(1.05);

			.rainbow-glow {
				opacity: 0.8;
				animation: rainbow-glow-animation 3s linear infinite;
			}

			.game-icon-container {
				border-color: rgba(255, 255, 255, 0.4);
				box-shadow: 0 12px 32px rgba(0, 0, 0, 0.5);

				.play-overlay {
					opacity: 1;
					backdrop-filter: blur(4px);
				}

				img.game-icon {
					transform: scale(1.1);
				}
			}

			.game-footer .game-name {
				color: #fff;
				text-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
			}
		}

		&.running {
			.rainbow-glow {
				opacity: 1;
				animation: rainbow-glow-animation 2s linear infinite;
			}

			.game-icon-container {
				border-color: var(--success, #22c55e);
				box-shadow: 0 0 20px rgba(34, 197, 94, 0.3);
			}
		}
	}

	.game-icon-container {
		aspect-ratio: 1;
		background: #000;
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 28px;
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
		cursor: pointer;
		transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
		z-index: 1;
	}

	.icon-wrapper {
		position: absolute;
		inset: 2px;
		background: #111;
		border-radius: 26px;
		z-index: 2;
		overflow: hidden;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.rainbow-glow {
		position: absolute;
		inset: -2px;
		background: conic-gradient(from 0deg, #ff0000, #ffff00, #00ff00, #00ffff, #0000ff, #ff00ff, #ff0000);
		opacity: 0;
		transition: opacity 0.4s;
		z-index: 1;
		filter: blur(3px) brightness(1.5);
		border-radius: 32px;
	}

	.game-icon {
		width: 100%;
		height: 100%;
		object-fit: cover;
		transition: transform 0.6s cubic-bezier(0.23, 1, 0.32, 1);
	}

	.game-icon-fallback {
		width: 64px;
		height: 64px;
		opacity: 0.8;
		filter: brightness(0) invert(1) drop-shadow(0 0 10px rgba(255, 255, 255, 0.4));
	}

	.play-overlay {
		position: absolute;
		inset: 0;
		background: transparent;
		display: flex;
		align-items: center;
		justify-content: center;
		opacity: 0;
		transition: opacity 0.3s;
		z-index: 3;

		.launch-icon-large {
			width: 60px;
			height: 60px;
			filter: brightness(0) invert(1) drop-shadow(0 0 15px rgba(0, 0, 0, 0.9));
			transform: scale(0.8);
			transition: transform 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
		}
	}

	.game-card:hover .launch-icon-large {
		transform: scale(1);
	}

	.running-indicator {
		position: absolute;
		top: 12px;
		right: 12px;
		background: var(--success, #22c55e);
		color: #fff;
		padding: 4px 10px;
		border-radius: 20px;
		font-size: 0.65rem;
		font-weight: 900;
		display: flex;
		align-items: center;
		gap: 6px;
		letter-spacing: 0.5px;
		z-index: 10;
		box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3);

		.pulse {
			width: 6px;
			height: 6px;
			background: #fff;
			border-radius: 50%;
			display: inline-block;
			animation: blink 1s infinite;
		}
	}

	.game-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0 8px;

		.game-name {
			font-size: 0.95rem;
			font-weight: 700;
			color: rgba(255, 255, 255, 0.9);
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
			flex: 1;
			transition: all 0.3s;
			letter-spacing: -0.2px;
		}

		.config-btn {
			background: rgba(255, 255, 255, 0.05);
			border: 1px solid rgba(255, 255, 255, 0.1);
			padding: 6px;
			border-radius: 10px;
			cursor: pointer;
			display: flex;
			align-items: center;
			justify-content: center;
			transition: all 0.3s;

			&:hover {
				background: rgba(255, 255, 255, 0.15);
				transform: rotate(45deg);
				border-color: rgba(255, 255, 255, 0.3);
			}

			img {
				width: 16px;
				height: 16px;
				filter: brightness(0) invert(1);
			}
		}
	}

	@keyframes rainbow-glow-animation {
		0% {
			filter: blur(3px) hue-rotate(0deg);
		}
		50% {
			filter: blur(3px) hue-rotate(180deg);
			transform: scale(1.02);
		}
		100% {
			filter: blur(3px) hue-rotate(360deg);
		}
	}

	@keyframes blink {
		0%,
		100% {
			opacity: 1;
			transform: scale(1);
		}
		50% {
			opacity: 0.5;
			transform: scale(0.8);
		}
	}
</style>
