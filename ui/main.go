package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"goproton-wails/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if len(os.Args) > 1 {
		gamePath := os.Args[1]

		if _, err := os.Stat(gamePath); err == nil {
			if absPath, err := filepath.Abs(gamePath); err == nil {
				os.Setenv("GOPROTON_LAUNCHER_PATH", absPath)
				fmt.Printf("Pre-selecting launcher path: %s\n", absPath)
			}
		}
	}

	app := backend.NewApp()

	err := wails.Run(&options.App{
		Title:  "GoProton",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 24, G: 24, B: 27, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
