<h1 align="center">
  <!-- <img src="docs" alt="Logo" width="128" height="128" style="filter: brightness(0) invert(1);"/><br> -->
  ✦ GoProton ✦
</h1>

<p align="center">
  <strong>◈ Proton Instance Manager for Linux ◈</strong>
  <br>
  <strong>◈ Powered by Go & umu-run ◈</strong>
</p>

<p align="center">
  <!-- <a href="https://github.com/AzPepoze/goproton/releases/latest">
    <img src="https://img.shields.io/github/v/release/AzPepoze/goproton?style=for-the-badge&label=%E2%97%88%20RELEASE%20%E2%97%88&labelColor=%23181818&color=%23ffffff" alt="Latest Release">
  </a> -->
  <a href="LICENSE">
    <img src="https://img.shields.io/github/license/AzPepoze/goproton?style=for-the-badge&label=%E2%97%88%20LICENSE%20%E2%97%88&labelColor=%23181818&color=%23ffffff" alt="License">
  </a>
  <a href="https://github.com/AzPepoze/goproton/stargazers">
    <img src="https://img.shields.io/github/stars/AzPepoze/goproton?style=for-the-badge&label=%E2%97%88%20STARS%20%E2%97%88&labelColor=%23181818&color=%23ffffff" alt="Stars">
  </a>
</p>

> [!WARNING]
> 
> This project is still in **early development**. You may encounter bugs or breaking changes. Feel free to report issues or contribute!

> [!NOTE]
>
> The name **GoProton** comes from the word **"Go"** (as in "Go for it!"), representing the ability to launch Proton immediately. It's also a happy coincidence that the project is built using the **Go (Golang)** programming language! XD

## CONTENTS

- [CONTENTS](#contents)
- [FEATURES](#features)
- [ARCHITECTURE \& EFFICIENCY](#architecture--efficiency)
- [PREREQUISITES](#prerequisites)
- [BUILD](#build)
- [USAGE](#usage)
- [STONKS!](#stonks)

## FEATURES

- **Independent Instance Manager** – Each game runs as a standalone detached process. Closing the main launcher does not affect running games.
- **Multi-Game Support** – Run multiple Windows applications simultaneously with unique Proton configurations and prefixes.
- **Process Isolation** – Every game gets its own System Tray icon for individual management (Graceful Stop/Status).
- **Native Terminal Integration** – Real-time logs are piped to your preferred terminal (Kitty, Alacritty, etc.) for live debugging.
- **Automatic Log Management** – Persistent logging to `~/GoProton/logs` with automatic rotation (keeps last 10 sessions)
- **umu-run Core** – Full compatibility with the Unified Linux Runtime (umu) for superior non-Steam game execution.

## ARCHITECTURE & EFFICIENCY

1. **UI (Wails/Go):** Used only for configuration. It **closes completely** after launching a game, freeing all WebKit memory.
2. **Instance Manager (`goproton-instance`):** A tiny Go binary that manages the process life-cycle and tray.

## PREREQUISITES

- **umu-launcher** (Required for execution)
- **Steam** (Installed and configured)
- **ProtonPlus** (Recommended for managing and adding Proton versions)

> [!TIP]
> 
> Use **ProtonPlus** or **Steam** to download and install different Proton versions. GoProton will automatically detect them in your Steam compatibility tools directory.

## BUILD

```bash
# Build the wrapper, instance manager, and UI
make build
```

## USAGE

```bash
# Open Launcher
./goproton

# CLI Direct Launch
./goproton-instance --game "path/to/game.exe" --prefix "path/to/prefix" --proton-pattern "GE-Proton"
```

## STONKS!

<div align="center">
  <a href="https://www.star-history.com/#AzPepoze/goproton&type=date&legend=top-left">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=AzPepoze/goproton&type=date&theme=dark&legend=top-left" />
      <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=AzPepoze/goproton&type=date&legend=top-left" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=AzPepoze/goproton&type=date&legend=top-left" width="600" />
    </picture>
  </a>
  <br>
  <br>
  <strong>✦ Made with ♥︎ by AzPepoze ✦</strong>
</div>
