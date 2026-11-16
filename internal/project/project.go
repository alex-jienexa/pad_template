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
	if name == "" {
		return "", errors.New("unsafe template path")
	}

	// Normalize both Unix and Windows separators.
	// This allows native relative paths produced by filepath.Join.
	portable := strings.ReplaceAll(name, `\`, "/")

	// Reject rooted paths, including Windows paths rooted at the
	// current drive and UNC paths.
	if strings.HasPrefix(portable, "/") {
		return "", errors.New("unsafe template path")
	}

	// Reject Windows drive-prefixed paths on every OS.
	if len(portable) >= 2 &&
		((portable[0] >= 'A' && portable[0] <= 'Z') ||
			(portable[0] >= 'a' && portable[0] <= 'z')) &&
		portable[1] == ':' {
		return "", errors.New("unsafe template path")
	}

	// Normalize using portable slash-based path semantics.
	clean := path.Clean(portable)

	if clean == "." ||
		clean == ".." ||
		strings.HasPrefix(clean, "../") {
		return "", errors.New("unsafe template path")
	}

	result := filepath.FromSlash(clean)

	if filepath.IsAbs(result) || filepath.VolumeName(result) != "" {
		return "", errors.New("unsafe template path")
	}

	return result, nil
}
