//go:build windows

package install

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func ensureWindowsPath(dir string) error {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return nil
	}
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return fmt.Errorf("open machine PATH (admin required): %w", err)
	}
	defer key.Close()
	current, _, err := key.GetStringValue("Path")
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("read machine PATH: %w", err)
	}
	for _, part := range strings.Split(current, ";") {
		if strings.EqualFold(strings.TrimSpace(part), dir) {
			return nil
		}
	}
	updated := strings.TrimSpace(current)
	if updated == "" {
		updated = dir
	} else {
		updated += ";" + dir
	}
	if err := key.SetStringValue("Path", updated); err != nil {
		return fmt.Errorf("update machine PATH: %w", err)
	}
	// Best-effort broadcast; ignore failure.
	_ = os.Getenv("Path")
	return nil
}
