package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go-proton/internal/ui/pages"
)

// MainLayout manages the sidebar and content area
type MainLayout struct {
	ContentContainer *fyne.Container
	CurrentPage      pages.Page
	
	// Page Instances
	HomePage      *pages.HomePage
	RunPage       *pages.RunPage
	ProtonMgrPage *pages.ProtonManagerPage
	PrefixMgrPage *pages.PrefixManagerPage
}

func NewMainLayout() *MainLayout {
	l := &MainLayout{
		ContentContainer: container.NewMax(),
		HomePage:         pages.NewHomePage(),
		RunPage:          pages.NewRunPage(),
		ProtonMgrPage:    pages.NewProtonManagerPage(),
		PrefixMgrPage:    pages.NewPrefixManagerPage(),
	}
	
	// Default Page
	l.SetPage(l.HomePage)
	
	return l
}

func (l *MainLayout) SetPage(p pages.Page) {
	l.CurrentPage = p
	l.ContentContainer.Objects = []fyne.CanvasObject{p.Content()}
	l.ContentContainer.Refresh()
}

func (l *MainLayout) LoadUI() fyne.CanvasObject {
	// --- Sidebar ---
	// We use a List for selection, or just a VBox of Buttons if we want more control styling
	// For "Modern", a List with custom item rendering is good, but simple buttons are safer for now.
	// Let's use a nice vertical box with icons.

	navItems := []struct {
		Name string
		Icon fyne.Resource
		Page pages.Page
	}{
		{"Home", theme.HomeIcon(), l.HomePage},
		{"Run", theme.MediaPlayIcon(), l.RunPage},
		{"Versions", theme.StorageIcon(), l.ProtonMgrPage},
		{"Prefix", theme.FolderIcon(), l.PrefixMgrPage},
	}

	var buttons []fyne.CanvasObject
	for _, item := range navItems {
		item := item // Capture for closure
		btn := widget.NewButtonWithIcon(item.Name, item.Icon, func() {
			l.SetPage(item.Page)
		})
		btn.Alignment = widget.ButtonAlignLeading
		// btn.Importance = widget.LowImportance // Flat look
		buttons = append(buttons, btn)
	}

	sidebar := container.NewVBox(
		container.NewPadded(widget.NewLabelWithStyle("GoProton", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		container.NewVBox(buttons...),
		layout.NewSpacer(), // Push items up
		widget.NewLabelWithStyle("v1.0.0", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
	)

	// Wrap sidebar in a background if needed, but SplitContainer handles divider
	split := container.NewHSplit(
		container.NewPadded(sidebar), 
		container.NewPadded(l.ContentContainer),
	)
	split.Offset = 0.25 // 25% width for sidebar

	return split
}
