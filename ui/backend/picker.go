package backend

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) runSystemPicker(title string, folder bool, filters []runtime.FileFilter) (string, bool) {
	if _, err := exec.LookPath("zenity"); err == nil {
		args := []string{"--file-selection", "--title=" + title}
		if folder {
			args = append(args, "--directory")
		}
		if len(filters) > 0 {
			for _, f := range filters {
				pattern := strings.ReplaceAll(f.Pattern, ";", " ")
				args = append(args, "--file-filter="+f.DisplayName+"|"+pattern)
			}
		}
		cmd := exec.Command("zenity", args...)
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output)), true
		}
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return "", true
		}
	}

	if _, err := exec.LookPath("kdialog"); err == nil {
		var args []string
		if folder {
			args = []string{"--getexistingdirectory", ".", "--title", title}
		} else {
			filterStr := ""
			if len(filters) > 0 {
				var parts []string
				for _, f := range filters {
					pattern := strings.ReplaceAll(f.Pattern, ";", " ")
					parts = append(parts, f.DisplayName+" ("+pattern+")")
				}
				filterStr = strings.Join(parts, ";;")
			}
			args = []string{"--getopenfilename", ".", filterStr, "--title", title}
		}
		cmd := exec.Command("kdialog", args...)
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output)), true
		}
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return "", true
		}
	}

	return "", false
}

func (a *App) SearchExecutables(folderPath string, maxDepth int, excludeNames []string) ([]string, error) {
	var executables []string
	cleanPath := filepath.Clean(folderPath)

	var excludeRegex []*regexp.Regexp
	for _, pattern := range excludeNames {
		if pattern == "" {
			continue
		}
		// Convert simple glob * to regex .* for user convenience
		regStr := pattern
		if strings.Contains(regStr, "*") {
			regStr = strings.ReplaceAll(regexp.QuoteMeta(regStr), "\\*", ".*")
		}
		// Case insensitive match
		re, err := regexp.Compile("(?i)" + regStr)
		if err == nil {
			excludeRegex = append(excludeRegex, re)
		}
	}

	err := filepath.WalkDir(cleanPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(cleanPath, path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			if rel == "." {
				return nil
			}
			depth := len(strings.Split(rel, string(os.PathSeparator)))
			if maxDepth != -1 && depth > maxDepth {
				return filepath.SkipDir
			}
			return nil
		}

		fileName := filepath.Base(path)
		for _, re := range excludeRegex {
			if re.MatchString(fileName) {
				return nil
			}
		}

		if strings.EqualFold(filepath.Ext(path), ".exe") {
			executables = append(executables, path)
		}

		return nil
	})

	return executables, err
}

func (a *App) PickFile() (string, error) {
	if path, ok := a.runSystemPicker("Select Game Executable", false, []runtime.FileFilter{
		{DisplayName: "Executables (*.exe)", Pattern: "*.exe"},
		{DisplayName: "All Files", Pattern: "*.*"},
	}); ok {
		return path, nil
	}

	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Game Executable",
		Filters: []runtime.FileFilter{
			{DisplayName: "Executables (*.exe)", Pattern: "*.exe"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
}

func (a *App) PickFolder() (string, error) {
	if path, ok := a.runSystemPicker("Select Directory", true, nil); ok {
		return path, nil
	}

	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Directory",
	})
}

func (a *App) PickFileCustom(title string, filters []runtime.FileFilter) (string, error) {
	if path, ok := a.runSystemPicker(title, false, filters); ok {
		return path, nil
	}

	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   title,
		Filters: filters,
	})
}
