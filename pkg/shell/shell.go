package shell

import (
	"bufio"
	"fmt"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"github.com/fatih/color"
	"os"
)

func NewShell(sink chan error, f func(commandStr string) error) chan error {
	reader := bufio.NewReader(os.Stdin)
	for {
		cyan := color.New(color.FgCyan)
		boldCyan := cyan.Add(color.Bold)
		_, err := boldCyan.Print(fmt.Sprintf("➜ %s: ", util.AppName))
		if err != nil {
			sink <- err
			return sink
		}
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			sink <- err
			return sink
		}
		err = f(cmdString)
		if err != nil {
			sink <- err
			return sink
		}
	}
}
