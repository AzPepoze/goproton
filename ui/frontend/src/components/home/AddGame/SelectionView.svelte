<script lang="ts">
	import { fade } from "svelte/transition";

	export let onAddFile: () => void;
	export let onAddFolder: () => void;
</script>

<div class="selection-view" transition:fade={{ duration: 200 }}>
	<div class="mode-buttons">
		<button class="mode-btn" on:click={onAddFile}>
			<div class="icon-box">
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
					><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline
						points="14 2 14 8 20 8"
					></polyline><line x1="12" y1="18" x2="12" y2="12"></line><line x1="9" y1="15" x2="15" y2="15"
					></line></svg
				>
			</div>
			<div class="btn-text">
				<span class="title">Single File</span>
				<span class="desc">Add a specific .exe game file</span>
			</div>
		</button>

		<button class="mode-btn" on:click={onAddFolder}>
			<div class="icon-box">
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
					><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"
					></path><line x1="12" y1="11" x2="12" y2="17"></line><line x1="9" y1="14" x2="15" y2="14"
					></line></svg
				>
			</div>
			<div class="btn-text">
				<span class="title">Scan Folder</span>
				<span class="desc">Batch add games from directory</span>
			</div>
		</button>
	</div>
</div>

<style lang="scss">
	.selection-view {
		display: flex;
		flex-direction: column;
	}

	.mode-buttons {
		display: grid;
		grid-template-columns: 1fr;
		gap: 16px;
	}

	.mode-btn {
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.05) 0%, rgba(255, 255, 255, 0.01) 100%);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 20px;
		padding: 24px;
		display: flex;
		align-items: center;
		gap: 24px;
		cursor: pointer;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		text-align: left;
		position: relative;
		overflow: hidden;

		&::after {
			content: "";
			position: absolute;
			top: 0;
			left: 0;
			right: 0;
			bottom: 0;
			background: radial-gradient(circle at top left, rgba(255, 255, 255, 0.1), transparent 70%);
			opacity: 0;
			transition: opacity 0.3s;
		}

		&:hover {
			background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.03) 100%);
			border-color: rgba(255, 255, 255, 0.3);
			transform: translateY(-4px);
			box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);

			&::after {
				opacity: 1;
			}

			.icon-box {
				background: #fff;
				color: #000;
				transform: rotate(-10deg) scale(1.1);
				box-shadow: 0 0 20px rgba(255, 255, 255, 0.2);
			}
		}

		.icon-box {
			width: 64px;
			height: 64px;
			background: rgba(255, 255, 255, 0.05);
			border-radius: 18px;
			display: flex;
			align-items: center;
			justify-content: center;
			color: #fff;
			transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
			border: 1px solid rgba(255, 255, 255, 0.1);
		}

		.btn-text {
			.title {
				display: block;
				font-size: 1.25rem;
				font-weight: 800;
				color: #fff;
				margin-bottom: 4px;
			}
			.desc {
				font-size: 0.9rem;
				color: var(--text-dim);
			}
		}
	}
</style>
