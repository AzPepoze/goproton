<script lang="ts">
	import { notifications } from "../notificationStore";
	import { flip } from "svelte/animate";
	import { fade, fly } from "svelte/transition";
</script>

<div class="notification-container">
	{#each $notifications as n (n.id)}
		<div
			class="notification-card {n.type} glass"
			animate:flip={{ duration: 300 }}
			in:fly={{ y: 20, duration: 300 }}
			out:fade={{ duration: 200 }}
		>
			<div class="icon">
				{#if n.type === "error"}
					❌
				{:else if n.type === "success"}
					✅
				{:else}
					ℹ️
				{/if}
			</div>
			<div class="message">{n.message}</div>
			<button class="close" on:click={() => notifications.remove(n.id)}>&times;</button>
		</div>
	{/each}
</div>

<style lang="scss">
	.notification-container {
		position: fixed;
		bottom: 24px;
		right: 24px;
		z-index: 9999;
		display: flex;
		flex-direction: column;
		gap: 12px;
		max-width: 400px;
		pointer-events: none;
	}

	.notification-card {
		pointer-events: auto;
		padding: 16px 20px;
		border-radius: 12px;
		display: flex;
		align-items: center;
		gap: 12px;
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.4);
		border: 1px solid var(--glass-border);
		background: #1a1a1a;
		color: white;

		&.error {
			border-left: 4px solid #ef4444;
		}
		&.success {
			border-left: 4px solid #10b981;
		}
		&.info {
			border-left: 4px solid #3b82f6;
		}

		.icon {
			font-size: 1.25rem;
		}
		.message {
			flex: 1;
			font-size: 0.9rem;
			font-weight: 500;
		}

		.close {
			background: none;
			border: none;
			color: var(--text-dim);
			font-size: 1.5rem;
			cursor: pointer;
			padding: 0 4px;
			line-height: 1;
			&:hover {
				color: white;
			}
		}
	}
</style>
