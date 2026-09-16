package screenshot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Capture() (string, error) {
	dir := os.TempDir()
	filename := fmt.Sprintf("devctl-screenshot_%s.png", time.Now().Format("02-Jan-2006_15:04"))
	path := filepath.Join(dir, filename)

	cmd := exec.Command("screencapture", "-i", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("screencapture failed: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("screenshot cancelled by user")
	}

	return path, nil
}
