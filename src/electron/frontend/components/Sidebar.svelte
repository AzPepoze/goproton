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

	function setActive(id: string) {
		onNavigate(id);
	}
</script>

<div class="sidebar">
	<div class="header">
		<h1>GoProton</h1>
	</div>

	<nav>
		{#each navItems as item}
			<button class:active={activePage === item.id} on:click={() => setActive(item.id)}>
				<img src={item.icon} alt={item.label} class="icon" />
				<span class="label">{item.label}</span>
			</button>
		{/each}
	</nav>

	<div class="footer">v1.0.0</div>
</div>

<style lang="scss">
	.sidebar {
		width: 240px;
		height: 100%;
		background: rgba(18, 18, 18, 0.5);
		border-right: 1px solid var(--glass-border);
		display: flex;
		flex-direction: column;
	}

	.header {
		padding: 48px 24px 32px 24px;
		text-align: center;

		h1 {
			font-size: 1.5rem;
			font-weight: 900;
			color: #fff;
			margin: 0;
			letter-spacing: 2px;
			text-transform: uppercase;
			opacity: 0.9;
		}
	}

	nav {
		flex: 1;
		padding: 0 16px;
		display: flex;
		flex-direction: column;
		gap: 4px;

		button {
			width: 100%;
			display: flex;
			align-items: center;
			gap: 14px;
			padding: 12px 16px;
			border-radius: 12px;
			border: 1px solid transparent;
			background: transparent;
			color: var(--text-muted);
			cursor: pointer;
			text-align: left;
			transition: all 0.2s;

			&:hover {
				background: rgba(255, 255, 255, 0.04);
				color: var(--text-main);
			}

			&.active {
				background: rgba(255, 255, 255, 0.07);
				border-color: rgba(255, 255, 255, 0.1);
				color: #fff;
				font-weight: 600;
			}

			.icon {
				width: 20px;
				height: 20px;
				opacity: 0.7;
				/* Force icon to be solid white */
				filter: brightness(0) invert(1);
				transition: opacity 0.2s;
			}

			&.active .icon {
				opacity: 1;
				/* Add a subtle glow to active icon */
				filter: brightness(0) invert(1) drop-shadow(0 0 2px rgba(255, 255, 255, 0.5));
			}
		}
	}

	.footer {
		padding: 24px;
		text-align: center;
		font-size: 0.65rem;
		color: var(--text-dim);
		letter-spacing: 1px;
		text-transform: uppercase;
	}
</style>
