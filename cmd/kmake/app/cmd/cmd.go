package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewKMakeCommand(_ io.Reader, out, err io.Writer) *cobra.Command {
	c := &cobra.Command{
		Use:   "kmake",
		Short: "make kubernetes object",
	}

	c.AddCommand(NewCmdVersion(out))
	return c
}
