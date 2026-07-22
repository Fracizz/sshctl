package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// PrimaryConfigPath returns the canonical inventory path ~/.sshctl/servers.json.
func PrimaryConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".sshctl", "servers.json")
	}
	return filepath.Join(home, ".sshctl", "servers.json")
}

// DefaultConfigPath is an alias for PrimaryConfigPath.
func DefaultConfigPath() string {
	return PrimaryConfigPath()
}

func legacyConfigPaths(home string) []string {
	dirs := []string{".sshfrac", ".invossh"}
	out := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		out = append(out, filepath.Join(home, dir, "servers.json"))
	}
	return out
}

// MigrateLegacy copies the first existing legacy inventory into PrimaryConfigPath.
// The source file is renamed to servers.json.bak after a successful copy.
// Returns the legacy source path when migration ran.
func MigrateLegacy() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", nil
	}
	dest := PrimaryConfigPath()
	if _, err := os.Stat(dest); err == nil {
		return "", nil
	}
	for _, src := range legacyConfigPaths(home) {
		if _, err := os.Stat(src); err != nil {
			continue
		}
		loaded, err := Load(src)
		if err != nil {
			return "", fmt.Errorf("migrate %s: %w", src, err)
		}
		if err := Save(dest, loaded); err != nil {
			return "", fmt.Errorf("write %s: %w", dest, err)
		}
		backup := src + ".bak"
		if err := os.Rename(src, backup); err != nil {
			return "", fmt.Errorf("backup legacy config %s -> %s: %w", src, backup, err)
		}
		return src, nil
	}
	return "", nil
}
