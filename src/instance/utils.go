package main

import (
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandPath expands ~ to the user's home directory.
// This allows users to specify paths like ~/Games in configuration files.
func ExpandPath(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		u, err := user.Current()
		if err != nil {
			return path
		}
		if path == "~" {
			return u.HomeDir
		}
		return filepath.Join(u.HomeDir, path[2:])
	}
	return path
}

