package backend

import (
	"goproton/pkg/core"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) ScanProtonVersions() ([]core.ProtonTool, error) {
	return core.GetProtonTools()
}

func (a *App) GetProtonVariants() []core.ProtonVariant {
	return core.GetKnownVariants()
}

func (a *App) GetProtonReleases(variantID string) ([]core.GitHubRelease, error) {
	return core.FetchReleases(variantID)
}

func (a *App) InstallProtonVersion(url, version string) error {
	return core.InstallProton(url, version, func(percent int, msg string) {
		runtime.EventsEmit(a.ctx, "install-proton-progress", map[string]interface{}{
			"percent": percent,
			"message": msg,
		})
	})
}
