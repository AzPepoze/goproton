<script lang="ts">
	import { fade, scale } from "svelte/transition";

	export let show = false;
	export let title = "Settings";
	export let fullscreen = false;
	export let onClose: () => void = () => {};
	export let showDone = true;
	export let contentClass = "";

	function close() {
		onClose();
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === "Escape") {
			close();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

{#if show}
	<div
		class="modal-backdrop"
		on:click={close}
		on:keydown={(e) => e.key === "Escape" && close()}
		transition:fade={{ duration: 200 }}
		role="presentation"
	>
		<div
			class="modal-content glass {contentClass}"
			class:fullscreen
			on:click|stopPropagation
			on:keydown|stopPropagation={handleKeydown}
			transition:scale={{ duration: 200, start: 0.95 }}
			role="dialog"
			tabindex="0"
			aria-modal="true"
		>
			<div class="modal-header">
				<h3>{title}</h3>
				<button class="close-btn" on:click={close} aria-label="Close modal">
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
						><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"
						></line></svg
					>
				</button>
			</div>
			<div class="modal-body">
				<slot></slot>
			</div>
			{#if $$slots.footer || showDone}
				<div class="modal-footer">
					<slot name="footer">
						<button class="btn primary" on:click={close}>Done</button>
					</slot>
				</div>
			{/if}
		</div>
	</div>
{/if}

<style lang="scss">
	.modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.85);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal-content {
		width: 90%;
		max-width: 600px;
		background: #1a1a1a;
		border: 1px solid #333;
		border-radius: 16px;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.6);
		display: flex;
		flex-direction: column;
		max-height: 90vh; /* Default max height */

		&.fullscreen {
			width: 100%;
			height: 100%;
			max-width: none;
			max-height: none;
			border-radius: 0;
			border: none;
			background: var(--glass-bg); /* Match app bg */
		}
	}

	.modal-header {
		padding: 20px 32px;
		border-bottom: 1px solid var(--glass-border);
		display: flex;
		align-items: center;
		justify-content: space-between;
		flex-shrink: 0;

		h3 {
			margin: 0;
			font-size: 1.5rem;
			color: var(--text-main);
		}

		.close-btn {
			background: rgba(255, 255, 255, 0.05);
			border: 1px solid rgba(255, 255, 255, 0.1);
			color: var(--text-dim);
			width: 32px;
			height: 32px;
			border-radius: 10px;
			display: flex;
			align-items: center;
			justify-content: center;
			cursor: pointer;
			transition: all 0.2s;

			&:hover {
				color: var(--text-main);
				background: rgba(255, 255, 255, 0.1);
				transform: scale(1.1);
			}
		}
	}

	.modal-body {
		padding: 32px;
		overflow-y: auto;
		flex: 1; /* Take remaining space */
		display: flex;
		flex-direction: column;
	}

	.modal-footer {
		padding: 20px 32px;
		border-top: 1px solid var(--glass-border);
		display: flex;
		justify-content: flex-end;
		flex-shrink: 0;
	}
</style>
