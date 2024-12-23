package util

import (
	"fmt"
	"os"
	"path/filepath"
)

type TempDir struct {
	Path string
}

func NewTempDir(name string) (*TempDir, error) {
	path, err := os.MkdirTemp("", fmt.Sprintf("knit.%s", name))
	if err != nil {
		return nil, err
	}

	return &TempDir{Path: path}, nil
}

func (d *TempDir) CreatePath(path string) (string, error) {
	newPath := filepath.Join(d.Path, path)
	err := os.MkdirAll(filepath.Dir(newPath), 0744)
	if err != nil {
		return "", err
	}
	return newPath, nil
}

func (d *TempDir) CreateFile(name string) (*os.File, error) {
	f, err := os.Create(filepath.Join(d.Path, name))
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (d *TempDir) Remove() error {
	return os.Remove(d.Path)
}

func FindFileUpward(filename string, errorNotFound bool) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		filePath := filepath.Join(dir, filename)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if errorNotFound {
		return "", fmt.Errorf("file %q not found in any parent directory", filename)
	} else {
		return "", nil
	}
}

func FindModuleRoot() (string, error) {
	modPath, err := FindFileUpward("kcl.mod", false)
	if err != nil {
		return "", err
	}

	return filepath.Dir(modPath), nil
}
