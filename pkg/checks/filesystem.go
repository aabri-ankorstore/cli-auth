package checks

import (
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"os"
)

const pattern = "*-auth.lock"

type FileSystem struct {
	PluginFolder string
}

func NewFilesystem(p string) *FileSystem {
	return &FileSystem{
		PluginFolder: p,
	}
}

func (f *FileSystem) CreateTmpFile() (string, error) {
	plugin := fmt.Sprintf("%s/%s", f.PluginFolder, utils.PluginPath)
	file, err := os.CreateTemp(plugin, pattern)
	f.CheckError(err)
	return file.Name(), nil
}

func (f *FileSystem) RemoveFile() {
	f.CheckError(os.Remove(utils.LockFile))
}

func (f *FileSystem) CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
