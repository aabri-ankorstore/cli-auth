package shell

import (
	"bufio"
	"fmt"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"github.com/fatih/color"
	"os"
)

func NewShell(sink chan error, f func(commandStr string) error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		cyan := color.New(color.FgCyan)
		boldCyan := cyan.Add(color.Bold)
		_, err := boldCyan.Print(fmt.Sprintf("➜ %s: ", util.AppName))
		if err != nil {
			sink <- err
		}
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			sink <- err
		}
		err = f(cmdString)
		if err != nil {
			sink <- err
		}
	}
}
