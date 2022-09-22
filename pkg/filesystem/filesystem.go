package filesystem

import (
	"encoding/json"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"github.com/ankorstore/ankorstore-cli-core/pkg/plugin"
	"os"
)

const Pattern = "*-auth.lock"

type FileSystem struct {
	PluginFolder string
}

func NewFilesystem(p string) *FileSystem {
	return &FileSystem{
		PluginFolder: p,
	}
}

func (f *FileSystem) CreateTmpFile() (string, error) {
	p := fmt.Sprintf("%s/%s", f.PluginFolder, utils.PluginPath)
	file, err := os.CreateTemp(p, Pattern)
	if err != nil {
		return "", err
	}
	status := utils.AuthStatus{
		IsConnected: true,
	}
	statusByte, _ := json.Marshal(status)
	jsonContent := string(statusByte)
	_, err = file.WriteString(plugin.Encode(jsonContent))
	if err != nil {
		return "", err
	}
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
