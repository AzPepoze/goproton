<script lang="ts">
	import { fade, fly } from "svelte/transition";

	export let icon: string;
	export let title: string;
	export let subtitle: string;
	export let isPulsing = false;
	export let showSuccess = false;
	export let btnClass = "";
	export let onclick: () => void = () => {};
</script>

<button class="util-btn {btnClass}" class:pulsing={isPulsing} on:click={onclick}>
	<div class="content-left">
		<img src={icon} alt="" />
		<div class="btn-stack">
			<span>{title}</span>
			<small>{subtitle}</small>
		</div>
	</div>
	{#if showSuccess}
		<div class="check-indicator" in:fly={{ x: 20, duration: 400 }} out:fade={{ duration: 200 }}>
			<svg
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="3"
				stroke-linecap="round"
				stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg
			>
		</div>
	{/if}
</button>

<style lang="scss">
	.util-btn {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 12px 16px;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.05);
		border-radius: 12px;
		color: #fff;
		cursor: pointer;
		transition: all 0.3s cubic-bezier(0.23, 1, 0.32, 1);
		position: relative;
		overflow: hidden;

		.content-left {
			display: flex;
			align-items: center;
			gap: 12px;
		}

		img {
			width: 20px;
			height: 20px;
			opacity: 0.9;
			z-index: 2;
			transition: all 0.3s;
		}

		.btn-stack {
			display: flex;
			flex-direction: column;
			align-items: flex-start;
			z-index: 2;
		}

		span {
			font-size: 0.85rem;
			font-weight: 700;
			transition: all 0.3s;
		}

		small {
			font-size: 0.7rem;
			color: var(--text-dim);
			display: block;
		}

		&:hover {
			transform: translateY(-2px);
			img {
				opacity: 1;
			}
		}

		&.cleanup {
			background: rgba(254, 240, 138, 0.1);
			border-color: rgba(234, 179, 8, 0.3);

			img {
				filter: invert(79%) sepia(50%) saturate(1000%) hue-rotate(10deg) brightness(100%) contrast(100%);
			}
			span {
				color: #facc15;
			}
			small {
				color: rgba(250, 204, 21, 0.8);
			}

			&:hover {
				background: rgba(254, 240, 138, 0.15);
				border-color: rgba(234, 179, 8, 0.5);
				box-shadow: 0 4px 15px rgba(234, 179, 8, 0.2);
			}
			&.pulsing::after {
				background: linear-gradient(90deg, transparent, rgba(234, 179, 8, 0.4), transparent);
			}
		}

		&.cache {
			background: rgba(59, 130, 246, 0.08);
			border-color: rgba(59, 130, 246, 0.3);

			img {
				filter: invert(53%) sepia(90%) saturate(3000%) hue-rotate(200deg) brightness(100%) contrast(105%);
			}
			span {
				color: #60a5fa;
			}
			small {
				color: rgba(96, 165, 250, 0.8);
			}

			&:hover {
				background: rgba(59, 130, 246, 0.12);
				border-color: rgba(59, 130, 246, 0.5);
				box-shadow: 0 4px 15px rgba(59, 130, 246, 0.2);
			}
			&.pulsing::after {
				background: linear-gradient(90deg, transparent, rgba(59, 130, 246, 0.4), transparent);
			}
		}

		&.pulsing {
			pointer-events: none;

			&::after {
				content: "";
				position: absolute;
				top: 0;
				left: -150%;
				width: 150%;
				height: 100%;
				z-index: 1;
				animation: pulse-scan 1.5s cubic-bezier(0.4, 0, 0.2, 1);
			}
		}
	}

	.check-indicator {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 20px;
		height: 20px;
		background: rgba(255, 255, 255, 0.15);
		border-radius: 50%;
		color: #fff;
		padding: 3px;
		z-index: 2;

		svg {
			width: 100%;
			height: 100%;
		}
	}

	@keyframes pulse-scan {
		0% {
			transform: translateX(0);
		}
		100% {
			transform: translateX(250%);
		}
	}
</style>
