package util

import (
	"fmt"
	"os"
	"path"
)

type TempDir struct {
	path string
}

func NewTempDir(name string) (*TempDir, error) {
	path, err := os.MkdirTemp("", fmt.Sprintf("knit.%s", name))
	if err != nil {
		return nil, err
	}

	return &TempDir{path: path}, nil
}

func (d *TempDir) CreateFile(name string) (*os.File, error) {
	f, err := os.Create(path.Join(d.path, name))
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (d *TempDir) Remove() error {
	return os.Remove(d.path)
}
