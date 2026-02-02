<script lang="ts">
	import Sidebar from "./components/Sidebar.svelte";
	import Home from "./pages/Home.svelte";
	import Run from "./pages/Run.svelte";
	import Versions from "./pages/Versions.svelte";
	import Prefix from "./pages/Prefix.svelte";
	import Utils from "./pages/Utils.svelte";
	import EditLsfg from "./pages/EditLsfg.svelte";
	import NotificationHost from "./components/NotificationHost.svelte";
	import { fade, fly } from "svelte/transition";
	import { GetInitialLauncherPath, GetShouldEditLsfg } from "../wailsjs/go/main/App";
	import { onMount } from "svelte";
	import { navigationCommand } from "./stores/navigationStore";
	import { runState } from "./stores/runState";

	let activePage = "home";
	let editLsfgGamePath = "";

	onMount(async () => {
		try {
			const shouldEditLsfg = await GetShouldEditLsfg();
			const launcherPath = await GetInitialLauncherPath();

			if (shouldEditLsfg && launcherPath) {
				editLsfgGamePath = launcherPath;
				activePage = "editlsfg";
			} else if (launcherPath) {
				runState.update((state) => ({
					...state,
					options: {
						...state.options,
						LauncherPath: launcherPath,
					},
				}));
				activePage = "run";
			}
		} catch (e) {
			console.error("Error in App onMount:", e);
		}
	});

	// Subscribe to navigation commands
	navigationCommand.subscribe((cmd) => {
		if (cmd) {
			if (cmd.page === "editlsfg" && cmd.gamePath) {
				editLsfgGamePath = cmd.gamePath;
				activePage = "editlsfg";
			} else if (cmd.page) {
				activePage = cmd.page;
			}
			navigationCommand.set(null);
		}
	});

	function handleNavigate(page: string) {
		activePage = page;
	}
</script>

<main>
	{#if activePage !== "editlsfg"}
		<Sidebar {activePage} onNavigate={handleNavigate} />
	{/if}

	<div class="content" class:fullscreen={activePage === "editlsfg"}>
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
				{:else if activePage === "editlsfg"}
					<EditLsfg gamePath={editLsfgGamePath} />
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

		&.fullscreen {
			width: 100vw;
		}
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
