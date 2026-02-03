package main

import (
	"goproton/pkg/launcher"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) ScanProtonVersions() ([]launcher.ProtonTool, error) {
	return launcher.GetProtonTools()
}

func (a *App) GetProtonVariants() []launcher.ProtonVariant {
	return launcher.GetKnownVariants()
}

func (a *App) GetProtonReleases(variantID string) ([]launcher.GitHubRelease, error) {
	return launcher.FetchReleases(variantID)
}

func (a *App) InstallProtonVersion(url, version string) error {
	return launcher.InstallProton(url, version, func(percent int, msg string) {
		runtime.EventsEmit(a.ctx, "install-proton-progress", map[string]interface{}{
			"percent": percent,
			"message": msg,
		})
	})
}
