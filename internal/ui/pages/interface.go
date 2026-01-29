package pages

import (
	"fyne.io/fyne/v2"
)

// Page defines the contract for a UI page
type Page interface {
	// Content returns the canvas object to be displayed
	Content() fyne.CanvasObject
}
