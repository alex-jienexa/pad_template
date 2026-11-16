package project

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Discover walks upward from start until it finds pad.yaml.
func Discover(start string) (string, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		abs = filepath.Dir(abs)
	}
	for {
		if _, err := os.Stat(filepath.Join(abs, "pad.yaml")); err == nil {
			return abs, nil
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return "", errors.New("could not find pad.yaml in this directory or its parents")
		}
		abs = parent
	}
}

// IsEmpty reports whether a directory has no entries.
func IsEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func SafeRelativePath(name string) (string, error) {
	// Template paths use portable slash-separated paths.
	// Reject OS-specific separators and rooted paths before conversion.
	if name == "" || strings.Contains(name, `\`) || strings.HasPrefix(name, "/") {
		return "", errors.New("unsafe template path")
	}

	// Reject Windows drive-prefixed paths on every operating system.
	if len(name) >= 2 &&
		((name[0] >= 'A' && name[0] <= 'Z') ||
		(name[0] >= 'a' && name[0] <= 'z')) &&
		name[1] == ':' {
			return "", errors.New("unsafe template path")
		}

		// Normalize using slash-based path semantics.
		clean := path.Clean(name)

		if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
			return "", errors.New("unsafe template path")
		}

		// Convert only after portable path validation.
		result := filepath.FromSlash(clean)

		if filepath.IsAbs(result) || filepath.VolumeName(result) != "" {
			return "", errors.New("unsafe template path")
		}

		return result, nil
}
