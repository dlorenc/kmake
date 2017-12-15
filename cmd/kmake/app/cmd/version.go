package cmd

import (
	"fmt"
	"io"

	"github.com/r2d4/kmake/pkg/kmake/version"
	"github.com/spf13/cobra"
)

func NewCmdVersion(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the version of kmake",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVersion(out, cmd)
		},
	}
	return cmd
}

func RunVersion(out io.Writer, cmd *cobra.Command) error {
	fmt.Fprintln(out, version.GetVersion())
	return nil
}
