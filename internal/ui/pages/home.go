package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type HomePage struct{}

func NewHomePage() *HomePage {
	return &HomePage{}
}

func (p *HomePage) Content() fyne.CanvasObject {
	// Hero Image / Icon
	heroIcon := widget.NewIcon(theme.HomeIcon())
	heroIcon.SetResource(theme.ComputerIcon()) // Placeholder
	
	title := widget.NewLabel("Welcome to GoProton")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Your minimalist Proton manager.")
	subtitle.Alignment = fyne.TextAlignCenter

	// Info Cards
	stats := container.NewGridWithColumns(2,
		p.createCard("Proton Versions", "Scanning...", theme.StorageIcon()),
		p.createCard("Active Prefix", "Default", theme.FolderIcon()),
	)

	return container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(heroIcon),
		title,
		subtitle,
		layout.NewSpacer(),
		stats,
		layout.NewSpacer(),
	)
}

func (p *HomePage) createCard(title, value string, icon fyne.Resource) fyne.CanvasObject {
	return widget.NewCard(title, value, nil)
}
