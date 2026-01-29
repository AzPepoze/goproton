package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PrefixManagerPage struct{}

func NewPrefixManagerPage() *PrefixManagerPage {
	return &PrefixManagerPage{}
}

func (p *PrefixManagerPage) Content() fyne.CanvasObject {
	label := widget.NewLabel("Prefix Manager - Coming Soon")
	return container.NewCenter(label)
}
