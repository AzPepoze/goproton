<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let checked = false;
  export let label = "";
  export let subtitle = "";
  export let hasConfig = false;

  const dispatch = createEventDispatcher();

  function toggle() {
    checked = !checked;
    dispatch('change', checked);
  }

  function openConfig(e: MouseEvent) {
    e.stopPropagation();
    dispatch('config');
  }
</script>

<div class="slide-button-card" class:active={checked} on:click={toggle}>
  <div class="info">
    <div class="title">{label}</div>
    {#if subtitle}
      <div class="subtitle">{subtitle}</div>
    {/if}
  </div>
  <div class="actions">
    {#if hasConfig}
      <button class="config-btn" on:click={openConfig} title="Configure">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
      </button>
    {/if}
    <div class="switch-container">
      <input type="checkbox" checked={checked} on:change|stopPropagation={toggle} />
      <span class="slider"></span>
    </div>
  </div>
</div>

<style lang="scss">
  .slide-button-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    padding: 16px 20px;
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid var(--glass-border);
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.2s;
    user-select: none;

    &:hover {
      background: rgba(255, 255, 255, 0.05);
      border-color: var(--glass-border-bright);

      .config-btn {
        opacity: 1;
      }
    }

    &.active {
      border-color: var(--accent-secondary);
      background: rgba(255, 255, 255, 0.04);
    }
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .config-btn {
    background: none;
    border: none;
    color: white;
    padding: 6px;
    border-radius: 6px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0.6;
    transition: all 0.2s;

    &:hover {
      background: rgba(255, 255, 255, 0.1);
      opacity: 1;
    }
  }

  .switch-container {
    position: relative;
    width: 38px;
    height: 20px;
    flex-shrink: 0;

    input {
      opacity: 0;
      width: 0;
      height: 0;

      &:checked + .slider {
        background-color: var(--accent-primary);
      }

      &:checked + .slider:before {
        transform: translateX(18px);
        background-color: #000;
      }
    }
  }

  .slider {
    position: absolute;
    cursor: pointer;
    top: 0; left: 0; right: 0; bottom: 0;
    background-color: rgba(255, 255, 255, 0.1);
    transition: .3s;
    border-radius: 34px;

    &:before {
      position: absolute;
      content: "";
      height: 14px;
      width: 14px;
      left: 3px;
      bottom: 3px;
      background-color: var(--text-dim);
      transition: .3s;
      border-radius: 50%;
    }
  }

  .info {
    .title {
      font-weight: 600;
      color: var(--text-main);
      font-size: 0.9rem;
    }
    .subtitle {
      font-size: 0.7rem;
      color: var(--text-dim);
      margin-top: 2px;
    }
  }
</style>