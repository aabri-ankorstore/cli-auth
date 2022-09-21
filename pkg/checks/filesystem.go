package checks

import (
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"os"
	"path/filepath"
)

const pattern = "*-auth.lock"

type FileSystem struct {
	Type       string
	PluginPath string
}

func NewFilesystem() CheckManager {
	return &FileSystem{
		Type: "filesystem",
	}
}

func (f *FileSystem) Listen() error {
	panic("Not Implemented")
}

func (f *FileSystem) CreateTmpFile() (string, error) {
	file, err := os.CreateTemp(f.PluginPath, pattern)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func (f *FileSystem) IsAuthenticatedOffline() bool {
	file := fmt.Sprintf("%s/%s", f.PluginPath, pattern)
	matches, err := filepath.Glob(file)
	if err != nil {
		return false
	}
	if len(matches) > 0 {
		return true
	}
	return false
}

func (f *FileSystem) RemoveAuth() {
	_ = os.Remove(utils.LockFile)
}

func (f *FileSystem) HandleClient() error {
	panic("Not Implemented")
}

func (f *FileSystem) CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
