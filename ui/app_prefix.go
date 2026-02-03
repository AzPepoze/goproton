package main

import (
	"goproton/pkg/launcher"
)

func (a *App) ListPrefixes() ([]string, error) {
	return launcher.ListPrefixes()
}

func (a *App) CreatePrefix(name string) error {
	return launcher.CreatePrefix(name)
}

func (a *App) GetPrefixBaseDir() string {
	return launcher.GetPrefixBaseDir()
}
