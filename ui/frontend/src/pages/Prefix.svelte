<script lang="ts">
  import { onMount } from 'svelte';
  import { PickFolder, RunPrefixTool, ScanProtonVersions, ListPrefixes, GetPrefixBaseDir, CreatePrefix } from '../../wailsjs/go/main/App';
  import Dropdown from '../components/Dropdown.svelte';
  import type { launcher } from '../../wailsjs/go/models';

  let prefixPath = "";
  let baseDir = "";
  let availablePrefixes: string[] = [];
  let protonVersions: launcher.ProtonTool[] = [];
  let protonOptions: string[] = [];
  let selectedProton = "";
  let isLoading = false;
  let newPrefixName = "";

  async function refreshPrefixes() {
    try {
      const list = await ListPrefixes();
      availablePrefixes = Array.isArray(list) ? list : [];
      baseDir = await GetPrefixBaseDir();
      if (!prefixPath && availablePrefixes.length > 0) {
        selectPrefix(availablePrefixes[0]);
      } else if (!prefixPath) {
        prefixPath = baseDir + "/Default";
      }
    } catch (err) {
      console.error(err);
      availablePrefixes = [];
    }
  }

  onMount(async () => {
    try {
      protonVersions = await ScanProtonVersions();
      protonOptions = protonVersions.map(t => t.DisplayName);
      if (protonVersions.length > 0) {
        selectedProton = protonOptions[0];
      }
      await refreshPrefixes();
    } catch (err) {
      console.error(err);
    }
  });

  function selectPrefix(name: string) {
    prefixPath = baseDir + "/" + name;
  }

  async function handleBrowse() {
    try {
      const path = await PickFolder();
      if (path) prefixPath = path;
    } catch (err) {
      console.error(err);
    }
  }

  async function handleCreatePrefix() {
    if (!newPrefixName) return;
    const name = newPrefixName;
    try {
      await CreatePrefix(name);
      newPrefixName = "";
      await refreshPrefixes();
      selectPrefix(name);
    } catch (err) {
      alert(`Failed to create prefix: ${err}`);
    }
  }

  function handleProtonChange(event: CustomEvent<string>) {
    selectedProton = event.detail;
  }

  async function runTool(tool: string) {
    if (isLoading) return;
    if (!prefixPath) {
      alert("Please select or create a prefix first.");
      return;
    }
    isLoading = true;
    let cleanName = selectedProton;
    if (cleanName.startsWith("(Steam) ")) {
      cleanName = cleanName.substring(8);
    }
    try {
      await RunPrefixTool(prefixPath, tool, cleanName);
    } catch (err) {
      alert(`Failed to run ${tool}: ${err}`);
    } finally {
      isLoading = false;
    }
  }

  $: currentPrefixName = prefixPath.startsWith(baseDir) 
    ? prefixPath.replace(baseDir + "/", "") 
    : prefixPath;
</script>

