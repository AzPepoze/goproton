package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CatppuccinMocha implements fyne.Theme
type CatppuccinMocha struct{}

var (
	// Catppuccin Mocha Colors
	colBase     = color.RGBA{R: 0x1e, G: 0x1e, B: 0x2e, A: 0xff} // Background
	colText     = color.RGBA{R: 0xcd, G: 0xd6, B: 0xf4, A: 0xff} // Foreground
	colMauve    = color.RGBA{R: 0xcb, G: 0xa6, B: 0xf7, A: 0xff} // Primary/Focus
	colSurface0 = color.RGBA{R: 0x31, G: 0x32, B: 0x44, A: 0xff} // Input BG
	colSurface1 = color.RGBA{R: 0x45, G: 0x47, B: 0x5a, A: 0xff} // Button
	colSurface2 = color.RGBA{R: 0x58, G: 0x5b, B: 0x70, A: 0xff} // Selection
	colOverlay0 = color.RGBA{R: 0x6c, G: 0x70, B: 0x86, A: 0xff} // Placeholder
	colRed      = color.RGBA{R: 0xf3, G: 0x8b, B: 0xa8, A: 0xff} // Error
	colGreen    = color.RGBA{R: 0xa6, G: 0xe3, B: 0xa1, A: 0xff} // Success
)

func (m CatppuccinMocha) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground, theme.ColorNameMenuBackground:
		return colBase
	case theme.ColorNameForeground:
		return colText
	case theme.ColorNamePlaceHolder:
		return colOverlay0
	case theme.ColorNameInputBackground:
		return colSurface0
	case theme.ColorNamePrimary, theme.ColorNameFocus:
		return colMauve
	case theme.ColorNameSelection:
		return colSurface2
	case theme.ColorNameButton, theme.ColorNameOverlayBackground:
		return colSurface1
	case theme.ColorNameError:
		return colRed
	case theme.ColorNameSuccess:
		return colGreen
	case theme.ColorNameScrollBar:
		return colSurface0
	}
	// Fallback to default dark theme for others
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (m CatppuccinMocha) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m CatppuccinMocha) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m CatppuccinMocha) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
