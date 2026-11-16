package project

import (
	"errors"
	"os"
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

// SafeRelativePath rejects paths that could escape the destination root.
func SafeRelativePath(path string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(path))
	if clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errors.New("unsafe template path")
	}
	return clean, nil
}
