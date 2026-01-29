package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ProtonManagerPage struct{}

func NewProtonManagerPage() *ProtonManagerPage {
	return &ProtonManagerPage{}
}

func (p *ProtonManagerPage) Content() fyne.CanvasObject {
	label := widget.NewLabel("Proton Version Manager - Coming Soon")
	return container.NewCenter(label)
}
