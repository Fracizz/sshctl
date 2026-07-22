package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const AppName = "sshctl"

// DefaultDir returns the recommended system install directory for this OS.
func DefaultDir() string {
	switch runtime.GOOS {
	case "windows":
		pf := os.Getenv("ProgramFiles")
		if pf == "" {
			pf = `C:\Program Files`
		}
		return filepath.Join(pf, AppName)
	default:
		return "/usr/local/bin"
	}
}

// TargetPath returns the installed binary path inside dir.
func TargetPath(dir string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(dir, AppName+".exe")
	}
	if filepath.Base(dir) == AppName || strings.HasSuffix(dir, "/bin") {
		return filepath.Join(dir, AppName)
	}
	return filepath.Join(dir, AppName)
}

// Install copies the current executable to dir and optionally adds dir to PATH.
func Install(dir string, addPath bool) (string, error) {
	src, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locate current binary: %w", err)
	}
	src, err = filepath.EvalSymlinks(src)
	if err != nil {
		return "", fmt.Errorf("resolve binary path: %w", err)
	}
	if dir == "" {
		dir = DefaultDir()
	}
	dest := TargetPath(dir)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", fmt.Errorf("create install dir: %w", err)
	}
	if err := copyFile(src, dest); err != nil {
		return "", err
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(dest, 0o755); err != nil {
			return "", fmt.Errorf("chmod %s: %w", dest, err)
		}
	}
	if addPath {
		if err := ensurePath(filepath.Dir(dest)); err != nil {
			return dest, err
		}
	}
	return dest, nil
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return fmt.Errorf("create %s: %w", dest, err)
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy binary: %w", err)
	}
	return out.Close()
}

func ensurePath(dir string) error {
	switch runtime.GOOS {
	case "windows":
		return ensureWindowsPath(dir)
	default:
		return fmt.Errorf("automatic PATH update on %s is not implemented; add %s to PATH manually", runtime.GOOS, dir)
	}
}
