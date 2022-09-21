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

func NewFilesystem() *FileSystem {
	return &FileSystem{
		Type: "filesystem",
	}
}

func (f *FileSystem) CreateTmpFile() (string, error) {
	file, err := os.CreateTemp(f.PluginPath, pattern)
	f.CheckError(err)
	return file.Name(), nil
}

func (f *FileSystem) IsAuthenticated() bool {
	file := fmt.Sprintf("%s/%s", f.PluginPath, pattern)
	matches, err := filepath.Glob(file)
	f.CheckError(err)
	if len(matches) > 0 {
		return true
	}
	return false
}

func (f *FileSystem) RemoveAuth() {
	f.CheckError(os.Remove(utils.LockFile))
}

func (f *FileSystem) CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
