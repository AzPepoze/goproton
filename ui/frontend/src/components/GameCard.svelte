<script lang="ts">
	import GridCard from "./gameCards/GridCard.svelte";
	import ListGridCard from "./gameCards/ListGridCard.svelte";
	import PerspectiveCard from "./gameCards/PerspectiveCard.svelte";

	export let game: any;
	export let icon: string = "";
	export let isRunning: boolean = false;
	export let view: "grid" | "list-grid" | "perspective" = "grid";
	export let active: boolean = false; // For perspective view
</script>

{#if view === "grid"}
	<GridCard {game} {icon} {isRunning} on:launch on:configure />
{:else if view === "list-grid"}
	<ListGridCard {game} {icon} {isRunning} on:launch on:configure />
{:else}
	<PerspectiveCard {game} {icon} {isRunning} {active} on:launch on:configure />
{/if}
