package shell

import (
	"bufio"
	"fmt"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"github.com/fatih/color"
	"os"
)

func NewShell(f func(commandStr string) error) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		cyan := color.New(color.FgCyan)
		boldCyan := cyan.Add(color.Bold)
		_, err := boldCyan.Print(fmt.Sprintf("➜ %s: ", util.AppName))
		if err != nil {
			return err
		}
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		err = f(cmdString)
		if err != nil {
			return err
		}
	}
}
