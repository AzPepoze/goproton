<script lang="ts">
  import Sidebar from './components/Sidebar.svelte';
  import Home from './pages/Home.svelte';
  import Run from './pages/Run.svelte';
  import Versions from './pages/Versions.svelte';
  import Prefix from './pages/Prefix.svelte';
  import Utils from './pages/Utils.svelte';
  import NotificationHost from './components/NotificationHost.svelte';

  let activePage = "home";

  function handleNavigate(event: CustomEvent<string>) {
    activePage = event.detail;
  }
</script>

<main>
  <Sidebar {activePage} on:navigate={handleNavigate} />

  <div class="content">
    {#if activePage === 'home'}
      <Home />
    {:else if activePage === 'run'}
      <Run />
    {:else if activePage === 'versions'}
      <Versions />
    {:else if activePage === 'prefix'}
      <Prefix />
    {:else if activePage === 'utils'}
      <Utils />
    {:else}
      <div class="placeholder">
        Page "{activePage}" not implemented yet.
      </div>
    {/if}
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
  }

  .content {
    flex: 1;
    height: 100%;
    overflow: hidden;
    /* Subtle neutral gradient for depth */
    background: radial-gradient(circle at 50% 50%, rgba(255, 255, 255, 0.02) 0%, transparent 100%);
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