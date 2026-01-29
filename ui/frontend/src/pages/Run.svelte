<script lang="ts">
  import { onMount } from 'svelte';
  import { WindowHide } from '../../wailsjs/runtime/runtime';
  import { PickFile, PickFolder, ScanProtonVersions, RunGame } from '../../wailsjs/go/main/App';
  import type { launcher } from '../../wailsjs/go/models';
  import SlideButton from '../components/SlideButton.svelte';
  import Dropdown from '../components/Dropdown.svelte';

  // State
  let gamePath = "";
  let prefixPath = "~/GoProton/Prefixes/Default";
  let protonVersions: launcher.ProtonTool[] = [];
  let protonOptions: string[] = [];
  let selectedProton = "";
  let isLoadingProton = true;

  // Toggles
  let enableMango = false;
  let enableGamemode = false;
  let enableGamescope = false;
  let showLogsWindow = true; 
  let gamescopeW = "1920";
  let gamescopeH = "1080";
  let gamescopeR = "60";

  onMount(async () => {
    // 1. Load Proton Versions
    try {
      const tools = await ScanProtonVersions();
      if (tools) {
        protonVersions = tools;
        protonOptions = tools.map(t => t.DisplayName);
        if (protonOptions.length > 0) {
          selectedProton = protonOptions[0];
        }
      }
    } catch (err) {
      console.error("Failed to scan proton:", err);
      protonOptions = ["Error scanning tools"];
    } finally {
      isLoadingProton = false;
    }
  });

  async function handleBrowseGame() {
    try {
      const path = await PickFile();
      if (path) gamePath = path;
    } catch (err) {
      console.error(err);
    }
  }

  async function handleBrowsePrefix() {
    try {
      const path = await PickFolder();
      if (path) prefixPath = path;
    } catch (err) {
      console.error(err);
    }
  }

  function handleProtonChange(event: CustomEvent<string>) {
    selectedProton = event.detail;
  }

  async function handleLaunch() {
    if (!gamePath) {
      alert("Please select a game executable.");
      return;
    }

    // Find tool by display name
    const tool = protonVersions.find(p => p.DisplayName === selectedProton);
    
    // Clean name for Pattern
    let cleanName = selectedProton;
    if (cleanName.startsWith("(Steam) ")) {
      cleanName = cleanName.substring(8);
    }

    const opts: launcher.LaunchOptions = {
      GamePath: gamePath,
      PrefixPath: prefixPath,
      ProtonPattern: cleanName,
      ProtonPath: tool ? tool.Path : "",
      EnableGamescope: enableGamescope,
      GamescopeW: gamescopeW,
      GamescopeH: gamescopeH,
      GamescopeR: gamescopeR,
      EnableMangoHud: enableMango,
      EnableGamemode: enableGamemode,
    };

    try {
      await RunGame(opts, showLogsWindow);
      // Hide the UI immediately as requested
      WindowHide();
    } catch (err) {
      console.error("Launch failed:", err);
      alert(`Launch failed: ${err}`);
    }
  }
</script>

<div class="run-container">
  <div class="header-row">
    <h1 class="page-title">Launch Configuration</h1>
  </div>

  <div class="form-container">
    
    <!-- Game Path -->
    <div class="form-group">
      <label>Game Executable</label>
      <div class="input-group">
        <input 
          type="text" 
          bind:value={gamePath} 
          placeholder="Select .exe file..." 
          class="input"
        />
        <button on:click={handleBrowseGame} class="btn">Browse</button>
      </div>
    </div>

    <!-- Prefix Path -->
    <div class="form-group">
      <label>WINEPREFIX</label>
      <div class="input-group">
        <input 
          type="text" 
          bind:value={prefixPath} 
          class="input"
        />
        <button on:click={handleBrowsePrefix} class="btn">Browse</button>
      </div>
    </div>

    <!-- Proton Version -->
    <div class="form-group">
      <label>Proton Version</label>
      <Dropdown 
        options={protonOptions} 
        bind:value={selectedProton} 
        placeholder={isLoadingProton ? "Scanning..." : "Select Version"}
        disabled={isLoadingProton}
        on:change={handleProtonChange}
      />
    </div>

    <!-- Toggles Grid -->
    <div class="toggles-grid">
      <SlideButton 
        bind:checked={enableMango} 
        label="MangoHud" 
        subtitle="Performance overlay" 
      />
      
      <SlideButton 
        bind:checked={enableGamemode} 
        label="GameMode" 
        subtitle="Optimize priorities" 
      />

      <SlideButton 
        bind:checked={enableGamescope} 
        label="Gamescope" 
        subtitle="Micro-compositor" 
      />

      <SlideButton 
        bind:checked={showLogsWindow} 
        label="Show Logs" 
        subtitle="Open logs in terminal" 
      />
    </div>

    <!-- Gamescope Settings -->
    {#if enableGamescope}
      <div class="gamescope-settings glass">
        <div class="form-group">
          <label>Width</label>
          <input type="text" class="input sm" bind:value={gamescopeW} />
        </div>
        <div class="form-group">
          <label>Height</label>
          <input type="text" class="input sm" bind:value={gamescopeH} />
        </div>
        <div class="form-group">
          <label>Refresh Rate</label>
          <input type="text" class="input sm" bind:value={gamescopeR} />
        </div>
      </div>
    {/if}

    <!-- Launch Button -->
    <div class="action-area">
      <button class="btn primary launch-btn" on:click={handleLaunch}>
        🚀 LAUNCH GAME
      </button>
    </div>
  </div>
</div>

<style lang="scss">
  .run-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 32px;
    overflow: hidden;
  }

  .form-container {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 24px;
    overflow-y: auto;
    padding-right: 8px;
  }

  .form-group label {
    display: block;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--text-muted);
    margin-bottom: 8px;
  }

  .input-group {
    display: flex;
    gap: 12px;
    width: 100%;
    .input { flex: 1; }
  }

  .toggles-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 16px;
    margin-top: 8px;
  }

  .gamescope-settings {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 16px;
    padding: 16px;
    border-radius: 12px;
  }

  .action-area {
    margin-top: 16px;
    margin-bottom: 32px;
  }

  .launch-btn {
    width: 100%;
    padding: 16px;
    font-size: 1.125rem;
    font-weight: bold;
    border-radius: 12px;
    text-transform: uppercase;
    letter-spacing: 1px;
  }

  .input.sm { padding: 6px 10px; font-size: 0.875rem; }
  .page-title { font-size: 2rem; font-weight: bold; color: var(--text-main); margin: 0 0 24px 0; }
</style>