<script lang="ts">
	import rocketIcon from "../../../icons/rocket.svg";

	export let game: any;
	export let icon: string = "";
	export let isRunning: boolean = false;
	export let active: boolean = false;
	export let onLaunch: (game: any) => void = () => {};
	export let onConfigure: (game: any) => void = () => {};

	function handleLaunch() {
		onLaunch(game);
	}

	function handleConfigure() {
		onConfigure(game);
	}
</script>
<div class="perspective-card" class:active class:running={isRunning}>
	<div
		class="card-inner"
		on:click={handleLaunch}
		role="button"
		tabindex="0"
		on:keydown={(e) => e.key === "Enter" && handleLaunch()}
	>
		<div class="glow"></div>

		<div class="image-container">
			{#if icon}
				<img src={icon} alt={game.name} class="game-icon" />
			{:else}
				<div class="fallback">
					<img src={rocketIcon} alt="rocket" />
				</div>
			{/if}
		</div>

		<div class="info-overlay">
			<div class="name-row">
				<span class="name">{game.name}</span>
				<button class="config-btn-small" title="Configure" on:click|stopPropagation={handleConfigure}>
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
						><circle cx="12" cy="12" r="3"></circle><path
							d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
						></path></svg
					>
				</button>
			</div>
			{#if isRunning}
				<div class="status">RUNNING</div>
			{/if}
		</div>
	</div>
</div>

<style lang="scss">
	.perspective-card {
		width: 220px;
		height: 310px;
		cursor: pointer;
		transition: all 0.6s cubic-bezier(0.23, 1, 0.32, 1);
		position: relative;
		perspective: 1000px;
		flex-shrink: 0;
		transform: rotateY(-15deg);

		&:hover {
			transform: scale(1.1) translateY(-10px) rotateY(0deg);
			z-index: 10;

			.card-inner {
				box-shadow:
					0 30px 60px rgba(0, 0, 0, 0.8),
					0 0 20px rgba(255, 255, 255, 0.1);
				border-color: rgba(255, 255, 255, 0.4);
			}

			.glow {
				opacity: 0.8;
			}

			.info-overlay {
				opacity: 1;
				transform: translateY(0);
			}
		}

		&.running {
			.card-inner {
				border-color: var(--success, #22c55e);
				box-shadow: 0 0 30px rgba(34, 197, 94, 0.2);
			}
			.status {
				color: var(--success, #22c55e);
			}
		}
	}

	.card-inner {
		width: 100%;
		height: 100%;
		background: #111;
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 20px;
		overflow: hidden;
		position: relative;
		transition: all 0.4s;
		box-shadow: 0 10px 20px rgba(0, 0, 0, 0.4);
	}

	.glow {
		position: absolute;
		inset: 0;
		background: radial-gradient(circle at 50% 0%, rgba(255, 255, 255, 0.2), transparent);
		opacity: 0.2;
		transition: opacity 0.4s;
		pointer-events: none;
	}

	.image-container {
		width: 100%;
		height: 100%;

		.game-icon {
			width: 100%;
			height: 100%;
			object-fit: cover;
		}

		.fallback {
			width: 100%;
			height: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
			background: linear-gradient(135deg, #1a1a1a 0%, #000 100%);

			img {
				width: 48px;
				height: 48px;
				opacity: 0.1;
			}
		}
	}

	.info-overlay {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		padding: 20px;
		background: linear-gradient(transparent, rgba(0, 0, 0, 0.9));
		display: flex;
		flex-direction: column;
		gap: 4px;
		opacity: 0.8;
		transform: translateY(5px);
		transition: all 0.4s;

		.name {
			font-weight: 800;
			color: #fff;
			font-size: 1rem;
			text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}

		.name-row {
			display: flex;
			justify-content: space-between;
			align-items: center;
			gap: 8px;
		}

		.config-btn-small {
			background: rgba(255, 255, 255, 0.1);
			border: none;
			color: #fff;
			padding: 4px;
			border-radius: 6px;
			cursor: pointer;
			display: flex;
			transition: all 0.2s;

			&:hover {
				background: rgba(255, 255, 255, 0.2);
				transform: rotate(45deg);
			}
		}

		.status {
			font-size: 0.7rem;
			font-weight: 900;
			letter-spacing: 1px;
		}
	}
</style>
