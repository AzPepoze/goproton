package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type RunPage struct {
	// Fields that need to be accessed externally (e.g. CLI args) can be public
	InitialExePath string
}

func NewRunPage() *RunPage {
	return &RunPage{}
}

func (p *RunPage) SetExePath(path string) {
	p.InitialExePath = path
}

func (p *RunPage) Content() fyne.CanvasObject {
	label := widget.NewLabel("Run Game - Logic Migration In Progress")
	return container.NewCenter(label)
}
