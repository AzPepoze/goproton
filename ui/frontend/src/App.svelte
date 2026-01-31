<script lang="ts">
	import Sidebar from "./components/Sidebar.svelte";
	import Home from "./pages/Home.svelte";
	import Run from "./pages/Run.svelte";
	import Versions from "./pages/Versions.svelte";
	import Prefix from "./pages/Prefix.svelte";
	import Utils from "./pages/Utils.svelte";
	import NotificationHost from "./components/NotificationHost.svelte";
	import { fade, fly } from "svelte/transition";

	let activePage = "home";

	function handleNavigate(page: string) {
		activePage = page;
	}
</script>

<main>
	<Sidebar {activePage} onNavigate={handleNavigate} />

	<div class="content">
		{#key activePage}
			<div class="page-wrapper" in:fly={{ y: 10, duration: 300, delay: 150 }} out:fade={{ duration: 150 }}>
				{#if activePage === "home"}
					<Home />
				{:else if activePage === "run"}
					<Run />
				{:else if activePage === "versions"}
					<Versions />
				{:else if activePage === "prefix"}
					<Prefix />
				{:else if activePage === "utils"}
					<Utils />
				{:else}
					<div class="placeholder">
						Page "{activePage}" not implemented yet.
					</div>
				{/if}
			</div>
		{/key}
	</div>

	<NotificationHost />
</main>

<style lang="scss">
	main {
		display: flex;
		height: 100vh;
		width: 100vw;
		background-color: var(--glass-bg);
		color: var(--text-main);
		user-select: none;
		overflow: hidden; /* Prevent scrollbar flicker during transition */
	}

	.content {
		flex: 1;
		height: 100%;
		position: relative; /* Important for absolute positioning of page wrapper */
		/* Subtle neutral gradient for depth */
		background: radial-gradient(circle at 50% 50%, rgba(255, 255, 255, 0.02) 0%, transparent 100%);
	}

	/* Wrapper to handle transition positioning */
	.page-wrapper {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		overflow-y: auto; /* Allow scrolling inside the page */
	}

	.placeholder {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100%;
		color: var(--text-dim);
		font-size: 0.9rem;
		font-style: italic;
	}
</style>
