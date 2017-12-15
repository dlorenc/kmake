package app

import (
	"os"

	"github.com/r2d4/kmake/cmd/kmake/app/cmd"
)

func Run() error {
	c := cmd.NewKMakeCommand(os.Stdin, os.Stdout, os.Stderr)
	return c.Execute()
}