<div class="prefix-container">
  <h1 class="page-title">Prefix Manager</h1>
  <div class="main-layout">
    <div class="sidebar-section glass">
      <div class="section-header">
        <h2>Available Prefixes</h2>
      </div>
      <div class="prefix-list">
        {#each availablePrefixes as name}
          <button 
            class="prefix-item" 
            class:active={currentPrefixName === name}
            on:click={() => selectPrefix(name)}
          >
            <span class="folder-icon">📁</span>
            <span class="name">{name}</span>
          </button>
        {/each}
        {#if availablePrefixes.length === 0}
          <div class="empty-state">No prefixes found in default directory.</div>
        {/if}
      </div>
      <div class="add-prefix-area">
        <input 
          type="text" 
          placeholder="New prefix..." 
          bind:value={newPrefixName} 
          class="input sm"
          on:keydown={(e) => e.key === 'Enter' && handleCreatePrefix()}
        />
        <button class="btn primary sm" on:click={handleCreatePrefix}>Create</button>
      </div>
    </div>
    <div class="content-section">
      <div class="config-card glass">
        <div class="form-group">
          <label>Selected Prefix Path</label>
          <div class="input-group">
            <input type="text" class="input" bind:value={prefixPath} readonly />
            <button class="btn" on:click={handleBrowse}>Browse Other</button>
          </div>
        </div>
        <div class="form-group">
          <label>Runtime Environment (Proton)</label>
          <Dropdown 
            options={protonOptions} 
            bind:value={selectedProton} 
            on:change={handleProtonChange}
          />
        </div>
      </div>
      <div class="tools-grid">
        <button class="tool-card" on:click={() => runTool('winecfg')}>
          <div class="icon">⚙️</div>
          <div class="text">
            <h3>Winecfg</h3>
            <p>Settings</p>
          </div>
        </button>
        <button class="tool-card" on:click={() => runTool('regedit')}>
          <div class="icon">📝</div>
          <div class="text">
            <h3>Regedit</h3>
            <p>Registry</p>
          </div>
        </button>
        <button class="tool-card" on:click={() => runTool('cmd')}>
          <div class="icon">💻</div>
          <div class="text">
            <h3>CMD</h3>
            <p>Terminal</p>
          </div>
        </button>
        <button class="tool-card" on:click={() => runTool('winetricks')}>
          <div class="icon">🪄</div>
          <div class="text">
            <h3>Winetricks</h3>
            <p>Extras</p>
          </div>
        </button>
      </div>
    </div>
  </div>
</div>

<style lang="scss">
  .prefix-container {
    padding: 24px;
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    box-sizing: border-box;
  }
  .page-title {
    font-size: 1.75rem;
    font-weight: 800;
    color: var(--text-main);
    margin: 0 0 20px 0;
    flex-shrink: 0;
  }
  .main-layout {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 20px;
    flex: 1;
    min-height: 0;
    overflow: hidden;
    padding-bottom: 4px; /* Space for bottom borders */
  }
  .sidebar-section {
    display: flex;
    flex-direction: column;
    border-radius: 16px;
    overflow: hidden;
    border: 1px solid var(--glass-border);
    background: rgba(255, 255, 255, 0.02);
    height: 100%;
    box-sizing: border-box;

    .section-header {
      padding: 16px;
      border-bottom: 1px solid var(--glass-border);
      h2 { font-size: 0.9rem; margin: 0; color: var(--text-dim); text-transform: uppercase; letter-spacing: 1px; }
    }
  }
  .prefix-list {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  .prefix-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: transparent;
    border: none;
    border-radius: 8px;
    color: var(--text-main);
    cursor: pointer;
    text-align: left;
    transition: all 0.2s;
    &:hover { background: rgba(255, 255, 255, 0.05); }
    &.active {
      background: var(--accent-primary);
      color: black;
      font-weight: 600;
      .folder-icon { color: black; }
    }
    .folder-icon { font-size: 1.1rem; }
    .name { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  }
  .add-prefix-area {
    padding: 16px;
    border-top: 1px solid var(--glass-border);
    display: flex;
    flex-direction: column;
    gap: 8px;
    background: rgba(0, 0, 0, 0.1);
  }
  .content-section {
    display: flex;
    flex-direction: column;
    gap: 24px;
    overflow-y: auto;
    padding: 2px 12px 24px 2px; /* Extra padding for bottom and right to prevent clipping */
    height: 100%;
    box-sizing: border-box;
  }
  .config-card {
    padding: 24px;
    border-radius: 16px;
    border: 1px solid var(--glass-border);
    display: flex;
    flex-direction: column;
    gap: 20px;
  }
  .form-group {
    label {
      display: block;
      font-size: 0.8rem;
      font-weight: 600;
      color: var(--text-dim);
      margin-bottom: 8px;
    }
  }
  .input-group {
    display: flex;
    gap: 8px;
    input { flex: 1; }
  }
  .tools-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
  }
  .tool-card {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--glass-border);
    border-radius: 16px;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 16px;
    cursor: pointer;
    transition: all 0.2s;
    color: var(--text-main);
    text-align: left;
    &:hover {
      background: rgba(255, 255, 255, 0.08);
      border-color: var(--accent-primary);
      transform: translateY(-2px);
    }
    .icon { font-size: 1.75rem; }
    h3 { font-size: 1rem; margin: 0; font-weight: 700; }
    p { font-size: 0.75rem; margin: 2px 0 0 0; color: var(--text-dim); }
  }
  .empty-state {
    padding: 32px;
    text-align: center;
    color: var(--text-dim);
    font-size: 0.8rem;
  }
  .input.sm { padding: 8px 12px; font-size: 0.85rem; }
  .btn.sm { padding: 8px; font-size: 0.85rem; }
</style>
