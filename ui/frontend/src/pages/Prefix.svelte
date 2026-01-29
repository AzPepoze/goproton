<script lang="ts">
  import { onMount } from 'svelte';
  import { PickFolder, RunPrefixTool, ScanProtonVersions } from '../../wailsjs/go/main/App';
  import Dropdown from '../components/Dropdown.svelte';
  import type { launcher } from '../../wailsjs/go/models';

  let prefixPath = "~/GoProton/Prefixes/Default";
  let protonVersions: launcher.ProtonTool[] = [];
  let protonOptions: string[] = [];
  let selectedProton = "";
  let isLoading = false;

  onMount(async () => {
    // Load protons for running tools
    try {
      protonVersions = await ScanProtonVersions();
      protonOptions = protonVersions.map(t => t.DisplayName);
      if (protonVersions.length > 0) {
        selectedProton = protonVersions[0].DisplayName;
      }
    } catch (err) {
      console.error(err);
    }
  });

  async function handleBrowse() {
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

  async function runTool(tool: string) {
    if (isLoading) return;
    isLoading = true;
    
    // Clean name logic same as launcher
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
</script>

<div class="prefix-container">
  <h1 class="page-title">Prefix Manager</h1>

  <div class="config-section">
    <!-- Path Selection -->
    <div class="form-group">
      <label>Target Prefix</label>
      <div class="input-group">
        <input type="text" class="input" bind:value={prefixPath} />
        <button class="btn" on:click={handleBrowse}>Browse</button>
      </div>
    </div>

    <!-- Proton Selection (Env) -->
    <div class="form-group">
      <label>Environment (Proton Version)</label>
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
      <h3>Wine Configuration</h3>
      <p>Manage drives, audio, and graphics settings.</p>
    </button>

    <button class="tool-card" on:click={() => runTool('regedit')}>
      <div class="icon">📝</div>
      <h3>Registry Editor</h3>
      <p>Modify the Windows registry keys.</p>
    </button>

    <button class="tool-card" on:click={() => runTool('cmd')}>
      <div class="icon">💻</div>
      <h3>Command Prompt</h3>
      <p>Open a terminal inside the prefix.</p>
    </button>

    <button class="tool-card" on:click={() => runTool('winetricks')}>
      <div class="icon">🪄</div>
      <h3>Winetricks</h3>
      <p>Install libraries and dependencies.</p>
    </button>
  </div>
</div>

<style lang="scss">
  .prefix-container {
    padding: 32px;
    height: 100%;
    display: flex;
    flex-direction: column;
    background-color: var(--glass-bg);
  }

  .page-title {
    font-size: 2rem;
    font-weight: 800;
    color: var(--text-main);
    margin: 0 0 32px 0;
  }

  .config-section {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--glass-border);
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 32px;
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .form-group {
    label {
      display: block;
      font-size: 0.875rem;
      font-weight: 600;
      color: var(--text-muted);
      margin-bottom: 8px;
    }
  }

  .input-group {
    display: flex;
    gap: 8px;
  }

  .tools-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 24px;
  }

  .tool-card {
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid var(--glass-border);
    border-radius: 16px;
    padding: 24px;
    text-align: left;
    cursor: pointer;
    transition: all 0.2s;
    color: var(--text-main);

    &:hover {
      background: rgba(255, 255, 255, 0.06);
      border-color: var(--accent-primary);
      transform: translateY(-4px);
    }

    &:active {
      transform: translateY(-2px);
    }

    .icon {
      font-size: 2rem;
      margin-bottom: 16px;
    }

    h3 {
      font-size: 1.1rem;
      font-weight: 700;
      margin: 0 0 8px 0;
    }

    p {
      font-size: 0.85rem;
      color: var(--text-muted);
      margin: 0;
      line-height: 1.4;
    }
  }
</style>
