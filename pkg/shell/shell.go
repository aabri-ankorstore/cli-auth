package shell

import (
	"bufio"
	"github.com/fatih/color"
	"os"
)

func NewShell(f func(commandStr string) error) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		cyan := color.New(color.FgCyan)
		boldCyan := cyan.Add(color.Bold)
		_, err := boldCyan.Print("➜ ")
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
