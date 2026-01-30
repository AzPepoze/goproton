<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { fade, scale } from 'svelte/transition';

  export let show = false;
  export let title = "Settings";

  const dispatch = createEventDispatcher();

  function close() {
    dispatch('close');
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      close();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown}/>

{#if show}
  <div class="modal-backdrop" on:click={close} transition:fade={{ duration: 200 }}>
    <div class="modal-content glass" on:click|stopPropagation transition:scale={{ duration: 200, start: 0.95 }}>
      <div class="modal-header">
        <h3>{title}</h3>
        <button class="close-btn" on:click={close}>&times;</button>
      </div>
      <div class="modal-body">
        <slot></slot>
      </div>
      <div class="modal-footer">
        <button class="btn primary" on:click={close}>Done</button>
      </div>
    </div>
  </div>
{/if}

<style lang="scss">
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(2px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal-content {
    width: 90%;
    max-width: 500px;
    background: #1a1a1a;
    border: 1px solid #333;
    border-radius: 16px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.6);
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    padding: 20px 24px;
    border-bottom: 1px solid var(--glass-border);
    display: flex;
    align-items: center;
    justify-content: space-between;

    h3 {
      margin: 0;
      font-size: 1.25rem;
      color: var(--text-main);
    }

    .close-btn {
      background: none;
      border: none;
      color: var(--text-dim);
      font-size: 1.5rem;
      cursor: pointer;
      padding: 0;
      line-height: 1;
      
      &:hover {
        color: var(--text-main);
      }
    }
  }

  .modal-body {
    padding: 24px;
    max-height: 70vh;
    overflow-y: auto;
  }

  .modal-footer {
    padding: 16px 24px;
    border-top: 1px solid var(--glass-border);
    display: flex;
    justify-content: flex-end;
  }
</style>
