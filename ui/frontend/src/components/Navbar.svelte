<script lang="ts">
	import homeIcon from "../icons/home.svg";
	import runIcon from "../icons/run.svg";
	import versionsIcon from "../icons/versions.svg";
	import prefixIcon from "../icons/prefix.svg";

	export let activePage: string = "home";
	export let onNavigate: (page: string) => void = () => {};

	const navItems = [
		{ id: "home", label: "Home", icon: homeIcon },
		{ id: "run", label: "Run", icon: runIcon },
		{ id: "versions", label: "Versions", icon: versionsIcon },
		{ id: "prefix", label: "Prefix", icon: prefixIcon },
		{ id: "utils", label: "Utils", icon: runIcon },
	];

	let navbarRef: HTMLElement;
	let indicatorStyle = "";

	function updateIndicator(id: string) {
		if (!navbarRef) return;
		const activeEl = navbarRef.querySelector(`button[data-id="${id}"]`) as HTMLElement;
		if (activeEl) {
			const top = activeEl.offsetTop;
			const height = activeEl.offsetHeight;
			indicatorStyle = `top: ${top}px; height: ${height}px; opacity: 1;`;
		}
	}

	$: {
		if (activePage && navbarRef) {
			setTimeout(() => updateIndicator(activePage), 0);
		}
	}

	function setActive(id: string) {
		onNavigate(id);
	}
</script>

<div class="navbar-wrapper">
	<nav class="navbar" bind:this={navbarRef}>
		<div class="moving-indicator" style={indicatorStyle}></div>
		{#each navItems as item}
			<button
				class="nav-item"
				data-id={item.id}
				class:active={activePage === item.id}
				on:click={() => setActive(item.id)}
			>
				<div class="icon-container">
					<img src={item.icon} alt={item.label} class="icon" />
					{#if activePage === item.id}
						<div class="glow-ring"></div>
					{/if}
				</div>
				<span class="label">{item.label}</span>
			</button>
		{/each}
	</nav>
</div>

<style lang="scss">
	.navbar-wrapper {
		z-index: 1000;
		display: flex;
		justify-content: center;
	}

	.navbar {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8px;
		padding: 16px 0;
		background: transparent;
		border: none;
		box-shadow: none;
	}

	.moving-indicator {
		position: absolute;
		left: 0;
		width: 3px;
		background: #ffffff;
		border-radius: 0 4px 4px 0;
		transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
		pointer-events: none;
		opacity: 0;
		z-index: 0;
	}

	.nav-item {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 2px;
		padding: 12px 0;
		border: none;
		background: transparent;
		color: var(--text-dim);
		cursor: pointer;
		border-radius: 20px;
		transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
		min-width: 68px; /* Narrower items */
		z-index: 1;

		&:hover {
			color: var(--text-main);

			.icon {
				transform: translateY(-2px);
				opacity: 1;
			}
		}

		&.active {
			color: #fff;

			.icon {
				transform: scale(1.1);
				opacity: 1;
				filter: brightness(0) invert(1) drop-shadow(0 0 8px rgba(255, 255, 255, 0.4));
			}

			.label {
				opacity: 1;
			}
		}
	}

	.icon-container {
		position: relative;
		width: 22px;
		height: 22px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.icon {
		width: 100%;
		height: 100%;
		filter: brightness(0) invert(1);
		opacity: 0.6;
		transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
	}

	.label {
		font-size: 0.6rem;
		font-weight: 800;
		text-transform: uppercase;
		letter-spacing: 0.8px;
		opacity: 0.6;
		transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
	}

	.glow-ring {
		position: absolute;
		width: 140%;
		height: 140%;
		border-radius: 50%;
		background: radial-gradient(circle, rgba(255, 255, 255, 0.15) 0%, transparent 70%);
		animation: pulse 2s infinite;
	}

	@keyframes pulse {
		0% {
			transform: scale(0.8);
			opacity: 0;
		}
		50% {
			transform: scale(1.2);
			opacity: 0.4;
		}
		100% {
			transform: scale(1.4);
			opacity: 0;
		}
	}
</style>
