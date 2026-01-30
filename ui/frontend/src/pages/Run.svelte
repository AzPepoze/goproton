<script lang="ts">
  import { onMount } from 'svelte';
  import { WindowHide } from '../../wailsjs/runtime/runtime';
  import { PickFile, PickFolder, ScanProtonVersions, RunGame, GetConfig, ListPrefixes, GetPrefixBaseDir, GetSystemToolsStatus, DetectLosslessDll, PickFileCustom } from '../../wailsjs/go/main/App';
  import type { launcher } from '../../wailsjs/go/models';
  import SlideButton from '../components/SlideButton.svelte';
  import Dropdown from '../components/Dropdown.svelte';
  import Modal from '../components/Modal.svelte';
  import { notifications } from '../notificationStore';

  // State
  let gamePath = "";
  let prefixPath = "";
  let baseDir = "";
  let availablePrefixes: string[] = [];
  let systemStatus: launcher.SystemToolsStatus = { hasGamescope: false, hasMangoHud: false, hasGameMode: false };
  let protonVersions: launcher.ProtonTool[] = [];
  let protonOptions: string[] = [];
  let selectedProton = "";
  let selectedPrefixName = "Default";
  let isLoadingProton = true;

  // Toggles
  let enableMango = false;
  let enableGamemode = false;
  let enableLsfg = false;
  let lsfgMultiplier = "2";
  let lsfgPerfMode = false;
  let lsfgDllPath = "";
  let enableGamescope = false;
  let showLogsWindow = false; 
  let gamescopeW = "1920";
  let gamescopeH = "1080";
  let gamescopeR = "60";

  let showGamescopeModal = false;
  let showLsfgModal = false;
  let showValidationModal = false;
  let missingToolsList: string[] = [];

  async function loadConfigForGame(path: string) {
    try {
      const config = await GetConfig(path);
      if (config) {
        prefixPath = config.PrefixPath;
        if (prefixPath.startsWith(baseDir)) {
          selectedPrefixName = prefixPath.replace(baseDir + "/", "");
        } else {
          selectedPrefixName = "Custom...";
        }
        const match = protonVersions.find(p => p.Path === config.ProtonPath);
        if (match) {
          selectedProton = match.DisplayName;
        } else if (config.ProtonPattern) {
          selectedProton = config.ProtonPattern;
        }
        enableMango = config.EnableMangoHud;
        enableGamemode = config.EnableGamemode;
        enableLsfg = config.EnableLsfgVk;
        lsfgMultiplier = config.LsfgMultiplier || "2";
        lsfgPerfMode = config.LsfgPerfMode;
        lsfgDllPath = config.LsfgDllPath || "";
        enableGamescope = config.EnableGamescope;
        gamescopeW = config.GamescopeW;
        gamescopeH = config.GamescopeH;
        gamescopeR = config.GamescopeR;
      }
      
      // Auto detect DLL if still empty
      if (!lsfgDllPath) {
        lsfgDllPath = await DetectLosslessDll();
      }
    } catch (err) {}
  }

  onMount(async () => {
    try {
      const [tools, prefixes, base, sysStatus, detectedDll] = await Promise.all([
        ScanProtonVersions(),
        ListPrefixes(),
        GetPrefixBaseDir(),
        GetSystemToolsStatus(),
        DetectLosslessDll()
      ]);
      if (tools) {
        protonVersions = tools;
        protonOptions = tools.map(t => t.DisplayName);
        if (protonOptions.length > 0) {
          selectedProton = protonOptions[0];
        }
      }
      availablePrefixes = Array.isArray(prefixes) ? prefixes : ["Default"];
      baseDir = base;
      systemStatus = sysStatus;
      lsfgDllPath = detectedDll;
      
      if (!prefixPath) {
        prefixPath = baseDir + "/Default";
        selectedPrefixName = "Default";
      }
    } catch (err) {
      console.error("Failed to initialize:", err);
    } finally {
      isLoadingProton = false;
    }
  });

  function handlePrefixChange(event: CustomEvent<string>) {
    const name = event.detail;
    if (name !== "Custom...") {
      prefixPath = baseDir + "/" + name;
    }
  }

  async function handleBrowseGame() {
    try {
      const path = await PickFile();
      if (path) {
        gamePath = path;
        await loadConfigForGame(path);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function handleBrowsePrefix() {
    try {
      const path = await PickFolder();
      if (path) {
        prefixPath = path;
        selectedPrefixName = "Custom...";
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function handleBrowseDll() {
    try {
      const path = await PickFileCustom("Select Lossless.dll", [{ DisplayName: "Lossless.dll", Pattern: "Lossless.dll" }]);
      if (path) lsfgDllPath = path;
    } catch (err) {
      console.error(err);
    }
  }

  function handleProtonChange(event: CustomEvent<string>) {
    selectedProton = event.detail;
  }

  async function handleLaunch() {
    if (!gamePath) {
      notifications.add("Please select a game executable.", "error");
      return;
    }
    
    if (enableLsfg && !lsfgDllPath) {
      notifications.add("LSFG-VK requires Lossless.dll. Please set the path in LSFG settings.", "error");
      showLsfgModal = true;
      return;
    }

    missingToolsList = [];
    if (enableGamescope && !systemStatus.hasGamescope) missingToolsList.push("Gamescope");
    if (enableMango && !systemStatus.hasMangoHud) missingToolsList.push("MangoHud");
    if (enableGamemode && !systemStatus.hasGameMode) missingToolsList.push("GameMode");
    if (missingToolsList.length > 0) {
      showValidationModal = true;
      return;
    }
    await proceedToLaunch();
  }

  async function proceedToLaunch() {
    showValidationModal = false;
    const tool = protonVersions.find(p => p.DisplayName === selectedProton);
    let cleanName = selectedProton;
    if (cleanName.startsWith("(Steam) ")) {
      cleanName = cleanName.substring(8);
    }
    const opts: launcher.LaunchOptions = {
      GamePath: gamePath,
      PrefixPath: prefixPath,
      ProtonPattern: cleanName,
      ProtonPath: tool ? tool.Path : "",
      EnableGamescope: enableGamescope && systemStatus.hasGamescope,
      GamescopeW: gamescopeW,
      GamescopeH: gamescopeH,
      GamescopeR: gamescopeR,
      EnableMangoHud: enableMango && systemStatus.hasMangoHud,
      EnableGamemode: enableGamemode && systemStatus.hasGameMode,
      EnableLsfgVk: enableLsfg,
      LsfgMultiplier: lsfgMultiplier,
      LsfgPerfMode: lsfgPerfMode,
      LsfgDllPath: lsfgDllPath,
    };
    try {
      await RunGame(opts, showLogsWindow);
      WindowHide();
    } catch (err) {
      console.error("Launch failed:", err);
      notifications.add(`Launch failed: ${err}`, "error");
    }
  }
</script>

<div class="run-container">
  <div class="header-row">
    <h1 class="page-title">Launch Configuration</h1>
  </div>

  <div class="form-container">
    <div class="form-group">
      <label>Game Executable</label>
      <div class="input-group">
        <input type="text" bind:value={gamePath} placeholder="Select .exe file..." class="input" />
        <button on:click={handleBrowseGame} class="btn">Browse</button>
      </div>
    </div>

    <div class="form-group">
      <label>WINEPREFIX</label>
      <div class="input-group">
        <div class="dropdown-wrapper">
          <Dropdown options={[...availablePrefixes, "Custom..."]} bind:value={selectedPrefixName} on:change={handlePrefixChange} />
        </div>
        <button on:click={handleBrowsePrefix} class="btn">Browse</button>
      </div>
      {#if selectedPrefixName === "Custom..." || !prefixPath.startsWith(baseDir)}
        <div class="path-display">{prefixPath}</div>
      {/if}
    </div>

    <div class="form-group">
      <label>Proton Version</label>
      <Dropdown options={protonOptions} bind:value={selectedProton} placeholder={isLoadingProton ? "Scanning..." : "Select Version"} disabled={isLoadingProton} on:change={handleProtonChange} />
    </div>

    <div class="toggles-grid">
      <SlideButton bind:checked={enableMango} label="MangoHud" subtitle="Performance overlay" />
      <SlideButton bind:checked={enableGamemode} label="GameMode" subtitle="Optimize priorities" />
      <SlideButton bind:checked={enableLsfg} label="LSFG-VK" subtitle="Lossless Scaling Frame Generation" hasConfig={true} on:config={() => showLsfgModal = true} />
      <SlideButton bind:checked={enableGamescope} label="Gamescope" subtitle="Micro-compositor" hasConfig={true} on:config={() => showGamescopeModal = true} />
      <SlideButton bind:checked={showLogsWindow} label="Show Logs" subtitle="Open logs in terminal" />
    </div>

    <!-- LSFG Settings Modal -->
    <Modal show={showLsfgModal} title="LSFG-VK Configuration" on:close={() => showLsfgModal = false}>
      <div class="modal-form">
        <div class="form-group">
          <label>
            Lossless.dll Path (from Steam)
            {#if lsfgDllPath}
              <span class="status-badge success">Detected</span>
            {:else}
              <span class="status-badge error">Not Found</span>
            {/if}
          </label>
          <div class="input-group">
            <input type="text" class="input sm" bind:value={lsfgDllPath} placeholder="Path to Lossless.dll..." />
            <button class="btn sm" on:click={handleBrowseDll}>Browse</button>
          </div>
        </div>
        <div class="form-group">
          <label>FPS Multiplier</label>
          <Dropdown options={["2", "3", "4"]} bind:value={lsfgMultiplier} />
        </div>
        <div class="form-group">
          <SlideButton bind:checked={lsfgPerfMode} label="Performance Mode" subtitle="Faster frame generation, slight quality loss" />
        </div>
      </div>
    </Modal>

    <!-- Gamescope Settings Modal -->
    <Modal show={showGamescopeModal} title="Gamescope Configuration" on:close={() => showGamescopeModal = false}>
      <div class="modal-form">
        <div class="form-group">
          <label>Width (px)</label>
          <input type="text" class="input" bind:value={gamescopeW} placeholder="e.g. 1920" />
        </div>
        <div class="form-group">
          <label>Height (px)</label>
          <input type="text" class="input" bind:value={gamescopeH} placeholder="e.g. 1080" />
        </div>
        <div class="form-group">
          <label>Refresh Rate (Hz)</label>
          <input type="text" class="input" bind:value={gamescopeR} placeholder="e.g. 60" />
        </div>
        <p class="note">Note: Mouse visibility fix enabled automatically.</p>
      </div>
    </Modal>

    <Modal show={showValidationModal} title="Missing Dependencies" on:close={() => showValidationModal = false}>
      <div class="warning-modal-content">
        <div class="warning-icon">⚠️</div>
        <p>The following requested features are not installed on your system:</p>
        <div class="missing-list">
          {#each missingToolsList as tool}
            <span class="tool-tag">{tool}</span>
          {/each}
        </div>
        <p class="question">Do you want to launch the game without these features?</p>
        <div class="modal-actions">
          <button class="btn secondary" on:click={() => showValidationModal = false}>Cancel</button>
          <button class="btn primary" on:click={proceedToLaunch}>Launch Anyway</button>
        </div>
      </div>
    </Modal>

    <div class="action-area">
      <button class="btn primary launch-btn" on:click={handleLaunch}>🚀 LAUNCH GAME</button>
    </div>
  </div>
</div>

<style lang="scss">
  .run-container { display: flex; flex-direction: column; height: 100%; padding: 32px; overflow: hidden; }
  .form-container { width: 100%; display: flex; flex-direction: column; gap: 24px; overflow-y: auto; padding-right: 8px; }
  .form-group label { display: block; font-size: 0.875rem; font-weight: 600; color: var(--text-muted); margin-bottom: 8px; }
  .input-group { display: flex; gap: 12px; width: 100%; .input { flex: 1; } .dropdown-wrapper { flex: 1; } }
  .path-display { margin-top: 8px; font-size: 0.75rem; color: var(--text-dim); word-break: break-all; padding: 8px; background: rgba(0, 0, 0, 0.2); border-radius: 6px; }
  .toggles-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 16px; margin-top: 8px; }
  .modal-form { display: flex; flex-direction: column; gap: 16px; }
  .warning-modal-content { text-align: center; .warning-icon { font-size: 3rem; margin-bottom: 16px; } p { color: var(--text-main); line-height: 1.5; }
    .missing-list { margin: 16px 0; display: flex; flex-wrap: wrap; justify-content: center; gap: 12px;
      .tool-tag { background: rgba(239, 68, 68, 0.1); color: #ef4444; padding: 6px 16px; border-radius: 20px; font-size: 0.9rem; border: 1px solid rgba(239, 68, 68, 0.2); font-weight: bold; }
    }
    .question { margin-top: 24px; font-weight: 600; color: var(--accent-secondary); }
  }
  .modal-actions { display: flex; gap: 12px; margin-top: 32px; button { flex: 1; padding: 12px; font-weight: 600; } }
  .action-area { margin-top: 16px; margin-bottom: 32px; }
  .launch-btn { width: 100%; padding: 16px; font-size: 1.125rem; font-weight: bold; border-radius: 12px; text-transform: uppercase; letter-spacing: 1px; }
  .input.sm { padding: 6px 10px; font-size: 0.875rem; }
  .btn.sm { padding: 6px 12px; font-size: 0.8rem; }
  
  .status-badge {
    font-size: 0.65rem;
    padding: 2px 6px;
    border-radius: 4px;
    margin-left: 8px;
    text-transform: uppercase;
    font-weight: bold;
    &.success { background: rgba(16, 185, 129, 0.2); color: #10b981; }
    &.error { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
  }

  .page-title { font-size: 2rem; font-weight: bold; color: var(--text-main); margin: 0 0 24px 0; }
  .note { font-size: 0.75rem; color: var(--text-muted); font-style: italic; margin-top: 8px; }
</style>
