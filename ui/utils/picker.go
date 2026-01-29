package utils

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/ncruces/zenity"
)

// BrowseFile tries to open a native file picker, falling back to Fyne's dialog if necessary.
// extensions: list of extensions (e.g., ".exe", "*.png"). The function handles formatting differences.
func BrowseFile(w fyne.Window, title string, extensions []string, onSelected func(string)) {
	// 1. Try Native (Zenity)
	var zOptions []zenity.Option
	zOptions = append(zOptions, zenity.Title(title))

	if len(extensions) > 0 {
		var patterns []string
		for _, ext := range extensions {
			// Zenity expects "*.ext"
			clean := strings.TrimSpace(ext)
			if !strings.HasPrefix(clean, "*") {
				if strings.HasPrefix(clean, ".") {
					clean = "*" + clean // .exe -> *.exe
				} else {
					clean = "*." + clean // exe -> *.exe
				}
			}
			patterns = append(patterns, clean)
		}
		zOptions = append(zOptions, zenity.FileFilters{
			{Name: "Supported Files", Patterns: patterns},
		})
	}

	path, err := zenity.SelectFile(zOptions...)

	// Case A: Success Native
	if err == nil && path != "" {
		onSelected(path)
		return
	}

	// Case B: User Cancelled Native (Do NOT show fallback, just stop)
	if err == zenity.ErrCanceled {
		return
	}

	// Case C: Native Error/Not Supported -> Fallback to Fyne
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader != nil {
			onSelected(reader.URI().Path())
		}
	}, w)

	if len(extensions) > 0 {
		// Fyne expects ".ext"
		var fyneExts []string
		for _, ext := range extensions {
			clean := strings.TrimSpace(ext)
			clean = strings.TrimPrefix(clean, "*") // Remove *
			if !strings.HasPrefix(clean, ".") {
				clean = "." + clean // exe -> .exe
			}
			fyneExts = append(fyneExts, clean)
		}
		fd.SetFilter(storage.NewExtensionFileFilter(fyneExts))
	}

	fd.Show()
}

// BrowseFolder tries to open a native folder picker, falling back to Fyne's dialog.
func BrowseFolder(w fyne.Window, title string, onSelected func(string)) {
	// 1. Try Native
	path, err := zenity.SelectFile(
		zenity.Title(title),
		zenity.Directory(),
	)

	if err == nil && path != "" {
		onSelected(path)
		return
	}

	if err == zenity.ErrCanceled {
		return
	}

	// 2. Fallback
	fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err == nil && uri != nil {
			onSelected(uri.Path())
		}
	}, w)
	
	fd.Show()
}
