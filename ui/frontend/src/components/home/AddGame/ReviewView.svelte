<script lang="ts">
	import { fade } from "svelte/transition";

	export let foundExecutables: { path: string; name: string; icon: string | null }[];
	export let discardedExecutables: Set<string>;
	export let onToggleDiscard: (path: string) => void;
</script>

<div class="review-view" transition:fade={{ duration: 200 }}>
	<div class="found-grid">
		{#each foundExecutables as exe}
			<div
				class="exe-card"
				class:discarded={discardedExecutables.has(exe.path)}
				on:click={() => onToggleDiscard(exe.path)}
				role="button"
				tabindex="0"
				on:keydown={(e) => e.key === "Enter" && onToggleDiscard(exe.path)}
			>
				<div class="card-content">
					<div class="icon-container">
						{#if exe.icon}
							<img src={exe.icon} alt={exe.name} class="game-icon" />
						{:else}
							<div class="icon-placeholder">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									width="24"
									height="24"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									><rect x="2" y="2" width="20" height="20" rx="5" ry="5"></rect><path
										d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"
									></path><line x1="17.5" y1="6.5" x2="17.51" y2="6.5"></line></svg
								>
							</div>
						{/if}
					</div>
					<div class="name" title={exe.path}>{exe.name}</div>
				</div>
				<div class="discard-overlay">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="24"
						height="24"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						><polyline points="3 6 5 6 21 6"></polyline><path
							d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
						></path></svg
					>
				</div>
			</div>
		{/each}
	</div>
</div>

<style lang="scss">
	.review-view {
		display: flex;
		flex-direction: column;
	}

	.found-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
		gap: 16px;
	}

	.exe-card {
		position: relative;
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.05) 0%, rgba(255, 255, 255, 0.01) 100%);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 20px;
		padding: 20px 16px;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12px;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		cursor: pointer;
		overflow: hidden;

		&::after {
			content: "";
			position: absolute;
			inset: 0;
			background: radial-gradient(circle at top left, rgba(255, 255, 255, 0.05), transparent 70%);
			opacity: 0;
			transition: opacity 0.3s;
		}

		.discard-overlay {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background: rgba(239, 68, 68, 0.2);
			display: flex;
			align-items: center;
			justify-content: center;
			opacity: 0;
			transition: all 0.2s;
			color: #fff;
		}

		&.discarded {
			filter: grayscale(1);
			opacity: 0.4;
			border-color: rgba(239, 68, 68, 0.3);

			.discard-overlay {
				opacity: 1;
			}
		}

		&:hover:not(.discarded) {
			background: rgba(255, 255, 255, 0.1);
			border-color: var(--accent-primary);
			transform: translateY(-4px);
		}

		.card-content {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 12px;
			width: 100%;
		}

		.icon-container {
			width: 56px;
			height: 56px;
			display: flex;
			align-items: center;
			justify-content: center;

			.game-icon {
				width: 100%;
				height: 100%;
				object-fit: contain;
			}

			.icon-placeholder {
				color: rgba(255, 255, 255, 0.15);
			}
		}

		.name {
			font-size: 0.85rem;
			font-weight: 700;
			color: #fff;
			text-align: center;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
			width: 100%;
		}
	}
</style>
