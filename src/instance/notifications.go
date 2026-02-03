package main

import (
	"os/exec"
)

// sendNotification sends a desktop notification using notify-send
func sendNotification(title, message string) {
	_ = exec.Command("notify-send", "-a", "GoProton", title, message).Run()
}
