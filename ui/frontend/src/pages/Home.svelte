<script lang="ts">
</script>

<div class="home-container">
  <div class="hero-section">
    <h1 class="welcome-text">WELCOME TO <span class="highlight">GOPROTON!</span></h1>
  </div>

  <div class="status-bar">
    <div class="status-badge">
      <span class="dot">●</span> System Ready
    </div>
  </div>

  <div class="grid">
    <div class="card hoverable">
      <div class="card-header">
        <span class="card-icon">📦</span>
        <h3>Proton Versions</h3>
      </div>
      <p>Scan and manage your installed compatibility tools from Steam and ProtonUp.</p>
    </div>

    <div class="card hoverable">
      <div class="card-header">
        <span class="card-icon">🍷</span>
        <h3>Active Prefix</h3>
      </div>
      <p>Current: <span class="badge">Default</span></p>
      <p class="mt-2 text-xs">Manage Wine configuration and prefixes.</p>
    </div>
  </div>
</div>

<style lang="scss">
  .home-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 64px 48px;
    overflow-y: auto;
    background-color: transparent;
  }

  .hero-section {
    text-align: center;
    margin-bottom: 48px;

    .welcome-text {
      font-size: 3.5rem;
      font-weight: 900;
      color: #fff;
      margin: 0;
      letter-spacing: -1px;
      line-height: 1.1;
      
      .highlight {
        color: transparent;
        -webkit-text-stroke: 1.5px rgba(255, 255, 255, 0.8);
        background: linear-gradient(180deg, #fff 0%, rgba(255,255,255,0.2) 100%);
        -webkit-background-clip: text;
      }
    }
  }

  .status-bar {
    display: flex;
    justify-content: center;
    margin-bottom: 48px;
  }

  .status-badge {
    background-color: rgba(255, 255, 255, 0.03);
    padding: 8px 20px;
    border-radius: 100px;
    border: 1px solid var(--glass-border);
    font-size: 0.8rem;
    font-weight: 700;
    color: var(--success);
    display: flex;
    align-items: center;
    gap: 10px;
    text-transform: uppercase;
    letter-spacing: 1px;

    .dot {
      font-size: 0.8rem;
    }
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 32px;
    max-width: 1000px;
    margin: 0 auto;
    width: 100%;
  }

  .card {
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid var(--glass-border);
    border-radius: 20px;
    padding: 32px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    .card-header {
      display: flex;
      align-items: center;
      gap: 16px;
      margin-bottom: 16px;
    }

    .card-icon {
      font-size: 1.5rem;
    }

    h3 {
      font-size: 1.25rem;
      font-weight: 700;
      margin: 0;
      color: #fff;
    }

    p {
      font-size: 0.95rem;
      color: var(--text-muted);
      margin: 0;
      line-height: 1.5;
    }

    &.hoverable {
      cursor: pointer;
      
      &:hover {
        background: rgba(255, 255, 255, 0.05);
        border-color: var(--glass-border-bright);
        transform: translateY(-4px);
      }
    }
  }

  .badge {
    background-color: rgba(255, 255, 255, 0.1);
    color: #fff;
    padding: 2px 8px;
    border-radius: 6px;
    font-family: monospace;
    font-weight: 600;
  }

  .mt-2 { margin-top: 8px; }
  .text-xs { font-size: 0.75rem; }
</style>
