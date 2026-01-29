package launcher

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
)

// ProtonTool represents a discovered compatibility tool
type ProtonTool struct {
	Name        string
	Path        string
	IsSteam     bool
	DisplayName string
}

// GetProtonTools scans for Proton/Compatibility tools in specified directories
func GetProtonTools() ([]ProtonTool, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	userSteamPath := filepath.Join(u.HomeDir, ".steam/root/compatibilitytools.d")
	systemSteamPath := "/usr/share/steam/compatibilitytools.d"

	var tools []ProtonTool

	// Scan User Steam directory
	userTools, _ := scanDir(userSteamPath, true)
	tools = append(tools, userTools...)

	// Scan System directory
	systemTools, _ := scanDir(systemSteamPath, false)
	tools = append(tools, systemTools...)

	// Sorting: (Steam) items first, then alphabetical
	sort.Slice(tools, func(i, j int) bool {
		if tools[i].IsSteam && !tools[j].IsSteam {
			return true
		}
		if !tools[i].IsSteam && tools[j].IsSteam {
			return false
		}
		return tools[i].Name < tools[j].Name
	})

	return tools, nil
}

func scanDir(path string, isSteam bool) ([]ProtonTool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var tools []ProtonTool
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			// Verification: Only include if it has a 'proton' executable
			protonPath := filepath.Join(path, name, "proton")
			if _, err := os.Stat(protonPath); os.IsNotExist(err) {
				continue // Not a real Proton tool
			}

			displayName := name
			if isSteam {
				displayName = fmt.Sprintf("(Steam) %s", name)
			}
			tools = append(tools, ProtonTool{
				Name:        name,
				Path:        filepath.Join(path, name),
				IsSteam:     isSteam,
				DisplayName: displayName,
			})
		}
	}
	return tools, nil
}

// ExpandPath expands ~ to the user's home directory
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		u, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(u.HomeDir, path[2:])
	}
	return path
}
